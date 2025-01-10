package manager

import ()

type Settings struct {
	MinPlayerCount uint
	IdLength       uint
}

func CreateSettings() *Settings {
	return &Settings{
		MinPlayerCount: 3,
		IdLength:       3,
	}
}
