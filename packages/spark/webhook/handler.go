package webhook

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/bedrock/packages/contracts/events"
	"github.com/bedrock/packages/spark"
)

// Handler dispatches payment provider webhook events to domain actions.
type Handler struct {
	subscriptions spark.SubscriptionStore
	customers     spark.CustomerStore
	transactions  spark.TransactionStore
	events        events.Dispatcher
}

// NewHandler creates a Handler.
func NewHandler(
	subscriptions spark.SubscriptionStore,
	customers spark.CustomerStore,
	transactions spark.TransactionStore,
	events events.Dispatcher,
) *Handler {
	return &Handler{
		subscriptions: subscriptions,
		customers:     customers,
		transactions:  transactions,
		events:        events,
	}
}

// Handle is the main entry point for webhook processing.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	var payload map[string]any

	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)

		return
	}

	_ = h.events.Dispatch(ctx, spark.WebhookReceivedEvent{Payload: payload})

	eventType, _ := payload["event_type"].(string)

	switch eventType {
	case "customer.updated":
		h.handleCustomerUpdated(ctx, payload)
	case "subscription.created":
		h.handleSubscriptionCreated(ctx, payload)
	case "subscription.updated":
		h.handleSubscriptionUpdated(ctx, payload)
	case "subscription.canceled":
		h.handleSubscriptionCanceled(ctx, payload)
	case "subscription.paused":
		h.handleSubscriptionPaused(ctx, payload)
	case "transaction.completed":
		h.handleTransactionCompleted(ctx, payload)
	case "transaction.updated":
		h.handleTransactionUpdated(ctx, payload)
	}

	_ = h.events.Dispatch(ctx, spark.WebhookHandledEvent{Payload: payload})
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleCustomerUpdated(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, spark.CustomerUpdatedEvent{Payload: payload})
}

func (h *Handler) handleSubscriptionCreated(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, spark.SubscriptionCreatedEvent{Payload: payload})
}

func (h *Handler) handleSubscriptionUpdated(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, spark.SubscriptionUpdatedEvent{Payload: payload})
}

func (h *Handler) handleSubscriptionCanceled(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, spark.SubscriptionCanceledEvent{Payload: payload})
}

func (h *Handler) handleSubscriptionPaused(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, spark.SubscriptionPausedEvent{Payload: payload})
}

func (h *Handler) handleTransactionCompleted(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, spark.TransactionCompletedEvent{Payload: payload})
}

func (h *Handler) handleTransactionUpdated(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, spark.TransactionUpdatedEvent{Payload: payload})
}
