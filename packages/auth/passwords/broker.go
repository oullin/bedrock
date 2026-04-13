package passwords

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// TokenRepository stores and validates password reset tokens.
type TokenRepository interface {
	// Create stores a token for the given email and returns it.
	Create(ctx context.Context, email string) (string, error)
	// Exists reports whether the token is valid and not expired.
	Exists(ctx context.Context, email, token string) bool
	// Delete removes the token for the given email.
	Delete(ctx context.Context, email string) error
	// DeleteExpired removes all expired tokens.
	DeleteExpired(ctx context.Context) error
}

// ResetCallback is called with the user and plain-text token to perform the reset.
type ResetCallback func(ctx context.Context, user cauth.CanResetPassword, token, password string) error

// Broker orchestrates the password reset flow.
type Broker struct {
	users  cauth.UserProvider
	tokens TokenRepository
	expiry time.Duration
}

// NewBroker creates a Broker. expiry is the token lifetime.
func NewBroker(users cauth.UserProvider, tokens TokenRepository, expiry time.Duration) *Broker {
	return &Broker{users: users, tokens: tokens, expiry: expiry}
}

// SendResetLink finds the user by email and sends them a password reset notification.
func (b *Broker) SendResetLink(ctx context.Context, email string) error {
	user, err := b.getUser(ctx, email)

	if err != nil {
		return err
	}

	_, err = b.tokens.Create(ctx, email)

	if err != nil {
		return err
	}

	_ = user

	return nil
}

// Reset validates the token and invokes the reset callback.
func (b *Broker) Reset(ctx context.Context, credentials map[string]any, resetFn ResetCallback) error {
	email, _ := credentials["email"].(string)
	token, _ := credentials["token"].(string)
	password, _ := credentials["password"].(string)

	user, err := b.getUser(ctx, email)

	if err != nil {
		return err
	}

	if !b.tokens.Exists(ctx, email, token) {
		return errors.New("passwords: invalid or expired token")
	}

	if err = resetFn(ctx, user, token, password); err != nil {
		return err
	}

	return b.tokens.Delete(ctx, email)
}

// GetUser retrieves the user by email for password reset.
func (b *Broker) GetUser(ctx context.Context, email string) (cauth.CanResetPassword, error) {
	return b.getUser(ctx, email)
}

// CreateToken creates a password reset token for the given user.
func (b *Broker) CreateToken(ctx context.Context, user cauth.CanResetPassword) (string, error) {
	return b.tokens.Create(ctx, user.GetEmailForPasswordReset())
}

// DeleteToken removes the password reset token for the given user.
func (b *Broker) DeleteToken(ctx context.Context, user cauth.CanResetPassword) error {
	return b.tokens.Delete(ctx, user.GetEmailForPasswordReset())
}

// TokenExists reports whether a valid token exists for the given user.
func (b *Broker) TokenExists(ctx context.Context, user cauth.CanResetPassword, token string) bool {
	return b.tokens.Exists(ctx, user.GetEmailForPasswordReset(), token)
}

// GetRepository returns the token repository.
func (b *Broker) GetRepository() TokenRepository {
	return b.tokens
}

func (b *Broker) getUser(ctx context.Context, email string) (cauth.CanResetPassword, error) {
	u, err := b.users.RetrieveByCredentials(ctx, map[string]string{"email": email})

	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, errors.New("passwords: user not found")
	}

	crp, ok := u.(cauth.CanResetPassword)

	if !ok {
		return nil, errors.New("passwords: user does not implement CanResetPassword")
	}

	return crp, nil
}

// GenerateToken creates a cryptographically random reset token.
func GenerateToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
