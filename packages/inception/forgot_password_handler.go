package inception

import "net/http"

// ForgotPasswordHandler handles POST /forgot-password requests.
type ForgotPasswordHandler struct {
	app *Inception
}

// NewForgotPasswordHandler creates a new forgot password handler.
func NewForgotPasswordHandler(app *Inception) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{app: app}
}

// ServeHTTP handles the forgot password request.
func (h *ForgotPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.ResetPasswords {
		http.Error(w, "password resets are disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()
	input := RequestInput(r, h.app.config.IdentifierField)

	if err := h.ensureNotRateLimited(input[h.app.config.IdentifierField], r); err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)

		return
	}

	if err := h.app.broker.SendResetLink(ctx, input); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	h.app.responder.PasswordResetLinkSentResponse(w, r)
}

func (h *ForgotPasswordHandler) ensureNotRateLimited(identifier string, r *http.Request) error {
	if h.app.limiter == nil {
		return nil
	}

	key := "password_reset|" + ThrottleKey(identifier, RequestIP(r))

	if h.app.limiter.TooManyAttempts(key, 1) {
		return ErrTooManyAttempts
	}

	h.app.limiter.Hit(key, h.app.config.LoginRateDecay)

	return nil
}
