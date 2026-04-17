package middleware_test

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/config"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/inertia/middleware"
)

func testKey(t *testing.T) []byte {
	t.Helper()

	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	return key
}

func csrfMiddleware(t *testing.T) (func(http.Handler) http.Handler, []byte) {
	t.Helper()

	key := testKey(t)

	return middleware.CSRF(config.CSRFConfig{}, key), key
}

func issueCSRFCookie(t *testing.T, mw func(http.Handler) http.Handler) (*http.Cookie, string) {
	t.Helper()

	var rawToken string

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken = protocol.CSRFTokenFromContext(r.Context())

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	return findCookie(t, w, "XSRF-TOKEN"), rawToken
}

func TestCSRF_SetsCookieOnGET(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	cookies := resp.Cookies()
	found := false

	for _, c := range cookies {
		if c.Name == "XSRF-TOKEN" {
			found = true

			if c.HttpOnly {
				t.Error("cookie should not be HttpOnly")
			}
		}
	}

	if !found {
		t.Error("missing CSRF cookie")
	}
}

func TestCSRF_SafeMethodsPassThrough(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
		r := httptest.NewRequest(method, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("%s: status = %d, want %d", method, w.Code, http.StatusOK)
		}
	}
}

func TestCSRF_POSTWithoutToken_419(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != 419 {
		t.Errorf("status = %d, want 419", w.Code)
	}
}

func TestCSRF_POSTWithValidToken_XCSRF(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)
	tokenCookie, rawToken := issueCSRFCookie(t, mw)

	// Step 2: POST with cookie + X-CSRF-TOKEN header (raw token).
	postHandler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	postReq := httptest.NewRequest(http.MethodPost, "/submit", nil)

	postReq.AddCookie(tokenCookie)

	postReq.Header.Set("X-CSRF-TOKEN", rawToken)

	postW := httptest.NewRecorder()

	postHandler.ServeHTTP(postW, postReq)

	if postW.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", postW.Code, http.StatusOK)
	}
}

func TestCSRF_POSTWithValidToken_XXSRF(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)
	tokenCookie, _ := issueCSRFCookie(t, mw)

	// POST with cookie + X-XSRF-TOKEN header (encrypted value from cookie).
	postReq := httptest.NewRequest(http.MethodPost, "/submit", nil)

	postReq.AddCookie(tokenCookie)

	postReq.Header.Set("X-XSRF-TOKEN", tokenCookie.Value)

	postW := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(postW, postReq)

	if postW.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", postW.Code, http.StatusOK)
	}
}

func TestCSRF_POSTWithValidToken_XXSRF_URLEncoded(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)
	tokenCookie, _ := issueCSRFCookie(t, mw)

	postReq := httptest.NewRequest(http.MethodPost, "/submit", nil)

	postReq.AddCookie(tokenCookie)

	postReq.Header.Set("X-XSRF-TOKEN", url.QueryEscape(tokenCookie.Value))

	postW := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(postW, postReq)

	if postW.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", postW.Code, http.StatusOK)
	}
}

func TestCSRF_POSTWithFormField(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)
	tokenCookie, rawToken := issueCSRFCookie(t, mw)

	// POST with _token form field.
	form := url.Values{}

	form.Set("_token", rawToken)

	postReq := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(form.Encode()))

	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	postReq.AddCookie(tokenCookie)

	postW := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(postW, postReq)

	if postW.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", postW.Code, http.StatusOK)
	}
}

func TestCSRF_POSTWithWrongToken_419(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)

	getReq := httptest.NewRequest(http.MethodGet, "/", nil)
	getW := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(getW, getReq)

	tokenCookie := findCookie(t, getW, "XSRF-TOKEN")

	postReq := httptest.NewRequest(http.MethodPost, "/submit", nil)

	postReq.AddCookie(tokenCookie)

	postReq.Header.Set("X-CSRF-TOKEN", "wrong-token")

	postW := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(postW, postReq)

	if postW.Code != 419 {
		t.Errorf("status = %d, want 419", postW.Code)
	}
}

func TestCSRF_MutationMethods(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		r := httptest.NewRequest(method, "/", nil)
		w := httptest.NewRecorder()

		mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(w, r)

		if w.Code != 419 {
			t.Errorf("%s: status = %d, want 419", method, w.Code)
		}
	}
}

func TestCSRFFromFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	csrfPath := filepath.Join(dir, "csrf.yml")
	cryptoPath := filepath.Join(dir, "crypto.yml")

	// Generate a valid 32-byte key.
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	csrfContent := `
cookie_name: "_my_csrf"
secure: false
same_site: "strict"
`

	cryptoContent := `
key: "` + encodeKey(key) + `"
`

	if err := os.WriteFile(csrfPath, []byte(csrfContent), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(cryptoPath, []byte(cryptoContent), 0644); err != nil {
		t.Fatal(err)
	}

	mw, err := middleware.CSRFFromFile(csrfPath, cryptoPath)

	if err != nil {
		t.Fatal(err)
	}

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	// Check that it used the custom cookie name.
	found := false

	for _, c := range w.Result().Cookies() {
		if c.Name == "_my_csrf" {
			found = true
		}
	}

	if !found {
		t.Error("expected custom cookie name _my_csrf")
	}
}

func TestCSRFFromFile_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := middleware.CSRFFromFile("/nonexistent/csrf.yml", "/nonexistent/crypto.yml")

	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCSRF_SameSiteStrict(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	mw := middleware.CSRF(config.CSRFConfig{
		SameSite: "strict",
	}, key)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCSRF_SameSiteNone(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	mw := middleware.CSRF(config.CSRFConfig{
		SameSite: "none",
		Secure:   true,
	}, key)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCSRF_POSTWithMalformedCookie_RegeneratesToken(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Send a GET with a cookie that has invalid encrypted value.
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "not-valid-encrypted-payload"})

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	// Should succeed on GET (new token generated and cookie set).
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	// Should have set a new cookie.
	cookies := w.Result().Cookies()
	found := false

	for _, c := range cookies {
		if c.Name == "XSRF-TOKEN" && c.Value != "not-valid-encrypted-payload" {
			found = true
		}
	}

	if !found {
		t.Error("expected new CSRF cookie to be set after malformed cookie")
	}
}

func TestCSRF_StoresTokenInContext(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)

	var ctxToken string

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxToken = protocol.CSRFTokenFromContext(r.Context())

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if strings.TrimSpace(ctxToken) == "" {
		t.Error("expected CSRF token to be stored in context")
	}
}

func TestCSRF_POSTWithXXSRFToken_InvalidPayload_419(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)

	// GET to get a valid cookie.
	getReq := httptest.NewRequest(http.MethodGet, "/", nil)
	getW := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(getW, getReq)

	tokenCookie := findCookie(t, getW, "XSRF-TOKEN")

	// POST with tampered X-XSRF-TOKEN header.
	postReq := httptest.NewRequest(http.MethodPost, "/submit", nil)

	postReq.AddCookie(tokenCookie)

	postReq.Header.Set("X-XSRF-TOKEN", "tampered-encrypted-value")

	postW := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(postW, postReq)

	if postW.Code != 419 {
		t.Errorf("status = %d, want 419", postW.Code)
	}
}

func TestCSRF_TokenSourcePrecedence(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)
	tokenCookie, rawToken := issueCSRFCookie(t, mw)

	// POST with valid _token form field and invalid X-CSRF-TOKEN header.
	// The _token field should take precedence.
	form := url.Values{}

	form.Set("_token", rawToken)

	postReq := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(form.Encode()))

	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postReq.Header.Set("X-CSRF-TOKEN", "wrong-token")

	postReq.AddCookie(tokenCookie)

	postW := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(postW, postReq)

	if postW.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (form field should take precedence)", postW.Code, http.StatusOK)
	}
}

func TestCSRF_POSTWithSameOriginFetchSite_SkipsTokenValidation(t *testing.T) {
	t.Parallel()

	mw, _ := csrfMiddleware(t)

	req := httptest.NewRequest(http.MethodPost, "/submit", nil)

	req.Header.Set("Sec-Fetch-Site", "same-origin")

	w := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCSRF_OriginOnlyFailedVerification_Returns403(t *testing.T) {
	t.Parallel()

	key := testKey(t)
	mw := middleware.CSRF(config.CSRFConfig{
		OriginOnly: true,
	}, key)

	req := httptest.NewRequest(http.MethodPost, "/submit", nil)
	w := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestCSRFFromFile_AllEnvOverrides(t *testing.T) {
	t.Setenv("INERTIA_CSRF_COOKIE_NAME", "_env_csrf")
	t.Setenv("INERTIA_CSRF_SECURE", "true")
	t.Setenv("INERTIA_CSRF_SAME_SITE", "strict")

	dir := t.TempDir()
	csrfPath := filepath.Join(dir, "csrf.yml")
	cryptoPath := filepath.Join(dir, "crypto.yml")

	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	csrfContent := `
cookie_name: "_file_csrf"
secure: false
same_site: "lax"
`

	cryptoContent := `
key: "` + encodeKey(key) + `"
`

	if err := os.WriteFile(csrfPath, []byte(csrfContent), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(cryptoPath, []byte(cryptoContent), 0644); err != nil {
		t.Fatal(err)
	}

	mw, err := middleware.CSRFFromFile(csrfPath, cryptoPath)

	if err != nil {
		t.Fatal(err)
	}

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	// Verify env-overridden cookie name.
	found := false

	for _, c := range w.Result().Cookies() {
		if c.Name == "_env_csrf" {
			found = true
		}
	}

	if !found {
		t.Error("expected env-overridden cookie name _env_csrf")
	}
}

func TestCSRF_SameSiteFetchSite_AllowSameSiteTrue(t *testing.T) {
	t.Parallel()

	key := testKey(t)
	mw := middleware.CSRF(config.CSRFConfig{
		AllowSameSite: true,
	}, key)

	req := httptest.NewRequest(http.MethodPost, "/submit", nil)

	req.Header.Set("Sec-Fetch-Site", "same-site")

	w := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCSRFFromFile_InvalidCryptoPath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	csrfPath := filepath.Join(dir, "csrf.yml")

	csrfContent := `cookie_name: "TOKEN"`

	if err := os.WriteFile(csrfPath, []byte(csrfContent), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := middleware.CSRFFromFile(csrfPath, "/nonexistent/crypto.yml")

	if err == nil {
		t.Error("expected error for missing crypto config")
	}
}

func TestCSRF_SameSiteFetchSite_AllowSameSiteFalse(t *testing.T) {
	t.Parallel()

	key := testKey(t)
	mw := middleware.CSRF(config.CSRFConfig{
		AllowSameSite: false,
	}, key)

	req := httptest.NewRequest(http.MethodPost, "/submit", nil)

	req.Header.Set("Sec-Fetch-Site", "same-site")

	w := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(w, req)

	// Should fail because same-site not allowed and no token provided.
	if w.Code != 419 {
		t.Errorf("status = %d, want 419", w.Code)
	}
}

// --- helpers ---

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

func encodeKey(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}
