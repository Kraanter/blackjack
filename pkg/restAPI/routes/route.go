package routes

import (
	"net/http"

	"github.com/kraanter/blackjack/pkg/restAPI/middleware"
)

type ApiRoute struct {
	Pattern string
	Handler http.HandlerFunc
	noAuth  bool
}

func (r *ApiRoute) TotalHandler() http.HandlerFunc {
	isAllowed := middleware.AuthMiddleware(r.noAuth)
	return func(res http.ResponseWriter, req *http.Request) {
		isYes := isAllowed(res, req)
		println("allowed", isYes)
		if !isYes {
			http.Error(res, "Unauthorized: Invalid User", http.StatusUnauthorized)
		} else {
			r.Handler(res, req)
		}
	}
}
