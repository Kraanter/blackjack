package main

import (
	"blackjack/pkg/blackjack"
	"fmt"
)

func main() {
	fmt.Println(*blackjack.CreateCard(blackjack.Two, blackjack.Clubs))
}
