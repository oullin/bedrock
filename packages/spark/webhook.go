package spark

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// WebhookHandler dispatches payment provider webhook events to domain actions.
type WebhookHandler struct {
	subscriptions SubscriptionStore
	customers     CustomerStore
	transactions  TransactionStore
	events        EventDispatcher
}

// NewWebhookHandler creates a WebhookHandler.
func NewWebhookHandler(
	subscriptions SubscriptionStore,
	customers CustomerStore,
	transactions TransactionStore,
	events EventDispatcher,
) *WebhookHandler {
	return &WebhookHandler{
		subscriptions: subscriptions,
		customers:     customers,
		transactions:  transactions,
		events:        events,
	}
}

// Handle is the main entry point for webhook processing.
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	_ = h.events.Dispatch(ctx, WebhookReceivedEvent{Payload: payload})

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

	_ = h.events.Dispatch(ctx, WebhookHandledEvent{Payload: payload})
	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) handleCustomerUpdated(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, CustomerUpdatedEvent{Payload: payload})
}

func (h *WebhookHandler) handleSubscriptionCreated(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, SubscriptionCreatedEvent{Payload: payload})
}

func (h *WebhookHandler) handleSubscriptionUpdated(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, SubscriptionUpdatedEvent{Payload: payload})
}

func (h *WebhookHandler) handleSubscriptionCanceled(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, SubscriptionCanceledEvent{Payload: payload})
}

func (h *WebhookHandler) handleSubscriptionPaused(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, SubscriptionPausedEvent{Payload: payload})
}

func (h *WebhookHandler) handleTransactionCompleted(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, TransactionCompletedEvent{Payload: payload})
}

func (h *WebhookHandler) handleTransactionUpdated(ctx context.Context, payload map[string]any) {
	_ = h.events.Dispatch(ctx, TransactionUpdatedEvent{Payload: payload})
}
