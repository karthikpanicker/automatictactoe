package main

import (
	"net/http"
)

type pageHandler struct {
	handlerCom    *handlerCommon
	emptyString   string
	requestSecret string
	etsyManager   *etsyDataManager
	trelloManger  *trelloDataManager
	dCache        dataStore
}

func newPageHandler(cache dataStore) *pageHandler {
	ph := new(pageHandler)
	ph.handlerCom = newHandlerCommon()
	ph.etsyManager = newEtsyDataManager()
	ph.trelloManger = newTrelloDataManager()
	ph.dCache = cache
	return ph
}

func (ph *pageHandler) getLoginPage(w http.ResponseWriter, r *http.Request) {
	ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", nil)
}

func (ph *pageHandler) redirectToEtsy(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, ph.etsyManager.getAuthorizationURL(), http.StatusFound)
}

func (ph *pageHandler) etsyAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	info, err := ph.etsyManager.getAndPopulateEtsyDetails(r)
	ph.dCache.saveDetailsToCache(info.UserID, *info)
	if err != nil {
		Error("Error in login page template.", err)
	} else {
		ph.handlerCom.SaveUserIDInSession(r, w, info.UserID)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", info)
	}
}

func (ph *pageHandler) redirectToTrello(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, ph.trelloManger.getAuthorizationURL(), http.StatusFound)
}

func (ph *pageHandler) trelloAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	userID := ph.handlerCom.GetUserIDFromSession(r)
	info, _ := ph.dCache.getUserInfo(userID)
	err := ph.trelloManger.getAndPopulateTrelloDetails(r, info)
	ph.dCache.saveDetailsToCache(info.UserID, *info)
	if err != nil {
		Error("Error in login page template.", err)
	} else {
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", info)
	}
}
