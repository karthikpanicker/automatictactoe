package main

import (
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
)

const (
	trelloBaseURL = "https://api.trello.com/1/"
)

type trelloDataManager struct {
	config               oauth1.Config
	requestSecret        string
	trelloConsumerKey    string
	trelloConsumerSecret string
}

func newTrelloDataManager() *trelloDataManager {
	tdm := new(trelloDataManager)
	tdm.config = oauth1.Config{
		ConsumerKey:    os.Getenv("TRELLO_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TRELLO_CONSUMER_SECRET"),
		CallbackURL:    "http://localhost:8900/callback-trello",
		Endpoint: oauth1.Endpoint{
			AccessTokenURL:  "https://trello.com/1/OAuthGetAccessToken",
			AuthorizeURL:    "https://trello.com/1/OAuthAuthorizeToken",
			RequestTokenURL: "https://trello.com/1/OAuthGetRequestToken",
		},
	}
	return tdm
}

func (tm *trelloDataManager) getAuthorizationURL() string {
	requestToken, requestSecret, err := tm.config.RequestToken()
	if err != nil {
		Error(err)
	}
	tm.requestSecret = requestSecret
	authorizationURL, err := tm.config.AuthorizationURL(requestToken)
	return authorizationURL.String() + "&name=" + "Etsello - an etsy order capture for trello" + "&expiration=never&scope=read,write"
}

func (tm *trelloDataManager) getAndPopulateTrelloDetails(r *http.Request, userInfo *userInfo) error {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	accessToken, accessSecret, err := tm.config.AccessToken(requestToken, tm.requestSecret, verifier)
	if err != nil {
		return err
	}
	userInfo.TrelloDetails = trelloDetails{
		trelloAccessToken:  accessToken,
		trelloAccessSecret: accessSecret,
	}

	trelloBoardIds, _ := tm.getUserBoards(userInfo)
	for _, boardID := range trelloBoardIds {
		boardDetails, _ := tm.getBoardInfo(userInfo, boardID)
		userInfo.TrelloDetails.TrelloBoards = append(userInfo.TrelloDetails.TrelloBoards, *boardDetails)
	}
	userInfo.CurrentStep = 2
	return nil
}

func (tm *trelloDataManager) addCard(info *userInfo, card trelloCardDetails) error {
	path := trelloBaseURL + "cards"
	httpOAuthClient := newHTTPOAuthClient(info.TrelloDetails.trelloAccessToken,
		info.TrelloDetails.trelloAccessSecret, tm.config)
	var result string
	err := httpOAuthClient.postResource(path, card, &result)
	Info(result)
	return err
}

func (tm *trelloDataManager) getUserBoards(info *userInfo) ([]string, error) {
	path := trelloBaseURL + "members/me"
	var result map[string]interface{}
	httpOAuthClient := newHTTPOAuthClient(info.TrelloDetails.trelloAccessToken,
		info.TrelloDetails.trelloAccessSecret, tm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	boardIds := make([]string, 0)
	for _, idBoard := range result["idBoards"].([]interface{}) {
		boardIds = append(boardIds, idBoard.(string))
	}
	return boardIds, nil
}

func (tm *trelloDataManager) getBoardInfo(info *userInfo, boardID string) (*boardDetails, error) {
	path := trelloBaseURL + "boards/" + boardID
	var result boardDetails
	httpOAuthClient := newHTTPOAuthClient(info.TrelloDetails.trelloAccessToken,
		info.TrelloDetails.trelloAccessSecret, tm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	return &result, nil
}

func (tm *trelloDataManager) getBoardLists(info *userInfo, boardID string) ([]boardList, error) {
	path := trelloBaseURL + "boards/" + boardID + "/lists"
	var result []boardList
	httpOAuthClient := newHTTPOAuthClient(info.TrelloDetails.trelloAccessToken,
		info.TrelloDetails.trelloAccessSecret, tm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	return result, nil
}
