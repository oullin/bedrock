package actions_test

import (
	"context"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/fortify/actions"
	"github.com/gollin/packages/auth/fortify/contracts"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/memory"
)

func TestCreateUserNormalizesEmailAndValidates(t *testing.T) {
	t.Parallel()

	users := memory.NewInMemoryUserRepository()
	action := actions.CreateUser{
		Users:  users,
		Hasher: auth.DefaultPasswordHasher{},
		IDs:    memory.NewSequenceIDGenerator("user"),
		Clock:  memory.NewFixedClock(time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)),
	}

	user, err := action.Create(context.Background(), contracts.RegisterInput{
		Name:                 "Test User",
		Email:                " TEST@Example.com ",
		Password:             "password-123",
		PasswordConfirmation: "password-123",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	profile := user.(auth.UserProfile)
	if profile.GetEmail() != "test@example.com" {
		t.Fatalf("unexpected normalized email: %s", profile.GetEmail())
	}

	if _, err := action.Create(context.Background(), contracts.RegisterInput{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestResetAndUpdatePassword(t *testing.T) {
	t.Parallel()

	hasher := auth.DefaultPasswordHasher{}
	hash, err := hasher.Hash(context.Background(), "password-123")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	users := memory.NewInMemoryUserRepository()
	user := &foundation.User{
		ID:            "user-1",
		Name:          "User",
		Email:         "user@example.com",
		PasswordHash:  hash,
		RememberToken: "remember",
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	reset := actions.ResetUserPassword{Users: users, Hasher: hasher}
	if err := reset.Reset(context.Background(), user, "new-password-123"); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if user.GetRememberToken() != "" {
		t.Fatal("expected remember token to be cleared")
	}

	update := actions.UpdateUserPassword{Users: users, Hasher: hasher}
	if err := update.Update(context.Background(), user, contracts.UpdatePasswordInput{
		CurrentPassword:      "new-password-123",
		Password:             "another-password-123",
		PasswordConfirmation: "another-password-123",
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	if err := update.Update(context.Background(), user, contracts.UpdatePasswordInput{
		CurrentPassword:      "wrong",
		Password:             "another-password-123",
		PasswordConfirmation: "another-password-123",
	}); err != auth.ErrInvalidCredentials {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestUpdateProfileInformationMarksEmailUnverified(t *testing.T) {
	t.Parallel()

	users := memory.NewInMemoryUserRepository()
	now := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)
	user := &foundation.User{
		ID:        "user-1",
		Name:      "User",
		Email:     "user@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}
	user.MarkEmailAsVerified(now)
	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	action := actions.UpdateUserProfileInformation{Users: users}
	if err := action.Update(context.Background(), user, contracts.UpdateProfileInformationInput{
		Name:  "Updated User",
		Email: "UPDATED@example.com",
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if user.GetName() != "Updated User" {
		t.Fatalf("unexpected name: %s", user.GetName())
	}
	if user.GetEmail() != "updated@example.com" {
		t.Fatalf("unexpected email: %s", user.GetEmail())
	}
	if user.HasVerifiedEmail() {
		t.Fatal("expected email verification to be cleared")
	}
}
