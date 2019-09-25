package web

import (
	"etsello/common"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetLoginPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	dStore := NewTestDataStore()
	handler := http.HandlerFunc(newPageHandler(dStore, "../templates/*.html").getLoginPage)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRedirectToAppLogin(t *testing.T) {
	req, err := http.NewRequest("GET", "/apps/etsy/authorize", nil)
	if err != nil {
		t.Fatal(err)
	}
	requestParams := make(map[string]string)
	requestParams["appType"] = "etsy"
	req = mux.SetURLVars(req, requestParams)
	rr := httptest.NewRecorder()
	dStore := NewTestDataStore()
	handler := http.HandlerFunc(newPageHandler(dStore, "../templates/*.html").redirectToAppLogin)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAppAuthorizationCallback(t *testing.T) {
	req, err := http.NewRequest("GET", "/apps/etsy/callback", nil)
	if err != nil {
		t.Fatal(err)
	}
	requestParams := make(map[string]string)
	requestParams["appType"] = "etsy"
	req = mux.SetURLVars(req, requestParams)
	rr := httptest.NewRecorder()
	dStore := NewTestDataStore()
	ph := newPageHandler(dStore, "../templates/*.html")
	ph.handlerCom.SaveKeyValueToSession(req, rr, common.ActiveReqSecret, "abcd")
	ph.handlerCom.SaveKeyValueToSession(req, rr, common.UserID, 12345)
	handler := http.HandlerFunc(ph.appAuthorizationCallback)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAppAuthorizationLogout(t *testing.T) {
	req, err := http.NewRequest("GET", "/apps/etsy/callback", nil)
	if err != nil {
		t.Fatal(err)
	}
	requestParams := make(map[string]string)
	requestParams["appType"] = "etsy"
	req = mux.SetURLVars(req, requestParams)
	rr := httptest.NewRecorder()
	dStore := NewTestDataStore()
	ph := newPageHandler(dStore, "../templates/*.html")
	handler := http.HandlerFunc(ph.logout)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// A dummy implementation of datastore be be used for testing
type TestDataCache struct {
}

func (tdc *TestDataCache) SaveDetailsToCache(userID int, userInfo common.UserInfo) {

}

func (tdc *TestDataCache) GetUserInfo(userID int) (*common.UserInfo, error) {
	return &common.UserInfo{}, nil
}

func (tdc *TestDataCache) GetUserMap() map[int]common.UserInfo {
	return make(map[int]common.UserInfo)
}

func (tdc *TestDataCache) DisconnectCache() {

}

func NewTestDataStore() *TestDataCache {
	tds := new(TestDataCache)
	return tds
}
