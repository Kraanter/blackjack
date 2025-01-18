package routes

import (
	"net/http"

	"github.com/kraanter/blackjack/pkg/restAPI/users"
)

var buyHandRoute = ApiRoute{
	Pattern: "PUT /hand",
	Handler: buyHandHandler,
}

type buyHandBody struct {
	BetAmount uint `json:"bet-amount"`
}

func buyHandHandler(res http.ResponseWriter, req *http.Request) {
	user := users.GetUserFromReq(req)
	if user == nil {
		res.Write([]byte("No user found"))
		return
	}

	var reqBody buyHandBody
	err := bodyToStruct(req.Body, &reqBody)
	if err != nil {
		handleError(res, err, http.StatusInternalServerError)
	}

	err = user.Player.Bet(reqBody.BetAmount)
	if err != nil {
		handleError(res, err, http.StatusPaymentRequired)
		return
	}
}
