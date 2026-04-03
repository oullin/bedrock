package passwords

import (
	"context"
	"fmt"
	"strings"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/support/crypto"
	configpkg "github.com/gollin/packages/config"
)

// Config controls password broker behavior.
type Config struct {
	Name     string
	Expire   time.Duration
	Throttle time.Duration
}

// ConfigFromRepository loads a password broker configuration.
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
	Mailer auth.Mailer
	Clock  auth.Clock
}

// SendResetLink creates and emails a password reset token.
func (b *Broker) SendResetLink(ctx context.Context, email string) (string, error) {
	user, err := b.Users.RetrieveByCredentials(ctx, map[string]string{"email": email})
	if err != nil {
		return "", err
	}

	profile, ok := user.(auth.UserProfile)
	if !ok {
		return "", fmt.Errorf("passwords: user does not expose email")
	}

	recentlyCreated, err := b.Tokens.RecentlyCreated(ctx, user.GetAuthIdentifier(), b.Clock.Now().Add(-b.Config.Throttle))
	if err != nil {
		return "", err
	}
	if recentlyCreated {
		return "", &auth.ThrottleError{Scope: "password-reset", RetryAfter: b.Config.Throttle}
	}

	token, err := crypto.RandomString(24)
	if err != nil {
		return "", fmt.Errorf("generate reset token: %w", err)
	}

	if err := b.Tokens.Save(ctx, &Token{
		UserID:    user.GetAuthIdentifier(),
		TokenHash: crypto.HashString(token),
		CreatedAt: b.Clock.Now(),
		ExpiresAt: b.Clock.Now().Add(b.Config.Expire),
	}); err != nil {
		return "", fmt.Errorf("save reset token: %w", err)
	}

	if err := b.Mailer.Send(ctx, auth.MailMessage{
		To:      profile.GetEmail(),
		Subject: "Reset your password",
		Body:    fmt.Sprintf("Use this password reset token for %s: %s", profile.GetEmail(), token),
		Metadata: map[string]string{
			"token": token,
		},
	}); err != nil {
		return "", fmt.Errorf("send reset email: %w", err)
	}

	return token, nil
}

// Reset validates a token and updates the user's password.
func (b *Broker) Reset(ctx context.Context, email string, token string, password string, action ResetsUserPasswords) (auth.Authenticatable, error) {
	tokenHash := crypto.HashString(token)
	record, err := b.Tokens.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if b.Clock.Now().After(record.ExpiresAt) {
		return nil, auth.ErrTokenExpired
	}

	user, err := b.Users.RetrieveByCredentials(ctx, map[string]string{"email": email})
	if err != nil {
		return nil, err
	}
	if user.GetAuthIdentifier() != record.UserID {
		return nil, auth.ErrInvalidToken
	}

	if err := action.Reset(ctx, user, password); err != nil {
		return nil, err
	}
	if err := b.Tokens.DeleteByTokenHash(ctx, tokenHash); err != nil {
		return nil, err
	}
	return user, nil
}

// NormalizeEmail normalizes the reset email field.
func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
