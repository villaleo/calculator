package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var cmd = cobra.Command{
	Use:     "run [-A address] [-F logfile]",
	Short:   "Run the server",
	Example: "run -A localhost:3000 -F server.log",
	Run: func(cmd *cobra.Command, args []string) {
		addr := cmd.Flag("addr").Value.String()
		logFile := cmd.Flag("file").Value.String()
		runServer(addr, logFile)
	},
}

func main() {
	cmd.Flags().StringP("addr", "A", "localhost:8080", "address to listen on")
	cmd.Flags().StringP("file", "F", "/dev/stdout", "file to write logs to")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runServer(addr string, logFile string) {
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("failed to open file %q: %w", logFile, err))
	}
	logger := slog.NewTextHandler(file, nil)

	mux := http.NewServeMux()
	srv := NewServerAdapter(
		WithHandler(mux),
		WithAddr(addr),
		WithLogger(logger),
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
