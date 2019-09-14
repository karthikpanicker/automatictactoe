package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

func TestGetGTasksAuthorizationURL(t *testing.T) {
	gotenv.OverApply(strings.NewReader("GTASKS_CLIENT_ID=abc"))
	gotenv.OverApply(strings.NewReader("GTASKS_CLIENT_SECRET=abc"))
	gotenv.OverApply(strings.NewReader("HOST_URL=http://localhost:80/"))
	gotenv.OverApply(strings.NewReader("GTASKS_AUTH_URL=http://localhost/o/oauth2/auth"))
	gotenv.OverApply(strings.NewReader("GTASKS_TOKEN_URL=http://localhost/token"))
	gtm := newGTasksDataManager()
	authURL := gtm.getAuthorizationURL()
	assert.Equal(t, "http://localhost/o/oauth2/auth?access_type=offline&"+
		"client_id=abc&"+
		"redirect_uri=http%3A%2F%2Flocalhost%3A80%2Fcallback-google&response_type=code&"+
		"scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Ftasks&state=state-token", authURL)
}

func TestGetAndPopulateGTasksDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/token")
		// Send response to be tested
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"access_token":"ya29.GluCBy-lBDcy9k7vSE1k0Ixh1uBqVC-fYksrKjIcNFjtQxlcVwMTK4jcXqL978bhhShPUU2FZ9_miwB4556d-Da3HheqHxk4FdwYqQ2PO1skjGlp7pvUAogAvbR6",
		"token_type":"Bearer","refresh_token":"1/sStalRv7dB1YLxIDjZVxlhQ205yRhmG7tbWKT5bXMkQ",
		"expiry":"2019-09-13T23:37:45.282532+05:30"}`))
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("GTASKS_AUTH_URL=" + server.URL + "/o/oauth2/auth"))
	gotenv.OverApply(strings.NewReader("GTASKS_TOKEN_URL=" + server.URL + "/token"))
	info := buildDummyUserInfo()
	gtm := newGTasksDataManager()
	err := gtm.getAndPopulateGTasksDetails("hello", info)
	assert.Nil(t, err)
}

func TestGetGTasksService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/token")
		// Send response to be tested
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"access_token":"ya29.GluCBy-lBDcy9k7vSE1k0Ixh1uBqVC-fYksrKjIcNFjtQxlcVwMTK4jcXqL978bhhShPUU2FZ9_miwB4556d-Da3HheqHxk4FdwYqQ2PO1skjGlp7pvUAogAvbR6",
		"token_type":"Bearer","refresh_token":"1/sStalRv7dB1YLxIDjZVxlhQ205yRhmG7tbWKT5bXMkQ",
		"expiry":"2019-09-13T23:37:45.282532+05:30"}`))
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("GTASKS_AUTH_URL=" + server.URL + "/o/oauth2/auth"))
	gotenv.OverApply(strings.NewReader("GTASKS_TOKEN_URL=" + server.URL + "/token"))
	info := buildDummyUserInfo()
	gtm := newGTasksDataManager()
	err := gtm.getAndPopulateGTasksDetails("hello", info)
	assert.Nil(t, err)
	_, err = gtm.getGTasksService(info)
	assert.Nil(t, err)
}

func TestGetTaskLists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.Header().Set("Content-Type", "application/json")
		if req.URL.String() == "/token" {
			rw.Write([]byte(`{"access_token":"ya29.GluCBy-lBDcy9k7vSE1k0Ixh1uBqVC-fYksrKjIcNFjtQxlcVwMTK4jcXqL978bhhShPUU2FZ9_miwB4556d-Da3HheqHxk4FdwYqQ2PO1skjGlp7pvUAogAvbR6",
			"token_type":"Bearer","refresh_token":"1/sStalRv7dB1YLxIDjZVxlhQ205yRhmG7tbWKT5bXMkQ",
			"expiry":"2019-09-13T23:37:45.282532+05:30"}`))
			return
		} else if req.URL.String() == "/users/@me/lists?alt=json&maxResults=10&prettyPrint=false" {
			return
		}
		assert.Fail(t, "Unknow URL to be handled")
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("GTASKS_AUTH_URL=" + server.URL + "/o/oauth2/auth"))
	gotenv.OverApply(strings.NewReader("GTASKS_TOKEN_URL=" + server.URL + "/token"))
	info := buildDummyUserInfo()
	gtm := newGTasksDataManager()
	err := gtm.getAndPopulateGTasksDetails("hello", info)
	assert.Nil(t, err)
	svc, err := gtm.getGTasksService(info)
	svc.BasePath = server.URL
	lists, err := gtm.getTaskLists(info, svc)
	Info(lists)
	assert.Nil(t, err)
}
