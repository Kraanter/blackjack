package main

import (
	"fmt"
	"time"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

func main() {
	game := blackjack.CreateGame(blackjack.CreateSettings())
	game.Settings.TimeBetweenRounds = 100 * time.Millisecond
	players := make([]*blackjack.Player, 0, 3)

	players = append(players, game.AddPlayerWithBalance(10))
	players = append(players, game.AddPlayerWithBalance(10))
	players = append(players, game.AddPlayerWithBalance(10))

	game.OnGameUpdate = func(game *blackjack.BlackjackGame) {
		fmt.Printf("\n---\ngame_update: %v\n\nplayers:", game.GameState)

		for _, player := range players {
			fmt.Println(player.String())
		}

		fmt.Println("Dealer: ", game.Dealer.String())
	}
	game.OnPlayerTurn = func(pi blackjack.PlayerId) {
		game.PlayerHit(pi)
		game.PlayerStand(pi)
	}

	go func() {
		time.Sleep(10 * time.Millisecond)
		err := game.SetPlayerBet(players[0].PlayerNum, 5)
		err = game.SkipPlayerBet(players[1].PlayerNum)
		err = game.SetPlayerBet(players[2].PlayerNum, 2)
		if err != nil {
			println(err.Error())
		}
	}()

	game.Start()

	for game.GameState != blackjack.NoState {
		println("Waiting")
		time.Sleep(1 * time.Second)
	}
}
