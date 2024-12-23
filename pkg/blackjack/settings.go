package blackjack

import "time"

type Settings struct {
	DealCardTime time.Duration
}

func CreateSettings() *Settings {
	return &Settings{
		DealCardTime: 500 * time.Millisecond,
	}
}
