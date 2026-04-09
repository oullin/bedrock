package fortify

import (
	"context"
	"errors"
	"net/http"
)

// ErrTooManyAttempts is returned when login attempts are rate limited.
var ErrTooManyAttempts = errors.New("fortify: too many login attempts")

// ErrInvalidCredentials is returned when login credentials are invalid.
var ErrInvalidCredentials = errors.New("fortify: invalid credentials")

// LoginHandler handles POST /login requests.
type LoginHandler struct {
	fortify *Fortify
}

// NewLoginHandler creates a new login handler.
func NewLoginHandler(f *Fortify) *LoginHandler {
	return &LoginHandler{fortify: f}
}

// ServeHTTP handles the login request.
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := h.fortify.config
	input := RequestInput(r, config.IdentifierField, "password")
	identifier := input[config.IdentifierField]
	remember := RequestBool(r, "remember")

	if err := h.ensureNotRateLimited(identifier, r); err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	h.dispatchEvent(ctx, EventLoginAttempted, LoginAttemptedPayload{
		Identifier: identifier,
		Remember:   remember,
	})

	user, err := h.authenticate(ctx, input)
	if err != nil {
		h.onFailure(ctx, w, r, identifier)
		return
	}

	if tfa, ok := user.(TwoFactorAuthenticatable); ok && tfa.IsTwoFactorEnabled() && tfa.GetTwoFactorConfirmedAt() != nil {
		if err := h.fortify.guard.LoginWithPendingTwoFactor(ctx, w, user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.fortify.responder.TwoFactorChallengeResponse(w, r)

		return
	}

	h.onSuccess(ctx, w, r, user, remember)
}

func (h *LoginHandler) authenticate(ctx context.Context, input map[string]string) (Authenticatable, error) {
	if custom := h.fortify.authenticator; custom != nil {
		return custom.Authenticate(ctx, input)
	}

	config := h.fortify.config
	identifier := input[config.IdentifierField]
	password := input["password"]

	credentials := map[string]string{
		config.IdentifierField: identifier,
	}

	user, err := h.fortify.provider.RetrieveByCredentials(ctx, credentials)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	valid, err := h.fortify.provider.ValidateCredentials(ctx, user, map[string]string{
		config.IdentifierField: identifier,
		"password":             password,
	})

	if err != nil || !valid {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (h *LoginHandler) ensureNotRateLimited(identifier string, r *http.Request) error {
	limiter := h.fortify.limiter
	if limiter == nil {
		return nil
	}

	key := ThrottleKey(identifier, RequestIP(r))

	if limiter.TooManyAttempts(key, h.fortify.config.LoginRateLimit) {
		return ErrTooManyAttempts
	}

	return nil
}

func (h *LoginHandler) onSuccess(ctx context.Context, w http.ResponseWriter, r *http.Request, user Authenticatable, remember bool) {
	if limiter := h.fortify.limiter; limiter != nil {
		limiter.Clear(ThrottleKey(user.GetAuthIdentifier(), RequestIP(r)))
	}

	if err := h.fortify.guard.Login(ctx, w, user, remember); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.dispatchEvent(ctx, EventLoginSucceeded, LoginSucceededPayload{
		User:     user,
		Remember: remember,
	})

	h.fortify.responder.LoginResponse(w, r)
}

func (h *LoginHandler) onFailure(ctx context.Context, w http.ResponseWriter, r *http.Request, identifier string) {
	if limiter := h.fortify.limiter; limiter != nil {
		key := ThrottleKey(identifier, RequestIP(r))
		limiter.Hit(key, h.fortify.config.LoginRateDecay)
	}

	h.dispatchEvent(ctx, EventLoginFailed, LoginFailedPayload{
		Identifier: identifier,
	})

	http.Error(w, ErrInvalidCredentials.Error(), http.StatusUnprocessableEntity)
}

func (h *LoginHandler) dispatchEvent(ctx context.Context, name string, payload any) {
	if h.fortify.events == nil {
		return
	}

	_ = h.fortify.events.Dispatch(ctx, Event{Name: name, Payload: payload})
}
