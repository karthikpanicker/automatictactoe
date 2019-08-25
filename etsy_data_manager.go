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
	config          oauth1.Config
	requestSecret   string
	httpOAuthClient *httpOAuthClient
}

type etsyProfileResponse struct {
	Count   int           `json:"count"`
	Results []userProfile `json:"results"`
}

type etsyShopResponse struct {
	Count   int           `json:"count"`
	Results []shopDetails `json:"results"`
}

type etsyTransactionResponse struct {
	Count   int                  `json:"count"`
	Results []transactionDetails `json:"results"`
}

type transactionDetails struct {
	ID             int    `json:"transaction_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	BuyerUserID    string `json:"seller_user_id"`
	CreationTime   int64  `json:"creation_tsz"`
	Price          string `json:"price"`
	Currency       string `json:"currency_code"`
	ShippingPrice  string `json:"shipping_cost"`
	ImageListingID string `json:"image_listing_id"`
	EtsyURL        string `json:"url"`
}

type userProfile struct {
	EmailID        string `json:"primary_email"`
	EtsyUserID     int    `json:"user_id"`
	UserProfileURL string `json:"image_url_75x75"`
	UserName       string `json:"login_name"`
}

func newEtsyDataManager() *etsyDataManager {
	edm := new(etsyDataManager)
	edm.config = oauth1.Config{
		ConsumerKey:    os.Getenv("ETSY_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("ETSY_CONSUMER_SECRET"),
		CallbackURL:    "http://localhost:8900/callback-etsy",
		Endpoint: oauth1.Endpoint{
			AccessTokenURL:  etsyBaseURL + "oauth/access_token",
			AuthorizeURL:    "https://www.etsy.com/oauth/signin?oauth_consumer_key=" + os.Getenv("ETSY_CONSUMER_KEY") + "&service=v2_prod",
			RequestTokenURL: etsyBaseURL + "oauth/request_token?scope=email_r%20listings_r%20transactions_r",
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
	edm.httpOAuthClient = newHTTPOAuthClient(accessToken, accessSecret, edm.config)
	if err != nil {
		return nil, err
	}
	userInfo, err := edm.getUserProfileInfo(accessToken, accessSecret)
	edm.getShops(userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (edm *etsyDataManager) getUserProfileInfo(accessToken string, accessSecret string) (*userInfo, error) {
	path := etsyBaseURL + "users/__SELF__"
	result := etsyProfileResponse{}
	edm.httpOAuthClient.getMarshalledAPIResponse(path, &result)
	userInfo := &userInfo{
		EmailID: result.Results[0].EmailID,
		UserID:  result.Results[0].EtsyUserID,
		EtsyDetails: etsyDetails{
			etsyAccessToken:  accessToken,
			etsyAccessSecret: accessSecret,
		},
		CurrentStep: 1,
	}
	edm.getProfileDetails(userInfo)
	return userInfo, nil
}

func (edm *etsyDataManager) getShops(info *userInfo) error {
	path := etsyBaseURL + "/users/" + strconv.Itoa(info.UserID) + "/shops"
	var result etsyShopResponse
	edm.httpOAuthClient.getMarshalledAPIResponse(path, &result)
	info.EtsyDetails.UserShopDetails = result.Results[0]
	return nil
}

func (edm *etsyDataManager) getTransactionList(info userInfo) (*etsyTransactionResponse, error) {
	path := etsyBaseURL + "shops/" + strconv.Itoa(info.EtsyDetails.UserShopDetails.ShopID) + "/transactions"
	var result etsyTransactionResponse
	err := edm.httpOAuthClient.getMarshalledAPIResponse(path, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (edm *etsyDataManager) getProfileDetails(info *userInfo) error {
	path := etsyBaseURL + "users/" + strconv.Itoa(info.UserID) + "/profile"
	var result etsyProfileResponse
	edm.httpOAuthClient.getMarshalledAPIResponse(path, &result)
	info.EtsyDetails.UserProfileURL = result.Results[0].UserProfileURL
	info.EtsyDetails.UserName = result.Results[0].UserName
	return nil
}
