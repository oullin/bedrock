package webhook

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/spark"
)

// Handler processes incoming webhook payloads from payment providers.
// Mirrors Laravel\Paddle\Http\Controllers\WebhookController and
// Spark\Http\Controllers\WebhookController.
type Handler struct {
	subscriptions spark.SubscriptionStore
	customers     spark.CustomerStore
	transactions  spark.TransactionStore
	events        spark.EventDispatcher
}

// NewHandler creates a webhook Handler.
func NewHandler(
	subs spark.SubscriptionStore,
	customers spark.CustomerStore,
	txns spark.TransactionStore,
	events spark.EventDispatcher,
) *Handler {
	return &Handler{
		subscriptions: subs,
		customers:     customers,
		transactions:  txns,
		events:        events,
	}
}

// Handle dispatches webhook payloads by event type.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)

		return
	}

	if h.events != nil {
		h.events.Dispatch(spark.WebhookReceivedEvent{Payload: payload})
	}

	eventType, _ := payload["event_type"].(string)

	switch eventType {
	case "customer.updated":
		h.handleCustomerUpdated(r, payload)
	case "subscription.created":
		h.handleSubscriptionCreated(r, payload)
	case "subscription.updated":
		h.handleSubscriptionUpdated(r, payload)
	case "subscription.canceled":
		h.handleSubscriptionCanceled(r, payload)
	case "subscription.paused":
		h.handleSubscriptionPaused(r, payload)
	case "transaction.completed":
		h.handleTransactionCompleted(r, payload)
	case "transaction.updated":
		h.handleTransactionUpdated(r, payload)
	}

	if h.events != nil {
		h.events.Dispatch(spark.WebhookHandledEvent{Payload: payload})
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleCustomerUpdated(_ *http.Request, payload map[string]any) {
	data, _ := payload["data"].(map[string]any)

	if data == nil {
		return
	}

	paddleID, _ := data["id"].(string)

	if paddleID == "" {
		return
	}

	ctx := context.Background()
	customer, err := h.customers.FindByProviderID(ctx, paddleID)

	if err != nil || customer == nil {
		return
	}

	if name, ok := data["name"].(string); ok {
		customer.Name = name
	}

	if email, ok := data["email"].(string); ok {
		customer.Email = email
	}

	h.customers.Save(ctx, customer)
}

func (h *Handler) handleSubscriptionCreated(_ *http.Request, _ map[string]any)  {}
func (h *Handler) handleSubscriptionUpdated(_ *http.Request, _ map[string]any)  {}
func (h *Handler) handleSubscriptionCanceled(_ *http.Request, _ map[string]any) {}
func (h *Handler) handleSubscriptionPaused(_ *http.Request, _ map[string]any)   {}
func (h *Handler) handleTransactionCompleted(_ *http.Request, _ map[string]any) {}
func (h *Handler) handleTransactionUpdated(_ *http.Request, _ map[string]any)   {}
