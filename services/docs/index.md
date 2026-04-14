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
    link: /packages/auth
    type: secondary
features:
  - title: auth
    details: Authentication, authorization, and password management. Supports session, token, and request guards with configurable user providers.
  - title: bus
    details: Command and event bus with synchronous, async, deferred, and chained dispatch. Includes pipeline middleware and distributed unique-job locking.
  - title: cache
    details: Two-level caching abstraction — Store for backend operations and Repository for high-level helpers including tags and distributed locks.
  - title: cookie
    details: Laravel-inspired cookie management with a queuing jar, factory interfaces, and HTTP middleware for transparent encryption and decryption.
  - title: fortify
    details: Application fortification with rate limiting, two-factor authentication, and pipeline-based authentication flows.
  - title: httpx
    details: Rich HTTP request and response primitives. Fluent response writers, file upload handling, outbound client with test fakes, and JSON API resources.
  - title: jetstream
    details: Team and organization management layer with member invitations, role assignment, and multi-tenancy support.
  - title: queue
    details: Job queue management with multiple driver implementations (sync, database, Redis, SQS, beanstalkd) and a configurable background Worker.
  - title: routing
    details: A faithful 1:1 Go port of Laravel's routing layer. Full resource routing, named routes, middleware groups, and URL generation.
  - title: session
    details: HTTP session management with flash data, CSRF tokens, and swappable backends — array, file, database, cache, cookie, and encrypting handlers.
  - title: spark
    details: Subscription billing, checkout, and entitlement management. Stripe/Paddle-style cashier with plans, prices, orders, and webhooks.
footer: MIT Licensed
---
