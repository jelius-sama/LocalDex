package main

import (
	"LocalDex/api"
	"LocalDex/logger"
	"net/http"
)

var Port = "6969"
var Environment = "development"
var Version string

func main() {
	startServer := func() error {
		logger.Info("Server started on port :" + Port)

		routeHandler := api.HandleRouting()
		return http.ListenAndServe(":"+Port, routeHandler)
	}

	if err := startServer(); err != nil {
		logger.Panic("Could not start the server on port :"+Port, "\n", err)
	}
}
