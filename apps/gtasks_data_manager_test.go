package apps

import (
	"etsello/common"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
	"google.golang.org/api/tasks/v1"
)

func TestGetGTasksAuthorizationURL(t *testing.T) {
	gotenv.OverApply(strings.NewReader("GTASKS_CLIENT_ID=abc"))
	gotenv.OverApply(strings.NewReader("GTASKS_CLIENT_SECRET=abc"))
	gotenv.OverApply(strings.NewReader("HOST_URL=http://localhost:80/"))
	gotenv.OverApply(strings.NewReader("GTASKS_AUTH_URL=http://localhost/o/oauth2/auth"))
	gotenv.OverApply(strings.NewReader("GTASKS_TOKEN_URL=http://localhost/token"))
	gtm := GetAppManager(Gtask)
	authURL, _, err := gtm.GetAuthorizationURL()
	assert.Nil(t, err)
	assert.Equal(t, "http://localhost/o/oauth2/auth?access_type=offline&"+
		"client_id=abc&"+
		"redirect_uri=http%3A%2F%2Flocalhost%3A80%2Fapps%2Fgtask%2Fcallback&response_type=code&"+
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
	gtm := GetAppManager(Gtask)
	err := gtm.GetAndPopulateAppDetails(info, httptest.NewRequest("GET",
		"http://localhost/apps/gtask/callback?code=abcd", nil), "abcd")
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
	gtm := GetAppManager(Gtask)
	err := gtm.GetAndPopulateAppDetails(info, httptest.NewRequest("GET",
		"http://localhost/apps/gtask/callback?code=abcd", nil), "abcd")
	assert.Nil(t, err)
	_, err = gtm.GetAppData(info, gTaskGetServiceRequest, nil)
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
			rw.Write([]byte(`{
				"kind": "tasks#taskLists",
				"etag": "\"NeaKRry_JhdhtP9cvsTMfcXJ-gY/b_H2Wj5LD31e6JwBbkaJ2LvwrQE\"",
				"items": [
				  {
					"kind": "tasks#taskList",
					"id": "MTAwNDM3NTExNzg4OTEzMzQ4Njk6MDow",
					"title": "My Tasks",
					"updated": "2019-09-10T16:35:12.067Z",
					"selfLink": "https://www.googleapis.com/tasks/v1/users/@me/lists/MTAwNDM3NTExNzg4OTEzMzQ4Njk6MDow"
				  },
				  {
					"kind": "tasks#taskList",
					"id": "MTAwNDM3NTExNzg4OTEzMzQ4Njk6NDg2MjMzNjY1MjAyNTQxNzow",
					"title": "karthik's list",
					"updated": "2019-09-10T17:21:35.444Z",
					"selfLink": "https://www.googleapis.com/tasks/v1/users/@me/lists/MTAwNDM3NTExNzg4OTEzMzQ4Njk6NDg2MjMzNjY1MjAyNTQxNzow"
				  },
				  {
					"kind": "tasks#taskList",
					"id": "dTVMX2k3ZThPRXh2bGJQTw",
					"title": "Etsy List",
					"updated": "2019-09-11T18:57:31.216Z",
					"selfLink": "https://www.googleapis.com/tasks/v1/users/@me/lists/dTVMX2k3ZThPRXh2bGJQTw"
				  }
				]
			  }`))
			return
		}
		assert.Fail(t, "Unknow URL to be handled")
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("GTASKS_AUTH_URL=" + server.URL + "/o/oauth2/auth"))
	gotenv.OverApply(strings.NewReader("GTASKS_TOKEN_URL=" + server.URL + "/token"))
	info := buildDummyUserInfo()
	gtm := GetAppManager(Gtask)
	err := gtm.GetAndPopulateAppDetails(info, httptest.NewRequest("GET",
		"http://localhost/apps/gtask/callback?code=abcd", nil), "abcd")
	assert.Nil(t, err)
	svc, err := gtm.GetAppData(info, gTaskGetServiceRequest, nil)
	svc.(*tasks.Service).BasePath = server.URL
	requestParams := make(map[string]interface{})
	requestParams[gTaskServiceKey] = svc
	lists, err := gtm.GetAppData(info, GTaskGetListsRequest, requestParams)
	assert.Nil(t, err)
	taskLists := lists.([]*common.GTasksListDetails)
	assert.Equal(t, 3, len(taskLists))
}

func TestAddItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		common.Info(req.URL.String())
		rw.Header().Set("Content-Type", "application/json")
		if req.URL.String() == "/token" {
			rw.Write([]byte(`{"access_token":"ya29.GluCBy-lBDcy9k7vSE1k0Ixh1uBqVC-fYksrKjIcNFjtQxlcVwMTK4jcXqL978bhhShPUU2FZ9_miwB4556d-Da3HheqHxk4FdwYqQ2PO1skjGlp7pvUAogAvbR6",
			"token_type":"Bearer","refresh_token":"1/sStalRv7dB1YLxIDjZVxlhQ205yRhmG7tbWKT5bXMkQ",
			"expiry":"2019-09-13T23:37:45.282532+05:30"}`))
			return
		} else if req.URL.String() == "/lists//tasks?alt=json&prettyPrint=false" && req.Method == "POST" {
			rw.Write([]byte(`{"access_token":"ya29.GluCBy-lBDcy9k7vSE1k0Ixh1uBqVC-fYksrKjIcNFjtQxlcVwMTK4jcXqL978bhhShPUU2FZ9_miwB4556d-Da3HheqHxk4FdwYqQ2PO1skjGlp7pvUAogAvbR6",
			"token_type":"Bearer","refresh_token":"1/sStalRv7dB1YLxIDjZVxlhQ205yRhmG7tbWKT5bXMkQ",
			"expiry":"2019-09-13T23:37:45.282532+05:30"}`))
			return
		}
		assert.Fail(t, "Unknow URL to be handled")
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("GTASKS_AUTH_URL=" + server.URL + "/o/oauth2/auth"))
	gotenv.OverApply(strings.NewReader("GTASKS_TOKEN_URL=" + server.URL + "/token"))
	info := buildDummyUserInfo()
	gtm := GetAppManager(Gtask)
	err := gtm.GetAndPopulateAppDetails(info, httptest.NewRequest("GET",
		"http://localhost/apps/gtask/callback?code=abcd", nil), "abcd")
	assert.Nil(t, err)
	svc, err := gtm.GetAppData(info, gTaskGetServiceRequest, nil)
	svc.(*tasks.Service).BasePath = server.URL
	task := &tasks.Task{
		Title: "Hello world",
	}
	var response tasks.Task
	requestParams := make(map[string]interface{})
	requestParams[gTaskServiceKey] = svc
	err = gtm.AddItem(info, task, requestParams, &response)
	assert.Nil(t, err)
	assert.NotNil(t, response)
}
