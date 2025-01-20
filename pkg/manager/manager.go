package manager

import (
	"time"
)

type Manager struct {
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
	for id, game := range m.gameMap {
		if game.GetPlayerCount() == 0 {
			m.RemoveGame(id)
		}
	}
}

func (m *Manager) GetGameCount() uint {
	return uint(len(m.gameMap))
}
