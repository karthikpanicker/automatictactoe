package apps

import (
	"errors"
	"etsello/common"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
)

const (
	// TrelloShouldAttachImage is a key to store value to specifiy whether an image should
	// be attached
	TrelloShouldAttachImage = "trelloShouldAttachImage"
	// TrelloBoardListRequest is used to specify a request to get lists assiciated with a board
	TrelloBoardListRequest  = "trelloBoardListRequest"
	trelloUserBoardsRequest = "trelloUserBoardsRequest"
	trelloBoardInfoRequest  = "trelloBoardInfoRequest"
	// TrelloBoardIDKey is the key used to store trello board ID in maps
	TrelloBoardIDKey = "trelloBoardID"
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

func (tm *trelloDataManager) GetAuthorizationURL() (string, string, error) {
	requestToken, requestSecret, err := tm.config.RequestToken()
	if err != nil {
		common.Error(err)
		return "", "", err
	}
	authorizationURL, err := tm.config.AuthorizationURL(requestToken)
	return authorizationURL.String() +
			"&name=" + "Automatictactoe - an etsy sales capture for trello" + "&expiration=never&scope=read,write",
		requestSecret, nil
}

func (tm *trelloDataManager) GetAndPopulateAppDetails(info *common.UserInfo, r *http.Request, requestSecret string) error {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	accessToken, accessSecret, err := tm.config.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		return err
	}
	info.TrelloDetails = common.TrelloDetails{
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

func (tm *trelloDataManager) AddItem(info *common.UserInfo, appItemDetails interface{},
	requestParams map[string]interface{}, appItemResponse interface{}) error {
	path := tm.trelloBaseURL + "cards"
	httpOAuthClient := common.NewHTTPOAuth1Client(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	err := httpOAuthClient.PostResource(path, appItemDetails, appItemResponse)
	if err != nil {
		return err
	}
	if requestParams[TrelloShouldAttachImage] != nil &&
		requestParams[TrelloShouldAttachImage].(bool) {
		// no need to bother if there is an error while attaching the image
		tm.attachImage(info, appItemResponse.(*common.TrelloCardDetailsResponse),
			requestParams[EtsyImageDetailsKey].(common.EtsyImageDetails))
	}
	return nil
}

func (tm *trelloDataManager) GetAppData(info *common.UserInfo, requestType string,
	requestParams map[string]interface{}) (interface{}, error) {
	switch requestType {
	case TrelloBoardListRequest:
		return tm.getBoardLists(info, requestParams[TrelloBoardIDKey].(string))
	case trelloUserBoardsRequest:
		return tm.getUserBoards(info)
	case trelloBoardInfoRequest:
		return tm.getBoardInfo(info, requestParams[TrelloBoardIDKey].(string))
	default:
		return nil, errors.New("unknown request type provided")
	}
}

func (tm *trelloDataManager) attachImage(info *common.UserInfo, resultCard *common.TrelloCardDetailsResponse,
	etsyImage common.EtsyImageDetails) error {
	path := tm.trelloBaseURL + "cards/" + resultCard.ID + "/attachments"
	httpOAuthClient := common.NewHTTPOAuth1Client(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	err := httpOAuthClient.PostResource(path, common.TrelloImageAttachment{
		Name: "Primary Image",
		URL:  etsyImage.ImageURL,
	}, nil)
	return err
}

func (tm *trelloDataManager) getUserBoards(info *common.UserInfo) ([]string, error) {
	path := tm.trelloBaseURL + "members/me"
	var result map[string]interface{}
	httpOAuthClient := common.NewHTTPOAuth1Client(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	err := httpOAuthClient.GetMarshalledAPIResponse(path, &result)
	if err != nil {
		return nil, err
	}
	boardIds := make([]string, 0)
	for _, idBoard := range result["idBoards"].([]interface{}) {
		boardIds = append(boardIds, idBoard.(string))
	}
	return boardIds, nil
}

func (tm *trelloDataManager) getBoardInfo(info *common.UserInfo, boardID string) (*common.BoardDetails, error) {
	path := tm.trelloBaseURL + "boards/" + boardID
	var result common.BoardDetails
	httpOAuthClient := common.NewHTTPOAuth1Client(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	httpOAuthClient.GetMarshalledAPIResponse(path, &result)
	return &result, nil
}

func (tm *trelloDataManager) getBoardLists(info *common.UserInfo, boardID string) ([]common.BoardList, error) {
	path := tm.trelloBaseURL + "boards/" + boardID + "/lists"
	var result []common.BoardList
	httpOAuthClient := common.NewHTTPOAuth1Client(info.TrelloDetails.TrelloAccessToken,
		info.TrelloDetails.TrelloAccessSecret, tm.config)
	err := httpOAuthClient.GetMarshalledAPIResponse(path, &result)
	if err != nil {
		return nil, err
	}
	for index, list := range result {
		if list.ID == info.TrelloDetails.SelectedListID {
			result[index].IsSelected = true
		}
	}
	return result, nil
}
