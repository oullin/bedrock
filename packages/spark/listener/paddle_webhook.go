package listener

import (
	"context"
	"time"

	"github.com/bedrock/packages/spark"
)

// PaddleWebhookListener handles Paddle webhook events.
// Mirrors app/Listeners/PaddleWebhookListener.php.
type PaddleWebhookListener struct {
	orders spark.OrderStore
}

// NewPaddleWebhookListener creates a PaddleWebhookListener.
func NewPaddleWebhookListener(orders spark.OrderStore) *PaddleWebhookListener {
	return &PaddleWebhookListener{orders: orders}
}

// HandleTransactionCompleted marks the associated order as completed
// using the order_id from custom_data.
func (l *PaddleWebhookListener) HandleTransactionCompleted(ctx context.Context, payload map[string]any) error {
	data, _ := payload["data"].(map[string]any)

	if data == nil {
		return nil
	}

	customData, _ := data["custom_data"].(map[string]any)

	if customData == nil {
		return nil
	}

	orderIDRaw, ok := customData["order_id"]

	if !ok {
		return nil
	}

	orderID, ok := orderIDRaw.(float64)

	if !ok {
		return nil
	}

	order, err := l.orders.FindByID(ctx, int64(orderID))

	if err != nil || order == nil {
		return err
	}

	paddleID, _ := data["id"].(string)
	order.MarkAsCompleted(paddleID, time.Now())

	return l.orders.Save(ctx, order)
}
