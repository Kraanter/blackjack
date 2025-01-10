package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kraanter/blackjack/pkg/blackjack"
	"github.com/kraanter/blackjack/pkg/manager"
)

func main() {
	man := manager.CreateManager(nil)
	player := man.JoinRandomGame(context.Background(), 10)
	var printMutex sync.Mutex

	player.Game.OnGameUpdate = func(game *blackjack.BlackjackGame) {
		printMutex.Lock()
		defer printMutex.Unlock()
		fmt.Printf("\n---\ngame_update: %v\n\nplayers: \n", game.GameState)

		fmt.Println(player.String())

		fmt.Println("Dealer: ", game.Dealer.String())
	}
	player.Game.OnPlayerTurn = func(pi blackjack.PlayerId) {
		player.Hit()
		player.Stand()
	}

	go func() {
		time.Sleep(10 * time.Millisecond)
		err := player.Bet(10)
		if err != nil {
			println("Error while betting", err.Error())
		}
	}()

	player.Game.Start()

	for player.Game.GameState != blackjack.NoState {
	}
	printMutex.Lock()
	printMutex.Unlock()
}
