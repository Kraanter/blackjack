package routes

import (
	"context"
	"net/http"

	"github.com/kraanter/blackjack/pkg/restAPI/users"
)

var getHandRoute = createRoute("GET /hand", getHandHandler)
var buyHandRoute = createRoute("POST /hand", buyHandHandler)

func getHandHandler(res http.ResponseWriter, req *http.Request) {
	user := users.GetUserFromReq(req)
	if user == nil {
		handleUnauthenticated(res)
		return
	}

	writeStructToResponse(res, user.Player.Player.Hands, http.StatusOK)
}

type buyHandBody struct {
	BetAmount uint `json:"bet-amount"`
}

func buyHandHandler(res http.ResponseWriter, req *http.Request) {
	user := users.GetUserFromReq(req)
	if user == nil {
		handleUnauthenticated(res)
		return
	}

	var reqBody buyHandBody
	err := bodyToStruct(req.Body, &reqBody)
	if err != nil {
		handleError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Player.Game.Initialize()

	err = user.Player.Bet(reqBody.BetAmount)
	if err != nil {
		handleError(res, err.Error(), http.StatusPaymentRequired)
		return
	}

	err = user.Player.Game.Start(context.TODO())
	if err != nil {
		handleError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	writeStructToResponse(res, user.Player.Player.Hands, http.StatusCreated)
}
