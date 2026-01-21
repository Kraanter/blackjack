package blackjack

import (
	"fmt"
	"slices"
	"strconv"
	"sync"

	"github.com/kraanter/blackjack/pkg/blackjack/util"
)

type PlayerId = uint

type BlackjackGame struct {
	Dealer      *Hand                `json:"dealer"`
	PlayerMap   map[PlayerId]*Player `json:"players"`
	GameState   util.Cell[GameState] `json:"gameState"`
	CurrentTurn util.Cell[PlayerId]  `json:"current-turn"`

	hiddenDealerCard util.Cell[*Card]
	playerCount      PlayerId
	shoe             Shoe
	playerMapMutex   sync.Mutex
	sendMutex        sync.Mutex

	bettingStateChannel chan (struct{})

	OnPlayerTurn   func(PlayerId)          `json:"-"`
	OnGameUpdate   func(*BlackjackGame)    `json:"-"`
	OnGameFinished func(map[PlayerId]uint) `json:"-"`
}

func CreateGame() *BlackjackGame {
	game := BlackjackGame{
		Dealer:              CreateHand(0),
		PlayerMap:           make(map[PlayerId]*Player, 0),
		shoe:                *CreateShoe(1),
		bettingStateChannel: make(chan struct{}),
	}

	game.CurrentTurn.After(game.sendPlayerTurn)
	game.GameState.After(func(_ GameState) { game.sendGameUpdate() })

	return &game
}

func (game *BlackjackGame) createGameStateError(expected GameState) error {
	return createGameStateError(expected, game.GameState.Get())
}

func (game *BlackjackGame) AddPlayerWithBalance(balance PlayerId) *Player {
	game.playerCount++
	newPlayer := CreatePlayer(game.playerCount, balance)

	game.playerMapMutex.Lock()
	game.PlayerMap[game.playerCount] = newPlayer
	game.playerMapMutex.Unlock()

	game.sendGameUpdate()

	return newPlayer
}

func (game *BlackjackGame) SetPlayerBet(playerId PlayerId, betAmount uint) error {
	if game.GameState.Get() != BettingState {
		return game.createGameStateError(BettingState)
	}

	player, ok := game.GetPlayer(playerId)
	if !ok {
		return PlayerNotFoundError
	}

	err := player.PlaceBet(betAmount)
	if err != nil {
		return err
	}

	game.gameloopTick()

	return nil
}

func (game *BlackjackGame) SkipPlayerBet(playerId PlayerId) error {
	if game.GameState.Get() != BettingState {
		return game.createGameStateError(BettingState)
	}

	player, ok := game.GetPlayer(playerId)
	if !ok {
		return PlayerNotFoundError
	}

	player.playing = true

	game.gameloopTick()

	return nil
}

func (game *BlackjackGame) sendGameUpdate() {
	if game.OnGameUpdate != nil {
		game.OnGameUpdate(game)
	}
}

func (game *BlackjackGame) gameloopTick() {
	game.sendGameUpdate()

	switch game.GameState.Get() {
	case PlayingState:
		dealer, nextNum := game.nextPlayersTurn()
		if dealer {
			game.GameState.Set(DealerState)
			err := game.DealerTurn()
			if err != nil {
				panic(fmt.Sprintf("Error during dealer turn: %v", err))
			}
		} else if nextNum != game.CurrentTurn.Get() {
			game.CurrentTurn.Set(nextNum)
		}
	case BettingState:
		if len(game.GetPlayersWihoutBets()) == 0 {
			select {
			case <-game.bettingStateChannel:
			default:
				close(game.bettingStateChannel)
			}
		}
	}
}

func (game *BlackjackGame) sendPlayerTurn(playerId PlayerId) {
	if game.OnPlayerTurn != nil {
		go game.OnPlayerTurn(playerId)
	}

	game.sendGameUpdate()
}

var PlayerNotFoundError error = fmt.Errorf("Could not find player")

func (game *BlackjackGame) RemovePlayer(playerNum PlayerId) (balance uint, err error) {
	{
		game.playerMapMutex.Lock()
		defer game.playerMapMutex.Unlock()
		playerToDelete, ok := game.PlayerMap[playerNum]
		if !ok {
			return 0, PlayerNotFoundError
		}

		delete(game.PlayerMap, playerNum)
		balance = playerToDelete.Destroy()
	}

	game.gameloopTick()

	return balance, nil
}

func (game *BlackjackGame) GetPlayer(playerNum PlayerId) (*Player, bool) {
	game.playerMapMutex.Lock()
	defer game.playerMapMutex.Unlock()
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
	payoutMap := game.payoutBets()

	game.Dealer = nil
	game.CurrentTurn.Set(0)
	game.GameState.Set(NoState)

	game.gameloopTick()

	if game.OnGameFinished != nil {
		game.OnGameFinished(payoutMap)
	}

}

func (game *BlackjackGame) payoutBets() map[PlayerId]uint {
	if game.GameState.Get() != PayoutState || game.Dealer == nil || !game.Dealer.IsLocked() {
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
			case dealerBlackjack && playerBlackjack, !playerBust && (playerTotal == dealerTotal):
				winnings = hand.Bet
			case playerBlackjack && !dealerBlackjack:
				// Blackjack pays 3:2
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

	return fmt.Sprintf("GameState: %v Playercount: %v NextPlayer: %v\nDealer: %v hidden: %v\nHands:\n%v", game.GameState.Get(), game.GetPlayerCount(), nextString, game.Dealer.String(), game.hiddenDealerCard.Get(), playerStrings)
}

func (b *BlackjackGame) forEachPlayer(fn func(id PlayerId, player *Player)) {
	b.playerMapMutex.Lock()
	defer b.playerMapMutex.Unlock()

	for k, v := range b.PlayerMap {
		fn(k, v)
	}
}
