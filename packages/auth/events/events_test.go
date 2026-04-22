package events_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/auth"
	"github.com/bedrock/packages/auth/events"
	cauth "github.com/bedrock/packages/contracts/auth"
)

type resettableEventUser struct {
	*auth.GenericUser
	email string
}

var _ cauth.ResettableAuthenticatable = (*resettableEventUser)(nil)

func (u *resettableEventUser) GetEmailForPasswordReset() string { return u.email }

func TestAttemptingFields(t *testing.T) {
	e := events.Attempting{
		Guard:       "web",
		Credentials: map[string]string{"email": "a@b.com"},
		Remember:    true,
	}

	if e.Guard != "web" {
		t.Errorf("Guard = %q, want %q", e.Guard, "web")
	}

	if e.Credentials["email"] != "a@b.com" {
		t.Errorf("Credentials[email] = %v, want %q", e.Credentials["email"], "a@b.com")
	}

	if !e.Remember {
		t.Error("Remember = false, want true")
	}
}

func TestAuthenticatedFields(t *testing.T) {
	e := events.Authenticated{Guard: "api", User: nil}

	if e.Guard != "api" {
		t.Errorf("Guard = %q, want %q", e.Guard, "api")
	}
}

func TestFailedNilUser(t *testing.T) {
	e := events.Failed{
		Guard:       "web",
		User:        nil,
		Credentials: map[string]string{"email": "x@y.com"},
	}

	if e.User != nil {
		t.Error("User should be nil for unknown user failures")
	}
}

func TestLoginFields(t *testing.T) {
	e := events.Login{Guard: "web", User: nil, Remember: false}

	if e.Guard != "web" {
		t.Errorf("Guard = %q, want %q", e.Guard, "web")
	}

	if e.Remember {
		t.Error("Remember = true, want false")
	}
}

func TestLogoutFields(t *testing.T) {
	e := events.Logout{Guard: "web", User: nil}

	if e.Guard != "web" {
		t.Errorf("Guard = %q, want %q", e.Guard, "web")
	}
}

func TestLockoutFields(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	e := events.Lockout{Request: req}

	if e.Request != req {
		t.Error("Request mismatch")
	}
}

func TestCurrentDeviceLogoutFields(t *testing.T) {
	e := events.CurrentDeviceLogout{Guard: "web", User: nil}

	if e.Guard != "web" {
		t.Errorf("Guard = %q, want %q", e.Guard, "web")
	}
}

func TestOtherDeviceLogoutFields(t *testing.T) {
	e := events.OtherDeviceLogout{Guard: "web", User: nil}

	if e.Guard != "web" {
		t.Errorf("Guard = %q, want %q", e.Guard, "web")
	}
}

func TestPasswordResetFields(t *testing.T) {
	e := events.PasswordReset{User: nil}

	if e.User != nil {
		t.Error("expected nil user")
	}
}

func TestPasswordResetLinkSentAcceptsResettableAuthenticatable(t *testing.T) {
	user := &resettableEventUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}
	e := events.PasswordResetLinkSent{User: user}

	if e.User.GetAuthIdentifier() != "1" {
		t.Errorf("GetAuthIdentifier = %q, want %q", e.User.GetAuthIdentifier(), "1")
	}

	if e.User.GetEmailForPasswordReset() != "test@example.com" {
		t.Errorf("GetEmailForPasswordReset = %q, want %q", e.User.GetEmailForPasswordReset(), "test@example.com")
	}
}

func TestRegisteredFields(t *testing.T) {
	e := events.Registered{User: nil}

	if e.User != nil {
		t.Error("expected nil user")
	}
}

func TestValidatedFields(t *testing.T) {
	e := events.Validated{Guard: "web", User: nil}

	if e.Guard != "web" {
		t.Errorf("Guard = %q, want %q", e.Guard, "web")
	}
}

func TestVerifiedAcceptsMustVerifyEmail(t *testing.T) {
	e := events.Verified{User: nil}

	if e.User != nil {
		t.Error("expected nil user")
	}
}
