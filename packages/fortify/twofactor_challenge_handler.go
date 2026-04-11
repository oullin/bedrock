package fortify

import (
	"errors"
	"net/http"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/fortify/twofactor"
)

// ErrInvalidTwoFactorCode is returned when a TOTP or recovery code is invalid.

// TwoFactorChallengeHandler handles POST /two-factor-challenge requests.
// It completes the login flow for users with 2FA enabled by validating
// either a TOTP code or a recovery code.
type TwoFactorChallengeHandler struct {
	fortify *Fortify
}

var ErrInvalidTwoFactorCode = errors.New("fortify: invalid two-factor authentication code")

// NewTwoFactorChallengeHandler creates a new two-factor challenge handler.
func NewTwoFactorChallengeHandler(f *Fortify) *TwoFactorChallengeHandler {
	return &TwoFactorChallengeHandler{fortify: f}
}

// ServeHTTP handles the two-factor challenge request.
func (h *TwoFactorChallengeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.fortify.config.Features.TwoFactorAuthentication {
		http.Error(w, "two-factor authentication is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()
	input := RequestInput(r, "code", "recovery_code")
	code := input["code"]
	recoveryCode := input["recovery_code"]

	user, err := h.fortify.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	tfa, ok := user.(cauth.TwoFactorAuthenticatable)

	if !ok {
		http.Error(w, "user does not support two-factor authentication", http.StatusBadRequest)

		return
	}

	if code != "" {
		if !twofactor.Validate(code, tfa.GetTwoFactorSecret()) {
			h.onFailure(w)

			return
		}
	} else if recoveryCode != "" {
		codes := tfa.GetTwoFactorRecoveryCodes()
		idx := twofactor.ValidateRecoveryCode(recoveryCode, codes)

		if idx == -1 {
			h.onFailure(w)

			return
		}

		tfa.SetTwoFactorRecoveryCodes(twofactor.ConsumeRecoveryCode(codes, idx))

		if h.fortify.events != nil {
			_ = h.fortify.events.Dispatch(ctx, EventRecoveryCodeUsed)
		}
	} else {
		http.Error(w, "a code or recovery_code is required", http.StatusUnprocessableEntity)

		return
	}

	remember := RequestBool(r, "remember")

	if err := h.fortify.guard.Login(ctx, w, user, remember); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if h.fortify.events != nil {
		_ = h.fortify.events.Dispatch(ctx, LoginSucceededPayload{User: user, Remember: remember})
	}

	h.fortify.responder.LoginResponse(w, r)
}

func (h *TwoFactorChallengeHandler) onFailure(w http.ResponseWriter) {
	http.Error(w, ErrInvalidTwoFactorCode.Error(), http.StatusUnprocessableEntity)
}
