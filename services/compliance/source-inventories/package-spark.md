### `laravel/spark-paddle` Source Codebase Inventory

Private source baseline: `/Users/gocanto/Sites/madora/vendor/laravel/spark-paddle`

This package is private and closed source, so this inventory maps the installed Composer artifact instead of a public GitHub source tree. The artifact contains 62 files and no upstream test files.

| Source Path | Status | Bedrock Evidence | Notes |
| --- | --- | --- | --- |
| `.DS_Store` | excluded | n/a | Local filesystem metadata; not a package surface. |
| `RELEASE.md` | excluded | n/a | Release metadata; no runtime behavior. |
| `composer.json` | adapted | `packages/spark/go.mod`, `packages/spark/package.json` | PHP package metadata maps to Go module and package metadata. |
| `package.json` | adapted | `packages/spark/package.json` | Frontend package metadata is represented only where Bedrock ships package metadata. |
| `postcss.config.js` | excluded | n/a | Build tooling for shipped Vue/Inertia assets; Bedrock does not own a Spark asset pipeline. |
| `testbench.yaml` | excluded | n/a | Orchestra Testbench setup; no Go runtime surface. |
| `vite.config.js` | excluded | n/a | Build tooling for shipped Vue/Inertia assets; Bedrock does not own a Spark asset pipeline. |
| `config/spark.php` | adapted | `packages/spark/config.go` | Path, middleware, proration, date, brand, billable, plan, and feature config are represented as explicit Go config. |
| `database/migrations/2019_05_03_000001_create_customers_table.php` | adapted | `packages/spark/customer.go`, `packages/spark/store.go` | Customer schema is represented by the `Customer` model and `CustomerStore` interface. |
| `database/migrations/2019_05_03_000002_create_subscriptions_table.php` | adapted | `packages/spark/subscription.go`, `packages/spark/store.go` | Subscription schema is represented by the `Subscription` model and `SubscriptionStore` interface. |
| `database/migrations/2019_05_03_000003_create_subscription_items_table.php` | adapted | `packages/spark/subscription_item.go`, `packages/spark/store.go` | Subscription item schema is represented by the `SubscriptionItem` model and `SubscriptionItemStore` interface. |
| `database/migrations/2019_05_03_000004_create_transactions_table.php` | adapted | `packages/spark/transaction.go`, `packages/spark/store.go` | Transaction schema is represented by the `Transaction` model and `TransactionStore` interface. |
| `public/css/app.css` | adapted | `packages/spark/handler/billing_portal.go`, `services/docs/packages/spark.md` | Laravel publishes compiled UI assets; Bedrock exposes billing state and leaves UI assets to the host app. |
| `public/js/app.js` | adapted | `packages/spark/handler/billing_portal.go`, `services/docs/packages/spark.md` | Laravel publishes compiled UI assets; Bedrock exposes billing state and leaves UI assets to the host app. |
| `resources/css/spark.css` | adapted | `services/docs/packages/spark.md` | Source CSS is host-owned in Bedrock. |
| `resources/js/app.js` | adapted | `packages/spark/handler/billing_portal.go` | Inertia app bootstrap maps to JSON billing portal state. |
| `resources/js/Components/Button.vue` | adapted | `packages/spark/handler/billing_portal.go` | Vue component is host-owned; backing state is exposed by Go handlers. |
| `resources/js/Components/ErrorMessages.vue` | adapted | `packages/spark/errors.go`, `packages/spark/handler/*` | Error display is host-owned; Go handlers return HTTP errors and validation payloads. |
| `resources/js/Components/InfoMessages.vue` | adapted | `packages/spark/handler/billing_portal.go` | UI messaging is host-owned. |
| `resources/js/Components/IntervalSelector.vue` | adapted | `packages/spark/state/frontend.go`, `packages/spark/plan.go` | Monthly/yearly plan grouping is represented in frontend state. |
| `resources/js/Components/InvoiceList.vue` | adapted | `packages/spark/handler/download_invoice.go`, `packages/spark/transaction.go` | Invoice behavior is exposed through transactions and invoice download handlers. |
| `resources/js/Components/Modal.vue` | adapted | `packages/spark/handler/*` | Modal UI is host-owned; handler commands are exposed as JSON endpoints. |
| `resources/js/Components/Paginator.vue` | adapted | `packages/spark/handler/order.go`, `packages/spark/store.go` | Pagination UI is host-owned; stores expose limited result queries. |
| `resources/js/Components/Plan.vue` | adapted | `packages/spark/plan.go`, `packages/spark/state/frontend.go` | Plan serialization and active plan state are represented in Go. |
| `resources/js/Components/PlanList.vue` | adapted | `packages/spark/plan.go`, `packages/spark/state/frontend.go` | Plan lists map to monthly/yearly frontend state. |
| `resources/js/Components/PlanSectionHeading.vue` | adapted | `packages/spark/state/frontend.go` | Heading UI is host-owned; section data is exposed in state. |
| `resources/js/Components/SecondaryButton.vue` | adapted | `packages/spark/handler/*` | Button UI is host-owned; commands are Go handlers. |
| `resources/js/Components/SectionHeading.vue` | adapted | `packages/spark/state/frontend.go` | Heading UI is host-owned. |
| `resources/js/Components/SuccessMessage.vue` | adapted | `packages/spark/handler/*` | Flash UI is host-owned; handlers return command status. |
| `resources/js/Icons/ChevronIcon.vue` | adapted | n/a | Icon asset is host-owned. |
| `resources/js/Icons/ClearIcon.vue` | adapted | n/a | Icon asset is host-owned. |
| `resources/js/Icons/LoadingIcon.vue` | adapted | n/a | Icon asset is host-owned. |
| `resources/js/Mixins/Base.js` | adapted | `packages/spark/state/frontend.go`, `packages/spark/handler/*` | Shared UI behavior maps to explicit state and handler contracts. |
| `resources/js/Pages/BillingPortal.vue` | adapted | `packages/spark/handler/billing_portal.go`, `packages/spark/state/frontend.go` | Inertia page maps to a caller-owned UI over Go billing portal state. |
| `resources/views/app.blade.php` | adapted | `packages/spark/handler/billing_portal.go` | Blade shell is host-owned; Bedrock exposes state instead of rendering the shell. |
| `routes/web.php` | ported | `packages/spark/handler/routes.go` | Spark subscription, pending checkout, invoice, payment method, and portal routes are registered explicitly. |
| `src/Actions/GenerateCheckoutSession.php` | ported | `packages/spark/action/create_subscription.go`, `packages/spark/checkout.go` | Checkout generation, plan eligibility, and per-seat quantity are represented in Go. |
| `src/Billable.php` | adapted | `packages/spark/billable.go`, `packages/spark/customer.go`, `packages/spark/subscription.go`, `packages/spark/service/billing.go` | Eloquent trait behavior maps to explicit billable interfaces, customer records, subscriptions, and billing services. |
| `src/BillableConfigurationBuilder.php` | ported | `packages/spark/manager.go` | Resolve, authorize, eligibility, per-seat, and plan builder methods are implemented. |
| `src/Console/InstallCommand.php` | adapted | `packages/spark/spark_service_provider.go`, `services/docs/packages/spark.md` | Laravel installer command maps to explicit provider setup and documentation. |
| `src/Contracts/.DS_Store` | excluded | n/a | Local filesystem metadata; not a package surface. |
| `src/Contracts/Actions/GeneratesCheckoutSessions.php` | adapted | `packages/spark/action/create_subscription.go`, `packages/spark/checkout.go` | PHP container contract maps to an explicit Go action type and checkout value. |
| `src/FrontendState.php` | adapted | `packages/spark/state/frontend.go`, `packages/spark/madora_inventory_test.go` | Billing portal state, pending-checkout state, plan grouping, subscription state, CTA, dashboard URL, terms URL, and seat label are represented; user/profile/Paddle token UI data is host-owned. |
| `src/GuessesBillableTypes.php` | ported | `packages/spark/manager.go` | Single-billable defaulting maps to `Manager.DefaultBillableType`. |
| `src/Http/Controllers/BillingPortalController.php` | adapted | `packages/spark/handler/billing_portal.go`, `packages/spark/state/frontend.go` | Portal rendering maps to JSON state; Inertia/Blade rendering is host-owned. |
| `src/Http/Controllers/CancelSubscriptionController.php` | ported | `packages/spark/handler/cancel_subscription.go`, `packages/spark/service/billing.go` | Active subscription validation and cancel flow are implemented through the billing service. |
| `src/Http/Controllers/DownloadInvoiceController.php` | adapted | `packages/spark/handler/download_invoice.go`, `packages/spark/transaction.go` | Invoice route scoping and transaction lookup are implemented; provider PDF streaming remains host/provider-owned. |
| `src/Http/Controllers/NewPendingCheckoutController.php` | ported | `packages/spark/handler/pending_checkout.go`, `packages/spark/customer.go` | Pending checkout ID is persisted on the customer record for the resolved billable. |
| `src/Http/Controllers/NewSubscriptionController.php` | ported | `packages/spark/handler/new_subscription.go`, `packages/spark/action/create_subscription.go` | Plan validation, eligibility checks, per-seat checkout quantity, and checkout serialization are implemented. |
| `src/Http/Controllers/ResumeSubscriptionController.php` | ported | `packages/spark/handler/resume_subscription.go`, `packages/spark/service/billing.go` | Grace-period resume is implemented through the billing service. |
| `src/Http/Controllers/RetrievesBillableModels.php` | adapted | `packages/spark/billable.go`, `packages/spark/manager.go`, `packages/spark/handler/*` | Eloquent lookup and trait validation map to injected resolvers, explicit interfaces, and authorization callbacks. |
| `src/Http/Controllers/UpdatePaymentMethodController.php` | adapted | `packages/spark/handler/payment_methods.go`, `packages/spark/handler/routes.go` | The upstream route is registered; provider-specific transaction/session creation remains host/provider-owned. |
| `src/Http/Controllers/UpdateSubscriptionController.php` | ported | `packages/spark/handler/update_subscription.go`, `packages/spark/action/update_subscription.go` | Existing subscription validation, plan validation, eligibility, and plan swap persistence are implemented. |
| `src/Http/Middleware/HandleInertiaRequests.php` | adapted | `packages/spark/handler/billing_portal.go` | Inertia asset versioning and flash sharing are not part of the Go transport; state is returned directly. |
| `src/Http/Middleware/VerifyBillableIsSubscribed.php` | ported | `packages/spark/handler/middleware.go` | Subscribed checks, redirects, and JSON payment-required responses are implemented. |
| `src/Listeners/SubscriptionCreatedListener.php` | adapted | `packages/spark/reconcile.go`, `packages/spark/webhook/handler.go` | Pending checkout cleanup and duplicate subscription handling map to reconciliation and webhook subscription sync. |
| `src/Plan.php` | ported | `packages/spark/plan.go` | Fluent plan intervals, incentives, descriptions, features, options, active/archive state, and serialization are implemented. |
| `src/Spark.php` | adapted | `packages/spark/manager.go`, `packages/spark/spark_service_provider.go` | Facade access maps to explicit manager and container registration. |
| `src/SparkManager.php` | ported | `packages/spark/manager.go` | Billable resolution, authorization, eligibility, per-seat billing, plan registration, config-backed plan expansion, and model lookup are represented by explicit Go APIs. |
| `src/SparkServiceProvider.php` | adapted | `packages/spark/spark_service_provider.go`, `packages/spark/handler/routes.go` | Laravel provider bootstrapping maps to explicit container registration and route setup. |
| `src/ValidPlan.php` | ported | `packages/spark/manager.go`, `packages/spark/handler/new_subscription.go`, `packages/spark/handler/update_subscription.go` | Active-plan validation is implemented. |
| `stubs/SparkServiceProvider.php` | adapted | `packages/spark/spark_service_provider.go`, `services/docs/packages/spark.md` | Laravel stub maps to explicit host registration in Go. |
| `stubs/en.json` | adapted | `packages/spark/handler/*`, `services/docs/packages/spark.md` | Translation strings are host-owned; handlers expose status/errors for host localization. |
