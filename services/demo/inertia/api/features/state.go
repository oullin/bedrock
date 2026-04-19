package features

import (
	"log/slog"
	"net/http"

	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
)

func (a app) rememberHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/State/Remember", protocol.Props{})
}

func (a app) flashDataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/State/FlashData", protocol.Props{})
	case http.MethodPost:
		if err := protocol.ParseForm(r); err != nil {
			slog.Error("flash: parse form", "error", err)

			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}

		kind := r.FormValue("kind")

		var msg flash.Message

		switch kind {
		case "error":
			msg = flash.Message{Kind: "error", Title: "Error", Message: "Something went wrong!"}
		case "warning":
			msg = flash.Message{Kind: "warning", Title: "Warning", Message: "Proceed with caution."}
		default:
			msg = flash.Message{Kind: "success", Title: "Flash sent", Message: "This is a success flash message."}
		}

		if err := a.container.SetFlash(w, msg); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.state.flash-data", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) flashDataActionHandler(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")

	var msg flash.Message

	switch action {
	case "error":
		msg = flash.Message{Kind: "error", Title: "Error", Message: "Something went wrong!"}
	case "warning":
		msg = flash.Message{Kind: "warning", Title: "Warning", Message: "Proceed with caution."}
	default:
		msg = flash.Message{Kind: "info", Title: "Info", Message: "Here's some information."}
	}

	if err := a.container.SetFlash(w, msg); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.container.Redirect(w, r, a.container.RouteURL("features.state.flash-data", nil))
}

func (a app) sharedPropsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/State/SharedProps", protocol.Props{})
}
