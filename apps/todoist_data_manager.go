package apps

import (
	"encoding/json"
	"errors"
	"etsello/common"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	// TodoistProjectsRequest is used to specify a request to get all projects associated with the todoist
	// account
	TodoistProjectsRequest = "todoistProjectsRequest"
)

type todoistDataManager struct {
	config  *oauth2.Config
	baseURL string
}

func (tdm *todoistDataManager) initDataManager() {
	tdm.config = &oauth2.Config{
		ClientID:     os.Getenv("TODOIST_CLIENT_ID"),
		ClientSecret: os.Getenv("TODOIST_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   os.Getenv("TODOIST_AUTH_URL"),
			TokenURL:  os.Getenv("TODOIST_TOKEN_URL"),
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: os.Getenv("HOST_URL") + "apps/todoist/callback",
		Scopes:      []string{"data:read_write"},
	}
	tdm.baseURL = os.Getenv("TODOIST_API_BASE_URL")
}

func (tdm *todoistDataManager) GetAuthorizationURL() (string, string, error) {
	authURL := tdm.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL, "", nil
}

func (tdm *todoistDataManager) GetAndPopulateAppDetails(info *common.UserInfo, r *http.Request, requestSecret string) error {
	authCode := r.URL.Query().Get("code")
	tok, err := tdm.config.Exchange(context.TODO(), authCode)
	if err != nil {
		common.Error("Unable to retrieve token from web", err)
		return err
	}
	tokBytes, err := json.Marshal(tok)
	if err != nil {
		return err
	}
	info.TodoistDetails.Token = string(tokBytes)
	info.TodoistDetails.IsLinked = true
	return nil
}

func (tdm *todoistDataManager) GetAppData(info *common.UserInfo, requestType string,
	requestParams map[string]interface{}) (interface{}, error) {
	switch requestType {
	case TodoistProjectsRequest:
		return tdm.getProjects(info)
	default:
		return nil, errors.New("Unknown request type")
	}
}

func (tdm *todoistDataManager) AddItem(info *common.UserInfo, appItemDetails interface{},
	requestParams map[string]interface{}, appItemResponse interface{}) error {
	path := tdm.baseURL + "/tasks"
	client := common.NewHTTPOauth2Client(tdm.config)
	err := client.PostResource(info.TodoistDetails.Token, path, appItemDetails, appItemResponse)
	return err
}

func (tdm *todoistDataManager) getProjects(info *common.UserInfo) ([]common.TodoistProject, error) {
	path := tdm.baseURL + "/projects"
	result := make([]common.TodoistProject, 0)
	client := common.NewHTTPOauth2Client(tdm.config)
	err := client.GetMarshalledAPIResponse(path, info.TodoistDetails.Token, &result)
	for index, project := range result {
		if project.ID == info.TodoistDetails.SelectedProjectID {
			result[index].IsSelected = true
		}
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}
