package inception

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/inception/twofactor"
)

func TestTwoFactorChallengeWithTOTP(t *testing.T) {
	secret, _ := twofactor.GenerateSecret(0)
	code := twofactor.CurrentCode(secret)

	user := &testTwoFactorUser{
		testUser: testUser{id: "1", twoFA: true},
		secret:   secret,
	}

	f, guard, events, responder := build2FAInception(user)

	handler := NewTwoFactorChallengeHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/two-factor-challenge", "code="+code)

	handler.ServeHTTP(w, r)

	if guard.loggedIn == nil {
		t.Fatal("expected user to be logged in after TOTP challenge")
	}

	if !responder.loginCalled {
		t.Fatal("expected login response")
	}

	if len(events.dispatched) != 1 {
		t.Fatal("expected TwoFactorChallenge event")
	}

	if _, ok := events.dispatched[0].(LoginSucceededPayload); !ok {
		t.Fatalf("expected LoginSucceededPayload event, got %T", events.dispatched[0])
	}
}

func TestTwoFactorChallengeWithRecoveryCode(t *testing.T) {
	user := &testTwoFactorUser{
		testUser: testUser{id: "1", twoFA: true},
		secret:   "JBSWY3DPEHPK3PXP",
		codes:    []string{"aaaa-bbbb", "cccc-dddd", "eeee-ffff"},
	}

	f, guard, events, _ := build2FAInception(user)

	handler := NewTwoFactorChallengeHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/two-factor-challenge", "recovery_code=cccc-dddd")

	handler.ServeHTTP(w, r)

	if guard.loggedIn == nil {
		t.Fatal("expected user to be logged in after recovery code")
	}

	if len(user.codes) != 2 {
		t.Fatalf("expected 2 remaining codes, got %d", len(user.codes))
	}

	// Should have RecoveryCodeUsed (string) + LoginSucceededPayload events
	if len(events.dispatched) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events.dispatched))
	}

	if s, ok := events.dispatched[0].(string); !ok || s != EventRecoveryCodeUsed {
		t.Fatalf("expected RecoveryCodeUsed event string, got %T", events.dispatched[0])
	}
}

func TestTwoFactorChallengeInvalidTOTP(t *testing.T) {
	user := &testTwoFactorUser{
		testUser: testUser{id: "1", twoFA: true},
		secret:   "JBSWY3DPEHPK3PXP",
	}

	f, guard, _, _ := build2FAInception(user)

	handler := NewTwoFactorChallengeHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/two-factor-challenge", "code=000000")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}

	if guard.loggedIn != nil {
		t.Fatal("user should not be logged in with invalid TOTP")
	}
}

func TestTwoFactorChallengeInvalidRecoveryCode(t *testing.T) {
	user := &testTwoFactorUser{
		testUser: testUser{id: "1", twoFA: true},
		codes:    []string{"aaaa-bbbb"},
	}

	f, guard, _, _ := build2FAInception(user)

	handler := NewTwoFactorChallengeHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/two-factor-challenge", "recovery_code=xxxx-yyyy")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}

	if guard.loggedIn != nil {
		t.Fatal("user should not be logged in with invalid recovery code")
	}
}

func TestTwoFactorChallengeNoCodeProvided(t *testing.T) {
	user := &testTwoFactorUser{
		testUser: testUser{id: "1", twoFA: true},
	}

	f, _, _, _ := build2FAInception(user)

	handler := NewTwoFactorChallengeHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/two-factor-challenge", "")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestTwoFactorChallengeUnauthenticated(t *testing.T) {
	f, _, _, _ := build2FAInception(nil)

	handler := NewTwoFactorChallengeHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/two-factor-challenge", "code=123456")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestTwoFactorChallengeFeatureDisabled(t *testing.T) {
	f, _, _, _ := build2FAInception(&testTwoFactorUser{})
	f.config.Features.TwoFactorAuthentication = false

	handler := NewTwoFactorChallengeHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/two-factor-challenge", "code=123456")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
