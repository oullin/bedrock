package authflows

import (
	"net/http"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// SendVerificationHandler handles POST /email/verification-notification requests.
type SendVerificationHandler struct {
	authflows *AuthFlows
}

// NewSendVerificationHandler creates a new send verification handler.

// ServeHTTP handles the send verification notification request.

// VerifyEmailHandler handles GET /verify-email/{id}/{hash} requests.
type VerifyEmailHandler struct {
	authflows *AuthFlows
}

func NewSendVerificationHandler(f *AuthFlows) *SendVerificationHandler {
	return &SendVerificationHandler{authflows: f}
}

func (h *SendVerificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.EmailVerification {
		http.Error(w, "email verification is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.authflows.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	if verifiable, ok := user.(cauth.MustVerifyEmail); ok && verifiable.HasVerifiedEmail() {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if h.authflows.limiter != nil {
		key := "verification|" + user.GetAuthIdentifier()

		if h.authflows.limiter.TooManyAttempts(key, 1) {
			http.Error(w, ErrTooManyAttempts.Error(), http.StatusTooManyRequests)

			return
		}

		h.authflows.limiter.Hit(key, h.authflows.config.LoginRateDecay)
	}

	if err := h.authflows.verifier.SendVerificationNotification(ctx, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	h.authflows.responder.EmailVerificationSentResponse(w, r)
}

// NewVerifyEmailHandler creates a new verify email handler.
func NewVerifyEmailHandler(f *AuthFlows) *VerifyEmailHandler {
	return &VerifyEmailHandler{authflows: f}
}

// ServeHTTP handles the email verification request.
func (h *VerifyEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.EmailVerification {
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

	if err := h.authflows.verifier.Verify(ctx, id, hash); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)

		return
	}

	if h.authflows.events != nil {
		_, _ = h.authflows.events.Dispatch(ctx, EventVerified)
	}

	w.WriteHeader(http.StatusOK)
}
