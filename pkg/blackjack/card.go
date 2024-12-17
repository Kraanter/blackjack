package blackjack

import (
	"strconv"

	"github.com/kraanter/blackjack/pkg/assert"
)

type Suit int
type Face int

const (
	Ace Face = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King

	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

type Card struct {
	Face Face
	Suit Suit
}

func CreateCard(face Face, suit Suit) *Card {
	return &Card{
		Face: face,
		Suit: suit,
	}
}

func (c Card) String() string {
	return c.Face.String() + c.Suit.String()
}

func (c Card) GetValue() int {
	return c.Face.toValue()
}

func (face Face) String() string {
	switch face {
	case Ace:
		return "A"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return strconv.Itoa(face.toValue())
	}
}

func (face Face) toValue() int {
	switch face {
	case Ace:
		return 11
	case Jack, Queen, King:
		return 10
	}

	value := int(face) + 1

	return value
}

func (suit Suit) String() string {
	switch suit {
	case Hearts:
		return ""
	case Diamonds:
		return "◆"
	case Clubs:
		return "♣"
	case Spades:
		return "♠"
	}

	assert.Never("Can't convert suit (%v) to a string", suit)
	return ""
}
