package routes

import (
	"net/http"

	"github.com/kraanter/blackjack/pkg/blackjack"
	"github.com/kraanter/blackjack/pkg/restAPI/users"
)

var hitRoute = createRoute("PATCH /hit", hitHandler)
var standRoute = createRoute("PATCH /stand", standHandler)

func hitHandler(w http.ResponseWriter, r *http.Request) {
	user := users.GetUserFromReq(r)
	if user == nil {
		handleUnauthenticated(w)
		return
	}

	_, err := user.Player.Hit()
	if err != nil {
		handleError(w, err.Error(), http.StatusTooEarly)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func standHandler(w http.ResponseWriter, r *http.Request) {
	user := users.GetUserFromReq(r)
	if user == nil {
		handleUnauthenticated(w)
		return
	}

	err := user.Player.Stand()
	if err != nil {
		if err == blackjack.WrongGameStateError {
			handleError(w, err.Error(), http.StatusTooEarly)
		} else if err == blackjack.PlayerNotFoundError {
			handleError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
