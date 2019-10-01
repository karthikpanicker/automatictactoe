package web

import (
	"etsello/common"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetBoardLists(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/users/12345/trello-boards/54321/lists", nil)
	if err != nil {
		t.Fatal(err)
	}
	requestParams := make(map[string]string)
	requestParams["userId"] = "12345"
	requestParams["boardId"] = "54321"
	req = mux.SetURLVars(req, requestParams)
	rr := httptest.NewRecorder()
	dStore := NewTestDataStore()
	ah := newAPIHandler(dStore, "../templates/*.html")
	ah.handlerCom.SaveKeyValueToSession(req, rr, common.UserID, 12345)
	handler := http.HandlerFunc(ah.getBoardLists)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	assert.Equal(t, status, http.StatusBadRequest)
}
