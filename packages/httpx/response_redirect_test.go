package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx"
)

func TestRedirectResponseSend(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/dashboard", http.StatusFound)

	if err := resp.Send(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	if rec.Header().Get("Location") != "/dashboard" {
		t.Fatalf("expected Location: /dashboard, got %s", rec.Header().Get("Location"))
	}
}

func TestRedirectResponseStatus(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/home", http.StatusMovedPermanently)

	resp.Send()

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected 301, got %d", rec.Code)
	}
}

func TestRedirectResponseWithFragment(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/page", http.StatusFound).
		WithFragment("section1")

	resp.Send()

	if rec.Header().Get("Location") != "/page#section1" {
		t.Fatalf("expected /page#section1, got %s", rec.Header().Get("Location"))
	}
}

func TestRedirectResponseWithoutFragment(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/page", http.StatusFound).
		WithFragment("section").
		WithoutFragment()

	if resp.GetTargetURL() != "/page" {
		t.Fatalf("expected /page without fragment, got %s", resp.GetTargetURL())
	}
}

func TestRedirectResponseWithFlashData(t *testing.T) {
	t.Parallel()

	session := newStubSession()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/home", http.StatusFound).
		SetSession(session).
		With("status", "success").
		With("message", "Welcome back!")

	resp.Send()

	if session.flashed["status"] != "success" {
		t.Fatal("expected status flash data")
	}

	if session.flashed["message"] != "Welcome back!" {
		t.Fatal("expected message flash data")
	}
}

func TestRedirectResponseWithData(t *testing.T) {
	t.Parallel()

	session := newStubSession()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/home", http.StatusFound).
		SetSession(session).
		WithData(map[string]any{"a": 1, "b": 2})

	resp.Send()

	if session.flashed["a"] != 1 || session.flashed["b"] != 2 {
		t.Fatal("expected flash data from WithData")
	}
}

func TestRedirectResponseWithErrors(t *testing.T) {
	t.Parallel()

	session := newStubSession()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/form", http.StatusFound).
		SetSession(session).
		WithErrors(map[string][]string{"name": {"required"}})

	resp.Send()

	errs, ok := session.flashed["errors"]

	if !ok {
		t.Fatal("expected errors flash data")
	}

	errMap := errs.(map[string][]string)

	if len(errMap["name"]) != 1 || errMap["name"][0] != "required" {
		t.Fatalf("unexpected errors: %v", errMap)
	}
}

func TestRedirectResponseWithInput(t *testing.T) {
	t.Parallel()

	session := newStubSession()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/?name=Taylor&email=test@test.com", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/form", http.StatusFound).
		SetSession(session).
		WithInput()

	resp.Send()

	if session.oldInput["name"] != "Taylor" {
		t.Fatal("expected name to be flashed as input")
	}

	if session.oldInput["email"] != "test@test.com" {
		t.Fatal("expected email to be flashed as input")
	}
}

func TestRedirectResponseOnlyInput(t *testing.T) {
	t.Parallel()

	session := newStubSession()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/?name=Taylor&password=secret", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/form", http.StatusFound).
		SetSession(session).
		OnlyInput("name")

	resp.Send()

	if session.oldInput["name"] != "Taylor" {
		t.Fatal("expected name to be flashed")
	}

	if _, ok := session.oldInput["password"]; ok {
		t.Fatal("expected password to NOT be flashed")
	}
}

func TestRedirectResponseExceptInput(t *testing.T) {
	t.Parallel()

	session := newStubSession()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/?name=Taylor&password=secret", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/form", http.StatusFound).
		SetSession(session).
		ExceptInput("password")

	resp.Send()

	if session.oldInput["name"] != "Taylor" {
		t.Fatal("expected name to be flashed")
	}

	if _, ok := session.oldInput["password"]; ok {
		t.Fatal("expected password to NOT be flashed")
	}
}

func TestRedirectResponseHeaders(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/home", http.StatusFound).
		Header("X-Custom", "value").
		WithHeaders(map[string]string{"X-Another": "yes"})

	resp.Send()

	if rec.Header().Get("X-Custom") != "value" {
		t.Fatal("expected X-Custom header")
	}
}

func TestRedirectResponseCookie(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/home", http.StatusFound).
		Cookie(&http.Cookie{Name: "flash", Value: "yes", Path: "/"})

	resp.Send()

	cookies := rec.Result().Cookies()

	if len(cookies) != 1 || cookies[0].Name != "flash" {
		t.Fatal("expected cookie to be set")
	}
}

func TestRedirectResponseEnforceSameOrigin(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Host = "example.com"

	// Same origin - should succeed.
	resp := httpx.NewRedirectResponse(rec, raw, "http://example.com/path", http.StatusFound)

	if err := resp.EnforceSameOrigin(); err != nil {
		t.Fatalf("expected same-origin check to pass: %v", err)
	}

	// Different origin - should fail.
	resp2 := httpx.NewRedirectResponse(rec, raw, "http://evil.com/path", http.StatusFound)

	if err := resp2.EnforceSameOrigin(); err == nil {
		t.Fatal("expected same-origin check to fail")
	}
}

func TestRedirectResponseEnforceSameOriginRelative(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)

	resp := httpx.NewRedirectResponse(rec, raw, "/relative/path", http.StatusFound)

	if err := resp.EnforceSameOrigin(); err != nil {
		t.Fatalf("relative URLs should always pass same-origin: %v", err)
	}
}

func TestRedirectResponseAway(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/initial", http.StatusFound).
		Away("https://external.com", http.StatusTemporaryRedirect)

	resp.Send()

	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307, got %d", rec.Code)
	}

	if rec.Header().Get("Location") != "https://external.com" {
		t.Fatalf("expected https://external.com, got %s", rec.Header().Get("Location"))
	}
}

func TestRedirectResponseGetTargetURL(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/target", http.StatusFound)

	if resp.GetTargetURL() != "/target" {
		t.Fatalf("expected /target, got %s", resp.GetTargetURL())
	}
}

func TestRedirectResponseGetStatusCode(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/", http.StatusFound)

	if resp.GetStatusCode() != http.StatusFound {
		t.Fatalf("expected 302, got %d", resp.GetStatusCode())
	}
}

func TestSecureRedirect(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Host = "example.com"

	resp := httpx.Secure(rec, raw, "/login")
	resp.Send()

	if rec.Header().Get("Location") != "https://example.com/login" {
		t.Fatalf("expected https://example.com/login, got %s", rec.Header().Get("Location"))
	}
}

func TestRedirectWithoutSession(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httpx.NewRedirectResponse(rec, raw, "/home", http.StatusFound).
		With("key", "value")

	// Should not panic when no session is attached.
	resp.Send()

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
}
