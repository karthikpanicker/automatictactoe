package apps

import (
	"etsello/common"
	"net/http"
)

// AppType custom type represents various apps supported by the application
type AppType int

const (
	// Trello is a constant to represent trello app
	Trello AppType = 1
	// Gtask is a constant to represent google tasks app
	Gtask AppType = 2
	// Todoist is a constant to represent todoist app
	Todoist AppType = 3
	// Etsy is a constant to represent etsy app
	Etsy AppType = 4
	// DefaultAppType that would be used for unsupported app types
	DefaultAppType AppType = 5
)

// AppDataManager is the interface that abstracts app implementations.
type AppDataManager interface {
	initDataManager()
	GetAuthorizationURL() (string, string, error)
	GetAndPopulateAppDetails(info *common.UserInfo, r *http.Request, requestSecret string) error
	AddItem(info *common.UserInfo, appItemDetails interface{},
		requestParams map[string]interface{}, appItemResponse interface{}) error
	GetAppData(info *common.UserInfo, requestType string,
		requestParams map[string]interface{}) (interface{}, error)
}

type defaultAppManager struct {
}

// GetAppManager is a function to get the app manger implementation for a given AppType value.
func GetAppManager(aT AppType) AppDataManager {
	var adm AppDataManager
	switch aT {
	case Trello:
		adm = new(trelloDataManager)
	case Gtask:
		adm = new(gTasksDataManager)
	case Etsy:
		adm = new(etsyDataManager)
	case Todoist:
		adm = new(todoistDataManager)
	default:
		adm = new(defaultAppManager)
	}
	adm.initDataManager()
	return adm
}

// GetAppTypeForString is a function to get the AppType value for given string
func GetAppTypeForString(aT string) (AppType, error) {
	switch aT {
	case "trello":
		return Trello, nil
	case "gtask":
		return Gtask, nil
	case "etsy":
		return Etsy, nil
	case "todoist":
		return Todoist, nil
	default:
		return DefaultAppType, nil
	}
}

func (dam *defaultAppManager) initDataManager() {

}
func (dam *defaultAppManager) GetAuthorizationURL() (string, string, error) {
	return "", "", nil
}

func (dam *defaultAppManager) GetAndPopulateAppDetails(info *common.UserInfo, r *http.Request, requestSecret string) error {
	return nil
}
func (dam *defaultAppManager) AddItem(info *common.UserInfo, appItemDetails interface{},
	requestParams map[string]interface{}, appItemResponse interface{}) error {
	return nil
}
func (dam *defaultAppManager) GetAppData(info *common.UserInfo, requestType string,
	requestParams map[string]interface{}) (interface{}, error) {
	return info, nil
}
