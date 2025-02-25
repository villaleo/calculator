package lx

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// MiddlewareFunc is an adaptor that enables introducing middleware into the
// server.
type MiddlewareFunc func(http.Handler) http.Handler

// ctxKeyRequestId is a key for storing a request ID in context.
type ctxKeyRequestId struct{}

// ctxKeyLogger is a key for storing a logger in context.
type ctxKeyLogger struct{}

// RequestId extracts the request ID from ctx.
func RequestId(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ctxKeyRequestId{}).(string)
	return id, ok
}

// RequestIdMiddleware generates or propagates a request ID for each request
// received.
func (s *Server) RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-Request-ID")
		if requestId == "" {
			requestId = uuid.NewString()
		}

		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestId)
		// Use server's logger with the requestId attached
		logger := s.logger.With("requestId", requestId)

		// Store request ID in context
		ctx := context.WithValue(r.Context(), ctxKeyRequestId{}, requestId)
		ctx = context.WithValue(ctx, ctxKeyLogger{}, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
