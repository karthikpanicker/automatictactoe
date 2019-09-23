package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// HTTPOAuth2Client can be used to fire requests to services with oauth2 header set.
type HTTPOAuth2Client struct {
	accessToken  string
	accessSecret string
	oauthConfig  *oauth2.Config
}

// NewHTTPOauth2Client will create a new instance of oauth2 client
func NewHTTPOauth2Client(config *oauth2.Config) *HTTPOAuth2Client {
	hoc := new(HTTPOAuth2Client)
	hoc.oauthConfig = config
	return hoc
}

// GetMarshalledAPIResponse to fire GET requests on services with oauth2 headers.
func (hoc *HTTPOAuth2Client) GetMarshalledAPIResponse(url string,
	authToken string, responseContainer interface{}) error {
	token := oauth2.Token{}
	json.Unmarshal([]byte(authToken), &token)
	client := hoc.oauthConfig.Client(context.Background(), &token)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	Info(string(body))
	err = json.Unmarshal(body, responseContainer)
	if err != nil {
		return err
	}
	return nil
}

// PostResource to fire POST requests on services with oauth2 headers.
func (hoc *HTTPOAuth2Client) PostResource(authToken string, url string, resource interface{},
	responseContainer interface{}) error {
	token := oauth2.Token{}
	json.Unmarshal([]byte(authToken), &token)
	client := hoc.oauthConfig.Client(context.Background(), &token)
	bodyBytes, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	resp, err := client.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if responseContainer != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, responseContainer)
	}
	return nil
}
