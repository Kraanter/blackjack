package blackjack

import (
	"math/rand"
	"slices"
)

type Shoe struct {
	cards     []*Card
	deckCount uint
}

func CreateShoe(amountOf uint) *Shoe {
	totalCards := createStackOfNRandomizedDecks(amountOf)
	// TODO: Implement the 70-90% cut in the shoe
	return &Shoe{
		cards:     totalCards,
		deckCount: amountOf,
	}
}

func createStackOfNRandomizedDecks(deckCount uint) []*Card {
	deck := make([]*Card, 0)
	for range deckCount {
		deck = slices.Concat(deck, CreateDeckOfCards())
	}

	for i := range deck {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	return deck
}

func (s *Shoe) DrawCard() *Card {
	if len(s.cards) == 0 {
		s.cards = createStackOfNRandomizedDecks(s.deckCount)
	}

	card := s.cards[len(s.cards)-1]
	s.cards = slices.Delete(s.cards, len(s.cards)-1, len(s.cards))

	return card
}

func (s *Shoe) Size() int {
	return len(s.cards)
}
