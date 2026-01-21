package blackjack_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

func setupGame(playerCount uint) *blackjack.BlackjackGame {
	game := blackjack.CreateGame()

	for range playerCount {
		game.AddPlayerWithBalance(10)
	}

	return game
}

func TestPlayerJoiningDuringBettingState(t *testing.T) {
	game := setupGame(1)

	game.Initialize()

	players := game.GetPlayersWihoutBets()
	for _, id := range players {
		err := game.SetPlayerBet(id, 10)
		if err != nil {
			t.Fatalf("Failed to set player (%v) bet (%v)", id, game.PlayerMap)
		}
	}

	gameContext, cancelGame := context.WithTimeout(context.Background(), time.Millisecond)
	game.Start(gameContext)
	cancelGame()

	newPlayer := game.AddPlayerWithBalance(10)

	expected := 2
	actual := len(game.PlayerMap)
	if expected != actual {
		t.Fatalf("Expected playermap to contain %v players but had %v", expected, actual)
	}

	err := game.SetPlayerBet(newPlayer.PlayerNum, 5)
	if !errors.Is(err, blackjack.WrongGameStateError) {
		t.Fatalf("New player should not be able to bet if joined after betting state")
	}

	for _, id := range players {
		err := game.PlayerStand(id)
		if err != nil {
			t.Fatalf("Failed to stand player after adding a new one during the round (%s)", err)	
		}
	}

	actual = len(game.PlayerMap)
	if expected != actual {
		t.Fatalf("Expected playermap to contain %v players but had %v", expected, actual)
	}
}
