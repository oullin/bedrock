package passwords

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	authevents "github.com/bedrock/packages/auth/events"
	cauth "github.com/bedrock/packages/contracts/auth"
	cevents "github.com/bedrock/packages/contracts/events"
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

// RecentTokenRepository can report whether a token was created recently enough
// to throttle repeated reset-link requests.
type RecentTokenRepository interface {
	RecentlyCreated(ctx context.Context, email string, within time.Duration) bool
}

// ResetCallback is called with the user and plain-text token to perform the reset.
type ResetCallback func(ctx context.Context, user cauth.CanResetPassword, token, password string) error

// ResetLinkCallback is called after a reset token is created to customize how
// the reset link notification is sent.
type ResetLinkCallback func(ctx context.Context, user cauth.CanResetPassword, token string) error

// Broker orchestrates the password reset flow.
type Broker struct {
	users    cauth.UserProvider
	tokens   TokenRepository
	expiry   time.Duration
	throttle time.Duration
	events   cevents.Dispatcher
}

// NewBroker creates a Broker. expiry is the token lifetime.
func NewBroker(users cauth.UserProvider, tokens TokenRepository, expiry time.Duration) *Broker {
	return &Broker{users: users, tokens: tokens, expiry: expiry}
}

// WithThrottle configures how long reset-link creation should be throttled for
// a user after a token has already been created.
func (b *Broker) WithThrottle(throttle time.Duration) *Broker {
	b.throttle = throttle

	return b
}

// WithEventDispatcher configures the broker's auth event dispatcher.
func (b *Broker) WithEventDispatcher(dispatcher cevents.Dispatcher) *Broker {
	b.events = dispatcher

	return b
}

// SetEventDispatcher configures the broker's auth event dispatcher.
func (b *Broker) SetEventDispatcher(dispatcher cevents.Dispatcher) {
	b.events = dispatcher
}

func (b *Broker) dispatch(ctx context.Context, event any) {
	if b.events != nil {
		_, _ = b.events.Dispatch(ctx, event)
	}
}

// SendResetLink finds the user by email and sends them a password reset notification.
func (b *Broker) SendResetLink(ctx context.Context, email string) error {
	return b.SendResetLinkUsing(ctx, email, nil)
}

// SendResetLinkUsing finds the user by email, creates a reset token, and lets
// the callback customize how the notification is sent.
func (b *Broker) SendResetLinkUsing(ctx context.Context, email string, callback ResetLinkCallback) error {
	user, err := b.getUser(ctx, email)

	if err != nil {
		return err
	}

	if b.throttle > 0 {
		if recent, ok := b.tokens.(RecentTokenRepository); ok && recent.RecentlyCreated(ctx, email, b.throttle) {
			return errors.New("passwords: token recently created")
		}
	}

	token, err := b.tokens.Create(ctx, email)

	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ctx, user, token)
	}

	if sender, ok := user.(cauth.PasswordResetNotificationSender); ok {
		sender.SendPasswordResetNotification(token)
	}

	b.dispatch(ctx, authevents.PasswordResetLinkSent{User: user})

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
