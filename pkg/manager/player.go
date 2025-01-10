package manager

import (
	"context"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

type ManagedPlayer struct {
	Player *blackjack.Player
	GameId GameId
	Ctx    context.Context

	Game *blackjack.BlackjackGame
}

func createPlayer(ctx context.Context, game *blackjack.BlackjackGame, gameId GameId, player *blackjack.Player) *ManagedPlayer {
	return &ManagedPlayer{
		Player: player,
		Game:   game,
		GameId: gameId,
		Ctx:    ctx,
	}
}

// Returns true if player got a new card
// Returns false if game is not in state for player to receive card
// Returns error if any other reason like player could not be found
func (p *ManagedPlayer) Hit() (bool, error) {
	return p.Game.PlayerHit(p.Player.PlayerNum)
}

func (p *ManagedPlayer) Stand() error {
	return p.Game.PlayerStand(p.Player.PlayerNum)
}

func (p *ManagedPlayer) Bet(amount uint) error {
	return p.Game.SetPlayerBet(p.Player.PlayerNum, amount)
}

func (p *ManagedPlayer) SkipBet() error {
	return p.Game.SkipPlayerBet(p.Player.PlayerNum)
}

func (p *ManagedPlayer) Leave() (balance uint, err error) {
	return p.Game.RemovePlayer(p.Player.PlayerNum)
}

func (p *Manager) JoinGame(ctx context.Context, balance uint, gameId GameId) *ManagedPlayer {
	game, _ := p.GetGameWithId(gameId)
	if game == nil {
		return nil
	}

	player := game.AddPlayerWithBalance(balance)

	return createPlayer(ctx, game, gameId, player)
}

func (p *ManagedPlayer) String() string {
	return p.Player.String()
}
