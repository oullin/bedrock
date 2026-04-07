package passwords_test

import (
	"context"
	"errors"
	"testing"
	"time"

	auth "github.com/bedrock/packages/auth"
	"github.com/bedrock/packages/auth/passwords"
)

// ---------- test helpers ----------

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

func (u *passwordUser) GetRememberTokenName() string    { return "remember_token" }
func (u *passwordUser) GetEmailForPasswordReset() string { return u.email }

type passwordProvider struct {
	user   *passwordUser
	hasher auth.PasswordHasher
}

func (p *passwordProvider) RetrieveByID(_ context.Context, id string) (auth.Authenticatable, error) {
	if p.user != nil && p.user.id == id {
		return p.user, nil
	}
	return nil, auth.ErrUserNotFound
}

func (p *passwordProvider) RetrieveByToken(_ context.Context, id string, token string) (auth.Authenticatable, error) {
	if p.user != nil && p.user.id == id && p.user.rememberToken == token {
		return p.user, nil
	}
	return nil, auth.ErrUnauthorized
}

func (p *passwordProvider) RetrieveByCredentials(_ context.Context, credentials map[string]string) (auth.Authenticatable, error) {
	if p.user == nil {
		return nil, auth.ErrUserNotFound
	}
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

func newTestBroker(user *passwordUser) (*passwords.Broker, auth.PasswordHasher) {
	hasher, _ := auth.NewDefaultPasswordHasher()
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}

	provider := &passwordProvider{user: user, hasher: hasher}
	if user == nil {
		provider = &passwordProvider{hasher: hasher}
	}

	return &passwords.Broker{
		Config: passwords.Config{
			Name:     "users",
			Expire:   time.Hour,
			Throttle: time.Minute,
		},
		Users:  provider,
		Tokens: passwords.NewMemoryTokenRepository(),
		Clock:  clock,
	}, hasher
}

// ======================== BROKER TESTS ========================

// Upstream: testBrokerCreatesTokenAndRedirectsWithoutError
func TestBrokerCreateTokenSuccess(t *testing.T) {
	t.Parallel()

	user := &passwordUser{id: "user-1", email: "user@example.com"}
	broker, _ := newTestBroker(user)

	token, err := broker.CreateToken(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

// Upstream: testUserIsRetrievedByCredentials + testBrokerCreatesTokenAndRedirectsWithoutError
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
	broker, _ := newTestBroker(user)

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

// Upstream: testIfUserIsNotFoundErrorRedirectIsReturned
func TestBrokerResetReturnsErrorWhenUserNotFound(t *testing.T) {
	t.Parallel()

	broker, _ := newTestBroker(nil)

	_, err := broker.Reset(context.Background(), "missing@example.com", "any-token", "password", resetter{})
	if err == nil {
		t.Fatal("expected error for missing user")
	}
}

// Upstream: testIfTokenIsRecentlyCreated
func TestBrokerThrottlesTokenCreation(t *testing.T) {
	t.Parallel()

	user := &passwordUser{id: "user-1", email: "user@example.com"}
	broker, _ := newTestBroker(user)

	// Create first token
	_, err := broker.CreateToken(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}

	// Immediately try again — should be throttled
	canCreate, err := broker.CanCreateToken(context.Background(), user)
	if err != nil {
		t.Fatalf("CanCreateToken: %v", err)
	}
	if canCreate {
		t.Fatal("expected token creation to be throttled")
	}
}

// Upstream: testRedirectIsReturnedByResetWhenUserCredentialsInvalid
func TestBrokerResetReturnsErrorForInvalidToken(t *testing.T) {
	t.Parallel()

	user := &passwordUser{id: "user-1", email: "user@example.com"}
	broker, _ := newTestBroker(user)

	// Create a valid token
	broker.CreateToken(context.Background(), user)

	// Try to reset with wrong token
	_, err := broker.Reset(context.Background(), user.email, "wrong-token", "new-pass", resetter{})
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

// Upstream: testRedirectReturnedByRemindWhenRecordDoesntExistInTable
func TestBrokerTokenExistsReturnsFalseForMissingToken(t *testing.T) {
	t.Parallel()

	user := &passwordUser{id: "user-1", email: "user@example.com"}
	broker, _ := newTestBroker(user)

	valid, _ := broker.TokenExists(context.Background(), user, "nonexistent-token")
	if valid {
		t.Fatal("expected false for missing token")
	}
}

// Upstream: testResetRemovesRecordOnReminderTableAndCallsCallback
func TestBrokerResetDeletesTokenAfterSuccess(t *testing.T) {
	t.Parallel()

	hasher, _ := auth.NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "old-pass")
	user := &passwordUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	broker, _ := newTestBroker(user)

	token, _ := broker.CreateToken(context.Background(), user)

	_, err := broker.Reset(context.Background(), user.email, token, "new-pass", resetter{hasher: hasher})
	if err != nil {
		t.Fatalf("Reset: %v", err)
	}

	// Token should be deleted after successful reset
	valid, _ := broker.TokenExists(context.Background(), user, token)
	if valid {
		t.Fatal("expected token to be deleted after reset")
	}
}

// Test DeleteToken explicitly
func TestBrokerDeleteToken(t *testing.T) {
	t.Parallel()

	user := &passwordUser{id: "user-1", email: "user@example.com"}
	broker, _ := newTestBroker(user)

	token, _ := broker.CreateToken(context.Background(), user)

	// Verify token exists
	valid, _ := broker.TokenExists(context.Background(), user, token)
	if !valid {
		t.Fatal("expected token to exist before delete")
	}

	// Delete token
	err := broker.DeleteToken(context.Background(), user)
	if err != nil {
		t.Fatalf("DeleteToken: %v", err)
	}

	// Verify token is gone
	valid, _ = broker.TokenExists(context.Background(), user, token)
	if valid {
		t.Fatal("expected token to be deleted")
	}
}

// Test token expiration
func TestBrokerTokenExistsReturnsFalseForExpiredToken(t *testing.T) {
	t.Parallel()

	user := &passwordUser{id: "user-1", email: "user@example.com"}
	now := time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)

	repo := passwords.NewMemoryTokenRepository()
	broker := &passwords.Broker{
		Config: passwords.Config{Name: "users", Expire: time.Hour, Throttle: time.Minute},
		Users:  &passwordProvider{user: user},
		Tokens: repo,
		Clock:  fixedClock{now: now},
	}

	token, _ := broker.CreateToken(context.Background(), user)

	// Advance clock past expiry
	broker.Clock = fixedClock{now: now.Add(2 * time.Hour)}

	_, err := broker.TokenExists(context.Background(), user, token)
	if !errors.Is(err, auth.ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

// Test CanCreateToken returns true when no token exists
func TestBrokerCanCreateTokenWhenNoTokenExists(t *testing.T) {
	t.Parallel()

	user := &passwordUser{id: "user-1", email: "user@example.com"}
	broker, _ := newTestBroker(user)

	canCreate, err := broker.CanCreateToken(context.Background(), user)
	if err != nil {
		t.Fatalf("CanCreateToken: %v", err)
	}
	if !canCreate {
		t.Fatal("expected CanCreateToken to return true when no token exists")
	}
}

// Test NormalizeEmail
func TestNormalizeEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"User@Example.COM", "user@example.com"},
		{"  user@test.com  ", "user@test.com"},
		{"", ""},
		{"UPPER@CASE.COM", "upper@case.com"},
	}

	for _, tc := range tests {
		if got := passwords.NormalizeEmail(tc.input); got != tc.expected {
			t.Fatalf("NormalizeEmail(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

// ======================== MEMORY TOKEN REPOSITORY TESTS ========================

// Upstream: testCreateInsertsNewRecordIntoTable
func TestMemoryRepoSave(t *testing.T) {
	t.Parallel()

	repo := passwords.NewMemoryTokenRepository()
	now := time.Now().UTC()

	token := &passwords.Token{
		UserID:    "user-1",
		TokenHash: "hash-1",
		CreatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}

	if err := repo.Save(context.Background(), token); err != nil {
		t.Fatalf("Save: %v", err)
	}

	found, err := repo.FindByTokenHash(context.Background(), "hash-1")
	if err != nil {
		t.Fatalf("FindByTokenHash: %v", err)
	}
	if found.UserID != "user-1" {
		t.Fatalf("expected user-1, got %q", found.UserID)
	}
}

// Upstream: testExistReturnsFalseIfNoRowFoundForUser
func TestMemoryRepoFindByTokenHashReturnsErrorWhenMissing(t *testing.T) {
	t.Parallel()

	repo := passwords.NewMemoryTokenRepository()
	_, err := repo.FindByTokenHash(context.Background(), "nonexistent")
	if !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

// Upstream: testDeleteMethodDeletesByToken
func TestMemoryRepoDeleteByTokenHash(t *testing.T) {
	t.Parallel()

	repo := passwords.NewMemoryTokenRepository()
	now := time.Now().UTC()
	repo.Save(context.Background(), &passwords.Token{
		UserID: "u1", TokenHash: "h1", CreatedAt: now, ExpiresAt: now.Add(time.Hour),
	})

	err := repo.DeleteByTokenHash(context.Background(), "h1")
	if err != nil {
		t.Fatalf("DeleteByTokenHash: %v", err)
	}

	_, err = repo.FindByTokenHash(context.Background(), "h1")
	if !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatal("expected token to be deleted")
	}
}

// Upstream: testDeleteExpiredMethodDeletesExpiredTokens (via DeleteByUserID)
func TestMemoryRepoDeleteByUserID(t *testing.T) {
	t.Parallel()

	repo := passwords.NewMemoryTokenRepository()
	now := time.Now().UTC()
	repo.Save(context.Background(), &passwords.Token{
		UserID: "u1", TokenHash: "h1", CreatedAt: now, ExpiresAt: now.Add(time.Hour),
	})

	err := repo.DeleteByUserID(context.Background(), "u1")
	if err != nil {
		t.Fatalf("DeleteByUserID: %v", err)
	}

	_, err = repo.FindByTokenHash(context.Background(), "h1")
	if !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatal("expected token to be deleted by user ID")
	}
}

// Upstream: testRecentlyCreatedReturnsTrueIfRecordIsRecentlyCreated
func TestMemoryRepoRecentlyCreatedTrue(t *testing.T) {
	t.Parallel()

	repo := passwords.NewMemoryTokenRepository()
	now := time.Now().UTC()
	repo.Save(context.Background(), &passwords.Token{
		UserID: "u1", TokenHash: "h1", CreatedAt: now, ExpiresAt: now.Add(time.Hour),
	})

	recent, err := repo.RecentlyCreated(context.Background(), "u1", now.Add(-time.Minute))
	if err != nil {
		t.Fatalf("RecentlyCreated: %v", err)
	}
	if !recent {
		t.Fatal("expected recently created to be true")
	}
}

// Upstream: testRecentlyCreatedReturnsFalseIfNoRowFoundForUser
func TestMemoryRepoRecentlyCreatedFalseForMissingUser(t *testing.T) {
	t.Parallel()

	repo := passwords.NewMemoryTokenRepository()
	recent, err := repo.RecentlyCreated(context.Background(), "missing", time.Now().UTC())
	if err != nil {
		t.Fatalf("RecentlyCreated: %v", err)
	}
	if recent {
		t.Fatal("expected recently created to be false for missing user")
	}
}

// Upstream: testRecentlyCreatedReturnsFalseIfValidRecordExists (old enough)
func TestMemoryRepoRecentlyCreatedFalseWhenOldEnough(t *testing.T) {
	t.Parallel()

	repo := passwords.NewMemoryTokenRepository()
	past := time.Now().UTC().Add(-10 * time.Minute)
	repo.Save(context.Background(), &passwords.Token{
		UserID: "u1", TokenHash: "h1", CreatedAt: past, ExpiresAt: past.Add(time.Hour),
	})

	// Check if recently created since "now" — the token was created 10 min ago
	recent, err := repo.RecentlyCreated(context.Background(), "u1", time.Now().UTC())
	if err != nil {
		t.Fatalf("RecentlyCreated: %v", err)
	}
	if recent {
		t.Fatal("expected recently created to be false when token is old enough")
	}
}

// Save replaces existing token for same user
func TestMemoryRepoSaveReplacesExisting(t *testing.T) {
	t.Parallel()

	repo := passwords.NewMemoryTokenRepository()
	now := time.Now().UTC()

	repo.Save(context.Background(), &passwords.Token{
		UserID: "u1", TokenHash: "old-hash", CreatedAt: now, ExpiresAt: now.Add(time.Hour),
	})
	repo.Save(context.Background(), &passwords.Token{
		UserID: "u1", TokenHash: "new-hash", CreatedAt: now, ExpiresAt: now.Add(time.Hour),
	})

	// New hash should be findable
	found, err := repo.FindByTokenHash(context.Background(), "new-hash")
	if err != nil {
		t.Fatalf("FindByTokenHash: %v", err)
	}
	if found.UserID != "u1" {
		t.Fatalf("expected u1, got %q", found.UserID)
	}
}

// Delete nonexistent token is a no-op
func TestMemoryRepoDeleteNonexistent(t *testing.T) {
	t.Parallel()

	repo := passwords.NewMemoryTokenRepository()
	err := repo.DeleteByTokenHash(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = repo.DeleteByUserID(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
