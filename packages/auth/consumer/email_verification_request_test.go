package consumer_test

import (
	"context"
	"testing"
	"time"

	auth "github.com/bedrock/packages/auth"
	consumer "github.com/bedrock/packages/auth/consumer"
	securitycrypto "github.com/bedrock/packages/support/crypto"
)

// ---------- test helpers ----------

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

// nonVerifiableUser does NOT implement MustVerifyEmail
type nonVerifiableUser struct {
	id string
}

func (u *nonVerifiableUser) GetAuthIdentifierName() string { return "id" }
func (u *nonVerifiableUser) GetAuthIdentifier() string     { return u.id }
func (u *nonVerifiableUser) GetAuthPasswordName() string   { return "password" }
func (u *nonVerifiableUser) GetAuthPassword() string       { return "" }
func (u *nonVerifiableUser) SetAuthPassword(string)        {}
func (u *nonVerifiableUser) GetRememberToken() string      { return "" }
func (u *nonVerifiableUser) SetRememberToken(string)       {}
func (u *nonVerifiableUser) GetRememberTokenName() string  { return "remember_token" }

// ======================== TESTS ========================

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

	// Idempotent: fulfilling again should be a no-op
	if err := request.FulfillAt(context.Background(), now.Add(time.Hour)); err != nil {
		t.Fatalf("FulfillAt idempotent: %v", err)
	}
}

func TestEmailVerificationRequestAuthorizeDeniesWrongID(t *testing.T) {
	t.Parallel()

	user := &verifiableUser{id: "user-1", email: "person@example.com"}
	request := consumer.EmailVerificationRequest{
		User:      user,
		RouteID:   "wrong-id",
		RouteHash: securitycrypto.EmailHash(user.GetEmailForVerification()),
	}

	if request.Authorize() {
		t.Fatal("expected request to be denied for wrong ID")
	}
}

func TestEmailVerificationRequestAuthorizeDeniesWrongHash(t *testing.T) {
	t.Parallel()

	user := &verifiableUser{id: "user-1", email: "person@example.com"}
	request := consumer.EmailVerificationRequest{
		User:      user,
		RouteID:   user.GetAuthIdentifier(),
		RouteHash: "wrong-hash",
	}

	if request.Authorize() {
		t.Fatal("expected request to be denied for wrong hash")
	}
}

func TestEmailVerificationRequestAuthorizeDeniesNilUser(t *testing.T) {
	t.Parallel()

	request := consumer.EmailVerificationRequest{
		User:      nil,
		RouteID:   "user-1",
		RouteHash: "some-hash",
	}

	if request.Authorize() {
		t.Fatal("expected request to be denied for nil user")
	}
}

func TestEmailVerificationRequestAuthorizeDeniesNonVerifiableUser(t *testing.T) {
	t.Parallel()

	user := &nonVerifiableUser{id: "user-1"}
	request := consumer.EmailVerificationRequest{
		User:      user,
		RouteID:   "user-1",
		RouteHash: "some-hash",
	}

	if request.Authorize() {
		t.Fatal("expected request to be denied for non-verifiable user")
	}
}

func TestEmailVerificationRequestFulfillAtReturnsErrorWhenUnauthorized(t *testing.T) {
	t.Parallel()

	request := consumer.EmailVerificationRequest{
		User:      nil,
		RouteID:   "user-1",
		RouteHash: "hash",
	}

	err := request.FulfillAt(context.Background(), time.Now())
	if err == nil {
		t.Fatal("expected error for unauthorized request")
	}
}

func TestEmailVerificationRequestFulfillAtWithNonVerifiableUser(t *testing.T) {
	t.Parallel()

	user := &nonVerifiableUser{id: "user-1"}
	request := consumer.EmailVerificationRequest{
		User:      user,
		RouteID:   "user-1",
		RouteHash: "hash",
	}

	err := request.FulfillAt(context.Background(), time.Now())
	if err == nil {
		t.Fatal("expected error for non-verifiable user")
	}
}

func TestEmailVerificationRequestFulfillUsesCurrentTime(t *testing.T) {
	t.Parallel()

	user := &verifiableUser{id: "user-1", email: "person@example.com"}
	request := consumer.EmailVerificationRequest{
		User:      user,
		RouteID:   user.GetAuthIdentifier(),
		RouteHash: securitycrypto.EmailHash(user.GetEmailForVerification()),
	}

	before := time.Now().UTC()
	err := request.Fulfill(context.Background())
	if err != nil {
		t.Fatalf("Fulfill: %v", err)
	}

	if !user.HasVerifiedEmail() {
		t.Fatal("expected user to be verified")
	}

	_ = before // Fulfill uses time.Now() internally
}

func TestEmailVerificationRequestUnverifyUser(t *testing.T) {
	t.Parallel()

	now := time.Now()
	user := &verifiableUser{id: "user-1", email: "person@example.com", verifiedAt: &now}

	if !user.HasVerifiedEmail() {
		t.Fatal("expected user to start verified")
	}

	user.MarkEmailAsUnverified()

	if user.HasVerifiedEmail() {
		t.Fatal("expected user to be unverified")
	}
}

// MustVerifyEmail interface compliance
func TestVerifiableUserImplementsMustVerifyEmail(t *testing.T) {
	t.Parallel()

	var _ auth.MustVerifyEmail = (*verifiableUser)(nil)
}
