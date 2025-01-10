package manager

import ()

type Settings struct {
	MinPlayerCount uint
	IdLength       uint
}

func createSettings() *Settings {
	return &Settings{
		MinPlayerCount: 3,
		IdLength:       3,
	}
}
