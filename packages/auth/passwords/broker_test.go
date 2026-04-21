package passwords_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/auth"
	authevents "github.com/bedrock/packages/auth/events"
	"github.com/bedrock/packages/auth/passwords"
	cauth "github.com/bedrock/packages/contracts/auth"
	cevents "github.com/bedrock/packages/contracts/events"
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
	notifiedToken string
}

type brokerDispatcher struct {
	events []any
}

var _ cauth.PasswordResetNotificationSender = (*notifyingResetUser)(nil)

func (u *resetUser) GetEmailForPasswordReset() string { return u.email }

func (u *notifyingResetUser) SendPasswordResetNotification(token string) {
	u.notifiedToken = token
}

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

// Port of Illuminate\Tests\Auth\AuthPasswordBrokerTest::testBrokerCreatesTokenAndRedirectsWithoutError
func TestBrokerSendResetLinkCreatesTokenSendsNotificationAndDispatchesEvent(t *testing.T) {
	user := &notifyingResetUser{resetUser: &resetUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}}
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{"1": user}}
	repo := passwords.NewMemoryRepository(time.Hour)
	dispatcher := &brokerDispatcher{}
	broker := passwords.NewBroker(provider, repo, time.Hour).WithEventDispatcher(dispatcher)

	err := broker.SendResetLink(context.Background(), "test@example.com")

	if err != nil {
		t.Fatal(err)
	}

	if user.notifiedToken == "" {
		t.Fatal("SendResetLink should send the generated token to the user notification hook")
	}

	if !broker.TokenExists(context.Background(), user, user.notifiedToken) {
		t.Error("SendResetLink should store the generated token")
	}

	if len(dispatcher.events) != 1 {
		t.Fatalf("expected one dispatched event, got %d", len(dispatcher.events))
	}

	if _, ok := dispatcher.events[0].(authevents.PasswordResetLinkSent); !ok {
		t.Fatalf("expected PasswordResetLinkSent event, got %T", dispatcher.events[0])
	}
}

// Port of Illuminate\Tests\Auth\AuthPasswordBrokerTest::testExecutesCallbackInsteadOfSendingNotification
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
