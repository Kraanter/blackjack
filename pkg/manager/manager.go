package manager

import (
	"github.com/kraanter/blackjack/pkg/blackjack"
)

type Manager struct {
	gameMap  map[GameId]*blackjack.BlackjackGame
	Settings *Settings
}

func CreateManager(settings *Settings) *Manager {
	if settings == nil {
		settings = CreateSettings()
	}

	return &Manager{
		gameMap:  make(map[GameId]*blackjack.BlackjackGame),
		Settings: settings,
	}
}

func (m *Manager) GetGameCount() uint {
	return uint(len(m.gameMap))
}
