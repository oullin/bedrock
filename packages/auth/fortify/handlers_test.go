package fortify_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/fortify"
	"github.com/gollin/packages/auth/internal/totp"
)

func TestRegisterLoginProtectedLogoutFlow(t *testing.T) {
	t.Parallel()

	manager, mailer, client, server := newHarness(t)
	defer server.Close()

	postJSON(t, client, server.URL+"/register", map[string]any{
		"email":                 "flow@example.com",
		"password":              "password-123",
		"password_confirmation": "password-123",
	}, http.StatusCreated)

	if len(mailer.Messages()) != 1 {
		t.Fatalf("expected verification mail, got %d", len(mailer.Messages()))
	}

	postJSON(t, client, server.URL+"/login", map[string]any{
		"email":    "flow@example.com",
		"password": "password-123",
	}, http.StatusOK)

	getJSON(t, client, server.URL+"/protected", http.StatusOK)
	postJSON(t, client, server.URL+"/logout", map[string]any{}, http.StatusOK)
	getJSON(t, client, server.URL+"/protected", http.StatusUnauthorized)

	_ = manager
}

func TestForgotResetAndLoginFlow(t *testing.T) {
	t.Parallel()

	_, mailer, client, server := newHarness(t)
	defer server.Close()

	postJSON(t, client, server.URL+"/register", map[string]any{
		"email":                 "reset-http@example.com",
		"password":              "password-123",
		"password_confirmation": "password-123",
	}, http.StatusCreated)

	postJSON(t, client, server.URL+"/forgot-password", map[string]any{
		"email": "reset-http@example.com",
	}, http.StatusAccepted)

	messages := mailer.Messages()
	token := messages[len(messages)-1].Metadata["token"]
	if token == "" {
		t.Fatal("expected reset token metadata")
	}

	postJSON(t, client, server.URL+"/reset-password", map[string]any{
		"email":                 "reset-http@example.com",
		"token":                 token,
		"password":              "new-password-123",
		"password_confirmation": "new-password-123",
	}, http.StatusOK)

	postJSON(t, client, server.URL+"/login", map[string]any{
		"email":    "reset-http@example.com",
		"password": "new-password-123",
	}, http.StatusOK)
}

func TestVerificationResendAndVerifiedMiddleware(t *testing.T) {
	t.Parallel()

	_, mailer, client, server := newHarness(t)
	defer server.Close()

	postJSON(t, client, server.URL+"/register", map[string]any{
		"email":                 "verify-http@example.com",
		"password":              "password-123",
		"password_confirmation": "password-123",
	}, http.StatusCreated)
	postJSON(t, client, server.URL+"/login", map[string]any{
		"email":    "verify-http@example.com",
		"password": "password-123",
	}, http.StatusOK)

	getJSON(t, client, server.URL+"/verified", http.StatusForbidden)
	postJSON(t, client, server.URL+"/email/verification-notification", map[string]any{}, http.StatusAccepted)

	link := mailer.Messages()[len(mailer.Messages())-1].Metadata["link"]
	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("Parse verification link: %v", err)
	}

	getJSON(t, client, server.URL+parsed.RequestURI(), http.StatusOK)
	getJSON(t, client, server.URL+"/verified", http.StatusOK)
}

func TestTwoFactorChallengeAndPasswordConfirmationGate(t *testing.T) {
	t.Parallel()

	_, _, client, server := newHarness(t)
	defer server.Close()

	postJSON(t, client, server.URL+"/register", map[string]any{
		"email":                 "2fa-http@example.com",
		"password":              "password-123",
		"password_confirmation": "password-123",
	}, http.StatusCreated)
	postJSON(t, client, server.URL+"/login", map[string]any{
		"email":    "2fa-http@example.com",
		"password": "password-123",
	}, http.StatusOK)

	postJSON(t, client, server.URL+"/user/two-factor-authentication", map[string]any{}, http.StatusForbidden)
	postJSON(t, client, server.URL+"/user/confirm-password", map[string]any{
		"password": "password-123",
	}, http.StatusOK)

	enableBody := postJSON(t, client, server.URL+"/user/two-factor-authentication", map[string]any{}, http.StatusCreated)
	secret := enableBody["secret"].(string)
	recoveryCodes := enableBody["recoveryCodes"].([]any)
	if len(recoveryCodes) == 0 {
		t.Fatal("expected recovery codes")
	}

	postJSON(t, client, server.URL+"/logout", map[string]any{}, http.StatusOK)
	loginBody := postJSON(t, client, server.URL+"/login", map[string]any{
		"email":    "2fa-http@example.com",
		"password": "password-123",
	}, http.StatusOK)
	if !loginBody["requiresTwoFactor"].(bool) {
		t.Fatal("expected two-factor challenge requirement")
	}

	code, err := totp.Code(secret, time.Now().UTC())
	if err != nil {
		t.Fatalf("totp.Code: %v", err)
	}

	postJSON(t, client, server.URL+"/two-factor-challenge", map[string]any{
		"code": code,
	}, http.StatusOK)

	getJSON(t, client, server.URL+"/protected", http.StatusOK)

	postJSON(t, client, server.URL+"/logout", map[string]any{}, http.StatusOK)
	postJSON(t, client, server.URL+"/login", map[string]any{
		"email":    "2fa-http@example.com",
		"password": "password-123",
	}, http.StatusOK)
	postJSON(t, client, server.URL+"/two-factor-challenge", map[string]any{
		"recovery_code": recoveryCodes[0].(string),
	}, http.StatusOK)
}

func newHarness(t *testing.T) (*auth.Manager, *auth.InMemoryMailer, *http.Client, *httptest.Server) {
	t.Helper()

	mailer := &auth.InMemoryMailer{}
	manager, err := auth.NewManager(auth.Config{
		Cookies: auth.CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
			HTTPOnly:     true,
		},
		SigningKey: []byte("test-signing-key"),
	}, auth.Dependencies{
		Mailer: mailer,
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	mux := http.NewServeMux()
	fortify.RegisterRoutes(mux, manager, fortify.RouteOptions{})
	mux.Handle("GET /protected", manager.RequireAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})))
	mux.Handle("GET /verified", manager.RequireVerifiedEmail(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"verified"}`))
	})))

	server := httptest.NewServer(mux)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New: %v", err)
	}

	client := server.Client()
	client.Jar = jar
	return manager, mailer, client, server
}

func postJSON(t *testing.T, client *http.Client, endpoint string, payload map[string]any, wantStatus int) map[string]any {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("NewRequestWithContext: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("client.Do: %v", err)
	}
	defer resp.Body.Close()

	return decodeBody(t, resp, wantStatus)
}

func getJSON(t *testing.T, client *http.Client, endpoint string, wantStatus int) map[string]any {
	t.Helper()

	resp, err := client.Get(endpoint)
	if err != nil {
		t.Fatalf("client.Get: %v", err)
	}
	defer resp.Body.Close()

	return decodeBody(t, resp, wantStatus)
}

func decodeBody(t *testing.T, resp *http.Response, wantStatus int) map[string]any {
	t.Helper()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if resp.StatusCode != wantStatus {
		t.Fatalf("unexpected status %d want %d body=%s", resp.StatusCode, wantStatus, strings.TrimSpace(string(raw)))
	}

	if len(raw) == 0 {
		return map[string]any{}
	}

	body := map[string]any{}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("json.Unmarshal: %v body=%s", err, string(raw))
	}
	return body
}
