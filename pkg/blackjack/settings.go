package blackjack

import "time"

type Settings struct {
	DealCardTime      time.Duration
	TimeBetweenRounds time.Duration
}

func CreateSettings() *Settings {
	return &Settings{
		DealCardTime:      500 * time.Millisecond,
		TimeBetweenRounds: 10 * time.Second,
	}
}
