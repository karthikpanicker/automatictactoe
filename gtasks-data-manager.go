package main

import (
	"encoding/json"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/tasks/v1"
)

type gTasksDataManager struct {
	config        *oauth2.Config
	requestSecret string
}

func newGTasksDataManager() *gTasksDataManager {
	gtm := new(gTasksDataManager)
	gtm.config = &oauth2.Config{
		ClientID:     "598823989725-il1f05povf8cfo2osgvu57oij5950tv0.apps.googleusercontent.com",
		ClientSecret: "flux3hoFfgbXwNrAEO3jnID-",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://accounts.google.com/o/oauth2/auth",
			TokenURL:  "https://oauth2.googleapis.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: "http://localhost/callback-google",
		Scopes:      []string{tasks.TasksScope},
	}
	return gtm
}

func (gtm *gTasksDataManager) getAuthorizationURL() string {
	return gtm.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (gtm *gTasksDataManager) getAndPopulateGTasksDetails(authCode string, info *userInfo) error {
	tok, err := gtm.config.Exchange(context.TODO(), authCode)
	if err != nil {
		Error("Unable to retrieve token from web: %v", err)
		return err
	}
	tokBytes, _ := json.Marshal(tok)
	info.GTasksDetails.Token = string(tokBytes)
	info.GTasksDetails.IsLinked = true
	return nil
}

func (gtm *gTasksDataManager) getGTasksService(info *userInfo) (*tasks.Service, error) {
	res := oauth2.Token{}
	json.Unmarshal([]byte(info.GTasksDetails.Token), &res)
	client := gtm.config.Client(context.Background(), &res)
	srv, err := tasks.New(client)
	if err != nil {
		Error("Unable to retrieve tasks Client %v", err)
		return nil, err
	}
	return srv, nil
}

func (gtm *gTasksDataManager) addToDoItem(info *userInfo, todoItem *tasks.Task) (*tasks.Task, error) {
	srv, err := gtm.getGTasksService(info)
	if err != nil {
		return nil, err
	}
	task, err := srv.Tasks.Insert(info.GTasksDetails.SelectedTaskListID, todoItem).Do()
	return task, err
}

func (gtm *gTasksDataManager) getTaskLists(info *userInfo) (*tasks.TaskLists, error) {
	srv, err := gtm.getGTasksService(info)
	if err != nil {
		return nil, err
	}
	lists, err := srv.Tasklists.List().MaxResults(10).Do()
	if err != nil {
		return nil, err
	}
	return lists, nil
}
