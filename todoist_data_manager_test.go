package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

func TestGetTodoistAuthorizationURL(t *testing.T) {
	gotenv.OverApply(strings.NewReader("TODOIST_CLIENT_ID=abc"))
	gotenv.OverApply(strings.NewReader("TODOIST_CLIENT_SECRET=abc"))
	gotenv.OverApply(strings.NewReader("HOST_URL=http://localhost:80/"))
	gotenv.OverApply(strings.NewReader("TODOIST_AUTH_URL=http://localhost/oauth/authorize"))
	gotenv.OverApply(strings.NewReader("TODOIST_TOKEN_URL=http://localhost/oauth/access_token"))
	gtm := newToDoistDataManager()
	authURL := gtm.getAuthorizationURL()
	assert.Equal(t, "http://localhost/oauth/authorize?access_type=offline&"+
		"client_id=abc&"+
		"redirect_uri=http%3A%2F%2Flocalhost%3A80%2Fcallback-todoist&response_type=code&"+
		"scope=data%3Aread_write&state=state-token", authURL)
}
