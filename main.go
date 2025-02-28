package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	srv := NewServerAdapter(
		WithHandler(mux),
		WithAddr("localhost:3000"),
	)

	// Register calculation endpoints
	operations := []string{"add", "subtract", "multiply", "divide"}
	for _, operation := range operations {
		mux.HandleFunc("POST /"+operation, srv.handleCalculation)
	}

	mux.HandleFunc("POST /sum", srv.handleSum)
	srv.ListenAndServe()
}
