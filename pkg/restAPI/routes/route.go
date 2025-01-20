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

func registerRoute(route *ApiRoute) {
	ApiRoutes = append(ApiRoutes, *route)
}

func createRoute(pattern string, handler http.HandlerFunc) *ApiRoute {
	route := &ApiRoute{pattern, handler, false}

	registerRoute(route)

	return route
}

func createNoAuthRoute(pattern string, handler http.HandlerFunc) *ApiRoute {
	route := &ApiRoute{pattern, handler, true}

	registerRoute(route)

	return route
}

func (r *ApiRoute) GetRouteHandler() http.HandlerFunc {
	isAllowed := users.AuthMiddleware(r.noAuth)
	return func(res http.ResponseWriter, req *http.Request) {
		hasAccess := isAllowed(res, req)
		if !hasAccess {
			handleUnauthenticated(res)
			return
		}

		r.Handler(res, req)
	}
}
