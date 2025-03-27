package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MatchesMethod ensures that the request method matches the specified HTTP
// method. A non-nil error means the method does not match the request method.
func MatchesMethod(r *http.Request, method string) error {
	if r.Method != method {
		err := fmt.Errorf("method %s not allowed", r.Method)
		return err
	}

	return nil
}

// DecodeRequestBody decodes the request body into val using json.Decode. val
// should be passed by reference to ensure the decoded value is stored
// correctly.
func DecodeRequestBody(r *http.Request, val any) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&val); err != nil {
		return fmt.Errorf("received malformed request: %w", err)
	}

	return nil
}
