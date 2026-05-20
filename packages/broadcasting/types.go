package broadcasting

import (
	"context"

	contractsauth "github.com/bedrock/packages/contracts/auth"
)

// AuthRequest is the server-side channel authentication request.
type AuthRequest struct {
	ChannelName  string
	SocketID     string
	Callback     string
	UserResolver UserResolver
}

// UserResolver returns the authenticated user for a guard. An empty guard
// represents the upstream default request user lookup.
type UserResolver interface {
	User(guard string) contractsauth.Authenticatable
}

// UserResolverFunc adapts a function into a UserResolver.
type UserResolverFunc func(guard string) contractsauth.Authenticatable

// User returns the authenticated user for guard.

// ChannelHandler authorizes a channel subscription.
type ChannelHandler func(user contractsauth.Authenticatable, params ...any) (any, error)

// ChannelJoiner is the Go equivalent of upstream class-based channel joiners.
type ChannelJoiner interface {
	Join(user contractsauth.Authenticatable, params ...any) (any, error)
}

// ChannelOptions configures channel authorization.
type ChannelOptions struct {
	Guards []string
}

// WithGuards returns ChannelOptions with guard lookup order.

// BindingFunc resolves a named channel parameter.
type BindingFunc func(value string) (any, error)

// Broadcaster sends broadcast events to channels.
type Broadcaster interface {
	Broadcast(ctx context.Context, channels []string, event string, payload map[string]any) error
}

// Authenticator authenticates private and presence channel requests.
type Authenticator interface {
	Auth(request AuthRequest) (any, error)
	ValidAuthenticationResponse(request AuthRequest, result any) (any, error)
}

// Factory resolves named broadcaster connections.
type Factory interface {
	Connection(name string) (Broadcaster, error)
}

func (fn UserResolverFunc) User(guard string) contractsauth.Authenticatable {
	return fn(guard)
}

func WithGuards(guards ...string) ChannelOptions {
	return ChannelOptions{Guards: guards}
}
