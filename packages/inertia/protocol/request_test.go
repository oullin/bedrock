package protocol_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
)

func TestIsInertiaRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header string
		want   bool
	}{
		{"with header", "true", true},
		{"without header", "", false},
		{"wrong value", "false", false},
		{"with surrounding whitespace", "  true  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", nil)

			if strings.TrimSpace(tt.header) != "" {
				r.Header.Set(protocol.HeaderInertia, tt.header)
			}

			if got := protocol.IsInertiaRequest(r); got != tt.want {
				t.Errorf("IsInertiaRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsHTTPPreviewRequest(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("POST", "/", nil)

	r.Header.Set("HTTPPreview", "true")

	if !protocol.IsHTTPPreviewRequest(r) {
		t.Error("IsHTTPPreviewRequest() = false, want true")
	}
}

func TestIsHTTPPreviewRequest_Missing(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("POST", "/", nil)

	if protocol.IsHTTPPreviewRequest(r) {
		t.Error("IsHTTPPreviewRequest() = true, want false")
	}
}

func TestValidateOnly(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("POST", "/", nil)

	r.Header.Set("Validate-Only", "name,email,phone")

	fields := protocol.ValidateOnly(r)

	if len(fields) != 3 {
		t.Fatalf("len = %d, want 3", len(fields))
	}

	expected := []string{"name", "email", "phone"}

	for i, f := range fields {
		if f != expected[i] {
			t.Errorf("fields[%d] = %q, want %q", i, f, expected[i])
		}
	}
}

func TestValidateOnly_WithSpaces(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("POST", "/", nil)

	r.Header.Set("Validate-Only", " name , email ")

	fields := protocol.ValidateOnly(r)

	if len(fields) != 2 {
		t.Fatalf("len = %d, want 2", len(fields))
	}

	if fields[0] != "name" || fields[1] != "email" {
		t.Errorf("fields = %v, want [name email]", fields)
	}
}

func TestValidateOnly_Missing(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("POST", "/", nil)

	if fields := protocol.ValidateOnly(r); fields != nil {
		t.Errorf("ValidateOnly() = %v, want nil", fields)
	}
}

func TestValidateOnly_Empty(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("POST", "/", nil)

	r.Header.Set("Validate-Only", "")

	if fields := protocol.ValidateOnly(r); fields != nil {
		t.Errorf("ValidateOnly() = %v, want nil", fields)
	}
}
