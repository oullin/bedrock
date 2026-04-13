package inception

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/inception/twofactor"
)

// --- 2FA test user ---

type testTwoFactorUser struct {
	testUser
	secret       string
	codes        []string
	confirmedAt  *time.Time
	emailAddress string
}

func (u *testTwoFactorUser) IsTwoFactorEnabled() bool             { return u.twoFA }
func (u *testTwoFactorUser) SetTwoFactorEnabled(e bool)           { u.twoFA = e }
func (u *testTwoFactorUser) GetTwoFactorSecret() string           { return u.secret }
func (u *testTwoFactorUser) SetTwoFactorSecret(s string)          { u.secret = s }
func (u *testTwoFactorUser) GetTwoFactorRecoveryCodes() []string  { return u.codes }
func (u *testTwoFactorUser) SetTwoFactorRecoveryCodes(c []string) { u.codes = c }
func (u *testTwoFactorUser) GetTwoFactorConfirmedAt() *time.Time  { return u.confirmedAt }
func (u *testTwoFactorUser) SetTwoFactorConfirmedAt(t *time.Time) { u.confirmedAt = t }
func (u *testTwoFactorUser) HasVerifiedEmail() bool               { return true }
func (u *testTwoFactorUser) MarkEmailAsVerified(_ time.Time)      {}
func (u *testTwoFactorUser) MarkEmailAsUnverified()               {}
func (u *testTwoFactorUser) GetEmailForVerification() string      { return u.emailAddress }

// --- helpers ---

func build2FAInception(user cauth.Authenticatable) (*Inception, *testGuard, *testEvents, *testResponder) {
	guard := &testGuard{authenticatedUser: user}
	events := &testEvents{}
	responder := &testResponder{}

	f := buildTestInception(guard, &testProvider{}, events, responder)
	f.config.Features.TwoFactorAuthentication = true

	return f, guard, events, responder
}

func postWithBody(path string, body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

// --- enable 2FA tests ---

func TestEnableTwoFactorSuccess(t *testing.T) {
	user := &testTwoFactorUser{
		testUser: testUser{id: "1", email: "user@example.com"},
	}

	f, _, events, _ := build2FAInception(user)

	handler := NewEnableTwoFactorHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/user/two-factor-authentication", nil)

	handler.ServeHTTP(w, r)

	if user.secret == "" {
		t.Fatal("expected secret to be set")
	}

	if len(user.codes) == 0 {
		t.Fatal("expected recovery codes to be set")
	}

	if !user.twoFA {
		t.Fatal("expected 2FA to be enabled")
	}

	if user.confirmedAt != nil {
		t.Fatal("expected confirmedAt to be nil until confirmation")
	}

	if len(events.dispatched) != 1 {
		t.Fatal("expected TwoFactorEnabled event")
	}

	if s, ok := events.dispatched[0].(string); !ok || s != EventTwoFactorEnabled {
		t.Fatal("expected TwoFactorEnabled event string")
	}
}

func TestEnableTwoFactorDisabledFeature(t *testing.T) {
	f, _, _, _ := build2FAInception(&testTwoFactorUser{})
	f.config.Features.TwoFactorAuthentication = false

	handler := NewEnableTwoFactorHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/user/two-factor-authentication", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestEnableTwoFactorUnauthenticated(t *testing.T) {
	f, _, _, _ := build2FAInception(nil)

	handler := NewEnableTwoFactorHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/user/two-factor-authentication", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// --- confirm 2FA tests ---

func TestConfirmTwoFactorSuccess(t *testing.T) {
	secret, _ := twofactor.GenerateSecret(0)
	code := twofactor.CurrentCode(secret)

	user := &testTwoFactorUser{
		testUser: testUser{id: "1", twoFA: true},
		secret:   secret,
	}

	f, _, events, _ := build2FAInception(user)

	handler := NewConfirmTwoFactorHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/confirmed-two-factor-authentication", "code="+code)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if user.confirmedAt == nil {
		t.Fatal("expected confirmedAt to be set")
	}

	if len(events.dispatched) != 1 {
		t.Fatal("expected TwoFactorConfirmed event")
	}

	if s, ok := events.dispatched[0].(string); !ok || s != EventTwoFactorConfirmed {
		t.Fatal("expected TwoFactorConfirmed event string")
	}
}

func TestConfirmTwoFactorInvalidCode(t *testing.T) {
	user := &testTwoFactorUser{
		testUser: testUser{id: "1", twoFA: true},
		secret:   "JBSWY3DPEHPK3PXP",
	}

	f, _, _, _ := build2FAInception(user)

	handler := NewConfirmTwoFactorHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/confirmed-two-factor-authentication", "code=000000")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

// --- disable 2FA tests ---

func TestDisableTwoFactorSuccess(t *testing.T) {
	now := time.Now()
	user := &testTwoFactorUser{
		testUser:    testUser{id: "1", twoFA: true},
		secret:      "secret",
		codes:       []string{"code1"},
		confirmedAt: &now,
	}

	f, _, events, _ := build2FAInception(user)

	handler := NewDisableTwoFactorHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/user/two-factor-authentication", nil)

	handler.ServeHTTP(w, r)

	if user.twoFA {
		t.Fatal("expected 2FA to be disabled")
	}

	if user.secret != "" {
		t.Fatal("expected secret to be cleared")
	}

	if user.codes != nil {
		t.Fatal("expected recovery codes to be cleared")
	}

	if user.confirmedAt != nil {
		t.Fatal("expected confirmedAt to be cleared")
	}

	if len(events.dispatched) != 1 {
		t.Fatal("expected TwoFactorDisabled event")
	}

	if s, ok := events.dispatched[0].(string); !ok || s != EventTwoFactorDisabled {
		t.Fatal("expected TwoFactorDisabled event string")
	}
}

// --- QR code handler tests ---

func TestTwoFactorQRCodeSuccess(t *testing.T) {
	user := &testTwoFactorUser{
		testUser:     testUser{id: "1", twoFA: true},
		secret:       "JBSWY3DPEHPK3PXP",
		emailAddress: "user@example.com",
	}

	f, _, _, _ := build2FAInception(user)

	handler := NewTwoFactorQRCodeHandler(f, "Bedrock")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/user/two-factor-qr-code", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]string

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !strings.Contains(body["url"], "otpauth://totp/") {
		t.Fatalf("expected otpauth URI, got %s", body["url"])
	}

	if !strings.Contains(body["url"], "JBSWY3DPEHPK3PXP") {
		t.Fatalf("expected secret in URI, got %s", body["url"])
	}
}

func TestTwoFactorQRCodeNoSecret(t *testing.T) {
	user := &testTwoFactorUser{
		testUser: testUser{id: "1"},
		secret:   "",
	}

	f, _, _, _ := build2FAInception(user)

	handler := NewTwoFactorQRCodeHandler(f, "Bedrock")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/user/two-factor-qr-code", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// --- recovery codes handler tests ---

func TestRecoveryCodesGet(t *testing.T) {
	user := &testTwoFactorUser{
		testUser: testUser{id: "1"},
		codes:    []string{"code1-code1", "code2-code2"},
	}

	f, _, _, _ := build2FAInception(user)

	handler := NewTwoFactorRecoveryCodesHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/user/two-factor-recovery-codes", nil)

	handler.ServeHTTP(w, r)

	var codes []string

	if err := json.NewDecoder(w.Body).Decode(&codes); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(codes) != 2 {
		t.Fatalf("expected 2 codes, got %d", len(codes))
	}
}

func TestRecoveryCodesRegenerate(t *testing.T) {
	user := &testTwoFactorUser{
		testUser: testUser{id: "1"},
		codes:    []string{"old-code"},
	}

	f, _, _, _ := build2FAInception(user)

	handler := NewTwoFactorRecoveryCodesHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/user/two-factor-recovery-codes", nil)

	handler.ServeHTTP(w, r)

	if len(user.codes) != 8 {
		t.Fatalf("expected 8 new codes, got %d", len(user.codes))
	}
}
