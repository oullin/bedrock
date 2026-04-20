package routes

import (
	"fmt"
	"net/http"

	"github.com/bedrock/packages/container"
	democonfig "github.com/bedrock/services/demo/api/config"
)

// RegisterWeb mounts the Upstream skeleton web and health routes.
func RegisterWeb(mux *http.ServeMux, application *container.Application) {
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		app := democonfig.DefaultApp("", "")

		if raw, err := application.Make("demo.config.app"); err == nil {
			if cfg, ok := raw.(democonfig.App); ok {
				app = cfg
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprintf(w, "<!doctype html><html><head><title>%s</title></head><body><main><h1>%s</h1><p>Upstream skeleton port running on Bedrock.</p></main></body></html>\n", app.Name, app.Name)
	})

	mux.HandleFunc("GET /up", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("OK\n"))
	})
}
