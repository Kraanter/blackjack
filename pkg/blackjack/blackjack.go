package blackjack

import (
	"fmt"
	"slices"
	"strconv"
	"sync"
)

type PlayerId = uint

type BlackjackGame struct {
	Dealer      *Hand                `json:"dealer"`
	PlayerMap   map[PlayerId]*Player `json:"players"`
	GameState   GameState            `json:"gameState"`
	CurrentTurn PlayerId             `json:"current-turn"`

	playerCount PlayerId
	shoe        Shoe
	mutex       sync.Mutex

	OnPlayerTurn func(PlayerId)       `json:"-"`
	OnGameUpdate func(*BlackjackGame) `json:"-"`
}

func CreateGame() *BlackjackGame {
	return &BlackjackGame{
		Dealer:    CreateHand(0),
		PlayerMap: make(map[PlayerId]*Player, 0),
		shoe:      *CreateShoe(1),
	}
}

func (game *BlackjackGame) AddPlayerWithBalance(balance PlayerId) *Player {
	game.playerCount++
	newPlayer := CreatePlayer(game.playerCount, balance)

	game.mutex.Lock()
	game.PlayerMap[game.playerCount] = newPlayer
	game.mutex.Unlock()

	game.sendGameUpdate()

	return newPlayer
}

func (game *BlackjackGame) SetPlayerBet(playerId PlayerId, betAmount uint) error {
	if game.GameState == NoState {
		game.Start()
	}

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
		} else if nextNum != game.CurrentTurn {
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
	game.mutex.Lock()
	playerToDelete, ok := game.PlayerMap[playerNum]
	if !ok {
		return 0, PlayerNotFoundError
	}

	delete(game.PlayerMap, playerNum)
	game.mutex.Unlock()
	playersBalance := playerToDelete.Destroy()

	game.sendGameUpdate()

	return playersBalance, nil
}

func (game *BlackjackGame) GetPlayer(playerNum PlayerId) (*Player, bool) {
	game.mutex.Lock()
	defer game.mutex.Unlock()
	player, ok := game.PlayerMap[playerNum]

	return player, ok
}

// Get list of all people that still need to bet
func (b *BlackjackGame) GetPlayersWihoutBets() []PlayerId {
	peopleArr := make([]PlayerId, 0)
	b.forEachPlayer(func(k PlayerId, v *Player) {
		if len(v.Hands) == 0 && v.playing == false {
			peopleArr = append(peopleArr, k)
		}
	})

	return peopleArr
}

// Returns true if dealers turn is next
// Returns false, playerId of the player that is next
func (b *BlackjackGame) nextPlayersTurn() (isDealersTurn bool, turnPlayerId PlayerId) {
	players := make([]PlayerId, 0, len(b.PlayerMap))
	b.forEachPlayer(func(k PlayerId, player *Player) {
		if len(player.Hands) == 0 {
			return
		}
		if player.GetActiveHand() == nil {
			return
		}
		players = append(players, k)
	})

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
	game.forEachPlayer(func(playerId PlayerId, player *Player) {
		for _, hand := range player.Hands {
			playerTotal := hand.Total()
			playerBust := playerTotal > 21
			playerBlackjack := isBlackjack(hand)
			defer player.reset()
			if !hand.IsLocked() {
				continue
			}

			winnings := uint(0)
			switch {
			case dealerBlackjack && playerBlackjack, playerTotal == dealerTotal:
				winnings = hand.Bet
			case playerBlackjack && !dealerBlackjack:
				// Blackjack pays 2 to 3
				winnings = 5 * (hand.Bet / 2)
			case dealerBust && !playerBust, (!playerBust && playerTotal > dealerTotal):
				winnings = 2 * hand.Bet
			case playerBust:
				winnings = 0
			}
			player.Balance += winnings
			payoutMap[playerId] += winnings
		}
	})

	return payoutMap
}

func isBlackjack(hand *Hand) bool {
	if hand == nil {
		return false
	}
	return hand.Total() == 21 && len(hand.Cards) == 2
}

func (game *BlackjackGame) GetPlayerCount() uint {
	return uint(len(game.PlayerMap))
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
	game.forEachPlayer(func(k PlayerId, v *Player) {
		playerStrings += "  " + v.String() + "\n"
	})

	return fmt.Sprintf("GameState: %v Playercount: %v NextPlayer: %v\nDealer: %v\nHands:\n%v", game.GameState, game.GetPlayerCount(), nextString, game.Dealer.String(), playerStrings)
}

func (b *BlackjackGame) forEachPlayer(fn func(id PlayerId, player *Player)) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for k, v := range b.PlayerMap {
		fn(k, v)
	}
}
