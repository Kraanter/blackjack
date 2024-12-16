package blackjack_test

import (
	"blackjack/pkg/blackjack"
	"testing"
)

// This test does not do anything because it is just in the String Func
func TestAllCardStrings(t *testing.T) {
	forAllSuitFaceCombiniations(func(suit blackjack.Suit, face blackjack.Face) {
		card := blackjack.CreateCard(face, suit)
		want := face.String() + suit.String()

		cardString := card.String()

		if cardString != want {
			t.Fatalf(`CreateCard(%v, %v) = %q, want match for %q`, face, suit, cardString, want)
		}
	})
}
