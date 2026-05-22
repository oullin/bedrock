package billing

// Event types mirror the upstream event classes dispatched during billing
// lifecycle operations.

// SubscriptionCreatedEvent is dispatched when a subscription is created.
type SubscriptionCreatedEvent struct {
	Billable     Billable
	Subscription *Subscription
	Payload      map[string]any
}

// SubscriptionUpdatedEvent is dispatched when a subscription is updated.
type SubscriptionUpdatedEvent struct {
	Subscription *Subscription
	Payload      map[string]any
}

// SubscriptionCanceledEvent is dispatched when a subscription is canceled.
type SubscriptionCanceledEvent struct {
	Subscription *Subscription
	Payload      map[string]any
}

// SubscriptionPausedEvent is dispatched when a subscription is paused.
type SubscriptionPausedEvent struct {
	Subscription *Subscription
	Payload      map[string]any
}

// CustomerUpdatedEvent is dispatched when a customer record is updated.
type CustomerUpdatedEvent struct {
	Billable Billable
	Customer *Customer
	Payload  map[string]any
}

// TransactionCompletedEvent is dispatched when a transaction completes.
type TransactionCompletedEvent struct {
	Billable    Billable
	Transaction *Transaction
	Payload     map[string]any
}

// TransactionUpdatedEvent is dispatched when a transaction is updated.
type TransactionUpdatedEvent struct {
	Transaction *Transaction
	Payload     map[string]any
}

// WebhookReceivedEvent is dispatched when any webhook payload arrives.
type WebhookReceivedEvent struct {
	Payload map[string]any
}

// WebhookHandledEvent is dispatched after a webhook has been processed.
type WebhookHandledEvent struct {
	Payload map[string]any
}

// EventDispatcher dispatches billing events to listeners.
type EventDispatcher interface {
	Dispatch(event any) error
}
