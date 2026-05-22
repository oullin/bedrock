package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/routegen"
)

func TestRenderPortalShell_NilRoutesFallsBackToDefault(t *testing.T) {
	rec := httptest.NewRecorder()

	renderPortalShell(rec, nil, nil, "Billing", map[string]any{"ok": true})

	if rec.Code != 200 {
		t.Fatalf("code = %d", rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Fatalf("content-type = %q", got)
	}
}

func TestRenderPortalShell_HappyPath(t *testing.T) {
	routes := routegen.New()

	routes.Add("billing.state", "GET", "/billing/state")

	rec := httptest.NewRecorder()

	renderPortalShell(rec, nil, routes, "Portal", map[string]any{"state": "active"})

	if rec.Code != 200 {
		t.Fatalf("code = %d", rec.Code)
	}

	body := rec.Body.String()

	if len(body) == 0 {
		t.Fatal("empty body")
	}
}

func TestRenderPortalShell_UnmarshalableStateReturnsError(t *testing.T) {
	rec := httptest.NewRecorder()

	// Channels are not JSON-serializable.
	state := map[string]any{"bad": make(chan int)}

	renderPortalShell(rec, nil, nil, "Portal", state)

	if rec.Code != 500 {
		t.Fatalf("code = %d, want 500", rec.Code)
	}
}
