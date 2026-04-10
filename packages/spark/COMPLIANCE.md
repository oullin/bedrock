# Spark Port — Test Compliance Report

> Generated: 2026-04-07
> Baseline: Madora PHP billing tests (`tests/{Feature,Unit}/Billing/`)
> PHP sources: laravel/cashier-paddle, laravel/spark-paddle, Madora app/Billing
> All Bedrock tests pass (`go test ./...` green across all packages).

---

## Executive Summary

| Metric                                           | Count           |
| ------------------------------------------------ | --------------- |
| PHP source files ported                          | 88 / 95 (92.6%) |
| PHP source files not ported (framework plumbing) | 7               |
| **Go source files**                              | **57**          |
| Go test files                                    | 6               |
| **Bedrock Go tests (total)**                     | **26**          |
| PHP test files in scope                          | **30**          |
| PHP test files analyzed line-by-line             | 30              |
| PHP test methods analyzed                        | 162             |
| Analyzed: COVERED                                | 12              |
| Analyzed: MISSING                                | 122             |
| Analyzed: INTENTIONAL-SKIP                       | 28              |
| Analyzed: BEDROCK-ONLY                           | 14              |
| **Coverage of analyzed tests (excl. skips)**     | **9%**          |

### What This Means

Bedrock has **26 passing Go tests** across 6 files in the root `spark` package. All 9 subdirectory packages have **zero test coverage**. For the 162 PHP test methods analyzed line-by-line, we cover 9% of the portable surface. The 122 missing methods span critical business logic: subscription lifecycle, checkout reconciliation, entitlement gating, webhook verification, and HTTP handlers.

---

## Per-Area Scorecard

### Unit Tests (9 PHP files, 39 methods)

| PHP Test File                            | PHP Tests | Covered | Missing | Skip  | Bedrock-Only | Portable %          |
| ---------------------------------------- | --------- | ------- | ------- | ----- | ------------ | ------------------- |
| **SubscriptionStatusEvaluatorTest**      | 9         | 9       | 0       | 0     | 0            | **100%**            |
| **PlanPeriodPriceTest**                  | 3         | 2       | 1       | 0     | 0            | **67%**             |
| **BillingPeriodTest**                    | 3         | 1       | 0       | 2     | 0            | **100%** (portable) |
| **SubscriptionGateTest**                 | 7         | 0       | 7       | 0     | 0            | **0%**              |
| **PaddleWebhookSignatureInspectorTest**  | 4         | 0       | 4       | 0     | 0            | **0%**              |
| **TransactionSnapshotTest**              | 3         | 0       | 3       | 0     | 0            | **0%**              |
| **SubscriptionEntitlementServiceTest**   | 2         | 0       | 2       | 0     | 0            | **0%**              |
| **BillingLocalDevUrlTest**               | 5         | 0       | 5       | 0     | 0            | **0%**              |
| **PaddleWebhookSignatureMiddlewareTest** | 1         | 0       | 1       | 0     | 0            | **0%**              |
| **Subtotal**                             | **37**    | **12**  | **23**  | **2** | **0**        | **34%**             |

### Feature Tests (21 PHP files, ~125 methods)

| PHP Test File                                | PHP Tests | Covered | Missing | Skip   | Bedrock-Only | Portable % |
| -------------------------------------------- | --------- | ------- | ------- | ------ | ------------ | ---------- |
| **BillingLifecycleTest**                     | 28        | 0       | 28      | 0      | 0            | **0%**     |
| **BillingControllerTest**                    | 19        | 0       | 19      | 0      | 0            | **0%**     |
| **SubscriptionStatusTransitionTest**         | 12        | 0       | 12      | 0      | 0            | **0%**     |
| **ReconcileSubscriptionsJobTest**            | 10        | 0       | 10      | 0      | 0            | **0%**     |
| **CustomPlanInquiryTest**                    | 7         | 0       | 7       | 0      | 0            | **0%**     |
| **CanonicalMoneySchemaContractTest**         | 7         | 0       | 0       | 7      | 0            | N/A (skip) |
| **ReconcileSubscriptionAfterCheckoutTest**   | 6         | 0       | 6       | 0      | 0            | **0%**     |
| **NormaliseBillableRouteParameterTest**      | 5         | 0       | 5       | 0      | 0            | **0%**     |
| **PlanPeriodPriceSchemaContractTest**        | 5         | 0       | 0       | 5      | 0            | N/A (skip) |
| **ExpirePendingSubscriptionsCommandTest**    | 4         | 0       | 4       | 0      | 0            | **0%**     |
| **SparkPlanRegistrarTest**                   | 3         | 0       | 3       | 0      | 0            | **0%**     |
| **SparkPortalTest**                          | 3         | 0       | 3       | 0      | 0            | **0%**     |
| **SyncSubscriptionAfterUpdateTest**          | 3         | 0       | 3       | 0      | 0            | **0%**     |
| **PlanPeriodPriceManagerTest**               | 3         | 0       | 3       | 0      | 0            | **0%**     |
| **ReconcileSubscriptionPlanMatchingTest**    | 2         | 0       | 2       | 0      | 0            | **0%**     |
| **InspectPaddleWebhookSignatureCommandTest** | 2         | 0       | 2       | 0      | 0            | **0%**     |
| **PaddleTransactionModelTest**               | 2         | 0       | 0       | 2      | 0            | N/A (skip) |
| **CurrencyForeignKeysRestrictDeleteTest**    | 1         | 0       | 0       | 1      | 0            | N/A (skip) |
| **ReconcileSubscriptionsCommandTest**        | 1         | 0       | 1       | 0      | 0            | **0%**     |
| **RecoverExplicitSubscriptionCommandTest**   | 1         | 0       | 1       | 0      | 0            | **0%**     |
| **ReferenceDataSeederTest**                  | 1         | 0       | 0       | 1      | 0            | N/A (skip) |
| **Subtotal**                                 | **~125**  | **0**   | **109** | **16** | **0**        | **0%**     |

### Bedrock-Only Tests (no PHP equivalent)

| Go Test                               | Tests  | What it covers                                                              |
| ------------------------------------- | ------ | --------------------------------------------------------------------------- |
| `TestDefaultConfig`                   | 1      | Default config values (Path, DashboardURL, Prorates, Currency, WebhookPath) |
| `TestCustomerOnGenericTrial`          | 1      | Customer trial period detection                                             |
| `TestCustomerHasExpiredGenericTrial`  | 1      | Customer trial expiry detection                                             |
| `TestSubscriptionPlanValid`           | 1      | SubscriptionPlan enum validation                                            |
| `TestSubscriptionStatusValid`         | 1      | SubscriptionStatus enum validation                                          |
| `TestSubscriptionStatusIsTerminal`    | 1      | Terminal status detection (Canceled, Expired)                               |
| `TestPlanPricingModeDisplayLabel`     | 1      | PricingMode display labels                                                  |
| `TestSubscriptionFeatureCodeValid`    | 1      | FeatureCode enum validation                                                 |
| `TestProrationBehaviorValid`          | 1      | ProrationBehavior enum validation                                           |
| `TestTransactionStatusValid`          | 1      | TransactionStatus enum validation                                           |
| `TestSubscriptionIsNew`               | 1      | Zero-ID = new subscription                                                  |
| `TestSubscriptionOnTrial`             | 1      | Subscription trial detection                                                |
| `TestSubscriptionOnGracePeriod`       | 1      | Grace period detection                                                      |
| `TestSubscriptionOnPausedGracePeriod` | 1      | Paused grace period detection                                               |
| `TestSubscriptionHasProduct`          | 1      | Product lookup by ID                                                        |
| `TestSubscriptionHasPrice`            | 1      | Price lookup by ID                                                          |
| `TestSubscriptionProration`           | 1      | Proration behavior configuration                                            |
| `TestManagerBillableRegistration`     | 1      | Billable type registration and resolution                                   |
| `TestManagerPlans`                    | 1      | Plan management                                                             |
| `TestManagerDefaultBillableType`      | 1      | Default billable type selection                                             |
| `TestManagerPerSeatBilling`           | 1      | Per-seat billing configuration                                              |
| **Subtotal**                          | **14** | Go-specific domain model tests                                              |

---

## Covered Tests — Detail

### SubscriptionStatusEvaluatorTest.php → `enum_test.go` + `subscription_test.go` (9/9 = 100%)

| PHP Method                                                       | Go Test                              | Status  |
| ---------------------------------------------------------------- | ------------------------------------ | ------- |
| `test_active_subscription_grants_access`                         | `TestSubscriptionStatusGrantsAccess` | COVERED |
| `test_trialing_subscription_grants_access`                       | `TestSubscriptionStatusGrantsAccess` | COVERED |
| `test_past_due_subscription_grants_access`                       | `TestSubscriptionStatusGrantsAccess` | COVERED |
| `test_paused_subscription_denies_access`                         | `TestSubscriptionStatusGrantsAccess` | COVERED |
| `test_canceled_subscription_denies_access`                       | `TestSubscriptionStatusGrantsAccess` | COVERED |
| `test_expired_subscription_denies_access`                        | `TestSubscriptionStatusGrantsAccess` | COVERED |
| `test_pending_subscription_denies_access`                        | `TestSubscriptionStatusGrantsAccess` | COVERED |
| `test_active_with_expired_pending_expires_at_denies_access`      | `TestSubscriptionValid`              | COVERED |
| `test_statuses_granting_access_returns_active_trialing_past_due` | `TestSubscriptionStatusGrantsAccess` | COVERED |

### PlanPeriodPriceTest.php → `plan_test.go` (2/3 = 67%)

| PHP Method                                                              | Go Test                            | Status  |
| ----------------------------------------------------------------------- | ---------------------------------- | ------- |
| `test_display_amount_formats_valid_money_prices`                        | `TestPlanPeriodPriceDisplayAmount` | COVERED |
| `test_display_amount_returns_literal_labels_for_free_and_custom_prices` | `TestPlanPeriodPriceDisplayAmount` | COVERED |
| `test_display_amount_throws_for_incomplete_money_prices`                | —                                  | MISSING |

### BillingPeriodTest.php → `enum_test.go` (1/1 portable = 100%)

| PHP Method                                                                          | Go Test                         | Status                      |
| ----------------------------------------------------------------------------------- | ------------------------------- | --------------------------- |
| `test_validation_rule_accepts_known_periods_and_rejects_unknown_values`             | `TestBillingPeriodDisplayLabel` | COVERED (partial)           |
| `test_model_casts_hydrate_billing_period_enum_instances`                            | —                               | INTENTIONAL-SKIP (Eloquent) |
| `test_pluck_returns_billing_period_enum_instances_from_hydrated_plan_period_models` | —                               | INTENTIONAL-SKIP (Eloquent) |

---

## Missing Tests — By Go Package

### `entitlement/` (0 tests — 9 PHP methods missing)

**Needs: `entitlement/gate_test.go`, `entitlement/access_test.go`**

| PHP Source                         | PHP Method                                                                 | Go Target                                      |
| ---------------------------------- | -------------------------------------------------------------------------- | ---------------------------------------------- |
| SubscriptionGateTest               | `test_no_subscription_returns_payment_required`                            | `entitlement/gate.go` — Gate.CheckRequirements |
| SubscriptionGateTest               | `test_matching_plan_returns_granted`                                       | `entitlement/gate.go` — Gate.CheckRequirements |
| SubscriptionGateTest               | `test_wrong_plan_returns_forbidden`                                        | `entitlement/gate.go` — Gate.CheckRequirements |
| SubscriptionGateTest               | `test_no_plan_filter_allows_any_active_subscription`                       | `entitlement/gate.go` — Gate.CheckRequirements |
| SubscriptionGateTest               | `test_inactive_subscription_returns_payment_required`                      | `entitlement/gate.go` — Gate.CheckRequirements |
| SubscriptionGateTest               | `test_past_due_subscription_grants_access`                                 | `entitlement/gate.go` — Gate.CheckRequirements |
| SubscriptionGateTest               | `test_paused_subscription_returns_payment_required`                        | `entitlement/gate.go` — Gate.CheckRequirements |
| SubscriptionEntitlementServiceTest | `test_active_subscription_features_are_provisioned_from_the_selected_plan` | `entitlement/access.go` — Access               |
| SubscriptionEntitlementServiceTest | `test_expired_subscription_deprovisions_features`                          | `entitlement/access.go` — Access               |

### `webhook/` (0 tests — 5 PHP methods missing)

**Needs: `webhook/signature_test.go`**

| PHP Source                           | PHP Method                                                                 | Go Target                                           |
| ------------------------------------ | -------------------------------------------------------------------------- | --------------------------------------------------- |
| PaddleWebhookSignatureInspectorTest  | `test_inspector_accepts_a_valid_signature`                                 | `webhook/signature.go` — VerifySignature            |
| PaddleWebhookSignatureInspectorTest  | `test_inspector_reports_a_mismatch_for_a_truncated_secret`                 | `webhook/signature.go` — VerifySignature            |
| PaddleWebhookSignatureInspectorTest  | `test_inspector_reports_a_malformed_signature_header`                      | `webhook/signature.go` — VerifySignature            |
| PaddleWebhookSignatureInspectorTest  | `test_inspector_reports_when_the_timestamp_is_outside_the_variance_window` | `webhook/signature.go` — VerifySignature            |
| PaddleWebhookSignatureMiddlewareTest | `test_middleware_accepts_a_valid_signature`                                | `webhook/signature.go` — VerifySignature middleware |

### `subscription/` (0 tests — 40 PHP methods missing)

**Needs: `subscription/transitioner_test.go`, `subscription/repository_test.go`, `subscription/state_test.go`**

| PHP Source                              | Methods | Go Target                                                     |
| --------------------------------------- | ------- | ------------------------------------------------------------- |
| SubscriptionStatusTransitionTest        | 12      | `subscription/transitioner.go` — Transitioner state machine   |
| BillingLifecycleTest (lifecycle subset) | ~20     | `subscription/repository.go` — Repository + state transitions |
| BillingLifecycleTest (state DTO subset) | ~8      | `subscription/state.go` — StateResolver billing snapshots     |

### `checkout/` (0 tests — 19 PHP methods missing)

**Needs: `checkout/reconciler_test.go`, `checkout/syncer_test.go`, `checkout/recovery_test.go`**

| PHP Source                             | Methods | Go Target                                              |
| -------------------------------------- | ------- | ------------------------------------------------------ |
| ReconcileSubscriptionsJobTest          | 10      | `checkout/reconciler.go` — Reconciler                  |
| ReconcileSubscriptionAfterCheckoutTest | 6       | `checkout/reconciler.go` — Reconciler.ReconcileCreated |
| SyncSubscriptionAfterUpdateTest        | 3       | `checkout/syncer.go` — Syncer.SyncUpdated              |

### `handler/` (0 tests — 31 PHP methods missing)

**Needs: `handler/billing_test.go`, `handler/inquiry_test.go`, `handler/middleware_test.go`, `handler/portal_test.go`**

| PHP Source                          | Methods | Go Target                                             |
| ----------------------------------- | ------- | ----------------------------------------------------- |
| BillingControllerTest               | 19      | `handler/billing.go` — BillingHandler                 |
| CustomPlanInquiryTest               | 7       | `handler/inquiry.go` — InquiryHandler                 |
| NormaliseBillableRouteParameterTest | 5       | `handler/middleware.go` — NormaliseBillableRouteParam |

### `catalog/` (0 tests — 6 PHP methods missing)

**Needs: `catalog/presentation_test.go`, `catalog/price_manager_test.go`**

| PHP Source                 | Methods | Go Target                                           |
| -------------------------- | ------- | --------------------------------------------------- |
| SparkPlanRegistrarTest     | 3       | `catalog/presentation.go` — SparkPlanRegistry       |
| PlanPeriodPriceManagerTest | 3       | `catalog/price_manager.go` — PlanPeriodPriceManager |

### `command/` (0 tests — 5 PHP methods missing)

**Needs: `command/command_test.go`**

| PHP Source                             | Methods | Go Target                                          |
| -------------------------------------- | ------- | -------------------------------------------------- |
| ExpirePendingSubscriptionsCommandTest  | 4       | `command/command.go` — ExpirePendingSubscriptions  |
| RecoverExplicitSubscriptionCommandTest | 1       | `command/command.go` — RecoverExplicitSubscription |

### Root package — remaining gaps (9 PHP methods missing)

| PHP Source              | PHP Method                                                                      | Go Target                                            |
| ----------------------- | ------------------------------------------------------------------------------- | ---------------------------------------------------- |
| TransactionSnapshotTest | `test_empty_state_returns_no_transaction_placeholder_values`                    | `snapshot.go` — EmptyTransactionSnapshot             |
| TransactionSnapshotTest | `test_from_transaction_marks_money_fields_empty_when_transaction_is_incomplete` | `snapshot.go` — TransactionSnapshot                  |
| TransactionSnapshotTest | `test_from_transaction_preserves_money_fields_when_transaction_is_complete`     | `snapshot.go` — TransactionSnapshot                  |
| BillingLocalDevUrlTest  | `test_it_returns_the_configured_local_dev_url`                                  | URL resolution (config or URLResolver)               |
| BillingLocalDevUrlTest  | `test_it_falls_back_to_the_app_url_when_the_local_dev_url_is_blank`             | URL resolution                                       |
| BillingLocalDevUrlTest  | `test_it_builds_absolute_urls_from_relative_paths`                              | URL resolution                                       |
| BillingLocalDevUrlTest  | `test_it_preserves_absolute_urls`                                               | URL resolution                                       |
| BillingLocalDevUrlTest  | `test_it_sets_cashier_webhook_from_the_local_dev_url_when_missing`              | URL resolution                                       |
| PlanPeriodPriceTest     | `test_display_amount_throws_for_incomplete_money_prices`                        | `plan.go` — PlanPeriodPrice.DisplayAmount error path |

---

## Intentional Skips (28 methods)

These PHP tests verify Laravel/Eloquent framework behavior with no Go equivalent:

| PHP Test File                                   | Methods | Reason                                                        |
| ----------------------------------------------- | ------- | ------------------------------------------------------------- |
| CanonicalMoneySchemaContractTest                | 7       | Laravel database schema validation (migrations, column types) |
| PlanPeriodPriceSchemaContractTest               | 5       | Laravel database unique constraints and schema enforcement    |
| BillingPeriodTest                               | 2       | Eloquent model casting and pluck (PHP-specific ORM behavior)  |
| PaddleTransactionModelTest                      | 2       | Eloquent model event hooks (creating/updating observers)      |
| CurrencyForeignKeysRestrictDeleteTest           | 1       | DB foreign key cascade enforcement                            |
| ReferenceDataSeederTest                         | 1       | Laravel DB seeder infrastructure                              |
| ReconcileSubscriptionsCommandTest               | 1       | Artisan command description metadata                          |
| SparkPortalTest (partial)                       | 3       | Inertia.js view rendering, Blade markup                       |
| InspectPaddleWebhookSignatureCommandTest        | 2       | Artisan command I/O testing                                   |
| ReconcileSubscriptionPlanMatchingTest (partial) | 2       | Tests rely on Eloquent model factories and DB state           |

---

## Gap Categories

- **(A) Language difference** — PHP/Laravel-specific constructs with no Go equivalent. Permanent intentional skips. (28 methods)
- **(B) Missing test coverage** — Portable features implemented in Go but without test coverage. Action required. (122 methods)
- **(C) Partial coverage** — Feature exists with some tests but incomplete. Action required. (1 method — `PlanPeriodPriceTest`)

---

## Per-Package Coverage

| Package         | Go Source Files | Go Test Files | Go Tests | PHP Methods (portable) | Coverage     |
| --------------- | --------------- | ------------- | -------- | ---------------------- | ------------ |
| `spark` (root)  | 26              | 6             | 26       | 22                     | **55%**      |
| `billing/`      | 2               | 0             | 0        | 0                      | N/A (facade) |
| `catalog/`      | 4               | 0             | 0        | 6                      | **0%**       |
| `checkout/`     | 4               | 0             | 0        | 19                     | **0%**       |
| `command/`      | 2               | 0             | 0        | 5                      | **0%**       |
| `entitlement/`  | 3               | 0             | 0        | 9                      | **0%**       |
| `handler/`      | 8               | 0             | 0        | 31                     | **0%**       |
| `listener/`     | 1               | 0             | 0        | 0                      | N/A (wiring) |
| `subscription/` | 4               | 0             | 0        | 40                     | **0%**       |
| `webhook/`      | 2               | 0             | 0        | 5                      | **0%**       |
| **Total**       | **57**          | **6**         | **26**   | **134**                | **9%**       |

---

## PHP to Go Mapping — laravel/cashier-paddle

| PHP File                                   | Go File                     | Tested |
| ------------------------------------------ | --------------------------- | ------ |
| Billable.php                               | `billable.go`               | No     |
| Cashier.php                                | `manager.go`                | Yes    |
| CashierFake.php                            | _Not ported_                | —      |
| CashierServiceProvider.php                 | _Not ported_                | —      |
| Checkout.php                               | `checkout_session.go`       | No     |
| Components/Button.php                      | _Not applicable_            | —      |
| Components/Checkout.php                    | _Not applicable_            | —      |
| Concerns/ManagesCustomer.php               | `customer.go`               | Yes    |
| Concerns/ManagesSubscriptions.php          | `subscription.go`           | Yes    |
| Concerns/ManagesTransactions.php           | `transaction.go`            | No     |
| Concerns/PerformsCharges.php               | `checkout_session.go`       | No     |
| Concerns/Prorates.php                      | `subscription.go` (methods) | Yes    |
| Customer.php                               | `customer.go`               | Yes    |
| Events/CustomerUpdated.php                 | `event.go`                  | No     |
| Events/SubscriptionCanceled.php            | `event.go`                  | No     |
| Events/SubscriptionCreated.php             | `event.go`                  | No     |
| Events/SubscriptionPaused.php              | `event.go`                  | No     |
| Events/SubscriptionUpdated.php             | `event.go`                  | No     |
| Events/TransactionCompleted.php            | `event.go`                  | No     |
| Events/TransactionUpdated.php              | `event.go`                  | No     |
| Events/WebhookHandled.php                  | `event.go`                  | No     |
| Events/WebhookReceived.php                 | `event.go`                  | No     |
| Exceptions/PaddleException.php             | `errors.go`                 | No     |
| Http/Controllers/WebhookController.php     | `webhook/handler.go`        | No     |
| Http/Middleware/VerifyWebhookSignature.php | `webhook/signature.go`      | No     |
| Payment.php                                | `payment.go`                | No     |
| Price.php                                  | `price.go`                  | No     |
| PricePreview.php                           | `price_preview.go`          | No     |
| Subscription.php                           | `subscription.go`           | Yes    |
| SubscriptionBuilder.php                    | `subscription_builder.go`   | No     |
| SubscriptionItem.php                       | `subscription_item.go`      | No     |
| Transaction.php                            | `transaction.go`            | No     |

## PHP to Go Mapping — laravel/spark-paddle

| PHP File                                           | Go File                                   | Tested |
| -------------------------------------------------- | ----------------------------------------- | ------ |
| Actions/GenerateCheckoutSession.php                | `provider.go` (ProviderCheckoutGenerator) | No     |
| Billable.php                                       | `billable.go`                             | No     |
| BillableConfigurationBuilder.php                   | `manager.go` (BillableConfigBuilder)      | Yes    |
| Console/InstallCommand.php                         | _Not ported_                              | —      |
| Contracts/Actions/GeneratesCheckoutSessions.php    | `provider.go`                             | No     |
| FrontendState.php                                  | `subscription/frontend.go`                | No     |
| GuessesBillableTypes.php                           | `manager.go` (DefaultBillableType)        | Yes    |
| Http/Controllers/BillingPortalController.php       | `handler/portal.go`                       | No     |
| Http/Controllers/CancelSubscriptionController.php  | `handler/subscription.go`                 | No     |
| Http/Controllers/DownloadInvoiceController.php     | `handler/invoice.go`                      | No     |
| Http/Controllers/NewPendingCheckoutController.php  | `handler/payment.go`                      | No     |
| Http/Controllers/NewSubscriptionController.php     | `handler/subscription.go`                 | No     |
| Http/Controllers/ResumeSubscriptionController.php  | `handler/subscription.go`                 | No     |
| Http/Controllers/RetrievesBillableModels.php       | _Not ported (trait)_                      | —      |
| Http/Controllers/UpdatePaymentMethodController.php | `handler/payment.go`                      | No     |
| Http/Controllers/UpdateSubscriptionController.php  | `handler/subscription.go`                 | No     |
| Http/Middleware/HandleInertiaRequests.php          | _Not applicable_                          | —      |
| Http/Middleware/VerifyBillableIsSubscribed.php     | `handler/middleware.go`                   | No     |
| Listeners/SubscriptionCreatedListener.php          | `listener/listener.go`                    | No     |
| Plan.php                                           | `plan.go`                                 | Yes    |
| Spark.php                                          | `manager.go`                              | Yes    |
| SparkManager.php                                   | `manager.go`                              | Yes    |
| SparkServiceProvider.php                           | _Not ported_                              | —      |
| ValidPlan.php                                      | `handler/input.go`                        | No     |

## PHP to Go Mapping — Madora app/Billing

| PHP File                                            | Go File                                     | Tested |
| --------------------------------------------------- | ------------------------------------------- | ------ |
| Actions/GenerateSparkCheckoutSession.php            | `provider.go` (ProviderCheckoutGenerator)   | No     |
| BillingConstants.php                                | `config.go`                                 | Yes    |
| Catalog/PlanCatalog.php                             | `catalog/catalog.go`                        | No     |
| Catalog/SparkPlanRegistry.php                       | `catalog/presentation.go`                   | No     |
| Checkout/CheckoutStarter.php                        | `checkout/starter.go`                       | No     |
| Checkout/ExplicitSubscriptionRecovery.php           | `checkout/recovery.go`                      | No     |
| Checkout/PaddleSubscriptionFetcher.php              | `provider.go` (ProviderSubscriptionFetcher) | No     |
| Checkout/SubscriptionReconciler.php                 | `checkout/reconciler.go`                    | No     |
| Checkout/SubscriptionUpdateSyncer.php               | `checkout/syncer.go`                        | No     |
| Commands/ExpirePendingSubscriptionsCommand.php      | `command/command.go`                        | No     |
| Commands/ReconcileSubscriptionsCommand.php          | `command/command.go`                        | No     |
| Commands/RecoverExplicitSubscriptionCommand.php     | `command/command.go`                        | No     |
| Contracts/BillingLocalDevUrl.php                    | `contract.go` (URLResolver)                 | No     |
| Controllers/BillingController.php                   | `handler/billing.go`                        | No     |
| Controllers/CustomPlanInquiryController.php         | `handler/inquiry.go`                        | No     |
| Data/Read/BillingPageData.php                       | `view.go` (BillingPageView)                 | No     |
| Data/Read/BillingPriceData.php                      | `view.go` (BillingPriceView)                | No     |
| Data/Read/BillingSelectedPlanData.php               | `view.go` (BillingSelectedPlanView)         | No     |
| Data/Read/BillingStateData.php                      | `view.go` (BillingStateView)                | No     |
| Data/Read/BillingTeamData.php                       | `view.go` (BillingTeamView)                 | No     |
| Data/Read/FrontendPlanData.php                      | `view.go` (FrontendPlanView)                | No     |
| Data/Read/LandingPlanData.php                       | `view.go` (LandingPlanView)                 | No     |
| Data/Read/PlanFeatureData.php                       | `view.go` (PlanFeatureView)                 | No     |
| Data/Read/PlanPeriodData.php                        | `view.go` (PlanPeriodView)                  | No     |
| Data/Read/SubscriptionCtaData.php                   | `view.go` (SubscriptionCTAView)             | No     |
| Data/Read/SubscriptionStateData.php                 | `view.go` (SubscriptionStateView)           | No     |
| Data/Read/TransactionData.php                       | `view.go` (TransactionView)                 | No     |
| Data/Result/SubscriptionDecision.php                | `entitlement/gate.go` (Decision)            | No     |
| Data/Result/SubscriptionReconciliationData.php      | `checkout/reconciler.go`                    | No     |
| Data/State/BillingStateSnapshot.php                 | `snapshot.go` (BillingStateSnapshot)        | No     |
| Data/State/PlanConfig.php                           | `snapshot.go` (PlanConfig)                  | No     |
| Data/State/SubscriptionCtaSnapshot.php              | `snapshot.go` (SubscriptionCTASnapshot)     | No     |
| Data/State/SubscriptionStateSnapshot.php            | `snapshot.go` (SubscriptionStateSnapshot)   | No     |
| Data/State/TransactionSnapshot.php                  | `snapshot.go` (TransactionSnapshot)         | No     |
| Entitlements/EntitlementAccess.php                  | `entitlement/access.go`                     | No     |
| Entitlements/EntitlementSynchronizer.php            | `entitlement/sync.go`                       | No     |
| Entitlements/SubscriptionStatusEvaluator.php        | `enum.go` (SubscriptionStatus methods)      | Yes    |
| Enums/BillingPeriod.php                             | `enum.go`                                   | Yes    |
| Enums/PlanPricingMode.php                           | `enum.go`                                   | Yes    |
| Enums/SubscriptionFeatureCode.php                   | `enum.go`                                   | Yes    |
| Enums/SubscriptionPlan.php                          | `enum.go`                                   | Yes    |
| Enums/SubscriptionStatus.php                        | `enum.go`                                   | Yes    |
| Events/SubscriptionChanged.php                      | `event.go` (SubscriptionChangedEvent)       | No     |
| Exceptions/InvalidPlanPriceException.php            | `errors.go` (InvalidPlanPriceError)         | No     |
| Http/Middleware/NormaliseBillableRouteParameter.php | `handler/middleware.go`                     | No     |
| Jobs/ReconcileSubscriptions.php                     | `command/job.go`                            | No     |
| Lifecycle/BillingStateResolver.php                  | `subscription/state.go`                     | No     |
| Lifecycle/SubscriptionRepository.php                | `subscription/repository.go`                | No     |
| Listeners/ReconcileSubscriptionAfterCheckout.php    | `listener/listener.go`                      | No     |
| Listeners/SyncSubscriptionAfterUpdate.php           | `listener/listener.go`                      | No     |
| Listeners/SyncSubscriptionEntitlements.php          | `listener/listener.go`                      | No     |
| Mail/CustomPlanInquiryConfirmationMail.php          | `mail.go`                                   | No     |
| Mail/CustomPlanInquiryMail.php                      | `mail.go`                                   | No     |
| Models/Currency.php                                 | `currency.go`                               | No     |
| Models/Feature.php                                  | `plan.go` (Feature struct)                  | No     |
| Models/Invoice.php                                  | `invoice.go`                                | No     |
| Models/Plan.php                                     | `plan.go`                                   | Yes    |
| Models/PlanFeature.php                              | `plan.go` (PlanFeature struct)              | No     |
| Models/PlanPeriod.php                               | `plan.go` (PlanPeriod struct)               | No     |
| Models/PlanPeriodPrice.php                          | `plan.go` (PlanPeriodPrice struct)          | Yes    |
| Models/Subscription.php                             | `subscription.go`                           | Yes    |
| Models/SubscriptionFeature.php                      | `subscription_feature.go`                   | No     |
| Models/SubscriptionItem.php                         | `subscription_item.go`                      | No     |
| Requests/CheckoutRequest.php                        | `handler/input.go` (CheckoutInput)          | No     |
| Requests/CustomPlanInquiryRequest.php               | `handler/input.go` (InquiryInput)           | No     |
| Support/PlanPeriodPriceManager.php                  | `catalog/price_manager.go`                  | No     |
| Support/PlanPresentation.php                        | `catalog/presentation.go`                   | No     |
| Support/PlanSeeder.php                              | `catalog/seeder.go`                         | No     |
| Support/SerialisesBillingPrices.php                 | `view.go` (BillingPriceView)                | No     |
| Support/SubscriptionGate.php                        | `entitlement/gate.go`                       | No     |
| Support/SubscriptionTransitioner.php                | `subscription/transitioner.go`              | No     |
| Workflows/BillingWorkflow.php                       | `billing/workflow.go`                       | No     |
| Workflows/CatalogWorkflow.php                       | `billing/catalog.go`                        | No     |
| Workflows/SubscriptionStageWorkflow.php             | `subscription/transitioner.go`              | No     |

---

## Not Ported (by design)

These PHP files have no Go equivalent because they are Laravel/PHP framework plumbing:

| PHP File                                     | Reason                                                    |
| -------------------------------------------- | --------------------------------------------------------- |
| CashierFake.php                              | Testing utility — Go tests use interface mocks            |
| CashierServiceProvider.php                   | Laravel service container — Go uses explicit construction |
| SparkServiceProvider.php                     | Laravel service container — Go uses explicit construction |
| Console/InstallCommand.php                   | Package scaffolding command — not applicable              |
| Http/Controllers/RetrievesBillableModels.php | PHP trait — logic absorbed into handler helpers           |
| Http/Middleware/HandleInertiaRequests.php    | Inertia.js middleware — handled by inertia-go             |
| Components/Button.php                        | Blade UI component — frontend concern                     |
| Components/Checkout.php                      | Blade UI component — frontend concern                     |

---

## Recommended Next Steps (Priority Order)

### Tier 1: Unit-testable business logic (no DB needed)

These can be tested with interface mocks — highest value, lowest friction:

1. **Subscription Gate** — `entitlement/gate.go` (7 PHP methods)
   - Tests: CheckRequirements with various subscription states and plan filters
2. **Webhook Signature** — `webhook/signature.go` (5 PHP methods)
   - Tests: HMAC verification, malformed headers, time drift, middleware integration
3. **Transaction Snapshot** — `snapshot.go` (3 PHP methods)
   - Tests: Empty state, incomplete money fields, complete money fields
4. **Plan Period Price error path** — `plan.go` (1 PHP method)
   - Tests: DisplayAmount with missing amount/currency throws error

### Tier 2: State machine and lifecycle (mock stores)

5. **Subscription Transitioner** — `subscription/transitioner.go` (12 PHP methods)
   - Tests: All state transitions, entitlement side effects, idempotency
6. **Subscription Repository** — `subscription/repository.go` (~20 PHP methods)
   - Tests: Pending creation, starter trial, checkout, idempotency
7. **Billing State Resolver** — `subscription/state.go` (~8 PHP methods)
   - Tests: State snapshots for active, canceled, paused, starter trial, no subscription

### Tier 3: Checkout and reconciliation (mock stores + provider)

8. **Checkout Reconciler** — `checkout/reconciler.go` (18 PHP methods)
   - Tests: Pending→active matching, plan copying, feature provisioning, error handling
9. **Checkout Syncer** — `checkout/syncer.go` (3 PHP methods)
   - Tests: Plan/period sync after provider update
10. **Checkout Recovery** — `checkout/recovery.go` (1 PHP method)
    - Tests: Explicit provider subscription recovery

### Tier 4: Catalog and presentation

11. **Spark Plan Registry** — `catalog/presentation.go` (3 PHP methods)
12. **Price Manager** — `catalog/price_manager.go` (3 PHP methods)

### Tier 5: HTTP handlers (integration-level)

13. **Billing Handler** — `handler/billing.go` (19 PHP methods)
14. **Inquiry Handler** — `handler/inquiry.go` (7 PHP methods)
15. **Middleware** — `handler/middleware.go` (5 PHP methods)
16. **Portal Handler** — `handler/portal.go` (3 PHP methods)

### Tier 6: Commands and jobs

17. **Expire Command** — `command/command.go` (4 PHP methods)
18. **Recover Command** — `command/command.go` (1 PHP method)
