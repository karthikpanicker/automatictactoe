package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/thedevsaddam/renderer"
)

type handlerCommon struct {
	rnd *renderer.Render
}

const (
	messageInvalidBoardID   = "Invalid trello board id specified"
	messageInvalidBoardList = "Invalid board or list specified"
	messageSavedBoardList   = "Linked board and list to etsy shop"
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

func (hc *handlerCommon) ExtractSessionID(r *http.Request) int {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		Error(err)
	}
	sessionID, _ := strconv.Atoi(cookie.Value)
	return sessionID
}
