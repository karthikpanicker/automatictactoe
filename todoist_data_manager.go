package main

import (
	"encoding/json"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type todoistDataManager struct {
	config *oauth2.Config
}

func newToDoistDataManager() *todoistDataManager {
	tdm := new(todoistDataManager)
	tdm.config = &oauth2.Config{
		ClientID:     os.Getenv("TODOIST_CLIENT_ID"),
		ClientSecret: os.Getenv("TODOIST_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   os.Getenv("TODOIST_AUTH_URL"),
			TokenURL:  os.Getenv("TODOIST_TOKEN_URL"),
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: os.Getenv("HOST_URL") + "callback-todoist",
		Scopes:      []string{"data:read_write"},
	}
	return tdm
}

func (tdm *todoistDataManager) getAuthorizationURL() string {
	authURL := tdm.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	Info(authURL)
	return authURL
}

func (tdm *todoistDataManager) getAndPopulateTodoistDetails(authCode string, info *userInfo) error {
	tok, err := tdm.config.Exchange(context.TODO(), authCode)
	if err != nil {
		Error("Unable to retrieve token from web", err)
		return err
	}
	tokBytes, _ := json.Marshal(tok)
	info.TodoistDetails.Token = string(tokBytes)
	info.TodoistDetails.IsLinked = true
	return nil
}
