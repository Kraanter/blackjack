package blackjack

import (
	"testing"
)

func TestPlayerPlacesBetWithdrawsFundsAndCreatesHand(t *testing.T) {
	player := CreatePlayer(0, 10)
	want := 0

	err := player.placeBet(10)

	if err != nil {
		t.Fatalf(`player.PlaceBet(10) on balance of 10, returned error, %q`, err.Error())
	}

	if player.Balance != 0 {
		t.Fatalf(`player.Balance = %v, want match for %v`, player.Balance, want)
	}
}

func TestPlayerPlacesBetThrowsErrorIfNotEnoughBalance(t *testing.T) {
	player := CreatePlayer(0, 0)

	err := player.placeBet(1)

	if err == nil {
		t.Fatalf(`player.PlaceBet(1) on balance of 0, should return error`)
	}
}

func TestPlayerPlacesBetThrowsErrorIfBetToLow(t *testing.T) {
	player := CreatePlayer(0, 0)

	err := player.placeBet(0)

	if err == nil {
		t.Fatalf(`player.PlaceBet(0), should return error`)
	}
}

func TestPlayerDestroyReturnsBalance(t *testing.T) {
	want := uint(10)
	player := CreatePlayer(0, want)

	balance := player.Destroy()
	if balance != want {
		t.Fatalf(`player.Destory() = %v, want match for %v`, balance, want)
	}
}

func TestPlayerHasNoActiveHandAt21OrBust(t *testing.T) {
	for _, cards := range [][]*Card{
		{CreateCard(Ace, Hearts), CreateCard(King, Spades)},
		{CreateCard(King, Hearts), CreateCard(Queen, Spades), CreateCard(Two, Clubs)},
	} {
		player := CreatePlayer(0, 10)
		if err := player.placeBet(1); err != nil {
			t.Fatal(err)
		}
		for _, card := range cards {
			player.Hands[0].AddCard(card)
		}
		if hand := player.GetActiveHand(); hand != nil {
			t.Fatalf("GetActiveHand() = %v, want nil", hand)
		}
	}
}
