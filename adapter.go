package main

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// OptsFunc is a function that provides flexibility with Opts.
type OptsFunc func(*Opts)

// Opts contains options for ServerAdapter, providing flexibility.
type Opts struct {
	addr    string
	handler http.Handler
	writer  io.Writer
	logger  slog.Handler
}

// Default returns the default ServerAdapter configuration.
func Default() Opts {
	return Opts{
		addr:    ":8080",
		handler: http.NewServeMux(),
		writer:  os.Stdout,
		logger:  slog.NewJSONHandler(os.Stdout, nil),
	}
}

// WithAddr sets the ServerAdapter's address.
func WithAddr(addr string) OptsFunc {
	return func(opt *Opts) {
		opt.addr = addr
	}
}

// WithHandler sets the ServerAdapter's handler.
func WithHandler(handler http.Handler) OptsFunc {
	return func(opt *Opts) {
		opt.handler = handler
	}
}

// WithWriter sets the ServerAdapter's writer.
func WithWriter(w io.Writer) OptsFunc {
	return func(opt *Opts) {
		opt.writer = w
		opt.logger = slog.NewJSONHandler(w, nil)
	}
}

// WithLogger sets the ServerAdapter's logger.
func WithLogger(logger slog.Handler) OptsFunc {
	return func(opt *Opts) {
		opt.logger = logger
	}
}

// ServerAdapter wraps an http.Server, providing structured logging.
type ServerAdapter struct {
	*http.Server
	Opts

	// Logger enables structured logging, slog.New() with a slog.NewJSONHandler()
	// if nil.
	Logger *slog.Logger
}

// NewServerAdapter creates a new ServerAdapter.
func NewServerAdapter(opts ...OptsFunc) *ServerAdapter {
	options := Default()

	for _, optionFunc := range opts {
		optionFunc(&options)
	}

	return &ServerAdapter{
		Server: &http.Server{
			Addr:    options.addr,
			Handler: options.handler,
		},
		Logger: slog.New(options.logger),
	}
}

// ListenAndServe starts the server and gracefully shuts down on termination
// signals.
func (s *ServerAdapter) ListenAndServe() {
	s.Logger.Info("server started", "addr", s.Addr)

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	serverErrCh := make(chan error, 1)

	go func() {
		err := s.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	// Wait for either an interrupt signal or a server error.
	select {
	case <-interruptCh:
		s.Logger.Info("shutting down server")
		s.Cleanup()
		os.Exit(0)
	case err := <-serverErrCh:
		s.Logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

// Cleanup frees any resources (closing database connections, etc). It is
// automatically called before the server gracefully shuts down.
func (s *ServerAdapter) Cleanup() {
	s.Logger.Info("cleaning up resources")
}
