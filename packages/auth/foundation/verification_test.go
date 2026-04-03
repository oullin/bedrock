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

func TestVerificationServiceSendAndVerify(t *testing.T) {
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

	if len(parts) != 2 {
		t.Fatalf("unexpected path: %s", parsed.Path)
	}

	verified, err := service.Verify(context.Background(), parts[0], parts[1], mustInt64(t, parsed.Query().Get("expires")), parsed.Query().Get("signature"))

	if err != nil {
		t.Fatalf("Verify: %v", err)
	}

	verifiable := verified.(auth.MustVerifyEmail)

	if !verifiable.HasVerifiedEmail() {
		t.Fatal("expected verified email")
	}
}

func mustInt64(t *testing.T, value string) int64 {
	t.Helper()
	parsed, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		t.Fatalf("ParseInt: %v", err)
	}

	return parsed
}
