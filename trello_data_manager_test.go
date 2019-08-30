package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
func TestAddCard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/cards")
		// Send response to be tested
		rw.Write([]byte(`{
			"id": "58e800aa9ebaaa01c586f630",
			"badges": {
			  "votes": 0,
			  "viewingMemberVoted": false,
			  "subscribed": false,
			  "fogbugz": "",
			  "checkItems": 0,
			  "checkItemsChecked": 0,
			  "comments": 0,
			  "attachments": 1,
			  "description": true,
			  "due": null,
			  "dueComplete": false
			},
			"checkItemStates": [
			  
			],
			"closed": false,
			"dateLastActivity": "2017-04-07T21:26:00.365Z",
			"desc": "Allows you to sink etsy orders with Trello",
			"descData": {
			  "emoji": {
				
			  }
			},
			"due": null,
			"dueComplete": false,
			"idBoard": "4d5ea62fd76aa1136000000c",
			"idChecklists": [
			  
			],
			"idLabels": [
			  
			],
			"idList": "58e7fee3e06e4001f1cc3658",
			"idMembers": [
			  
			],
			"idShort": 1579,
			"idAttachmentCover": "58e801c30f51ca01a6dd6745",
			"manualCoverAttachment": false,
			"labels": [
			  
			],
			"name": "Test trello for Etsy",
			"pos": 131071,
			"shortUrl": "https://trello.com/c/SDmaAwz9",
			"url": "https://trello.com/c/SDmaAwz9/1579-trello-alerts-for-slack"
		  }`))
	}))
	defer server.Close()
	gotenv.Apply(strings.NewReader("TRELLO_API_BASE_URL=" + server.URL + "/"))
	tdm := newTrelloDataManager()
	card := trelloCardDetails{
		Name:       "Test trello for Etsy",
		Descripton: "Allows you to sink etsy orders with Trello",
		ListID:     "58e7fee3e06e4001f1cc3658",
	}
	var resultCard trelloCardDetailsResponse
	err := tdm.addCard(buildDummyUserInfo(), card, &resultCard)
	assert.Nil(t, err)
	assert.Equal(t, "58e800aa9ebaaa01c586f630", resultCard.ID)
}

func buildDummyUserInfo() *userInfo {
	info := &userInfo{
		UserID:  1234,
		EmailID: "karthik.panicker@gmail.com",
	}
	info.TrelloDetails = trelloDetails{
		TrelloAccessToken:  "token",
		TrelloAccessSecret: "secret",
		SelectedBoardID:    "boardId",
		SelectedListID:     "listId",
	}
	return info
}
