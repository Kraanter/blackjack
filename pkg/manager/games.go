package manager

import (
	"fmt"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

func (m *Manager) GetJoinableGame() *blackjack.BlackjackGame {
	// Join a random game
	for k, v := range m.gameMap {
		if v.GetPlayerCount() < m.settings.MinPlayerCount {
			game, err := m.GetGameWithId(&k)
			if err != nil {
				return nil
			}

			return game
		}
	}

	return m.createNewGame()
}

var GameNotFoundError = fmt.Errorf("Could not find Game")

func (m *Manager) GetGameWithId(id *GameId) (*blackjack.BlackjackGame, error) {
	if id == nil {
		return nil, GameNotFoundError
	}

	game, ok := m.gameMap[*id]

	if !ok {
		return nil, GameNotFoundError
	}

	return game, nil

}

func (m *Manager) createNewGame() *blackjack.BlackjackGame {
	newGame := blackjack.CreateGame(&m.settings.BlackjackSettings)
	for {
		gameId := CreateRandomGameId(m.settings.IdLength)

		if m.addGameWithID(gameId, newGame) {
			return newGame
		}
	}
}

func (m *Manager) addGameWithID(id GameId, game *blackjack.BlackjackGame) bool {
	_, ok := m.gameMap[id]
	if ok {
		return false
	}

	m.gameMap[id] = game

	return true
}
