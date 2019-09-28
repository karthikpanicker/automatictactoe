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
	userID := ph.handlerCom.GetValueForKeyFromSession(r, common.UserID)
	if userID != nil {
		info, err := ph.dCache.GetUserInfo(userID.(int))
		if err == nil {
			ph.handlerCom.rnd.HTML(w, http.StatusOK, "home", info)
			return
		}
	}
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
	userID := ph.handlerCom.GetValueForKeyFromSession(r, common.UserID)
	var info *common.UserInfo
	if userID != nil {
		info, err = ph.dCache.GetUserInfo(userID.(int))
		if err != nil {
			info = &common.UserInfo{}
		}
	} else {
		info = &common.UserInfo{}
	}
	err = aDataMgr.GetAndPopulateAppDetails(info, r, requestSecret)
	if err != nil {
		common.Error("Error processing etsy authorization callback.", err)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", nil)
	} else {
		if userID == nil {
			storedInfo, _ := ph.dCache.GetUserInfo(info.UserID)
			//This users information is already stored, set the etsy information
			if storedInfo != nil {
				storedInfo.EtsyDetails = info.EtsyDetails
				ph.dCache.SaveDetailsToCache(info.UserID, *storedInfo)
			}
		} else {
			ph.dCache.SaveDetailsToCache(info.UserID, *info)
		}
		ph.handlerCom.SaveKeyValueToSession(r, w, common.UserID, info.UserID)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "callbacksuccess", nil)
	}
}

func (ph *pageHandler) showDetails(w http.ResponseWriter, r *http.Request) {
	userID := ph.handlerCom.GetValueForKeyFromSession(r, common.UserID)
	if userID != nil {
		info, err := ph.dCache.GetUserInfo(userID.(int))
		if err == nil {
			ph.handlerCom.rnd.HTML(w, http.StatusOK, "home", info)
			return
		}
	}
	ph.handlerCom.rnd.HTML(w, http.StatusOK, "home", nil)
}

func (ph *pageHandler) logout(w http.ResponseWriter, r *http.Request) {
	ph.handlerCom.DestroySession(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
