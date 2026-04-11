package fortify

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// --- test doubles ---

type stubGuard struct{}

type stubProvider struct{}

type stubHasher struct{}

type stubResponder struct{}

type stubCreatesNewUsers struct{}

type stubUpdatesProfile struct{}

type stubUpdatesPasswords struct{}

type stubBroker struct{}

type stubResets struct{}

type stubVerifier struct{}

type stubConfirms struct{}

func (s *stubGuard) Name() string { return "web" }
func (s *stubGuard) AuthenticateRequest(_ context.Context, _ http.ResponseWriter, _ *http.Request) (cauth.Authenticatable, error) {
	return nil, nil
}
func (s *stubGuard) Login(_ context.Context, _ http.ResponseWriter, _ cauth.Authenticatable, _ bool) error {
	return nil
}
func (s *stubGuard) LoginWithPendingTwoFactor(_ context.Context, _ http.ResponseWriter, _ cauth.Authenticatable) error {
	return nil
}
func (s *stubGuard) Logout(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (s *stubProvider) RetrieveByID(_ context.Context, _ string) (cauth.Authenticatable, error) {
	return nil, nil
}
func (s *stubProvider) RetrieveByToken(_ context.Context, _ string, _ string) (cauth.Authenticatable, error) {
	return nil, nil
}
func (s *stubProvider) RetrieveByCredentials(_ context.Context, _ map[string]string) (cauth.Authenticatable, error) {
	return nil, nil
}
func (s *stubProvider) UpdateRememberToken(_ context.Context, _ cauth.Authenticatable, _ string) error {
	return nil
}
func (s *stubProvider) ValidateCredentials(_ context.Context, _ cauth.Authenticatable, _ map[string]string) (bool, error) {
	return false, nil
}
func (s *stubProvider) RehashPasswordIfRequired(_ context.Context, _ cauth.Authenticatable, _ map[string]string, _ bool) error {
	return nil
}

func (s *stubHasher) Hash(_ context.Context, _ string) (string, error) { return "hashed", nil }
func (s *stubHasher) Check(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}
func (s *stubHasher) NeedsRehash(_ string) bool { return false }

func (s *stubResponder) LoginResponse(_ http.ResponseWriter, _ *http.Request)                     {}
func (s *stubResponder) LogoutResponse(_ http.ResponseWriter, _ *http.Request)                    {}
func (s *stubResponder) RegisterResponse(_ http.ResponseWriter, _ *http.Request)                  {}
func (s *stubResponder) PasswordResetLinkSentResponse(_ http.ResponseWriter, _ *http.Request)     {}
func (s *stubResponder) PasswordResetResponse(_ http.ResponseWriter, _ *http.Request)             {}
func (s *stubResponder) PasswordUpdateResponse(_ http.ResponseWriter, _ *http.Request)            {}
func (s *stubResponder) PasswordConfirmResponse(_ http.ResponseWriter, _ *http.Request)           {}
func (s *stubResponder) ProfileInformationUpdatedResponse(_ http.ResponseWriter, _ *http.Request) {}
func (s *stubResponder) EmailVerificationSentResponse(_ http.ResponseWriter, _ *http.Request)     {}
func (s *stubResponder) TwoFactorChallengeResponse(_ http.ResponseWriter, _ *http.Request)        {}
func (s *stubResponder) TwoFactorEnabledResponse(_ http.ResponseWriter, _ *http.Request)          {}
func (s *stubResponder) TwoFactorDisabledResponse(_ http.ResponseWriter, _ *http.Request)         {}

func (s *stubCreatesNewUsers) Create(_ context.Context, _ map[string]string) (cauth.Authenticatable, error) {
	return nil, nil
}

func (s *stubUpdatesProfile) Update(_ context.Context, _ cauth.Authenticatable, _ map[string]string) error {
	return nil
}

func (s *stubUpdatesPasswords) Update(_ context.Context, _ cauth.Authenticatable, _ map[string]string) error {
	return nil
}

func (s *stubBroker) SendResetLink(_ context.Context, _ map[string]string) error { return nil }
func (s *stubBroker) Reset(_ context.Context, _ map[string]string, _ func(cauth.Authenticatable, string) error) error {
	return nil
}

func (s *stubResets) Reset(_ context.Context, _ cauth.Authenticatable, _ string) error { return nil }

func (s *stubVerifier) SendVerificationNotification(_ context.Context, _ cauth.Authenticatable) error {
	return nil
}
func (s *stubVerifier) Verify(_ context.Context, _ string, _ string) error { return nil }

func (s *stubConfirms) Confirm(_ context.Context, _ cauth.Authenticatable, _ string) error {
	return nil
}

// --- tests ---

func TestBuildMinimalFortify(t *testing.T) {
	config := DefaultConfig()
	config.Features = Features{}

	f, err := NewBuilder().
		WithConfig(config).
		WithGuard(&stubGuard{}).
		WithProvider(&stubProvider{}).
		WithHasher(&stubHasher{}).
		WithResponder(&stubResponder{}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Guard().Name() != "web" {
		t.Fatal("expected guard name to be web")
	}
}

func TestBuildFullFortify(t *testing.T) {
	f, err := NewBuilder().
		WithConfig(DefaultConfig()).
		WithGuard(&stubGuard{}).
		WithProvider(&stubProvider{}).
		WithHasher(&stubHasher{}).
		WithResponder(&stubResponder{}).
		WithCreateUser(&stubCreatesNewUsers{}).
		WithUpdateProfile(&stubUpdatesProfile{}).
		WithUpdatePassword(&stubUpdatesPasswords{}).
		WithBroker(&stubBroker{}).
		WithResetPassword(&stubResets{}).
		WithVerifier(&stubVerifier{}).
		WithConfirmPassword(&stubConfirms{}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Config().Features.Registration != true {
		t.Fatal("expected registration to be enabled")
	}
}

func TestBuildFailsWithoutGuard(t *testing.T) {
	config := DefaultConfig()
	config.Features = Features{}

	_, err := NewBuilder().
		WithConfig(config).
		WithProvider(&stubProvider{}).
		WithHasher(&stubHasher{}).
		WithResponder(&stubResponder{}).
		Build()

	if err == nil {
		t.Fatal("expected error for missing guard")
	}

	if !strings.Contains(err.Error(), "guard is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildFailsWithMissingAction(t *testing.T) {
	_, err := NewBuilder().
		WithConfig(DefaultConfig()).
		WithGuard(&stubGuard{}).
		WithProvider(&stubProvider{}).
		WithHasher(&stubHasher{}).
		WithResponder(&stubResponder{}).
		Build()

	if err == nil {
		t.Fatal("expected error for missing actions")
	}

	if !strings.Contains(err.Error(), "CreatesNewUsers") {
		t.Fatalf("expected CreatesNewUsers error, got: %v", err)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Guard != "web" {
		t.Fatalf("expected guard web, got %s", config.Guard)
	}

	if config.PasswordTimeout != 3*time.Hour {
		t.Fatalf("expected 3h timeout, got %v", config.PasswordTimeout)
	}

	if config.LoginRateLimit != 5 {
		t.Fatalf("expected rate limit 5, got %d", config.LoginRateLimit)
	}

	if config.IdentifierField != "email" {
		t.Fatalf("expected email identifier, got %s", config.IdentifierField)
	}
}

func TestDefaultFeaturesAllEnabled(t *testing.T) {
	f := DefaultFeatures()

	if !f.Registration || !f.ResetPasswords || !f.EmailVerification ||
		!f.UpdateProfileInformation || !f.UpdatePasswords ||
		!f.TwoFactorAuthentication || !f.ConfirmPassword {
		t.Fatal("expected all features enabled by default")
	}
}
