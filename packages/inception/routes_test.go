package inception

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bedrock/packages/inception/twofactor"
)

func buildFullInception() *Inception {
	user := &testTwoFactorUser{
		testUser:     testUser{id: "1", email: "user@example.com", password: "hashed"},
		emailAddress: "user@example.com",
	}

	provider := &testProvider{user: user, validPassword: "secret"}

	return &Inception{
		config:        DefaultConfig(),
		guard:         &testGuard{authenticatedUser: user},
		provider:      provider,
		hasher:        &stubHasher{},
		broker:        &testBroker{},
		verifier:      &testVerifier{},
		events:        &testEvents{},
		limiter:       nil,
		responder:     &testResponder{},
		createUser:    &testCreatesUsers{returnUser: user},
		updateProfile: &testUpdatesProfile{},
		updatePass:    &testUpdatesPasswords{},
		resetPass:     &stubResets{},
		confirmPass:   &stubConfirms{},
	}
}

func TestRegisterRoutesAllFeatures(t *testing.T) {
	f := buildFullInception()
	mux := http.NewServeMux()
	router := &StdMuxRouter{Mux: mux}

	RegisterRoutes(router, f, "Bedrock", RouteConfig{})

	routes := []struct {
		method string
		path   string
		want   int
	}{
		{"POST", "/login", http.StatusOK},
		{"POST", "/logout", http.StatusOK},
		{"POST", "/register", http.StatusCreated},
		{"POST", "/forgot-password", http.StatusOK},
		{"POST", "/reset-password", http.StatusUnprocessableEntity}, // no token
		{"POST", "/email/verification-notification", http.StatusOK},
		{"POST", "/user/two-factor-authentication", http.StatusOK},
		{"DELETE", "/user/two-factor-authentication", http.StatusOK},
		{"GET", "/user/two-factor-qr-code", http.StatusBadRequest}, // no secret
		{"POST", "/two-factor-challenge", http.StatusUnprocessableEntity},
		{"PUT", "/user/password", http.StatusUnprocessableEntity},
		{"POST", "/user/confirm-password", http.StatusUnprocessableEntity},
		{"PUT", "/user/profile-information", http.StatusOK},
	}

	for _, tc := range routes {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			var body *strings.Reader

			switch tc.path {
			case "/login":
				body = strings.NewReader("email=user@example.com&password=secret")
			case "/register":
				body = strings.NewReader("email=new@example.com&name=Test&password=pass")
			case "/forgot-password":
				body = strings.NewReader("email=user@example.com")
			default:
				body = strings.NewReader("")
			}

			req := httptest.NewRequest(tc.method, tc.path, body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.RemoteAddr = "127.0.0.1:1234"

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code == http.StatusMethodNotAllowed || w.Code == 404 {
				t.Fatalf("route %s %s not registered, got %d", tc.method, tc.path, w.Code)
			}
		})
	}
}

func TestRegisterRoutesMinimalFeatures(t *testing.T) {
	f := buildFullInception()
	f.config.Features = Features{}

	mux := http.NewServeMux()
	router := &StdMuxRouter{Mux: mux}

	RegisterRoutes(router, f, "Bedrock", RouteConfig{})

	disabledRoutes := []struct {
		method string
		path   string
	}{
		{"POST", "/register"},
		{"POST", "/forgot-password"},
		{"POST", "/reset-password"},
		{"POST", "/email/verification-notification"},
		{"POST", "/user/two-factor-authentication"},
		{"PUT", "/user/password"},
		{"POST", "/user/confirm-password"},
		{"PUT", "/user/profile-information"},
	}

	for _, tc := range disabledRoutes {
		t.Run(tc.method+" "+tc.path+" disabled", func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			// Should get 404 (not registered) not 200
			if w.Code == http.StatusOK || w.Code == http.StatusCreated {
				t.Fatalf("route %s %s should not be registered when feature disabled, got %d", tc.method, tc.path, w.Code)
			}
		})
	}

	// Login/logout should always be registered
	alwaysRoutes := []struct {
		method string
		path   string
	}{
		{"POST", "/login"},
		{"POST", "/logout"},
	}

	for _, tc := range alwaysRoutes {
		t.Run(tc.method+" "+tc.path+" always", func(t *testing.T) {
			body := strings.NewReader("email=user@example.com&password=secret")
			req := httptest.NewRequest(tc.method, tc.path, body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.RemoteAddr = "127.0.0.1:1234"

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code == 404 || w.Code == http.StatusMethodNotAllowed {
				t.Fatalf("route %s %s should always be registered, got %d", tc.method, tc.path, w.Code)
			}
		})
	}
}

func TestFullLoginLifecycle(t *testing.T) {
	secret, _ := twofactor.GenerateSecret(0)
	now := time.Now()

	user := &testTwoFactorUser{
		testUser:     testUser{id: "1", email: "user@example.com", password: "hashed", twoFA: true},
		secret:       secret,
		codes:        []string{"recovery-code1"},
		confirmedAt:  &now,
		emailAddress: "user@example.com",
	}

	provider := &testProvider{user: user, validPassword: "secret"}
	guard := &testGuard{authenticatedUser: user}
	events := &testEvents{}

	f := &Inception{
		config:    DefaultConfig(),
		guard:     guard,
		provider:  provider,
		hasher:    &stubHasher{},
		events:    events,
		responder: &testResponder{},
	}

	mux := http.NewServeMux()
	router := &StdMuxRouter{Mux: mux}
	RegisterRoutes(router, f, "Bedrock", RouteConfig{})

	// Step 1: Login -> should get 2FA challenge (user has 2FA)
	loginBody := strings.NewReader("email=user@example.com&password=secret")
	loginReq := httptest.NewRequest("POST", "/login", loginBody)
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginReq.RemoteAddr = "127.0.0.1:1234"
	loginW := httptest.NewRecorder()
	mux.ServeHTTP(loginW, loginReq)

	if !guard.pendingTwoFactor {
		t.Fatal("step 1: expected pending two factor after login with 2FA user")
	}

	// Step 2: Complete 2FA challenge with TOTP
	code := twofactor.CurrentCode(secret)
	challengeBody := strings.NewReader("code=" + code)
	challengeReq := httptest.NewRequest("POST", "/two-factor-challenge", challengeBody)
	challengeReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	challengeW := httptest.NewRecorder()
	mux.ServeHTTP(challengeW, challengeReq)

	if guard.loggedIn == nil {
		t.Fatal("step 2: expected user to be fully logged in after 2FA challenge")
	}

	// Step 3: Logout
	logoutReq := httptest.NewRequest("POST", "/logout", nil)
	logoutW := httptest.NewRecorder()
	mux.ServeHTTP(logoutW, logoutReq)

	if !guard.loggedOut {
		t.Fatal("step 3: expected logout")
	}
}
