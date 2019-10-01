package web

import "testing"

func TestRegisterRoutes(t *testing.T) {
	dStore := NewTestDataStore()
	rm := newRouteManager(dStore, "../templates/*.html")
	rm.registerRoutes()

}
