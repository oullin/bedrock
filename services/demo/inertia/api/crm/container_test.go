package crm

import (
	"database/sql"
	"net/http"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/services/demo/inertia/api/internal/database"
)

func TestContainerValidate_ZeroValue(t *testing.T) {
	err := Container{}.Validate()

	if err == nil {
		t.Fatal("expected error for zero-value container")
	}

	for _, field := range []string{"DB", "RequireAuth", "Render", "Redirect", "RouteURL", "SetFlash", "CurrentUser"} {
		if !strings.Contains(err.Error(), field) {
			t.Errorf("error should mention %s, got: %s", field, err.Error())
		}
	}
}

func TestContainerValidate_Valid(t *testing.T) {
	requireAuthFn := func(h http.Handler) http.Handler { return h }
	renderFn := func(http.ResponseWriter, *http.Request, string, protocol.Props) {}
	redirectFn := func(http.ResponseWriter, *http.Request, string) {}
	routeURLFn := func(string, map[string]string) string { return "" }
	setFlashFn := func(http.ResponseWriter, flash.Message) error { return nil }
	currentUserFn := func(*http.Request) *database.User { return nil }

	c := Container{
		DB:          &sql.DB{},
		RequireAuth: requireAuthFn,
		Render:      renderFn,
		Redirect:    redirectFn,
		RouteURL:    routeURLFn,
		SetFlash:    setFlashFn,
		CurrentUser: currentUserFn,
	}

	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
