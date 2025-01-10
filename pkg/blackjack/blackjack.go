package blackjack

import (
	"fmt"
	"slices"
)

type PlayerId = uint

type BlackjackGame struct {
	Dealer      *Hand
	playerMap   map[PlayerId]*Player
	playerCount PlayerId
	currentTurn PlayerId
	shoe        Shoe
	GameState   GameState

	OnPlayerTurn func(PlayerId)
	OnGameUpdate func(*BlackjackGame)
}

func CreateGame() *BlackjackGame {
	return &BlackjackGame{
		Dealer:    CreateHand(0),
		playerMap: make(map[PlayerId]*Player, 0),
		shoe:      *CreateShoe(1),
	}
}

func (game *BlackjackGame) AddPlayerWithBalance(balance PlayerId) *Player {
	game.playerCount++
	newPlayer := CreatePlayer(game.playerCount, balance)
	game.playerMap[game.playerCount] = newPlayer

	game.sendGameUpdate()

	return newPlayer
}

func (game *BlackjackGame) SetPlayerBet(playerId PlayerId, betAmount uint) error {
	if game.GameState != BettingState {
		return WrongGameStateError
	}

	player, ok := game.GetPlayer(playerId)
	if !ok {
		return PlayerNotFoundError
	}

	err := player.PlaceBet(betAmount)
	if err != nil {
		return err
	}

	game.sendGameUpdate()

	return nil
}

func (game *BlackjackGame) SkipPlayerBet(playerId PlayerId) error {
	if game.GameState != BettingState {
		return WrongGameStateError
	}

	player, ok := game.GetPlayer(playerId)
	if !ok {
		return PlayerNotFoundError
	}

	player.playing = true

	game.sendGameUpdate()

	return nil
}

func (game *BlackjackGame) sendGameUpdate() {
	if game.OnGameUpdate != nil {
		game.OnGameUpdate(game)
	}

	if game.GameState == PlayingState {
		dealer, nextNum := game.nextPlayersTurn()
		if dealer {
			game.GameState = DealerState
			go game.DealerTurn()
		} else if nextNum != game.currentTurn {
			game.sendPlayerTurn(nextNum)
		}
	}
}

func (game *BlackjackGame) sendPlayerTurn(playerId PlayerId) {
	if game.OnPlayerTurn != nil {
		go game.OnPlayerTurn(playerId)
	}
}

var PlayerNotFoundError error = fmt.Errorf("Could not find player")

func (game *BlackjackGame) RemovePlayer(playerNum PlayerId) (PlayerId, error) {
	playerToDelete, ok := game.playerMap[playerNum]
	if !ok {
		return 0, PlayerNotFoundError
	}

	delete(game.playerMap, playerNum)
	playersBalance := playerToDelete.Destroy()

	game.sendGameUpdate()

	return playersBalance, nil
}

func (game *BlackjackGame) GetPlayer(playerNum PlayerId) (*Player, bool) {
	player, ok := game.playerMap[playerNum]

	return player, ok
}

// Get list of all people that still need to bet
func (b *BlackjackGame) GetPlayersWihoutBets() []PlayerId {
	peopleArr := make([]PlayerId, 0)
	for k, v := range b.playerMap {
		if v.Hand == nil && v.playing == false {
			peopleArr = append(peopleArr, k)
		}
	}

	return peopleArr
}

// Returns true if dealers turn is next
// Returns false, playerId of the player that is next
func (b *BlackjackGame) nextPlayersTurn() (isDealersTurn bool, turnPlayerId PlayerId) {
	players := make([]PlayerId, 0, len(b.playerMap))
	for k, player := range b.playerMap {
		if player.Hand == nil {
			continue
		}
		if player.Hand.IsLocked() {
			continue
		}
		players = append(players, k)
	}

	slices.Sort(players)

	if len(players) > 0 {
		return false, players[0]
	}

	return true, PlayerId(0)
}

func (game *BlackjackGame) reset() {
	// TODO: Maybe something with the payoutMap from this function
	game.payoutBets()

	game.Dealer = nil
	game.GameState = NoState
	game.sendGameUpdate()
}

func (game *BlackjackGame) payoutBets() map[PlayerId]uint {
	if game.GameState != PayoutState || game.Dealer == nil || !game.Dealer.IsLocked() {
		return nil
	}

	dealerTotal := game.Dealer.Total()
	dealerBust := dealerTotal > 21
	payoutMap := make(map[PlayerId]uint)
	for playerId, player := range game.playerMap {
		playerTotal := player.Hand.Total()
		defer player.reset()
		if player.Hand == nil || !player.Hand.IsLocked() {
			continue
		}

		winnings := uint(0)
		if playerBust := playerTotal > 21; playerBust {
			winnings = 0
		} else if isDraw := !dealerBust && playerTotal == dealerTotal; isDraw {
			winnings = player.Hand.Bet
		} else if isWin := dealerBust || playerTotal > dealerTotal; isWin {
			winnings = 2 * player.Hand.Bet
		}

		player.Balance += winnings
		payoutMap[playerId] += winnings
	}

	return payoutMap
}

func (game *BlackjackGame) GetPlayerCount() uint {
	return game.playerCount
}
