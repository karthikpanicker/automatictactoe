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
	dCache        dataStore
}

func newRouteManager(dCache dataStore) *routeManager {
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
			"Site private policy",
			[]string{"GET"},
			"/privacy-policy",
			rm.pageHandler.showPrivacyPolicy,
		},
		{
			"Terms and Conditions",
			[]string{"GET"},
			"/terms-and-conditions",
			rm.pageHandler.showTermsAndConditions,
		},
		{
			"Login page",
			[]string{"GET"},
			"/",
			rm.pageHandler.getLoginPage,
		},
		{
			"Etsy login page",
			[]string{"GET"},
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
			[]string{"GET"},
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
			"Redirection to google for authorizaton",
			[]string{"GET"},
			"/authorize-gtask",
			rm.pageHandler.redirectToGTask,
		},
		{
			"Callback from google after successful authorization",
			[]string{"GET"},
			"/callback-google",
			rm.pageHandler.gTasksAuthorizationCallback,
		},
		{
			"Show details page after successful login",
			[]string{"GET"},
			"/details",
			rm.pageHandler.showDetails,
		},
		{
			"Get a list of boards associated with trello",
			[]string{"GET"},
			"/api/users/{userId}/trello-boards/{boardId}/lists",
			rm.apiHandler.getBordLists,
		},
		{
			"Save trello configuration details",
			[]string{"POST"},
			"/api/users/{userId}/trello-details",
			rm.apiHandler.saveBoardAndList,
		},
		{
			"Get google tasks list",
			[]string{"GET"},
			"/api/users/{userId}/gtask-lists",
			rm.apiHandler.getGTasksLists,
		},
		{
			"Save google tasks configuration details",
			[]string{"POST"},
			"/api/users/{userId}/gtasks-details",
			rm.apiHandler.saveGTasksConfig,
		},
	}
}
