package restapi

import (
	"net/http"

	"github.com/kraanter/blackjack/pkg/restAPI/routes"
)

func Start() error {
	for _, route := range routes.ApiRoutes {
		http.HandleFunc(route.Pattern, route.GetRouteHandler())
	}

	return http.ListenAndServe(":42069", nil)
}
