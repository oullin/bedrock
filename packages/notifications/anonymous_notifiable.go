package notifications

import (
	"context"

	cn "github.com/bedrock/packages/contracts/notifications"
)

// AnonymousNotifiable enables on-demand notifications without requiring a
// persistent model. Routes are registered per channel.
type AnonymousNotifiable struct {
	routes map[string]any
}

// compile-time interface check.
var _ cn.Notifiable = (*AnonymousNotifiable)(nil)

// NewAnonymousNotifiable creates an AnonymousNotifiable.
func NewAnonymousNotifiable() *AnonymousNotifiable {
	return &AnonymousNotifiable{routes: make(map[string]any)}
}

// Route registers a routing value for the named channel and returns the
// notifiable for chaining. The database channel is not allowed for anonymous
// notifiables.
func (n *AnonymousNotifiable) Route(channel string, route any) *AnonymousNotifiable {
	if channel == "database" {
		return n
	}

	n.routes[channel] = route

	return n
}

// RouteNotificationFor returns the routing information for the given channel.
func (n *AnonymousNotifiable) RouteNotificationFor(_ context.Context, channel string) any {
	if route, ok := n.routes[channel]; ok {
		return route
	}

	return nil
}

// GetKey returns an empty string since anonymous notifiables have no identity.
func (n *AnonymousNotifiable) GetKey() string { return "" }

// Notify dispatches a notification to this anonymous notifiable via the
// provided manager.
func (n *AnonymousNotifiable) Notify(ctx context.Context, manager cn.Dispatcher, notification any) error {
	return manager.Send(ctx, []cn.Notifiable{n}, notification)
}

// NotifyNow dispatches a notification synchronously to this anonymous
// notifiable via the provided manager.
func (n *AnonymousNotifiable) NotifyNow(ctx context.Context, manager cn.Dispatcher, notification any, channels ...string) error {
	return manager.SendNow(ctx, []cn.Notifiable{n}, notification, channels...)
}
