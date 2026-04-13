package inception

import "net/http"

// UpdatePasswordHandler handles PUT /user/password requests.
type UpdatePasswordHandler struct {
	app *Inception
}

// NewUpdatePasswordHandler creates a new update password handler.
func NewUpdatePasswordHandler(app *Inception) *UpdatePasswordHandler {
	return &UpdatePasswordHandler{app: app}
}

// ServeHTTP handles the password update request.
func (h *UpdatePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.UpdatePasswords {
		http.Error(w, "password updates are disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.app.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	input := RequestInput(r, "current_password", "password", "password_confirmation")

	if err := h.app.updatePass.Update(ctx, user, input); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(ctx, EventPasswordUpdated)
	}

	h.app.responder.PasswordUpdateResponse(w, r)
}

// ConfirmPasswordHandler handles POST /user/confirm-password requests.
type ConfirmPasswordHandler struct {
	app *Inception
}

// NewConfirmPasswordHandler creates a new confirm password handler.
func NewConfirmPasswordHandler(app *Inception) *ConfirmPasswordHandler {
	return &ConfirmPasswordHandler{app: app}
}

// ServeHTTP handles the password confirmation request.
func (h *ConfirmPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.ConfirmPassword {
		http.Error(w, "password confirmation is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.app.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	input := RequestInput(r, "password")
	password := input["password"]

	if err := h.app.confirmPass.Confirm(ctx, user, password); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	SetPasswordConfirmedAt(r)

	h.app.responder.PasswordConfirmResponse(w, r)
}
