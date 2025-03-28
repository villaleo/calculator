package main

import (
	"mime"
	"net/http"
)

// EnforceJSONHandler enforces incoming HTTP requests to have their
// Content-Type header set to 'application/json'.
func EnforceJSONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				WriteError(w, http.StatusBadRequest, "malformed Content-Type header")
				return
			}

			if mt != "application/json" {
				WriteError(w, http.StatusUnsupportedMediaType, "Content-Type header must be application/json")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// LogRequestHandler logs the request to this ServerAdapter's logger.
func (s *ServerAdapter) LogRequestHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("request received", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
