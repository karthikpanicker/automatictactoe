package main

import (
	"net/http"
)

type pageHandler struct {
	handlerCom     *handlerCommon
	emptyString    string
	requestSecret  string
	etsyManager    *etsyDataManager
	trelloManger   *trelloDataManager
	gTaskManager   *gTasksDataManager
	todoistManager *todoistDataManager
	dCache         dataStore
}

func newPageHandler(cache dataStore) *pageHandler {
	ph := new(pageHandler)
	ph.handlerCom = newHandlerCommon()
	ph.etsyManager = newEtsyDataManager()
	ph.trelloManger = newTrelloDataManager()
	ph.gTaskManager = newGTasksDataManager()
	ph.todoistManager = newToDoistDataManager()
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

func (ph *pageHandler) redirectToEtsy(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, ph.etsyManager.getAuthorizationURL(), http.StatusFound)
}

func (ph *pageHandler) etsyAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	info, err := ph.etsyManager.getAndPopulateEtsyDetails(r)
	info.EtsyDetails.UserShopDetails, _ = ph.etsyManager.getShops(info)
	if err != nil {
		Error("Error processing etsy authorization callback.", err)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", nil)
	} else {
		ph.dCache.saveDetailsToCache(info.UserID, *info)
		ph.handlerCom.SaveUserIDInSession(r, w, info.UserID)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "callbacksuccess", nil)
	}
}

func (ph *pageHandler) showDetails(w http.ResponseWriter, r *http.Request) {
	userID := ph.handlerCom.GetUserIDFromSession(r)
	info, err := ph.dCache.getUserInfo(userID)
	if err != nil {
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "home", userInfo{})
	}
	ph.handlerCom.rnd.HTML(w, http.StatusOK, "home", info)
}

func (ph *pageHandler) redirectToTrello(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, ph.trelloManger.getAuthorizationURL(), http.StatusFound)
}

func (ph *pageHandler) trelloAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	userID := ph.handlerCom.GetUserIDFromSession(r)
	info, _ := ph.dCache.getUserInfo(userID)
	err := ph.trelloManger.getAndPopulateTrelloDetails(r, info)
	if err != nil {
		Error("Error in login page template.", err)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", userInfo{})
	} else {
		ph.dCache.saveDetailsToCache(info.UserID, *info)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "callbacksuccess", nil)
	}
}

func (ph *pageHandler) redirectToGTask(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, ph.gTaskManager.getAuthorizationURL(), http.StatusFound)
}

func (ph *pageHandler) gTasksAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	userID := ph.handlerCom.GetUserIDFromSession(r)
	info, _ := ph.dCache.getUserInfo(userID)
	code := r.URL.Query().Get("code")
	err := ph.gTaskManager.getAndPopulateGTasksDetails(code, info)
	if err != nil {
		Error("Error in login page template.", err)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", userInfo{})
	} else {
		ph.dCache.saveDetailsToCache(info.UserID, *info)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "callbacksuccess", nil)
	}
}

func (ph *pageHandler) redirectToTodoist(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, ph.todoistManager.getAuthorizationURL(), http.StatusFound)
}

func (ph *pageHandler) todoistAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	userID := ph.handlerCom.GetUserIDFromSession(r)
	info, _ := ph.dCache.getUserInfo(userID)
	code := r.URL.Query().Get("code")
	err := ph.todoistManager.getAndPopulateTodoistDetails(code, info)
	if err != nil {
		Error("Error in login page template.", err)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "details", userInfo{})
	} else {
		ph.dCache.saveDetailsToCache(info.UserID, *info)
		ph.handlerCom.rnd.HTML(w, http.StatusOK, "callbacksuccess", nil)
	}
}
