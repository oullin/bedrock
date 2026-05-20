# spark

<!-- laravel-docs: cashier-paddle.md#laravel-cashier-paddle -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

The spark package provides Bedrock's Go implementation for this Laravel-aligned surface.

<div class="docs-callout docs-callout-laravel">
  <strong>Laravel baseline.</strong>
  This page follows the Laravel 13.x documentation structure for the matching feature area, then rewrites the examples and edge cases for Bedrock's Go packages.
</div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  Bedrock replaces Laravel facades, service container magic, PHP traits, and Artisan commands with explicit Go constructors, interfaces, structs, context propagation, and ordinary package tests.
</div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/spark@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/spark/...
```

## Source Coverage

| Package    | Purpose                                      |
| ---------- | -------------------------------------------- |
| `spark`    | Public spark API surface for this module.    |
| `action`   | Public action API surface for this module.   |
| `handler`  | Public handler API surface for this module.  |
| `listener` | Public listener API surface for this module. |
| `service`  | Public service API surface for this module.  |
| `state`    | Public state API surface for this module.    |
| `webhook`  | Public webhook API surface for this module.  |

## Core Concepts

The spark reference is organized around the exported Go surface for package `spark`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `ActiveSubscription`, `AuthorizerFunc`, `Billable`, `BillableConfig`, `BillableConfigBuilder`, `BillableResolver`, `BillingService`, `CancelSubscriptionHandler`, `Checkout`, `CheckoutItem`, `CheckoutReconciliationStore`, `Config`, `Customer`, `CustomerCreateOptions`, `CustomerStore`, `CustomerUpdatedEvent`, `DownloadInvoiceHandler`, `EligibilityFunc`, `EventDispatcher`, `ExpirableSubscriptionStore`, and 62 more                      |
| Constructors and functions | `Activate`, `Active`, `AddPlan`, `AddSeats`, `AppName`, `Archive`, `Authorize`, `Billable`, `BillableConfig`, `Billables`, `BootstrapCustomer`, `BrandColor`, `BrandLogo`, `Build`, `CanUpdateSeats`, `Cancel`, `CancelSubscription`, `Canceled`, `ChargePerSeat`, `ChargesPerSeat`, and 171 more                                                                                                                                                   |
| Variables                  | `ErrAlreadySubscribed`, `ErrBillableRequired`, `ErrInvalidProvider`, `ErrNotFound`, `ErrNotSubscribed`, `SupportedProviders`                                                                                                                                                                                                                                                                                                                        |
| Constants                  | `DefaultCurrency`, `DefaultPendingExpiryDays`, `DefaultSubscriptionType`, `DoNotBill`, `FullImmediately`, `FullNextBilling`, `IntervalDay`, `IntervalMonth`, `IntervalWeek`, `IntervalYear`, `OrderStatusCompleted`, `OrderStatusPending`, `PlanPricingModeCustom`, `PlanPricingModeFree`, `PlanPricingModeMoney`, `ProductTypeOneTime`, `ProductTypeSubscription`, `ProrateImmediately`, `ProrateNextBilling`, `RouteInvoiceDownload`, and 26 more |

### Capability Matrix

| Capability                  | Documentation note                                                                                                   |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Events and listeners        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/spark"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/spark` cover the supported creation paths, default values, and Laravel parity behavior.

## Configuration

Laravel documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Laravel shape     | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Laravel parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If Laravel depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/spark/...
```

Laravel parity is tracked by these tests:

- `packages/spark/madora_inventory_test.go`

## API Reference

### Exported Types

| Type                             | Notes                                                                              |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| `ActiveSubscription`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthorizerFunc`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Billable`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BillableConfig`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BillableConfigBuilder`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BillableResolver`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BillingService`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CancelSubscriptionHandler`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Checkout`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheckoutItem`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheckoutReconciliationStore`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Config`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Customer`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CustomerCreateOptions`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CustomerStore`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CustomerUpdatedEvent`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DownloadInvoiceHandler`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EligibilityFunc`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventDispatcher`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExpirableSubscriptionStore`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FeatureConfig`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Features`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FrontendState`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handler`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handlers`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvoiceDownload`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvoiceDownloader`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSubscriptionHandler`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Order`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrderStore`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaddleError`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaddleWebhookListener`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PayInvoice`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Payment`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaymentMethodUpdateTransaction` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaymentMethodUpdater`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaymentMethodsHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaymentReporter`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingCheckoutHandler`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Plan`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PlanPeriodPrice`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PlanPricingMode`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PortalHandler`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Price`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PricePreview`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Product`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProductStore`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProrationBehavior`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProviderCheckout`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProviderConfigurable`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProviderOperations`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QuantityUpdater`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolverFunc`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResumeSubscriptionHandler`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteSet`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SeatCountFunc`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SparkServiceProvider`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StripeWebhookListener`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subscription`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionBuilder`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionCanceledEvent`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionCreatedEvent`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionCreator`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionInterval`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionItem`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionItemStore`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionPausedEvent`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionStatus`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionStore`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionUpdatedEvent`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscriptionUpdater`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Transaction`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionCompletedEvent`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionStatus`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionStore`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionUpdatedEvent`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateSubscriptionHandler`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidationError`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidationErrors`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WebhookHandledEvent`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WebhookReceivedEvent`           | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                             | Notes                                                                              |
| ------------------------------------ | ---------------------------------------------------------------------------------- |
| `Activate`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Active`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddPlan`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddSeats`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AppName`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Archive`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Authorize`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Billable`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BillableConfig`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Billables`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BootstrapCustomer`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BrandColor`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BrandLogo`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Build`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CanUpdateSeats`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cancel`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CancelSubscription`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Canceled`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChargePerSeat`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChargesPerSeat`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheckPlanEligibility`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Checkout`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheckoutSessionOptions`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientSideToken`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CollectsBillingAddress`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CollectsEuVat`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Create`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateOneTimeOrder`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Current`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentAt`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CustomerCheckout`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Daily`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DashboardURL`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DateFormat`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultBillableType`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultConfig`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultConfigRepository`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisplayAmount`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DoNotBill`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Download`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Enabled`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnforcesAcceptingTerms`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnsurePlanEligibility`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Execute`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Expirable`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Expire`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExpireStaleSubscriptions`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FeatureFlags`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindItemByPrice`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FormatAmount`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FormattedAmount`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FormattedPrice`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetActiveSubscription`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GrantsAccess`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GuestCheckout`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandleCheckoutSessionCompleted`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandleTransactionCompleted`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasErrors`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasExpiredGenericTrial`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasPrice`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasProduct`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasTax`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImmediatelyWithoutProrate`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvoiceDownloadURL`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsAuthorized`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsCompleted`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsOneTime`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPending`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsSubscribedToAnyProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsSubscription`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsSupportedCurrency`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsValidProvider`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkAsCompleted`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkPastDue`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkPaymentReady`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Middleware`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Monthly`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBillingService`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCancelSubscriptionHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConfig`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConfigFromValues`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDownloadInvoiceHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFeatures`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFrontendState`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHandler`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNewSubscriptionHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPaddleWebhookListener`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPaymentMethodsHandler`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPendingCheckoutHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPendingSubscription`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPlan`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPortalHandler`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewResumeSubscriptionHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRouteRegistry`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRouteSet`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSparkServiceProvider`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStripeWebhookListener`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSubscriptionBuilder`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSubscriptionCreator`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSubscriptionUpdater`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTrialSubscription`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUpdateSubscriptionHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NoProrate`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NormalizeItems`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnGenericTrial`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnGracePeriod`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnTrial`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Option`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Options`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseMinorAmount`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PastDue`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Path`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pause`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Paused`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaymentMethodSessionOptions`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Plan`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Plans`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PriceIDForProvider`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prorate`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProrateImmediately`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prorates`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProrationBehavior`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Quantity`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RawAmount`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RawSubtotal`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RawTax`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RawTotal`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReconcileSubscriptionAfterCheckout` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Recurring`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterBillable`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterRoutes`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterWayfinderRoutes`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemoveSeats`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Repository`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolve`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveBillable`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resume`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResumeSubscription`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetainKey`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RotatePlanPrice`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sandbox`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SeatCount`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SeatName`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SellerID`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendsInvoiceEmails`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendsPaymentNotificationEmails`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCheckoutSessionOptions`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefault`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetFeatures`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetIncentive`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetInterval`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetOptions`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPaymentMethodSessionOptions`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetProrates`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetShortDescription`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetStatus`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTrialDays`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Setup`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Show`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SparkPlan`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `State`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TaxFormatted`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TermsURL`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TotalFormatted`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionCheckout`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Type`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unwrap`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Update`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateSeats`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Valid`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidPlan`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Validate`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidatePriceMoney`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidateTransactionMoney`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VerifyBillableIsSubscribed`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VerifySignature`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Weekly`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithCustomData`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithCustomerStore`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithPaymentReporter`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithReturnURL`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithRoutes`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTransactionStore`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Yearly`                             | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                             | Notes                                                                              |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| `DefaultCurrency`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultPendingExpiryDays`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultSubscriptionType`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DoNotBill`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrAlreadySubscribed`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrBillableRequired`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidProvider`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNotFound`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNotSubscribed`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FullImmediately`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FullNextBilling`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IntervalDay`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IntervalMonth`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IntervalWeek`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IntervalYear`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrderStatusCompleted`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrderStatusPending`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PlanPricingModeCustom`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PlanPricingModeFree`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PlanPricingModeMoney`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProductTypeOneTime`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProductTypeSubscription`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProrateImmediately`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProrateNextBilling`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteInvoiceDownload`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoutePendingCheckout`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoutePortal`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoutePortalForBillable`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoutePortalForType`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteState`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteSubscriptionCancel`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteSubscriptionPaymentMethod` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteSubscriptionResume`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteSubscriptionStore`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteSubscriptionUpdate`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteWayfinder`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatusActive`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatusAwaitingPayment`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatusCanceled`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatusExpired`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatusPastDue`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatusPaused`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatusPending`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatusTrialing`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportedProviders`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionBilled`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionCanceled`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionCompleted`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionDraft`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionPaid`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionPastDue`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionReady`               | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
