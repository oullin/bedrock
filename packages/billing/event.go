package billing

// ---------------------------------------------------------------------------
// Cashier events
// ---------------------------------------------------------------------------

// CustomerUpdatedEvent is dispatched when a customer record is updated via webhook.
type CustomerUpdatedEvent struct {
	BillableID int64
	Customer   *Customer
	Payload    map[string]any
}

// SubscriptionCreatedEvent is dispatched when a new subscription is created via webhook.
type SubscriptionCreatedEvent struct {
	BillableID   int64
	Subscription *Subscription
	Payload      map[string]any
}

// SubscriptionUpdatedEvent is dispatched when a subscription is updated via webhook.
type SubscriptionUpdatedEvent struct {
	Subscription *Subscription
	Payload      map[string]any
}

// SubscriptionCanceledEvent is dispatched when a subscription is canceled via webhook.
type SubscriptionCanceledEvent struct {
	Subscription *Subscription
	Payload      map[string]any
}

// SubscriptionPausedEvent is dispatched when a subscription is paused via webhook.
type SubscriptionPausedEvent struct {
	Subscription *Subscription
	Payload      map[string]any
}

// TransactionCompletedEvent is dispatched when a transaction completes via webhook.
type TransactionCompletedEvent struct {
	BillableID  int64
	Transaction *Transaction
	Payload     map[string]any
}

// TransactionUpdatedEvent is dispatched when a transaction is updated via webhook.
type TransactionUpdatedEvent struct {
	BillableID  int64
	Transaction *Transaction
	Payload     map[string]any
}

// WebhookReceivedEvent is dispatched for every incoming provider webhook.
type WebhookReceivedEvent struct {
	Payload map[string]any
}

// WebhookHandledEvent is dispatched after a webhook has been processed.
type WebhookHandledEvent struct {
	Payload map[string]any
}

// ---------------------------------------------------------------------------
// Madora domain events
// ---------------------------------------------------------------------------

// SubscriptionChangedEvent is dispatched after a subscription transition.
type SubscriptionChangedEvent struct {
	TeamID int64
	Plan   string
	Status string
}
