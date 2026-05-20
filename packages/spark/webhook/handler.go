package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bedrock/packages/spark"
)

// Handler processes incoming webhook payloads from payment providers.
// Mirrors upstream Paddle\Http\Controllers\WebhookController and
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

func (h *Handler) handleSubscriptionCreated(r *http.Request, payload map[string]any) {
	h.syncSubscription(r.Context(), payload, spark.StatusActive)
}

func (h *Handler) handleSubscriptionUpdated(r *http.Request, payload map[string]any) {
	h.syncSubscription(r.Context(), payload, "")
}

func (h *Handler) handleSubscriptionCanceled(r *http.Request, payload map[string]any) {
	h.syncSubscription(r.Context(), payload, spark.StatusCanceled)
}

func (h *Handler) handleSubscriptionPaused(r *http.Request, payload map[string]any) {
	h.syncSubscription(r.Context(), payload, spark.StatusPaused)
}

func (h *Handler) handleTransactionCompleted(r *http.Request, payload map[string]any) {
	h.syncTransaction(r.Context(), payload, spark.TransactionCompleted)
}

func (h *Handler) handleTransactionUpdated(r *http.Request, payload map[string]any) {
	h.syncTransaction(r.Context(), payload, "")
}

func (h *Handler) syncSubscription(ctx context.Context, payload map[string]any, fallback spark.SubscriptionStatus) {
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
		status = spark.SubscriptionStatus(value)
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
	case spark.StatusCanceled:
		_ = h.events.Dispatch(spark.SubscriptionCanceledEvent{Subscription: subscription, Payload: payload})
	case spark.StatusPaused:
		_ = h.events.Dispatch(spark.SubscriptionPausedEvent{Subscription: subscription, Payload: payload})
	default:
		_ = h.events.Dispatch(spark.SubscriptionUpdatedEvent{Subscription: subscription, Payload: payload})
	}
}

func (h *Handler) syncTransaction(ctx context.Context, payload map[string]any, fallback spark.TransactionStatus) {
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
		transaction = &spark.Transaction{PaddleID: paddleID}
	}

	status := fallback

	if value := stringField(data, "status"); value != "" {
		status = spark.TransactionStatus(value)
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
	case spark.TransactionCompleted:
		_ = h.events.Dispatch(spark.TransactionCompletedEvent{Transaction: transaction, Payload: payload})
	default:
		_ = h.events.Dispatch(spark.TransactionUpdatedEvent{Transaction: transaction, Payload: payload})
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
		return 0, 0, "", spark.ValidationErrors{{Field: "details.totals", Message: "is required"}}
	}

	total, err := spark.ParseMinorAmount(totals["total"])

	if err != nil {
		return 0, 0, "", err
	}

	tax, err := spark.ParseMinorAmount(totals["tax"])

	if err != nil {
		return 0, 0, "", err
	}

	currency := stringField(totals, "currency_code")

	if currency == "" {
		currency = stringField(data, "currency_code")
	}

	currency = strings.ToUpper(strings.TrimSpace(currency))

	if err := spark.ValidateTransactionMoney(total, tax, currency); err != nil {
		return 0, 0, "", err
	}

	return total, tax, currency, nil
}

func subscriptionItems(subscriptionID int64, data map[string]any) []spark.SubscriptionItem {
	rawItems, _ := data["items"].([]any)
	items := make([]spark.SubscriptionItem, 0, len(rawItems))

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

		items = append(items, spark.SubscriptionItem{
			SubscriptionID: subscriptionID,
			PriceID:        priceID,
			ProductID:      stringField(price, "product_id"),
			Quantity:       quantity,
			Status:         string(spark.StatusActive),
		})
	}

	return items
}
