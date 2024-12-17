package blackjack_test

import (
	"testing"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

func TestShoeDrawsLastCardNotInDeckAnymore(t *testing.T) {
	shoe := blackjack.CreateShoe(1)
	want := shoe.Size() - 1

	_ = shoe.DrawCard()
	shoeSize := shoe.Size()

	if shoeSize != want {
		t.Fatalf("shoe.DrawCard(1 Deck) contains %v cards, wants %v cards", shoeSize, want)
	}
}
