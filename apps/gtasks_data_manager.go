package apps

import (
	"encoding/json"
	"errors"
	"etsello/common"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/tasks/v1"
)

const (
	// GTaskGetListsRequest is a key used to specifiy a request to get all lists associated with the
	// gtask account.
	GTaskGetListsRequest   = "gTaskGetListsRequest"
	gTaskGetServiceRequest = "gTaskGetServiceRequest"
	gTaskServiceKey        = "gTaskServiceKey"
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

func (gtm *gTasksDataManager) GetAuthorizationURL() (string, string, error) {
	return gtm.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), "", nil
}

func (gtm *gTasksDataManager) GetAndPopulateAppDetails(info *common.UserInfo, r *http.Request, requestSecret string) error {
	authCode := r.URL.Query().Get("code")
	tok, err := gtm.config.Exchange(context.TODO(), authCode)
	if err != nil {
		common.Error("Unable to retrieve token from web", err)
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

func (gtm *gTasksDataManager) AddItem(info *common.UserInfo, appItemDetails interface{},
	requestParams map[string]interface{}, appItemResponse interface{}) error {
	var srvValue = requestParams[gTaskServiceKey]
	var err error
	if srvValue == nil {
		srvValue, err = gtm.getGTasksService(info)
	}
	if err != nil {
		return err
	}
	srv := srvValue.(*tasks.Service)
	task, err := srv.Tasks.Insert(info.GTasksDetails.SelectedTaskListID, appItemDetails.(*tasks.Task)).Do()
	appItemResponse = task
	return err
}
func (gtm *gTasksDataManager) GetAppData(info *common.UserInfo, requestType string,
	requestParams map[string]interface{}) (interface{}, error) {
	switch requestType {
	case GTaskGetListsRequest:
		var svc *tasks.Service
		param := requestParams[gTaskServiceKey]
		if param != nil {
			svc = param.(*tasks.Service)
		}
		return gtm.getTaskLists(info, svc)
	case gTaskGetServiceRequest:
		return gtm.getGTasksService(info)
	default:
		return nil, errors.New("Unknown request type")
	}
}

func (gtm *gTasksDataManager) getGTasksService(info *common.UserInfo) (*tasks.Service, error) {
	res := oauth2.Token{}
	json.Unmarshal([]byte(info.GTasksDetails.Token), &res)
	client := gtm.config.Client(context.Background(), &res)
	srv, err := tasks.New(client)
	if err != nil {
		common.Error("Unable to retrieve tasks Client", err)
		return nil, err
	}
	return srv, nil
}

func (gtm *gTasksDataManager) getTaskLists(info *common.UserInfo, service *tasks.Service) (*tasks.TaskLists, error) {
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
