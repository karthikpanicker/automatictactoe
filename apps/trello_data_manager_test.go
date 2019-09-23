package apps

import (
	"etsello/common"
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

func TestGetAuthorizationURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/1/OAuthGetRequestToken")
		rw.Write([]byte(`oauth_token=915132350c7d73b3daae0deea59e21d1` +
			`&oauth_token_secret=aa148e1719a0b6f02bca8515b7310283&oauth_callback_confirmed=true`))
	}))
	gotenv.OverApply(strings.NewReader("TRELLO_OAUTH_BASE_URL=" + server.URL + "/1/"))
	tdm := GetAppManager(Trello)
	authURL, _, err := tdm.GetAuthorizationURL()
	assert.Nil(t, err)
	assert.Equal(t, server.URL+"/1/OAuthAuthorizeToken?oauth_token=915132350c7d73b3daae0deea59e21d1"+
		"&name=Etsello - an etsy order capture for trello&expiration=never&scope=read,write",
		authURL)
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
	tdm := GetAppManager(Trello)
	card := common.TrelloCardDetails{
		Name:       "Test trello for Etsy",
		Descripton: "Allows you to sink etsy orders with Trello",
		ListID:     "58e7fee3e06e4001f1cc3658",
	}
	var resultCard common.TrelloCardDetailsResponse
	err := tdm.AddItem(buildDummyUserInfo(), card, nil, &resultCard)
	assert.Nil(t, err)
	assert.Equal(t, "58e800aa9ebaaa01c586f630", resultCard.ID)
}

func TestGetUserBoards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/members/me")
		// Send response to be tested
		rw.Write([]byte(`{
			"id": "513227c9d846f3834300649a",
			"activityBlocked": false,
			"avatarHash": null,
			"avatarUrl": null,
			"bio": "",
			"bioData": null,
			"confirmed": true,
			"fullName": "karthik panicker",
			"idEnterprise": null,
			"idEnterprisesDeactivated": [],
			"idMemberReferrer": null,
			"idPremOrgsAdmin": [],
			"initials": "KP",
			"memberType": "normal",
			"nonPublic": {},
			"nonPublicAvailable": false,
			"products": [],
			"url": "https://trello.com/karthikpanicker",
			"username": "karthikpanicker",
			"status": "disconnected",
			"aaEmail": null,
			"aaId": null,
			"avatarSource": "none",
			"email": null,
			"gravatarHash": "be665905eb6b06e1733b97497d0fbccd",
			"idBoards": [
			  "513227cad846f3834300649c",
			  "5d5b89f68c4e1144e3693792"
			],
			"idOrganizations": [],
			"idEnterprisesAdmin": [],
			"limits": {
			  "boards": {
				"totalPerMember": {
				  "status": "ok",
				  "disableAt": 4500,
				  "warnAt": 4050
				}
			  },
			  "orgs": {
				"totalPerMember": {
				  "status": "ok",
				  "disableAt": 850,
				  "warnAt": 765
				}
			  }
			},
			"loginTypes": null,
			"marketingOptIn": {
			  "optedIn": false,
			  "date": "2018-04-25T21:30:44.432Z"
			},
			"messagesDismissed": [
			  {
				"_id": "5d55b16ccfa45a38da6a4d0f",
				"name": "team-join-cta-banner-513227e90c7d58f12d005080",
				"count": 4,
				"lastDismissed": "2019-08-15T19:24:28.583Z"
			  }
			],
			"oneTimeMessagesDismissed": [
			  "nux-boards-page-ORG_TYPE_FREE-513227e90c7d58f12d005080-banner",
			  "close-menu-of-first-board",
			  "board-background-prompt"
			],
			"prefs": {
			  "privacy": {
				"fullName": "public",
				"avatar": "public"
			  },
			  "sendSummaries": true,
			  "minutesBetweenSummaries": 60,
			  "minutesBeforeDeadlineToNotify": 1440,
			  "colorBlind": false,
			  "locale": ""
			},
			"trophies": [],
			"uploadedAvatarHash": null,
			"uploadedAvatarUrl": null,
			"premiumFeatures": [],
			"isAaMastered": false,
			"ixUpdate": "20",
			"idBoardsPinned": null
		  }`))
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("TRELLO_API_BASE_URL=" + server.URL + "/"))
	tdm := GetAppManager(Trello)
	response, err := tdm.GetAppData(buildDummyUserInfo(), trelloUserBoardsRequest, nil)
	assert.Nil(t, err)
	assert.Equal(t, "513227cad846f3834300649c", response.([]string)[0])
}

func TestGetBoardLists(t *testing.T) {
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
	tdm := GetAppManager(Trello)
	requestParams := make(map[string]interface{})
	requestParams[TrelloBoardIDKey] = boardID
	boardLists, err := tdm.GetAppData(info, TrelloBoardListRequest, requestParams)
	assert.Nil(t, err)
	assert.Equal(t, "5d72e23d9e6aa902f8f8a701", boardLists.([]common.BoardList)[0].ID)
}

func TestGetBoardInfo(t *testing.T) {
	boardID := "513227cad846f3834300649c"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/boards/"+boardID)
		// Send response to be tested
		rw.Write([]byte(`{
			"id":"513227cad846f3834300649c",
			"name":"Welcome Board",
			"desc":"",
			"descData":null,
			"closed":false,
			"idOrganization":null,
			"pinned":false,
			"url":"https://trello.com/b/ASOx3Wag/welcome-board",
			"shortUrl":"https://trello.com/b/ASOx3Wag",
			"prefs":{
			   "permissionLevel":"private",
			   "hideVotes":false,
			   "voting":"members",
			   "comments":"members",
			   "invitations":"members",
			   "selfJoin":false,
			   "cardCovers":true,
			   "isTemplate":false,
			   "calendarFeedEnabled":false,
			   "background":"blue",
			   "backgroundImage":null,
			   "backgroundImageScaled":null,
			   "backgroundTile":false,
			   "backgroundBrightness":"dark",
			   "backgroundColor":"#0079BF",
			   "backgroundBottomColor":"#0079BF",
			   "backgroundTopColor":"#0079BF",
			   "canBePublic":true,
			   "canBeEnterprise":true,
			   "canBeOrg":true,
			   "canBePrivate":true,
			   "canInvite":true
			},
			"labelNames":{
			   "green":"",
			   "yellow":"",
			   "orange":"",
			   "red":"",
			   "purple":"",
			   "blue":"",
			   "sky":"",
			   "lime":"",
			   "pink":"",
			   "black":""
			}
		 }`))
	}))
	defer server.Close()
	gotenv.OverApply(strings.NewReader("TRELLO_API_BASE_URL=" + server.URL + "/"))
	info := buildDummyUserInfo()
	tdm := GetAppManager(Trello)
	requestParams := make(map[string]interface{})
	requestParams[TrelloBoardIDKey] = boardID
	boardInfo, err := tdm.GetAppData(info, trelloBoardInfoRequest, requestParams)
	assert.Nil(t, err)
	assert.Equal(t, "513227cad846f3834300649c", boardInfo.(*common.BoardDetails).ID)
}

func buildDummyUserInfo() *common.UserInfo {
	info := &common.UserInfo{
		UserID:  1234,
		EmailID: "karthik.panicker@gmail.com",
	}
	return info
}
