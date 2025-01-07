package manager

import (
	"github.com/kraanter/blackjack/pkg/blackjack"
)

type Manager struct {
	gameMap  map[GameId]*blackjack.BlackjackGame
	settings *Settings
}

func CreateManager(settings *Settings) *Manager {
	if settings == nil {
		settings = createSettings()
	}

	return &Manager{
		gameMap:  make(map[GameId]*blackjack.BlackjackGame),
		settings: settings,
	}
}

func (m *Manager) GetGameCount() uint {
	return uint(len(m.gameMap))
}
