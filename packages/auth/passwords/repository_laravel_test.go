package passwords_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/auth"
	"github.com/bedrock/packages/auth/passwords"
	cauth "github.com/bedrock/packages/contracts/auth"
)

// Port of Framework\Tests\Auth\AuthDatabaseTokenRepositoryTest::testCreateInsertsNewRecordIntoTable
// Port of Framework\Tests\Auth\AuthDatabaseTokenRepositoryTest::testExistReturnsTrueIfValidRecordExists
// Port of Framework\Tests\Auth\AuthDatabaseTokenRepositoryTest::testExistReturnsFalseIfInvalidToken
func TestUpstreamAuthDatabaseTokenRepositoryCreatesAndValidatesTokens(t *testing.T) {
	repo := passwords.NewMemoryRepository(time.Hour)
	ctx := context.Background()

	token, err := repo.Create(ctx, "taylor@example.com")

	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("Create should return a token")
	}

	if !repo.Exists(ctx, "taylor@example.com", token) {
		t.Fatal("Exists should accept a valid token")
	}

	if repo.Exists(ctx, "taylor@example.com", "invalid") {
		t.Fatal("Exists should reject an invalid token")
	}
}

// Port of Framework\Tests\Auth\AuthDatabaseTokenRepositoryTest::testExistReturnsFalseIfNoRowFoundForUser
// Port of Framework\Tests\Auth\AuthDatabaseTokenRepositoryTest::testExistReturnsFalseIfRecordIsExpired
func TestUpstreamAuthDatabaseTokenRepositoryRejectsMissingAndExpiredTokens(t *testing.T) {
	ctx := context.Background()
	repo := passwords.NewMemoryRepository(0)

	token, err := repo.Create(ctx, "taylor@example.com")

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond)

	if repo.Exists(ctx, "nobody@example.com", token) {
		t.Fatal("Exists should reject a missing user")
	}

	if repo.Exists(ctx, "taylor@example.com", token) {
		t.Fatal("Exists should reject an expired token")
	}
}

// Port of Framework\Tests\Auth\AuthDatabaseTokenRepositoryTest::testRecentlyCreatedReturnsFalseIfNoRowFoundForUser
// Port of Framework\Tests\Auth\AuthDatabaseTokenRepositoryTest::testRecentlyCreatedReturnsTrueIfRecordIsRecentlyCreated
// Port of Framework\Tests\Auth\AuthDatabaseTokenRepositoryTest::testRecentlyCreatedReturnsFalseIfValidRecordExists
func TestUpstreamAuthDatabaseTokenRepositoryRecentlyCreated(t *testing.T) {
	ctx := context.Background()
	repo := passwords.NewMemoryRepository(time.Hour)

	if repo.RecentlyCreated(ctx, "missing@example.com", time.Minute) {
		t.Fatal("RecentlyCreated should reject missing records")
	}

	if _, err := repo.Create(ctx, "taylor@example.com"); err != nil {
		t.Fatal(err)
	}

	if !repo.RecentlyCreated(ctx, "taylor@example.com", time.Minute) {
		t.Fatal("RecentlyCreated should accept fresh records")
	}

	if repo.RecentlyCreated(ctx, "taylor@example.com", 0) {
		t.Fatal("RecentlyCreated should reject records outside the throttle window")
	}
}

// Port of Framework\Tests\Auth\AuthDatabaseTokenRepositoryTest::testDeleteMethodDeletesByToken
// Port of Framework\Tests\Auth\AuthDatabaseTokenRepositoryTest::testDeleteExpiredMethodDeletesExpiredTokens
func TestUpstreamAuthDatabaseTokenRepositoryDeleteMethods(t *testing.T) {
	ctx := context.Background()
	repo := passwords.NewMemoryRepository(0)

	token, err := repo.Create(ctx, "taylor@example.com")

	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Delete(ctx, "taylor@example.com"); err != nil {
		t.Fatal(err)
	}

	if repo.Exists(ctx, "taylor@example.com", token) {
		t.Fatal("Delete should remove the token")
	}

	token, err = repo.Create(ctx, "taylor@example.com")

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond)

	if err := repo.DeleteExpired(ctx); err != nil {
		t.Fatal(err)
	}

	if repo.Exists(ctx, "taylor@example.com", token) {
		t.Fatal("DeleteExpired should remove expired tokens")
	}
}

// Port of Framework\Tests\Auth\AuthPasswordBrokerTest::testIfTokenIsRecentlyCreated
func TestUpstreamAuthPasswordBrokerRejectsRecentlyCreatedToken(t *testing.T) {
	user := &resetUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{"1": user}}
	repo := passwords.NewMemoryRepository(time.Hour)
	broker := passwords.NewBroker(provider, repo, time.Hour).WithThrottle(time.Minute)
	ctx := context.Background()

	if err := broker.SendResetLink(ctx, "test@example.com"); err != nil {
		t.Fatal(err)
	}

	if err := broker.SendResetLink(ctx, "test@example.com"); !errors.Is(err, passwords.ErrResetLinkThrottled) {
		t.Fatalf("SendResetLink error = %v, want ErrResetLinkThrottled", err)
	}
}

// Port of Framework\Tests\Auth\AuthPasswordBrokerTest::testResetRemovesRecordOnReminderTableAndCallsCallback
func TestUpstreamAuthPasswordBrokerResetDeletesTokenAndCallsCallback(t *testing.T) {
	user := &resetUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{"1": user}}
	repo := passwords.NewMemoryRepository(time.Hour)
	broker := passwords.NewBroker(provider, repo, time.Hour)
	ctx := context.Background()

	token, err := broker.CreateToken(ctx, user)

	if err != nil {
		t.Fatal(err)
	}

	called := false
	err = broker.Reset(ctx, map[string]any{
		"email":    "test@example.com",
		"token":    token,
		"password": "new-secret",
	}, func(_ context.Context, got cauth.CanResetPassword, gotToken, password string) error {
		called = true

		if got != user {
			t.Error("reset callback received the wrong user")
		}

		if gotToken != token {
			t.Errorf("token = %q, want %q", gotToken, token)
		}

		if password != "new-secret" {
			t.Errorf("password = %q, want %q", password, "new-secret")
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Fatal("reset callback was not called")
	}

	if broker.TokenExists(ctx, user, token) {
		t.Fatal("Reset should delete the used token")
	}
}

// Port of Framework\Tests\Auth\AuthPasswordBrokerTest::testRedirectIsReturnedByResetWhenUserCredentialsInvalid
func TestUpstreamAuthPasswordBrokerResetRejectsInvalidToken(t *testing.T) {
	user := &resetUser{
		GenericUser: auth.NewGenericUser(map[string]any{"id": "1"}),
		email:       "test@example.com",
	}
	provider := &brokerProvider{users: map[string]cauth.Authenticatable{"1": user}}
	repo := passwords.NewMemoryRepository(time.Hour)
	broker := passwords.NewBroker(provider, repo, time.Hour)

	err := broker.Reset(context.Background(), map[string]any{
		"email": "test@example.com",
		"token": "invalid",
	}, func(context.Context, cauth.CanResetPassword, string, string) error {
		return errors.New("callback should not run")
	})

	if err == nil {
		t.Fatal("Reset should reject invalid credentials")
	}
}
