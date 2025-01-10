package manager_test

import (
	"context"
	"testing"
	"time"

	"github.com/kraanter/blackjack/pkg/blackjack"
	"github.com/kraanter/blackjack/pkg/manager"
)

func TestManagedPlayerCanPlayFullGame(t *testing.T) {
	man := manager.CreateManager(nil)
	player := man.JoinRandomGame(context.Background(), 10)

	// NOTE: For debugging
	// var printMutex sync.Mutex
	//
	// player.Game.OnGameUpdate = func(game *blackjack.BlackjackGame) {
	// 	printMutex.Lock()
	// 	defer printMutex.Unlock()
	// 	fmt.Printf("\n---\ngame_update: %v\n\nplayers: \n", game.GameState)
	//
	// 	fmt.Println(player.String())
	//
	// 	fmt.Println("Dealer: ", game.Dealer.String())
	// }

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

	time.Sleep(10 * time.Millisecond)

	if player.Game.GameState != blackjack.NoState {
		t.Fatalf("Game should be in no state (%v) after a match is finished, was %v", blackjack.NoState, player.Game.GameState)
	}

	_, err := player.Leave()
	if err != nil {
		t.Fatalf("Leaving game after match is done should not return error. Error = %v", err)
	}
}
