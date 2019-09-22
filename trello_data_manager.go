package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
)

const (
	trelloShouldAttachImage = "trelloShouldAttachImage"
	trelloBoardListRequest  = "trelloBoardListRequest"
	trelloUserBoardsRequest = "trelloUserBoardsRequest"
	trelloBoardInfoRequest  = "trelloBoardInfoRequest"
	trelloBoardIDKey        = "trelloBoardID"
)

type trelloDataManager struct {
	config        oauth1.Config
	trelloBaseURL string
}

func (tm *trelloDataManager) initDataManager() {
	tm.trelloBaseURL = os.Getenv("TRELLO_API_BASE_URL")
	tm.config = oauth1.Config{
		ConsumerKey:    os.Getenv("TRELLO_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TRELLO_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("HOST_URL") + "apps/trello/callback",
		Endpoint: oauth1.Endpoint{
			AccessTokenURL:  os.Getenv("TRELLO_OAUTH_BASE_URL") + "OAuthGetAccessToken",
			AuthorizeURL:    os.Getenv("TRELLO_OAUTH_BASE_URL") + "OAuthAuthorizeToken",
			RequestTokenURL: os.Getenv("TRELLO_OAUTH_BASE_URL") + "OAuthGetRequestToken",
		},
	}
}

func (tm *trelloDataManager) getAuthorizationURL() (string, string, error) {
	requestToken, requestSecret, err := tm.config.RequestToken()
	if err != nil {
		Error(err)
		return "", "", err
	}
	authorizationURL, err := tm.config.AuthorizationURL(requestToken)
	return authorizationURL.String() +
			"&name=" + "Etsello - an etsy order capture for trello" + "&expiration=never&scope=read,write",
		requestSecret, nil
}

func (tm *trelloDataManager) getAndPopulateAppDetails(info *userInfo, r *http.Request, requestSecret string) error {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	accessToken, accessSecret, err := tm.config.AccessToken(requestToken, requestSecret, verifier)
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

func (tm *trelloDataManager) addItem(info *userInfo, appItemDetails interface{},
	requestParams map[string]interface{}, appItemResponse interface{}) error {
	path := tm.trelloBaseURL + "cards"
	httpOAuthClient := newHTTPOAuth1Client(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	err := httpOAuthClient.postResource(path, appItemDetails, appItemResponse)
	if err != nil {
		return err
	}
	if requestParams[trelloShouldAttachImage] != nil &&
		requestParams[trelloShouldAttachImage].(bool) {
		// no need to bother if there is an error while attaching the image
		tm.attachImage(info, appItemResponse.(*trelloCardDetailsResponse),
			requestParams[etsyImageDetailsKey].(etsyImageDetails))
	}
	return nil
}

func (tm *trelloDataManager) getAppData(info *userInfo, requestType string,
	requestParams map[string]interface{}) (interface{}, error) {
	switch requestType {
	case trelloBoardListRequest:
		return tm.getBoardLists(info, requestParams[trelloBoardIDKey].(string))
	case trelloUserBoardsRequest:
		return tm.getUserBoards(info)
	case trelloBoardInfoRequest:
		return tm.getBoardInfo(info, requestParams[trelloBoardIDKey].(string))
	default:
		return nil, errors.New("Unknown request type provided")
	}
}

func (tm *trelloDataManager) attachImage(info *userInfo, resultCard *trelloCardDetailsResponse,
	etsyImage etsyImageDetails) error {
	path := tm.trelloBaseURL + "cards/" + resultCard.ID + "/attachments"
	httpOAuthClient := newHTTPOAuth1Client(info.TrelloDetails.TrelloAccessToken,
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
	httpOAuthClient := newHTTPOAuth1Client(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	err := httpOAuthClient.getMarshalledAPIResponse(path, &result)
	if err != nil {
		return nil, err
	}
	boardIds := make([]string, 0)
	for _, idBoard := range result["idBoards"].([]interface{}) {
		boardIds = append(boardIds, idBoard.(string))
	}
	return boardIds, nil
}

func (tm *trelloDataManager) getBoardInfo(info *userInfo, boardID string) (*boardDetails, error) {
	path := tm.trelloBaseURL + "boards/" + boardID
	var result boardDetails
	httpOAuthClient := newHTTPOAuth1Client(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	return &result, nil
}

func (tm *trelloDataManager) getBoardLists(info *userInfo, boardID string) ([]boardList, error) {
	path := tm.trelloBaseURL + "boards/" + boardID + "/lists"
	var result []boardList
	httpOAuthClient := newHTTPOAuth1Client(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	for index, list := range result {
		if list.ID == info.TrelloDetails.SelectedListID {
			result[index].IsSelected = true
		}
	}
	return result, nil
}
