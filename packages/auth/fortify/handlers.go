package fortify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	auth "github.com/gollin/packages/auth"
)

// Responder renders Fortify-style HTTP responses.
type Responder interface {
	Success(w http.ResponseWriter, status int, payload any)
	Error(w http.ResponseWriter, err error)
}

// JSONResponder is the default JSON response renderer.
type JSONResponder struct{}

// Success renders a JSON success payload.
func (JSONResponder) Success(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// Error renders an error response with auth-aware status codes.
func (JSONResponder) Error(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	var validationErr *auth.ValidationError
	var throttleErr *auth.ThrottleError
	switch {
	case errors.As(err, &validationErr):
		status = http.StatusUnprocessableEntity
	case errors.Is(err, auth.ErrUnauthorized):
		status = http.StatusUnauthorized
	case errors.Is(err, auth.ErrInvalidCredentials):
		status = http.StatusUnauthorized
	case errors.Is(err, auth.ErrInvalidToken):
		status = http.StatusUnprocessableEntity
	case errors.Is(err, auth.ErrTokenExpired):
		status = http.StatusGone
	case errors.Is(err, auth.ErrTwoFactorInvalid):
		status = http.StatusUnauthorized
	case errors.Is(err, auth.ErrTwoFactorRequired):
		status = http.StatusConflict
	case errors.Is(err, auth.ErrPasswordConfirmationRequired):
		status = http.StatusForbidden
	case errors.Is(err, auth.ErrUserExists):
		status = http.StatusConflict
	case errors.As(err, &throttleErr):
		status = http.StatusTooManyRequests
	case errors.Is(err, auth.ErrEmailVerificationInvalid):
		status = http.StatusUnprocessableEntity
	}

	auth.WriteHTTPError(w, status, err)
}

// RouteOptions customizes route registration.
type RouteOptions struct {
	Prefix    string
	Responder Responder
}

// Handler serves built-in Fortify-compatible HTTP handlers.
type Handler struct {
	manager   *auth.Manager
	responder Responder
}

// NewHandler creates a new route handler.
func NewHandler(manager *auth.Manager, options RouteOptions) *Handler {
	responder := options.Responder
	if responder == nil {
		responder = JSONResponder{}
	}

	return &Handler{
		manager:   manager,
		responder: responder,
	}
}

// RegisterRoutes attaches all Fortify-inspired routes to a mux.
func RegisterRoutes(mux *http.ServeMux, manager *auth.Manager, options RouteOptions) {
	handler := NewHandler(manager, options)
	prefix := strings.TrimSuffix(options.Prefix, "/")

	mux.Handle("POST "+prefix+"/register", http.HandlerFunc(handler.register))
	mux.Handle("POST "+prefix+"/login", http.HandlerFunc(handler.login))
	mux.Handle("POST "+prefix+"/logout", manager.RequireAuthenticated(http.HandlerFunc(handler.logout)))
	mux.Handle("POST "+prefix+"/forgot-password", http.HandlerFunc(handler.forgotPassword))
	mux.Handle("POST "+prefix+"/reset-password", http.HandlerFunc(handler.resetPassword))
	mux.Handle("GET "+prefix+"/verify-email/{userID}/{hash}", http.HandlerFunc(handler.verifyEmail))
	mux.Handle("POST "+prefix+"/email/verification-notification", manager.RequireAuthenticated(http.HandlerFunc(handler.resendVerification)))
	mux.Handle("POST "+prefix+"/user/confirm-password", manager.RequireAuthenticated(http.HandlerFunc(handler.confirmPassword)))
	mux.Handle("POST "+prefix+"/user/two-factor-authentication", manager.RequirePasswordConfirmed(http.HandlerFunc(handler.enableTwoFactor)))
	mux.Handle("DELETE "+prefix+"/user/two-factor-authentication", manager.RequirePasswordConfirmed(http.HandlerFunc(handler.disableTwoFactor)))
	mux.Handle("POST "+prefix+"/two-factor-challenge", http.HandlerFunc(handler.twoFactorChallenge))
	mux.Handle("GET "+prefix+"/user/two-factor-recovery-codes", manager.RequirePasswordConfirmed(http.HandlerFunc(handler.recoveryCodes)))
	mux.Handle("POST "+prefix+"/user/two-factor-recovery-codes", manager.RequirePasswordConfirmed(http.HandlerFunc(handler.regenerateRecoveryCodes)))
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var input auth.RegisterInput
	if err := decodeJSON(r, &input); err != nil {
		h.responder.Error(w, err)
		return
	}

	user, err := h.manager.Register(r.Context(), input)
	if err != nil {
		h.responder.Error(w, err)
		return
	}

	h.responder.Success(w, http.StatusCreated, map[string]any{
		"status": "registered",
		"user":   user,
	})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var input auth.LoginInput
	if err := decodeJSON(r, &input); err != nil {
		h.responder.Error(w, err)
		return
	}

	result, err := h.manager.Login(r.Context(), input, throttleKey(r, input.Email))
	if err != nil {
		h.responder.Error(w, err)
		return
	}

	guard := h.manager.DefaultGuard()
	guard.SetSessionCookie(w, result.Session)
	if result.RememberToken != "" {
		guard.SetRememberCookie(w, result.User.ID, result.RememberToken)
	}

	payload := map[string]any{
		"status":             "authenticated",
		"requiresTwoFactor":  result.RequiresTwoFactor,
		"passwordConfirmed":  result.PasswordConfirmed,
		"authenticatedState": result.AuthenticatedState,
		"twoFactorPending":   result.Session.PendingTwoFactor,
	}
	if result.User != nil {
		payload["user"] = result.User
	}
	if result.RequiresTwoFactor {
		payload["status"] = "two_factor_required"
	}

	h.responder.Success(w, http.StatusOK, payload)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.CurrentUser(r)
	session, _ := auth.CurrentSession(r)
	if err := h.manager.Logout(r.Context(), user, session); err != nil {
		h.responder.Error(w, err)
		return
	}

	h.manager.DefaultGuard().ClearSessionCookies(w)
	h.responder.Success(w, http.StatusOK, map[string]any{"status": "logged_out"})
}

func (h *Handler) forgotPassword(w http.ResponseWriter, r *http.Request) {
	var input auth.ForgotPasswordInput
	if err := decodeJSON(r, &input); err != nil {
		h.responder.Error(w, err)
		return
	}

	if err := h.manager.SendPasswordResetLink(r.Context(), input, throttleKey(r, input.Email)); err != nil {
		h.responder.Error(w, err)
		return
	}

	h.responder.Success(w, http.StatusAccepted, map[string]any{"status": "password_reset_link_sent"})
}

func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request) {
	var input auth.ResetPasswordInput
	if err := decodeJSON(r, &input); err != nil {
		h.responder.Error(w, err)
		return
	}

	user, err := h.manager.ResetPassword(r.Context(), input)
	if err != nil {
		h.responder.Error(w, err)
		return
	}

	h.responder.Success(w, http.StatusOK, map[string]any{
		"status": "password_reset",
		"user":   user,
	})
}

func (h *Handler) verifyEmail(w http.ResponseWriter, r *http.Request) {
	expiresAt, err := strconv.ParseInt(r.URL.Query().Get("expires"), 10, 64)
	if err != nil {
		h.responder.Error(w, auth.ErrInvalidToken)
		return
	}

	user, err := h.manager.VerifyEmail(
		r.Context(),
		r.PathValue("userID"),
		r.PathValue("hash"),
		expiresAt,
		r.URL.Query().Get("signature"),
	)
	if err != nil {
		h.responder.Error(w, err)
		return
	}

	h.responder.Success(w, http.StatusOK, map[string]any{
		"status": "email_verified",
		"user":   user,
	})
}

func (h *Handler) resendVerification(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.CurrentUser(r)
	if err := h.manager.SendEmailVerification(r.Context(), user); err != nil {
		h.responder.Error(w, err)
		return
	}

	h.responder.Success(w, http.StatusAccepted, map[string]any{"status": "verification_notification_sent"})
}

func (h *Handler) confirmPassword(w http.ResponseWriter, r *http.Request) {
	var input auth.ConfirmPasswordInput
	if err := decodeJSON(r, &input); err != nil {
		h.responder.Error(w, err)
		return
	}

	user, _ := auth.CurrentUser(r)
	session, _ := auth.CurrentSession(r)
	if err := h.manager.ConfirmPassword(r.Context(), user, session, input); err != nil {
		h.responder.Error(w, err)
		return
	}

	h.responder.Success(w, http.StatusOK, map[string]any{"status": "password_confirmed"})
}

func (h *Handler) enableTwoFactor(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.CurrentUser(r)
	session, _ := auth.CurrentSession(r)

	result, err := h.manager.EnableTwoFactor(r.Context(), user, session)
	if err != nil {
		h.responder.Error(w, err)
		return
	}

	h.responder.Success(w, http.StatusCreated, map[string]any{
		"status":        "two_factor_enabled",
		"secret":        result.Secret,
		"recoveryCodes": result.RecoveryCodes,
		"otpAuthUrl":    result.OTPAuthURL,
	})
}

func (h *Handler) disableTwoFactor(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.CurrentUser(r)
	session, _ := auth.CurrentSession(r)
	if err := h.manager.DisableTwoFactor(r.Context(), user, session); err != nil {
		h.responder.Error(w, err)
		return
	}

	h.responder.Success(w, http.StatusOK, map[string]any{"status": "two_factor_disabled"})
}

func (h *Handler) recoveryCodes(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.CurrentUser(r)
	session, _ := auth.CurrentSession(r)
	codes, err := h.manager.RecoveryCodes(r.Context(), user, session)
	if err != nil {
		h.responder.Error(w, err)
		return
	}

	h.responder.Success(w, http.StatusOK, map[string]any{
		"status":        "recovery_codes_loaded",
		"recoveryCodes": codes,
	})
}

func (h *Handler) regenerateRecoveryCodes(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.CurrentUser(r)
	session, _ := auth.CurrentSession(r)
	codes, err := h.manager.RegenerateRecoveryCodes(r.Context(), user, session)
	if err != nil {
		h.responder.Error(w, err)
		return
	}

	h.responder.Success(w, http.StatusCreated, map[string]any{
		"status":        "recovery_codes_regenerated",
		"recoveryCodes": codes,
	})
}

func (h *Handler) twoFactorChallenge(w http.ResponseWriter, r *http.Request) {
	var input auth.TwoFactorChallengeInput
	if err := decodeJSON(r, &input); err != nil {
		h.responder.Error(w, err)
		return
	}

	session, _, err := h.manager.Authenticate(r.Context(), w, r)
	if err != nil || session == nil || !session.PendingTwoFactor {
		h.responder.Error(w, auth.ErrUnauthorized)
		return
	}

	result, err := h.manager.ChallengeTwoFactor(r.Context(), session, input)
	if err != nil {
		h.responder.Error(w, err)
		return
	}

	guard := h.manager.DefaultGuard()
	guard.SetSessionCookie(w, result.Session)
	if result.RememberToken != "" {
		guard.SetRememberCookie(w, result.User.ID, result.RememberToken)
	}

	h.responder.Success(w, http.StatusOK, map[string]any{
		"status": "authenticated",
		"user":   result.User,
	})
}

func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return &auth.ValidationError{Fields: map[string]string{"body": "request body is required"}}
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return &auth.ValidationError{Fields: map[string]string{"body": fmt.Sprintf("invalid JSON: %v", err)}}
	}
	return nil
}

func throttleKey(r *http.Request, identifier string) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	return strings.TrimSpace(strings.ToLower(identifier)) + "|" + host
}
