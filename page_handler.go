package main

import (
	"net/http"
	"strconv"
	"time"
)

type pageHandler struct {
	handlerCom    *handlerCommon
	emptyString   string
	requestSecret string
	etsyManager   *etsyDataManager
	trelloManger  *trelloDataManager
	userCache     *userCache
}

func newPageHandler(cache *userCache) *pageHandler {
	ph := new(pageHandler)
	ph.handlerCom = newHandlerCommon()
	ph.etsyManager = newEtsyDataManager()
	ph.trelloManger = newTrelloDataManager()
	ph.userCache = cache
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
	ph.userCache.saveDetailsToCache(info.UserID, *info)
	if err != nil {
		Error("Error in login page template.", err)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   strconv.Itoa(info.UserID),
			Expires: time.Now().AddDate(0, 0, 5),
		})
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", info)
	}
}

func (ph *pageHandler) redirectToTrello(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, ph.trelloManger.getAuthorizationURL(), http.StatusFound)
}

func (ph *pageHandler) trelloAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	userID := ph.handlerCom.ExtractSessionID(r)
	info := ph.userCache.getUserInfo(userID)
	err := ph.trelloManger.getAndPopulateTrelloDetails(r, &info)
	ph.userCache.saveDetailsToCache(info.UserID, info)
	if err != nil {
		Error("Error in login page template.", err)
	} else {
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", info)
	}
}
