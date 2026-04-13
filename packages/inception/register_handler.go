package inception

import (
	"net/http"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// RegisterHandler handles POST /register requests.
type RegisterHandler struct {
	app *Inception
}

// NewRegisterHandler creates a new registration handler.
func NewRegisterHandler(app *Inception) *RegisterHandler {
	return &RegisterHandler{app: app}
}

// ServeHTTP handles the registration request.
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.Registration {
		http.Error(w, "registration is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()
	config := h.app.config
	input := RequestInput(r, config.IdentifierField, "name", "password", "password_confirmation")

	user, err := h.app.createUser.Create(ctx, input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(ctx, RegisteredPayload{User: user})
	}

	if config.Features.EmailVerification {
		if verifiable, ok := user.(cauth.MustVerifyEmail); ok && !verifiable.HasVerifiedEmail() {
			if h.app.verifier != nil {
				_ = h.app.verifier.SendVerificationNotification(ctx, user)
			}
		}
	}

	if err := h.app.guard.Login(ctx, w, user, false); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	h.app.responder.RegisterResponse(w, r)
}
