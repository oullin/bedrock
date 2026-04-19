package features

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/services/demo/inertia/api/internal/httputil"
)

func (a app) globalEventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Events/GlobalEvents", protocol.Props{})
	case http.MethodPost:
		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Event", Message: "Global event action completed."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.events.global-events", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) visitCallbacksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Events/VisitCallbacks", protocol.Props{})
	case http.MethodPost:
		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Callback", Message: "Visit callback action completed."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.events.visit-callbacks", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) progressHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Events/Progress", protocol.Props{})
}

func (a app) progressSlowHandler(w http.ResponseWriter, r *http.Request) {
	if httputil.SleepCtx(r.Context(), 2*time.Second) != nil {
		return
	}

	a.container.Render(w, r, "Features/Events/Progress", protocol.Props{})
}
