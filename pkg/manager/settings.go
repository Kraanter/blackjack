package manager

import "time"

type Settings struct {
	MinPlayerCount     uint
	IdLength           uint
	CleanupTimerLength time.Duration
}

func CreateSettings() *Settings {
	return &Settings{
		MinPlayerCount:     3,
		IdLength:           3,
		CleanupTimerLength: 30 * time.Second,
	}
}
