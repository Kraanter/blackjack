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
	Dealer      *util.Cell[*Hand]     `json:"dealer"`
	PlayerMap   map[PlayerId]*Player  `json:"players"`
	GameState   *util.Cell[GameState] `json:"gameState"`
	CurrentTurn *util.Cell[PlayerId]  `json:"current-turn"`

	hiddenDealerCard *util.Cell[*Card]
	playerCount      uint
	shoe             Shoe
	mu               sync.Mutex

	bettingStateChannel chan (struct{})

	OnPlayerTurn   func(PlayerId)          `json:"-"`
	OnGameUpdate   func(*GameSnapshot)     `json:"-"`
	OnGameFinished func(map[PlayerId]uint) `json:"-"`
}

func CreateGame() *BlackjackGame {
	var hiddenCard *Card = nil
	game := BlackjackGame{
		Dealer:              util.New(CreateHand(0)),
		PlayerMap:           make(map[PlayerId]*Player, 0),
		shoe:                *CreateShoe(1),
		bettingStateChannel: make(chan struct{}),
		CurrentTurn:         util.New(PlayerId(0)),
		GameState:           util.New(NoState),
		hiddenDealerCard:    util.New(hiddenCard),
	}

	game.Initialize()

	return &game
}

// GameSnapshot is a copy of game state that can safely outlive a game update.
type GameSnapshot struct {
	Dealer      *Hand                `json:"dealer"`
	PlayerMap   map[PlayerId]*Player `json:"players"`
	GameState   GameState            `json:"gameState"`
	CurrentTurn PlayerId             `json:"current-turn"`
}

func (game *BlackjackGame) Snapshot() *GameSnapshot {
	game.mu.Lock()
	defer game.mu.Unlock()
	return game.snapshot()
}

func (game *BlackjackGame) snapshot() *GameSnapshot {
	snapshot := &GameSnapshot{
		Dealer:      cloneHand(game.Dealer.Get()),
		PlayerMap:   make(map[PlayerId]*Player, len(game.PlayerMap)),
		GameState:   game.GameState.Get(),
		CurrentTurn: game.CurrentTurn.Get(),
	}
	for id, player := range game.PlayerMap {
		copy := *player
		copy.Hands = make([]*Hand, len(player.Hands))
		for i, hand := range player.Hands {
			copy.Hands[i] = cloneHand(hand)
		}
		snapshot.PlayerMap[id] = &copy
	}
	return snapshot
}

func cloneHand(hand *Hand) *Hand {
	if hand == nil {
		return nil
	}
	copy := *hand
	copy.Cards = append([]*Card(nil), hand.Cards...)
	return &copy
}

func (game *BlackjackGame) createGameStateError(expected GameState) error {
	return createGameStateError(expected, game.GameState.Get())
}

func (game *BlackjackGame) AddPlayerWithBalance(balance PlayerId) *Player {
	game.mu.Lock()
	defer game.mu.Unlock()
	game.playerCount++
	newPlayer := CreatePlayer(game.playerCount, balance)

	game.PlayerMap[game.playerCount] = newPlayer

	game.sendGameUpdate()

	return newPlayer
}

func (game *BlackjackGame) SetPlayerBet(playerId PlayerId, betAmount uint) error {
	game.mu.Lock()
	defer game.mu.Unlock()
	if game.GameState.Get() != BettingState {
		return game.createGameStateError(BettingState)
	}

	player, ok := game.getPlayer(playerId)
	if !ok {
		return PlayerNotFoundError
	}

	err := player.placeBet(betAmount)
	if err != nil {
		return err
	}

	game.gameloopTick()

	return nil
}

func (game *BlackjackGame) SkipPlayerBet(playerId PlayerId) error {
	game.mu.Lock()
	defer game.mu.Unlock()
	if game.GameState.Get() != BettingState {
		return game.createGameStateError(BettingState)
	}

	player, ok := game.getPlayer(playerId)
	if !ok {
		return PlayerNotFoundError
	}

	player.playing = true

	game.gameloopTick()

	return nil
}

// sendGameUpdate must be called while game.mu is held.
func (game *BlackjackGame) sendGameUpdate() {
	if game.OnGameUpdate != nil {
		game.OnGameUpdate(game.snapshot())
	}
}

func (game *BlackjackGame) gameloopTick() {
	game.sendGameUpdate()

	switch game.GameState.Get() {
	case PlayingState:
		dealer, nextNum := game.nextPlayersTurn()
		if dealer {
			game.GameState.Set(DealerState)
			err := game.dealerTurn()
			if err != nil {
				panic(fmt.Sprintf("Error during dealer turn: %v", err))
			}
		} else if nextNum != game.CurrentTurn.Get() {
			game.CurrentTurn.Set(nextNum)
			game.sendPlayerTurn(nextNum)
		}
	case BettingState:
		if len(game.PlayerMap) > 0 && len(game.GetPlayersWihoutBets()) == 0 {
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

func (game *BlackjackGame) RemovePlayer(playerNum PlayerId) (uint, error) {
	game.mu.Lock()
	defer game.mu.Unlock()
	balance, err := func() (uint, error) {
		playerToDelete, ok := game.PlayerMap[playerNum]
		if !ok {
			return 0, PlayerNotFoundError
		}

		delete(game.PlayerMap, playerNum)
		return playerToDelete.Destroy(), nil
	}()
	if err != nil {
		return 0, err
	}

	game.gameloopTick()

	return balance, nil
}

func (game *BlackjackGame) GetPlayer(playerNum PlayerId) (*Player, bool) {
	game.mu.Lock()
	defer game.mu.Unlock()
	return game.getPlayer(playerNum)
}

func (game *BlackjackGame) getPlayer(playerNum PlayerId) (*Player, bool) {
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
	if b.playerCount == 0 {
		return false, PlayerId(0)
	}

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

	game.Dealer.Set(nil)
	game.CurrentTurn.Set(0)
	game.GameState.Set(NoState)

	game.gameloopTick()

	if game.OnGameFinished != nil {
		go game.OnGameFinished(payoutMap)
	}

}

func (game *BlackjackGame) payoutBets() map[PlayerId]uint {
	if game.GameState.Get() != PayoutState || game.Dealer.Get() == nil || !game.Dealer.Get().IsLocked() {
		return nil
	}

	dealerTotal := game.Dealer.Get().Total()
	dealerBust := dealerTotal > 21
	dealerBlackjack := isBlackjack(game.Dealer.Get())
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
	game.mu.Lock()
	defer game.mu.Unlock()
	return uint(len(game.PlayerMap))
}

func (game *BlackjackGame) String() string {
	game.mu.Lock()
	defer game.mu.Unlock()
	return game.string()
}

func (game *BlackjackGame) string() string {
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

	return fmt.Sprintf("GameState: %v Playercount: %v NextPlayer: %v\nDealer: %v hidden: %v\nHands:\n%v", game.GameState.Get(), len(game.PlayerMap), nextString, game.Dealer.Get().String(), game.hiddenDealerCard.Get(), playerStrings)
}

func (b *BlackjackGame) forEachPlayer(fn func(id PlayerId, player *Player)) {
	for k, v := range b.PlayerMap {
		fn(k, v)
	}
}
