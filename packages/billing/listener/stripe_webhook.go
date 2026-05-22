// Package listener provides event listeners for billing webhook events.
package listener

import (
	"context"
	"time"

	"github.com/bedrock/packages/billing"
)

// StripeWebhookListener handles Stripe webhook events, particularly
// checkout.session.completed for marking orders as complete.
type StripeWebhookListener struct {
	orders billing.OrderStore
}

// NewStripeWebhookListener creates a StripeWebhookListener.
func NewStripeWebhookListener(orders billing.OrderStore) *StripeWebhookListener {
	return &StripeWebhookListener{orders: orders}
}

// HandleCheckoutSessionCompleted processes a completed Stripe checkout
// session and marks the associated order as complete.
func (l *StripeWebhookListener) HandleCheckoutSessionCompleted(ctx context.Context, payload map[string]any) error {
	data, _ := payload["data"].(map[string]any)

	if data == nil {
		return nil
	}

	object, _ := data["object"].(map[string]any)

	if object == nil {
		return nil
	}

	metadata, _ := object["metadata"].(map[string]any)

	if metadata == nil {
		return nil
	}

	orderIDRaw, ok := metadata["order_id"]

	if !ok {
		return nil
	}

	orderID, ok := orderIDRaw.(float64) // JSON numbers decode as float64

	if !ok {
		return nil
	}

	order, err := l.orders.FindByID(ctx, int64(orderID))

	if err != nil || order == nil {
		return err
	}

	// Use payment_intent if available, otherwise session ID.
	txnID, _ := object["payment_intent"].(string)

	if txnID == "" {
		txnID, _ = object["id"].(string)
	}

	order.MarkAsCompleted(txnID, time.Now())

	return l.orders.Save(ctx, order)
}
