package manager

import (
	"github.com/kraanter/blackjack/pkg/blackjack"
)

type ManagedGame struct {
	blackjackGame *blackjack.BlackjackGame
	// TODO: Maybe player names
	Players map[blackjack.PlayerId]*ManagedPlayer
}

func createManagedGame() *ManagedGame {
	blackjackGame := blackjack.CreateGame()

	manGame := &ManagedGame{
		blackjackGame: blackjackGame,
		Players:       make(map[blackjack.PlayerId]*ManagedPlayer),
	}

	blackjackGame.OnGameUpdate = createGameUpdateHandler(manGame)

	return manGame
}

func (m *ManagedGame) GetPlayerCount() uint {
	return m.blackjackGame.GetPlayerCount()
}
