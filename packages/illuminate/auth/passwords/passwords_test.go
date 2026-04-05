package passwords_test

import (
	"context"
	"errors"
	"testing"
	"time"

	auth "github.com/gollin/packages/illuminate/auth"
	"github.com/gollin/packages/illuminate/auth/passwords"
)

type passwordUser struct {
	id            string
	email         string
	passwordHash  string
	rememberToken string
}

func (u *passwordUser) GetAuthIdentifierName() string { return "id" }
func (u *passwordUser) GetAuthIdentifier() string     { return u.id }
func (u *passwordUser) GetAuthPasswordName() string   { return "password" }
func (u *passwordUser) GetAuthPassword() string       { return u.passwordHash }
func (u *passwordUser) SetAuthPassword(password string) {
	u.passwordHash = password
}

func (u *passwordUser) GetRememberToken() string { return u.rememberToken }
func (u *passwordUser) SetRememberToken(token string) {
	u.rememberToken = token
}

func (u *passwordUser) GetRememberTokenName() string { return "remember_token" }
func (u *passwordUser) GetEmailForPasswordReset() string {
	return u.email
}

type passwordProvider struct {
	user   *passwordUser
	hasher auth.PasswordHasher
}

func (p *passwordProvider) RetrieveByID(_ context.Context, id string) (auth.Authenticatable, error) {
	if p.user.id == id {
		return p.user, nil
	}

	return nil, auth.ErrUserNotFound
}

func (p *passwordProvider) RetrieveByToken(_ context.Context, id string, token string) (auth.Authenticatable, error) {
	if p.user.id == id && p.user.rememberToken == token {
		return p.user, nil
	}

	return nil, auth.ErrUnauthorized
}

func (p *passwordProvider) RetrieveByCredentials(_ context.Context, credentials map[string]string) (auth.Authenticatable, error) {
	if email, ok := credentials["email"]; ok && passwords.NormalizeEmail(email) == passwords.NormalizeEmail(p.user.email) {
		return p.user, nil
	}

	return nil, auth.ErrUserNotFound
}

func (p *passwordProvider) UpdateRememberToken(_ context.Context, user auth.Authenticatable, token string) error {
	user.SetRememberToken(token)
	return nil
}

func (p *passwordProvider) ValidateCredentials(ctx context.Context, user auth.Authenticatable, credentials map[string]string) (bool, error) {
	if err := p.hasher.Compare(ctx, user.GetAuthPassword(), credentials["password"]); err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (p *passwordProvider) RehashPasswordIfRequired(context.Context, auth.Authenticatable, map[string]string, bool) error {
	return nil
}

type resetter struct {
	hasher auth.PasswordHasher
}

func (r resetter) Reset(ctx context.Context, user auth.Authenticatable, password string) error {
	hash, err := r.hasher.Hash(ctx, password)
	if err != nil {
		return err
	}

	user.SetAuthPassword(hash)
	return nil
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

func TestBrokerCreateValidateAndReset(t *testing.T) {
	t.Parallel()

	hasher, err := auth.NewDefaultPasswordHasher()
	if err != nil {
		t.Fatalf("NewDefaultPasswordHasher: %v", err)
	}

	passwordHash, err := hasher.Hash(context.Background(), "secret")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	user := &passwordUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	broker := &passwords.Broker{
		Config: passwords.Config{
			Name:     "users",
			Expire:   time.Hour,
			Throttle: time.Minute,
		},
		Users:  &passwordProvider{user: user, hasher: hasher},
		Tokens: passwords.NewMemoryTokenRepository(),
		Clock:  clock,
	}

	token, err := broker.CreateToken(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}

	valid, err := broker.TokenExists(context.Background(), user, token)
	if err != nil {
		t.Fatalf("TokenExists: %v", err)
	}

	if !valid {
		t.Fatal("expected token to be valid")
	}

	canCreate, err := broker.CanCreateToken(context.Background(), user)
	if err != nil {
		t.Fatalf("CanCreateToken: %v", err)
	}

	if canCreate {
		t.Fatal("expected throttle window to block immediate reissue")
	}

	if _, err := broker.Reset(context.Background(), user.email, token, "new-secret", resetter{hasher: hasher}); err != nil {
		t.Fatalf("Reset: %v", err)
	}

	if err := hasher.Compare(context.Background(), user.GetAuthPassword(), "new-secret"); err != nil {
		t.Fatalf("expected updated password hash: %v", err)
	}
}
