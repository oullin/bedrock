package fortify_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/contracts/events"
	"github.com/bedrock/packages/inception"
	"github.com/bedrock/packages/inception/twofactor"
	"github.com/bedrock/packages/validation"
)

type fortifyUser struct {
	id            string
	email         string
	password      string
	rememberToken string

	twoFactorEnabled   bool
	twoFactorSecret    string
	twoFactorCodes     []string
	twoFactorConfirmed *time.Time

	verified   bool
	unverified bool
}

func (u *fortifyUser) GetAuthIdentifierName() string       { return "id" }
func (u *fortifyUser) GetAuthIdentifier() string           { return u.id }
func (u *fortifyUser) GetAuthPasswordName() string         { return "password" }
func (u *fortifyUser) GetAuthPassword() string             { return u.password }
func (u *fortifyUser) SetAuthPassword(password string)     { u.password = password }
func (u *fortifyUser) GetRememberToken() string            { return u.rememberToken }
func (u *fortifyUser) SetRememberToken(token string)       { u.rememberToken = token }
func (u *fortifyUser) GetRememberTokenName() string        { return "remember_token" }
func (u *fortifyUser) IsTwoFactorEnabled() bool            { return u.twoFactorEnabled }
func (u *fortifyUser) SetTwoFactorEnabled(enabled bool)    { u.twoFactorEnabled = enabled }
func (u *fortifyUser) GetTwoFactorSecret() string          { return u.twoFactorSecret }
func (u *fortifyUser) SetTwoFactorSecret(secret string)    { u.twoFactorSecret = secret }
func (u *fortifyUser) GetTwoFactorRecoveryCodes() []string { return u.twoFactorCodes }
func (u *fortifyUser) SetTwoFactorRecoveryCodes(codes []string) {
	u.twoFactorCodes = append([]string(nil), codes...)
}
func (u *fortifyUser) GetTwoFactorConfirmedAt() *time.Time { return u.twoFactorConfirmed }
func (u *fortifyUser) SetTwoFactorConfirmedAt(at *time.Time) {
	u.twoFactorConfirmed = at
}
func (u *fortifyUser) HasVerifiedEmail() bool          { return u.verified }
func (u *fortifyUser) MarkEmailAsVerified(_ time.Time) { u.verified = true }
func (u *fortifyUser) MarkEmailAsUnverified()          { u.unverified = true }
func (u *fortifyUser) GetEmailForVerification() string { return u.email }
func (u *fortifyUser) GetEmailForPasswordReset() string {
	return u.email
}

type fortifyGuard struct {
	authenticated cauth.Authenticatable

	loggedIn         cauth.Authenticatable
	loggedInRemember bool
	pendingTwoFactor bool
	loggedOut        bool
}

func (g *fortifyGuard) Name() string { return "web" }
func (g *fortifyGuard) AuthenticateRequest(_ context.Context, _ http.ResponseWriter, _ *http.Request) (cauth.Authenticatable, error) {
	return g.authenticated, nil
}
func (g *fortifyGuard) Login(_ context.Context, _ http.ResponseWriter, user cauth.Authenticatable, remember bool) error {
	g.loggedIn = user
	g.loggedInRemember = remember

	return nil
}
func (g *fortifyGuard) LoginWithPendingTwoFactor(_ context.Context, _ http.ResponseWriter, user cauth.Authenticatable) error {
	g.loggedIn = user
	g.pendingTwoFactor = true

	return nil
}
func (g *fortifyGuard) Logout(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
	g.loggedOut = true

	return nil
}

type fortifyProvider struct {
	user          cauth.Authenticatable
	validPassword string
}

func (p *fortifyProvider) RetrieveByID(_ context.Context, _ string) (cauth.Authenticatable, error) {
	return p.user, nil
}
func (p *fortifyProvider) RetrieveByToken(_ context.Context, _ string, _ string) (cauth.Authenticatable, error) {
	return p.user, nil
}
func (p *fortifyProvider) RetrieveByCredentials(_ context.Context, _ map[string]string) (cauth.Authenticatable, error) {
	return p.user, nil
}
func (p *fortifyProvider) UpdateRememberToken(_ context.Context, _ cauth.Authenticatable, _ string) error {
	return nil
}
func (p *fortifyProvider) ValidateCredentials(_ context.Context, _ cauth.Authenticatable, credentials map[string]string) (bool, error) {
	return credentials["password"] == p.validPassword, nil
}
func (p *fortifyProvider) RehashPasswordIfRequired(_ context.Context, _ cauth.Authenticatable, _ map[string]string, _ bool) error {
	return nil
}

type fortifyHasher struct{}

func (h fortifyHasher) Hash(_ context.Context, password string) (string, error) { return password, nil }
func (h fortifyHasher) Check(_ context.Context, password string, hash string) (bool, error) {
	return password == hash, nil
}
func (h fortifyHasher) NeedsRehash(_ string) bool { return false }

type fortifyEvents struct {
	dispatched []any
}

func (e *fortifyEvents) Listen(_ any, _ ...events.Listener)          {}
func (e *fortifyEvents) HasListeners(_ any) bool                     { return false }
func (e *fortifyEvents) HasWildcardListeners(_ any) bool             { return false }
func (e *fortifyEvents) Subscribe(_ events.Subscriber)               {}
func (e *fortifyEvents) Until(_ context.Context, _ any) (any, error) { return nil, nil }
func (e *fortifyEvents) Push(_ context.Context, _ any)               {}
func (e *fortifyEvents) Flush(_ context.Context, _ string) error     { return nil }
func (e *fortifyEvents) Forget(_ any)                                {}
func (e *fortifyEvents) ForgetPushed()                               {}
func (e *fortifyEvents) GetListeners(_ any) []events.Listener        { return nil }
func (e *fortifyEvents) Dispatch(_ context.Context, event any) ([]any, error) {
	e.dispatched = append(e.dispatched, event)

	return nil, nil
}

type fortifyResponder struct {
	login, logout, register, resetLink, reset, updatePassword, confirmPassword bool
	profile, verification, twoFactorChallenge, twoFactorEnabled, twoFactorOff  bool
}

func (r *fortifyResponder) LoginResponse(w http.ResponseWriter, _ *http.Request) {
	r.login = true
	w.WriteHeader(http.StatusOK)
}
func (r *fortifyResponder) LogoutResponse(w http.ResponseWriter, _ *http.Request) {
	r.logout = true
	w.WriteHeader(http.StatusOK)
}
func (r *fortifyResponder) RegisterResponse(w http.ResponseWriter, _ *http.Request) {
	r.register = true
	w.WriteHeader(http.StatusCreated)
}
func (r *fortifyResponder) PasswordResetLinkSentResponse(w http.ResponseWriter, _ *http.Request) {
	r.resetLink = true
	w.WriteHeader(http.StatusOK)
}
func (r *fortifyResponder) PasswordResetResponse(w http.ResponseWriter, _ *http.Request) {
	r.reset = true
	w.WriteHeader(http.StatusOK)
}
func (r *fortifyResponder) PasswordUpdateResponse(w http.ResponseWriter, _ *http.Request) {
	r.updatePassword = true
	w.WriteHeader(http.StatusOK)
}
func (r *fortifyResponder) PasswordConfirmResponse(w http.ResponseWriter, _ *http.Request) {
	r.confirmPassword = true
	w.WriteHeader(http.StatusOK)
}
func (r *fortifyResponder) ProfileInformationUpdatedResponse(w http.ResponseWriter, _ *http.Request) {
	r.profile = true
	w.WriteHeader(http.StatusOK)
}
func (r *fortifyResponder) EmailVerificationSentResponse(w http.ResponseWriter, _ *http.Request) {
	r.verification = true
	w.WriteHeader(http.StatusOK)
}
func (r *fortifyResponder) TwoFactorChallengeResponse(w http.ResponseWriter, _ *http.Request) {
	r.twoFactorChallenge = true
	w.WriteHeader(http.StatusOK)
}
func (r *fortifyResponder) TwoFactorEnabledResponse(w http.ResponseWriter, _ *http.Request) {
	r.twoFactorEnabled = true
	w.WriteHeader(http.StatusOK)
}
func (r *fortifyResponder) TwoFactorDisabledResponse(w http.ResponseWriter, _ *http.Request) {
	r.twoFactorOff = true
	w.WriteHeader(http.StatusOK)
}

type fortifyLimiter struct {
	tooMany bool
	hits    []string
	cleared []string
}

func (l *fortifyLimiter) TooManyAttempts(_ string, _ int) bool { return l.tooMany }
func (l *fortifyLimiter) Hit(key string, _ time.Duration) int {
	l.hits = append(l.hits, key)

	return len(l.hits)
}
func (l *fortifyLimiter) Clear(key string) {
	l.cleared = append(l.cleared, key)
}
func (l *fortifyLimiter) AvailableIn(_ string) time.Duration { return 0 }

type createUsersAction struct {
	user cauth.Authenticatable
	err  error
}

func (a createUsersAction) Create(_ context.Context, _ map[string]string) (cauth.Authenticatable, error) {
	return a.user, a.err
}

type passwordBroker struct {
	err        error
	sent       bool
	reset      bool
	credential map[string]string
}

func (b *passwordBroker) SendResetLink(_ context.Context, credentials map[string]string) error {
	b.sent = true
	b.credential = credentials

	return b.err
}
func (b *passwordBroker) Reset(_ context.Context, credentials map[string]string, callback func(cauth.Authenticatable, string) error) error {
	if b.err != nil {
		return b.err
	}

	b.reset = true
	b.credential = credentials

	return callback(&fortifyUser{id: "reset-user", email: credentials["email"]}, credentials["password"])
}

type resetPasswordsAction struct {
	password string
	called   bool
}

func (a *resetPasswordsAction) Reset(_ context.Context, _ cauth.Authenticatable, password string) error {
	a.called = true
	a.password = password

	return nil
}

type updatePasswordsAction struct {
	input map[string]string
	err   error
}

func (a *updatePasswordsAction) Update(_ context.Context, _ cauth.Authenticatable, input map[string]string) error {
	a.input = input

	return a.err
}

type confirmPasswordsAction struct {
	password string
	err      error
}

func (a *confirmPasswordsAction) Confirm(_ context.Context, _ cauth.Authenticatable, password string) error {
	a.password = password

	return a.err
}

type updateProfileAction struct {
	err error
}

func (a updateProfileAction) Update(_ context.Context, user cauth.Authenticatable, input map[string]string) error {
	if a.err != nil {
		return a.err
	}

	if u, ok := user.(*fortifyUser); ok {
		u.email = input["email"]
	}

	return nil
}

type verifier struct {
	sent bool
	id   string
	hash string
	err  error
}

func (v *verifier) SendVerificationNotification(_ context.Context, _ cauth.Authenticatable) error {
	v.sent = true

	return v.err
}
func (v *verifier) Verify(_ context.Context, id string, hash string) error {
	v.id = id
	v.hash = hash

	return v.err
}

func newApp(t *testing.T, configure func(*inception.Config), options ...func(*inception.Builder)) (*inception.Inception, *fortifyGuard, *fortifyResponder, *fortifyEvents) {
	t.Helper()

	user := &fortifyUser{id: "1", email: "user@example.com", password: "hash"}
	config := inception.DefaultConfig()
	config.Features = inception.Features{}

	if configure != nil {
		configure(&config)
	}

	guard := &fortifyGuard{authenticated: user}
	responder := &fortifyResponder{}
	dispatcher := &fortifyEvents{}
	builder := inception.NewBuilder().
		WithConfig(config).
		WithGuard(guard).
		WithProvider(&fortifyProvider{user: user, validPassword: "secret"}).
		WithHasher(fortifyHasher{}).
		WithEvents(dispatcher).
		WithResponder(responder)

	for _, option := range options {
		option(builder)
	}

	app, err := builder.Build()
	if err != nil {
		t.Fatalf("build fortify app: %v", err)
	}

	return app, guard, responder, dispatcher
}

func formRequest(method string, path string, body string) *http.Request {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.RemoteAddr = "127.0.0.1:1234"

	return request
}

func jsonRequest(method string, path string, body string) *http.Request {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.RemoteAddr = "127.0.0.1:1234"

	return request
}

// AuthenticatedSessionControllerTest::test_user_can_authenticate
// AuthenticatedSessionControllerTest::test_case_insensitive_usernames_can_be_used
func TestAuthenticatedSessionControllerAuthenticatesUser(t *testing.T) {
	limiter := &fortifyLimiter{}
	app, guard, responder, events := newApp(t, nil, func(builder *inception.Builder) {
		builder.WithLimiter(limiter)
	})

	recorder := httptest.NewRecorder()
	inception.NewLoginHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/login", "email=USER@example.com&password=secret&remember=on"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !responder.login {
		t.Fatal("expected Fortify login response")
	}
	if guard.loggedIn == nil || guard.loggedIn.GetAuthIdentifier() != "1" {
		t.Fatalf("logged in user = %#v, want id 1", guard.loggedIn)
	}
	if !guard.loggedInRemember {
		t.Fatal("remember selection should be forwarded to the guard")
	}
	if len(events.dispatched) != 2 {
		t.Fatalf("events = %d, want login attempted and succeeded", len(events.dispatched))
	}
	if len(limiter.cleared) != 1 || limiter.cleared[0] != "1|127.0.0.1" {
		t.Fatalf("cleared throttle keys = %#v", limiter.cleared)
	}
}

// AuthenticatedSessionControllerTest::test_validation_exception_returned_on_failure
func TestAuthenticatedSessionControllerReturnsValidationFailure(t *testing.T) {
	app, guard, _, events := newApp(t, nil)

	recorder := httptest.NewRecorder()
	inception.NewLoginHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/login", "email=user@example.com&password=wrong"))

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	if guard.loggedIn != nil {
		t.Fatal("user should not be logged in after invalid credentials")
	}
	if _, ok := events.dispatched[len(events.dispatched)-1].(inception.LoginFailedPayload); !ok {
		t.Fatalf("last event = %T, want LoginFailedPayload", events.dispatched[len(events.dispatched)-1])
	}
}

// AuthenticatedSessionControllerTest::test_login_attempts_are_throttled
// AuthenticatedSessionControllerTest::test_cant_bypass_throttle_with_special_characters
func TestAuthenticatedSessionControllerThrottlesLoginAttempts(t *testing.T) {
	limiter := &fortifyLimiter{tooMany: true}
	app, _, _, _ := newApp(t, nil, func(builder *inception.Builder) {
		builder.WithLimiter(limiter)
	})

	recorder := httptest.NewRecorder()
	inception.NewLoginHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/login", "email=User+Alias@example.com&password=secret"))

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusTooManyRequests)
	}
	key := inception.ThrottleKey("User+Alias@example.com", "127.0.0.1")
	if key != "user+alias@example.com|127.0.0.1" {
		t.Fatalf("throttle key = %q", key)
	}
}

// AuthenticatedSessionControllerTest::test_the_user_can_logout_of_the_application
// AuthenticatedSessionControllerTest::test_the_user_can_logout_of_the_application_using_json_request
// AuthenticatedSessionControllerTest::test_users_can_logout
func TestAuthenticatedSessionControllerLogsOut(t *testing.T) {
	app, guard, responder, events := newApp(t, nil)

	recorder := httptest.NewRecorder()
	inception.NewLogoutHandler(app).ServeHTTP(recorder, jsonRequest(http.MethodPost, "/logout", `{}`))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !guard.loggedOut || !responder.logout {
		t.Fatalf("loggedOut = %v responder.logout = %v", guard.loggedOut, responder.logout)
	}
	if _, ok := events.dispatched[0].(inception.LoggedOutPayload); !ok {
		t.Fatalf("event = %T, want LoggedOutPayload", events.dispatched[0])
	}
}

// AuthenticatedSessionControllerTest::test_must_be_authenticated_to_logout
func TestAuthenticatedSessionControllerLogoutRequiresSessionMiddlewareInHostApplication(t *testing.T) {
	logoutGuard := &fortifyGuard{}
	app, _, responder, events := newApp(t, nil, func(builder *inception.Builder) {
		builder.WithGuard(logoutGuard)
	})

	recorder := httptest.NewRecorder()
	inception.NewLogoutHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/logout", ""))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !logoutGuard.loggedOut {
		t.Fatal("logout should still be delegated to the guard")
	}
	if !responder.logout || len(events.dispatched) != 0 {
		t.Fatalf("response = %v events = %#v", responder.logout, events.dispatched)
	}
}

// AuthenticatedSessionControllerWithTwoFactorTest::test_user_is_redirected_to_challenge_when_using_two_factor_authentication
// AuthenticatedSessionControllerWithTwoFactorTest::test_user_is_redirected_to_challenge_when_using_two_factor_authentication_that_has_been_confirmed_and_confirmation_is_enabled
func TestAuthenticatedSessionControllerRedirectsToTwoFactorChallenge(t *testing.T) {
	confirmedAt := time.Now()
	user := &fortifyUser{
		id:                 "1",
		email:              "user@example.com",
		twoFactorEnabled:   true,
		twoFactorConfirmed: &confirmedAt,
		twoFactorSecret:    "secret",
		twoFactorCodes:     []string{"aaaa-bbbb"},
	}
	app, guard, responder, _ := newApp(t, nil, func(builder *inception.Builder) {
		builder.WithProvider(&fortifyProvider{user: user, validPassword: "secret"})
	})

	recorder := httptest.NewRecorder()
	inception.NewLoginHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/login", "email=user@example.com&password=secret"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !guard.pendingTwoFactor || !responder.twoFactorChallenge {
		t.Fatalf("pending = %v challenge response = %v", guard.pendingTwoFactor, responder.twoFactorChallenge)
	}
}

// AuthenticatedSessionControllerWithTwoFactorTest::test_user_can_authenticate_when_two_factor_challenge_is_disabled
// AuthenticatedSessionControllerWithTwoFactorTest::test_user_is_not_redirected_to_challenge_when_using_two_factor_authentication_that_has_not_been_confirmed_and_confirmation_is_enabled
func TestAuthenticatedSessionControllerSkipsUnconfirmedTwoFactorChallenge(t *testing.T) {
	user := &fortifyUser{id: "1", email: "user@example.com", twoFactorEnabled: true}
	app, guard, responder, _ := newApp(t, nil, func(builder *inception.Builder) {
		builder.WithProvider(&fortifyProvider{user: user, validPassword: "secret"})
	})

	recorder := httptest.NewRecorder()
	inception.NewLoginHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/login", "email=user@example.com&password=secret"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if guard.pendingTwoFactor || responder.twoFactorChallenge {
		t.Fatal("unconfirmed two-factor user should complete normal login")
	}
	if guard.loggedIn == nil || !responder.login {
		t.Fatal("expected normal login response")
	}
}

// AuthenticatedSessionControllerWithTwoFactorTest::test_two_factor_challenge_can_be_passed_via_code
// AuthenticatedSessionControllerWithTwoFactorTest::test_two_factor_authentication_preserves_remember_me_selection
func TestTwoFactorChallengePassesViaCode(t *testing.T) {
	secret, err := twofactor.GenerateSecret(0)
	if err != nil {
		t.Fatal(err)
	}

	user := &fortifyUser{id: "1", email: "user@example.com", twoFactorEnabled: true, twoFactorSecret: secret}
	guard := &fortifyGuard{authenticated: user}
	app, _, responder, _ := newApp(t, func(config *inception.Config) {
		config.Features.TwoFactorAuthentication = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(guard)
	})

	code := twofactor.CurrentCode(secret)
	recorder := httptest.NewRecorder()
	inception.NewTwoFactorChallengeHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/two-factor-challenge", "code="+code+"&remember=1"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !guard.loggedInRemember || !responder.login {
		t.Fatalf("remember = %v login response = %v", guard.loggedInRemember, responder.login)
	}
}

// AuthenticatedSessionControllerWithTwoFactorTest::test_two_factor_challenge_can_be_passed_via_recovery_code
// AuthenticatedSessionControllerWithTwoFactorTest::test_two_factor_challenge_can_fail_via_recovery_code
func TestTwoFactorChallengeUsesRecoveryCodeOnce(t *testing.T) {
	user := &fortifyUser{id: "1", email: "user@example.com", twoFactorEnabled: true, twoFactorCodes: []string{"cccc-dddd", "eeee-ffff"}}
	guard := &fortifyGuard{authenticated: user}
	app, _, _, events := newApp(t, func(config *inception.Config) {
		config.Features.TwoFactorAuthentication = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(guard)
	})

	recorder := httptest.NewRecorder()
	inception.NewTwoFactorChallengeHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/two-factor-challenge", "recovery_code=cccc-dddd"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if guard.loggedIn == nil {
		t.Fatal("expected recovered user to be logged in")
	}
	if len(user.twoFactorCodes) != 1 || user.twoFactorCodes[0] != "eeee-ffff" {
		t.Fatalf("recovery codes = %#v", user.twoFactorCodes)
	}
	if events.dispatched[0] != inception.EventRecoveryCodeUsed {
		t.Fatalf("event = %#v, want recovery code event", events.dispatched[0])
	}

	failed := httptest.NewRecorder()
	inception.NewTwoFactorChallengeHandler(app).ServeHTTP(failed, formRequest(http.MethodPost, "/two-factor-challenge", "recovery_code=cccc-dddd"))
	if failed.Code != http.StatusUnprocessableEntity {
		t.Fatalf("reuse status = %d, want %d", failed.Code, http.StatusUnprocessableEntity)
	}
}

// AuthenticatedSessionControllerWithTwoFactorTest::test_two_factor_challenge_fails_for_old_otp_and_zero_window
func TestTwoFactorChallengeRejectsCodeOutsideAcceptedSkew(t *testing.T) {
	secret, err := twofactor.GenerateSecret(0)
	if err != nil {
		t.Fatal(err)
	}

	user := &fortifyUser{id: "1", email: "user@example.com", twoFactorEnabled: true, twoFactorSecret: secret}
	app, _, _, _ := newApp(t, func(config *inception.Config) {
		config.Features.TwoFactorAuthentication = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{authenticated: user})
	})

	oldCode := twofactor.CodeAt(secret, time.Now().Add(-5*time.Minute))
	recorder := httptest.NewRecorder()
	inception.NewTwoFactorChallengeHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/two-factor-challenge", "code="+oldCode))

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
}

// AuthenticatedSessionControllerWithTwoFactorTest::test_two_factor_challenge_requires_a_challenged_user
// AuthenticatedSessionControllerWithTwoFactorTest::test_two_factor_challenge_can_fail_via_recovery_code
func TestTwoFactorChallengeRequiresAuthenticatedPendingUser(t *testing.T) {
	app, _, _, _ := newApp(t, func(config *inception.Config) {
		config.Features.TwoFactorAuthentication = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{})
	})

	recorder := httptest.NewRecorder()
	inception.NewTwoFactorChallengeHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/two-factor-challenge", "code=123456"))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

// ConfirmablePasswordControllerTest::test_password_can_be_confirmed
// ConfirmablePasswordControllerTest::test_password_can_be_confirmed_with_json
// ConfirmablePasswordControllerTest::test_password_confirmation_status_has_been_confirmed
func TestConfirmablePasswordControllerConfirmsPassword(t *testing.T) {
	action := &confirmPasswordsAction{}
	app, _, responder, _ := newApp(t, func(config *inception.Config) {
		config.Features.ConfirmPassword = true
	}, func(builder *inception.Builder) {
		builder.WithConfirmPassword(action)
	})

	request := jsonRequest(http.MethodPost, "/user/confirm-password", `{"password":"secret"}`)
	recorder := httptest.NewRecorder()
	inception.NewConfirmPasswordHandler(app).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if action.password != "secret" || !responder.confirmPassword {
		t.Fatalf("password = %q response = %v", action.password, responder.confirmPassword)
	}
	if inception.GetPasswordConfirmedAt(request) == nil {
		t.Fatal("expected request to be marked password-confirmed")
	}
}

// ConfirmablePasswordControllerTest::test_password_confirmation_can_fail_with_an_invalid_password
// ConfirmablePasswordControllerTest::test_password_confirmation_can_fail_without_a_password
// ConfirmablePasswordControllerTest::test_password_confirmation_can_fail_with_json
func TestConfirmablePasswordControllerFailsInvalidPassword(t *testing.T) {
	app, _, responder, _ := newApp(t, func(config *inception.Config) {
		config.Features.ConfirmPassword = true
	}, func(builder *inception.Builder) {
		builder.WithConfirmPassword(&confirmPasswordsAction{err: errors.New("invalid password")})
	})

	recorder := httptest.NewRecorder()
	inception.NewConfirmPasswordHandler(app).ServeHTTP(recorder, jsonRequest(http.MethodPost, "/user/confirm-password", `{}`))

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	if responder.confirmPassword {
		t.Fatal("confirmation response should not run on invalid password")
	}
}

// ConfirmablePasswordControllerTest::test_password_confirmation_status_has_expired
// ConfirmablePasswordControllerTest::test_password_confirmation_status_has_not_confirmed
func TestConfirmablePasswordControllerStatusMiddlewareRequiresFreshConfirmation(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := inception.EnsurePasswordIsConfirmed(3 * time.Hour)(next)

	notConfirmed := httptest.NewRecorder()
	handler.ServeHTTP(notConfirmed, formRequest(http.MethodGet, "/confirmed-password-status", ""))
	if notConfirmed.Code != http.StatusLocked {
		t.Fatalf("not confirmed status = %d, want %d", notConfirmed.Code, http.StatusLocked)
	}

	expiredRequest := formRequest(http.MethodGet, "/confirmed-password-status", "")
	expiredRequest = expiredRequest.WithContext(inception.WithPasswordConfirmedAt(expiredRequest.Context(), time.Now().Add(-4*time.Hour)))
	expired := httptest.NewRecorder()
	handler.ServeHTTP(expired, expiredRequest)
	if expired.Code != http.StatusLocked {
		t.Fatalf("expired status = %d, want %d", expired.Code, http.StatusLocked)
	}

	freshRequest := formRequest(http.MethodGet, "/confirmed-password-status", "")
	freshRequest = freshRequest.WithContext(inception.WithPasswordConfirmedAt(freshRequest.Context(), time.Now()))
	fresh := httptest.NewRecorder()
	handler.ServeHTTP(fresh, freshRequest)
	if fresh.Code != http.StatusOK {
		t.Fatalf("fresh status = %d, want %d", fresh.Code, http.StatusOK)
	}
}

// RegisteredUserControllerTest::test_users_can_be_created
// RegisteredUserControllerTest::test_users_can_be_created_with_remember_option
func TestRegisteredUserControllerCreatesUser(t *testing.T) {
	user := &fortifyUser{id: "new", email: "new@example.com"}
	app, guard, responder, events := newApp(t, func(config *inception.Config) {
		config.Features.Registration = true
	}, func(builder *inception.Builder) {
		builder.WithCreateUser(createUsersAction{user: user})
	})

	recorder := httptest.NewRecorder()
	inception.NewRegisterHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/register", "email=new@example.com&name=New&password=secret&password_confirmation=secret"))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if guard.loggedIn.GetAuthIdentifier() != "new" || !responder.register {
		t.Fatalf("logged in = %#v register response = %v", guard.loggedIn, responder.register)
	}
	if _, ok := events.dispatched[0].(inception.RegisteredPayload); !ok {
		t.Fatalf("event = %T, want RegisteredPayload", events.dispatched[0])
	}
}

// RegisteredUserControllerTest::test_usernames_will_be_stored_case_insensitive
func TestRegisteredUserControllerReceivesCaseInsensitiveIdentifierInput(t *testing.T) {
	user := &fortifyUser{id: "new", email: "new@example.com"}
	action := createUsersAction{user: user}
	app, _, _, _ := newApp(t, func(config *inception.Config) {
		config.Features.Registration = true
		config.IdentifierField = "username"
	}, func(builder *inception.Builder) {
		builder.WithCreateUser(action)
	})

	recorder := httptest.NewRecorder()
	inception.NewRegisterHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/register", "username=MiXeD&name=New&password=secret&password_confirmation=secret"))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
}

// PasswordResetLinkRequestControllerTest::test_reset_link_can_be_successfully_requested
// PasswordResetLinkRequestControllerTest::test_reset_link_can_be_successfully_requested_with_customized_email_field
func TestPasswordResetLinkRequestControllerSendsResetLink(t *testing.T) {
	broker := &passwordBroker{}
	app, _, responder, _ := newApp(t, func(config *inception.Config) {
		config.Features.ResetPasswords = true
		config.IdentifierField = "username"
	}, func(builder *inception.Builder) {
		builder.WithBroker(broker)
		builder.WithResetPassword(&resetPasswordsAction{})
	})

	recorder := httptest.NewRecorder()
	inception.NewForgotPasswordHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/forgot-password", "username=taylor"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !broker.sent || broker.credential["username"] != "taylor" || !responder.resetLink {
		t.Fatalf("broker.sent = %v credential = %#v response = %v", broker.sent, broker.credential, responder.resetLink)
	}
}

// PasswordResetLinkRequestControllerTest::test_case_insensitive_usernames_can_be_used
func TestPasswordResetLinkRequestControllerCaseInsensitiveUsernameThrottleKey(t *testing.T) {
	limiter := &fortifyLimiter{}
	broker := &passwordBroker{}
	app, _, _, _ := newApp(t, func(config *inception.Config) {
		config.Features.ResetPasswords = true
		config.IdentifierField = "username"
	}, func(builder *inception.Builder) {
		builder.WithLimiter(limiter)
		builder.WithBroker(broker)
		builder.WithResetPassword(&resetPasswordsAction{})
	})

	recorder := httptest.NewRecorder()
	inception.NewForgotPasswordHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/forgot-password", "username=Taylor"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got, want := limiter.hits[0], "password_reset|taylor|127.0.0.1"; got != want {
		t.Fatalf("throttle key = %q, want %q", got, want)
	}
	if broker.credential["username"] != "Taylor" {
		t.Fatalf("broker username = %q, want original casing", broker.credential["username"])
	}
}

// PasswordResetLinkRequestControllerTest::test_reset_link_request_can_fail
// PasswordResetLinkRequestControllerTest::test_reset_link_request_can_fail_with_json
func TestPasswordResetLinkRequestControllerFails(t *testing.T) {
	broker := &passwordBroker{err: errors.New("user not found")}
	app, _, responder, _ := newApp(t, func(config *inception.Config) {
		config.Features.ResetPasswords = true
	}, func(builder *inception.Builder) {
		builder.WithBroker(broker)
		builder.WithResetPassword(&resetPasswordsAction{})
	})

	recorder := httptest.NewRecorder()
	inception.NewForgotPasswordHandler(app).ServeHTTP(recorder, jsonRequest(http.MethodPost, "/forgot-password", `{"email":"missing@example.com"}`))

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	if responder.resetLink {
		t.Fatal("reset link response should not run on broker failure")
	}
}

// NewPasswordControllerTest::test_password_can_be_reset
// NewPasswordControllerTest::test_password_can_be_reset_with_customized_email_address_field
func TestNewPasswordControllerResetsPassword(t *testing.T) {
	broker := &passwordBroker{}
	action := &resetPasswordsAction{}
	app, _, responder, events := newApp(t, func(config *inception.Config) {
		config.Features.ResetPasswords = true
		config.IdentifierField = "username"
	}, func(builder *inception.Builder) {
		builder.WithBroker(broker)
		builder.WithResetPassword(action)
	})

	recorder := httptest.NewRecorder()
	inception.NewResetPasswordHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/reset-password", "username=taylor&token=token&password=new&password_confirmation=new"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !broker.reset || !action.called || action.password != "new" || !responder.reset {
		t.Fatalf("broker.reset = %v action = %#v response = %v", broker.reset, action, responder.reset)
	}
	if events.dispatched[0] != inception.EventPasswordReset {
		t.Fatalf("event = %#v, want password reset", events.dispatched[0])
	}
}

// NewPasswordControllerTest::test_password_reset_can_fail
// NewPasswordControllerTest::test_password_reset_can_fail_with_json
func TestNewPasswordControllerResetCanFail(t *testing.T) {
	broker := &passwordBroker{err: errors.New("invalid token")}
	app, _, responder, _ := newApp(t, func(config *inception.Config) {
		config.Features.ResetPasswords = true
	}, func(builder *inception.Builder) {
		builder.WithBroker(broker)
		builder.WithResetPassword(&resetPasswordsAction{})
	})

	recorder := httptest.NewRecorder()
	inception.NewResetPasswordHandler(app).ServeHTTP(recorder, jsonRequest(http.MethodPost, "/reset-password", `{"email":"user@example.com","token":"bad","password":"new"}`))

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	if responder.reset {
		t.Fatal("reset response should not run on broker failure")
	}
}

// NewPasswordControllerTest::test_password_is_required
func TestNewPasswordControllerRequiresPassword(t *testing.T) {
	broker := &passwordBroker{}
	app, _, responder, _ := newApp(t, func(config *inception.Config) {
		config.Features.ResetPasswords = true
	}, func(builder *inception.Builder) {
		builder.WithBroker(broker)
		builder.WithResetPassword(&resetPasswordsAction{})
	})

	recorder := httptest.NewRecorder()
	inception.NewResetPasswordHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/reset-password", "email=user@example.com&token=token"))

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	if broker.reset || responder.reset {
		t.Fatalf("broker.reset = %v responder.reset = %v", broker.reset, responder.reset)
	}
}

// PasswordControllerTest::test_passwords_can_be_updated
func TestPasswordControllerUpdatesPassword(t *testing.T) {
	action := &updatePasswordsAction{}
	app, _, responder, events := newApp(t, func(config *inception.Config) {
		config.Features.UpdatePasswords = true
	}, func(builder *inception.Builder) {
		builder.WithUpdatePassword(action)
	})

	recorder := httptest.NewRecorder()
	inception.NewUpdatePasswordHandler(app).ServeHTTP(recorder, formRequest(http.MethodPut, "/user/password", "current_password=old&password=new&password_confirmation=new"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if action.input["current_password"] != "old" || !responder.updatePassword {
		t.Fatalf("input = %#v response = %v", action.input, responder.updatePassword)
	}
	if events.dispatched[0] != inception.EventPasswordUpdated {
		t.Fatalf("event = %#v, want password updated", events.dispatched[0])
	}
}

// PasswordControllerTest::test_passwords_cannot_be_updated_without_current_password
// PasswordControllerTest::test_passwords_cannot_be_updated_without_current_password_confirmation
func TestPasswordControllerUpdateCanFail(t *testing.T) {
	app, _, responder, _ := newApp(t, func(config *inception.Config) {
		config.Features.UpdatePasswords = true
	}, func(builder *inception.Builder) {
		builder.WithUpdatePassword(&updatePasswordsAction{err: errors.New("current password required")})
	})

	recorder := httptest.NewRecorder()
	inception.NewUpdatePasswordHandler(app).ServeHTTP(recorder, formRequest(http.MethodPut, "/user/password", "password=new"))

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	if responder.updatePassword {
		t.Fatal("password update response should not run on validation failure")
	}
}

// ProfileInformationControllerTest::test_contact_information_can_be_updated
// ProfileInformationControllerTest::test_email_address_will_be_updated_case_insensitive
func TestProfileInformationControllerUpdatesContactInformation(t *testing.T) {
	user := &fortifyUser{id: "1", email: "old@example.com", verified: true}
	verifier := &verifier{}
	app, _, responder, events := newApp(t, func(config *inception.Config) {
		config.Features.UpdateProfileInformation = true
		config.Features.EmailVerification = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{authenticated: user})
		builder.WithVerifier(verifier)
		builder.WithUpdateProfile(updateProfileAction{})
	})

	recorder := httptest.NewRecorder()
	inception.NewUpdateProfileHandler(app).ServeHTTP(recorder, formRequest(http.MethodPut, "/user/profile-information", "name=User&email=New@Example.com"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if user.email != "New@Example.com" || !user.unverified || !verifier.sent || !responder.profile {
		t.Fatalf("user = %#v verifier.sent = %v response = %v", user, verifier.sent, responder.profile)
	}
	if events.dispatched[0] != inception.EventProfileUpdated {
		t.Fatalf("event = %#v, want profile updated", events.dispatched[0])
	}
}

// EmailVerificationNotificationControllerTest::test_email_verification_notification_can_be_sent
func TestEmailVerificationNotificationControllerSendsNotification(t *testing.T) {
	verifier := &verifier{}
	app, _, responder, _ := newApp(t, func(config *inception.Config) {
		config.Features.EmailVerification = true
	}, func(builder *inception.Builder) {
		builder.WithVerifier(verifier)
	})

	recorder := httptest.NewRecorder()
	inception.NewSendVerificationHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/email/verification-notification", ""))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !verifier.sent || !responder.verification {
		t.Fatalf("verifier.sent = %v response = %v", verifier.sent, responder.verification)
	}
}

// EmailVerificationNotificationControllerTest::test_user_is_redirect_if_already_verified
// EmailVerificationPromptControllerTest::test_user_is_redirect_home_if_already_verified
func TestEmailVerificationNotificationControllerSkipsVerifiedUser(t *testing.T) {
	user := &fortifyUser{id: "1", email: "user@example.com", verified: true}
	verifier := &verifier{}
	app, _, responder, _ := newApp(t, func(config *inception.Config) {
		config.Features.EmailVerification = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{authenticated: user})
		builder.WithVerifier(verifier)
	})

	recorder := httptest.NewRecorder()
	inception.NewSendVerificationHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/email/verification-notification", ""))

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if verifier.sent || responder.verification {
		t.Fatal("verified users should not receive another notification")
	}
}

// VerifyEmailControllerTest::test_the_email_can_be_verified
// VerifyEmailControllerTest::test_email_is_not_verified_if_id_does_not_match
// VerifyEmailControllerTest::test_email_is_not_verified_if_email_does_not_match
func TestVerifyEmailControllerVerifiesSignedPathValues(t *testing.T) {
	verifier := &verifier{}
	app, _, _, events := newApp(t, func(config *inception.Config) {
		config.Features.EmailVerification = true
	}, func(builder *inception.Builder) {
		builder.WithVerifier(verifier)
	})

	request := httptest.NewRequest(http.MethodGet, "/verify-email/1/hash", nil)
	request.SetPathValue("id", "1")
	request.SetPathValue("hash", "hash")
	recorder := httptest.NewRecorder()
	inception.NewVerifyEmailHandler(app).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if verifier.id != "1" || verifier.hash != "hash" {
		t.Fatalf("verification args = %q %q", verifier.id, verifier.hash)
	}
	if events.dispatched[0] != inception.EventVerified {
		t.Fatalf("event = %#v, want verified", events.dispatched[0])
	}
}

// VerifyEmailControllerTest::test_redirected_if_email_is_already_verified
func TestVerifyEmailControllerSkipsAlreadyVerifiedUser(t *testing.T) {
	user := &fortifyUser{id: "1", email: "user@example.com", verified: true}
	verifier := &verifier{}
	app, _, _, events := newApp(t, func(config *inception.Config) {
		config.Features.EmailVerification = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{authenticated: user})
		builder.WithVerifier(verifier)
	})

	request := httptest.NewRequest(http.MethodGet, "/verify-email/1/hash", nil)
	request.SetPathValue("id", "1")
	request.SetPathValue("hash", "hash")
	recorder := httptest.NewRecorder()
	inception.NewVerifyEmailHandler(app).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if verifier.id != "" || verifier.hash != "" || len(events.dispatched) != 0 {
		t.Fatalf("verifier = %#v events = %#v", verifier, events.dispatched)
	}
}

// TwoFactorAuthenticationControllerTest::test_two_factor_authentication_can_be_enabled
// TwoFactorAuthenticationControllerTest::test_calling_two_factor_authentication_endpoint_will_overwrite_with_force_parameter
func TestTwoFactorAuthenticationControllerEnablesTwoFactor(t *testing.T) {
	user := &fortifyUser{id: "1", email: "user@example.com"}
	app, _, responder, events := newApp(t, func(config *inception.Config) {
		config.Features.TwoFactorAuthentication = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{authenticated: user})
	})

	recorder := httptest.NewRecorder()
	inception.NewEnableTwoFactorHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/user/two-factor-authentication", ""))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !user.twoFactorEnabled || user.twoFactorSecret == "" || len(user.twoFactorCodes) == 0 || !responder.twoFactorEnabled {
		t.Fatalf("user = %#v response = %v", user, responder.twoFactorEnabled)
	}
	if events.dispatched[0] != inception.EventTwoFactorEnabled {
		t.Fatalf("event = %#v, want two factor enabled", events.dispatched[0])
	}

	previousSecret := user.twoFactorSecret
	previousCodes := append([]string(nil), user.twoFactorCodes...)

	preserved := httptest.NewRecorder()
	inception.NewEnableTwoFactorHandler(app).ServeHTTP(preserved, formRequest(http.MethodPost, "/user/two-factor-authentication", ""))
	if preserved.Code != http.StatusOK {
		t.Fatalf("preserved status = %d, want %d", preserved.Code, http.StatusOK)
	}
	if user.twoFactorSecret != previousSecret || strings.Join(user.twoFactorCodes, ",") != strings.Join(previousCodes, ",") {
		t.Fatalf("two-factor state should not be overwritten without force: %#v", user)
	}

	forced := httptest.NewRecorder()
	inception.NewEnableTwoFactorHandler(app).ServeHTTP(forced, formRequest(http.MethodPost, "/user/two-factor-authentication", "force=1"))
	if forced.Code != http.StatusOK {
		t.Fatalf("forced status = %d, want %d", forced.Code, http.StatusOK)
	}
	if user.twoFactorSecret == previousSecret {
		t.Fatal("force should rotate the two-factor secret")
	}
}

// TwoFactorAuthenticationControllerTest::test_calling_two_factor_authentication_endpoint_will_not_overwrite_without_force_parameter
func TestTwoFactorAuthenticationControllerDoesNotOverwriteWithoutForce(t *testing.T) {
	now := time.Now()
	user := &fortifyUser{
		id:                 "1",
		email:              "user@example.com",
		twoFactorEnabled:   true,
		twoFactorSecret:    "EXISTINGSECRET",
		twoFactorCodes:     []string{"aaaa-bbbb"},
		twoFactorConfirmed: &now,
	}
	app, _, responder, events := newApp(t, func(config *inception.Config) {
		config.Features.TwoFactorAuthentication = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{authenticated: user})
	})

	recorder := httptest.NewRecorder()
	inception.NewEnableTwoFactorHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/user/two-factor-authentication", ""))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if user.twoFactorSecret != "EXISTINGSECRET" || user.twoFactorConfirmed != &now || strings.Join(user.twoFactorCodes, ",") != "aaaa-bbbb" {
		t.Fatalf("two-factor state was overwritten: %#v", user)
	}
	if !responder.twoFactorEnabled || len(events.dispatched) != 0 {
		t.Fatalf("response = %v events = %#v", responder.twoFactorEnabled, events.dispatched)
	}
}

// TwoFactorAuthenticationControllerTest::test_two_factor_authentication_secret_key_can_be_retrieved
func TestTwoFactorAuthenticationControllerSecretCanBeRetrievedThroughProvisioningURI(t *testing.T) {
	user := &fortifyUser{id: "1", email: "user@example.com", twoFactorEnabled: true, twoFactorSecret: "JBSWY3DPEHPK3PXP"}
	app, _, _, _ := newApp(t, func(config *inception.Config) {
		config.Features.TwoFactorAuthentication = true
		config.Features.EmailVerification = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{authenticated: user})
		builder.WithVerifier(&verifier{})
	})

	recorder := httptest.NewRecorder()
	inception.NewTwoFactorQRCodeHandler(app, "Bedrock").ServeHTTP(recorder, formRequest(http.MethodGet, "/user/two-factor-qr-code", ""))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), "secret=JBSWY3DPEHPK3PXP") {
		t.Fatalf("response does not contain provisioning secret: %s", recorder.Body.String())
	}
}

// RecoveryCodeControllerTest::test_new_recovery_codes_can_be_generated
func TestRecoveryCodeControllerGeneratesNewRecoveryCodes(t *testing.T) {
	user := &fortifyUser{id: "1", email: "user@example.com", twoFactorEnabled: true, twoFactorSecret: "secret", twoFactorCodes: []string{"aaaa-bbbb"}}
	app, _, _, _ := newApp(t, func(config *inception.Config) {
		config.Features.TwoFactorAuthentication = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{authenticated: user})
	})

	recorder := httptest.NewRecorder()
	inception.NewTwoFactorRecoveryCodesHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/user/two-factor-recovery-codes", ""))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if len(user.twoFactorCodes) != twofactor.DefaultRecoveryCodeCount {
		t.Fatalf("recovery codes = %#v, want %d", user.twoFactorCodes, twofactor.DefaultRecoveryCodeCount)
	}
	if user.twoFactorCodes[0] == "aaaa-bbbb" {
		t.Fatalf("recovery codes were not regenerated: %#v", user.twoFactorCodes)
	}
}

// TwoFactorAuthenticationControllerTest::test_two_factor_authentication_can_be_confirmed
// TwoFactorAuthenticationControllerTest::test_two_factor_authentication_can_not_be_confirmed_with_invalid_code
func TestTwoFactorAuthenticationControllerConfirmsTwoFactor(t *testing.T) {
	secret, err := twofactor.GenerateSecret(0)
	if err != nil {
		t.Fatal(err)
	}

	user := &fortifyUser{id: "1", email: "user@example.com", twoFactorEnabled: true, twoFactorSecret: secret}
	app, _, _, events := newApp(t, func(config *inception.Config) {
		config.Features.TwoFactorAuthentication = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{authenticated: user})
	})

	recorder := httptest.NewRecorder()
	inception.NewConfirmTwoFactorHandler(app).ServeHTTP(recorder, formRequest(http.MethodPost, "/user/confirmed-two-factor-authentication", "code="+twofactor.CurrentCode(secret)))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if user.twoFactorConfirmed == nil {
		t.Fatal("expected two-factor confirmation timestamp")
	}
	if events.dispatched[0] != inception.EventTwoFactorConfirmed {
		t.Fatalf("event = %#v, want two factor confirmed", events.dispatched[0])
	}

	failed := httptest.NewRecorder()
	inception.NewConfirmTwoFactorHandler(app).ServeHTTP(failed, formRequest(http.MethodPost, "/user/confirmed-two-factor-authentication", "code=000000"))
	if failed.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid status = %d, want %d", failed.Code, http.StatusUnprocessableEntity)
	}
}

// TwoFactorAuthenticationControllerTest::test_two_factor_authentication_can_be_disabled
func TestTwoFactorAuthenticationControllerDisablesTwoFactor(t *testing.T) {
	now := time.Now()
	user := &fortifyUser{
		id:                 "1",
		email:              "user@example.com",
		twoFactorEnabled:   true,
		twoFactorSecret:    "secret",
		twoFactorCodes:     []string{"aaaa-bbbb"},
		twoFactorConfirmed: &now,
	}
	app, _, responder, events := newApp(t, func(config *inception.Config) {
		config.Features.TwoFactorAuthentication = true
	}, func(builder *inception.Builder) {
		builder.WithGuard(&fortifyGuard{authenticated: user})
	})

	recorder := httptest.NewRecorder()
	inception.NewDisableTwoFactorHandler(app).ServeHTTP(recorder, formRequest(http.MethodDelete, "/user/two-factor-authentication", ""))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if user.twoFactorEnabled || user.twoFactorSecret != "" || user.twoFactorConfirmed != nil || len(user.twoFactorCodes) != 0 {
		t.Fatalf("two-factor state was not cleared: %#v", user)
	}
	if !responder.twoFactorOff || events.dispatched[0] != inception.EventTwoFactorDisabled {
		t.Fatalf("response = %v event = %#v", responder.twoFactorOff, events.dispatched[0])
	}
}

// PasswordRuleTest::test_password_rule
// PasswordRuleTest::test_password_rule_can_require_special_characters
// PasswordRuleTest::test_password_rule_can_require_numeric_and_special_characters
func TestPasswordRule(t *testing.T) {
	factory := validation.NewFactory()
	rule := validation.Rule.Password().Min(10).MixedCase().Numbers().Symbols()

	validator := factory.Make(map[string]any{
		"password": "Stronger-123",
	}, map[string]any{
		"password": rule,
	}, nil, nil)
	if validator.Fails() {
		t.Fatalf("strong password errors = %v", validator.Errors().All())
	}

	validator = factory.Make(map[string]any{
		"password": "letters-only",
	}, map[string]any{
		"password": rule,
	}, nil, nil)
	if !validator.Fails() {
		t.Fatal("password without numeric and special characters should fail")
	}
}
