package blackjack

import (
	"fmt"
	"slices"
	"strconv"
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
	dealerBlackjack := isBlackjack(game.Dealer)
	payoutMap := make(map[PlayerId]uint)
	for playerId, player := range game.playerMap {
		playerTotal := player.Hand.Total()
		playerBust := playerTotal > 21
		playerBlackjack := isBlackjack(player.Hand)
		defer player.reset()
		if player.Hand == nil || !player.Hand.IsLocked() {
			continue
		}

		winnings := uint(0)
		switch {
		case dealerBlackjack && playerBlackjack, playerTotal == dealerTotal:
			winnings = player.Hand.Bet
		case playerBlackjack && !dealerBlackjack:
			// Blackjack pays 2 to 3
			winnings = 5 * (player.Hand.Bet / 2)
		case dealerBust && !playerBust, (!playerBust && playerTotal > dealerTotal):
			winnings = 2 * player.Hand.Bet
		case playerBust:
			winnings = 0
		}

		player.Balance += winnings
		payoutMap[playerId] += winnings
	}

	return payoutMap
}

func isBlackjack(hand *Hand) bool {
	if hand == nil {
		return false
	}
	return hand.Total() == 21 && len(hand.Cards) == 2
}

func (game *BlackjackGame) GetPlayerCount() uint {
	return uint(len(game.playerMap))
}

func (game *BlackjackGame) String() string {
	dealerTurn, nextPlayer := game.nextPlayersTurn()
	var nextString string
	if dealerTurn {
		nextString = "dealer"
	} else {
		nextString = strconv.Itoa(int(nextPlayer))
	}

	playerStrings := ""
	for _, v := range game.playerMap {
		playerStrings += "  " + v.String() + "\n"
	}

	return fmt.Sprintf("GameState: %v Playercount: %v NextPlayer: %v\nDealer: %v\nHands:\n%v", game.GameState, game.GetPlayerCount(), nextString, game.Dealer.String(), playerStrings)
}
