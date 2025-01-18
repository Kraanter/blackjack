package routes

import (
	"encoding/json"
	"net/http"

	"github.com/kraanter/blackjack/pkg/restAPI/users"
)

var authRoute = ApiRoute{
	Pattern: "GET /auth",
	Handler: authHandler,
}

func authHandler(res http.ResponseWriter, req *http.Request) {
	user := users.GetUserFromReq(req)
	if user == nil {
		res.Write([]byte("No user found"))
		return
	}

	data, err := json.Marshal(user)

	if err != nil {
		println(err.Error())
	}
	res.Write(data)
}
