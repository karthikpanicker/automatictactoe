package main

import (
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
)

type trelloDataManager struct {
	config               oauth1.Config
	requestSecret        string
	trelloConsumerKey    string
	trelloConsumerSecret string
	trelloBaseURL        string
}

func newTrelloDataManager() *trelloDataManager {
	tdm := new(trelloDataManager)
	tdm.trelloBaseURL = os.Getenv("TRELLO_API_BASE_URL")
	tdm.config = oauth1.Config{
		ConsumerKey:    os.Getenv("TRELLO_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TRELLO_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("HOST_URL") + "callback-trello",
		Endpoint: oauth1.Endpoint{
			AccessTokenURL:  os.Getenv("TRELLO_OAUTH_BASE_URL") + "OAuthGetAccessToken",
			AuthorizeURL:    os.Getenv("TRELLO_OAUTH_BASE_URL") + "OAuthAuthorizeToken",
			RequestTokenURL: os.Getenv("TRELLO_OAUTH_BASE_URL") + "OAuthGetRequestToken",
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

func (tm *trelloDataManager) getAndPopulateTrelloDetails(r *http.Request, info *userInfo) error {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	accessToken, accessSecret, err := tm.config.AccessToken(requestToken, tm.requestSecret, verifier)
	if err != nil {
		return err
	}
	info.TrelloDetails = trelloDetails{
		TrelloAccessToken:  accessToken,
		TrelloAccessSecret: accessSecret,
	}

	trelloBoardIds, _ := tm.getUserBoards(info)
	for _, boardID := range trelloBoardIds {
		boardDetails, _ := tm.getBoardInfo(info, boardID)
		info.TrelloDetails.TrelloBoards = append(info.TrelloDetails.TrelloBoards, *boardDetails)
	}
	info.TrelloDetails.IsLinked = true
	return nil
}

func (tm *trelloDataManager) addCard(info *userInfo, card trelloCardDetails,
	resultCard *trelloCardDetailsResponse) error {
	path := tm.trelloBaseURL + "cards"
	httpOAuthClient := newHTTPOAuthClient(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	err := httpOAuthClient.postResource(path, card, resultCard)
	return err
}

func (tm *trelloDataManager) attachImage(info *userInfo, resultCard *trelloCardDetailsResponse,
	etsyImage etsyImageDetails) error {
	path := tm.trelloBaseURL + "cards/" + resultCard.ID + "/attachments"
	httpOAuthClient := newHTTPOAuthClient(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	err := httpOAuthClient.postResource(path, trelloImageAttachment{
		Name: "Primary Image",
		URL:  etsyImage.ImageURL,
	}, nil)
	return err
}

func (tm *trelloDataManager) getUserBoards(info *userInfo) ([]string, error) {
	path := tm.trelloBaseURL + "members/me"
	var result map[string]interface{}
	httpOAuthClient := newHTTPOAuthClient(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	boardIds := make([]string, 0)
	for _, idBoard := range result["idBoards"].([]interface{}) {
		boardIds = append(boardIds, idBoard.(string))
	}
	return boardIds, nil
}

func (tm *trelloDataManager) getBoardInfo(info *userInfo, boardID string) (*boardDetails, error) {
	path := tm.trelloBaseURL + "boards/" + boardID
	var result boardDetails
	httpOAuthClient := newHTTPOAuthClient(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	return &result, nil
}

func (tm *trelloDataManager) getBoardLists(info *userInfo, boardID string) ([]boardList, error) {
	path := tm.trelloBaseURL + "boards/" + boardID + "/lists"
	var result []boardList
	httpOAuthClient := newHTTPOAuthClient(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	for index, list := range result {
		if list.ID == info.TrelloDetails.SelectedListID {
			result[index].IsSelected = true
		}
	}
	return result, nil
}
