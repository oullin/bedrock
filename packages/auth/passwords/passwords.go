package passwords

import (
	"context"
	"fmt"
	"strings"
	"time"

	auth "github.com/gollin/packages/auth"
	configpkg "github.com/gollin/packages/config"
	securitycrypto "github.com/gollin/packages/security/crypto"
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

// Broker issues password reset links and tokens.
type Broker struct {
	Config    Config
	Users     auth.UserProvider
	Tokens    TokenRepository
	Mailer    auth.Mailer
	Clock     auth.Clock
	CreateURL func(user auth.Authenticatable, token string) string
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

// SendResetLink creates a reset token and sends a message.
func (b *Broker) SendResetLink(ctx context.Context, email string) (string, error) {
	user, err := b.Users.RetrieveByCredentials(ctx, map[string]string{"email": NormalizeEmail(email)})
	if err != nil {
		return "", err
	}

	profile, ok := user.(auth.UserProfile)
	if !ok {
		return "", fmt.Errorf("passwords: user does not expose email")
	}

	if b.Clock == nil {
		b.Clock = auth.SystemClock{}
	}

	recentlyCreated, err := b.Tokens.RecentlyCreated(ctx, user.GetAuthIdentifier(), b.Clock.Now().Add(-b.Config.Throttle))
	if err != nil {
		return "", err
	}

	if recentlyCreated {
		return "", auth.ErrInvalidToken
	}

	token, err := b.CreateToken(ctx, user)
	if err != nil {
		return "", err
	}

	url := ""
	if b.CreateURL != nil {
		url = b.CreateURL(user, token)
	}

	if b.Mailer != nil {
		if err := b.Mailer.Send(ctx, auth.MailMessage{
			To:      profile.GetEmail(),
			Subject: "Reset your password",
			Body:    "Reset your password by visiting " + url,
			Metadata: map[string]string{
				"token": token,
				"url":   url,
			},
		}); err != nil {
			return "", err
		}
	}

	return token, nil
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

	if err := b.Tokens.Save(ctx, &Token{
		UserID:    user.GetAuthIdentifier(),
		TokenHash: securitycrypto.HashString(token),
		CreatedAt: b.Clock.Now(),
		ExpiresAt: b.Clock.Now().Add(b.Config.Expire),
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
