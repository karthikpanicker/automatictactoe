package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/dghubble/oauth1"
)

const (
	etsyBaseURL = "https://openapi.etsy.com/v2/"
)

type etsyDataManager struct {
	config        oauth1.Config
	requestSecret string
}

func newEtsyDataManager() *etsyDataManager {
	edm := new(etsyDataManager)
	edm.config = oauth1.Config{
		ConsumerKey:    os.Getenv("ETSY_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("ETSY_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("HOST_URL") + "callback-etsy",
		Endpoint: oauth1.Endpoint{
			AccessTokenURL:  etsyBaseURL + "oauth/access_token",
			AuthorizeURL:    "https://www.etsy.com/oauth/signin?oauth_consumer_key=" + os.Getenv("ETSY_CONSUMER_KEY") + "&service=v2_prod",
			RequestTokenURL: etsyBaseURL + "oauth/request_token?scope=email_r%20listings_r%20transactions_r%20address_r",
		},
	}
	return edm
}

func (edm *etsyDataManager) getAuthorizationURL() string {
	requestToken, requestSecret, err := edm.config.RequestToken()
	if err != nil {
		Error(err)
	}
	edm.requestSecret = requestSecret
	authorizationURL, err := edm.config.AuthorizationURL(requestToken)
	return authorizationURL.String() + "&oauth_token=" + requestToken
}

func (edm *etsyDataManager) getAndPopulateEtsyDetails(r *http.Request) (*userInfo, error) {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	accessToken, accessSecret, err := edm.config.AccessToken(requestToken, edm.requestSecret, verifier)
	if err != nil {
		return nil, err
	}
	userInfo, err := edm.getUserProfileInfo(accessToken, accessSecret)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (edm *etsyDataManager) getUserProfileInfo(accessToken string, accessSecret string) (*userInfo, error) {
	path := etsyBaseURL + "users/__SELF__"
	result := etsyProfileResponse{}
	httpOAuthClient := newHTTPOAuthClient(accessToken, accessSecret, edm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	store := newDataStore()
	info, err := store.getUserInfo(result.Results[0].EtsyUserID)
	if err != nil || info == nil {
		info = &userInfo{
			EmailID: result.Results[0].EmailID,
			UserID:  result.Results[0].EtsyUserID,
			EtsyDetails: etsyDetails{
				EtsyAccessToken:  accessToken,
				EtsyAccessSecret: accessSecret,
			},
			CurrentStep: 1,
		}
	}
	profileDetails, err := edm.getProfileDetails(info.UserID, info)
	if err != nil {
		return nil, err
	}
	info.EtsyDetails.UserName = profileDetails.UserName
	info.EtsyDetails.UserProfileURL = profileDetails.UserProfileURL
	return info, nil
}

func (edm *etsyDataManager) getShops(info *userInfo) error {
	path := etsyBaseURL + "/users/" + strconv.Itoa(info.UserID) + "/shops"
	var result etsyShopResponse
	httpOAuthClient := newHTTPOAuthClient(info.EtsyDetails.EtsyAccessToken,
		info.EtsyDetails.EtsyAccessSecret, edm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	return nil
}

func (edm *etsyDataManager) getTransactionList(info userInfo) (*etsyTransactionResponse, error) {
	path := etsyBaseURL + "shops/" + strconv.Itoa(info.EtsyDetails.UserShopDetails.ShopID) + "/transactions"
	var result etsyTransactionResponse
	httpOAuthClient := newHTTPOAuthClient(info.EtsyDetails.EtsyAccessToken,
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
	httpOAuthClient := newHTTPOAuthClient(info.EtsyDetails.EtsyAccessToken,
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
	httpOAuthClient := newHTTPOAuthClient(info.EtsyDetails.EtsyAccessToken,
		info.EtsyDetails.EtsyAccessSecret, edm.config)
	httpOAuthClient.getMarshalledAPIResponse(path, &result)
	return result.Results[0], nil
}
