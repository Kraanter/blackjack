package main

import (
	"github.com/kraanter/blackjack/pkg/manager"
)

func main() {
	settings := manager.CreateSettings()
	settings.MinPlayerCount = 1
	manager := manager.CreateManager(nil)

	println(manager.GetGameCount())
}
