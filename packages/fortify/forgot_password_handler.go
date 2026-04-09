package fortify

import "net/http"

// ForgotPasswordHandler handles POST /forgot-password requests.
type ForgotPasswordHandler struct {
	fortify *Fortify
}

// NewForgotPasswordHandler creates a new forgot password handler.
func NewForgotPasswordHandler(f *Fortify) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{fortify: f}
}

// ServeHTTP handles the forgot password request.
func (h *ForgotPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.fortify.config.Features.ResetPasswords {
		http.Error(w, "password resets are disabled", http.StatusNotFound)
		return
	}

	ctx := r.Context()
	input := RequestInput(r, h.fortify.config.IdentifierField)

	if err := h.ensureNotRateLimited(input[h.fortify.config.IdentifierField], r); err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	if err := h.fortify.broker.SendResetLink(ctx, input); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	h.fortify.responder.PasswordResetLinkSentResponse(w, r)
}

func (h *ForgotPasswordHandler) ensureNotRateLimited(identifier string, r *http.Request) error {
	if h.fortify.limiter == nil {
		return nil
	}

	key := "password_reset|" + ThrottleKey(identifier, RequestIP(r))

	if h.fortify.limiter.TooManyAttempts(key, 1) {
		return ErrTooManyAttempts
	}

	h.fortify.limiter.Hit(key, h.fortify.config.LoginRateDecay)

	return nil
}
