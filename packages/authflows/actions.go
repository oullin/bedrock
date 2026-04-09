package authflows

import "context"

// CreatesNewUsers creates a new user from registration input.
type CreatesNewUsers interface {
	Create(ctx context.Context, input map[string]string) (Authenticatable, error)
}

// AuthenticatesUsers customizes the authentication attempt.
// Implement this to override the default credential-based login.
type AuthenticatesUsers interface {
	Authenticate(ctx context.Context, input map[string]string) (Authenticatable, error)
}

// UpdatesUserProfileInformation updates a user's profile.
type UpdatesUserProfileInformation interface {
	Update(ctx context.Context, user Authenticatable, input map[string]string) error
}

// UpdatesUserPasswords changes an authenticated user's password.
type UpdatesUserPasswords interface {
	Update(ctx context.Context, user Authenticatable, input map[string]string) error
}

// ResetsUserPasswords sets a new password during the reset flow.
type ResetsUserPasswords interface {
	Reset(ctx context.Context, user Authenticatable, password string) error
}

// ConfirmsPasswords validates the user's current password.
type ConfirmsPasswords interface {
	Confirm(ctx context.Context, user Authenticatable, password string) error
}
