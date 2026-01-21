package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/kraanter/blackjack/pkg/blackjack"
	"github.com/kraanter/blackjack/pkg/manager"
)

func main() {
	man := manager.CreateManager(nil)
	player := man.JoinRandomGame(10)
	var printMutex sync.Mutex

	player.Game.OnGameUpdate = func(game *blackjack.BlackjackGame) {
		printMutex.Lock()
		defer printMutex.Unlock()
		fmt.Printf("\n---\ngame_update: %v\n\nplayers: \n", game.GameState.Get())

		fmt.Println(player.String())

		fmt.Println("Dealer: ", game.Dealer.Get().String())
	}
	player.Game.OnPlayerTurn = func(pi blackjack.PlayerId) {
		player.Hit()
		player.Stand()
	}

	err := player.Bet(10)
	if err != nil {
		println("Error while betting", err.Error())
		panic(1)
	}

	player.Game.Start(context.Background())

	for player.Game.GameState.Get() != blackjack.NoState {
	}
	printMutex.Lock()
}
