package fortify_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/authflows"
	"github.com/gollin/packages/auth/authflows/contracts"
	"github.com/gollin/packages/auth/authflows/responses"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/memory"
	"github.com/gollin/packages/auth/passwords"
	"github.com/gollin/packages/auth/support/otp"
	configpkg "github.com/gollin/packages/config"
	"github.com/gollin/packages/config/foundation/configuration"
)

// Confirmation is required before two-factor login challenges begin.

type fortifyEnv struct {
	URL    string
	Client *http.Client
	Clock  *memory.FixedClock
	Mailer *memory.InMemoryMailer
	Close  func()
}

type testResponse struct {
	status  int
	payload map[string]any
}

func TestAuthFlowsViewsDisabledAndRouteOverrides(t *testing.T) {
	t.Parallel()

	registry := responses.DefaultRegistry()
	registry.LoginViewResponse = testResponse{status: http.StatusTeapot, payload: map[string]any{"view": "custom-login"}}

	env := newAuthFlowsEnv(t, func(repo *configpkg.Repository) {
		repo.Set(map[string]any{
			"authflows.paths.login": "/sign-in",
			"authflows.views":       true,
		})
	}, func(deps *authflows.Dependencies) {
		deps.Responses = registry
	})

	defer env.Close()

	getJSON(t, env.Client, env.URL+"/sign-in", http.StatusTeapot)

	resp, err := env.Client.Get(env.URL + "/login")

	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for default route, got %d", resp.StatusCode)
	}

	envNoViews := newAuthFlowsEnv(t, func(repo *configpkg.Repository) {
		repo.Set("authflows.views", false)
	}, nil)

	defer envNoViews.Close()

	resp, err = envNoViews.Client.Get(envNoViews.URL + "/login")

	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 when views are disabled, got %d", resp.StatusCode)
	}

	postJSON(t, envNoViews.Client, envNoViews.URL+"/login", map[string]any{}, http.StatusUnprocessableEntity)
}

func TestAuthFlowsConfirmedPasswordStatusBranches(t *testing.T) {
	t.Parallel()

	env := newAuthFlowsEnv(t, nil, nil)

	defer env.Close()

	postJSON(t, env.Client, env.URL+"/register", map[string]any{
		"name":                  "Test User",
		"email":                 "status@example.com",
		"password":              "password-123",
		"password_confirmation": "password-123",
	}, http.StatusCreated)
	postJSON(t, env.Client, env.URL+"/login", map[string]any{
		"email":    "status@example.com",
		"password": "password-123",
	}, http.StatusOK)

	body := getJSON(t, env.Client, env.URL+"/user/confirmed-password-status", http.StatusOK)

	if body["confirmed"] != false {
		t.Fatalf("expected unconfirmed status, got %#v", body)
	}

	postJSON(t, env.Client, env.URL+"/user/confirm-password", map[string]any{
		"password": "password-123",
	}, http.StatusOK)

	body = getJSON(t, env.Client, env.URL+"/user/confirmed-password-status", http.StatusOK)

	if body["confirmed"] != true {
		t.Fatalf("expected confirmed status, got %#v", body)
	}

	env.Clock.Advance(4 * time.Hour)
	body = getJSON(t, env.Client, env.URL+"/user/confirmed-password-status", http.StatusOK)

	if body["confirmed"] != false {
		t.Fatalf("expected stale confirmation to be false, got %#v", body)
	}
}

func TestAuthFlowsTwoFactorConfirmationAndChallengeFlow(t *testing.T) {
	t.Parallel()

	env := newAuthFlowsEnv(t, nil, nil)

	defer env.Close()

	postJSON(t, env.Client, env.URL+"/register", map[string]any{
		"name":                  "Two Factor User",
		"email":                 "2fa@example.com",
		"password":              "password-123",
		"password_confirmation": "password-123",
	}, http.StatusCreated)
	postJSON(t, env.Client, env.URL+"/login", map[string]any{
		"email":    "2fa@example.com",
		"password": "password-123",
	}, http.StatusOK)
	postJSON(t, env.Client, env.URL+"/user/confirm-password", map[string]any{
		"password": "password-123",
	}, http.StatusOK)

	enabled := postJSON(t, env.Client, env.URL+"/user/two-factor-authentication", map[string]any{}, http.StatusCreated)
	secret, _ := enabled["secret"].(string)

	if secret == "" {
		t.Fatalf("expected two-factor secret, got %#v", enabled)
	}

	postJSON(t, env.Client, env.URL+"/logout", map[string]any{}, http.StatusOK)

	loggedIn := postJSON(t, env.Client, env.URL+"/login", map[string]any{
		"email":    "2FA@example.com",
		"password": "password-123",
	}, http.StatusOK)

	if loggedIn["status"] != "authenticated" {
		t.Fatalf("expected direct authentication before confirmation, got %#v", loggedIn)
	}

	postJSON(t, env.Client, env.URL+"/user/confirm-password", map[string]any{
		"password": "password-123",
	}, http.StatusOK)
	code, err := otp.Code(secret, env.Clock.Now())

	if err != nil {
		t.Fatalf("Code: %v", err)
	}

	postJSON(t, env.Client, env.URL+"/user/confirmed-two-factor-authentication", map[string]any{
		"code": code,
	}, http.StatusOK)

	postJSON(t, env.Client, env.URL+"/logout", map[string]any{}, http.StatusOK)
	challenge := postJSON(t, env.Client, env.URL+"/login", map[string]any{
		"email":    "2FA@example.com",
		"password": "password-123",
		"remember": true,
	}, http.StatusOK)

	if challenge["status"] != "two_factor_required" {
		t.Fatalf("expected two-factor challenge, got %#v", challenge)
	}

	code, err = otp.Code(secret, env.Clock.Now())

	if err != nil {
		t.Fatalf("Code: %v", err)
	}

	authenticated := postJSON(t, env.Client, env.URL+"/two-factor-challenge", map[string]any{
		"code": code,
	}, http.StatusOK)

	if authenticated["status"] != "authenticated" {
		t.Fatalf("expected completed authentication, got %#v", authenticated)
	}
}

func TestAuthFlowsPasswordResetFailureModes(t *testing.T) {
	t.Parallel()

	env := newAuthFlowsEnv(t, nil, nil)

	defer env.Close()

	postJSON(t, env.Client, env.URL+"/register", map[string]any{
		"name":                  "Reset User",
		"email":                 "reset@example.com",
		"password":              "password-123",
		"password_confirmation": "password-123",
	}, http.StatusCreated)

	postJSON(t, env.Client, env.URL+"/forgot-password", map[string]any{
		"email": "missing@example.com",
	}, http.StatusUnprocessableEntity)

	postJSON(t, env.Client, env.URL+"/forgot-password", map[string]any{
		"email": "RESET@example.com",
	}, http.StatusAccepted)

	token := env.Mailer.Messages()[1].Metadata["token"]
	postJSON(t, env.Client, env.URL+"/reset-password", map[string]any{
		"email":                 "reset@example.com",
		"token":                 token,
		"password":              "new-password-123",
		"password_confirmation": "different-password",
	}, http.StatusUnprocessableEntity)

	postJSON(t, env.Client, env.URL+"/reset-password", map[string]any{
		"email":                 "reset@example.com",
		"token":                 "invalid",
		"password":              "new-password-123",
		"password_confirmation": "new-password-123",
	}, http.StatusUnprocessableEntity)
}

func newAuthFlowsEnv(t *testing.T, mutateRepo func(*configpkg.Repository), mutateDeps func(*authflows.Dependencies)) fortifyEnv {
	t.Helper()

	repo, err := configuration.NewBuilder("/Users/gocanto/Sites/gollin/packages/config/config").Build(context.Background())

	if err != nil {
		t.Fatalf("Build config: %v", err)
	}

	if mutateRepo != nil {
		mutateRepo(repo)
	}

	authConfig, err := auth.ConfigFromRepository(repo)

	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	users := memory.NewInMemoryUserRepository()
	sessions := memory.NewInMemorySessionStore()
	mailer := &memory.InMemoryMailer{}
	clock := memory.NewFixedClock(time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC))
	manager, err := auth.NewManager(authConfig, map[string]auth.UserProvider{"users": users}, sessions, auth.ManagerDependencies{
		Hasher: auth.DefaultPasswordHasher{},
		Clock:  clock,
		IDs:    memory.NewSequenceIDGenerator("session"),
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
		Clock:  clock,
	}
	verification := &foundation.VerificationService{
		Config: authConfig,
		Users:  users,
		Signer: auth.HMACLinkSigner{Key: authConfig.SigningKey, Clock: clock},
		Mailer: mailer,
		Clock:  clock,
	}
	deps := authflows.Dependencies{
		AuthManager:    manager,
		Users:          users,
		PasswordBroker: broker,
		Verification:   verification,
		Clock:          clock,
	}

	if mutateDeps != nil {
		mutateDeps(&deps)
	}

	server, err := authflows.NewServerFromRepository(repo, deps)

	if err != nil {
		t.Fatalf("NewServerFromRepository: %v", err)
	}

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)
	httpServer := httptest.NewServer(mux)

	jar, err := cookiejar.New(nil)

	if err != nil {
		t.Fatalf("cookiejar.New: %v", err)
	}

	client := httpServer.Client()
	client.Jar = jar

	return fortifyEnv{
		URL:    httpServer.URL,
		Client: client,
		Clock:  clock,
		Mailer: mailer,
		Close:  httpServer.Close,
	}
}

func (r testResponse) ToResponse(w http.ResponseWriter, _ *http.Request, _ any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.status)
	_ = json.NewEncoder(w).Encode(r.payload)
}

var _ contracts.LoginViewResponse = testResponse{}
