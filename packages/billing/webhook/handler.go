package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bedrock/packages/billing"
)

// Handler processes incoming webhook payloads from payment providers.
// Mirrors Upstream\Paddle\Http\Controllers\WebhookController and
// Billing\Http\Controllers\WebhookController.
type Handler struct {
	subscriptions billing.SubscriptionStore
	customers     billing.CustomerStore
	transactions  billing.TransactionStore
	events        billing.EventDispatcher
}

// NewHandler creates a webhook Handler.
func NewHandler(
	subs billing.SubscriptionStore,
	customers billing.CustomerStore,
	txns billing.TransactionStore,
	events billing.EventDispatcher,
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
		h.events.Dispatch(billing.WebhookReceivedEvent{Payload: payload})
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
		h.events.Dispatch(billing.WebhookHandledEvent{Payload: payload})
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

func (h *Handler) handleSubscriptionCreated(r *http.Request, payload map[string]any) {
	h.syncSubscription(r.Context(), payload, billing.StatusActive)
}

func (h *Handler) handleSubscriptionUpdated(r *http.Request, payload map[string]any) {
	h.syncSubscription(r.Context(), payload, "")
}

func (h *Handler) handleSubscriptionCanceled(r *http.Request, payload map[string]any) {
	h.syncSubscription(r.Context(), payload, billing.StatusCanceled)
}

func (h *Handler) handleSubscriptionPaused(r *http.Request, payload map[string]any) {
	h.syncSubscription(r.Context(), payload, billing.StatusPaused)
}

func (h *Handler) handleTransactionCompleted(r *http.Request, payload map[string]any) {
	h.syncTransaction(r.Context(), payload, billing.TransactionCompleted)
}

func (h *Handler) handleTransactionUpdated(r *http.Request, payload map[string]any) {
	h.syncTransaction(r.Context(), payload, "")
}

func (h *Handler) syncSubscription(ctx context.Context, payload map[string]any, fallback billing.SubscriptionStatus) {
	if h.subscriptions == nil {
		return
	}

	data := payloadData(payload)
	paddleID := stringField(data, "id")

	if paddleID == "" {
		return
	}

	subscription, err := h.subscriptions.FindByProviderID(ctx, paddleID)

	if err != nil || subscription == nil {
		return
	}

	status := fallback

	if value := stringField(data, "status"); value != "" {
		status = billing.SubscriptionStatus(value)
	}

	if status != "" {
		subscription.Status = status
	}

	if items := subscriptionItems(subscription.ID, data); len(items) > 0 {
		subscription.Items = items
	}

	_ = h.subscriptions.Save(ctx, subscription)

	if h.events == nil {
		return
	}

	switch status {
	case billing.StatusCanceled:
		_ = h.events.Dispatch(billing.SubscriptionCanceledEvent{Subscription: subscription, Payload: payload})
	case billing.StatusPaused:
		_ = h.events.Dispatch(billing.SubscriptionPausedEvent{Subscription: subscription, Payload: payload})
	default:
		_ = h.events.Dispatch(billing.SubscriptionUpdatedEvent{Subscription: subscription, Payload: payload})
	}
}

func (h *Handler) syncTransaction(ctx context.Context, payload map[string]any, fallback billing.TransactionStatus) {
	if h.transactions == nil {
		return
	}

	data := payloadData(payload)
	paddleID := stringField(data, "id")

	if paddleID == "" {
		return
	}

	transaction, err := h.transactions.FindByProviderID(ctx, paddleID)

	if err != nil {
		return
	}

	create := transaction == nil

	if transaction == nil {
		transaction = &billing.Transaction{PaddleID: paddleID}
	}

	status := fallback

	if value := stringField(data, "status"); value != "" {
		status = billing.TransactionStatus(value)
	}

	total, tax, currency, err := transactionMoney(data)

	if err != nil {
		return
	}

	if status != "" {
		transaction.Status = status
	}

	transaction.Total = total
	transaction.Tax = tax
	transaction.Currency = currency

	if create {
		_ = h.transactions.Create(ctx, transaction)
	} else {
		_ = h.transactions.Save(ctx, transaction)
	}

	if h.events == nil {
		return
	}

	switch status {
	case billing.TransactionCompleted:
		_ = h.events.Dispatch(billing.TransactionCompletedEvent{Transaction: transaction, Payload: payload})
	default:
		_ = h.events.Dispatch(billing.TransactionUpdatedEvent{Transaction: transaction, Payload: payload})
	}
}

func payloadData(payload map[string]any) map[string]any {
	data, _ := payload["data"].(map[string]any)

	if data == nil {
		return map[string]any{}
	}

	return data
}

func stringField(data map[string]any, key string) string {
	value, _ := data[key].(string)

	return value
}

func transactionMoney(data map[string]any) (int64, int64, string, error) {
	details, _ := data["details"].(map[string]any)
	totals, _ := details["totals"].(map[string]any)

	if totals == nil {
		return 0, 0, "", billing.ValidationErrors{{Field: "details.totals", Message: "is required"}}
	}

	total, err := billing.ParseMinorAmount(totals["total"])

	if err != nil {
		return 0, 0, "", err
	}

	tax, err := billing.ParseMinorAmount(totals["tax"])

	if err != nil {
		return 0, 0, "", err
	}

	currency := stringField(totals, "currency_code")

	if currency == "" {
		currency = stringField(data, "currency_code")
	}

	currency = strings.ToUpper(strings.TrimSpace(currency))

	if err := billing.ValidateTransactionMoney(total, tax, currency); err != nil {
		return 0, 0, "", err
	}

	return total, tax, currency, nil
}

func subscriptionItems(subscriptionID int64, data map[string]any) []billing.SubscriptionItem {
	rawItems, _ := data["items"].([]any)
	items := make([]billing.SubscriptionItem, 0, len(rawItems))

	for _, raw := range rawItems {
		item, _ := raw.(map[string]any)

		if item == nil {
			continue
		}

		price, _ := item["price"].(map[string]any)

		if price == nil {
			continue
		}

		priceID := stringField(price, "id")

		if priceID == "" {
			continue
		}

		quantity := 1

		if rawQuantity, ok := item["quantity"].(float64); ok && rawQuantity > 0 {
			quantity = int(rawQuantity)
		}

		items = append(items, billing.SubscriptionItem{
			SubscriptionID: subscriptionID,
			PriceID:        priceID,
			ProductID:      stringField(price, "product_id"),
			Quantity:       quantity,
			Status:         string(billing.StatusActive),
		})
	}

	return items
}
