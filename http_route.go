package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routeManager struct {
	apiRoutes     []route
	webPageRoutes []route
	pageHandler   *pageHandler
	apiHandler    *apiHandler
	Router        *mux.Router
	dCache        dataCache
}

func newRouteManager(dCache dataCache) *routeManager {
	rm := new(routeManager)
	rm.pageHandler = newPageHandler(dCache)
	rm.apiHandler = newAPIHandler(dCache)
	rm.Router = rm.registerRoutes()
	return rm
}

func (rm *routeManager) registerRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	rm.routeMapping()
	for _, route := range rm.webPageRoutes {
		router.
			Methods(route.Method...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/assets/"))))
	return router
}

func (rm *routeManager) routeMapping() {
	rm.webPageRoutes = []route{
		{
			"Login page",
			[]string{"GET"},
			"/",
			rm.pageHandler.getLoginPage,
		},
		{
			"Etsy login page",
			[]string{"POST"},
			"/authorize-etsy",
			rm.pageHandler.redirectToEtsy,
		},
		{
			"Callback from etsy after successful authorization",
			[]string{"GET"},
			"/callback-etsy",
			rm.pageHandler.etsyAuthorizationCallback,
		},
		{
			"Redirection to trello for authorizaton",
			[]string{"POST"},
			"/authorize-trello",
			rm.pageHandler.redirectToTrello,
		},
		{
			"Callback from trello after successful authorization",
			[]string{"GET"},
			"/callback-trello",
			rm.pageHandler.trelloAuthorizationCallback,
		},
		{
			"Get a list of boards associated with trello",
			[]string{"GET"},
			"/api/trello-boards/{boardId}/lists",
			rm.apiHandler.getBordLists,
		},
		{
			"Save board and list and link it to etsy",
			[]string{"POST"},
			"/api/user-info",
			rm.apiHandler.saveBoardAndList,
		},
	}
}
