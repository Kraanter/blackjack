package blackjack_test

import (
	"slices"
	"testing"
	"time"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

func TestGameLoopSingleRound(t *testing.T) {
	settings := blackjack.CreateSettings()
	settings.DealCardTime = 0 * time.Millisecond
	game := blackjack.CreateGame(settings)
	players := make([]*blackjack.Player, 0, 3)

	players = append(players, game.AddPlayerWithBalance(10))
	players = append(players, game.AddPlayerWithBalance(10))
	players = append(players, game.AddPlayerWithBalance(10))

	// NOTE: Logging for debugging
	// game.OnGameUpdate = func(game *blackjack.BlackjackGame) {
	// 	fmt.Printf("\n---\ngame_update: %v\n\nplayers:\n", game.GameState)
	//
	// 	for _, player := range players {
	// 		fmt.Println(player.String())
	// 	}
	//
	// 	fmt.Println("Dealer: ", game.Dealer.String())
	// }

	game.OnPlayerTurn = func(pi blackjack.PlayerId) {
		ok, err := game.PlayerHit(pi)
		if !ok || err != nil {
			t.Fatalf("Player (%v) decision 'hit' went wrong: ok: %v err: %v", pi, ok, err)
		}
		err = game.PlayerStand(pi)
		if err != nil {
			t.Fatalf("Player (%v) decision 'stand' went wrong: %v", pi, err)
		}

	}

	go func() {
		err := game.SetPlayerBet(players[0].PlayerNum, 5)
		err = game.SetPlayerBet(players[2].PlayerNum, 2)
		if err != nil {
			t.Fatalf("Setting player bets went wrong: %v", err.Error())
		}
		err = game.SkipPlayerBet(players[1].PlayerNum)
		if err != nil {
			t.Fatalf("Skipping player bet went wrong: %v", err.Error())
		}
	}()

	game.Start()

	time.Sleep(1 * time.Millisecond)

	if game.GameState != blackjack.NoState {
		t.Fatalf("Game should be in no state (%v) after a match is finished, was %v", blackjack.NoState, game.GameState)
	}
	if game.Dealer != nil {
		t.Fatalf("Dealer should not have a hand at the end of the game")
	}
	if slices.ContainsFunc(players, func(player *blackjack.Player) bool { return player.Hand != nil }) {
		t.Fatalf("All players should nto have a hand at the end of the game")
	}

}
