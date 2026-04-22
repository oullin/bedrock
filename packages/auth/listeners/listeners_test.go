package listeners_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/auth/events"
	"github.com/bedrock/packages/auth/listeners"
	cauth "github.com/bedrock/packages/contracts/auth"
)

// stubUser implements Authenticatable but not MustVerifyEmail.
type stubUser struct {
	id            string
	password      string
	rememberToken string
}

// verifiableUser implements both Authenticatable and EmailVerificationSender.
type verifiableUser struct {
	stubUser
	verified                 bool
	notificationSent         bool
	notificationContext      context.Context
	notificationContextErr   error
	notificationContextValue any
}

type contextKey string

func (u *stubUser) GetAuthIdentifierName() string   { return "id" }
func (u *stubUser) GetAuthIdentifier() string       { return u.id }
func (u *stubUser) GetAuthPasswordName() string     { return "password" }
func (u *stubUser) GetAuthPassword() string         { return u.password }
func (u *stubUser) SetAuthPassword(password string) { u.password = password }
func (u *stubUser) GetRememberToken() string        { return u.rememberToken }
func (u *stubUser) SetRememberToken(token string)   { u.rememberToken = token }
func (u *stubUser) GetRememberTokenName() string    { return "remember_token" }

func (u *verifiableUser) HasVerifiedEmail() bool          { return u.verified }
func (u *verifiableUser) MarkEmailAsVerified(_ time.Time) { u.verified = true }
func (u *verifiableUser) MarkEmailAsUnverified()          { u.verified = false }
func (u *verifiableUser) GetEmailForVerification() string { return "test@example.com" }
func (u *verifiableUser) SendEmailVerificationNotification(ctx context.Context) {
	u.notificationSent = true
	u.notificationContext = ctx
	u.notificationContextErr = ctx.Err()
	u.notificationContextValue = ctx.Value(contextKey("notification"))
}

// Port of Illuminate\Tests\Auth\AuthListenersSendEmailVerificationNotificationHandleFunctionTest::testWillExecuted
func TestSendEmailVerificationNotification_UnverifiedUser(t *testing.T) {
	user := &verifiableUser{
		stubUser: stubUser{id: "1"},
		verified: false,
	}

	var _ cauth.EmailVerificationNotificationSender = user

	listener := &listeners.SendEmailVerificationNotification{}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), contextKey("notification"), "verification"))
	cancel()
	listener.Handle(ctx, events.Registered{User: user})

	if !user.notificationSent {
		t.Error("expected notification to be sent for unverified user")
	}

	if user.notificationContext != ctx {
		t.Error("expected notification context to match listener context")
	}

	if user.notificationContextErr != context.Canceled {
		t.Errorf("notification context err = %v, want %v", user.notificationContextErr, context.Canceled)
	}

	if user.notificationContextValue != "verification" {
		t.Errorf("notification context value = %v, want %q", user.notificationContextValue, "verification")
	}
}

// Port of Illuminate\Tests\Auth\AuthListenersSendEmailVerificationNotificationHandleFunctionTest::testHasVerifiedEmailAsTrue
func TestSendEmailVerificationNotification_AlreadyVerified(t *testing.T) {
	user := &verifiableUser{
		stubUser: stubUser{id: "1"},
		verified: true,
	}

	listener := &listeners.SendEmailVerificationNotification{}
	listener.Handle(context.Background(), events.Registered{User: user})

	if user.notificationSent {
		t.Error("notification should not be sent for already-verified user")
	}
}

// Port of Illuminate\Tests\Auth\AuthListenersSendEmailVerificationNotificationHandleFunctionTest::testUserIsNotInstanceOfMustVerifyEmail
func TestSendEmailVerificationNotification_NonVerifiableUser(t *testing.T) {
	user := &stubUser{id: "1"}

	listener := &listeners.SendEmailVerificationNotification{}
	// Should not panic when user does not implement MustVerifyEmail.
	listener.Handle(context.Background(), events.Registered{User: user})
}
