package manager

import (
	"context"
	"fmt"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

// Returns ID 0 if no game is found and game pointer will be nil
func (m *Manager) GetJoinableGame() (GameId, *ManagedGame) {
	// Join a random game
	m.mu.RLock()
	games := make(map[GameId]*ManagedGame, len(m.gameMap))
	for k, v := range m.gameMap {
		games[k] = v
	}
	m.mu.RUnlock()
	for k, v := range games {
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
	m.mu.RLock()
	defer m.mu.RUnlock()
	game, ok := m.gameMap[id]

	if !ok {
		return nil, GameNotFoundError
	}

	return game, nil

}

func (m *Manager) createNewGame() (GameId, *ManagedGame) {
	newGame := createManagedGame(context.Background())
	for {
		gameId := CreateRandomGameId(m.Settings.IdLength)

		if m.addGameWithID(gameId, newGame) {
			return gameId, newGame
		}
	}
}

func (m *Manager) addGameWithID(id GameId, game *ManagedGame) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.gameMap[id]
	if ok {
		return false
	}

	m.gameMap[id] = game

	return true
}

func (m *Manager) RemoveGame(id GameId) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	game := m.gameMap[id]
	if game == nil {
		return false
	}

	if game.GetPlayerCount() != 0 {
		return false
	}

	delete(m.gameMap, id)
	return true
}

func createGameUpdateHandler(manGame *ManagedGame) func(game *blackjack.GameSnapshot) {
	return func(game *blackjack.GameSnapshot) {
		manGame.mu.RLock()
		handlers := make([]func(*blackjack.GameSnapshot), 0, len(manGame.Players))
		for _, player := range manGame.Players {
			if handler := player.gameUpdateHandler(); handler != nil {
				handlers = append(handlers, handler)
			}
		}
		manGame.mu.RUnlock()
		for _, handler := range handlers {
			handler(game)
		}
	}
}
