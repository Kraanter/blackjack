package blackjack_test

import (
	"blackjack/pkg/blackjack"
	"testing"
)

func TestPlayerPlacesBetWithdrawsFundsAndCreatesHand(t *testing.T) {
	player := blackjack.CreatePlayer(10)
	want := 0

	err := player.PlaceBet(10)

	if err != nil {
		t.Fatalf(`player.PlaceBet(10) on balance of 10, returned error, %q`, err.Error())
	}

	if player.Balance != 0 {
		t.Fatalf(`player.Balance = %v, want match for %v`, player.Balance, want)
	}
}

func TestPlayerPlacesBetThrowsErrorIfNotEnoughBalance(t *testing.T) {
	player := blackjack.CreatePlayer(0)

	err := player.PlaceBet(1)

	if err == nil {
		t.Fatalf(`player.PlaceBet(1) on balance of 0, should return error`)
	}
}

func TestPlayerPlacesBetThrowsErrorIfBetToLow(t *testing.T) {
	player := blackjack.CreatePlayer(0)

	err := player.PlaceBet(0)

	if err == nil {
		t.Fatalf(`player.PlaceBet(0), should return error`)
	}
}
