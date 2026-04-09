package authflows

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bedrock/packages/authflows/twofactor"
)

// EnableTwoFactorHandler handles POST /user/two-factor-authentication.
type EnableTwoFactorHandler struct {
	authflows *AuthFlows
}

// NewEnableTwoFactorHandler creates a new enable two-factor handler.
func NewEnableTwoFactorHandler(f *AuthFlows) *EnableTwoFactorHandler {
	return &EnableTwoFactorHandler{authflows: f}
}

// ServeHTTP enables two-factor authentication for the authenticated user.
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

	tfa, ok := user.(TwoFactorAuthenticatable)
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
		_ = h.authflows.events.Dispatch(ctx, Event{Name: EventTwoFactorEnabled})
	}

	h.authflows.responder.TwoFactorEnabledResponse(w, r)
}

// ConfirmTwoFactorHandler handles POST /user/confirmed-two-factor-authentication.
type ConfirmTwoFactorHandler struct {
	authflows *AuthFlows
}

// NewConfirmTwoFactorHandler creates a new confirm two-factor handler.
func NewConfirmTwoFactorHandler(f *AuthFlows) *ConfirmTwoFactorHandler {
	return &ConfirmTwoFactorHandler{authflows: f}
}

// ServeHTTP confirms two-factor authentication by validating a TOTP code.
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

	tfa, ok := user.(TwoFactorAuthenticatable)
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
		_ = h.authflows.events.Dispatch(ctx, Event{Name: EventTwoFactorConfirmed})
	}

	w.WriteHeader(http.StatusOK)
}

// DisableTwoFactorHandler handles DELETE /user/two-factor-authentication.
type DisableTwoFactorHandler struct {
	authflows *AuthFlows
}

// NewDisableTwoFactorHandler creates a new disable two-factor handler.
func NewDisableTwoFactorHandler(f *AuthFlows) *DisableTwoFactorHandler {
	return &DisableTwoFactorHandler{authflows: f}
}

// ServeHTTP disables two-factor authentication for the authenticated user.
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

	tfa, ok := user.(TwoFactorAuthenticatable)
	if !ok {
		http.Error(w, "user does not support two-factor authentication", http.StatusBadRequest)
		return
	}

	tfa.SetTwoFactorEnabled(false)
	tfa.SetTwoFactorSecret("")
	tfa.SetTwoFactorRecoveryCodes(nil)
	tfa.SetTwoFactorConfirmedAt(nil)

	if h.authflows.events != nil {
		_ = h.authflows.events.Dispatch(ctx, Event{Name: EventTwoFactorDisabled})
	}

	h.authflows.responder.TwoFactorDisabledResponse(w, r)
}

// TwoFactorQRCodeHandler handles GET /user/two-factor-qr-code.
type TwoFactorQRCodeHandler struct {
	authflows *AuthFlows
	issuer  string
}

// NewTwoFactorQRCodeHandler creates a new QR code handler.
func NewTwoFactorQRCodeHandler(f *AuthFlows, issuer string) *TwoFactorQRCodeHandler {
	return &TwoFactorQRCodeHandler{authflows: f, issuer: issuer}
}

// ServeHTTP returns the TOTP provisioning URI for QR code rendering.
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

	tfa, ok := user.(TwoFactorAuthenticatable)
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
	if verifiable, ok := user.(MustVerifyEmail); ok {
		email = verifiable.GetEmailForVerification()
	}

	uri := twofactor.ProvisioningURI(secret, email, h.issuer)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"svg": uri,
		"url": uri,
	})
}

// TwoFactorRecoveryCodesHandler handles GET/POST /user/two-factor-recovery-codes.
type TwoFactorRecoveryCodesHandler struct {
	authflows *AuthFlows
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

	tfa, ok := user.(TwoFactorAuthenticatable)
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
