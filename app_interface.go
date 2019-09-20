package main

import (
	"errors"
	"net/http"
)

type appType int

const (
	trello  appType = 1
	gtask   appType = 2
	todoist appType = 3
	etsy    appType = 4
)

type appDataManager interface {
	initDataManager()
	getAuthorizationURL() (string, string, error)
	getAndPopulateAppDetails(info *userInfo, r *http.Request, requestSecret string) error
	addItem(info *userInfo, appItemDetails interface{},
		requestParams map[string]interface{}, appItemResponse interface{}) error
	getAppData(info *userInfo, requestType string,
		requestParams map[string]interface{}) (interface{}, error)
}

func getAppManager(aT string) appDataManager {
	adm, _ := getAppMgrInstanceForApp(aT)
	adm.initDataManager()
	return adm
}

func getAppMgrInstanceForApp(aT string) (appDataManager, error) {
	switch aT {
	case "trello":
		return new(trelloDataManager), nil
	case "gtask":
		return new(gTasksDataManager), nil
	case "etsy":
		return new(etsyDataManager), nil
	case "todoist":
		return new(todoistDataManager), nil
	default:
		return nil, errors.New("Unknown app type")
	}
}
