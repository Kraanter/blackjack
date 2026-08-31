package restapi

import (
	"embed"
	"net/http"
	"strconv"
	"time"

	"github.com/kraanter/blackjack/pkg/restAPI/routes"
)

//go:embed *
var index embed.FS

type ServerSettings struct {
	Port uint

	// If 2 consequtive game updates happen in short succession the updates will be sent with this duration in between
	// If update 1 and 2 happen 10 miliseconds apart and this is set to 100 miliseconds
	// The updates will be sent as following: 1 -> wait 100 milis -> 2
	TimeBetweenGameUpdates time.Duration
}

func CreateDefaultServerSettings() *ServerSettings {
	return &ServerSettings{
		Port: 42069,

		TimeBetweenGameUpdates: routes.TimeBetweenGameUpdates,
	}
}

func Start(settings *ServerSettings) error {
	for _, route := range routes.ApiRoutes {
		http.HandleFunc(route.Pattern, route.GetRouteHandler())
	}

	http.Handle("/", http.FileServerFS(index))

	if settings.Port == 0 {
		settings.Port = 42069
	}

	routes.TimeBetweenGameUpdates = settings.TimeBetweenGameUpdates

	println("Starting server on port", settings.Port)
	return http.ListenAndServe(":"+strconv.Itoa(int(settings.Port)), nil)
}
