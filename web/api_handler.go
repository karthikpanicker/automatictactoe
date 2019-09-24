package web

import (
	"encoding/json"
	"etsello/apps"
	"etsello/common"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type apiHandler struct {
	handlerCom     *handlerCommon
	trelloManager  apps.AppDataManager
	gTManager      apps.AppDataManager
	todoistManager apps.AppDataManager
	dCache         common.DataStore
}

type trialInfo struct {
	SelectedBoardID string `json:"boardId"`
	SelectedListID  string `json:"listId"`
}

func newAPIHandler(cache common.DataStore, templatePattern string) *apiHandler {
	ah := new(apiHandler)
	ah.handlerCom = newHandlerCommon(templatePattern)
	ah.trelloManager = apps.GetAppManager(apps.Trello)
	ah.gTManager = apps.GetAppManager(apps.Gtask)
	ah.todoistManager = apps.GetAppManager(apps.Todoist)
	ah.dCache = cache
	return ah
}

func (ah *apiHandler) getBordLists(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, common.UserID).(int)
	info, _ := ah.dCache.GetUserInfo(userID)
	params := mux.Vars(r)
	boardID := params["boardId"]
	if boardID == "" {
		ah.handlerCom.ProcessErrorMessage(messageInvalidBoardID, w)
		return
	}
	requestParams := make(map[string]interface{})
	requestParams[apps.TrelloBoardIDKey] = boardID
	boardLists, _ := ah.trelloManager.GetAppData(info, apps.TrelloBoardListRequest, requestParams)
	ah.handlerCom.ProcessResponse(boardLists, w)
}

func (ah *apiHandler) saveTrelloConfiguration(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, common.UserID).(int)
	info, _ := ah.dCache.GetUserInfo(userID)
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
	// Mark the selected board in the list of boards as selected and mark others as unselected
	for index, board := range info.TrelloDetails.TrelloBoards {
		if board.ID == info.TrelloDetails.SelectedBoardID {
			info.TrelloDetails.TrelloBoards[index].IsSelected = true
		} else {
			info.TrelloDetails.TrelloBoards[index].IsSelected = false
		}
	}
	info.TrelloDetails.FromDate = ah.setFromDate(info.TrelloDetails.TransactionFilter)
	ah.dCache.SaveDetailsToCache(userID, *info)
	ah.handlerCom.ProcessSuccessMessage(messageSavedTrello, w)
}

func (ah *apiHandler) getGTasksLists(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, common.UserID).(int)
	info, _ := ah.dCache.GetUserInfo(userID)
	tasks, err := ah.gTManager.GetAppData(info, apps.GTaskGetListsRequest, nil)
	if err != nil {
		ah.handlerCom.ProcessErrorMessage(err.Error(), w)
	}
	ah.handlerCom.ProcessResponse(tasks, w)
}

func (ah *apiHandler) saveGTasksConfig(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, common.UserID).(int)
	info, _ := ah.dCache.GetUserInfo(userID)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&info.GTasksDetails)
	if err != nil {
		ah.handlerCom.ProcessErrorMessage(messageInvalidBoardID, w)
		return
	}
	info.GTasksDetails.FromDate = ah.setFromDate(info.GTasksDetails.TransactionFilter)
	ah.dCache.SaveDetailsToCache(userID, *info)
	ah.handlerCom.ProcessSuccessMessage(messageSavedGTasks, w)
}

func (ah *apiHandler) getTodoistProjects(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, common.UserID).(int)
	info, _ := ah.dCache.GetUserInfo(userID)
	projects, err := ah.todoistManager.GetAppData(info, apps.TodoistProjectsRequest, nil)
	if err != nil {
		ah.handlerCom.ProcessErrorMessage(err.Error(), w)
	}
	ah.handlerCom.ProcessResponse(projects, w)
}
func (ah *apiHandler) saveTodoistConfig(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, common.UserID).(int)
	info, _ := ah.dCache.GetUserInfo(userID)
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&info.TodoistDetails)
	info.TodoistDetails.FromDate = ah.setFromDate(info.TodoistDetails.TransactionFilter)
	ah.dCache.SaveDetailsToCache(userID, *info)
	ah.handlerCom.ProcessSuccessMessage(messageSavedGTasks, w)
}

func (ah *apiHandler) setFromDate(filter int) int {
	fromDate := 0
	switch filter {
	case 1:
		fromDate = int(time.Now().Unix())
		break
	case 2:
		fromDate = int(time.Now().AddDate(0, -1, 0).Unix())
		break
	case 3:
		fromDate = int(time.Now().AddDate(0, 0, -7).Unix())
		break
	case 4:
		fromDate = 0
	default:
		fromDate = 0
	}
	return fromDate
}
