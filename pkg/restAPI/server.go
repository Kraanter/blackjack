package restapi

import (
	"embed"
	"net/http"
	"strconv"

	"github.com/kraanter/blackjack/pkg/restAPI/routes"
)

//go:embed *
var index embed.FS

func Start(port int) error {
	for _, route := range routes.ApiRoutes {
		http.HandleFunc(route.Pattern, route.GetRouteHandler())
	}

	http.Handle("/", http.FileServerFS(index))

	if port == 0 {
		port = 42069
	}

	println("Starting server on port", port)
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
