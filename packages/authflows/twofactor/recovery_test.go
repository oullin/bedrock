package twofactor

import (
	"strings"
	"testing"
)

func TestGenerateRecoveryCodes(t *testing.T) {
	codes, err := GenerateRecoveryCodes(DefaultRecoveryCodeCount)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(codes) != DefaultRecoveryCodeCount {
		t.Fatalf("expected %d codes, got %d", DefaultRecoveryCodeCount, len(codes))
	}

	for _, code := range codes {
		if !strings.Contains(code, "-") {
			t.Fatalf("code should contain dash separator: %s", code)
		}

		if len(code) != 21 {
			t.Fatalf("expected 21 char code (10-10), got %d: %s", len(code), code)
		}
	}
}

func TestGenerateRecoveryCodesUnique(t *testing.T) {
	codes, _ := GenerateRecoveryCodes(DefaultRecoveryCodeCount)
	seen := make(map[string]bool)

	for _, code := range codes {
		if seen[code] {
			t.Fatalf("duplicate recovery code: %s", code)
		}

		seen[code] = true
	}
}

func TestGenerateRecoveryCodesDefaultCount(t *testing.T) {
	codes, err := GenerateRecoveryCodes(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(codes) != DefaultRecoveryCodeCount {
		t.Fatalf("expected %d codes with default, got %d", DefaultRecoveryCodeCount, len(codes))
	}
}

func TestValidateRecoveryCode(t *testing.T) {
	codes := []string{"aaaaa-bbbbb", "ccccc-ddddd", "eeeee-fffff"}

	idx := ValidateRecoveryCode("ccccc-ddddd", codes)
	if idx != 1 {
		t.Fatalf("expected index 1, got %d", idx)
	}
}

func TestValidateRecoveryCodeNotFound(t *testing.T) {
	codes := []string{"aaaaa-bbbbb", "ccccc-ddddd"}

	idx := ValidateRecoveryCode("xxxxx-yyyyy", codes)
	if idx != -1 {
		t.Fatalf("expected -1, got %d", idx)
	}
}

func TestValidateRecoveryCodeTrimsWhitespace(t *testing.T) {
	codes := []string{"aaaaa-bbbbb"}

	idx := ValidateRecoveryCode("  aaaaa-bbbbb  ", codes)
	if idx != 0 {
		t.Fatalf("expected index 0 with trimmed whitespace, got %d", idx)
	}
}

func TestConsumeRecoveryCode(t *testing.T) {
	codes := []string{"aaa", "bbb", "ccc"}

	result := ConsumeRecoveryCode(codes, 1)
	if len(result) != 2 {
		t.Fatalf("expected 2 codes, got %d", len(result))
	}

	if result[0] != "aaa" || result[1] != "ccc" {
		t.Fatalf("unexpected codes after consume: %v", result)
	}
}

func TestConsumeRecoveryCodeInvalidIndex(t *testing.T) {
	codes := []string{"aaa", "bbb"}

	result := ConsumeRecoveryCode(codes, -1)
	if len(result) != 2 {
		t.Fatal("should return original codes for invalid index")
	}

	result = ConsumeRecoveryCode(codes, 5)
	if len(result) != 2 {
		t.Fatal("should return original codes for out-of-range index")
	}
}
