package manager

import (
	"sync"
	"time"
)

type Manager struct {
	mu       sync.RWMutex
	gameMap  map[GameId]*ManagedGame
	Settings *Settings
}

func CreateManager(settings *Settings) *Manager {
	if settings == nil {
		settings = CreateSettings()
	}

	manager := &Manager{
		gameMap:  make(map[GameId]*ManagedGame),
		Settings: settings,
	}

	go manager.cleanupRoutine()

	return manager
}

func (m *Manager) cleanupRoutine() {
	for {
		<-time.After(m.Settings.CleanupTimerLength)

		m.cleanupEmptyGames()
	}

}

func (m *Manager) cleanupEmptyGames() {
	m.mu.RLock()
	games := make(map[GameId]*ManagedGame, len(m.gameMap))
	for id, game := range m.gameMap {
		games[id] = game
	}
	m.mu.RUnlock()
	for id, game := range games {
		if game.GetPlayerCount() == 0 {
			m.RemoveGame(id)
		}
	}
}

func (m *Manager) GetGameCount() uint {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return uint(len(m.gameMap))
}
