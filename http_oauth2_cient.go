package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type httpOAuth2Client struct {
	accessToken  string
	accessSecret string
	oauthConfig  *oauth2.Config
}

func newHTTPOauth2Client(config *oauth2.Config) *httpOAuth2Client {
	hoc := new(httpOAuth2Client)
	hoc.oauthConfig = config
	return hoc
}

func (hoc *httpOAuth2Client) getMarshalledAPIResponse(url string,
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

func (hoc *httpOAuth2Client) postResource(authToken string, url string, resource interface{},
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
