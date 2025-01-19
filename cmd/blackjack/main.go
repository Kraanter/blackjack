package main

import (
	"fmt"
	"sync"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

func main() {
	game := blackjack.CreateGame()
	var printMutex sync.Mutex
	players := make([]*blackjack.Player, 0, 3)

	players = append(players, game.AddPlayerWithBalance(10))
	players = append(players, game.AddPlayerWithBalance(10))
	players = append(players, game.AddPlayerWithBalance(10))

	game.OnGameUpdate = func(game *blackjack.BlackjackGame) {
		printMutex.Lock()
		defer printMutex.Unlock()
		fmt.Printf("\n---\ngame_update: %v\n\nplayers: \n", game.GameState)

		for _, player := range players {
			fmt.Println(player.String())
		}

		fmt.Println("Dealer: ", game.Dealer.String())
	}
	game.OnPlayerTurn = func(pi blackjack.PlayerId) {
		game.PlayerHit(pi)
		game.PlayerStand(pi)
	}

	game.Start()

	err := game.SetPlayerBet(players[0].PlayerNum, 5)
	err = game.SkipPlayerBet(players[1].PlayerNum)
	err = game.SetPlayerBet(players[2].PlayerNum, 2)
	if err != nil {
		println(err.Error())
	}

	for game.GameState != blackjack.NoState {
	}
	printMutex.Lock()
	printMutex.Unlock()
}
