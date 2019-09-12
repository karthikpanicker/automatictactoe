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
	gtm := newGTasksDataManager()
	authURL := gtm.getAuthorizationURL()
	assert.Equal(t, "https://accounts.google.com/o/oauth2/auth?access_type=offline&"+
		"client_id=abc&"+
		"redirect_uri=http%3A%2F%2Flocalhost%3A80%2Fcallback-google&response_type=code&"+
		"scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Ftasks&state=state-token", authURL)
}

func TestAddToDoItem(t *testing.T) {
	boardID := "513227cad846f3834300649c"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/boards/"+boardID+"/lists")
		// Send response to be tested
		rw.Write([]byte(`[
			{
			  "id": "5d72e23d9e6aa902f8f8a701",
			  "name": "New Order",
			  "closed": false,
			  "idBoard": "5d5b89f68c4e1144e3693792",
			  "pos": 18432,
			  "subscribed": false,
			  "softLimit": null
			},
			{
			  "id": "5d5b89f68c4e1144e3693794",
			  "name": "Vendor Procurement",
			  "closed": false,
			  "idBoard": "5d5b89f68c4e1144e3693792",
			  "pos": 32768,
			  "subscribed": false,
			  "softLimit": null
			},
			{
			  "id": "5d5b89f68c4e1144e3693795",
			  "name": "Handwork",
			  "closed": false,
			  "idBoard": "5d5b89f68c4e1144e3693792",
			  "pos": 49152,
			  "subscribed": false,
			  "softLimit": null
			},
			{
			  "id": "5d5b8a351ff2e25d6fa8d0b7",
			  "name": "Tailoring",
			  "closed": false,
			  "idBoard": "5d5b89f68c4e1144e3693792",
			  "pos": 114688,
			  "subscribed": false,
			  "softLimit": null
			},
			{
			  "id": "5d5b8a43d58da58cd9e3729e",
			  "name": "Shipped",
			  "closed": false,
			  "idBoard": "5d5b89f68c4e1144e3693792",
			  "pos": 180224,
			  "subscribed": false,
			  "softLimit": null
			}
		  ]`))
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("TRELLO_API_BASE_URL=" + server.URL + "/"))
	info := buildDummyUserInfo()
	info.TrelloDetails.SelectedListID = "5d72e23d9e6aa902f8f8a701"
	info.TrelloDetails.SelectedBoardID = "513227cad846f3834300649c"
	tdm := newTrelloDataManager()
	boardLists, err := tdm.getBoardLists(info, boardID)
	assert.Nil(t, err)
	assert.Equal(t, "5d72e23d9e6aa902f8f8a701", boardLists[0].ID)
}
