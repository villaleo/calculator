package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ErrorResponse is the general repsonse object for sending errors back to the
// client.
type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// WriteJSONWithStatus writes val to w (encoded using json.Marshal) and sets
// the status code. If json.Marshal fails for val, an internal server error
// is written instead.
func WriteJSONWithStatus(w http.ResponseWriter, statusCode int, val any) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(val)
	if err != nil {
		http.Error(w, `{"status":"internal server error","error":"an internal server error occurred"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(data)
}

// WriteError writes an error to w with a message and status code.
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	WriteJSONWithStatus(w, statusCode, ErrorResponse{
		Status: strings.ToLower(http.StatusText(statusCode)),
		Error:  strings.ToLower(strings.TrimSpace(message)),
	})
}

// WriteJSON writes val to w (encoded using json.Marshal) and sets the status
// code to http.StatusOk (200). If json.Marshal fails for val, an internal
// server error is written instead.
func WriteJSON(w http.ResponseWriter, val any) {
	WriteJSONWithStatus(w, http.StatusOK, val)
}
