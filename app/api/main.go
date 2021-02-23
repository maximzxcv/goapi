package main

import (
	"goapi/app/api/handlers"
	"goapi/app/api/middle"
	"goapi/bal"
	"net/http"
)

const (
	serverAddress = "0.0.0.0:8080"
)

func main() {
	logg := bal.NewLogg()
	mux := http.NewServeMux()

	mux.HandleFunc("/users", handlers.GetUsers)

	loggMiddle := middle.LoggMiddle(logg)
	configuredMux := loggMiddle(mux)

	api := http.Server{
		Addr:    serverAddress,
		Handler: configuredMux,
	}

	logg.Debug("API is running on %v", serverAddress)
	if err := api.ListenAndServe(); err != nil {
	}

}
