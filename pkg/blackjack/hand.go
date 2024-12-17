package blackjack

import ()

type Hand struct {
	cards []*Card
	bet   uint
	// True if all cards have been dealed
	locked bool
}

func CreateHand(bet uint) *Hand {
	return &Hand{
		cards: make([]*Card, 0),
		bet:   bet,
	}
}

func (hand *Hand) Total() int {
	total := 0

	aceCount := 0
	for _, card := range hand.cards {
		if card.Face == Ace {
			aceCount++
		}
		total += card.GetValue()
	}

	for aceCount > 0 && total > 21 {
		total -= 10
		aceCount--
	}

	return total
}

func (hand *Hand) AddCard(card *Card) {
	if hand.locked {
		return
	}

	hand.cards = append(hand.cards, card)
}

func (hand *Hand) lock() {
	hand.locked = true
}
