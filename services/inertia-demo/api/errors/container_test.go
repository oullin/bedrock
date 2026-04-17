package errors

import (
	"net/http"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
)

func TestContainerValidate_ZeroValue(t *testing.T) {
	err := Container{}.Validate()

	if err == nil {
		t.Fatal("expected error for zero-value container")
	}

	for _, field := range []string{"RequireAuth", "Render"} {
		if !strings.Contains(err.Error(), field) {
			t.Errorf("error should mention %s, got: %s", field, err.Error())
		}
	}
}

func TestContainerValidate_Valid(t *testing.T) {
	requireAuthFn := func(h http.Handler) http.Handler { return h }
	renderFn := func(http.ResponseWriter, *http.Request, string, protocol.Props) {}

	c := Container{
		RequireAuth: requireAuthFn,
		Render:      renderFn,
	}

	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
