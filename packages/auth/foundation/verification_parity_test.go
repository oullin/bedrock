package foundation_test

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/memory"
)

func TestVerificationServiceRejectsBadHashSignatureAndExpiry(t *testing.T) {
	t.Parallel()

	clock := memory.NewFixedClock(time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC))
	users := memory.NewInMemoryUserRepository()
	mailer := &memory.InMemoryMailer{}
	user := &foundation.User{
		ID:        "user-1",
		Name:      "Verify User",
		Email:     "verify@example.com",
		CreatedAt: clock.Now(),
		UpdatedAt: clock.Now(),
	}

	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	service := &foundation.VerificationService{
		Config: auth.Config{
			BaseURL:         "https://example.test",
			VerificationTTL: time.Hour,
		},
		Users:  users,
		Signer: auth.HMACLinkSigner{Key: []byte("signing-key"), Clock: clock},
		Mailer: mailer,
		Clock:  clock,
	}

	if err := service.Send(context.Background(), user); err != nil {
		t.Fatalf("Send: %v", err)
	}

	link := mailer.Messages()[0].Metadata["link"]
	parsed, err := url.Parse(link)

	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	parts := strings.Split(strings.TrimPrefix(parsed.Path, "/email/verify/"), "/")
	expiresAt := mustInt64(t, parsed.Query().Get("expires"))
	signature := parsed.Query().Get("signature")

	if _, err := service.Verify(context.Background(), parts[0], "bad-hash", expiresAt, signature); err != auth.ErrEmailVerificationInvalid {
		t.Fatalf("expected ErrEmailVerificationInvalid, got %v", err)
	}

	if _, err := service.Verify(context.Background(), parts[0], parts[1], expiresAt, "bad-signature"); err != auth.ErrEmailVerificationInvalid {
		t.Fatalf("expected bad signature error, got %v", err)
	}

	clock.Advance(2 * time.Hour)

	if _, err := service.Verify(context.Background(), parts[0], parts[1], expiresAt, signature); err != auth.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestVerificationServiceResendAndAlreadyVerified(t *testing.T) {
	t.Parallel()

	clock := memory.NewFixedClock(time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC))
	users := memory.NewInMemoryUserRepository()
	mailer := &memory.InMemoryMailer{}
	user := &foundation.User{
		ID:        "user-1",
		Name:      "Verify User",
		Email:     "verify@example.com",
		CreatedAt: clock.Now(),
		UpdatedAt: clock.Now(),
	}

	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	service := &foundation.VerificationService{
		Config: auth.Config{
			BaseURL:         "https://example.test",
			VerificationTTL: time.Hour,
		},
		Users:  users,
		Signer: auth.HMACLinkSigner{Key: []byte("signing-key"), Clock: clock},
		Mailer: mailer,
		Clock:  clock,
	}

	if err := service.Send(context.Background(), user); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if err := service.Send(context.Background(), user); err != nil {
		t.Fatalf("Send again: %v", err)
	}

	if len(mailer.Messages()) != 2 {
		t.Fatalf("expected two messages, got %d", len(mailer.Messages()))
	}

	link := mailer.Messages()[0].Metadata["link"]
	parsed, err := url.Parse(link)

	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	parts := strings.Split(strings.TrimPrefix(parsed.Path, "/email/verify/"), "/")
	expiresAt, err := strconv.ParseInt(parsed.Query().Get("expires"), 10, 64)

	if err != nil {
		t.Fatalf("ParseInt: %v", err)
	}

	verified, err := service.Verify(context.Background(), parts[0], parts[1], expiresAt, parsed.Query().Get("signature"))

	if err != nil {
		t.Fatalf("Verify: %v", err)
	}

	if !verified.(auth.MustVerifyEmail).HasVerifiedEmail() {
		t.Fatal("expected user to be verified")
	}
}
