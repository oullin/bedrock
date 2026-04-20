package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/bedrock/packages/httpx/routingx"
	"github.com/bedrock/packages/routing"
	democonfig "github.com/bedrock/services/demo/api/config"
	"github.com/bedrock/services/demo/api/routes"
)

// NewHandler builds the runnable Laravel-skeleton-equivalent HTTP handler.
func NewHandler(opts ...Options) (http.Handler, error) {
	o := resolveOptions(opts)
	application, err := newApplication(o)

	if err != nil {
		return nil, err
	}

	application.Container.Instance("demo.config.app", democonfig.DefaultApp(o.Env, o.AppKey))

	rawRouter, err := application.Make("router")

	if err != nil {
		return nil, fmt.Errorf("demo: resolve router: %w", err)
	}

	router, ok := rawRouter.(*routing.Router)

	if !ok {
		return nil, fmt.Errorf("demo: router binding has type %T", rawRouter)
	}

	routes.RegisterWeb(router, application)

	return routingx.NewHandler(router), nil
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
