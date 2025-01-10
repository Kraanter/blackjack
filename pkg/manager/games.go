package manager

import (
	"fmt"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

// Returns ID 0 if no game is found and game pointer will be nil
func (m *Manager) GetJoinableGame() (GameId, *blackjack.BlackjackGame) {
	// Join a random game
	for k, v := range m.gameMap {
		if v.GetPlayerCount() < m.Settings.MinPlayerCount {
			game, err := m.GetGameWithId(k)
			if err != nil {
				return 0, nil
			}

			return k, game
		}
	}

	return m.createNewGame()
}

var GameNotFoundError = fmt.Errorf("Could not find Game")

func (m *Manager) GetGameWithId(id GameId) (*blackjack.BlackjackGame, error) {
	game, ok := m.gameMap[id]

	if !ok {
		return nil, GameNotFoundError
	}

	return game, nil

}

func (m *Manager) createNewGame() (GameId, *blackjack.BlackjackGame) {
	newGame := blackjack.CreateGame()
	for {
		gameId := CreateRandomGameId(m.Settings.IdLength)

		if m.addGameWithID(gameId, newGame) {
			return gameId, newGame
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
