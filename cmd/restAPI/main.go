package main

import (
	"flag"

	"github.com/kraanter/blackjack/pkg/restAPI"
)

func main() {
	portPtr := flag.Int("port", 42069, "The port to run the server on")
	flag.Parse()

	restapi.Start(*portPtr)
}
