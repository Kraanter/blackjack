package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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
		writeStructToResponse(w, user.Player.Snapshot(), http.StatusOK)
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
	user.Player.SetGameUpdateHandler(playerUpdateHandler(user))

	select {
	case <-r.Context().Done():
		user.Player.SetGameUpdateHandler(nil)
	case <-user.WriteContext().Done():
		sendSSEvent(w, "close", nil)
	case <-user.Ctx.Done():
		user.Player.SetGameUpdateHandler(nil)
	}
}

var TimeBetweenGameUpdates time.Duration = 100 * time.Millisecond

func playerUpdateHandler(user *users.AuthUser) func(game *blackjack.GameSnapshot) {
	type update struct {
		event string
		data  []byte
	}
	// A full one-deck dealer draw burst fits in the queue instead of being coalesced.
	updates := make(chan update, 52)
	done := user.WriteContext().Done()

	go func() {
		for {
			select {
			case <-done:
				return
			case update := <-updates:
				sendSSEvent(user.GetUserWriter(), update.event, json.RawMessage(update.data))
				time.Sleep(TimeBetweenGameUpdates)
			}
		}
	}()
	updates <- update{"initial", structToString(user.Player.Snapshot())}

	return func(game *blackjack.GameSnapshot) {
		update := update{"update", structToString(game)}
		select {
		case updates <- update:
		default:
			<-updates
			updates <- update
		}
	}
}
