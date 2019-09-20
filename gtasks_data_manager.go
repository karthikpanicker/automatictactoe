package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/tasks/v1"
)

const (
	gTaskGetListsRequest = "gTaskGetListsRequest"
)

type gTasksDataManager struct {
	config *oauth2.Config
}

func (gtm *gTasksDataManager) initDataManager() {
	gtm.config = &oauth2.Config{
		ClientID:     os.Getenv("GTASKS_CLIENT_ID"),
		ClientSecret: os.Getenv("GTASKS_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   os.Getenv("GTASKS_AUTH_URL"),
			TokenURL:  os.Getenv("GTASKS_TOKEN_URL"),
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: os.Getenv("HOST_URL") + "apps/gtask/callback",
		Scopes:      []string{tasks.TasksScope},
	}
}

func (gtm *gTasksDataManager) getAuthorizationURL() (string, string, error) {
	return gtm.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), "", nil
}

func (gtm *gTasksDataManager) getAndPopulateAppDetails(info *userInfo, r *http.Request, requestSecret string) error {
	authCode := r.URL.Query().Get("code")
	tok, err := gtm.config.Exchange(context.TODO(), authCode)
	if err != nil {
		Error("Unable to retrieve token from web", err)
		return err
	}
	tokBytes, err := json.Marshal(tok)
	if err != nil {
		return err
	}
	info.GTasksDetails.Token = string(tokBytes)
	info.GTasksDetails.IsLinked = true
	return nil
}

func (gtm *gTasksDataManager) addItem(info *userInfo, appItemDetails interface{},
	requestParams map[string]interface{}, appItemResponse interface{}) error {
	return nil
}
func (gtm *gTasksDataManager) getAppData(info *userInfo, requestType string,
	requestParams map[string]interface{}) (interface{}, error) {
	switch requestType {
	case gTaskGetListsRequest:
		return gtm.getTaskLists(info, nil)
	default:
		return nil, errors.New("Unknown request type")
	}
}

func (gtm *gTasksDataManager) getGTasksService(info *userInfo) (*tasks.Service, error) {
	res := oauth2.Token{}
	json.Unmarshal([]byte(info.GTasksDetails.Token), &res)
	client := gtm.config.Client(context.Background(), &res)
	srv, err := tasks.New(client)
	if err != nil {
		Error("Unable to retrieve tasks Client", err)
		return nil, err
	}
	return srv, nil
}

func (gtm *gTasksDataManager) addToDoItem(info *userInfo, todoItem *tasks.Task, service *tasks.Service) (*tasks.Task, error) {
	var srv *tasks.Service = service
	var err error
	if srv == nil {
		srv, err = gtm.getGTasksService(info)
	}
	if err != nil {
		return nil, err
	}
	task, err := srv.Tasks.Insert(info.GTasksDetails.SelectedTaskListID, todoItem).Do()
	return task, err
}

func (gtm *gTasksDataManager) getTaskLists(info *userInfo, service *tasks.Service) (*tasks.TaskLists, error) {
	var srv *tasks.Service = service
	var err error
	if srv == nil {
		srv, err = gtm.getGTasksService(info)
	}
	if err != nil {
		return nil, err
	}
	lists, err := srv.Tasklists.List().MaxResults(10).Do()
	if err != nil {
		return nil, err
	}
	return lists, nil
}
