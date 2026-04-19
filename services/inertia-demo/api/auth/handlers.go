package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/bedrock/packages/inertia"
	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
)

func (a App) loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Auth/Login", protocol.Props{
			"status": r.URL.Query().Get("status"),
		})
	case http.MethodPost:
		a.loginSubmitHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a App) loginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if err := protocol.ParseForm(r); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	form := newLoginForm(r)
	errorsByField := form.validate()

	if len(errorsByField) == 0 {
		user, err := a.service.authenticate(form)

		switch {
		case err == nil:
			if err := a.setSession(w, user.ID, form.Remember); err != nil {
				slog.Error("session: set", "error", err)

				http.Error(w, "internal server error", http.StatusInternalServerError)

				return
			}

			if err := a.container.SetFlash(w, flash.Message{
				Kind:    "success",
				Title:   "Signed in",
				Message: "The demo session is now authenticated.",
			}); err != nil {
				slog.Error("flash: set", "error", err)
			}

			a.container.Redirect(w, r, a.container.RouteURL("dashboard", nil))

			return
		case errors.Is(err, errInvalidCredentials):
			errorsByField = protocol.ValidationErrors{
				"email": "Use test@example.com and password to sign in.",
			}
		default:
			slog.Error("authenticate", "error", err)

			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		}
	}

	ctx := inertia.SetValidationErrors(r.Context(), errorsByField)

	a.container.Render(w, r.WithContext(ctx), "Auth/Login", protocol.Props{})
}

func (a App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	a.clearSession(w)

	if err := a.container.SetFlash(w, flash.Message{
		Kind:    "info",
		Title:   "Signed out",
		Message: "Your demo session has been cleared.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.container.Redirect(w, r, a.container.RouteURL("login", nil))
}
