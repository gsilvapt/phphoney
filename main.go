package main

import (
	"log/slog"
	"net/http"
	"os"
)

const (
	version = "dev"
	env     = "dev"
)

func main() {
	m := http.NewServeMux()

	for path, handler := range handlers {
		m.Handle(path, logRequest(handler))
	}

	slog.Info("initializing server, listening on :8080", "version", version)
	if err := http.ListenAndServe(":8080", m); err != nil {
		slog.Error("failed starting server", "error", err)
		os.Exit(1)
	}
}
