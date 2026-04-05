package passwords

import (
	"context"
	"strings"
	"time"

	auth "github.com/gollin/packages/illuminate/auth"
	configpkg "github.com/gollin/packages/illuminate/config"
	securitycrypto "github.com/gollin/packages/illuminate/support/crypto"
)

// Config controls password broker behavior.
type Config struct {
	Name     string
	Expire   time.Duration
	Throttle time.Duration
}

// Token stores password reset token metadata.
type Token struct {
	UserID    string
	TokenHash string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// TokenRepository persists password reset tokens.
type TokenRepository interface {
	Save(ctx context.Context, token *Token) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*Token, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteByUserID(ctx context.Context, userID string) error
	RecentlyCreated(ctx context.Context, userID string, since time.Time) (bool, error)
}

// ResetsUserPasswords updates a user's password.
type ResetsUserPasswords interface {
	Reset(ctx context.Context, user auth.Authenticatable, password string) error
}

// Broker issues and validates password reset tokens.
type Broker struct {
	Config Config
	Users  auth.UserProvider
	Tokens TokenRepository
	Clock  auth.Clock
}

// ConfigFromRepository loads a broker config from the shared config repository.
func ConfigFromRepository(repo *configpkg.Repository, brokerName string) (Config, error) {
	expire, err := repo.Duration("auth.passwords." + brokerName + ".expire")
	if err != nil {
		return Config{}, err
	}

	throttle, err := repo.Duration("auth.passwords." + brokerName + ".throttle")
	if err != nil {
		return Config{}, err
	}

	return Config{
		Name:     brokerName,
		Expire:   expire,
		Throttle: throttle,
	}, nil
}

// CreateToken creates and stores a new reset token.
func (b *Broker) CreateToken(ctx context.Context, user auth.Authenticatable) (string, error) {
	if b.Clock == nil {
		b.Clock = auth.SystemClock{}
	}

	token, err := securitycrypto.RandomString(24)
	if err != nil {
		return "", err
	}

	if err := b.Tokens.DeleteByUserID(ctx, user.GetAuthIdentifier()); err != nil {
		return "", err
	}

	now := b.Clock.Now()
	if err := b.Tokens.Save(ctx, &Token{
		UserID:    user.GetAuthIdentifier(),
		TokenHash: securitycrypto.HashString(token),
		CreatedAt: now,
		ExpiresAt: now.Add(b.Config.Expire),
	}); err != nil {
		return "", err
	}

	return token, nil
}

// DeleteToken removes all reset tokens for a user.
func (b *Broker) DeleteToken(ctx context.Context, user auth.Authenticatable) error {
	return b.Tokens.DeleteByUserID(ctx, user.GetAuthIdentifier())
}

// TokenExists reports whether the token is valid for a user.
func (b *Broker) TokenExists(ctx context.Context, user auth.Authenticatable, token string) (bool, error) {
	if b.Clock == nil {
		b.Clock = auth.SystemClock{}
	}

	record, err := b.Tokens.FindByTokenHash(ctx, securitycrypto.HashString(token))
	if err != nil {
		if err == auth.ErrInvalidToken {
			return false, nil
		}

		return false, err
	}

	if b.Clock.Now().After(record.ExpiresAt) {
		return false, auth.ErrTokenExpired
	}

	return record.UserID == user.GetAuthIdentifier(), nil
}

// CanCreateToken reports whether the broker throttle window allows a new token.
func (b *Broker) CanCreateToken(ctx context.Context, user auth.Authenticatable) (bool, error) {
	if b.Clock == nil {
		b.Clock = auth.SystemClock{}
	}

	recentlyCreated, err := b.Tokens.RecentlyCreated(ctx, user.GetAuthIdentifier(), b.Clock.Now().Add(-b.Config.Throttle))
	if err != nil {
		return false, err
	}

	return !recentlyCreated, nil
}

// Reset validates a token and updates the user's password.
func (b *Broker) Reset(ctx context.Context, email string, token string, password string, action ResetsUserPasswords) (auth.Authenticatable, error) {
	user, err := b.Users.RetrieveByCredentials(ctx, map[string]string{"email": NormalizeEmail(email)})
	if err != nil {
		return nil, err
	}

	valid, err := b.TokenExists(ctx, user, token)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, auth.ErrInvalidToken
	}

	if err := action.Reset(ctx, user, password); err != nil {
		return nil, err
	}

	if err := b.DeleteToken(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// NormalizeEmail normalizes a password-reset email address.
func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
