package authflows

import "net/http"

// ForgotPasswordHandler handles POST /forgot-password requests.
type ForgotPasswordHandler struct {
	authflows *AuthFlows
}

// NewForgotPasswordHandler creates a new forgot password handler.
func NewForgotPasswordHandler(f *AuthFlows) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{authflows: f}
}

// ServeHTTP handles the forgot password request.
func (h *ForgotPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.ResetPasswords {
		http.Error(w, "password resets are disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()
	input := RequestInput(r, h.authflows.config.IdentifierField)

	if err := h.ensureNotRateLimited(input[h.authflows.config.IdentifierField], r); err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)

		return
	}

	if err := h.authflows.broker.SendResetLink(ctx, input); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	h.authflows.responder.PasswordResetLinkSentResponse(w, r)
}

func (h *ForgotPasswordHandler) ensureNotRateLimited(identifier string, r *http.Request) error {
	if h.authflows.limiter == nil {
		return nil
	}

	key := "password_reset|" + ThrottleKey(identifier, RequestIP(r))

	if h.authflows.limiter.TooManyAttempts(key, 1) {
		return ErrTooManyAttempts
	}

	h.authflows.limiter.Hit(key, h.authflows.config.LoginRateDecay)

	return nil
}
