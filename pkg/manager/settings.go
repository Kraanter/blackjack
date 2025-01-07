package manager

import "github.com/kraanter/blackjack/pkg/blackjack"

type Settings struct {
	MinPlayerCount    uint
	IdLength          uint
	BlackjackSettings blackjack.Settings
}

func createSettings() *Settings {
	return &Settings{
		MinPlayerCount:    3,
		IdLength:          3,
		BlackjackSettings: *blackjack.CreateSettings(),
	}
}
