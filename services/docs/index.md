---
layout: HomeLayout
home: true
heroText: Bedrock
tagline: A collection of well-tested, Laravel-inspired libraries for building production-grade Go web applications.
actions:
  - text: Get Started
    link: /getting-started
    type: primary
  - text: Browse Packages
    link: /packages/container
    type: secondary
features:
  - title: container
    details: IoC service container and application bootstrap. Factory and singleton bindings, contextual resolution, tagging, and a service-provider lifecycle.
  - title: routing
    details: A faithful 1:1 Go port of Laravel's routing layer. Named routes, resource routing, middleware groups, implicit model binding, and URL generation.
  - title: httpx
    details: Rich HTTP request and response primitives. Fluent response writers, file upload handling, outbound client with test fakes, and JSON API resources.
  - title: validation
    details: Rule-based input validation with 80+ built-in rules. Pipe-delimited syntax, custom messages, MessageBag error collection, and extensible rule registry.
  - title: auth
    details: Authentication, authorization, and password management. Session, token, and request guards with Gate-based access control and password reset broker.
  - title: cache
    details: Two-level caching abstraction — Store for backend operations and Repository for high-level helpers including tags and distributed locks.
  - title: redis
    details: Full Redis command surface with pipelines, transactions, pub/sub, Lua scripts, cluster and sentinel modes — a 1:1 port of Illuminate/Redis.
  - title: events
    details: Event dispatcher with named and typed events, wildcard listeners, subscribers, queued delivery, and transaction-deferred dispatch.
  - title: queue
    details: Job queue management with sync, database, Redis, SQS, beanstalkd, and failover drivers, plus a configurable background Worker.
  - title: bus
    details: Command and event bus with synchronous, async, deferred, and chained dispatch. Pipeline middleware and distributed unique-job locking.
  - title: mailx
    details: Driver-based email sending with SMTP, log, and array transports. Rich message construction with attachments, embeds, headers, and events.
  - title: notifications
    details: Multi-channel notification delivery — mail, database, broadcast, and custom channels. Queued delivery and per-channel suppression.
  - title: encryption
    details: AES encryption with CBC and GCM mode support. Key rotation, MAC verification, and Laravel-compatible payload format.
  - title: hashing
    details: Driver-based password hashing with bcrypt, Argon2i, and Argon2id. Rehash detection and hash-format recognition.
  - title: session
    details: HTTP session management with flash data, CSRF tokens, and swappable backends — array, file, database, cache, cookie, and encrypting handlers.
  - title: filesystem
    details: Local filesystem operations — reads, writes, copy, move, delete, directory management, MIME detection, file hashing, and locked file access.
  - title: pagination
    details: Offset-based and cursor-based pagination with generic item slices. URL windows, JSON serialization, and Laravel-compatible API.
  - title: collection
    details: Fluent, type-safe collection helpers for slices, ordered maps, lazy sequences, and one-off array or key-value operations.
  - title: support
    details: Global helpers, Fluent dynamic bags, Optional[T], MessageBag, string builders, Lottery, Sleep, Timebox — a Go port of Illuminate/Support.
  - title: log
    details: Driver-based structured logging with channels and stack aggregation. Stream, rotating, stderr, and syslog handlers with shared context.
  - title: translation
    details: Localisation and i18n with key-based lookup, namespace fallback chains, CLDR pluralisation, and file or in-memory loaders.
  - title: inception
    details: Unified authentication scaffolding — Laravel Fortify + Jetstream ported. Registration, 2FA, teams, API tokens, profile management with feature flags.
  - title: spark
    details: Subscription billing, checkout, and entitlement management. Stripe/Paddle-style cashier with plans, prices, orders, and webhooks.
footer: MIT Licensed
---
