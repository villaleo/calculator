package main

import (
	"github.com/villaleo/calculator/lx"
)

func main() {
	server := NewServerAdapter(lx.WithAddr("localhost:3000"))
	server.UseMiddleware(server.RequestIdMiddleware)

	// Register arithmetic handlers
	operations := []string{"add", "subtract", "multiply", "divide"}
	for _, operation := range operations {
		server.HandleFunc("POST /"+operation, server.handleCalculation)
	}

	server.HandleFunc("POST /sum", server.handleSum)
	server.ListenAndServe()
}
