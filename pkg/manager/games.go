package manager

import (
	"fmt"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

// Returns ID 0 if no game is found and game pointer will be nil
func (m *Manager) GetJoinableGame() (GameId, *ManagedGame) {
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

func (m *Manager) GetGameWithId(id GameId) (*ManagedGame, error) {
	game, ok := m.gameMap[id]

	if !ok {
		return nil, GameNotFoundError
	}

	return game, nil

}

func (m *Manager) createNewGame() (GameId, *ManagedGame) {
	newGame := createManagedGame()
	for {
		gameId := CreateRandomGameId(m.Settings.IdLength)

		if m.addGameWithID(gameId, newGame) {
			return gameId, newGame
		}
	}
}

func (m *Manager) addGameWithID(id GameId, game *ManagedGame) bool {
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

func createGameUpdateHandler(manGame *ManagedGame) func(game *blackjack.BlackjackGame) {
	return func(game *blackjack.BlackjackGame) {
		for _, player := range manGame.Players {
			if player.OnGameUpdate != nil {
				player.OnGameUpdate(game)
			}
		}
	}
}
