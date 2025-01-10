package manager_test

import (
	"context"
	"testing"
	"time"

	"github.com/kraanter/blackjack/pkg/blackjack"
	"github.com/kraanter/blackjack/pkg/manager"
)

func TestPlayerContextDoneRemovesPlayerAndGameIfEmpty(t *testing.T) {
	man := manager.CreateManager(nil)
	playerContext, cancel := context.WithCancel(context.Background())
	player := man.JoinRandomGame(playerContext, 10)

	want := uint(1)
	gameCount := man.GetGameCount()
	if gameCount != want {
		t.Fatalf("GetGameCount() = %v, want gamecount to be %v", gameCount, want)
	}

	cancel()
	time.Sleep(10 * time.Millisecond)
	gameCount = man.GetGameCount()
	want = 0
	if gameCount != want {
		t.Fatalf("GetGameCount() = %v, want gamecount to be %v because of cancelled context", gameCount, want)
	}
	if player.Game != nil {
		t.Fatalf("ManagedPlayer should not be connected to any game because context is gone")
	}
	if player.Player != nil {
		t.Fatalf("ManagedPlayer should not be connected to any blackjack player because context is gone")
	}
}

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
