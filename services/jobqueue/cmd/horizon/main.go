// Package main runs the Bedrock JobQueue dashboard HTTP server.
//
// Usage:
//
//	go run ./services/jobqueue/cmd/jobqueue
//
// The server listens on $PORT (default 8080) and serves the JSON API used by
// the Vue SPA in services/jobqueue/app.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bedrock/services/jobqueue/api"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	addr := ":8080"

	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		addr = ":" + port
	}

	handler := api.NewHandler(api.Options{
		Auth:        api.AllowAll,
		Supervisors: api.NewSliceSupervisors(nil),
		Batches:     api.NewInMemoryBatches(),
	})

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		_ = server.Shutdown(shutdownCtx)
	}()

	log.Printf("jobqueue dashboard listening on %s", addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
