package main

import (
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	file, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}

	srv := NewServerAdapter(
		WithHandler(mux),
		WithAddr("localhost:3000"),
		WithWriter(file),
	)

	// Register calculation endpoints
	operations := []string{"add", "subtract", "multiply", "divide"}
	for _, operation := range operations {
		calculationHandler := http.HandlerFunc(srv.handleCalculation)
		mux.Handle("POST /"+operation, EnforceJSONHandler(srv.LogRequestHandler(calculationHandler)))
	}

	sumHandler := http.HandlerFunc(srv.handleSum)
	mux.Handle("POST /sum", EnforceJSONHandler(srv.LogRequestHandler(sumHandler)))
	srv.ListenAndServe()
}
