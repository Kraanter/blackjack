package main

import (
	"flag"
	"fmt"

	"github.com/kraanter/blackjack/pkg/restAPI"
)

func main() {
	portPtr := flag.Int("port", 42069, "The port to run the server on")
	flag.Parse()

	fmt.Printf("portPtr: %v\n", *portPtr)

	restapi.Start(*portPtr)
}
