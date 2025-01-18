package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kraanter/blackjack/pkg/manager"
	"github.com/kraanter/blackjack/pkg/restAPI/middleware"
)

var ApiRoutes = []ApiRoute{joinRoute, authRoute}

var joinRoute = ApiRoute{
	Pattern: "GET /join",
	Handler: joinGameHandler,
	noAuth:  true,
}

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
		http.Error(res, "Could not connect to game", http.StatusUnprocessableEntity)
		return
	}

	fmt.Printf("joined %v\n", player)

	userCookieValue := middleware.RegisterUser(player, context.Background())

	userCookie := http.Cookie{
		Name:     middleware.CookiePlayerIdKey,
		Value:    userCookieValue,
		Secure:   true,
		HttpOnly: true,
		Quoted:   false,
		Path:     "/",
	}

	http.SetCookie(res, &userCookie)

	res.Write([]byte("hi"))
}

var authRoute = ApiRoute{
	Pattern: "GET /auth",
	Handler: authHandler,
}

func authHandler(res http.ResponseWriter, req *http.Request) {
	user := middleware.GetUserFromReq(req)
	if user == nil {
		res.Write([]byte("No user found"))
		return
	}

	data, err := json.Marshal(user)

	println(data)
	if err != nil {
		println(err.Error())
	}
	res.Write(data)
}
