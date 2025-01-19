package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kraanter/blackjack/pkg/manager"
	"github.com/kraanter/blackjack/pkg/restAPI/users"
)

var joinRoute = createNoAuthRoute("POST /join", joinGameHandler)

var managerSettings = manager.CreateSettings()
var gameManager = manager.CreateManager(managerSettings)

func joinGameHandler(res http.ResponseWriter, req *http.Request) {
	gameIdStr := req.URL.Query().Get("code")
	gameId, err := strconv.Atoi(gameIdStr)
	var player *manager.ManagedPlayer
	if err != nil {
		player = gameManager.JoinRandomGame(context.Background(), 100)
	} else {
		player = gameManager.JoinGame(context.Background(), 100, manager.GameId(gameId))
	}

	if player == nil {
		handleError(res, "Could not join game", http.StatusUnprocessableEntity)
		return
	}

	fmt.Printf("Player has joined: %v\n", player)

	userCookieValue := users.RegisterUser(player, context.Background())

	userCookie := http.Cookie{
		Name:     users.CookiePlayerIdKey,
		Value:    userCookieValue,
		Secure:   true,
		HttpOnly: true,
		Quoted:   false,
		Path:     "/",
	}

	http.SetCookie(res, &userCookie)

	writeStructToResponse(res, player, http.StatusCreated)
}
