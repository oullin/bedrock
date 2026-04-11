package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx"
)

// stubSession implements httpx.SessionStore for testing.
type stubSession struct {
	data     map[string]any
	flashed  map[string]any
	oldInput map[string]any
	removed  []string
}

func newStubSession() *stubSession {
	return &stubSession{
		data:     make(map[string]any),
		flashed:  make(map[string]any),
		oldInput: make(map[string]any),
	}
}

func (s *stubSession) Get(key string, fallback any) any {
	if v, ok := s.data[key]; ok {
		return v
	}

	return fallback
}

func (s *stubSession) Put(key string, value any) {
	s.data[key] = value
}

func (s *stubSession) Flash(key string, value any) {
	s.flashed[key] = value
}

func (s *stubSession) GetOldInput(key string, fallback any) any {
	if v, ok := s.oldInput[key]; ok {
		return v
	}

	return fallback
}

func (s *stubSession) HasOldInput(key string) bool {
	_, ok := s.oldInput[key]

	return ok
}

func (s *stubSession) FlashInput(values map[string]any) {
	for k, v := range values {
		s.oldInput[k] = v
	}
}

func (s *stubSession) Remove(key string) any {
	s.removed = append(s.removed, key)
	v := s.data[key]
	delete(s.data, key)

	return v
}

func TestOldWithSession(t *testing.T) {
	t.Parallel()

	session := newStubSession()
	session.oldInput["name"] = "Taylor"

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httpx.NewRequest(raw)
	req.SetSession(session)

	if req.Old("name") != "Taylor" {
		t.Fatalf("expected Taylor, got %v", req.Old("name"))
	}

	if req.Old("missing", "default") != "default" {
		t.Fatal("expected fallback for missing old input")
	}
}

func TestOldWithoutSession(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httpx.NewRequest(raw)

	if req.Old("name") != nil {
		t.Fatal("expected nil when no session attached")
	}

	if req.Old("name", "default") != "default" {
		t.Fatal("expected fallback when no session attached")
	}
}

func TestHasOld(t *testing.T) {
	t.Parallel()

	session := newStubSession()
	session.oldInput["email"] = "test@test.com"

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httpx.NewRequest(raw)
	req.SetSession(session)

	if !req.HasOld("email") {
		t.Fatal("expected HasOld to return true")
	}

	if req.HasOld("missing") {
		t.Fatal("expected HasOld to return false for missing key")
	}
}

func TestFlash(t *testing.T) {
	t.Parallel()

	session := newStubSession()

	raw := httptest.NewRequest(http.MethodGet, "/?name=Taylor&email=test@test.com", nil)
	req := httpx.NewRequest(raw)
	req.SetSession(session)

	req.Flash()

	if session.oldInput["name"] != "Taylor" {
		t.Fatal("expected name to be flashed")
	}

	if session.oldInput["email"] != "test@test.com" {
		t.Fatal("expected email to be flashed")
	}
}

func TestFlashOnly(t *testing.T) {
	t.Parallel()

	session := newStubSession()

	raw := httptest.NewRequest(http.MethodGet, "/?name=Taylor&email=test@test.com", nil)
	req := httpx.NewRequest(raw)
	req.SetSession(session)

	req.FlashOnly("name")

	if session.oldInput["name"] != "Taylor" {
		t.Fatal("expected name to be flashed")
	}

	if _, ok := session.oldInput["email"]; ok {
		t.Fatal("expected email to NOT be flashed")
	}
}

func TestFlashExcept(t *testing.T) {
	t.Parallel()

	session := newStubSession()

	raw := httptest.NewRequest(http.MethodGet, "/?name=Taylor&password=secret", nil)
	req := httpx.NewRequest(raw)
	req.SetSession(session)

	req.FlashExcept("password")

	if session.oldInput["name"] != "Taylor" {
		t.Fatal("expected name to be flashed")
	}

	if _, ok := session.oldInput["password"]; ok {
		t.Fatal("expected password to NOT be flashed")
	}
}

func TestFlush(t *testing.T) {
	t.Parallel()

	session := newStubSession()
	session.data["_old_input"] = map[string]any{"name": "Taylor"}

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httpx.NewRequest(raw)
	req.SetSession(session)

	req.Flush()

	if len(session.removed) == 0 || session.removed[0] != "_old_input" {
		t.Fatal("expected _old_input to be removed")
	}
}

func TestFlashWithoutSession(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?name=Taylor", nil)
	req := httpx.NewRequest(raw)

	// Should not panic.
	req.Flash()
	req.FlashOnly("name")
	req.FlashExcept("name")
	req.Flush()
}
