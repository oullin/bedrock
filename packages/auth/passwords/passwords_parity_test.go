package passwords_test

import (
	"context"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/authflows/actions"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/memory"
	"github.com/gollin/packages/auth/passwords"
)

func TestBrokerRejectsInvalidRecentlyCreatedAndExpiredTokens(t *testing.T) {
	t.Parallel()

	clock := memory.NewFixedClock(time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC))
	users := memory.NewInMemoryUserRepository()
	tokens := memory.NewInMemoryTokenRepository()
	mailer := &memory.InMemoryMailer{}
	hasher := auth.DefaultPasswordHasher{}

	hash, err := hasher.Hash(context.Background(), "password-123")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	user := &foundation.User{
		ID:           "user-1",
		Name:         "Reset User",
		Email:        "reset@example.com",
		PasswordHash: hash,
		CreatedAt:    clock.Now(),
		UpdatedAt:    clock.Now(),
	}
	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	broker := &passwords.Broker{
		Config: passwords.Config{Name: "users", Expire: time.Hour, Throttle: time.Minute},
		Users:  users,
		Tokens: tokens,
		Mailer: mailer,
		Clock:  clock,
	}

	if _, err := broker.SendResetLink(context.Background(), "missing@example.com"); err != auth.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	token, err := broker.SendResetLink(context.Background(), user.Email)
	if err != nil {
		t.Fatalf("SendResetLink: %v", err)
	}
	if _, err := broker.SendResetLink(context.Background(), user.Email); err == nil {
		t.Fatal("expected throttle error")
	}

	if _, err := broker.Reset(context.Background(), user.Email, "bad-token", "new-password-123", actions.ResetUserPassword{
		Users:  users,
		Hasher: hasher,
		Clock:  clock,
	}); err != auth.ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}

	clock.Advance(2 * time.Hour)
	if _, err := broker.Reset(context.Background(), user.Email, token, "new-password-123", actions.ResetUserPassword{
		Users:  users,
		Hasher: hasher,
		Clock:  clock,
	}); err != auth.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}
