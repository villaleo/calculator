package lx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// ServerOpts contains the Server options which can be changed.
type ServerOpts struct {
	addr   string    // The server address. Default is "localhost:8080".
	writer io.Writer // The logging writer output. Default is os.Stdout.
}

// ServerOptsFunc is an adapter that enables changing the default values of
// the Server using ServerOpts.
type ServerOptsFunc func(*ServerOpts)

func defaultServerOpts() ServerOpts {
	// Ensure the default values remain in sync with the documentation in
	// ServerOpts.
	return ServerOpts{
		addr:   "localhost:8080",
		writer: os.Stdout,
	}
}

// WithAddr sets the Server address to addr instead of the default.
func WithAddr(addr string) ServerOptsFunc {
	return func(so *ServerOpts) {
		so.addr = addr
	}
}

// WithWriter sets the Server logger to log to w instead of the default.
func WithWriter(w io.Writer) ServerOptsFunc {
	return func(so *ServerOpts) {
		so.writer = w
	}
}

// Server wraps an http.ServeMux and a slog.Logger to provide structured
// logging for HTTP requests and operations.
type Server struct {
	ServerOpts
	mux        *http.ServeMux
	logger     *slog.Logger
	middleware []MiddlewareFunc
}

// NewServer initializes a new Server ready to be used.
func NewServer(opts ...ServerOptsFunc) *Server {
	o := defaultServerOpts()
	for _, optFunc := range opts {
		optFunc(&o)
	}

	s := &Server{
		ServerOpts: o,
		mux:        http.NewServeMux(),
		logger:     slog.New(slog.NewTextHandler(o.writer, nil)),
	}

	// Determine write destination
	var writeDestination string
	switch o.writer {
	case os.Stdout:
		writeDestination = "stdout"
	case os.Stderr:
		writeDestination = "stderr"
	case io.Discard:
		writeDestination = "<discarded>"
	default:
		if f, ok := o.writer.(*os.File); ok {
			writeDestination = f.Name()
		} else {
			writeDestination = "<unknown>"
		}
	}

	s.logger.Info("server initialized", "output", writeDestination)

	return s
}

// UseMiddleware registers global middleware for all handlers.
func (s *Server) UseMiddleware(mw MiddlewareFunc) {
	s.middleware = append(s.middleware, mw)
}

// HandleFunc is a wrapper over http.HandleFunc and behaves the same in every
// way.
func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	var finalHandler http.Handler = handler

	// Wrap the handler with middleware, starting from the most recently added
	// middleware.
	for i := len(s.middleware) - 1; i >= 0; i-- {
		finalHandler = s.middleware[i](finalHandler)
	}

	s.logger.Info("handler registered", "pattern", pattern)
	s.mux.HandleFunc(pattern, finalHandler.ServeHTTP)
}

// ListenAndServe starts the server and gracefully shuts down on termination signals.
func (s *Server) ListenAndServe() {
	s.logger.Info("server started", "addr", s.addr)

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	serverErrCh := make(chan error, 1)

	go func() {
		err := http.ListenAndServe(s.addr, s.mux)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	// Wait for either an interrupt signal or a server error.
	select {
	case <-interruptCh:
		s.logger.Info("shutting down server")
		s.Cleanup()
		os.Exit(0)

	case err := <-serverErrCh:
		s.logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

// Cleanup is called before shutting down the server to free resources.
func (s *Server) Cleanup() {
	s.logger.Info("cleaning up resources")
}

// DecodeRequestBody decodes the request body into val using json.Decode. val
// should be passed by reference to ensure the decoded value is stored
// correctly.
func (s *Server) DecodeRequestBody(r *http.Request, val any) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&val); err != nil {
		s.logger.Error("malformed request", "method", r.Method, "url", r.URL, "error", err)
		return fmt.Errorf("received malformed request: %w", err)
	}

	return nil
}

// MatchesMethod ensures that the request method matches the specified HTTP
// method. A non-nil error means the method does not match the request method.
func (s *Server) MatchesMethod(r *http.Request, method string) error {
	if r.Method != method {
		err := fmt.Errorf("method %s not allowed", r.Method)
		s.logger.Warn(err.Error(), "method", r.Method, "url", r.URL)
		return err
	}

	return nil
}

// Logger gets the server's logger instance. When logging from a handler, call
// with r.Context() to use the request's context.
func (s *Server) Logger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(ctxKeyLogger{}).(*slog.Logger)
	if !ok {
		return s.logger
	}
	return logger
}
