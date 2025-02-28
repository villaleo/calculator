package main

import (
	"errors"
	"net/http"
	"strings"
)

func (s *ServerAdapter) handleCalculation(w http.ResponseWriter, r *http.Request) {
	if err := MatchesMethod(r, http.MethodPost); err != nil {
		s.Logger.Error("method not allowed", "method", r.Method, "url", r.URL.String(), "error", err.Error())
		WriteError(w, http.StatusMethodNotAllowed, err.Error())
		return
	}

	defer r.Body.Close()

	var calcReq CalculationRequest
	if err := DecodeRequestBody(r, &calcReq); err != nil {
		s.Logger.Error("received bad request", "method", r.Method, "url", r.URL.String(), "error", err.Error())
		WriteError(w, http.StatusBadRequest, err.Error())
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
		s.Logger.Warn("unknown operation", "method", r.Method, "url", r.URL.String(), "operation", operation)
		WriteError(w, http.StatusNotFound, "unknown operation")
		return
	}

	if err != nil {
		s.Logger.Error("received unacceptable request", "method", r.Method, "url", r.URL.String(), "error", err)
		WriteError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	WriteJson(w, resp)
}

func (s *ServerAdapter) handleSum(w http.ResponseWriter, r *http.Request) {
	if err := MatchesMethod(r, http.MethodPost); err != nil {
		s.Logger.Error("method not allowed", "method", r.Method, "url", r.URL.String(), "error", err.Error())
		WriteError(w, http.StatusMethodNotAllowed, err.Error())
		return
	}

	defer r.Body.Close()

	var sumReq SumRequest
	if err := DecodeRequestBody(r, &sumReq); err != nil {
		s.Logger.Error("received bad request", "method", r.Method, "url", r.URL.String(), "error", err.Error())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var total float64
	for _, num := range sumReq.Numbers {
		total += num
	}

	WriteJson(w, CalculationResponse{
		Interpretation: sumReq.Interpret(),
		Result:         total,
	})
}
