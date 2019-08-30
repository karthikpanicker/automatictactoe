package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type apiHandler struct {
	handlerCom    *handlerCommon
	trelloManager *trelloDataManager
	dCache        dataStore
}

type trialInfo struct {
	SelectedBoardID string `json:"boardId"`
	SelectedListID  string `json:"listId"`
}

func newAPIHandler(cache dataStore) *apiHandler {
	ah := new(apiHandler)
	ah.handlerCom = newHandlerCommon()
	ah.trelloManager = newTrelloDataManager()
	ah.dCache = cache
	return ah
}

func (ah *apiHandler) getBordLists(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetUserIDFromSession(r)
	info, _ := ah.dCache.getUserInfo(userID)
	params := mux.Vars(r)
	boardID := params["boardId"]
	if boardID == "" {
		ah.handlerCom.ProcessErrorMessage(messageInvalidBoardID, w)
		return
	}
	boardLists, _ := ah.trelloManager.getBoardLists(info, boardID)
	ah.handlerCom.ProcessResponse(boardLists, w)
}

func (ah *apiHandler) saveBoardAndList(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetUserIDFromSession(r)
	info, _ := ah.dCache.getUserInfo(userID)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&info.TrelloDetails)
	if err != nil {
		ah.handlerCom.ProcessErrorMessage(messageInvalidBoardID, w)
		return
	}
	if info.TrelloDetails.SelectedBoardID == "" || info.TrelloDetails.SelectedListID == "" {
		ah.handlerCom.ProcessErrorMessage(messageInvalidBoardList, w)
		return
	}
	ah.dCache.saveDetailsToCache(userID, *info)
	ah.handlerCom.ProcessSuccessMessage(messageSavedBoardList, w)
}
