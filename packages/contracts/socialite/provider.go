package socialauth

import "context"

// Provider handles a single OAuth redirect/callback cycle.
type Provider interface {
	// Redirect returns the URL the user should be sent to in order to
	// begin the OAuth flow.
	Redirect(ctx context.Context) (string, error)

	// User completes the OAuth callback and returns the authenticated user.
	User(ctx context.Context) (User, error)
}
