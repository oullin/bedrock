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

func (s stubEncrypter) Encrypt(v string) (string, error) { return "enc:" + v, nil }
func (s stubEncrypter) Decrypt(v string) (string, error) {
	if !strings.HasPrefix(v, "enc:") {
		return "", errors.New("not encrypted")
	}

	return v[4:], nil
}

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

	// Jar should be flushed after middleware runs.
	if j.HasQueued("auth") {
		t.Fatal("expected jar to be flushed")
	}
}
