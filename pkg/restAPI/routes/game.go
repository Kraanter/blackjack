package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/kraanter/blackjack/pkg/blackjack"
	"github.com/kraanter/blackjack/pkg/manager"
	"github.com/kraanter/blackjack/pkg/restAPI/games"
	"github.com/kraanter/blackjack/pkg/restAPI/users"
)

var joinRoute = createNoAuthRoute("POST /join", joinGameHandler)
var leaveRoute = createRoute("DELETE /leave", leaveGameHandler)
var gameStateRoute = createRoute("GET /gamestate", gameStateHandler)

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

	player.Game.Initialize()

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

type writeContextKey string

var contextKey = writeContextKey("writer")

func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	user := users.GetUserFromReq(r)
	if user == nil {
		handleUnauthenticated(w)
		return
	}

	isSSE := r.URL.Query().Has("sse")

	if !isSSE {
		writeStructToResponse(w, user.Player.Game, http.StatusOK)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Bind event
	oldContext := user.WriteContext()
	user.SetUserWriter(w, r.Context())
	if oldContext != nil {
		<-oldContext.Done()
	}
	user.Player.OnGameUpdate = playerUpdateHandler(user)

	select {
	case <-r.Context().Done():
		user.Player.OnGameUpdate = nil
	case <-user.WriteContext().Done():
		sendSSEvent(w, "close", nil)
	case <-user.Ctx.Done():
		user.Player.OnGameUpdate = nil
	}
}

var TimeBetweenGameUpdates time.Duration = 100 * time.Millisecond

func playerUpdateHandler(user *users.AuthUser) func(game *blackjack.BlackjackGame) {
	var mutex sync.Mutex

	go sendSSEvent(user.GetUserWriter(), "initial", user)

	return func(game *blackjack.BlackjackGame) {
		mutex.Lock()
		defer mutex.Unlock()
		fmt.Printf("game: %v\n", game)
		sendSSEvent(user.GetUserWriter(), "update", game)
		time.Sleep(TimeBetweenGameUpdates)
	}
}
