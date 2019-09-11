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
	gotenv.OverApply(strings.NewReader("TRELLO_API_BASE_URL=" + server.URL + "/"))
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

func TestGetUserBoards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/members/me")
		// Send response to be tested
		rw.Write([]byte(`{"id":"513227c9d846f3834300649a",
		"activityBlocked":false,
		"avatarHash":null,
		"avatarUrl":null,
		"bio":"",
		"bioData":null,
		"confirmed":true,
		"fullName":"karthik panicker",
		"idEnterprise":null,
		"idEnterprisesDeactivated":[],
		"idMemberReferrer":null,
		"idPremOrgsAdmin":[],
		"initials":"KP",
		"memberType":"normal",
		"nonPublic":{},
		"nonPublicAvailable":false,
		"products":[],
		"url":"https://trello.com/karthikpanicker",
		"username":"karthikpanicker",
		"status":"disconnected",
		"aaEmail":null,
		"aaId":null,
		"avatarSource":
		"none",
		"email":null,
		"gravatarHash":"be665905eb6b06e1733b97497d0fbccd",
		"idBoards":["513227cad846f3834300649c","5d5b89f68c4e1144e3693792"],
		"idOrganizations":[],"idEnterprisesAdmin":[],
		"limits":{"boards":{"totalPerMember":{"status":"ok","disableAt":4500,"warnAt":4050}},
		"orgs":{"totalPerMember":{"status":"ok","disableAt":850,"warnAt":765}}},
		"loginTypes":null,"marketingOptIn":{"optedIn":false,"date":"2018-04-25T21:30:44.432Z"},
		"messagesDismissed":[{"_id":"5d55b16ccfa45a38da6a4d0f",
		"name":"team-join-cta-banner-513227e90c7d58f12d005080",
		"count":4,"lastDismissed":"2019-08-15T19:24:28.583Z"}],
		"oneTimeMessagesDismissed":["nux-boards-page-ORG_TYPE_FREE-513227e90c7d58f12d005080-banner","close-menu-of-first-board",
		"board-background-prompt"],"prefs":{"privacy":{"fullName":"public","avatar":"public"},
		"sendSummaries":true,"minutesBetweenSummaries":60,"minutesBeforeDeadlineToNotify":1440,
		"colorBlind":false,"locale":""},"trophies":[],"uploadedAvatarHash":null,
		"uploadedAvatarUrl":null,"premiumFeatures":[],"isAaMastered":false,
		"ixUpdate":"20","idBoardsPinned":null}`))
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("TRELLO_API_BASE_URL=" + server.URL + "/"))
	Info(server.URL)
	tdm := newTrelloDataManager()
	boardArray, err := tdm.getUserBoards(buildDummyUserInfo())
	assert.Nil(t, err)
	assert.Equal(t, "513227cad846f3834300649c", boardArray[0])
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
