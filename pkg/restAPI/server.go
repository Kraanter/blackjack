package restapi

import (
	"embed"
	"net/http"

	"github.com/kraanter/blackjack/pkg/restAPI/routes"
)

//go:embed *
var index embed.FS

func Start() error {
	for _, route := range routes.ApiRoutes {
		http.HandleFunc(route.Pattern, route.GetRouteHandler())
	}

	http.Handle("/", http.FileServerFS(index))

	println("Starting server on port 42069!")
	return http.ListenAndServe(":42069", nil)
}
