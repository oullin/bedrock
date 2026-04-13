package fortify

import "net/http"

// UpdatePasswordHandler handles PUT /user/password requests.
type UpdatePasswordHandler struct {
	fortify *Fortify
}

// NewUpdatePasswordHandler creates a new update password handler.

// ServeHTTP handles the password update request.

// ConfirmPasswordHandler handles POST /user/confirm-password requests.
type ConfirmPasswordHandler struct {
	fortify *Fortify
}

func NewUpdatePasswordHandler(f *Fortify) *UpdatePasswordHandler {
	return &UpdatePasswordHandler{fortify: f}
}

func (h *UpdatePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.fortify.config.Features.UpdatePasswords {
		http.Error(w, "password updates are disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.fortify.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	input := RequestInput(r, "current_password", "password", "password_confirmation")

	if err := h.fortify.updatePass.Update(ctx, user, input); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.fortify.events != nil {
		_, _ = h.fortify.events.Dispatch(ctx, EventPasswordUpdated)
	}

	h.fortify.responder.PasswordUpdateResponse(w, r)
}

// NewConfirmPasswordHandler creates a new confirm password handler.
func NewConfirmPasswordHandler(f *Fortify) *ConfirmPasswordHandler {
	return &ConfirmPasswordHandler{fortify: f}
}

// ServeHTTP handles the password confirmation request.
func (h *ConfirmPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.fortify.config.Features.ConfirmPassword {
		http.Error(w, "password confirmation is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.fortify.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	input := RequestInput(r, "password")
	password := input["password"]

	if err := h.fortify.confirmPass.Confirm(ctx, user, password); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	SetPasswordConfirmedAt(r)

	h.fortify.responder.PasswordConfirmResponse(w, r)
}
