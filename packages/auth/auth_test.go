package auth

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gollin/packages/auth/internal/totp"
)

type fixedClock struct {
	now time.Time
}

func (c *fixedClock) Now() time.Time {
	return c.now
}

func (c *fixedClock) Advance(duration time.Duration) {
	c.now = c.now.Add(duration)
}

func newTestManager(t *testing.T) (*Manager, *InMemoryMailer, *fixedClock) {
	t.Helper()

	clock := &fixedClock{now: time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)}
	mailer := &InMemoryMailer{}
	manager, err := NewManager(Config{
		BaseURL: "https://example.test",
		Cookies: CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
			HTTPOnly:     true,
		},
		SigningKey: []byte("test-signing-key"),
	}, Dependencies{
		Mailer: mailer,
		Clock:  clock,
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	return manager, mailer, clock
}

func TestDefaultPasswordHasher(t *testing.T) {
	t.Parallel()

	hasher := DefaultPasswordHasher{}
	ctx := context.Background()

	cases := []struct {
		name        string
		password    string
		compareWith string
		wantErr     bool
	}{
		{name: "matches original password", password: "secret-pass", compareWith: "secret-pass"},
		{name: "rejects wrong password", password: "secret-pass", compareWith: "wrong-pass", wantErr: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			encoded, err := hasher.Hash(ctx, tc.password)
			if err != nil {
				t.Fatalf("Hash: %v", err)
			}

			err = hasher.Compare(ctx, encoded, tc.compareWith)
			if tc.wantErr && err == nil {
				t.Fatal("Compare succeeded unexpectedly")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("Compare: %v", err)
			}
		})
	}
}

func TestPasswordResetFlow(t *testing.T) {
	t.Parallel()

	manager, mailer, _ := newTestManager(t)
	ctx := context.Background()

	user, err := manager.Register(ctx, RegisterInput{
		Email:                "reset@example.com",
		Password:             "password-123",
		PasswordConfirmation: "password-123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	if err := manager.SendPasswordResetLink(ctx, ForgotPasswordInput{Email: user.Email}, "reset@example.com|127.0.0.1"); err != nil {
		t.Fatalf("SendPasswordResetLink: %v", err)
	}

	messages := mailer.Messages()
	if len(messages) < 2 {
		t.Fatalf("expected reset email, got %d messages", len(messages))
	}

	token := messages[len(messages)-1].Metadata["token"]
	if token == "" {
		t.Fatal("expected reset token in mail metadata")
	}

	if _, err := manager.ResetPassword(ctx, ResetPasswordInput{
		Email:                user.Email,
		Token:                token,
		Password:             "new-password-123",
		PasswordConfirmation: "new-password-123",
	}); err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}

	if _, err := manager.Login(ctx, LoginInput{Email: user.Email, Password: "password-123"}, "old|127.0.0.1"); err == nil {
		t.Fatal("expected old password login to fail")
	}

	if _, err := manager.Login(ctx, LoginInput{Email: user.Email, Password: "new-password-123"}, "new|127.0.0.1"); err != nil {
		t.Fatalf("Login with new password: %v", err)
	}
}

func TestVerifyEmailLink(t *testing.T) {
	t.Parallel()

	manager, mailer, _ := newTestManager(t)
	ctx := context.Background()

	user, err := manager.Register(ctx, RegisterInput{
		Email:                "verify@example.com",
		Password:             "password-123",
		PasswordConfirmation: "password-123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	link := mailer.Messages()[0].Metadata["link"]
	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("Parse verification link: %v", err)
	}

	parts := strings.Split(strings.TrimPrefix(parsed.Path, "/verify-email/"), "/")
	if len(parts) != 2 {
		t.Fatalf("unexpected verification path: %s", parsed.Path)
	}

	user, err = manager.VerifyEmail(ctx, parts[0], parts[1], mustParseInt64(t, parsed.Query().Get("expires")), parsed.Query().Get("signature"))
	if err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
	if user.EmailVerifiedAt == nil {
		t.Fatal("expected email to be verified")
	}
}

func TestTwoFactorChallengeAndRecoveryCodes(t *testing.T) {
	t.Parallel()

	manager, _, clock := newTestManager(t)
	ctx := context.Background()

	user, err := manager.Register(ctx, RegisterInput{
		Email:                "2fa@example.com",
		Password:             "password-123",
		PasswordConfirmation: "password-123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	loginResult, err := manager.Login(ctx, LoginInput{Email: user.Email, Password: "password-123"}, "login|127.0.0.1")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if err := manager.ConfirmPassword(ctx, user, loginResult.Session, ConfirmPasswordInput{Password: "password-123"}); err != nil {
		t.Fatalf("ConfirmPassword: %v", err)
	}

	enableResult, err := manager.EnableTwoFactor(ctx, user, loginResult.Session)
	if err != nil {
		t.Fatalf("EnableTwoFactor: %v", err)
	}
	if len(enableResult.RecoveryCodes) != 8 {
		t.Fatalf("expected 8 recovery codes, got %d", len(enableResult.RecoveryCodes))
	}

	if err := manager.Logout(ctx, user, loginResult.Session); err != nil {
		t.Fatalf("Logout: %v", err)
	}

	pending, err := manager.Login(ctx, LoginInput{Email: user.Email, Password: "password-123", Remember: true}, "2fa|127.0.0.1")
	if err != nil {
		t.Fatalf("Login with 2FA: %v", err)
	}
	if !pending.RequiresTwoFactor {
		t.Fatal("expected two-factor requirement")
	}

	code, err := totp.Code(enableResult.Secret, clock.Now())
	if err != nil {
		t.Fatalf("totp.Code: %v", err)
	}

	completed, err := manager.ChallengeTwoFactor(ctx, pending.Session, TwoFactorChallengeInput{Code: code})
	if err != nil {
		t.Fatalf("ChallengeTwoFactor code: %v", err)
	}
	if completed.User == nil || completed.Session.PendingTwoFactor {
		t.Fatal("expected authenticated session after TOTP challenge")
	}
	if completed.RememberToken == "" {
		t.Fatal("expected remember token after successful challenge")
	}

	reLogin, err := manager.Login(ctx, LoginInput{Email: user.Email, Password: "password-123"}, "recovery|127.0.0.1")
	if err != nil {
		t.Fatalf("Login for recovery flow: %v", err)
	}

	if _, err := manager.ChallengeTwoFactor(ctx, reLogin.Session, TwoFactorChallengeInput{RecoveryCode: enableResult.RecoveryCodes[0]}); err != nil {
		t.Fatalf("ChallengeTwoFactor recovery: %v", err)
	}

	codes, err := manager.RecoveryCodes(ctx, user, loginResult.Session)
	if err != nil {
		t.Fatalf("RecoveryCodes: %v", err)
	}
	if len(codes) != 7 {
		t.Fatalf("expected one recovery code to be consumed, got %d remaining", len(codes))
	}
}

func TestLoginThrottle(t *testing.T) {
	t.Parallel()

	manager, _, clock := newTestManager(t)
	ctx := context.Background()

	if _, err := manager.Register(ctx, RegisterInput{
		Email:                "throttle@example.com",
		Password:             "password-123",
		PasswordConfirmation: "password-123",
	}); err != nil {
		t.Fatalf("Register: %v", err)
	}

	for attempt := 0; attempt < manager.Config().Throttle.LoginLimit; attempt++ {
		if _, err := manager.Login(ctx, LoginInput{Email: "throttle@example.com", Password: "wrong-password"}, "throttle|127.0.0.1"); err == nil {
			t.Fatalf("attempt %d unexpectedly succeeded", attempt+1)
		}
	}

	if _, err := manager.Login(ctx, LoginInput{Email: "throttle@example.com", Password: "wrong-password"}, "throttle|127.0.0.1"); err == nil {
		t.Fatal("expected throttle error")
	} else {
		var throttleErr *ThrottleError
		if !errors.As(err, &throttleErr) {
			t.Fatalf("expected ThrottleError, got %v", err)
		}
	}

	clock.Advance(manager.Config().Throttle.LoginWindow + time.Second)
	if _, err := manager.Login(ctx, LoginInput{Email: "throttle@example.com", Password: "password-123"}, "throttle|127.0.0.1"); err != nil {
		t.Fatalf("login after throttle window: %v", err)
	}
}

func mustParseInt64(t *testing.T, value string) int64 {
	t.Helper()

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		t.Fatalf("ParseInt(%q): %v", value, err)
	}
	return parsed
}
