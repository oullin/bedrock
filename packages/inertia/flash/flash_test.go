package flash

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var testKey = []byte("0123456789abcdef0123456789abcdef")

func TestNewCookieStoreDefaults(t *testing.T) {
	t.Parallel()

	s := NewCookieStore()

	if s.cookieName != "inertia_flash" {
		t.Errorf("expected cookie name %q, got %q", "inertia_flash", s.cookieName)
	}

	if s.path != "/" {
		t.Errorf("expected path %q, got %q", "/", s.path)
	}

	if !s.httpOnly {
		t.Error("expected httpOnly to be true")
	}

	if s.secure {
		t.Error("expected secure to be false by default")
	}
}

func TestNewCookieStoreWithOptions(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(
		WithCookieName("my_flash"),
		WithPath("/app"),
		WithSecure(true),
		WithHTTPOnly(false),
		WithSameSite(http.SameSiteStrictMode),
	)

	if s.cookieName != "my_flash" {
		t.Errorf("expected cookie name %q, got %q", "my_flash", s.cookieName)
	}

	if s.path != "/app" {
		t.Errorf("expected path %q, got %q", "/app", s.path)
	}

	if !s.secure {
		t.Error("expected secure to be true")
	}

	if s.httpOnly {
		t.Error("expected httpOnly to be false")
	}

	if s.sameSite != http.SameSiteStrictMode {
		t.Errorf("expected SameSiteStrictMode, got %v", s.sameSite)
	}
}

func TestSetAndConsume(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(WithCookieName("test_flash"))
	msg := Message{Kind: "success", Title: "Done", Message: "Contact created."}

	rec := httptest.NewRecorder()

	if err := s.Set(rec, msg); err != nil {
		t.Fatalf("Set returned unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(findCookie(t, rec, "test_flash"))

	rec2 := httptest.NewRecorder()
	got := s.Consume(rec2, req)

	if got == nil {
		t.Fatal("expected a flash message, got nil")
	}

	if got.Kind != msg.Kind {
		t.Errorf("expected Kind %q, got %q", msg.Kind, got.Kind)
	}

	if got.Title != msg.Title {
		t.Errorf("expected Title %q, got %q", msg.Title, got.Title)
	}

	if got.Message != msg.Message {
		t.Errorf("expected Message %q, got %q", msg.Message, got.Message)
	}

	if dc := findCookie(t, rec2, "test_flash"); dc.MaxAge != -1 {
		t.Errorf("expected MaxAge -1 (delete), got %d", dc.MaxAge)
	}
}

func TestConsumeMirrorsSecurityAttributes(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(
		WithCookieName("sec_flash"),
		WithSecure(true),
		WithHTTPOnly(true),
		WithSameSite(http.SameSiteStrictMode),
	)

	rec := httptest.NewRecorder()

	if err := s.Set(rec, Message{Kind: "info", Title: "Hi", Message: "Hello."}); err != nil {
		t.Fatalf("Set returned unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(findCookie(t, rec, "sec_flash"))

	rec2 := httptest.NewRecorder()

	s.Consume(rec2, req)

	dc := findCookie(t, rec2, "sec_flash")

	if !dc.HttpOnly {
		t.Error("expected deletion cookie to have HttpOnly set")
	}

	if !dc.Secure {
		t.Error("expected deletion cookie to have Secure set")
	}

	if dc.SameSite != http.SameSiteStrictMode {
		t.Errorf("expected deletion cookie SameSite %v, got %v", http.SameSiteStrictMode, dc.SameSite)
	}
}

func TestConsumeWithNoCookie(t *testing.T) {
	t.Parallel()

	s := NewCookieStore()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	got := s.Consume(rec, req)

	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestConsumeWithInvalidJSON(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(WithCookieName("test_flash"))
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{Name: "test_flash", Value: "not-json"})

	rec := httptest.NewRecorder()
	got := s.Consume(rec, req)

	if got != nil {
		t.Errorf("expected nil for invalid JSON, got %+v", got)
	}
}

func TestConsumeAfterConsumeReturnsNil(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(WithCookieName("test_flash"))
	msg := Message{Kind: "info", Title: "Hey", Message: "Hello."}

	rec := httptest.NewRecorder()

	if err := s.Set(rec, msg); err != nil {
		t.Fatalf("Set returned unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(findCookie(t, rec, "test_flash"))

	rec2 := httptest.NewRecorder()
	first := s.Consume(rec2, req)

	if first == nil {
		t.Fatal("expected first consume to return message")
	}

	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	rec3 := httptest.NewRecorder()
	second := s.Consume(rec3, req2)

	if second != nil {
		t.Errorf("expected second consume to return nil, got %+v", second)
	}
}

// 32 bytes

func TestSetAndConsumeEncrypted(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(WithCookieName("enc_flash"), WithKey(testKey))
	msg := Message{Kind: "success", Title: "Done", Message: "Encrypted."}

	rec := httptest.NewRecorder()

	if err := s.Set(rec, msg); err != nil {
		t.Fatalf("Set returned unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(findCookie(t, rec, "enc_flash"))

	rec2 := httptest.NewRecorder()
	got := s.Consume(rec2, req)

	if got == nil {
		t.Fatal("expected a flash message, got nil")
	}

	if got.Kind != msg.Kind || got.Title != msg.Title || got.Message != msg.Message {
		t.Errorf("round-trip mismatch: got %+v, want %+v", got, msg)
	}
}

func TestConsumeRejectsTamperedCookie(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(WithCookieName("enc_flash"), WithKey(testKey))
	msg := Message{Kind: "info", Title: "Hi", Message: "Hello."}

	rec := httptest.NewRecorder()

	if err := s.Set(rec, msg); err != nil {
		t.Fatalf("Set returned unexpected error: %v", err)
	}

	cookie := findCookie(t, rec, "enc_flash")
	cookie.Value = cookie.Value[:len(cookie.Value)-4] + "XXXX"

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(cookie)

	rec2 := httptest.NewRecorder()
	got := s.Consume(rec2, req)

	if got != nil {
		t.Errorf("expected nil for tampered cookie, got %+v", got)
	}
}

func TestConsumeRejectsForgedPlaintext(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(WithCookieName("enc_flash"), WithKey(testKey))

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{
		Name:  "enc_flash",
		Value: `%7B%22kind%22%3A%22error%22%2C%22title%22%3A%22Forged%22%2C%22message%22%3A%22Evil%22%7D`,
	})

	rec := httptest.NewRecorder()
	got := s.Consume(rec, req)

	if got != nil {
		t.Errorf("expected nil for forged plaintext cookie, got %+v", got)
	}
}

func findCookie(t *testing.T, w *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()

	for _, c := range w.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}

	t.Fatalf("cookie %q not found", name)

	return nil
}
