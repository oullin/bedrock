package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	democonfig "github.com/bedrock/services/demo/api/config"
	"github.com/bedrock/services/demo/api/routes"
)

// NewHandler builds the runnable Upstream-skeleton-equivalent HTTP handler.
func NewHandler(opts ...Options) (http.Handler, error) {
	o := resolveOptions(opts)
	application, err := newApplication(o)

	if err != nil {
		return nil, err
	}

	application.Container.Instance("demo.config.app", democonfig.DefaultApp(o.Env, o.AppKey))

	mux := http.NewServeMux()
	routes.RegisterWeb(mux, application)

	return mux, nil
}

// Run starts the skeleton demo app and shuts it down when ctx is cancelled.
func Run(ctx context.Context, opts ...Options) error {
	handler, err := NewHandler(opts...)

	if err != nil {
		return err
	}

	addr := ":8080"

	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		addr = ":" + port
	}

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return fmt.Errorf("demo: listen: %w", err)
	}

	errc := make(chan error, 1)

	go func() {
		errc <- server.Serve(listener)
	}()

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()

	if url := strings.TrimSpace(os.Getenv("PORTLESS_URL")); url != "" {
		fmt.Printf("Server running at %s\n", url)
	} else {
		fmt.Printf("Server running at http://localhost%s\n", addr)
	}

	if err := <-errc; err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("demo: serve: %w", err)
	}

	return nil
}
