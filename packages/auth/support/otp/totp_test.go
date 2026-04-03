package otp

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateSecretValidateAndURL(t *testing.T) {
	t.Parallel()

	secret, err := GenerateSecret()

	if err != nil {
		t.Fatalf("GenerateSecret: %v", err)
	}

	if len(secret) != 32 {
		t.Fatalf("unexpected secret length: %d", len(secret))
	}

	now := time.Date(2026, 4, 3, 10, 0, 0, 0, time.UTC)
	code, err := Code(secret, now)

	if err != nil {
		t.Fatalf("Code: %v", err)
	}

	if !Validate(secret, code, now, 0) {
		t.Fatal("expected code to validate")
	}

	if Validate(secret, "000000", now, 0) {
		t.Fatal("expected invalid code to fail")
	}

	url := OTPAuthURL("gollin", "user@example.com", secret)

	if !strings.Contains(url, "otpauth://totp/") || !strings.Contains(url, "issuer=gollin") {
		t.Fatalf("unexpected otpauth url: %s", url)
	}
}

func TestCodeRejectsInvalidSecret(t *testing.T) {
	t.Parallel()

	if _, err := Code("not-base32", time.Now()); err == nil {
		t.Fatal("expected invalid secret to fail")
	}
}
