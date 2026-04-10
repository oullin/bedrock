package authflows

import "net/http"

// RegisterHandler handles POST /register requests.
type RegisterHandler struct {
	authflows *AuthFlows
}

// NewRegisterHandler creates a new registration handler.
func NewRegisterHandler(f *AuthFlows) *RegisterHandler {
	return &RegisterHandler{authflows: f}
}

// ServeHTTP handles the registration request.
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.Registration {
		http.Error(w, "registration is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()
	config := h.authflows.config
	input := RequestInput(r, config.IdentifierField, "name", "password", "password_confirmation")

	user, err := h.authflows.createUser.Create(ctx, input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.authflows.events != nil {
		_ = h.authflows.events.Dispatch(ctx, Event{
			Name:    EventRegistered,
			Payload: RegisteredPayload{User: user},
		})
	}

	if config.Features.EmailVerification {
		if verifiable, ok := user.(MustVerifyEmail); ok && !verifiable.HasVerifiedEmail() {
			if h.authflows.verifier != nil {
				_ = h.authflows.verifier.SendVerificationNotification(ctx, user)
			}
		}
	}

	if err := h.authflows.guard.Login(ctx, w, user, false); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	h.authflows.responder.RegisterResponse(w, r)
}
