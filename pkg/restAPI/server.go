package restapi

import (
	"embed"
	"flag"
	"net/http"
	"strconv"

	"github.com/kraanter/blackjack/pkg/restAPI/routes"
)

//go:embed *
var index embed.FS

func Start() error {
	for _, route := range routes.ApiRoutes {
		http.HandleFunc(route.Pattern, route.GetRouteHandler())
	}

	portPtr := flag.Int("port", 42069, "The port to run the server on")

	http.Handle("/", http.FileServerFS(index))

	println("Starting server on port", *portPtr)
	return http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil)
}
