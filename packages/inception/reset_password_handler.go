package inception

import (
	"net/http"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// ResetPasswordHandler handles POST /reset-password requests.
type ResetPasswordHandler struct {
	app *Inception
}

// NewResetPasswordHandler creates a new reset password handler.
func NewResetPasswordHandler(app *Inception) *ResetPasswordHandler {
	return &ResetPasswordHandler{app: app}
}

// ServeHTTP handles the reset password request.
func (h *ResetPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.ResetPasswords {
		http.Error(w, "password resets are disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()
	config := h.app.config
	input := RequestInput(r, config.IdentifierField, "password", "password_confirmation", "token")

	if input["password"] == "" {
		http.Error(w, "password is required", http.StatusUnprocessableEntity)

		return
	}

	err := h.app.broker.Reset(ctx, input, func(user cauth.Authenticatable, password string) error {
		return h.app.resetPass.Reset(ctx, user, password)
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(ctx, EventPasswordReset)
	}

	h.app.responder.PasswordResetResponse(w, r)
}
