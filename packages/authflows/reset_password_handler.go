package authflows

import "net/http"

// ResetPasswordHandler handles POST /reset-password requests.
type ResetPasswordHandler struct {
	authflows *AuthFlows
}

// NewResetPasswordHandler creates a new reset password handler.
func NewResetPasswordHandler(f *AuthFlows) *ResetPasswordHandler {
	return &ResetPasswordHandler{authflows: f}
}

// ServeHTTP handles the reset password request.
func (h *ResetPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.ResetPasswords {
		http.Error(w, "password resets are disabled", http.StatusNotFound)
		return
	}

	ctx := r.Context()
	config := h.authflows.config
	input := RequestInput(r, config.IdentifierField, "password", "password_confirmation", "token")

	err := h.authflows.broker.Reset(ctx, input, func(user Authenticatable, password string) error {
		return h.authflows.resetPass.Reset(ctx, user, password)
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if h.authflows.events != nil {
		_ = h.authflows.events.Dispatch(ctx, Event{
			Name: EventPasswordReset,
		})
	}

	h.authflows.responder.PasswordResetResponse(w, r)
}
