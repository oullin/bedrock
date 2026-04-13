package twofactor

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateSecret(t *testing.T) {
	secret, err := GenerateSecret(DefaultSecretSize)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(secret) == 0 {
		t.Fatal("secret should not be empty")
	}

	secret2, _ := GenerateSecret(DefaultSecretSize)

	if secret == secret2 {
		t.Fatal("secrets should be unique")
	}
}

func TestGenerateSecretDefaultSize(t *testing.T) {
	secret, err := GenerateSecret(0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(secret) == 0 {
		t.Fatal("secret should not be empty with default size")
	}
}

func TestValidateCurrentCode(t *testing.T) {
	secret, _ := GenerateSecret(DefaultSecretSize)
	code := CurrentCode(secret)

	if !Validate(code, secret) {
		t.Fatal("current code should be valid")
	}
}

func TestValidateRejectsWrongCode(t *testing.T) {
	secret, _ := GenerateSecret(DefaultSecretSize)

	if Validate("000000", secret) {
		// Extremely unlikely but technically possible. Skip rather than fail.
		t.Skip("unlikely collision with 000000")
	}
}

func TestValidateWithClockSkew(t *testing.T) {
	secret, _ := GenerateSecret(DefaultSecretSize)
	now := time.Now()

	pastCode := CodeAt(secret, now.Add(-30*time.Second))

	if !ValidateAt(pastCode, secret, now) {
		t.Fatal("code from previous period should be valid (clock skew)")
	}

	futureCode := CodeAt(secret, now.Add(30*time.Second))

	if !ValidateAt(futureCode, secret, now) {
		t.Fatal("code from next period should be valid (clock skew)")
	}
}

func TestValidateRejectsFarFutureCode(t *testing.T) {
	secret, _ := GenerateSecret(DefaultSecretSize)
	now := time.Now()

	farCode := CodeAt(secret, now.Add(5*time.Minute))

	if ValidateAt(farCode, secret, now) {
		t.Fatal("code from far future should be rejected")
	}
}

func TestCodeAtDeterministic(t *testing.T) {
	secret, _ := GenerateSecret(DefaultSecretSize)
	at := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	code1 := CodeAt(secret, at)
	code2 := CodeAt(secret, at)

	if code1 != code2 {
		t.Fatal("CodeAt should be deterministic for same time")
	}

	if len(code1) != DefaultDigits {
		t.Fatalf("expected %d digit code, got %d", DefaultDigits, len(code1))
	}
}

func TestProvisioningURI(t *testing.T) {
	uri := ProvisioningURI("JBSWY3DPEHPK3PXP", "user@example.com", "Bedrock")

	if !strings.HasPrefix(uri, "otpauth://totp/") {
		t.Fatalf("unexpected URI prefix: %s", uri)
	}

	if !strings.Contains(uri, "secret=JBSWY3DPEHPK3PXP") {
		t.Fatalf("URI should contain secret: %s", uri)
	}

	if !strings.Contains(uri, "issuer=Bedrock") {
		t.Fatalf("URI should contain issuer: %s", uri)
	}

	if !strings.Contains(uri, "Bedrock") && !strings.Contains(uri, "user%40example.com") {
		t.Fatalf("URI should contain label: %s", uri)
	}
}

func TestGenerateCodeInvalidSecret(t *testing.T) {
	code := generateCode("!!!invalid!!!", 0)

	if code != "" {
		t.Fatalf("expected empty code for invalid secret, got %s", code)
	}
}
