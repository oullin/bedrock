package fortify

import "net/http"

// ResetPasswordHandler handles POST /reset-password requests.
type ResetPasswordHandler struct {
	fortify *Fortify
}

// NewResetPasswordHandler creates a new reset password handler.
func NewResetPasswordHandler(f *Fortify) *ResetPasswordHandler {
	return &ResetPasswordHandler{fortify: f}
}

// ServeHTTP handles the reset password request.
func (h *ResetPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.fortify.config.Features.ResetPasswords {
		http.Error(w, "password resets are disabled", http.StatusNotFound)
		return
	}

	ctx := r.Context()
	config := h.fortify.config
	input := RequestInput(r, config.IdentifierField, "password", "password_confirmation", "token")

	err := h.fortify.broker.Reset(ctx, input, func(user Authenticatable, password string) error {
		return h.fortify.resetPass.Reset(ctx, user, password)
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if h.fortify.events != nil {
		_ = h.fortify.events.Dispatch(ctx, Event{
			Name: EventPasswordReset,
		})
	}

	h.fortify.responder.PasswordResetResponse(w, r)
}
