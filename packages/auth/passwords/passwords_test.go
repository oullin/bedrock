package passwords_test

import (
	"context"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/fortify/actions"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/memory"
	"github.com/gollin/packages/auth/passwords"
)

func TestBrokerSendResetLinkAndReset(t *testing.T) {
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

	token, err := broker.SendResetLink(context.Background(), user.Email)

	if err != nil {
		t.Fatalf("SendResetLink: %v", err)
	}

	if token == "" {
		t.Fatal("expected reset token")
	}

	if len(mailer.Messages()) != 1 {
		t.Fatalf("expected one mail, got %d", len(mailer.Messages()))
	}

	updated, err := broker.Reset(context.Background(), user.Email, token, "new-password-123", actions.ResetUserPassword{
		Users:  users,
		Hasher: hasher,
		Clock:  clock,
	})

	if err != nil {
		t.Fatalf("Reset: %v", err)
	}

	if err := hasher.Compare(context.Background(), updated.GetAuthPassword(), "new-password-123"); err != nil {
		t.Fatalf("Compare new password: %v", err)
	}
}
