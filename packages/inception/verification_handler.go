package inception

import (
	"net/http"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// SendVerificationHandler handles POST /email/verification-notification requests.
type SendVerificationHandler struct {
	app *Inception
}

// NewSendVerificationHandler creates a new send verification handler.
func NewSendVerificationHandler(app *Inception) *SendVerificationHandler {
	return &SendVerificationHandler{app: app}
}

// ServeHTTP handles the send verification notification request.
func (h *SendVerificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.EmailVerification {
		http.Error(w, "email verification is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.app.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	if verifiable, ok := user.(cauth.MustVerifyEmail); ok && verifiable.HasVerifiedEmail() {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if h.app.limiter != nil {
		key := "verification|" + user.GetAuthIdentifier()

		if h.app.limiter.TooManyAttempts(key, 1) {
			http.Error(w, ErrTooManyAttempts.Error(), http.StatusTooManyRequests)

			return
		}

		h.app.limiter.Hit(key, h.app.config.LoginRateDecay)
	}

	if err := h.app.verifier.SendVerificationNotification(ctx, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	h.app.responder.EmailVerificationSentResponse(w, r)
}

// VerifyEmailHandler handles GET /verify-email/{id}/{hash} requests.
type VerifyEmailHandler struct {
	app *Inception
}

// NewVerifyEmailHandler creates a new verify email handler.
func NewVerifyEmailHandler(app *Inception) *VerifyEmailHandler {
	return &VerifyEmailHandler{app: app}
}

// ServeHTTP handles the email verification request.
func (h *VerifyEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.EmailVerification {
		http.Error(w, "email verification is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	id := r.PathValue("id")
	hash := r.PathValue("hash")

	if id == "" || hash == "" {
		http.Error(w, "invalid verification link", http.StatusBadRequest)

		return
	}

	if err := h.app.verifier.Verify(ctx, id, hash); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(ctx, EventVerified)
	}

	w.WriteHeader(http.StatusOK)
}
