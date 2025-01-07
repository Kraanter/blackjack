package main

import (
	"github.com/kraanter/blackjack/pkg/manager"
)

func main() {
	manager := manager.CreateManager(nil)

	game := manager.GetJoinableGame()

	game.AddPlayerWithBalance(100)
	game = manager.GetJoinableGame()
	game.AddPlayerWithBalance(200)
	game.AddPlayerWithBalance(300)
	game = manager.GetJoinableGame()
	game.AddPlayerWithBalance(400)

	println(manager.GetGameCount())
}
