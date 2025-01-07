package blackjack_test

import (
	"testing"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

func TestHandTotalWith3Aces(t *testing.T) {
	hand := blackjack.CreateHand(0)
	want := 13

	hand.AddCard(blackjack.CreateCard(blackjack.Ace, blackjack.Hearts))
	hand.AddCard(blackjack.CreateCard(blackjack.Ace, blackjack.Diamonds))
	hand.AddCard(blackjack.CreateCard(blackjack.Ace, blackjack.Spades))

	total := hand.Total()

	if total != want {
		t.Fatalf("hand.Total(2 aces) = %v, want match for %v", total, want)
	}
}

func TestHandTotalWith2AcesKing(t *testing.T) {
	hand := blackjack.CreateHand(0)
	want := 12

	hand.AddCard(blackjack.CreateCard(blackjack.Ace, blackjack.Hearts))
	hand.AddCard(blackjack.CreateCard(blackjack.Ace, blackjack.Diamonds))
	hand.AddCard(blackjack.CreateCard(blackjack.King, blackjack.Spades))

	total := hand.Total()

	if total != want {
		t.Fatalf("hand.Total(2 aces, 1 king) = %v, want match for %v", total, want)
	}
}

func TestHandTotalWithAceKing(t *testing.T) {
	hand := blackjack.CreateHand(0)
	want := 21

	hand.AddCard(blackjack.CreateCard(blackjack.Ace, blackjack.Hearts))
	hand.AddCard(blackjack.CreateCard(blackjack.King, blackjack.Spades))

	total := hand.Total()

	if total != want {
		t.Fatalf("hand.Total(1 ace, 1 king) = %v, want match for %v", total, want)
	}
}
