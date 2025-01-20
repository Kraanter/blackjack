package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kraanter/blackjack/pkg/blackjack"
	"github.com/kraanter/blackjack/pkg/manager"
	"github.com/kraanter/blackjack/pkg/restAPI/games"
	"github.com/kraanter/blackjack/pkg/restAPI/users"
)

var joinRoute = createNoAuthRoute("POST /join", joinGameHandler)
var leaveRoute = createRoute("DELETE /leave", leaveGameHandler)

func joinGameHandler(w http.ResponseWriter, r *http.Request) {
	gameIdStr := r.URL.Query().Get("code")
	gameId, err := strconv.Atoi(gameIdStr)
	var player *manager.ManagedPlayer
	if err != nil {
		player = games.GameManager.JoinRandomGame(100)
	} else {
		player = games.GameManager.JoinGame(100, manager.GameId(gameId))
	}

	if player == nil {
		handleError(w, "Could not join game", http.StatusUnprocessableEntity)
		return
	}

	userCookieValue := users.RegisterUser(player, context.Background())

	userCookie := http.Cookie{
		Name:     users.CookiePlayerIdKey,
		Value:    userCookieValue,
		Secure:   true,
		HttpOnly: true,
		Quoted:   false,
		Path:     "/",
	}

	http.SetCookie(w, &userCookie)

	player.OnGameUpdate = playerUpdateHandler(player)
	player.Game.Start()

	writeStructToResponse(w, player, http.StatusCreated)
}

func leaveGameHandler(w http.ResponseWriter, r *http.Request) {
	user := users.GetUserFromReq(r)
	if user == nil {
		handleUnauthenticated(w)
		return
	}

	users.RemoveAuthUser(user)

	w.WriteHeader(http.StatusOK)
}

func playerUpdateHandler(player *manager.ManagedPlayer) func(game *blackjack.BlackjackGame) {
	return func(game *blackjack.BlackjackGame) {
		fmt.Printf("%v %v\n", player, game)
	}
}
