package blackjack_test

import (
	"fmt"
	"testing"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

func TestGameLoopSingleRound(t *testing.T) {
	game := blackjack.CreateGame(blackjack.CreateSettings())

	player1 := game.AddPlayerWithBalance(10)
	player2 := game.AddPlayerWithBalance(20)
	player3 := game.AddPlayerWithBalance(30)

}
