package fortify

import (
	"net/http"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// RegisterHandler handles POST /register requests.
type RegisterHandler struct {
	fortify *Fortify
}

// NewRegisterHandler creates a new registration handler.
func NewRegisterHandler(f *Fortify) *RegisterHandler {
	return &RegisterHandler{fortify: f}
}

// ServeHTTP handles the registration request.
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.fortify.config.Features.Registration {
		http.Error(w, "registration is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()
	config := h.fortify.config
	input := RequestInput(r, config.IdentifierField, "name", "password", "password_confirmation")

	user, err := h.fortify.createUser.Create(ctx, input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.fortify.events != nil {
		_, _ = h.fortify.events.Dispatch(ctx, RegisteredPayload{User: user})
	}

	if config.Features.EmailVerification {
		if verifiable, ok := user.(cauth.MustVerifyEmail); ok && !verifiable.HasVerifiedEmail() {
			if h.fortify.verifier != nil {
				_ = h.fortify.verifier.SendVerificationNotification(ctx, user)
			}
		}
	}

	if err := h.fortify.guard.Login(ctx, w, user, false); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	h.fortify.responder.RegisterResponse(w, r)
}
