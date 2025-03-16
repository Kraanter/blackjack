package main

import (
	"flag"
	"time"

	"github.com/kraanter/blackjack/pkg/restAPI"
)

func main() {
	portPtr := flag.Int("port", 42069, "The port to run the server on")
	flag.Parse()

	settings := restapi.CreateDefaultServerSettings()
	settings.Port = uint(*portPtr)
	settings.TimeBetweenGameUpdates = 250 * time.Millisecond

	restapi.Start(settings)
}
