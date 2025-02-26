package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/villaleo/calculator/lx"
)

// ServerAdapter wraps an lx.Server to register HTTP handlers.
type ServerAdapter struct {
	*lx.Server
}

// NewServerAdapter creates a new ServerAdapter ready to use.
func NewServerAdapter(opts ...lx.ServerOptsFunc) *ServerAdapter {
	return &ServerAdapter{
		Server: lx.NewServer(opts...),
	}
}

func (s *ServerAdapter) handleCalculation(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger(r.Context())

	if err := s.MatchesMethod(r, http.MethodPost); err != nil {
		logger.Error("method not allowed", "method", r.Method, "url", r.URL, "error", err.Error())
		lx.WriteError(w, http.StatusMethodNotAllowed, err.Error())
		return
	}

	defer r.Body.Close()

	var calcReq CalculationRequest
	if err := s.DecodeRequestBody(r, &calcReq); err != nil {
		logger.Error("received bad request", "method", r.Method, "url", r.URL, "error", err.Error())
		lx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var resp CalculationResponse
	var err error

	// Safely extract the operation name
	operation := strings.TrimPrefix(r.URL.String(), "/")

	switch operation {
	case "add":
		resp.Result = calcReq.X + calcReq.Y
		resp.Interpretation = calcReq.Interpret('+')
	case "subtract":
		resp.Result = calcReq.X - calcReq.Y
		resp.Interpretation = calcReq.Interpret('-')
	case "multiply":
		resp.Result = calcReq.X * calcReq.Y
		resp.Interpretation = calcReq.Interpret('*')
	case "divide":
		if calcReq.Y == 0 {
			err = errors.New("division by zero")
		} else {
			resp.Result = calcReq.X / calcReq.Y
			resp.Interpretation = calcReq.Interpret('/')
		}
	default:
		logger.Warn("unknown operation", "method", r.Method, "url", r.URL, "operation", operation)
		lx.WriteError(w, http.StatusNotFound, "unknown operation")
		return
	}

	if err != nil {
		logger.Error("received unacceptable request", "method", r.Method, "url", r.URL, "error", err)
		lx.WriteError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	lx.WriteJson(w, resp)
}

func (s *ServerAdapter) handleSum(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger(r.Context())

	if err := s.MatchesMethod(r, http.MethodPost); err != nil {
		logger.Error("method not allowed", "method", r.Method, "url", r.URL, "error", err.Error())
		lx.WriteError(w, http.StatusMethodNotAllowed, err.Error())
		return
	}

	defer r.Body.Close()

	var sumReq SumRequest
	if err := s.DecodeRequestBody(r, &sumReq); err != nil {
		logger.Error("received bad request", "method", r.Method, "url", r.URL, "error", err.Error())
		lx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var total float64
	for _, num := range sumReq.Numbers {
		total += num
	}

	lx.WriteJson(w, CalculationResponse{
		Interpretation: sumReq.Interpret(),
		Result:         total,
	})
}
