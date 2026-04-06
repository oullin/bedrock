package consumer_test

import (
	"context"
	"testing"
	"time"

	consumer "github.com/bedrock/packages/auth/consumer"
	securitycrypto "github.com/bedrock/packages/support/crypto"
)

type verifiableUser struct {
	id         string
	email      string
	verifiedAt *time.Time
}

func (u *verifiableUser) GetAuthIdentifierName() string { return "id" }
func (u *verifiableUser) GetAuthIdentifier() string     { return u.id }
func (u *verifiableUser) GetAuthPasswordName() string   { return "password" }
func (u *verifiableUser) GetAuthPassword() string       { return "" }
func (u *verifiableUser) SetAuthPassword(string)        {}
func (u *verifiableUser) GetRememberToken() string      { return "" }
func (u *verifiableUser) SetRememberToken(string)       {}
func (u *verifiableUser) GetRememberTokenName() string  { return "remember_token" }
func (u *verifiableUser) HasVerifiedEmail() bool        { return u.verifiedAt != nil }
func (u *verifiableUser) MarkEmailAsVerified(at time.Time) {
	u.verifiedAt = &at
}

func (u *verifiableUser) MarkEmailAsUnverified() {
	u.verifiedAt = nil
}

func (u *verifiableUser) GetEmailForVerification() string {
	return u.email
}

func TestEmailVerificationRequestAuthorizeAndFulfill(t *testing.T) {
	t.Parallel()

	user := &verifiableUser{id: "user-1", email: "person@example.com"}
	request := consumer.EmailVerificationRequest{
		User:      user,
		RouteID:   user.GetAuthIdentifier(),
		RouteHash: securitycrypto.EmailHash(user.GetEmailForVerification()),
	}

	if !request.Authorize() {
		t.Fatal("expected request to authorize")
	}

	now := time.Date(2026, 4, 5, 1, 0, 0, 0, time.UTC)
	if err := request.FulfillAt(context.Background(), now); err != nil {
		t.Fatalf("FulfillAt: %v", err)
	}

	if !user.HasVerifiedEmail() {
		t.Fatal("expected user to be verified")
	}

	if err := request.FulfillAt(context.Background(), now.Add(time.Hour)); err != nil {
		t.Fatalf("FulfillAt idempotent: %v", err)
	}
}
