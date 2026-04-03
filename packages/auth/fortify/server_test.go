package fortify_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/authflows"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/memory"
	"github.com/gollin/packages/auth/passwords"
	"github.com/gollin/packages/config/foundation/configuration"
)

func TestAuthFlowsRegisterLoginVerificationAndProfileFlow(t *testing.T) {
	t.Parallel()

	repo, err := configuration.NewBuilder("/Users/gocanto/Sites/gollin/packages/config/config").Build(context.Background())
	if err != nil {
		t.Fatalf("Build config: %v", err)
	}
	authConfig, err := auth.ConfigFromRepository(repo)
	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}
	users := memory.NewInMemoryUserRepository()
	sessions := memory.NewInMemorySessionStore()
	mailer := &memory.InMemoryMailer{}
	manager, err := auth.NewManager(authConfig, map[string]auth.UserProvider{"users": users}, sessions, auth.ManagerDependencies{
		Hasher: auth.DefaultPasswordHasher{},
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	passwordConfig, err := passwords.ConfigFromRepository(repo, "users")
	if err != nil {
		t.Fatalf("password config: %v", err)
	}
	broker := &passwords.Broker{
		Config: passwordConfig,
		Users:  users,
		Tokens: memory.NewInMemoryTokenRepository(),
		Mailer: mailer,
		Clock:  auth.SystemClock{},
	}
	verification := &foundation.VerificationService{
		Config: authConfig,
		Users:  users,
		Signer: auth.HMACLinkSigner{Key: authConfig.SigningKey, Clock: auth.SystemClock{}},
		Mailer: mailer,
		Clock:  auth.SystemClock{},
	}
	server, err := authflows.NewServerFromRepository(repo, authflows.Dependencies{
		AuthManager:    manager,
		Users:          users,
		PasswordBroker: broker,
		Verification:   verification,
	})
	if err != nil {
		t.Fatalf("NewServerFromRepository: %v", err)
	}

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)
	httpServer := httptest.NewServer(mux)
	defer httpServer.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New: %v", err)
	}
	client := httpServer.Client()
	client.Jar = jar

	postJSON(t, client, httpServer.URL+"/register", map[string]any{
		"name":                  "Test User",
		"email":                 "flow@example.com",
		"password":              "password-123",
		"password_confirmation": "password-123",
	}, http.StatusCreated)

	postJSON(t, client, httpServer.URL+"/login", map[string]any{
		"email":    "flow@example.com",
		"password": "password-123",
	}, http.StatusOK)

	postJSON(t, client, httpServer.URL+"/user/confirm-password", map[string]any{
		"password": "password-123",
	}, http.StatusOK)

	verifyLink := mailer.Messages()[0].Metadata["link"]
	parsed, err := url.Parse(verifyLink)
	if err != nil {
		t.Fatalf("url.Parse: %v", err)
	}
	getJSON(t, client, httpServer.URL+parsed.RequestURI(), http.StatusOK)

	putJSON(t, client, httpServer.URL+"/user/profile-information", map[string]any{
		"name":  "Updated User",
		"email": "updated@example.com",
	}, http.StatusOK)

	postJSON(t, client, httpServer.URL+"/logout", map[string]any{}, http.StatusOK)
}

func postJSON(t *testing.T, client *http.Client, endpoint string, payload map[string]any, wantStatus int) map[string]any {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()
	return decodeBody(t, resp, wantStatus)
}

func putJSON(t *testing.T, client *http.Client, endpoint string, payload map[string]any, wantStatus int) map[string]any {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, endpoint, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()
	return decodeBody(t, resp, wantStatus)
}

func getJSON(t *testing.T, client *http.Client, endpoint string, wantStatus int) map[string]any {
	t.Helper()
	resp, err := client.Get(endpoint)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()
	return decodeBody(t, resp, wantStatus)
}

func decodeBody(t *testing.T, resp *http.Response, wantStatus int) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if resp.StatusCode != wantStatus {
		t.Fatalf("unexpected status %d want %d body=%s", resp.StatusCode, wantStatus, strings.TrimSpace(toJSON(t, body)))
	}
	return body
}

func toJSON(t *testing.T, value any) string {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal body: %v", err)
	}
	return string(raw)
}
