package main

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

type httpManager struct {
}

func newHTTPManager() *httpManager {
	apiManager := httpManager{}
	return &apiManager
}

func (wam *httpManager) startServer(dCache dataCache, hostIP string, hostPort int) {
	rm := newRouteManager(dCache)
	address := fmt.Sprintf("%s:%d", hostIP, hostPort)
	Info("Starting http service at: ", address)
	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
		},

		AllowedHeaders: []string{
			"*",
		},
	})
	err := http.ListenAndServe(address, corsOpts.Handler(rm.Router))
	if err != nil {
		Fatal("Error starting http service", err)
	}
}
