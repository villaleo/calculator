package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	logFile := "server.log"
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("failed to open file %q: %w", logFile, err))
	}

	mux := http.NewServeMux()
	srv := NewServerAdapter(
		WithHandler(mux),
		WithAddr("localhost:3000"),
		WithWriter(file),
	)

	operations := []string{"add", "subtract", "multiply", "divide"}
	for _, operation := range operations {
		calculationHandler := http.HandlerFunc(srv.handleCalculation)
		mux.Handle("POST /"+operation, EnforceJSONHandler(srv.LogRequestHandler(calculationHandler)))
	}

	sumHandler := http.HandlerFunc(srv.handleSum)
	mux.Handle("POST /sum", EnforceJSONHandler(srv.LogRequestHandler(sumHandler)))
	srv.ListenAndServe()
}
