package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/dghubble/oauth1"
)

type httpOAuth1Client struct {
	accessToken  string
	accessSecret string
	oauthConfig  oauth1.Config
}

func newHTTPOAuth1Client(aT string, aS string, oC oauth1.Config) *httpOAuth1Client {
	hoc := new(httpOAuth1Client)
	hoc.accessToken = aT
	hoc.accessSecret = aS
	hoc.oauthConfig = oC
	return hoc
}

func (hoc *httpOAuth1Client) getMarshalledAPIResponse(url string, responseContainer interface{}) error {
	token := oauth1.NewToken(hoc.accessToken, hoc.accessSecret)
	// httpClient will automatically authorize http.Request's
	httpClient := hoc.oauthConfig.Client(oauth1.NoContext, token)
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//Info(string(body))
	err = json.Unmarshal(body, responseContainer)
	if err != nil {
		return err
	}
	return nil
}

func (hoc *httpOAuth1Client) postResource(url string, resource interface{},
	responseContainer interface{}) error {
	token := oauth1.NewToken(hoc.accessToken, hoc.accessSecret)
	httpClient := hoc.oauthConfig.Client(oauth1.NoContext, token)
	requestBody, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//Only if respose container is passed response need to be unmarshalled
	if responseContainer != nil {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(responseBody, responseContainer)
		if err != nil {
			return err
		}
	}
	return nil
}
