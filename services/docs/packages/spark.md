# spark

Subscription billing, checkout, and entitlement management.

## Overview

The `spark` package provides subscription billing and entitlement management
inspired by Laravel Spark.

**Module:** `github.com/gocanto/bedrock/packages/spark`

```bash
go get github.com/gocanto/bedrock/packages/spark@latest
```

## Core Concepts

| Concept        | Description                                              |
|----------------|----------------------------------------------------------|
| `Billable`     | Interface that a user/team type implements               |
| `Subscription` | An active recurring subscription with items             |
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
- Format amounts with `spark.FormatAmount`

## Sub-packages

| Sub-package      | Purpose                                               |
|------------------|-------------------------------------------------------|
| `spark/handler`  | HTTP webhook handler for payment processor events     |
| `spark/listener` | Default event listeners wired to webhook events       |
| `spark/action`   | Discrete billing actions (subscribe, swap, cancel)    |
| `spark/service`  | Service layer wrapping the payment processor SDK      |
| `spark/state`    | Subscription state machine constants and transitions  |
| `spark/webhook`  | Webhook signature verification and payload parsing    |
