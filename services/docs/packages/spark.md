# billing

<!-- upstream-docs: billing.md#upstream-cashier-stripe -->
<!-- upstream-docs: billing.md#subscriptions -->
<!-- upstream-docs: billing.md#checkout -->
<!-- upstream-docs: billing.md#customers -->

Subscription billing, checkout, and entitlement management.

## Overview

The `billing` package provides subscription billing and entitlement management
inspired by Upstream Billing.

**Module:** `github.com/bedrock/packages/billing`

```bash
go get github.com/bedrock/packages/billing@latest
```

## Core Concepts

| Concept        | Description                                              |
| -------------- | -------------------------------------------------------- |
| `Billable`     | Interface that a user/team type implements               |
| `Subscription` | An active recurring subscription with items              |
| `Plan`         | A product tier (monthly, annual, per-seat)               |
| `Price`        | A billing interval and currency-specific price           |
| `Order`        | A completed payment record                               |
| `Transaction`  | An individual charge associated with an order or invoice |
| `Checkout`     | A session representing a pending payment flow            |
| `Customer`     | The payment processor's customer record                  |

## Key Operations

- Create and cancel subscriptions
- Swap plans (upgrade/downgrade)
- Add, remove, and update subscription items (per-seat billing)
- Resume cancelled subscriptions
- Issue one-off charges and retrieve invoices
- Format amounts with `billing.FormatAmount`

## Usage

```go
service := billing.NewBillingService(paymentProcessor, subscriptionRepo)

// Subscribe
sub, err := service.NewSubscription(ctx, user, "pro-monthly")

// Swap plan (upgrade/downgrade)
err = sub.Swap(ctx, "pro-annual")

// Add per-seat item
err = sub.AddItem(ctx, "extra-seat", 2)

// Cancel (at period end)
err = sub.Cancel(ctx)

// Resume a cancelled subscription
err = sub.Resume(ctx)

// One-off charge
order, err := service.Charge(ctx, user, 2999, "USD", "Add-on purchase")

// Check entitlement
if user.Subscribed(ctx) {
    // has active subscription
}
if user.OnPlan(ctx, "pro-monthly") {
    // on a specific plan
}
```

## Sub-packages

| Sub-package      | Purpose                                              |
| ---------------- | ---------------------------------------------------- |
| `billing/handler`  | HTTP webhook handler for payment processor events    |
| `billing/listener` | Default event listeners wired to webhook events      |
| `billing/action`   | Discrete billing actions (subscribe, swap, cancel)   |
| `billing/service`  | Service layer wrapping the payment processor SDK     |
| `billing/state`    | Subscription state machine constants and transitions |
| `billing/webhook`  | Webhook signature verification and payload parsing   |
