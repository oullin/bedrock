package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx"
)

func TestIsPrecognitive(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("Precognition", "true")
	req := httpx.NewRequest(raw)

	if !req.IsPrecognitive() {
		t.Fatal("expected IsPrecognitive to return true")
	}
}

func TestIsPrecognitiveWithoutHeader(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	req := httpx.NewRequest(raw)

	if req.IsPrecognitive() {
		t.Fatal("expected IsPrecognitive to return false")
	}
}

func TestIsAttemptingPrecognition(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("Precognition", "true")
	req := httpx.NewRequest(raw)

	if !req.IsAttemptingPrecognition() {
		t.Fatal("expected IsAttemptingPrecognition to return true")
	}
}

func TestIsAttemptingPrecognitionWithValidateOnly(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("Precognition", "true")
	raw.Header.Set("Precognition-Validate-Only", "name,email")
	req := httpx.NewRequest(raw)

	if !req.IsAttemptingPrecognition() {
		t.Fatal("expected IsAttemptingPrecognition to return true with both headers")
	}
}

func TestIsAttemptingPrecognitionRequiresExactTrue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value string
		want  bool
	}{
		{"true", true},
		{"1", false},
		{"yes", false},
		{"on", false},
		{"True", false},
		{"false", false},
		{"", false},
	}

	for _, tt := range tests {
		raw := httptest.NewRequest(http.MethodPost, "/", nil)

		if tt.value != "" {
			raw.Header.Set("Precognition", tt.value)
		}

		req := httpx.NewRequest(raw)
		got := req.IsAttemptingPrecognition()

		if got != tt.want {
			t.Errorf("IsAttemptingPrecognition(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestPrecognitiveValidateOnly(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("Precognition", "true")
	raw.Header.Set("Precognition-Validate-Only", "name, email, age")
	req := httpx.NewRequest(raw)

	fields := req.PrecognitiveValidateOnly()

	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}

	if fields[0] != "name" || fields[1] != "email" || fields[2] != "age" {
		t.Fatalf("unexpected fields: %v", fields)
	}
}

func TestPrecognitiveValidateOnlyEmpty(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	req := httpx.NewRequest(raw)

	if fields := req.PrecognitiveValidateOnly(); fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
}

func TestFilterPrecognitiveRules(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("Precognition-Validate-Only", "name,address.*")
	req := httpx.NewRequest(raw)

	rules := map[string]any{
		"name":           "required",
		"email":          "required|email",
		"address.street": "required",
		"address.city":   "required",
		"phone":          "nullable",
	}

	filtered := req.FilterPrecognitiveRules(rules)

	if _, ok := filtered["name"]; !ok {
		t.Fatal("expected name rule to be kept")
	}

	if _, ok := filtered["address.street"]; !ok {
		t.Fatal("expected address.street rule to be kept")
	}

	if _, ok := filtered["address.city"]; !ok {
		t.Fatal("expected address.city rule to be kept")
	}

	if _, ok := filtered["email"]; ok {
		t.Fatal("expected email rule to be filtered out")
	}

	if _, ok := filtered["phone"]; ok {
		t.Fatal("expected phone rule to be filtered out")
	}
}

func TestFilterPrecognitiveRulesExactMatch(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("Precognition-Validate-Only", "name")
	req := httpx.NewRequest(raw)

	rules := map[string]any{
		"name":       "required",
		"name.first": "required",
		"email":      "required",
	}

	filtered := req.FilterPrecognitiveRules(rules)

	if _, ok := filtered["name"]; !ok {
		t.Fatal("expected name rule to be kept")
	}

	if _, ok := filtered["name.first"]; ok {
		t.Fatal("expected name.first to be filtered out (exact match only)")
	}

	if _, ok := filtered["email"]; ok {
		t.Fatal("expected email to be filtered out")
	}
}

func TestFilterPrecognitiveRulesWildcard(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("Precognition-Validate-Only", "user.profile.*")
	req := httpx.NewRequest(raw)

	rules := map[string]any{
		"user.profile.name":       "required",
		"user.profile.bio":        "required",
		"user.profile.avatar.url": "required",
		"user.email":              "required",
	}

	filtered := req.FilterPrecognitiveRules(rules)

	if _, ok := filtered["user.profile.name"]; !ok {
		t.Fatal("expected user.profile.name to match user.profile.*")
	}

	if _, ok := filtered["user.profile.bio"]; !ok {
		t.Fatal("expected user.profile.bio to match user.profile.*")
	}

	if _, ok := filtered["user.profile.avatar.url"]; ok {
		t.Fatal("expected user.profile.avatar.url NOT to match user.profile.* (extra dot segment)")
	}

	if _, ok := filtered["user.email"]; ok {
		t.Fatal("expected user.email NOT to match user.profile.*")
	}
}

func TestFilterPrecognitiveRulesNotPrecognitive(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	req := httpx.NewRequest(raw)

	rules := map[string]any{
		"name":  "required",
		"email": "required",
	}

	filtered := req.FilterPrecognitiveRules(rules)

	if len(filtered) != 2 {
		t.Fatalf("expected all rules when no Precognition-Validate-Only header, got %d", len(filtered))
	}
}
