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

func (m *Manager) JoinRandomGame(balance uint) *ManagedPlayer {
	gameId, _ := m.GetJoinableGame()
	if gameId == 0 {
		return nil
	}

	return m.JoinGame(balance, gameId)
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

func (m *Manager) RemoveGame(id GameId) bool {
	game, err := m.GetGameWithId(id)
	if err != nil {
		return false
	}

	if game.GetPlayerCount() != 0 {
		return false
	}

	delete(m.gameMap, id)
	return true
}
