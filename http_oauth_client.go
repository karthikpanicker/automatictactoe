package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dghubble/oauth1"
)

type httpOAuthClient struct {
	accessToken  string
	accessSecret string
	oauthConfig  oauth1.Config
}

func newHTTPOAuthClient(aT string, aS string, oC oauth1.Config) *httpOAuthClient {
	hoc := new(httpOAuthClient)
	hoc.accessToken = aT
	hoc.accessSecret = aS
	hoc.oauthConfig = oC
	return hoc
}

func (hoc *httpOAuthClient) getMarshalledAPIResponse(url string, responseContainer interface{}) error {
	token := oauth1.NewToken(hoc.accessToken, hoc.accessSecret)
	// httpClient will automatically authorize http.Request's
	httpClient := hoc.oauthConfig.Client(oauth1.NoContext, token)
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, responseContainer)
	if err != nil {
		return err
	}
	return nil
}
