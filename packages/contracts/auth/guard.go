package auth

import (
	"context"
	"net/http"
)

// Guard performs read-only authentication checks.
type Guard interface {
	User(ctx context.Context) (Authenticatable, error)
	Check(ctx context.Context) bool
	Guest(ctx context.Context) bool
	ID(ctx context.Context) any
}

// StatefulGuard extends Guard with login/logout state management.
type StatefulGuard interface {
	Guard
	Validate(ctx context.Context, credentials map[string]string) bool
	Attempt(ctx context.Context, credentials map[string]string, remember bool) bool
	Once(ctx context.Context, credentials map[string]string) bool
	Login(ctx context.Context, user Authenticatable, remember bool) error
	LoginUsingID(ctx context.Context, id string, remember bool) (Authenticatable, error)
	OnceUsingID(ctx context.Context, id string) (Authenticatable, error)
	ViaRemember(ctx context.Context) bool
	Logout(ctx context.Context) error
}

// SupportsBasicAuth is implemented by guards that can authenticate HTTP Basic
// credentials.
type SupportsBasicAuth interface {
	Basic(ctx context.Context, field string, extraConditions map[string]string) bool
	OnceBasic(ctx context.Context, field string, extraConditions map[string]string) bool
}

// HTTPGuard provides HTTP-aware authentication for web frameworks.
type HTTPGuard interface {
	Name() string
	AuthenticateRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (Authenticatable, error)
	Login(ctx context.Context, w http.ResponseWriter, user Authenticatable, remember bool) error
	LoginWithPendingTwoFactor(ctx context.Context, w http.ResponseWriter, user Authenticatable) error
	Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}
