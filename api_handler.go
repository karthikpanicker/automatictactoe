package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type apiHandler struct {
	handlerCom    *handlerCommon
	trelloManager *trelloDataManager
	userCache     *userCache
}

type trialInfo struct {
	SelectedBoardID string `json:"boardId"`
	SelectedListID  string `json:"listId"`
}

func newAPIHandler(cache *userCache) *apiHandler {
	ah := new(apiHandler)
	ah.handlerCom = newHandlerCommon()
	ah.trelloManager = newTrelloDataManager()
	ah.userCache = cache
	return ah
}

func (ah *apiHandler) getBordLists(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.ExtractSessionID(r)
	info := ah.userCache.getUserInfo(userID)
	params := mux.Vars(r)
	boardID := params["boardId"]
	if boardID == "" {
		ah.handlerCom.ProcessErrorMessage(messageInvalidBoardID, w)
		return
	}
	boardLists, _ := ah.trelloManager.getBoardLists(&info, boardID)
	ah.handlerCom.ProcessResponse(boardLists, w)
}

func (ah *apiHandler) saveBoardAndList(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.ExtractSessionID(r)
	info := ah.userCache.getUserInfo(userID)
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
	ah.userCache.saveDetailsToCache(userID, info)
	ah.handlerCom.ProcessSuccessMessage(messageSavedBoardList, w)
}
