package auth

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/bedrock/packages/encryption"
	"github.com/bedrock/packages/inertia"
	"github.com/bedrock/packages/inertia/assert"
	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/services/inertia-demo/api/internal/database"
	"github.com/bedrock/services/inertia-demo/api/internal/seed"
	"github.com/bedrock/services/inertia-demo/api/internal/testutil"
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
	t.Parallel()

	_, handler, _ := newAuthTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)

	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = "/login"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	page := assert.AssertFromBytes(t, w.Body.Bytes())

	page.AssertComponent(t, "Auth/Login")
}

func TestLoginHandlerCreatesSession(t *testing.T) {
	t.Parallel()

	_, handler, testDB := newAuthTestHandler(t)

	user, err := database.FindUserByEmail(testDB, "test@example.com")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	wantID := strconv.FormatInt(user.ID, 10)

	body := strings.NewReader(url.Values{
		"email":    {"test@example.com"},
		"password": {"password"},
		"remember": {"true"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/login", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/dashboard" {
		t.Fatalf("location = %q, want %q", got, "/dashboard")
	}

	cookie := testutil.FindCookie(t, w, SessionCookieName)

	decrypted := decryptForTest(t, cookie.Value)

	if decrypted != wantID {
		t.Fatalf("session cookie value = %q, want %q", decrypted, wantID)
	}
}

func TestLoginHandlerCreatesSessionFromJSON(t *testing.T) {
	t.Parallel()

	_, handler, testDB := newAuthTestHandler(t)

	user, err := database.FindUserByEmail(testDB, "test@example.com")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	wantID := strconv.FormatInt(user.ID, 10)

	body := strings.NewReader(`{"email":"test@example.com","password":"password","remember":true}`)

	req := httptest.NewRequest(http.MethodPost, "/login", body)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/dashboard" {
		t.Fatalf("location = %q, want %q", got, "/dashboard")
	}

	cookie := testutil.FindCookie(t, w, SessionCookieName)

	decrypted := decryptForTest(t, cookie.Value)

	if decrypted != wantID {
		t.Fatalf("session cookie value = %q, want %q", decrypted, wantID)
	}
}

func TestLoginHandlerRejectsInvalidPassword(t *testing.T) {
	t.Parallel()

	_, handler, _ := newAuthTestHandler(t)

	body := strings.NewReader(url.Values{
		"email":    {"test@example.com"},
		"password": {"wrong-password"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/login", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = "/login"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

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

func TestLogoutHandlerClearsSession(t *testing.T) {
	t.Parallel()

	_, handler, testDB := newAuthTestHandler(t)

	user, err := database.FindUserByEmail(testDB, "test@example.com")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	encrypted := encryptForTest(t, strconv.FormatInt(user.ID, 10))

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)

	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: encrypted})

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/login" {
		t.Fatalf("location = %q, want %q", got, "/login")
	}

	if testutil.FindCookie(t, w, SessionCookieName).MaxAge != -1 {
		t.Fatalf("logout should clear the session cookie")
	}
}

func TestWithCurrentUserLoadsUserFromCookie(t *testing.T) {
	t.Parallel()

	app, _, testDB := newAuthTestHandler(t)

	handler := app.WithCurrentUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.CurrentUser(r)

		if user == nil {
			t.Fatal("current user should be present")
		}

		if user.Email != "test@example.com" {
			t.Fatalf("email = %q, want %q", user.Email, "test@example.com")
		}
	}))

	user, err := database.FindUserByEmail(testDB, "test@example.com")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	encrypted := encryptForTest(t, strconv.FormatInt(user.ID, 10))

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)

	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: encrypted})

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}

func TestForgedCookieIsRejected(t *testing.T) {
	t.Parallel()

	app, _, _ := newAuthTestHandler(t)

	handler := app.WithCurrentUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.CurrentUser(r) != nil {
			t.Fatal("forged cookie should not authenticate a user")
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)

	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "1"})

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}

func newAuthTestHandler(t *testing.T) (App, http.Handler, *sql.DB) {
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

	t.Cleanup(func() {
		testDB.Close()
	})

	renderFn := func(w http.ResponseWriter, r *http.Request, component string, pageProps protocol.Props) {
		t.Helper()

		if err := testInertia.Render(w, r, component, pageProps); err != nil {
			t.Fatalf("render failed: %v", err)
		}
	}

	redirectFn := func(w http.ResponseWriter, r *http.Request, url string) {
		testInertia.Redirect(w, r, url)
	}

	setFlashFn := func(http.ResponseWriter, flash.Message) error { return nil }

	app, err := NewApp(Container{
		DB:        testDB,
		CryptoKey: testCryptoKey,
		Render:    renderFn,

		Redirect: redirectFn,

		RouteURL: authTestRouteURL,
		SetFlash: setFlashFn,
	})

	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}

	mux := http.NewServeMux()

	app.RegisterRoutes(mux)

	return app, app.WithCurrentUser(mux), testDB
}

func authTestRouteURL(name string, params map[string]string) string {
	pattern := map[string]string{
		"login":     "/login",
		"logout":    "/logout",
		"dashboard": "/dashboard",
	}[name]

	if strings.TrimSpace(pattern) == "" {
		return "#!wayfinder:unknown-route"
	}

	for key, value := range params {
		pattern = strings.ReplaceAll(pattern, "{"+key+"}", value)
	}

	return pattern
}
