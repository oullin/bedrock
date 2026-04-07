# Billing Port — Compliance Report

## Overview

This package ports three PHP codebases into a single provider-agnostic Go package:

| Source | Location | PHP Files |
|--------|----------|-----------|
| upstream/cashier-paddle | `vendor/upstream/cashier-paddle/src/` | 32 |
| upstream/billing-paddle | `vendor/upstream/billing-paddle/src/` | 24 |
| Madora app/Billing | `app/Billing/` | 39 |
| **Total** | | **95** |

**Ported**: 88 / 95 (92.6%)
**Not ported**: 7 (framework plumbing not applicable to Go — see below)

---

## Test Coverage

| Go File | Test File | Status |
|---------|-----------|--------|
| config.go | config_test.go | Covered |
| customer.go | customer_test.go | Covered |
| enums.go | enums_test.go | Covered |
| manager.go | manager_test.go | Covered |
| plan.go | plan_test.go | Covered |
| subscription.go | subscription_test.go | Covered |
| validation.go | validation_test.go | Covered |
| billable.go | — | Not covered |
| billing.go | — | Not covered |
| catalog.go | — | Not covered |
| catalog_presentation.go | — | Not covered |
| catalog_price_manager.go | — | Not covered |
| catalog_seeder.go | — | Not covered |
| checkout.go | — | Not covered |
| checkout_reconciler.go | — | Not covered |
| checkout_recovery.go | — | Not covered |
| checkout_session.go | — | Not covered |
| checkout_syncer.go | — | Not covered |
| command_expire.go | — | Not covered |
| command_reconcile.go | — | Not covered |
| command_recover.go | — | Not covered |
| currency.go | — | Not covered |
| entitlement.go | — | Not covered |
| entitlement_sync.go | — | Not covered |
| errors.go | — | Not covered |
| event.go | — | Not covered |
| frontend_state.go | — | Not covered |
| gate.go | — | Not covered |
| handler_billing.go | — | Not covered |
| handler_inquiry.go | — | Not covered |
| handler_invoice.go | — | Not covered |
| handler_payment.go | — | Not covered |
| handler_portal.go | — | Not covered |
| handler_subscription.go | — | Not covered |
| interfaces.go | — | Not covered |
| invoice.go | — | Not covered |
| job_reconcile.go | — | Not covered |
| lifecycle.go | — | Not covered |
| lifecycle_state.go | — | Not covered |
| lifecycle_transition.go | — | Not covered |
| listener_reconcile.go | — | Not covered |
| listener_subscription_created.go | — | Not covered |
| listener_sync_entitlements.go | — | Not covered |
| listener_sync_update.go | — | Not covered |
| mail.go | — | Not covered |
| manager_builder.go | — | Not covered |
| middleware.go | — | Not covered |
| payment.go | — | Not covered |
| price.go | — | Not covered |
| price_preview.go | — | Not covered |
| routes.go | — | Not covered |
| snapshot.go | — | Not covered |
| subscription_builder.go | — | Not covered |
| subscription_feature.go | — | Not covered |
| subscription_item.go | — | Not covered |
| transaction.go | — | Not covered |
| view.go | — | Not covered |
| webhook.go | — | Not covered |
| webhook_signature.go | — | Not covered |
| workflow.go | — | Not covered |
| workflow_catalog.go | — | Not covered |

**Files with tests**: 7 / 62 source files (11.3%)

---

## PHP to Go Mapping — upstream/cashier-paddle

| PHP File | Go File | Tested |
|----------|---------|--------|
| Billable.php | billable.go | No |
| Cashier.php | manager.go | Yes |
| CashierFake.php | *Not ported* | — |
| CashierServiceProvider.php | *Not ported* | — |
| Checkout.php | checkout_session.go | No |
| Components/Button.php | *Not applicable* | — |
| Components/Checkout.php | *Not applicable* | — |
| Concerns/ManagesCustomer.php | customer.go | Yes |
| Concerns/ManagesSubscriptions.php | subscription.go | Yes |
| Concerns/ManagesTransactions.php | transaction.go | No |
| Concerns/PerformsCharges.php | checkout_session.go | No |
| Concerns/Prorates.php | subscription.go (methods) | Yes |
| Customer.php | customer.go | Yes |
| Events/CustomerUpdated.php | event.go | No |
| Events/SubscriptionCanceled.php | event.go | No |
| Events/SubscriptionCreated.php | event.go | No |
| Events/SubscriptionPaused.php | event.go | No |
| Events/SubscriptionUpdated.php | event.go | No |
| Events/TransactionCompleted.php | event.go | No |
| Events/TransactionUpdated.php | event.go | No |
| Events/WebhookHandled.php | event.go | No |
| Events/WebhookReceived.php | event.go | No |
| Exceptions/PaddleException.php | errors.go | No |
| Http/Controllers/WebhookController.php | webhook.go | No |
| Http/Middleware/VerifyWebhookSignature.php | webhook_signature.go | No |
| Payment.php | payment.go | No |
| Price.php | price.go | No |
| PricePreview.php | price_preview.go | No |
| Subscription.php | subscription.go | Yes |
| SubscriptionBuilder.php | subscription_builder.go | No |
| SubscriptionItem.php | subscription_item.go | No |
| Transaction.php | transaction.go | No |

## PHP to Go Mapping — upstream/billing-paddle

| PHP File | Go File | Tested |
|----------|---------|--------|
| Actions/GenerateCheckoutSession.php | interfaces.go (ProviderCheckoutGenerator) | No |
| Billable.php | billable.go | No |
| BillableConfigurationBuilder.php | manager_builder.go | No |
| Console/InstallCommand.php | *Not ported* | — |
| Contracts/Actions/GeneratesCheckoutSessions.php | interfaces.go | No |
| FrontendState.php | frontend_state.go | No |
| GuessesBillableTypes.php | manager.go (DefaultBillableType) | Yes |
| Http/Controllers/BillingPortalController.php | handler_portal.go | No |
| Http/Controllers/CancelSubscriptionController.php | handler_subscription.go | No |
| Http/Controllers/DownloadInvoiceController.php | handler_invoice.go | No |
| Http/Controllers/NewPendingCheckoutController.php | handler_payment.go | No |
| Http/Controllers/NewSubscriptionController.php | handler_subscription.go | No |
| Http/Controllers/ResumeSubscriptionController.php | handler_subscription.go | No |
| Http/Controllers/RetrievesBillableModels.php | *Not ported (trait)* | — |
| Http/Controllers/UpdatePaymentMethodController.php | handler_payment.go | No |
| Http/Controllers/UpdateSubscriptionController.php | handler_subscription.go | No |
| Http/Middleware/HandleInertiaRequests.php | *Not applicable* | — |
| Http/Middleware/VerifyBillableIsSubscribed.php | middleware.go | No |
| Listeners/SubscriptionCreatedListener.php | listener_subscription_created.go | No |
| Plan.php | plan.go | Yes |
| Billing.php | manager.go | Yes |
| BillingManager.php | manager.go | Yes |
| BillingServiceProvider.php | *Not ported* | — |
| ValidPlan.php | validation.go | Yes |

## PHP to Go Mapping — Madora app/Billing

| PHP File | Go File | Tested |
|----------|---------|--------|
| Actions/GenerateBillingCheckoutSession.php | interfaces.go (ProviderCheckoutGenerator) | No |
| BillingConstants.php | billing.go | No |
| Catalog/PlanCatalog.php | catalog.go | No |
| Catalog/BillingPlanRegistry.php | catalog_presentation.go | No |
| Checkout/CheckoutStarter.php | checkout.go | No |
| Checkout/ExplicitSubscriptionRecovery.php | checkout_recovery.go | No |
| Checkout/PaddleSubscriptionFetcher.php | interfaces.go (ProviderSubscriptionFetcher) | No |
| Checkout/SubscriptionReconciler.php | checkout_reconciler.go | No |
| Checkout/SubscriptionUpdateSyncer.php | checkout_syncer.go | No |
| Commands/ExpirePendingSubscriptionsCommand.php | command_expire.go | No |
| Commands/ReconcileSubscriptionsCommand.php | command_reconcile.go | No |
| Commands/RecoverExplicitSubscriptionCommand.php | command_recover.go | No |
| Contracts/BillingLocalDevUrl.php | interfaces.go (URLResolver) | No |
| Controllers/BillingController.php | handler_billing.go | No |
| Controllers/CustomPlanInquiryController.php | handler_inquiry.go | No |
| Data/Read/BillingPageData.php | view.go (BillingPageView) | No |
| Data/Read/BillingPriceData.php | view.go (BillingPriceView) | No |
| Data/Read/BillingSelectedPlanData.php | view.go (BillingSelectedPlanView) | No |
| Data/Read/BillingStateData.php | view.go (BillingStateView) | No |
| Data/Read/BillingTeamData.php | view.go (BillingTeamView) | No |
| Data/Read/FrontendPlanData.php | view.go (FrontendPlanView) | No |
| Data/Read/LandingPlanData.php | view.go (LandingPlanView) | No |
| Data/Read/PlanFeatureData.php | view.go (PlanFeatureView) | No |
| Data/Read/PlanPeriodData.php | view.go (PlanPeriodView) | No |
| Data/Read/SubscriptionCtaData.php | view.go (SubscriptionCTAView) | No |
| Data/Read/SubscriptionStateData.php | view.go (SubscriptionStateView) | No |
| Data/Read/TransactionData.php | view.go (TransactionView) | No |
| Data/Result/SubscriptionDecision.php | gate.go (SubscriptionDecision) | No |
| Data/Result/SubscriptionReconciliationData.php | checkout_reconciler.go | No |
| Data/State/BillingStateSnapshot.php | snapshot.go (BillingStateSnapshot) | No |
| Data/State/PlanConfig.php | snapshot.go (PlanConfig) | No |
| Data/State/SubscriptionCtaSnapshot.php | snapshot.go (SubscriptionCTASnapshot) | No |
| Data/State/SubscriptionStateSnapshot.php | snapshot.go (SubscriptionStateSnapshot) | No |
| Data/State/TransactionSnapshot.php | snapshot.go (TransactionSnapshot) | No |
| Entitlements/EntitlementAccess.php | entitlement.go | No |
| Entitlements/EntitlementSynchronizer.php | entitlement_sync.go | No |
| Entitlements/SubscriptionStatusEvaluator.php | enums.go (SubscriptionStatus methods) | Yes |
| Enums/BillingPeriod.php | enums.go | Yes |
| Enums/PlanPricingMode.php | enums.go | Yes |
| Enums/SubscriptionFeatureCode.php | enums.go | Yes |
| Enums/SubscriptionPlan.php | enums.go | Yes |
| Enums/SubscriptionStatus.php | enums.go | Yes |
| Events/SubscriptionChanged.php | event.go (SubscriptionChangedEvent) | No |
| Exceptions/InvalidPlanPriceException.php | errors.go (InvalidPlanPriceError) | No |
| Http/Middleware/NormaliseBillableRouteParameter.php | middleware.go | No |
| Jobs/ReconcileSubscriptions.php | job_reconcile.go | No |
| Lifecycle/BillingStateResolver.php | lifecycle_state.go | No |
| Lifecycle/SubscriptionRepository.php | lifecycle.go | No |
| Listeners/ReconcileSubscriptionAfterCheckout.php | listener_reconcile.go | No |
| Listeners/SyncSubscriptionAfterUpdate.php | listener_sync_update.go | No |
| Listeners/SyncSubscriptionEntitlements.php | listener_sync_entitlements.go | No |
| Mail/CustomPlanInquiryConfirmationMail.php | mail.go | No |
| Mail/CustomPlanInquiryMail.php | mail.go | No |
| Models/Currency.php | currency.go | No |
| Models/Feature.php | plan.go (Feature struct) | No |
| Models/Invoice.php | invoice.go | No |
| Models/Plan.php | plan.go | Yes |
| Models/PlanFeature.php | plan.go (PlanFeature struct) | No |
| Models/PlanPeriod.php | plan.go (PlanPeriod struct) | No |
| Models/PlanPeriodPrice.php | plan.go (PlanPeriodPrice struct) | Yes |
| Models/Subscription.php | subscription.go | Yes |
| Models/SubscriptionFeature.php | subscription_feature.go | No |
| Models/SubscriptionItem.php | subscription_item.go | No |
| Requests/CheckoutRequest.php | validation.go (CheckoutInput) | Yes |
| Requests/CustomPlanInquiryRequest.php | validation.go (CustomPlanInquiryInput) | Yes |
| Support/PlanPeriodPriceManager.php | catalog_price_manager.go | No |
| Support/PlanPresentation.php | catalog_presentation.go | No |
| Support/PlanSeeder.php | catalog_seeder.go | No |
| Support/SerialisesBillingPrices.php | view.go (BillingPriceView) | No |
| Support/SubscriptionGate.php | gate.go | No |
| Support/SubscriptionTransitioner.php | lifecycle_transition.go | No |
| Workflows/BillingWorkflow.php | workflow.go | No |
| Workflows/CatalogWorkflow.php | workflow_catalog.go | No |
| Workflows/SubscriptionStageWorkflow.php | lifecycle_transition.go | No |

---

## Not Ported (by design)

These PHP files have no Go equivalent because they are Upstream/PHP framework plumbing with no meaningful Go counterpart:

| PHP File | Reason |
|----------|--------|
| CashierFake.php | Testing utility — Go tests use interface mocks instead |
| CashierServiceProvider.php | Upstream service container — Go uses explicit construction |
| BillingServiceProvider.php | Upstream service container — Go uses explicit construction |
| Console/InstallCommand.php | Package scaffolding command — not applicable in Go |
| Http/Controllers/RetrievesBillableModels.php | PHP trait — logic absorbed into handler helpers |
| Http/Middleware/HandleInertiaRequests.php | Inertia.js middleware — not applicable in Go |
| Components/Button.php | Template UI component — frontend concern |
| Components/Checkout.php | Template UI component — frontend concern |

---

## Test Coverage Gap Summary

**Covered (7 files)**:
- `config.go` — Default config values
- `customer.go` — Generic trial checking
- `enums.go` — All enum types, validity, display labels, GrantsAccess, IsTerminal
- `manager.go` — Billable registration, plan management, per-seat billing
- `plan.go` — PlanPeriodPrice.DisplayAmount
- `subscription.go` — Status methods, trial, grace period, proration, product/price checks
- `validation.go` — CheckoutInput, CustomPlanInquiryInput, ValidPlan

**High-priority gaps** (business logic that should be tested next):
1. `gate.go` — SubscriptionGate.CheckRequirements (authorization decisions)
2. `entitlement.go` — EntitlementAccess (feature access, capacity)
3. `lifecycle_transition.go` — SubscriptionTransitioner (state machine)
4. `lifecycle.go` — SubscriptionRepository (subscription creation)
5. `checkout_reconciler.go` — SubscriptionReconciler (post-checkout matching)
6. `webhook_signature.go` — HMAC verification
7. `errors.go` — InvalidPlanPriceError construction
8. `snapshot.go` — PlanConfig.MatchesPriceID, HasConfiguredProviderIDs
