package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type pageHandler struct {
	handlerCom    *handlerCommon
	emptyString   string
	requestSecret string
	dCache        dataStore
}

func newPageHandler(cache dataStore) *pageHandler {
	ph := new(pageHandler)
	ph.handlerCom = newHandlerCommon()
	ph.dCache = cache
	return ph
}

func (ph *pageHandler) getLoginPage(w http.ResponseWriter, r *http.Request) {
	ph.handlerCom.rnd.HTML(w, http.StatusOK, "home", nil)
}

func (ph *pageHandler) showPrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	ph.handlerCom.rnd.HTML(w, http.StatusOK, "privacy-policy", nil)
}

func (ph *pageHandler) showTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	ph.handlerCom.rnd.HTML(w, http.StatusOK, "terms-and-conditions", nil)
}

func (ph *pageHandler) redirectToAppLogin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	aTStr := params["appType"]
	aDataMgr := getAppManager(aTStr)
	authURL, requestSecret, _ := aDataMgr.getAuthorizationURL()
	ph.handlerCom.SaveKeyValueToSession(r, w, activeReqSecret, requestSecret)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (ph *pageHandler) appAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	aTStr := params["appType"]
	aDataMgr := getAppManager(aTStr)
	requestSecret := ph.handlerCom.GetValueForKeyFromSession(r, activeReqSecret).(string)
	// Fetch userinfo from db and if not available create a new userinfo instance
	userID := ph.handlerCom.GetValueForKeyFromSession(r, userID).(int)
	info, err := ph.dCache.getUserInfo(userID)
	if info == nil {
		info = new(userInfo)
	}
	err = aDataMgr.getAndPopulateAppDetails(info, r, requestSecret)
	if err != nil {
		Error("Error processing etsy authorization callback.", err)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", nil)
	} else {
		ph.dCache.saveDetailsToCache(info.UserID, *info)
		ph.handlerCom.SaveKeyValueToSession(r, w, userID, info.UserID)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "callbacksuccess", nil)
	}
}

func (ph *pageHandler) showDetails(w http.ResponseWriter, r *http.Request) {
	userID := ph.handlerCom.GetValueForKeyFromSession(r, userID).(int)
	info, err := ph.dCache.getUserInfo(userID)
	if err != nil {
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "home", userInfo{})
	}
	ph.handlerCom.rnd.HTML(w, http.StatusOK, "home", info)
}
