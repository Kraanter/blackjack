package manager

import (
	"context"
	"sync"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

type ManagedPlayer struct {
	mu      sync.RWMutex
	leaving bool
	Player  *blackjack.Player `json:"player"`
	GameId  GameId            `json:"game-id"`

	Game *blackjack.BlackjackGame `json:"game"`

	onGameUpdate func(*blackjack.GameSnapshot)
}

func (p *ManagedPlayer) SetGameUpdateHandler(handler func(*blackjack.GameSnapshot)) {
	p.mu.Lock()
	p.onGameUpdate = handler
	p.mu.Unlock()
}

func (p *ManagedPlayer) gameUpdateHandler() func(*blackjack.GameSnapshot) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.onGameUpdate
}

func createPlayer(game *ManagedGame, gameId GameId, player *blackjack.Player) *ManagedPlayer {
	return &ManagedPlayer{
		Player: player,
		Game:   game.blackjackGame,
		GameId: gameId,
	}
}

// Returns true if player got a new card
// Returns false if game is not in state for player to receive card
// Returns error if any other reason like player could not be found
func (p *ManagedPlayer) Hit() (bool, error) {
	p.mu.RLock()
	game, player := p.Game, p.Player
	p.mu.RUnlock()
	if game == nil || player == nil {
		return false, blackjack.PlayerNotFoundError
	}
	return game.PlayerHit(player.PlayerNum)
}

func (p *ManagedPlayer) Stand() error {
	p.mu.RLock()
	game, player := p.Game, p.Player
	p.mu.RUnlock()
	if game == nil || player == nil {
		return blackjack.PlayerNotFoundError
	}
	return game.PlayerStand(player.PlayerNum)
}

func (p *ManagedPlayer) GetBalance() uint {
	p.mu.RLock()
	game, player := p.Game, p.Player
	p.mu.RUnlock()
	if game == nil || player == nil {
		return 0
	}
	snapshotPlayer := game.Snapshot().PlayerMap[player.PlayerNum]
	if snapshotPlayer == nil {
		return 0
	}
	return snapshotPlayer.Balance
}

func (p *ManagedPlayer) Split() error {
	p.mu.RLock()
	game, player := p.Game, p.Player
	p.mu.RUnlock()
	if game == nil || player == nil {
		return blackjack.PlayerNotFoundError
	}
	return game.PlayerSplit(player.PlayerNum)
}

func (p *ManagedPlayer) Bet(amount uint) error {
	p.mu.RLock()
	game, player := p.Game, p.Player
	p.mu.RUnlock()
	if game == nil || player == nil {
		return blackjack.PlayerNotFoundError
	}
	err := game.SetPlayerBet(player.PlayerNum, amount)
	if err == nil {
		go game.Start(context.Background())
	}
	return err
}

func (p *ManagedPlayer) SkipBet() error {
	p.mu.RLock()
	game, player := p.Game, p.Player
	p.mu.RUnlock()
	if game == nil || player == nil {
		return blackjack.PlayerNotFoundError
	}
	err := game.SkipPlayerBet(player.PlayerNum)
	if err == nil {
		go game.Start(context.Background())
	}
	return err
}

func (p *ManagedPlayer) Leave() (balance uint, err error) {
	p.mu.Lock()
	if p.Game == nil || p.Player == nil || p.leaving {
		p.mu.Unlock()
		return 0, blackjack.PlayerNotFoundError
	}
	game, player := p.Game, p.Player
	p.leaving = true
	p.mu.Unlock()

	balance, err = game.RemovePlayer(player.PlayerNum)

	p.mu.Lock()
	p.Game = nil
	p.GameId = GameId(0)
	p.Player = nil
	p.leaving = false
	p.mu.Unlock()
	return balance, err
}

func (p *Manager) JoinGame(balance uint, gameId GameId) *ManagedPlayer {
	game, _ := p.GetGameWithId(gameId)
	if game == nil {
		return nil
	}

	player := game.blackjackGame.AddPlayerWithBalance(balance)

	manPlayer := createPlayer(game, gameId, player)

	game.mu.Lock()
	game.Players[manPlayer.Player.PlayerNum] = manPlayer
	game.mu.Unlock()

	return manPlayer
}

func (p *ManagedPlayer) String() string {
	p.mu.RLock()
	game, player := p.Game, p.Player
	p.mu.RUnlock()
	if game == nil || player == nil {
		return ""
	}
	snapshotPlayer := game.Snapshot().PlayerMap[player.PlayerNum]
	if snapshotPlayer == nil {
		return ""
	}
	return snapshotPlayer.String()
}

func (p *ManagedPlayer) Snapshot() *blackjack.GameSnapshot {
	p.mu.RLock()
	game := p.Game
	p.mu.RUnlock()
	if game == nil {
		return nil
	}
	return game.Snapshot()
}
