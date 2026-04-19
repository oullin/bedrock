# notifications

Multi-channel notification delivery.

## Overview

The `notifications` package provides a `Manager` (analogous to Laravel's
`NotificationSender`) that routes notifications across channels: mail, database,
broadcast, and custom drivers. Notifications can be sent synchronously or queued.

**Module:** `github.com/gocanto/bedrock/packages/notifications`

```bash
go get github.com/gocanto/bedrock/packages/notifications@latest
```

## Built-in Channels

| Channel     | Delivery mechanism                        |
| ----------- | ----------------------------------------- |
| `mail`      | Sends via the `mailx` package             |
| `database`  | Persists to a `DatabaseNotificationStore` |
| `broadcast` | Publishes to a broadcast backend          |

## Defining a Notification

```go
type InvoicePaid struct {
    InvoiceID string
    Amount    int
}

// Via declares which channels to use
func (n *InvoicePaid) Via(ctx context.Context, notifiable cn.Notifiable) []string {
    return []string{"mail", "database"}
}

// ToMail builds the mail message
func (n *InvoicePaid) ToMail(ctx context.Context, notifiable cn.Notifiable) (*notifications.MailMessage, error) {
    return notifications.NewMailMessage().
        Subject("Invoice Paid").
        Line(fmt.Sprintf("Invoice #%s for $%d has been paid.", n.InvoiceID, n.Amount)).
        Action("View Invoice", "https://example.com/invoices/"+n.InvoiceID), nil
}

// ToDatabase builds the database payload
func (n *InvoicePaid) ToDatabase(ctx context.Context, notifiable cn.Notifiable) (map[string]any, error) {
    return map[string]any{"invoice_id": n.InvoiceID, "amount": n.Amount}, nil
}
```

## Sending Notifications

```go
sender := notifications.NewSender(manager, events.NewDispatcher())

sender.Send(ctx, user, &InvoicePaid{InvoiceID: "INV-001", Amount: 99})

// Multiple notifiables at once
sender.SendToMany(ctx, []cn.Notifiable{user1, user2}, &InvoicePaid{...})
```

## Anonymous Notifiable

Send to an address without a user entity:

```go
notifications.Route("mail", "admin@example.com").
    Notify(ctx, sender, &InvoicePaid{...})
```

## Queued Notifications

Implement `notifications.ShouldQueue` to push delivery to the queue backend:

```go
func (n *InvoicePaid) ShouldQueue() {}
```

## Conditional Sending

Implement `ShouldSendNotification` to suppress delivery for specific channels:

```go
func (n *InvoicePaid) ShouldSend(ctx context.Context, notifiable cn.Notifiable, channel string) bool {
    return channel != "database" // only send via mail
}
```

## Database Channel

Persist and query notification records via `DatabaseNotificationStore`:

```go
store.ForNotifiable(ctx, "User", userID)
store.MarkAsRead(ctx, notificationID)
store.UnreadForNotifiable(ctx, "User", userID)
```
