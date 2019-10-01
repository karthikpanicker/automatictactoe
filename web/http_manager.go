package web

import (
	"etsello/common"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

// HTTPManager manages http endpoints for the application
type HTTPManager struct {
}

// NewHTTPManager creates a new instance of HTTP manager.
func NewHTTPManager() *HTTPManager {
	apiManager := HTTPManager{}
	return &apiManager
}

// StartServer starts a server on the host ip and port specified
func (wam *HTTPManager) StartServer(dCache common.DataStore, hostIP string, hostPort int) {
	rm := newRouteManager(dCache, "")
	address := fmt.Sprintf("%s:%d", hostIP, hostPort)
	common.Info("Starting http service at: ", address)
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
		common.Fatal("Error starting http service", err)
	}
}
