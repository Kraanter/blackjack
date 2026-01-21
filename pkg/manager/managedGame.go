package manager

import (
	"context"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

type ManagedGame struct {
	blackjackGame *blackjack.BlackjackGame
	// TODO: Maybe player names
	Players     map[blackjack.PlayerId]*ManagedPlayer
	gameContext context.Context
	cancelGame  context.CancelFunc
}

func createManagedGame(ctx context.Context) *ManagedGame {
	blackjackGame := blackjack.CreateGame()
	gameContext, cancel := context.WithCancel(ctx)

	manGame := &ManagedGame{
		blackjackGame: blackjackGame,
		Players:       make(map[blackjack.PlayerId]*ManagedPlayer),
		gameContext:   gameContext,
		cancelGame:    cancel,
	}

	blackjackGame.OnGameUpdate = createGameUpdateHandler(manGame)
	blackjackGame.OnGameFinished = func(m map[blackjack.PlayerId]uint) {
		go blackjackGame.Start(manGame.gameContext)
	}

	go blackjackGame.Start(manGame.gameContext)

	return manGame
}

func (m *ManagedGame) GetPlayerCount() uint {
	return m.blackjackGame.GetPlayerCount()
}
