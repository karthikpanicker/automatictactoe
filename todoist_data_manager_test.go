package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
	"golang.org/x/oauth2"
)

func TestGetTodoistAuthorizationURL(t *testing.T) {
	gotenv.OverApply(strings.NewReader("TODOIST_CLIENT_ID=abc"))
	gotenv.OverApply(strings.NewReader("TODOIST_CLIENT_SECRET=abc"))
	gotenv.OverApply(strings.NewReader("HOST_URL=http://localhost:80/"))
	gotenv.OverApply(strings.NewReader("TODOIST_AUTH_URL=http://localhost/oauth/authorize"))
	gotenv.OverApply(strings.NewReader("TODOIST_TOKEN_URL=http://localhost/oauth/access_token"))
	gtm := getAppManager(todoist)
	authURL, _, err := gtm.getAuthorizationURL()
	assert.Nil(t, err, "Error while fetching authorization URL")
	assert.Equal(t, "http://localhost/oauth/authorize?access_type=offline&"+
		"client_id=abc&"+
		"redirect_uri=http%3A%2F%2Flocalhost%3A80%2Fapps%2Ftodoist%2Fcallback&response_type=code&"+
		"scope=data%3Aread_write&state=state-token", authURL)
}

func TestGetAndPopulateTodoistDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/oauth/access_token")
		// Send response to be tested
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"access_token":"ya29.GluCBy-lBDcy9k7vSE1k0Ixh1uBqVC-fYksrKjIcNFjtQxlcVwMTK4jcXqL978bhhShPUU2FZ9_miwB4556d-Da3HheqHxk4FdwYqQ2PO1skjGlp7pvUAogAvbR6",
		"token_type":"Bearer","refresh_token":"",
		"expiry":"2019-09-13T23:37:45.282532+05:30"}`))
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("TODOIST_CLIENT_ID=abc"))
	gotenv.OverApply(strings.NewReader("TODOIST_CLIENT_SECRET=abc"))
	gotenv.OverApply(strings.NewReader("HOST_URL=http://localhost:80/"))
	gotenv.OverApply(strings.NewReader("TODOIST_AUTH_URL=" + server.URL + "/oauth/authorize"))
	gotenv.OverApply(strings.NewReader("TODOIST_TOKEN_URL=" + server.URL + "/oauth/access_token"))
	info := buildDummyUserInfo()
	tdm := getAppManager(todoist)
	err := tdm.getAndPopulateAppDetails(info, httptest.NewRequest("GET",
		"http://localhost/apps/todoist/callback?code=abcd", nil), "abcd")
	assert.Nil(t, err)
	assert.Equal(t, info.TodoistDetails.IsLinked, true)
	assert.NotEqual(t, info.TodoistDetails.Token, "")
}

func TestGetAndPopulateTodoistDetailsWithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/oauth/access_token")
		// Send response to be tested
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"access_token":"ya29.GluCBy-lBDcy9k7vSE1k0Ixh1uBqVC-fYksrKjIcNFjtQxlcVwMTK4jcXqL978bhhShPUU2FZ9_miwB4556d-Da3HheqHxk4FdwYqQ2PO1skjGlp7pvUAogAvbR6",
		"token_type":"Bearer","refresh_token":"",
		"expiry":"2019-09-13T23:37:45.282532+05:30"}`))
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("TODOIST_CLIENT_ID=abc"))
	gotenv.OverApply(strings.NewReader("TODOIST_CLIENT_SECRET=abc"))
	gotenv.OverApply(strings.NewReader("HOST_URL=http://localhost:80/"))
	gotenv.OverApply(strings.NewReader("TODOIST_AUTH_URL=" + server.URL + "/oauth/authorize"))
	// Wrong token URL to trigger an error
	gotenv.OverApply(strings.NewReader("TODOIST_TOKEN_URL=http://localhost/oauth/access_token"))
	info := buildDummyUserInfo()
	tdm := getAppManager(todoist)
	err := tdm.getAndPopulateAppDetails(info, httptest.NewRequest("GET",
		"http://localhost/apps/todoist/callback?code=abcd", nil), "abcd")
	assert.NotEqual(t, nil, err)
}

func TestGetProjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/projects")
		// Send response to be tested
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`[
			{
				"id": 1234,
				"name": "Inbox",
				"comment_count": 10,
				"order": 1
			}
		]`))
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("TODOIST_API_BASE_URL=" + server.URL))
	info := buildDummyUserInfo()
	token := oauth2.Token{
		AccessToken: "abc",
		Expiry:      time.Now().AddDate(0, 0, 1),
	}
	tokenBytes, err := json.Marshal(token)
	info.TodoistDetails.Token = string(tokenBytes)
	tdm := getAppManager(todoist)
	projects, err := tdm.getAppData(info, todoistProjectsRequest, nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(projects.([]todoistProject)))
}

func TestAddTodoistTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/tasks")
		// Send response to be tested
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{
			"comment_count": 0,
			"completed": false,
			"content": "Dress for Maria",
			"due": {
				"date": "2016-09-01",
				"datetime": "2016-09-01T11:00:00Z",
				"string": "2017-07-01 12:00",
				"timezone": "Europe/Lisbon"
			},
			"id": 123,
			"order": 20,
			"priority": 4,
			"project_id": 234,
			"url": "https://todoist.com/showTask?id=123"
		}`))
	}))

	defer server.Close()
	gotenv.OverApply(strings.NewReader("TODOIST_API_BASE_URL=" + server.URL))
	info := buildDummyUserInfo()
	token := oauth2.Token{
		AccessToken: "abc",
		Expiry:      time.Now().AddDate(0, 0, 1),
	}
	tokenBytes, _ := json.Marshal(token)
	info.TodoistDetails.Token = string(tokenBytes)
	tdm := getAppManager(todoist)
	task := &todoistTask{
		Content:   "Dress for Maria",
		ProjectID: 1,
	}
	err := tdm.addItem(info, task, nil, task)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, task.ID)
}

func TestAddTodoistTaskWithUnreachableURL(t *testing.T) {
	gotenv.OverApply(strings.NewReader("TODOIST_API_BASE_URL=http://localhost:9000"))
	info := buildDummyUserInfo()
	token := oauth2.Token{
		AccessToken: "abc",
		Expiry:      time.Now().AddDate(0, 0, 1),
	}
	tokenBytes, _ := json.Marshal(token)
	info.TodoistDetails.Token = string(tokenBytes)
	tdm := getAppManager(todoist)
	task := &todoistTask{
		Content:   "Dress for Maria",
		ProjectID: 1,
	}
	err := tdm.addItem(info, task, nil, task)
	assert.NotNil(t, err)
}
