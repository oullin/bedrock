package oauthserver

import "context"

// EventDispatcher is the minimal interface for dispatching OAuthServer lifecycle events.
// It accepts the contracts/events.Dispatcher implementation but avoids importing
// that package as a dependency, keeping the oauthserver module's dependency graph lean.
//
// Callers may pass nil when events are not needed — all dispatch sites guard
// against nil before calling.
type EventDispatcher interface {
	Dispatch(ctx context.Context, event any) ([]any, error)
}

// AccessTokenCreated is dispatched when a new access token is created.
// Mirrors Upstream OAuthServer's AccessTokenCreated event.
type AccessTokenCreated struct {
	TokenID  string
	UserID   string
	ClientID string
}

// AccessTokenRevoked is dispatched when an access token is revoked.
// Mirrors Upstream OAuthServer's AccessTokenRevoked event.
type AccessTokenRevoked struct {
	TokenID string
}

// RefreshTokenCreated is dispatched when a new refresh token is created.
// Mirrors Upstream OAuthServer's RefreshTokenCreated event.
type RefreshTokenCreated struct {
	ID            string
	AccessTokenID string
}

// dispatch is a nil-safe helper that fires an event if a dispatcher is set.
func dispatch(ctx context.Context, d EventDispatcher, event any) {
	if d == nil {
		return
	}

	_, _ = d.Dispatch(ctx, event)
}
