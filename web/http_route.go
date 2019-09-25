package web

import (
	"etsello/common"
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
	dCache        common.DataStore
}

func newRouteManager(dCache common.DataStore) *routeManager {
	rm := new(routeManager)
	rm.pageHandler = newPageHandler(dCache, "")
	rm.apiHandler = newAPIHandler(dCache, "")
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
			"App login page",
			[]string{"GET"},
			"/apps/{appType}/authorize",
			rm.pageHandler.redirectToAppLogin,
		},
		{
			"Callback from apps after authorization",
			[]string{"GET"},
			"/apps/{appType}/callback",
			rm.pageHandler.appAuthorizationCallback,
		},
		{
			"Show details page after successful login",
			[]string{"GET"},
			"/details",
			rm.pageHandler.showDetails,
		},
		{
			"Logout from session",
			[]string{"GET"},
			"/logout",
			rm.pageHandler.logout,
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
			rm.apiHandler.saveTrelloConfiguration,
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
		{
			"Get todoist projects",
			[]string{"GET"},
			"/api/users/{userId}/todoist-projects",
			rm.apiHandler.getTodoistProjects,
		},
		{
			"Get todoist projects",
			[]string{"POST"},
			"/api/users/{userId}/todoist-details",
			rm.apiHandler.saveTodoistConfig,
		},
	}
}
