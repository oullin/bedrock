package authflows

import (
	"encoding/json"
	"net/http"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/authflows/twofactor"
)

// EnableTwoFactorHandler handles POST /user/two-factor-authentication.
type EnableTwoFactorHandler struct {
	authflows *AuthFlows
}

// NewEnableTwoFactorHandler creates a new enable two-factor handler.

// ServeHTTP enables two-factor authentication for the authenticated user.

// ConfirmTwoFactorHandler handles POST /user/confirmed-two-factor-authentication.
type ConfirmTwoFactorHandler struct {
	authflows *AuthFlows
}

// NewConfirmTwoFactorHandler creates a new confirm two-factor handler.

// ServeHTTP confirms two-factor authentication by validating a TOTP code.

// DisableTwoFactorHandler handles DELETE /user/two-factor-authentication.
type DisableTwoFactorHandler struct {
	authflows *AuthFlows
}

// NewDisableTwoFactorHandler creates a new disable two-factor handler.

// ServeHTTP disables two-factor authentication for the authenticated user.

// TwoFactorQRCodeHandler handles GET /user/two-factor-qr-code.
type TwoFactorQRCodeHandler struct {
	authflows *AuthFlows
	issuer  string
}

// NewTwoFactorQRCodeHandler creates a new QR code handler.

// ServeHTTP returns the TOTP provisioning URI for QR code rendering.

// TwoFactorRecoveryCodesHandler handles GET/POST /user/two-factor-recovery-codes.
type TwoFactorRecoveryCodesHandler struct {
	authflows *AuthFlows
}

func NewEnableTwoFactorHandler(f *AuthFlows) *EnableTwoFactorHandler {
	return &EnableTwoFactorHandler{authflows: f}
}

func (h *EnableTwoFactorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.TwoFactorAuthentication {
		http.Error(w, "two-factor authentication is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.authflows.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	tfa, ok := user.(cauth.TwoFactorAuthenticatable)

	if !ok {
		http.Error(w, "user does not support two-factor authentication", http.StatusBadRequest)

		return
	}

	secret, err := twofactor.GenerateSecret(0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	tfa.SetTwoFactorSecret(secret)
	tfa.SetTwoFactorEnabled(true)
	tfa.SetTwoFactorConfirmedAt(nil)

	codes, err := twofactor.GenerateRecoveryCodes(0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	tfa.SetTwoFactorRecoveryCodes(codes)

	if h.authflows.events != nil {
		_ = h.authflows.events.Dispatch(ctx, EventTwoFactorEnabled)
	}

	h.authflows.responder.TwoFactorEnabledResponse(w, r)
}

func NewConfirmTwoFactorHandler(f *AuthFlows) *ConfirmTwoFactorHandler {
	return &ConfirmTwoFactorHandler{authflows: f}
}

func (h *ConfirmTwoFactorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.TwoFactorAuthentication {
		http.Error(w, "two-factor authentication is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.authflows.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	tfa, ok := user.(cauth.TwoFactorAuthenticatable)

	if !ok {
		http.Error(w, "user does not support two-factor authentication", http.StatusBadRequest)

		return
	}

	input := RequestInput(r, "code")
	code := input["code"]

	if !twofactor.Validate(code, tfa.GetTwoFactorSecret()) {
		http.Error(w, "invalid two-factor code", http.StatusUnprocessableEntity)

		return
	}

	now := time.Now()
	tfa.SetTwoFactorConfirmedAt(&now)

	if h.authflows.events != nil {
		_ = h.authflows.events.Dispatch(ctx, EventTwoFactorConfirmed)
	}

	w.WriteHeader(http.StatusOK)
}

func NewDisableTwoFactorHandler(f *AuthFlows) *DisableTwoFactorHandler {
	return &DisableTwoFactorHandler{authflows: f}
}

func (h *DisableTwoFactorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.TwoFactorAuthentication {
		http.Error(w, "two-factor authentication is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.authflows.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	tfa, ok := user.(cauth.TwoFactorAuthenticatable)

	if !ok {
		http.Error(w, "user does not support two-factor authentication", http.StatusBadRequest)

		return
	}

	tfa.SetTwoFactorEnabled(false)
	tfa.SetTwoFactorSecret("")
	tfa.SetTwoFactorRecoveryCodes(nil)
	tfa.SetTwoFactorConfirmedAt(nil)

	if h.authflows.events != nil {
		_ = h.authflows.events.Dispatch(ctx, EventTwoFactorDisabled)
	}

	h.authflows.responder.TwoFactorDisabledResponse(w, r)
}

func NewTwoFactorQRCodeHandler(f *AuthFlows, issuer string) *TwoFactorQRCodeHandler {
	return &TwoFactorQRCodeHandler{authflows: f, issuer: issuer}
}

func (h *TwoFactorQRCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.TwoFactorAuthentication {
		http.Error(w, "two-factor authentication is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.authflows.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	tfa, ok := user.(cauth.TwoFactorAuthenticatable)

	if !ok {
		http.Error(w, "user does not support two-factor authentication", http.StatusBadRequest)

		return
	}

	secret := tfa.GetTwoFactorSecret()

	if secret == "" {
		http.Error(w, "two-factor authentication is not enabled", http.StatusBadRequest)

		return
	}

	email := ""

	if verifiable, ok := user.(cauth.MustVerifyEmail); ok {
		email = verifiable.GetEmailForVerification()
	}

	uri := twofactor.ProvisioningURI(secret, email, h.issuer)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"svg": uri,
		"url": uri,
	})
}

// NewTwoFactorRecoveryCodesHandler creates a new recovery codes handler.
func NewTwoFactorRecoveryCodesHandler(f *AuthFlows) *TwoFactorRecoveryCodesHandler {
	return &TwoFactorRecoveryCodesHandler{authflows: f}
}

// ServeHTTP handles both GET (list) and POST (regenerate) for recovery codes.
func (h *TwoFactorRecoveryCodesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.TwoFactorAuthentication {
		http.Error(w, "two-factor authentication is disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.authflows.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	tfa, ok := user.(cauth.TwoFactorAuthenticatable)

	if !ok {
		http.Error(w, "user does not support two-factor authentication", http.StatusBadRequest)

		return
	}

	if r.Method == http.MethodPost {
		codes, err := twofactor.GenerateRecoveryCodes(0)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		tfa.SetTwoFactorRecoveryCodes(codes)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tfa.GetTwoFactorRecoveryCodes())
}
