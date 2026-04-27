package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"clearbill/mgr/server/internal/app"
)

// @title Clear Bill API
// @version 0.1.0
// @description Clear Bill backend service API.
// @BasePath /api/v1

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := app.NewServer()
	if err := server.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
