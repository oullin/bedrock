package passwords_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/auth"
	"github.com/bedrock/packages/auth/passwords"
	cauth "github.com/bedrock/packages/contracts/auth"
)

type resetUser struct {
	*auth.GenericUser
	email string
}

type brokerProvider struct {
	users map[string]cauth.Authenticatable
}

func (u *resetUser) GetEmailForPasswordReset() string { return u.email }

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
		if ru, ok := u.(*resetUser); ok && ru.email == email {
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
