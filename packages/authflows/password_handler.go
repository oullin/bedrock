package authflows

import "net/http"

// UpdatePasswordHandler handles PUT /user/password requests.
type UpdatePasswordHandler struct {
	authflows *AuthFlows
}

// NewUpdatePasswordHandler creates a new update password handler.

// ServeHTTP handles the password update request.

// ConfirmPasswordHandler handles POST /user/confirm-password requests.
type ConfirmPasswordHandler struct {
	authflows *AuthFlows
}

func NewUpdatePasswordHandler(f *AuthFlows) *UpdatePasswordHandler {
	return &UpdatePasswordHandler{authflows: f}
}

func (h *UpdatePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.UpdatePasswords {
		http.Error(w, "password updates are disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.authflows.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	input := RequestInput(r, "current_password", "password", "password_confirmation")

	if err := h.authflows.updatePass.Update(ctx, user, input); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.authflows.events != nil {
		_ = h.authflows.events.Dispatch(ctx, Event{Name: EventPasswordUpdated})
	}

	h.authflows.responder.PasswordUpdateResponse(w, r)
}

// NewConfirmPasswordHandler creates a new confirm password handler.
func NewConfirmPasswordHandler(f *AuthFlows) *ConfirmPasswordHandler {
	return &ConfirmPasswordHandler{authflows: f}
}

// ServeHTTP handles the password confirmation request.
func (h *ConfirmPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.ConfirmPassword {
		http.Error(w, "password confirmation is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.authflows.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	input := RequestInput(r, "password")
	password := input["password"]

	if err := h.authflows.confirmPass.Confirm(ctx, user, password); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	SetPasswordConfirmedAt(r)

	h.authflows.responder.PasswordConfirmResponse(w, r)
}
