package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bedrock/packages/encryption"
	"github.com/bedrock/packages/inertia"
	"github.com/bedrock/packages/inertia/assert"
	coreflash "github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/middleware"
	"github.com/bedrock/packages/inertia/protocol"
	corei18n "github.com/bedrock/packages/seo/i18n"
	"github.com/bedrock/services/demo/inertia/api/auth"
	"github.com/bedrock/services/demo/inertia/api/internal/database"
	"github.com/bedrock/services/demo/inertia/api/internal/seed"
	"github.com/bedrock/services/demo/inertia/api/internal/testutil"
)

// testCryptoKey is a zero-filled 32-byte key used only in tests.
var testCryptoKey = make([]byte, 32)

func encryptForTest(t *testing.T, plaintext string) string {
	t.Helper()

	enc, err := encryption.NewEncrypter(testCryptoKey, encryption.AES256CBC)

	if err != nil {
		t.Fatalf("encrypter: %v", err)
	}

	s, err := enc.EncryptString(plaintext)

	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	return s
}

func decryptForTest(t *testing.T, encoded string) string {
	t.Helper()

	enc, err := encryption.NewEncrypter(testCryptoKey, encryption.AES256CBC)

	if err != nil {
		t.Fatalf("encrypter: %v", err)
	}

	s, err := enc.DecryptString(encoded)

	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	return s
}

func TestLoginHandlerRendersPage(t *testing.T) {
	_, testMux := newPortTestMux(t)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)

	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = "/login"
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	page := assert.AssertFromBytes(t, w.Body.Bytes())

	page.AssertComponent(t, "Auth/Login")

	page.AssertHasProp(t, "auth")
	page.AssertHasProp(t, "routes")
}

func TestLoginHandlerCreatesSession(t *testing.T) {
	_, testMux := newPortTestMux(t)
	csrfCookie, rawToken := issuePortCSRFCookie(t, testMux, "/login")

	body := strings.NewReader(url.Values{
		"email":    {"test@example.com"},
		"password": {"password"},
		"remember": {"true"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/login", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-TOKEN", rawToken)

	req.AddCookie(csrfCookie)

	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/dashboard" {
		t.Fatalf("location = %q, want %q", got, "/dashboard")
	}

	cookie := testutil.FindCookie(t, w, auth.SessionCookieName)

	decrypted := decryptForTest(t, cookie.Value)

	if decrypted != "1" {
		t.Fatalf("session = %q, want %q", decrypted, "1")
	}
}

func TestLoginHandlerRejectsInvalidPassword(t *testing.T) {
	_, testMux := newPortTestMux(t)
	csrfCookie, rawToken := issuePortCSRFCookie(t, testMux, "/login")

	body := strings.NewReader(url.Values{
		"email":    {"test@example.com"},
		"password": {"wrong-password"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/login", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = "/login"

	req.AddCookie(csrfCookie)

	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	page := assert.AssertFromBytes(t, w.Body.Bytes())

	page.AssertComponent(t, "Auth/Login")

	errors, ok := page.Props["errors"].(map[string]any)

	if !ok {
		t.Fatal("errors prop not found or not a map")
	}

	if errors["email"] != "Use test@example.com and password to sign in." {
		t.Fatalf("errors[email] = %v", errors["email"])
	}
}

func TestDashboardRequiresSession(t *testing.T) {
	_, testMux := newPortTestMux(t)
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/login" {
		t.Fatalf("location = %q, want %q", got, "/login")
	}
}

func TestDashboardRendersForAuthenticatedUser(t *testing.T) {
	_, testMux := newPortTestMux(t)

	encrypted := mustEncryptSession(t, "1")

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)

	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = "/dashboard"

	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: encrypted})

	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	page := assert.AssertFromBytes(t, w.Body.Bytes())

	page.AssertComponent(t, "Crm/Dashboard")

	page.AssertHasProp(t, "recentActivity")
	page.AssertHasProp(t, "auth")
}

func TestLegacyDemoRoutesReturnNotFound(t *testing.T) {
	_, testMux := newPortTestMux(t)

	encrypted := mustEncryptSession(t, "1")

	tests := []struct {
		name string
		path string
		want int
	}{
		{name: "root", path: "/", want: http.StatusFound},
		{name: "dashboard navigation", path: "/dashboard/navigation", want: http.StatusNotFound},
		{name: "dashboard data", path: "/dashboard/data", want: http.StatusNotFound},
		{name: "dashboard state", path: "/dashboard/state", want: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: encrypted})

			w := httptest.NewRecorder()

			testMux.ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Fatalf("status = %d, want %d", w.Code, tt.want)
			}
		})
	}
}

func TestRootRedirectsAndFeatureErrorRoutes(t *testing.T) {
	t.Parallel()

	_, testMux := newPortTestMux(t)
	encrypted := mustEncryptSession(t, "1")

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: encrypted})

	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/dashboard" {
		t.Fatalf("root status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	req = httptest.NewRequest(http.MethodGet, "/features/errors/http-error", nil)

	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = "/features/errors/http-error"

	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: encrypted})

	w = httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("errors page status = %d, want %d", w.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/features/errors/http-error/403", nil)

	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: encrypted})

	w = httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("http-error/403 status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestFeatureExternalRedirectUsesInertiaLocation(t *testing.T) {
	t.Parallel()

	_, testMux := newPortTestMux(t)
	encrypted := mustEncryptSession(t, "1")
	csrfCookie, rawToken := issuePortCSRFCookie(t, testMux, "/login")

	req := httptest.NewRequest(http.MethodPost, "/features/navigation/redirects/external", nil)

	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = "/features/navigation/redirects/external"

	req.Header.Set("X-CSRF-TOKEN", rawToken)

	req.AddCookie(csrfCookie)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: encrypted})

	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusConflict || w.Header().Get(protocol.HeaderLocation) != "https://oullin.io" {
		t.Fatalf("status = %d, location = %q", w.Code, w.Header().Get(protocol.HeaderLocation))
	}
}

func TestStoreContactCreatesRecord(t *testing.T) {
	rt, testMux := newPortTestMux(t)
	csrfCookie, rawToken := issuePortCSRFCookie(t, testMux, "/login")

	encrypted := mustEncryptSession(t, "1")

	before, err := database.CountContacts(rt.db)

	if err != nil {
		t.Fatal(err)
	}

	body := strings.NewReader(url.Values{
		"organization_id": {"1"},
		"first_name":      {"Mina"},
		"last_name":       {"Cole"},
		"email":           {"mina@example.test"},
		"phone":           {"+1 555 0107"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/contacts", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-TOKEN", rawToken)

	req.AddCookie(csrfCookie)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: encrypted})

	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	after, err := database.CountContacts(rt.db)

	if err != nil {
		t.Fatal(err)
	}

	if after != before+1 {
		t.Fatalf("contact count = %d, want %d", after, before+1)
	}
}

func TestStoreNoteAppendsActivity(t *testing.T) {
	rt, testMux := newPortTestMux(t)
	csrfCookie, rawToken := issuePortCSRFCookie(t, testMux, "/login")

	encrypted := mustEncryptSession(t, "1")

	before, err := database.CountNotes(rt.db)

	if err != nil {
		t.Fatal(err)
	}

	body := strings.NewReader(url.Values{
		"body": {"Need legal review before the Friday call."},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/contacts/1/notes", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-TOKEN", rawToken)

	req.AddCookie(csrfCookie)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: encrypted})

	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	after, err := database.CountNotes(rt.db)

	if err != nil {
		t.Fatal(err)
	}

	if after != before+1 {
		t.Fatalf("notes = %d, want %d", after, before+1)
	}
}

func newPortTestMux(t *testing.T) (*runtime, http.Handler) {
	t.Helper()

	testInertia, err := inertia.New(testutil.TestTemplate, inertia.WithVersion("test"))

	if err != nil {
		t.Fatal(err)
	}

	testDB, err := database.Open(":memory:")

	if err != nil {
		t.Fatal(err)
	}

	if err := seed.Run(testDB); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { testDB.Close() })

	cfg := corei18n.DefaultI18n()
	cfg.URLPrefix = false

	rt := &runtime{
		db:         testDB,
		cryptoKey:  testCryptoKey,
		inertia:    testInertia,
		localeCfg:  cfg,
		flashStore: coreflash.NewCookieStore(coreflash.WithCookieName("beacon_flash")),
		routes:     initRoutes(),
	}

	mux := http.NewServeMux()
	authApp, err := rt.newAuth()

	if err != nil {
		t.Fatalf("newAuth() error = %v", err)
	}

	authApp.RegisterRoutes(mux)

	if err := rt.registerCRMRoutes(mux, authApp); err != nil {
		t.Fatalf("registerCRMRoutes() error = %v", err)
	}

	if err := rt.registerFeatureRoutes(mux, authApp); err != nil {
		t.Fatalf("registerFeatureRoutes() error = %v", err)
	}

	if err := rt.registerErrorRoutes(mux, authApp); err != nil {
		t.Fatalf("registerErrorRoutes() error = %v", err)
	}

	mux.Handle("GET /{$}", http.RedirectHandler("/dashboard", http.StatusFound))

	handler := rt.dashboardAppHandler(
		authApp.WithCurrentUser(rt.withDemoProps(authApp, mux)),
		middleware.CSRF(middleware.CSRFConfig{}, []byte("0123456789abcdef0123456789abcdef")),
	)

	return rt, handler
}

func mustEncryptSession(t *testing.T, value string) string {
	t.Helper()

	encrypted := encryptForTest(t, value)

	return encrypted
}

func issuePortCSRFCookie(t *testing.T, handler http.Handler, path string) (*http.Cookie, string) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	return testutil.FindCookie(t, w, "XSRF-TOKEN"), testutil.FindCSRFMetaToken(t, w.Body.String())
}
