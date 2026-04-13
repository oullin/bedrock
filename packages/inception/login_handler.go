package inception

import (
	"context"
	"errors"
	"net/http"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// ErrTooManyAttempts is returned when login attempts are rate limited.

// ErrInvalidCredentials is returned when login credentials are invalid.

// LoginHandler handles POST /login requests.
type LoginHandler struct {
	app *Inception
}

var ErrTooManyAttempts = errors.New("inception: too many login attempts")

var ErrInvalidCredentials = errors.New("inception: invalid credentials")

// NewLoginHandler creates a new login handler.
func NewLoginHandler(app *Inception) *LoginHandler {
	return &LoginHandler{app: app}
}

// ServeHTTP handles the login request.
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := h.app.config
	input := RequestInput(r, config.IdentifierField, "password")
	identifier := input[config.IdentifierField]
	remember := RequestBool(r, "remember")

	if err := h.ensureNotRateLimited(identifier, r); err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)

		return
	}

	h.dispatchEvent(ctx, LoginAttemptedPayload{
		Identifier: identifier,
		Remember:   remember,
	})

	user, err := h.authenticate(ctx, input)

	if err != nil {
		h.onFailure(ctx, w, r, identifier)

		return
	}

	if tfa, ok := user.(cauth.TwoFactorAuthenticatable); ok && tfa.IsTwoFactorEnabled() && tfa.GetTwoFactorConfirmedAt() != nil {
		if err := h.app.guard.LoginWithPendingTwoFactor(ctx, w, user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		h.app.responder.TwoFactorChallengeResponse(w, r)

		return
	}

	h.onSuccess(ctx, w, r, user, remember)
}

func (h *LoginHandler) authenticate(ctx context.Context, input map[string]string) (cauth.Authenticatable, error) {
	if custom := h.app.authenticator; custom != nil {
		return custom.Authenticate(ctx, input)
	}

	config := h.app.config
	identifier := input[config.IdentifierField]
	password := input["password"]

	credentials := map[string]string{
		config.IdentifierField: identifier,
	}

	user, err := h.app.provider.RetrieveByCredentials(ctx, credentials)

	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	valid, err := h.app.provider.ValidateCredentials(ctx, user, map[string]string{
		config.IdentifierField: identifier,
		"password":             password,
	})

	if err != nil || !valid {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (h *LoginHandler) ensureNotRateLimited(identifier string, r *http.Request) error {
	limiter := h.app.limiter

	if limiter == nil {
		return nil
	}

	key := ThrottleKey(identifier, RequestIP(r))

	if limiter.TooManyAttempts(key, h.app.config.LoginRateLimit) {
		return ErrTooManyAttempts
	}

	return nil
}

func (h *LoginHandler) onSuccess(ctx context.Context, w http.ResponseWriter, r *http.Request, user cauth.Authenticatable, remember bool) {
	if limiter := h.app.limiter; limiter != nil {
		limiter.Clear(ThrottleKey(user.GetAuthIdentifier(), RequestIP(r)))
	}

	if err := h.app.guard.Login(ctx, w, user, remember); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	h.dispatchEvent(ctx, LoginSucceededPayload{
		User:     user,
		Remember: remember,
	})

	h.app.responder.LoginResponse(w, r)
}

func (h *LoginHandler) onFailure(ctx context.Context, w http.ResponseWriter, r *http.Request, identifier string) {
	if limiter := h.app.limiter; limiter != nil {
		key := ThrottleKey(identifier, RequestIP(r))
		limiter.Hit(key, h.app.config.LoginRateDecay)
	}

	h.dispatchEvent(ctx, LoginFailedPayload{
		Identifier: identifier,
	})

	http.Error(w, ErrInvalidCredentials.Error(), http.StatusUnprocessableEntity)
}

func (h *LoginHandler) dispatchEvent(ctx context.Context, payload any) {
	if h.app.events == nil {
		return
	}

	_, _ = h.app.events.Dispatch(ctx, payload)
}
