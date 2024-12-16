package blackjack_test

import "blackjack/pkg/blackjack"

func forAllFaces(callback func(face blackjack.Face)) {
	faces := []blackjack.Face{blackjack.Ace, blackjack.Two, blackjack.Three, blackjack.Four, blackjack.Five, blackjack.Six, blackjack.Seven, blackjack.Eight, blackjack.Nine, blackjack.Ten, blackjack.Jack, blackjack.Queen, blackjack.King}

	for _, face := range faces {
		callback(face)
	}
}

func forAllSuits(callback func(suit blackjack.Suit)) {
	suits := []blackjack.Suit{blackjack.Hearts, blackjack.Diamonds, blackjack.Clubs, blackjack.Spades}
	for _, suit := range suits {
		callback(suit)
	}
}

func forAllSuitFaceCombiniations(callback func(suit blackjack.Suit, face blackjack.Face)) {
	forAllSuits(func(suit blackjack.Suit) {
		forAllFaces(func(face blackjack.Face) {
			callback(suit, face)
		})
	})
}
