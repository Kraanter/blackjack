package blackjack_test

import (
	"testing"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

func TestBlackJackAddPlayerIncrementCount(t *testing.T) {
	game := blackjack.CreateGame(nil)
	game.AddPlayerWithBalance(0)
	game.AddPlayerWithBalance(0)
	game.AddPlayerWithBalance(0)
	game.AddPlayerWithBalance(0)
	fifthPlayer := game.AddPlayerWithBalance(0)
	var want uint = 5

	playerNum := fifthPlayer.PlayerNum

	if playerNum != want {
		t.Errorf("5x AddPlayer playerCount = %v, want to be %v", playerNum, want)
	}
}

func TestBlackJackRemovePlayerRemovesPlayer(t *testing.T) {
	game := blackjack.CreateGame(nil)
	want := uint(10)
	game.AddPlayerWithBalance(0)
	secondPlayer := game.AddPlayerWithBalance(want)
	game.AddPlayerWithBalance(0)

	balance, err := game.RemovePlayer(secondPlayer.PlayerNum)
	if err != nil {
		t.Errorf("Removing a player that is just added should not return a error")
	}
	if balance != want {
		t.Errorf("Removing a player should return its balance")
	}

	_, ok := game.GetPlayer(secondPlayer.PlayerNum)
	if ok != false {
		t.Errorf("Searching a player after removing it should not find it")
	}
}

func TestBlackJackGetNonExistantPlayerReturnError(t *testing.T) {
	game := blackjack.CreateGame(nil)

	_, ok := game.GetPlayer(1)
	if ok {
		t.Fatalf("Searching a player that does not exist should not return ok")
	}
}
