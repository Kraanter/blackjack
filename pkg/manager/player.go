package manager

import (
	"github.com/kraanter/blackjack/pkg/blackjack"
)

type ManagedPlayer struct {
	Player *blackjack.Player `json:"player"`
	GameId GameId            `json:"game-id"`

	Game *blackjack.BlackjackGame `json:"game"`

	OnGameUpdate func(*blackjack.BlackjackGame) `json:"-"`
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
	return p.Game.PlayerHit(p.Player.PlayerNum)
}

func (p *ManagedPlayer) Stand() error {
	return p.Game.PlayerStand(p.Player.PlayerNum)
}

func (p *ManagedPlayer) GetBalance() uint {
	return p.Player.Balance
}

func (p *ManagedPlayer) Bet(amount uint) error {
	return p.Game.SetPlayerBet(p.Player.PlayerNum, amount)
}

func (p *ManagedPlayer) SkipBet() error {
	return p.Game.SkipPlayerBet(p.Player.PlayerNum)
}

func (p *ManagedPlayer) Leave() (balance uint, err error) {
	defer func() {
		p.Game = nil
		p.GameId = GameId(0)
		p.Player = nil
	}()

	return p.Game.RemovePlayer(p.Player.PlayerNum)
}

func (p *Manager) JoinGame(balance uint, gameId GameId) *ManagedPlayer {
	game, _ := p.GetGameWithId(gameId)
	if game == nil {
		return nil
	}

	player := game.blackjackGame.AddPlayerWithBalance(balance)

	manPlayer := createPlayer(game, gameId, player)

	game.Players[manPlayer.Player.PlayerNum] = manPlayer

	return manPlayer
}

func (p *ManagedPlayer) String() string {
	return p.Player.String()
}
