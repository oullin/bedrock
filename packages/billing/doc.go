// Package billing provides a provider-agnostic billing and subscription management
// system. It ports the Upstream Billing + Cashier Paddle billing stack to Go,
// offering plan catalog management, subscription lifecycle, entitlement access,
// checkout session handling, webhook processing, and a complete HTTP handler layer.
//
// The package abstracts the payment provider behind interfaces so that
// implementations for Paddle, Stripe, or any other gateway can be swapped
// without changing domain logic.
package billing
