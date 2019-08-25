package main

import (
	"encoding/json"
	"io/ioutil"
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
	return authorizationURL.String() + "&name=" + "Etsello - an etsy order capture app for trello" + "&expiration=never"
}

func (tm *trelloDataManager) getAndPopulateTrelloDetails(r *http.Request, userInfo *userInfo) error {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	accessToken, accessSecret, err := tm.config.AccessToken(requestToken, tm.requestSecret, verifier)
	if err != nil {
		return err
	}
	Info(accessToken, accessSecret)
	userInfo.TrelloDetails = trelloDetails{
		trelloAccessToken:  accessToken,
		trelloAccessSecret: accessSecret,
	}

	trelloBoardIds, _ := tm.getUserBoards(userInfo)
	for _, boardID := range trelloBoardIds {
		boardDetails, _ := tm.getBoardInfo(userInfo, boardID)
		userInfo.TrelloDetails.TrelloBoards = append(userInfo.TrelloDetails.TrelloBoards, *boardDetails)
	}
	return nil
}

func (tm *trelloDataManager) getUserBoards(info *userInfo) ([]string, error) {
	token := oauth1.NewToken(info.TrelloDetails.trelloAccessToken, info.TrelloDetails.trelloAccessSecret)
	httpClient := tm.config.Client(oauth1.NoContext, token)
	path := trelloBaseURL + "members/me"
	resp, err := httpClient.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	boardIds := make([]string, 0)
	for _, idBoard := range result["idBoards"].([]interface{}) {
		boardIds = append(boardIds, idBoard.(string))
	}
	return boardIds, nil
}

func (tm *trelloDataManager) getBoardInfo(info *userInfo, boardID string) (*boardDetails, error) {
	token := oauth1.NewToken(info.TrelloDetails.trelloAccessToken, info.TrelloDetails.trelloAccessSecret)
	httpClient := tm.config.Client(oauth1.NoContext, token)
	path := trelloBaseURL + "boards/" + boardID
	resp, err := httpClient.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result boardDetails
	json.Unmarshal(body, &result)
	return &result, nil
}

func (tm *trelloDataManager) getBoardLists(info *userInfo, boardID string) ([]boardList, error) {
	token := oauth1.NewToken(info.TrelloDetails.trelloAccessToken, info.TrelloDetails.trelloAccessSecret)
	httpClient := tm.config.Client(oauth1.NoContext, token)
	path := trelloBaseURL + "boards/" + boardID + "/lists"
	resp, err := httpClient.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result []boardList
	json.Unmarshal(body, &result)
	return result, nil
}
