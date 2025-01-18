package routes

import (
	"net/http"

	"github.com/kraanter/blackjack/pkg/restAPI/users"
)

type ApiRoute struct {
	Pattern string
	Handler http.HandlerFunc
	noAuth  bool
}

func (r *ApiRoute) GetRouteHandler() http.HandlerFunc {
	isAllowed := users.AuthMiddleware(r.noAuth)
	return func(res http.ResponseWriter, req *http.Request) {
		hasAccess := isAllowed(res, req)
		if !hasAccess {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
		} else {
			r.Handler(res, req)
		}
	}
}
