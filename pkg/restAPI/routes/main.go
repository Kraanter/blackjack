package routes

import (
	"fmt"
	"net/http"
)

type ApiRoute struct {
	Route      string
	Handler    func(http.ResponseWriter, *http.Request)
	MiddleWare *[]func(http.ResponseWriter, *http.Request)
}

var ApiRoutes = []ApiRoute{testRoute}

var testRoute = ApiRoute{
	Route:      "/test",
	Handler:    testRouteHandler,
	MiddleWare: nil,
}

func testRouteHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "test")
}
