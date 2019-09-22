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

func getAppManager(aT appType) appDataManager {
	var adm appDataManager
	switch aT {
	case trello:
		adm = new(trelloDataManager)
	case gtask:
		adm = new(gTasksDataManager)
	case etsy:
		adm = new(etsyDataManager)
	case todoist:
		adm = new(todoistDataManager)
	default:
		Fatal("Unknown app type")
	}
	adm.initDataManager()
	return adm
}

func getAppTypeForString(aT string) (appType, error) {
	switch aT {
	case "trello":
		return trello, nil
	case "gtask":
		return gtask, nil
	case "etsy":
		return etsy, nil
	case "todoist":
		return todoist, nil
	default:
		return 0, errors.New("Unknown app type")
	}
}
