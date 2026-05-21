package passwords_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/auth"
	authevents "github.com/bedrock/packages/auth/events"
	"github.com/bedrock/packages/auth/passwords"
	cauth "github.com/bedrock/packages/contracts/auth"
	cevents "github.com/bedrock/packages/contracts/events"
	clog "github.com/bedrock/packages/contracts/log"
)

type resetUser struct {
	*auth.GenericUser
	email string
}

type brokerProvider struct {
	users map[string]cauth.Authenticatable
}

type notifyingResetUser struct {
	*resetUser
	notifiedToken   string
	notifiedContext context.Context
}

type brokerDispatcher struct {
	events []any
}

type incompatibleTokenRepository struct {
	createCalls int
}

type brokerLogEntry struct {
	message string
	context []map[string]any
}

type brokerLogger struct {
	warnings []brokerLogEntry
}

type resetContextKey struct{}

var _ cauth.PasswordResetNotificationSender = (*notifyingResetUser)(nil)

func (u *resetUser) GetEmailForPasswordReset() string { return u.email }

func (u *notifyingResetUser) SendPasswordResetNotification(ctx context.Context, token string) {
	u.notifiedToken = token
	u.notifiedContext = ctx
}

func (r *incompatibleTokenRepository) Create(_ context.Context, _ string) (string, error) {
	r.createCalls++

	return "token", nil
}

func (r *incompatibleTokenRepository) Exists(_ context.Context, _, _ string) bool {
	return false
}

func (r *incompatibleTokenRepository) Delete(_ context.Context, _ string) error {
	return nil
}

func (r *incompatibleTokenRepository) DeleteExpired(_ context.Context) error {
	return nil
}

func (l *brokerLogger) Emergency(_ string, _ ...map[string]any) {}
func (l *brokerLogger) Alert(_ string, _ ...map[string]any)     {}
func (l *brokerLogger) Critical(_ string, _ ...map[string]any)  {}
func (l *brokerLogger) Error(_ string, _ ...map[string]any)     {}
func (l *brokerLogger) Notice(_ string, _ ...map[string]any)    {}
func (l *brokerLogger) Info(_ string, _ ...map[string]any)      {}
func (l *brokerLogger) Debug(_ string, _ ...map[string]any)     {}

func (l *brokerLogger) Warning(message string, context ...map[string]any) {
	l.warnings = append(l.warnings, brokerLogEntry{
		message: message,
		context: context,
	})
}

func (l *brokerLogger) Log(_ clog.Level, _ string, _ ...map[string]any) {}

func (p *brokerProvider) RetrieveByID(_ context.Context, id string) (cauth.Authenticatable, error) {
	return p.users[id], nil
}

func (p *brokerProvider) RetrieveByToken(_ context.Context, _ string, _ string) (cauth.Authenticatable, error) {
	return nil, nil
}

func (p *brokerProvider) UpdateRememberToken(_ context.Context, _ cauth.Authenticatable, _ string) error {
	return nil
}

func (p *brokerProvider) RetrieveByCredentials(_ context.Context, creds map[string]string) (cauth.Authenticatable, error) {
	email := creds["email"]

	for _, u := range p.users {
		if ru, ok := u.(cauth.CanResetPassword); ok && ru.GetEmailForPasswordReset() == email {
			return u, nil
		}
	}

	return nil, nil
}

func (p *brokerProvider) ValidateCredentials(_ context.Context, _ cauth.Authenticatable, _ map[string]string) (bool, error) {
	return true, nil
}

func (p *brokerProvider) RehashPasswordIfRequired(_ context.Context, _ cauth.Authenticatable, _ map[string]string, _ bool) error {
	return nil
}

func (d *brokerDispatcher) Listen(_ any, _ ...cevents.Listener)         {}
func (d *brokerDispatcher) HasListeners(_ any) bool                     { return false }
func (d *brokerDispatcher) HasWildcardListeners(_ any) bool             { return false }
func (d *brokerDispatcher) Subscribe(_ cevents.Subscriber)              {}
func (d *brokerDispatcher) Until(_ context.Context, _ any) (any, error) { return nil, nil }
func (d *brokerDispatcher) Dispatch(_ context.Context, event any) ([]any, error) {
	d.events = append(d.events, event)

	return nil, nil
}
func (d *brokerDispatcher) Push(_ context.Context, _ any)           {}
func (d *brokerDispatcher) Flush(_ context.Context, _ string) error { return nil }
func (d *brokerDispatcher) Forget(_ any)                            {}
func (d *brokerDispatcher) ForgetPushed()                           {}
func (d *brokerDispatcher) GetListeners(_ any) []cevents.Listener   { return nil }

func TestBrokerGetUser(t *testing.T) {
	user := &resetUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{"1": user}}
	repo := passwords.NewMemoryRepository(time.Hour)
	broker := passwords.NewBroker(provider, repo, time.Hour)

	got, err := broker.GetUser(context.Background(), "test@example.com")

	if err != nil {
		t.Fatal(err)
	}

	if got != user {
		t.Error("GetUser should return the matching user")
	}
}

func TestBrokerGetUserReturnsErrorForMissing(t *testing.T) {
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{}}
	repo := passwords.NewMemoryRepository(time.Hour)
	broker := passwords.NewBroker(provider, repo, time.Hour)

	_, err := broker.GetUser(context.Background(), "missing@example.com")

	if err == nil {
		t.Error("GetUser should return error for missing user")
	}
}

func TestBrokerCreateAndTokenExists(t *testing.T) {
	user := &resetUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{"1": user}}
	repo := passwords.NewMemoryRepository(time.Hour)
	broker := passwords.NewBroker(provider, repo, time.Hour)

	token, err := broker.CreateToken(context.Background(), user)

	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Error("CreateToken should return a non-empty token")
	}

	if !broker.TokenExists(context.Background(), user, token) {
		t.Error("TokenExists should return true for valid token")
	}

	if broker.TokenExists(context.Background(), user, "invalid") {
		t.Error("TokenExists should return false for invalid token")
	}
}

// Ref: @bedrock/code-0358
func TestBrokerSendResetLinkCreatesTokenSendsNotificationAndDispatchesEvent(t *testing.T) {
	user := &notifyingResetUser{resetUser: &resetUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}}
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{"1": user}}
	repo := passwords.NewMemoryRepository(time.Hour)
	dispatcher := &brokerDispatcher{}
	broker := passwords.NewBroker(provider, repo, time.Hour).WithEventDispatcher(dispatcher)
	ctx := context.WithValue(context.Background(), resetContextKey{}, "reset-request")

	err := broker.SendResetLink(ctx, "test@example.com")

	if err != nil {
		t.Fatal(err)
	}

	if user.notifiedToken == "" {
		t.Fatal("SendResetLink should send the generated token to the user notification hook")
	}

	if !broker.TokenExists(context.Background(), user, user.notifiedToken) {
		t.Error("SendResetLink should store the generated token")
	}

	if got := user.notifiedContext.Value(resetContextKey{}); got != "reset-request" {
		t.Errorf("SendResetLink should pass context to notification hook, got context value %v", got)
	}

	if len(dispatcher.events) != 1 {
		t.Fatalf("expected one dispatched event, got %d", len(dispatcher.events))
	}

	event, ok := dispatcher.events[0].(authevents.PasswordResetLinkSent)

	if !ok {
		t.Fatalf("expected PasswordResetLinkSent event, got %T", dispatcher.events[0])
	}

	if event.User.GetAuthIdentifier() != "1" {
		t.Errorf("event user auth identifier = %q, want %q", event.User.GetAuthIdentifier(), "1")
	}

	if event.User.GetEmailForPasswordReset() != "test@example.com" {
		t.Errorf("event user reset email = %q, want %q", event.User.GetEmailForPasswordReset(), "test@example.com")
	}
}

// Ref: @bedrock/code-0358
func TestBrokerSendResetLinkUsingExecutesCallbackInsteadOfNotification(t *testing.T) {
	user := &notifyingResetUser{resetUser: &resetUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}}
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{"1": user}}
	repo := passwords.NewMemoryRepository(time.Hour)
	dispatcher := &brokerDispatcher{}
	broker := passwords.NewBroker(provider, repo, time.Hour).WithEventDispatcher(dispatcher)

	var callbackToken string

	err := broker.SendResetLinkUsing(context.Background(), "test@example.com", func(_ context.Context, got cauth.CanResetPassword, token string) error {
		if got != user {
			t.Fatalf("callback user = %p, want %p", got, user)
		}

		callbackToken = token

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if callbackToken == "" {
		t.Fatal("SendResetLinkUsing should pass the generated token to the callback")
	}

	if user.notifiedToken != "" {
		t.Error("SendResetLinkUsing should not call the default notification hook when a callback is supplied")
	}

	if len(dispatcher.events) != 0 {
		t.Errorf("SendResetLinkUsing should not dispatch default reset-link event with callback, got %d events", len(dispatcher.events))
	}
}

func TestBrokerSendResetLinkWithThrottleRequiresRecentTokenRepository(t *testing.T) {
	user := &notifyingResetUser{resetUser: &resetUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}}
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{"1": user}}
	repo := &incompatibleTokenRepository{}
	dispatcher := &brokerDispatcher{}
	logger := &brokerLogger{}
	broker := passwords.NewBroker(provider, repo, time.Hour).
		WithThrottle(time.Minute).
		WithEventDispatcher(dispatcher).
		WithLogger(logger)

	callbackCalled := false

	err := broker.SendResetLinkUsing(context.Background(), "test@example.com", func(context.Context, cauth.CanResetPassword, string) error {
		callbackCalled = true

		return nil
	})

	if !errors.Is(err, passwords.ErrThrottleRepositoryUnsupported) {
		t.Fatalf("SendResetLinkUsing error = %v, want ErrThrottleRepositoryUnsupported", err)
	}

	if repo.createCalls != 0 {
		t.Fatalf("repository Create calls = %d, want 0", repo.createCalls)
	}

	if callbackCalled {
		t.Fatal("SendResetLinkUsing should not call the callback when the repository cannot support throttling")
	}

	if user.notifiedToken != "" {
		t.Fatal("SendResetLinkUsing should not send a notification when the repository cannot support throttling")
	}

	if len(dispatcher.events) != 0 {
		t.Fatalf("expected no dispatched events, got %d", len(dispatcher.events))
	}

	if len(logger.warnings) != 1 {
		t.Fatalf("expected one warning, got %d", len(logger.warnings))
	}

	warning := logger.warnings[0]

	if warning.message == "" {
		t.Fatal("warning message should not be empty")
	}

	if len(warning.context) != 1 {
		t.Fatalf("warning context count = %d, want 1", len(warning.context))
	}

	if warning.context[0]["repository"] == "" {
		t.Fatal("warning context should include the repository type")
	}

	if warning.context[0]["throttle"] != time.Minute.String() {
		t.Fatalf("warning throttle context = %v, want %q", warning.context[0]["throttle"], time.Minute.String())
	}

	if _, ok := warning.context[0]["email"]; ok {
		t.Fatal("warning context should not include email")
	}

	if _, ok := warning.context[0]["token"]; ok {
		t.Fatal("warning context should not include token")
	}
}

func TestBrokerDeleteToken(t *testing.T) {
	user := &resetUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{"1": user}}
	repo := passwords.NewMemoryRepository(time.Hour)
	broker := passwords.NewBroker(provider, repo, time.Hour)

	token, _ := broker.CreateToken(context.Background(), user)

	err := broker.DeleteToken(context.Background(), user)

	if err != nil {
		t.Fatal(err)
	}

	if broker.TokenExists(context.Background(), user, token) {
		t.Error("TokenExists should return false after DeleteToken")
	}
}

func TestBrokerGetRepository(t *testing.T) {
	repo := passwords.NewMemoryRepository(time.Hour)
	broker := passwords.NewBroker(nil, repo, time.Hour)

	if broker.GetRepository() != repo {
		t.Error("GetRepository should return the token repository")
	}
}
