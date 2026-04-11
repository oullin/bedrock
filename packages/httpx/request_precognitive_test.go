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
	raw.Header.Set("HTTPPreview", "true")
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

func TestIsAttemptingHTTPPreview(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("HTTPPreview", "true")
	raw.Header.Set("HTTPPreview-Validate-Only", "name,email")
	req := httpx.NewRequest(raw)

	if !req.IsAttemptingHTTPPreview() {
		t.Fatal("expected IsAttemptingHTTPPreview to return true")
	}
}

func TestIsAttemptingHTTPPreviewWithoutValidateOnly(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("HTTPPreview", "true")
	req := httpx.NewRequest(raw)

	if req.IsAttemptingHTTPPreview() {
		t.Fatal("expected IsAttemptingHTTPPreview to return false without Validate-Only")
	}
}

func TestPrecognitiveValidateOnly(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("HTTPPreview", "true")
	raw.Header.Set("HTTPPreview-Validate-Only", "name, email, age")
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
	raw.Header.Set("HTTPPreview", "true")
	raw.Header.Set("HTTPPreview-Validate-Only", "name,address")
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
		t.Fatalf("expected all rules when not precognitive, got %d", len(filtered))
	}
}
