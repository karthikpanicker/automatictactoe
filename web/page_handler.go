package web

import (
	"etsello/apps"
	"etsello/common"
	"net/http"

	"github.com/gorilla/mux"
)

type pageHandler struct {
	handlerCom    *handlerCommon
	emptyString   string
	requestSecret string
	dCache        common.DataStore
}

func newPageHandler(cache common.DataStore, templatePattern string) *pageHandler {
	ph := new(pageHandler)
	ph.handlerCom = newHandlerCommon(templatePattern)
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
	aT, err := apps.GetAppTypeForString(aTStr)
	if err != nil {
		common.Error("Error redirecting to app login.", err)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", nil)
		return
	}
	aDataMgr := apps.GetAppManager(aT)
	authURL, requestSecret, _ := aDataMgr.GetAuthorizationURL()
	ph.handlerCom.SaveKeyValueToSession(r, w, common.ActiveReqSecret, requestSecret)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (ph *pageHandler) appAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	aTStr := params["appType"]
	aT, err := apps.GetAppTypeForString(aTStr)
	if err != nil {
		common.Error("Error processing etsy authorization callback.", err)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", nil)
		return
	}
	aDataMgr := apps.GetAppManager(aT)
	requestSecret := ph.handlerCom.GetValueForKeyFromSession(r, common.ActiveReqSecret).(string)
	// Fetch userinfo from db and if not available create a new userinfo instance
	userID := ph.handlerCom.GetValueForKeyFromSession(r, common.UserID).(int)
	info, err := ph.dCache.GetUserInfo(userID)
	if info == nil {
		info = new(common.UserInfo)
	}
	err = aDataMgr.GetAndPopulateAppDetails(info, r, requestSecret)
	if err != nil {
		common.Error("Error processing etsy authorization callback.", err)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", nil)
	} else {
		ph.dCache.SaveDetailsToCache(info.UserID, *info)
		ph.handlerCom.SaveKeyValueToSession(r, w, userID, info.UserID)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "callbacksuccess", nil)
	}
}

func (ph *pageHandler) showDetails(w http.ResponseWriter, r *http.Request) {
	userID := ph.handlerCom.GetValueForKeyFromSession(r, common.UserID).(int)
	info, err := ph.dCache.GetUserInfo(userID)
	if err != nil {
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "home", common.UserInfo{})
	}
	ph.handlerCom.rnd.HTML(w, http.StatusOK, "home", info)
}
