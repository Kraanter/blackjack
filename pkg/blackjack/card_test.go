package blackjack_test

import (
	"testing"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

// This test does not do anything because it is just in the String Func
func TestAllCardStrings(t *testing.T) {
	blackjack.ForAllSuitFaceCombiniations(func(suit blackjack.Suit, face blackjack.Face) {
		card := blackjack.CreateCard(face, suit)
		want := face.String() + suit.String()

		cardString := card.String()

		if cardString != want {
			t.Fatalf(`CreateCard(%v, %v) = %q, want match for %q`, face, suit, cardString, want)
		}
	})
}
