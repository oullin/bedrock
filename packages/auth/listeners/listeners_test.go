package listeners_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/auth/events"
	"github.com/bedrock/packages/auth/listeners"
)

// stubUser implements Authenticatable but not MustVerifyEmail.
type stubUser struct {
	id            any
	password      string
	rememberToken string
}

func (u *stubUser) GetAuthIdentifierName() string  { return "id" }
func (u *stubUser) GetAuthIdentifier() any          { return u.id }
func (u *stubUser) GetAuthPassword() string         { return u.password }
func (u *stubUser) GetRememberToken() string        { return u.rememberToken }
func (u *stubUser) SetRememberToken(token string)   { u.rememberToken = token }
func (u *stubUser) GetRememberTokenName() string    { return "remember_token" }

// verifiableUser implements both Authenticatable and MustVerifyEmail.
type verifiableUser struct {
	stubUser
	verified          bool
	notificationSent  bool
}

func (u *verifiableUser) HasVerifiedEmail() bool                { return u.verified }
func (u *verifiableUser) MarkEmailAsVerified() error            { u.verified = true; return nil }
func (u *verifiableUser) SendEmailVerificationNotification()    { u.notificationSent = true }
func (u *verifiableUser) GetEmailForVerification() string       { return "test@example.com" }

func TestSendEmailVerificationNotification_UnverifiedUser(t *testing.T) {
	user := &verifiableUser{
		stubUser: stubUser{id: 1},
		verified: false,
	}

	listener := &listeners.SendEmailVerificationNotification{}
	listener.Handle(context.Background(), events.Registered{User: user})

	if !user.notificationSent {
		t.Error("expected notification to be sent for unverified user")
	}
}

func TestSendEmailVerificationNotification_AlreadyVerified(t *testing.T) {
	user := &verifiableUser{
		stubUser: stubUser{id: 1},
		verified: true,
	}

	listener := &listeners.SendEmailVerificationNotification{}
	listener.Handle(context.Background(), events.Registered{User: user})

	if user.notificationSent {
		t.Error("notification should not be sent for already-verified user")
	}
}

func TestSendEmailVerificationNotification_NonVerifiableUser(t *testing.T) {
	user := &stubUser{id: 1}

	listener := &listeners.SendEmailVerificationNotification{}
	// Should not panic when user does not implement MustVerifyEmail.
	listener.Handle(context.Background(), events.Registered{User: user})
}
