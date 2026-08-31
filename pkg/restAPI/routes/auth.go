package routes

import (
	"net/http"

	"github.com/kraanter/blackjack/pkg/restAPI/users"
)

var authRoute = createRoute("GET /auth", authHandler)

func authHandler(res http.ResponseWriter, req *http.Request) {
	user := users.GetUserFromReq(req)
	if user == nil {
		handleUnauthenticated(res)
		return
	}

	writeStructToResponse(res, user, http.StatusOK)
}
