package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/dghubble/oauth1"
)

const (
	etsyBaseURL                  = "https://openapi.etsy.com/v2/"
	etsyTransactionListRequest   = "etsyTransactionList"
	profileDetailsForUserRequest = "profileDetailsForUser"
	etsyImageDetailsRequest      = "etsyImageDetails"
	etsyUserIDKey                = "etsyUserId"
	etsyTranDetailsKey           = "etsyTranDetails"
	etsyImageDetailsKey          = "etsyImageDetails"
)

type etsyDataManager struct {
	config        oauth1.Config
	requestSecret string
}

func (edm *etsyDataManager) initDataManager() {
	edm.config = oauth1.Config{
		ConsumerKey:    os.Getenv("ETSY_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("ETSY_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("HOST_URL") + "apps/etsy/callback",
		Endpoint: oauth1.Endpoint{
			AccessTokenURL:  etsyBaseURL + "oauth/access_token",
			AuthorizeURL:    "https://www.etsy.com/oauth/signin?oauth_consumer_key=" + os.Getenv("ETSY_CONSUMER_KEY") + "&service=v2_prod",
			RequestTokenURL: etsyBaseURL + "oauth/request_token?scope=email_r%20listings_r%20transactions_r%20address_r",
		},
	}
}

func (edm *etsyDataManager) getAuthorizationURL() (string, string, error) {
	requestToken, requestSecret, err := edm.config.RequestToken()
	if err != nil {
		Error(err)
		return "", "", err
	}
	edm.requestSecret = requestSecret
	authorizationURL, err := edm.config.AuthorizationURL(requestToken)
	return authorizationURL.String() + "&oauth_token=" + requestToken, requestSecret, nil
}

func (edm *etsyDataManager) getAndPopulateAppDetails(info *userInfo, r *http.Request, requestSecret string) error {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	accessToken, accessSecret, err := edm.config.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		return err
	}
	err = edm.getUserProfileInfo(info, accessToken, accessSecret)
	return err
}

func (edm *etsyDataManager) addItem(info *userInfo, appItemDetails interface{},
	requestParams map[string]interface{}, appItemResponse interface{}) error {
	return errors.New("Primary app doesnt support add item")
}
func (edm *etsyDataManager) getAppData(info *userInfo, requestType string,
	requestParams map[string]interface{}) (interface{}, error) {
	switch requestType {
	case etsyTransactionListRequest:
		return edm.getTransactionList(info)
	case profileDetailsForUserRequest:
		return edm.getProfileDetails(requestParams[etsyUserIDKey].(int), info)
	case etsyImageDetailsRequest:
		return edm.getImageDetails(info, requestParams[etsyTranDetailsKey].(etsyTransactionDetails))
	default:
		return nil, errors.New("Unknown request type provided")
	}
}

func (edm *etsyDataManager) getUserProfileInfo(info *userInfo, accessToken string, accessSecret string) error {
	path := etsyBaseURL + "users/__SELF__"
	result := etsyProfileResponse{}
	httpOAuthClient := newHTTPOAuth1Client(accessToken, accessSecret, edm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	info.EmailID = result.Results[0].EmailID
	info.UserID = result.Results[0].EtsyUserID
	info.EtsyDetails = etsyDetails{
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

func (edm *etsyDataManager) getShops(info *userInfo) (shopDetails, error) {
	path := etsyBaseURL + "/users/" + strconv.Itoa(info.UserID) + "/shops"
	var result etsyShopResponse
	httpOAuthClient := newHTTPOAuth1Client(info.EtsyDetails.EtsyAccessToken,
		info.EtsyDetails.EtsyAccessSecret, edm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	return result.Results[0], nil
}

func (edm *etsyDataManager) getTransactionList(info *userInfo) (*etsyTransactionResponse, error) {
	path := etsyBaseURL + "shops/" + strconv.Itoa(info.EtsyDetails.UserShopDetails.ShopID) + "/transactions"
	var result etsyTransactionResponse
	httpOAuthClient := newHTTPOAuth1Client(info.EtsyDetails.EtsyAccessToken,
		info.EtsyDetails.EtsyAccessSecret, edm.config)
	err := httpOAuthClient.getMarshalledAPIResponse(path, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (edm *etsyDataManager) getProfileDetails(userid int, info *userInfo) (*etsyUserProfile, error) {
	path := etsyBaseURL + "users/" + strconv.Itoa(userid) + "/profile"
	var result etsyProfileResponse
	httpOAuthClient := newHTTPOAuth1Client(info.EtsyDetails.EtsyAccessToken,
		info.EtsyDetails.EtsyAccessSecret, edm.config)
	err := httpOAuthClient.getMarshalledAPIResponse(path, &result)
	if err != nil {
		return nil, err
	}
	return &result.Results[0], nil
}

func (edm *etsyDataManager) getImageDetails(info *userInfo,
	tranDetails etsyTransactionDetails) (etsyImageDetails, error) {
	path := etsyBaseURL + "listings/" + strconv.Itoa(tranDetails.ListingID) +
		"/images/" + strconv.Itoa(tranDetails.ImageListingID)
	var result etsyImageResponse
	httpOAuthClient := newHTTPOAuth1Client(info.EtsyDetails.EtsyAccessToken,
		info.EtsyDetails.EtsyAccessSecret, edm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	return result.Results[0], nil
}
