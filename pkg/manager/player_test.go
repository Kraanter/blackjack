package manager_test

import (
	"testing"

	"github.com/kraanter/blackjack/pkg/blackjack"
	"github.com/kraanter/blackjack/pkg/manager"
)

func TestManagedPlayerCanPlayFullGame(t *testing.T) {
	man := manager.CreateManager(nil)
	player := man.JoinRandomGame(10)

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

	player.Game.Start()
	err := player.Bet(10)
	if err != nil {
		t.Fatalf("Error while betting: %v", err)
	}

	for player.Game.GameState != blackjack.NoState {
	}

	_, err = player.Leave()
	if err != nil {
		t.Fatalf("Leaving game after match is done should not return error. Error = %v", err)
	}
}
