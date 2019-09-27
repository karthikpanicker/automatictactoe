package apps

import (
	"errors"
	"etsello/common"
	"net/http"
	"os"
	"strconv"

	"github.com/dghubble/oauth1"
)

const (
	// EtsyTransactionListRequest is to specify a request to get all transactions associated with the
	// linked etsy account
	EtsyTransactionListRequest = "etsyTransactionList"
	// ProfileDetailsForUserRequest is to specify a request to get user profile details
	ProfileDetailsForUserRequest = "profileDetailsForUser"
	//EtsyImageDetailsRequest is to specify a request to get image details from etsy
	EtsyImageDetailsRequest = "etsyImageDetails"
	// EtsyUserIDKey is a key used to store etsy user id
	EtsyUserIDKey = "etsyUserId"
	// EtsyTranDetailsKey is a key to store etsy transaction details
	EtsyTranDetailsKey = "etsyTranDetails"
	// EtsyImageDetailsKey is a key to store etsy image details
	EtsyImageDetailsKey = "etsyImageDetails"
)

type etsyDataManager struct {
	config        oauth1.Config
	requestSecret string
	etsyBaseURL   string
}

func (edm *etsyDataManager) initDataManager() {
	edm.etsyBaseURL = os.Getenv("ETSY_TOKEN_BASE_URL")
	edm.config = oauth1.Config{
		ConsumerKey:    os.Getenv("ETSY_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("ETSY_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("HOST_URL") + "apps/etsy/callback",
		Endpoint: oauth1.Endpoint{
			AccessTokenURL:  edm.etsyBaseURL + "oauth/access_token",
			AuthorizeURL:    "https://www.etsy.com/oauth/signin?oauth_consumer_key=" + os.Getenv("ETSY_CONSUMER_KEY") + "&service=v2_prod",
			RequestTokenURL: edm.etsyBaseURL + "oauth/request_token?scope=email_r%20listings_r%20transactions_r%20address_r",
		},
	}
}

func (edm *etsyDataManager) GetAuthorizationURL() (string, string, error) {
	requestToken, requestSecret, err := edm.config.RequestToken()
	if err != nil {
		common.Error(err)
		return "", "", err
	}
	edm.requestSecret = requestSecret
	authorizationURL, err := edm.config.AuthorizationURL(requestToken)
	return authorizationURL.String() + "&oauth_token=" + requestToken, requestSecret, nil
}

func (edm *etsyDataManager) GetAndPopulateAppDetails(info *common.UserInfo, r *http.Request, requestSecret string) error {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	accessToken, accessSecret, err := edm.config.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		return err
	}
	err = edm.getUserProfileInfo(info, accessToken, accessSecret)
	return err
}

func (edm *etsyDataManager) AddItem(info *common.UserInfo, appItemDetails interface{},
	requestParams map[string]interface{}, appItemResponse interface{}) error {
	return errors.New("Primary app doesnt support add item")
}
func (edm *etsyDataManager) GetAppData(info *common.UserInfo, requestType string,
	requestParams map[string]interface{}) (interface{}, error) {
	switch requestType {
	case EtsyTransactionListRequest:
		return edm.getTransactionList(info)
	case ProfileDetailsForUserRequest:
		return edm.getProfileDetails(requestParams[EtsyUserIDKey].(int), info)
	case EtsyImageDetailsRequest:
		return edm.getImageDetails(info, requestParams[EtsyTranDetailsKey].(common.EtsyTransactionDetails))
	default:
		return nil, errors.New("Unknown request type provided")
	}
}

func (edm *etsyDataManager) getUserProfileInfo(info *common.UserInfo, accessToken string, accessSecret string) error {
	path := edm.etsyBaseURL + "users/__SELF__"
	result := common.EtsyProfileResponse{}
	httpOAuthClient := common.NewHTTPOAuth1Client(accessToken, accessSecret, edm.config)
	httpOAuthClient.GetMarshalledAPIResponse(path, &result)
	info.EmailID = result.Results[0].EmailID
	info.UserID = result.Results[0].EtsyUserID
	info.EtsyDetails = common.EtsyDetails{
		EtsyAccessToken:  accessToken,
		EtsyAccessSecret: accessSecret,
	}

	profileDetails, err := edm.getProfileDetails(info.UserID, info)
	if err != nil {
		return err
	}
	info.EtsyDetails.UserName = profileDetails.UserName
	info.EtsyDetails.UserProfileURL = profileDetails.UserProfileURL

	//Get and populate shop details
	info.EtsyDetails.UserShopDetails, err = edm.getShops(info)
	return err
}

func (edm *etsyDataManager) getShops(info *common.UserInfo) (common.ShopDetails, error) {
	path := edm.etsyBaseURL + "users/" + strconv.Itoa(info.UserID) + "/shops"
	var result common.EtsyShopResponse
	httpOAuthClient := common.NewHTTPOAuth1Client(info.EtsyDetails.EtsyAccessToken,
		info.EtsyDetails.EtsyAccessSecret, edm.config)
	httpOAuthClient.GetMarshalledAPIResponse(path, &result)
	return result.Results[0], nil
}

func (edm *etsyDataManager) getTransactionList(info *common.UserInfo) (*common.EtsyTransactionResponse, error) {
	path := edm.etsyBaseURL + "shops/" + strconv.Itoa(info.EtsyDetails.UserShopDetails.ShopID) + "/transactions"
	var result common.EtsyTransactionResponse
	httpOAuthClient := common.NewHTTPOAuth1Client(info.EtsyDetails.EtsyAccessToken,
		info.EtsyDetails.EtsyAccessSecret, edm.config)
	err := httpOAuthClient.GetMarshalledAPIResponse(path, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (edm *etsyDataManager) getProfileDetails(userid int, info *common.UserInfo) (*common.EtsyUserProfile, error) {
	path := edm.etsyBaseURL + "users/" + strconv.Itoa(userid) + "/profile"
	var result common.EtsyProfileResponse
	httpOAuthClient := common.NewHTTPOAuth1Client(info.EtsyDetails.EtsyAccessToken,
		info.EtsyDetails.EtsyAccessSecret, edm.config)
	err := httpOAuthClient.GetMarshalledAPIResponse(path, &result)
	if err != nil {
		return nil, err
	}
	return &result.Results[0], nil
}

func (edm *etsyDataManager) getImageDetails(info *common.UserInfo,
	tranDetails common.EtsyTransactionDetails) (common.EtsyImageDetails, error) {
	path := edm.etsyBaseURL + "listings/" + strconv.Itoa(tranDetails.ListingID) +
		"/images/" + strconv.Itoa(tranDetails.ImageListingID)
	var result common.EtsyImageResponse
	httpOAuthClient := common.NewHTTPOAuth1Client(info.EtsyDetails.EtsyAccessToken,
		info.EtsyDetails.EtsyAccessSecret, edm.config)
	httpOAuthClient.GetMarshalledAPIResponse(path, &result)
	return result.Results[0], nil
}
