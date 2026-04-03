package memory

import (
	"context"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/passwords"
)

func TestInMemoryRepositoriesAndHelpers(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)
	clock := NewFixedClock(now)
	clock.Advance(time.Minute)
	if clock.Now() != now.Add(time.Minute) {
		t.Fatalf("unexpected clock time: %v", clock.Now())
	}

	mailer := &InMemoryMailer{}
	if err := mailer.Send(context.Background(), auth.MailMessage{To: "user@example.com"}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if len(mailer.Messages()) != 1 {
		t.Fatalf("unexpected messages: %d", len(mailer.Messages()))
	}

	ids := NewSequenceIDGenerator("user")
	if ids.NewID() != "user-1" || ids.NewID() != "user-2" {
		t.Fatal("unexpected sequence ids")
	}

	users := NewInMemoryUserRepository()
	user := &foundation.User{
		ID:           "user-1",
		Name:         "User One",
		Email:        "USER@example.com",
		PasswordHash: "hash",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := users.Create(context.Background(), user); err != auth.ErrUserExists {
		t.Fatalf("expected ErrUserExists, got %v", err)
	}

	byID, err := users.RetrieveByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("RetrieveByID: %v", err)
	}
	if byID.GetAuthIdentifier() != user.ID {
		t.Fatalf("unexpected user id: %s", byID.GetAuthIdentifier())
	}

	byCreds, err := users.RetrieveByCredentials(context.Background(), map[string]string{"email": " user@example.com "})
	if err != nil {
		t.Fatalf("RetrieveByCredentials: %v", err)
	}
	if byCreds.GetAuthIdentifier() != user.ID {
		t.Fatalf("unexpected credential user: %s", byCreds.GetAuthIdentifier())
	}

	if err := users.UpdateRememberToken(context.Background(), byCreds, "remember-token"); err != nil {
		t.Fatalf("UpdateRememberToken: %v", err)
	}
	if _, err := users.RetrieveByToken(context.Background(), user.ID, "remember-token"); err != nil {
		t.Fatalf("RetrieveByToken: %v", err)
	}

	sessionStore := NewInMemorySessionStore()
	session := &auth.Session{ID: "session-1", UserID: user.ID, ExpiresAt: now.Add(time.Hour), LastSeenAt: now, CreatedAt: now}
	if err := sessionStore.Create(context.Background(), session); err != nil {
		t.Fatalf("Create session: %v", err)
	}
	storedSession, err := sessionStore.FindByID(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	storedSession.PendingTwoFactor = true
	if err := sessionStore.Update(context.Background(), storedSession); err != nil {
		t.Fatalf("Update session: %v", err)
	}
	if err := sessionStore.Delete(context.Background(), session.ID); err != nil {
		t.Fatalf("Delete session: %v", err)
	}

	tokens := NewInMemoryTokenRepository()
	token := &passwords.Token{
		UserID:    user.ID,
		TokenHash: "token-hash",
		CreatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}
	if err := tokens.Save(context.Background(), token); err != nil {
		t.Fatalf("Save token: %v", err)
	}
	if _, err := tokens.FindByTokenHash(context.Background(), token.TokenHash); err != nil {
		t.Fatalf("FindByTokenHash: %v", err)
	}
	recent, err := tokens.RecentlyCreated(context.Background(), user.ID, now.Add(-time.Minute))
	if err != nil {
		t.Fatalf("RecentlyCreated: %v", err)
	}
	if !recent {
		t.Fatal("expected token to be recently created")
	}
	if err := tokens.DeleteByTokenHash(context.Background(), token.TokenHash); err != nil {
		t.Fatalf("DeleteByTokenHash: %v", err)
	}
	if _, err := tokens.FindByTokenHash(context.Background(), token.TokenHash); err != auth.ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}
