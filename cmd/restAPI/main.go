package main

import (
	"time"

	"github.com/kraanter/blackjack/pkg/restAPI"
	"github.com/kraanter/blackjack/pkg/restAPI/middleware"
)

func main() {
	go func() {
		for {
			println(len(middleware.UserMap))
			time.Sleep(1 * time.Second)
		}
	}()
	restapi.Start()

}
