package fortify

import (
	"net/http"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// SendVerificationHandler handles POST /email/verification-notification requests.
type SendVerificationHandler struct {
	fortify *Fortify
}

// NewSendVerificationHandler creates a new send verification handler.

// ServeHTTP handles the send verification notification request.

// VerifyEmailHandler handles GET /verify-email/{id}/{hash} requests.
type VerifyEmailHandler struct {
	fortify *Fortify
}

func NewSendVerificationHandler(f *Fortify) *SendVerificationHandler {
	return &SendVerificationHandler{fortify: f}
}

func (h *SendVerificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.fortify.config.Features.EmailVerification {
		http.Error(w, "email verification is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.fortify.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	if verifiable, ok := user.(cauth.MustVerifyEmail); ok && verifiable.HasVerifiedEmail() {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if h.fortify.limiter != nil {
		key := "verification|" + user.GetAuthIdentifier()

		if h.fortify.limiter.TooManyAttempts(key, 1) {
			http.Error(w, ErrTooManyAttempts.Error(), http.StatusTooManyRequests)

			return
		}

		h.fortify.limiter.Hit(key, h.fortify.config.LoginRateDecay)
	}

	if err := h.fortify.verifier.SendVerificationNotification(ctx, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	h.fortify.responder.EmailVerificationSentResponse(w, r)
}

// NewVerifyEmailHandler creates a new verify email handler.
func NewVerifyEmailHandler(f *Fortify) *VerifyEmailHandler {
	return &VerifyEmailHandler{fortify: f}
}

// ServeHTTP handles the email verification request.
func (h *VerifyEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.fortify.config.Features.EmailVerification {
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

	if err := h.fortify.verifier.Verify(ctx, id, hash); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)

		return
	}

	if h.fortify.events != nil {
		_ = h.fortify.events.Dispatch(ctx, EventVerified)
	}

	w.WriteHeader(http.StatusOK)
}
