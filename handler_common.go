package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/thedevsaddam/renderer"
)

var sessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type handlerCommon struct {
	rnd *renderer.Render
}

const (
	messageInvalidBoardID   = "Invalid trello board id specified"
	messageInvalidBoardList = "Invalid board or list specified"
	messageSavedTrello      = "Saved trello configuration"
	messageSavedGTasks      = "Saved google tasks configuration"
)

// Message is the message payload that would be send to client
type Message struct {
	ErrorMessage string `json:"message"`
}

func newHandlerCommon() *handlerCommon {
	h := &handlerCommon{}
	opts := renderer.Options{
		ParseGlobPattern: "./templates/*.html",
	}
	h.rnd = renderer.New(opts)
	return h
}

func (hc *handlerCommon) ProcessErrorMessage(message string, w http.ResponseWriter, values ...interface{}) {
	if len(values) > 0 {
		message = fmt.Sprintf(message, values...)
	}
	payload, _ := json.Marshal(&Message{ErrorMessage: message})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(payload)
}

func (hc *handlerCommon) processAuthorizationError(message string,
	w http.ResponseWriter, values ...interface{}) {
	if len(values) > 0 {
		message = fmt.Sprintf(message, values...)
	}
	payload, _ := json.Marshal(&Message{ErrorMessage: message})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(payload)
}

func (hc *handlerCommon) ProcessSuccessMessage(message string, w http.ResponseWriter) {
	payload, _ := json.Marshal(&Message{ErrorMessage: message})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (hc *handlerCommon) ProcessResponse(response interface{}, w http.ResponseWriter) {
	payload, _ := json.Marshal(&response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (hc *handlerCommon) GetUserIDFromSession(r *http.Request) int {
	session, err := sessionStore.Get(r, "userSession")
	if err != nil {
		Error(err)
	}
	userID := session.Values["userID"].(int)
	return userID
}

func (hc *handlerCommon) SaveUserIDInSession(r *http.Request, w http.ResponseWriter, userID int) {
	session, err := sessionStore.Get(r, "userSession")
	if err != nil {
		Error(err)
	}
	session.Values["userID"] = userID
	session.Save(r, w)
}
