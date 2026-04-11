package cookie_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/cookie"
)

// stubEncrypter is a trivial reversible encrypter for testing.
type stubEncrypter struct{}

// failEncrypter always fails on encrypt.
type failEncrypter struct{}

func (s stubEncrypter) Encrypt(v string) (string, error) { return "enc:" + v, nil }
func (s stubEncrypter) Decrypt(v string) (string, error) {
	if !strings.HasPrefix(v, "enc:") {
		return "", errors.New("not encrypted")
	}

	return v[4:], nil
}

func (f failEncrypter) Encrypt(string) (string, error)   { return "", errors.New("encrypt fail") }
func (f failEncrypter) Decrypt(v string) (string, error) { return v, nil }

// ---------------------------------------------------------------------------
// EncryptCookies — request decryption
// ---------------------------------------------------------------------------

func TestEncryptCookiesDecryptsRequest(t *testing.T) {
	t.Parallel()

	mw := cookie.NewEncryptCookies(stubEncrypter{})

	var got string

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("session")

		if c != nil {
			got = c.Value
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "enc:abc123"})

	h.ServeHTTP(httptest.NewRecorder(), req)

	if got != "abc123" {
		t.Fatalf("expected decrypted value 'abc123', got %q", got)
	}
}

func TestEncryptCookiesDecryptsMultipleRequestCookies(t *testing.T) {
	t.Parallel()

	mw := cookie.NewEncryptCookies(stubEncrypter{})

	var gotA, gotB string

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, _ := r.Cookie("a"); c != nil {
			gotA = c.Value
		}

		if c, _ := r.Cookie("b"); c != nil {
			gotB = c.Value
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "a", Value: "enc:alpha"})
	req.AddCookie(&http.Cookie{Name: "b", Value: "enc:beta"})

	h.ServeHTTP(httptest.NewRecorder(), req)

	if gotA != "alpha" {
		t.Fatalf("expected 'alpha', got %q", gotA)
	}

	if gotB != "beta" {
		t.Fatalf("expected 'beta', got %q", gotB)
	}
}

func TestEncryptCookiesKeepsOriginalOnDecryptFailure(t *testing.T) {
	t.Parallel()

	mw := cookie.NewEncryptCookies(stubEncrypter{})

	var got string

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("session")

		if c != nil {
			got = c.Value
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// Value without "enc:" prefix will fail decryption.
	req.AddCookie(&http.Cookie{Name: "session", Value: "not-encrypted"})

	h.ServeHTTP(httptest.NewRecorder(), req)

	if got != "not-encrypted" {
		t.Fatalf("expected original value 'not-encrypted', got %q", got)
	}
}

// ---------------------------------------------------------------------------
// EncryptCookies — except list (request)
// ---------------------------------------------------------------------------

func TestEncryptCookiesExceptPassThrough(t *testing.T) {
	t.Parallel()

	mw := cookie.NewEncryptCookies(stubEncrypter{}, "csrf")

	var got string

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("csrf")

		if c != nil {
			got = c.Value
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "csrf", Value: "raw-token"})

	h.ServeHTTP(httptest.NewRecorder(), req)

	if got != "raw-token" {
		t.Fatalf("expected raw token, got %q", got)
	}
}

func TestEncryptCookiesExceptDoesNotDecrypt(t *testing.T) {
	t.Parallel()

	// Even if value looks encrypted, excepted cookies pass through.
	mw := cookie.NewEncryptCookies(stubEncrypter{}, "special")

	var got string

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("special")

		if c != nil {
			got = c.Value
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "special", Value: "enc:should-not-decrypt"})

	h.ServeHTTP(httptest.NewRecorder(), req)

	if got != "enc:should-not-decrypt" {
		t.Fatalf("expected raw value, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// EncryptCookies — response encryption
// ---------------------------------------------------------------------------

func TestEncryptCookiesEncryptsResponse(t *testing.T) {
	t.Parallel()

	mw := cookie.NewEncryptCookies(stubEncrypter{})

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "abc123"})
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cookies := rec.Result().Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected at least one cookie in response")
	}

	var found *http.Cookie

	for _, c := range cookies {
		if c.Name == "session" {
			found = c

			break
		}
	}

	if found == nil {
		t.Fatal("expected session cookie in response")
	}

	if found.Value != "enc:abc123" {
		t.Fatalf("expected encrypted value 'enc:abc123', got %q", found.Value)
	}
}

func TestEncryptCookiesEncryptsMultipleResponseCookies(t *testing.T) {
	t.Parallel()

	mw := cookie.NewEncryptCookies(stubEncrypter{})

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "a", Value: "alpha"})
		http.SetCookie(w, &http.Cookie{Name: "b", Value: "beta"})
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cookies := rec.Result().Cookies()

	values := make(map[string]string)

	for _, c := range cookies {
		values[c.Name] = c.Value
	}

	if values["a"] != "enc:alpha" {
		t.Fatalf("expected 'enc:alpha', got %q", values["a"])
	}

	if values["b"] != "enc:beta" {
		t.Fatalf("expected 'enc:beta', got %q", values["b"])
	}
}

func TestEncryptCookiesExceptResponsePassThrough(t *testing.T) {
	t.Parallel()

	mw := cookie.NewEncryptCookies(stubEncrypter{}, "unencrypted")

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "encrypted", Value: "secret"})
		http.SetCookie(w, &http.Cookie{Name: "unencrypted", Value: "plain"})
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cookies := rec.Result().Cookies()
	values := make(map[string]string)

	for _, c := range cookies {
		values[c.Name] = c.Value
	}

	if values["encrypted"] != "enc:secret" {
		t.Fatalf("expected 'enc:secret', got %q", values["encrypted"])
	}

	if values["unencrypted"] != "plain" {
		t.Fatalf("expected 'plain', got %q", values["unencrypted"])
	}
}

func TestEncryptCookiesKeepsOriginalOnEncryptFailure(t *testing.T) {
	t.Parallel()

	mw := cookie.NewEncryptCookies(failEncrypter{})

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "plain"})
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cookies := rec.Result().Cookies()

	var found *http.Cookie

	for _, c := range cookies {
		if c.Name == "session" {
			found = c

			break
		}
	}

	if found == nil {
		t.Fatal("expected session cookie")
	}

	// On encrypt failure, the original value is kept.
	if found.Value != "plain" {
		t.Fatalf("expected original value 'plain', got %q", found.Value)
	}
}

// ---------------------------------------------------------------------------
// AttachQueued — writes cookies to response
// ---------------------------------------------------------------------------

func TestAttachQueuedWritesCookies(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "auth", Value: "token"})

	mw := cookie.NewAttachQueued(j)

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cookies := rec.Result().Cookies()
	found := false

	for _, c := range cookies {
		if c.Name == "auth" && c.Value == "token" {
			found = true
		}
	}

	if !found {
		t.Fatalf("expected auth cookie in response, got %v", cookies)
	}
}

func TestAttachQueuedMultipleCookies(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "a", Value: "1"})
	j.Queue(&http.Cookie{Name: "b", Value: "2"})
	j.Queue(&http.Cookie{Name: "c", Value: "3"})

	mw := cookie.NewAttachQueued(j)

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cookies := rec.Result().Cookies()

	if len(cookies) != 3 {
		t.Fatalf("expected 3 cookies, got %d", len(cookies))
	}

	names := make(map[string]bool)

	for _, c := range cookies {
		names[c.Name] = true
	}

	for _, name := range []string{"a", "b", "c"} {
		if !names[name] {
			t.Fatalf("expected cookie %q in response", name)
		}
	}
}

func TestAttachQueuedFlushesJar(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "auth", Value: "token"})

	mw := cookie.NewAttachQueued(j)

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	// Jar should be flushed after middleware runs.
	if j.HasQueued("auth") {
		t.Fatal("expected jar to be flushed")
	}

	if len(j.GetQueued()) != 0 {
		t.Fatalf("expected empty queue, got %d", len(j.GetQueued()))
	}
}

func TestAttachQueuedNoCookies(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	mw := cookie.NewAttachQueued(j)

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cookies := rec.Result().Cookies()

	if len(cookies) != 0 {
		t.Fatalf("expected no cookies, got %d", len(cookies))
	}
}

func TestAttachQueuedWithPathCookies(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar", Path: "/path"})
	j.Queue(&http.Cookie{Name: "foo", Value: "rab", Path: "/"})

	mw := cookie.NewAttachQueued(j)

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cookies := rec.Result().Cookies()

	if len(cookies) != 2 {
		t.Fatalf("expected 2 cookies, got %d", len(cookies))
	}
}
