package otp

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateSecretAndValidate(t *testing.T) {
	t.Parallel()

	secret, err := GenerateSecret()

	if err != nil {
		t.Fatalf("GenerateSecret: %v", err)
	}

	if len(secret) != 32 {
		t.Fatalf("unexpected secret length: %d", len(secret))
	}

	now := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)
	code, err := Code(secret, now)

	if err != nil {
		t.Fatalf("Code: %v", err)
	}

	if !Validate(secret, code, now, 1) {
		t.Fatal("expected code to validate")
	}

	if Validate(secret, "000000", now, 1) {
		t.Fatal("expected bad code to fail")
	}
}

func TestOTPAuthURL(t *testing.T) {
	t.Parallel()

	url := OTPAuthURL("gollin", "user@example.com", "ABC123")

	if !strings.Contains(url, "otpauth://totp/") || !strings.Contains(url, "issuer=gollin") {
		t.Fatalf("unexpected otpauth url: %s", url)
	}
}
