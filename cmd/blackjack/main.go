package main

import (
	"fmt"

	"github.com/kraanter/blackjack/pkg/blackjack"
)

func main() {
	fmt.Println(*blackjack.CreateCard(blackjack.Two, blackjack.Clubs))
}
