package blackjack

import (
	"fmt"
	"strings"
)

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
	if hand == nil {
		return total
	}

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

func (hand *Hand) AddCard(card *Card) bool {
	if hand == nil || hand.locked {
		return false
	}

	hand.cards = append(hand.cards, card)
	return true
}

func (hand *Hand) lock() {
	if hand == nil {
		return
	}
	hand.locked = true
}

func (hand *Hand) isLocked() bool {
	if hand == nil {
		return false
	}

	return hand.locked
}

func (hand *Hand) String() string {
	cards := make([]string, 0)

	if hand != nil {
		for _, card := range hand.cards {
			cards = append(cards, card.String())
		}
	} else {
		cards = append(cards, "-")
	}

	return strings.Join(cards, "  ") + fmt.Sprintf(" total: %v lock: %v", hand.Total(), hand.isLocked())
}
