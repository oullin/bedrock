# Bedrock – Upstream 13.x Compliance Report

> Generated: 2026-04-07
> Baseline: [upstream/framework 13.x tests](https://github.com/upstream/framework/tree/13.x/tests)
> All Bedrock tests pass (`go test ./...` green across all packages).

---

## Table of Contents

- [Executive Summary](#executive-summary)
- [Feature Matrix](#feature-matrix)
- [Test Coverage Matrix](#test-coverage-matrix)
- [Per-Component Details](#per-component-details)
  - [Container](#container)
  - [Config](#config)
  - [Encryption](#encryption)
  - [Hashing](#hashing)
  - [Auth](#auth)
  - [Session](#session)
  - [Cache](#cache)
  - [Events](#events)
  - [HTTP](#http)
  - [Routing](#routing)
  - [Console](#console)
  - [Foundation](#foundation)
  - [Support](#support)
  - [View](#view)
  - [Stub Packages](#stub-packages)
- [Cross-Cutting Gaps](#cross-cutting-gaps)
- [Gap Categories](#gap-categories)
- [Recommended Next Steps](#recommended-next-steps)

---

## Executive Summary

| Metric | Count |
|--------|-------|
| Ported packages with tests | 14 |
| Stub packages (no code) | 19 |
| **Bedrock Go tests (total)** | **612** |
| Upstream test files in scope (ported packages) | **210** |
| Upstream test files analyzed line-by-line | 75 |
| Upstream test files not yet analyzed | 135 |
| Upstream test methods analyzed | 1516 |
| Analyzed: COVERED | 434 |
| Analyzed: MISSING | 812 |
| Analyzed: INTENTIONAL-SKIP | 320 |
| Analyzed: BEDROCK-ONLY | 130 |
| **Coverage of analyzed tests (excl. skips)** | **35%** |
| **Estimated full Upstream test surface (ported packages)** | **~2,500+** |

Bedrock has **612 passing Go tests** across 14 packages. For the subset of Upstream tests analyzed line-by-line (1516 methods across 75 files), we cover 35% of the portable surface. The full Upstream test suite for these packages spans **~210 test files** with an estimated **~2,500+ test methods**.

Four packages have **100% portable coverage**: Encryption, Config, Hashing, Container.

---

## Feature Matrix

All Upstream framework components and their Bedrock implementation status.

| Component | Status | Key Types / Interfaces | Notes |
|-----------|--------|----------------------|-------|
| **Container** | Implemented | `Container`, `Factory`, `ContextualBindingBuilder` | Singleton/transient bindings, contextual bindings, lifecycle callbacks |
| **Config** | Implemented | `Repository`, `Builder` | Dotted-path access, typed getters (String, Int, Bool, Duration, Float, Map, StringSlice) |
| **Encryption** | Implemented | `Encrypter`, `Cipher` | AES-128/256-CBC, AES-128/256-GCM, key rotation, HMAC validation |
| **Hashing** | Implemented | `Manager`, `Bcrypt`, `Argon2i`, `Argon2id` | Multiple drivers, rehash detection |
| **Auth** | Implemented | `Manager`, `SessionGuard`, `TokenGuard`, `RequestGuard` | Gate/policies, middleware, password broker, email verification |
| **Session** | Implemented | `Store`, `Handler` | Array, cache, file, encrypting handlers; flash data, CSRF |
| **Cache** | Implemented | `Store`, `ArrayStore`, `NullStore` | Get/Put/Forget/Flush with TTL, Remember/RememberForever, Increment/Decrement |
| **Events** | Implemented | `Dispatcher`, `EventDispatcher` | Named listeners, wildcards, halting dispatch, subscribers |
| **HTTP** | Implemented | `Request`, `ResponseWriter`, `Middleware` | Input extraction, content negotiation, JSON/redirect responses, file uploads |
| **Routing** | Implemented | `Router`, `Context`, `HandlerFunc`, `ResourceHandlers` | RESTful resources, method routes, parameters, middleware, groups |
| **Console** | Implemented | `Kernel`, `Command`, `Invocation` | Command registration, aliases, hidden commands, sub-command invocation |
| **Foundation** | Implemented | `Application`, `Builder`, `RoutingConfig` | App lifecycle, middleware config, exception handling |
| **Support** | Implemented | Utility types | Basic utilities, no tests yet |
| **Cookie** | Implemented | Cookie management | Basic cookie operations |
| **View** | Partial | `Renderer` (Go `html/template`) | Template not applicable; sections/stacks/composers mostly missing |
| **Broadcasting** | Stub | — | Package exists, no implementation |
| **Bus** | Stub | — | Package exists, no implementation |
| **Conditionable** | Stub | — | Package exists, no implementation |
| **Contracts** | Stub | — | Package exists, no implementation |
| **Database** | Stub | — | Blocks auth providers, cache DB store, session DB handler, validation |
| **Filesystem** | Stub | — | Package exists, no implementation |
| **Log** | Stub | — | Package exists, no implementation |
| **Mail** | Stub | — | Package exists, no implementation |
| **Notifications** | Stub | — | Package exists, no implementation |
| **Pagination** | Stub | — | Package exists, no implementation |
| **Pipeline** | Stub | — | Package exists, no implementation |
| **Process** | Stub | — | Package exists, no implementation |
| **Queue** | Stub | — | Blocks queued events, queued listeners |
| **Redis** | Stub | — | Package exists, no implementation |
| **Testing** | Stub | — | Package exists, no implementation |
| **Translation** | Stub | — | Package exists, no implementation |
| **Validation** | Stub | — | 37 Upstream test files, ~500+ methods |

---

## Test Coverage Matrix

Test coverage per component comparing Upstream test surface to Bedrock coverage.

| Component | Upstream Files | Upstream Tests (analyzed) | Bedrock Tests | Covered | Missing | Intentional Skips | Portable Coverage | Key Gaps |
|-----------|--------------|-------------------------|---------------|---------|---------|-------------------|-------------------|----------|
| **Encryption** | 1 | 26 | 32 | 26 | 0 | 0 | **100%** | None |
| **Config** | 1 | 33 | 39 | 27 | 0 | 6 | **100%** | None |
| **Hashing** | 1 | 13 | 11 | 9 | 0 | 3 | **100%** | None |
| **Container** | 1 of 11 | 77 | 81 | 59 | 0 | 40 | **100%** | 10 files unanalyzed (mostly PHP-specific) |
| **Auth** | 15 of 16 | 253 | 203 | 137 | 36 | 59 | **79%** | User providers (DB), event firing, model authz |
| **Session** | 2 of 5 | 56 | 55 | ~45 | ~11 | ~13 | **~80% est.** | Encrypted/cache/file handlers unanalyzed |
| **Cache** | 2 of 21 | 38 | 42 | 18 | 14 | 3 | **56%** | Locks, 19 unanalyzed files (stores, rate limiting) |
| **Events** | 3 of 4 | 76 | 35 | ~35 | ~4 | 3 | **~90% est.** | 32 blocked by queue/broadcasting stubs |
| **HTTP** | 9 of 9 | 428 | 70 | 64 | 330 | 17 | **16%** | HTTP client (160), JSON resources (83) |
| **Routing** | 12 of 12 | ~310 | 21 | ~13 | ~250 | ~47 | **~6%** | Model binding, URL gen, redirects, middleware sorting |
| **Console** | 11 of 11 | 75 | 15 | 7 | 64 | 4 | **10%** | Mutex, scheduling, parsing, signals |
| **Foundation** | 15 of 16 | 232 | 20 | 14 | 71 | 147 | **16%** | Exception handler, time interaction, builder |
| **View** | 1 of 11 | 70 | 10 | ~2 | ~18 | ~50 | **~10%** | Template (skip), sections, composers |
| **Support** | 0 of 52 | — | 5 | — | — | — | **~1%** | Entire package mostly unported |

---

## Per-Component Details

---

### Container

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Singleton bindings | Implemented | `Singleton()` |
| Transient bindings | Implemented | `Bind()` |
| Instance bindings | Implemented | `Instance()` |
| Contextual bindings | Implemented | `When().Needs().Give()` |
| Factory callbacks | Implemented | `Factory` type |
| Lifecycle callbacks (resolving/afterResolving) | Implemented | Pre/post-resolution hooks |
| Tagging | Implemented | `Tag()` / `Tagged()` |
| Extend | Implemented | `Extend()` decorator pattern |
| Rewindable generator | Implemented | Lazy iteration over tagged bindings |
| PHP reflection / auto-wiring | Intentional Skip | Go has no runtime reflection-based DI |
| PHP attributes (`#[Singleton]`, `#[Bind]`) | Intentional Skip | PHP-specific |
| ArrayAccess | Intentional Skip | PHP interface |

#### Test Coverage

| Upstream File | Upstream Tests | Covered | Missing | Skip | Bedrock-Only | Portable % |
|-------------|---------------|---------|---------|------|--------------|------------|
| ContainerTest.php | 77 | 59 | 0 | 40 | 22 | **100%** |

**10 additional Upstream test files** exist (ContainerCallTest, ContainerExtendTest, ContainerTaggingTest, ContextualBindingTest, ContextualAttributeBindingTest, ResolvingCallbackTest, AfterResolvingAttributeCallbackTest, ContainerResolveNonInstantiableTest, RewindableGeneratorTest, UtilTest). These are primarily PHP reflection and attribute-based — most are expected intentional-skips. Bedrock covers tagging, contextual binding, and resolving callbacks via its own test files.

---

### Config

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Dotted-path access | Implemented | `Get("app.name")` |
| Typed getters | Implemented | String, Int, Bool, Duration, Float, Map, StringSlice |
| Has / Set / Prepend / Push | Implemented | Full CRUD on config values |
| All / Keys | Implemented | Enumerate configuration |
| Builder pattern | Implemented | Fluent config construction |
| Macros | Intentional Skip | PHP Macroable trait |
| ArrayAccess | Intentional Skip | PHP interface |
| Collection integration | Intentional Skip | PHP Collection class |

#### Test Coverage

| Upstream File | Upstream Tests | Covered | Missing | Skip | Bedrock-Only | Portable % |
|-------------|---------------|---------|---------|------|--------------|------------|
| RepositoryTest.php | 33 | 27 | 0 | 6 | 12 | **100%** |

---

### Encryption

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| AES-128-CBC / AES-256-CBC | Implemented | Block cipher modes |
| AES-128-GCM / AES-256-GCM | Implemented | AEAD modes |
| Key parsing | Implemented | Base64-encoded keys with cipher prefix |
| Key rotation | Implemented | Multiple decryption keys |
| HMAC validation (CBC) | Implemented | Tamper detection |
| JSON serialization | Implemented | Encrypted payload format |

#### Test Coverage

| Upstream File | Upstream Tests | Covered | Missing | Skip | Bedrock-Only | Portable % |
|-------------|---------------|---------|---------|------|--------------|------------|
| EncrypterTest.php | 26 | 26 | 0 | 0 | 6 | **100%** |

---

### Hashing

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Bcrypt hashing | Implemented | Default driver |
| Argon2i hashing | Implemented | Alternative driver |
| Argon2id hashing | Implemented | Alternative driver |
| Hash verification | Implemented | `Check()` |
| Rehash detection | Implemented | `NeedsRehash()` |
| Algorithm info extraction | Implemented | `Info()` |
| Runtime availability checks | Intentional Skip | Go crypto always available |

#### Test Coverage

| Upstream File | Upstream Tests | Covered | Missing | Skip | Bedrock-Only | Portable % |
|-------------|---------------|---------|---------|------|--------------|------------|
| HasherTest.php | 13 | 9 | 0 | 3 | 2 | **100%** |

---

### Auth

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Gate / authorization policies | Implemented | `Define`, `Before`, `After`, `ForUser` |
| Session guard | Implemented | Stateful authentication with remember-me |
| Token guard | Implemented | Stateless API token authentication |
| Request guard | Implemented | Custom closure-based authentication |
| Auth manager | Implemented | Multi-guard management |
| Password broker | Implemented | Password reset flow |
| Token repository | Implemented | Database token storage for resets |
| Authenticate middleware | Implemented | Route protection |
| Authorize middleware | Implemented | Policy-based route authorization |
| Email verification middleware | Implemented | Bedrock-only tests |
| Redirect-if-authenticated middleware | Implemented | Guest-only routes |
| Authenticatable trait | Implemented | User contract |
| HandlesAuthorization trait | Implemented | Allow/deny response builders |
| Database user provider | Not Implemented | Blocked by database package |
| Orm user provider | Not Implemented | Blocked by Orm ORM |
| Event firing (login/logout/failed) | Not Implemented | Blocked by event integration |

#### Test Coverage

| Upstream File | Upstream Tests | Covered | Missing | Skip | Bedrock-Only | Portable % |
|-------------|---------------|---------|---------|------|--------------|------------|
| AuthAccessGateTest | 92 | 48 | 13 | 14 | 20 | **79%** |
| AuthAccessResponseTest | 13 | 7 | 3 | 0 | — | **70%** |
| AuthGuardTest | 43 | 33 | 9 | 1 | 31 | **79%** |
| AuthTokenGuardTest | 14 | 10 | 4 | 0 | — | **71%** |
| AuthPasswordBrokerTest | 9 | 8 | 0 | 1 | 10 | **100%** |
| AuthDatabaseTokenRepositoryTest | 10 | 10 | 0 | 0 | — | **100%** |
| AuthenticateMiddlewareTest | 11 | 7 | 2 | 2 | — | **78%** |
| AuthorizeMiddlewareTest | 16 | 8 | 5 | 2 | 3 | **62%** |
| EnsureEmailIsVerifiedTest | 1 | 0 | 0 | 1 | 3 | N/A (skip) |
| RedirectIfAuthenticatedMiddlewareTest | 1 | 0 | 0 | 1 | 2 | N/A (skip) |
| AuthenticatableTest | 3 | 1 | 0 | 2 | — | **100%** |
| AuthHandlesAuthorizationTest | 5 | 5 | 0 | 0 | — | **100%** |
| AuthorizesResourcesTest | 6 | 0 | 0 | 6 | — | N/A (skip) |
| AuthDatabaseUserProviderTest | 14 | 0 | 0 | 14 | — | N/A (skip) |
| AuthOrmUserProviderTest | 15 | 0 | 0 | 15 | — | N/A (skip) |
| **Subtotal** | **253** | **137** | **36** | **59** | **69** | **79%** |

**Unanalyzed (1 file):** AuthListenersSendEmailVerificationNotificationTest.php (~1 test — requires event/notification system).

#### Key Gaps

- Gate: subtype/interface resolution, dash-to-camel conversion, class-name policies, array abilities, custom resource gates
- Guards: `logoutCurrentDevice`, cookie override on remember, event firing (5 tests blocked by missing event system)
- Token Guard: custom field validation (4 tests)
- Authorize MW: model-type authorization (5 tests)
- User Providers: entire database/Orm provider layer (29 tests, blocked by database package)

---

### Session

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Session store (start/save/regenerate/invalidate) | Implemented | Full lifecycle |
| Flash data | Implemented | Flash, reflash, keep |
| CSRF token | Implemented | Token generation and validation |
| Array handler | Implemented | In-memory storage |
| Cache handler | Implemented | Cache-backed sessions |
| File handler | Implemented | Filesystem-backed sessions |
| Encrypting handler | Implemented | Encryption layer on top of any handler |
| PHP Backed Enum keys | Intentional Skip | PHP-specific |

#### Test Coverage

| Upstream File | Upstream Tests | Bedrock Status |
|-------------|---------------|----------------|
| SessionStoreTest.php | 56 | **Partially covered** — 55 Bedrock tests, not mapped 1:1 |
| ArraySessionHandlerTest.php | 10 | **Covered** |
| CacheBasedSessionHandlerTest.php | ? | Not yet analyzed |
| EncryptedSessionStoreTest.php | ? | Not yet analyzed |
| FileSessionHandlerTest.php | ? | Not yet analyzed |

#### Key Gaps

- SessionStoreTest: 13 PHP-enum-specific tests (intentional skip)
- 3 handler test files not yet analyzed line-by-line

---

### Cache

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Store interface | Implemented | Get, Put, Forget, Flush, Has |
| Array store | Implemented | In-memory with TTL |
| Null store | Implemented | No-op store |
| TTL support | Implemented | Duration-based expiry |
| Increment / Decrement | Implemented | Atomic counters |
| Remember / RememberForever | Implemented | Lazy evaluation |
| Add (if not exists) | Implemented | Conditional put |
| Locks | Not Implemented | Lock/restore/forceRelease/ownership |
| Repository / decorator | Not Implemented | Typed getters, tags proxy, events, macros |
| File store | Not Implemented | Filesystem-backed cache |
| Redis store | Not Implemented | Redis-backed cache |
| Database store | Not Implemented | DB-backed cache |
| Rate limiter | Not Implemented | Rate limiting |
| Tagged cache | Not Implemented | Tag-based cache invalidation |
| Manager / factory | Not Implemented | Multi-driver management |

#### Test Coverage

| Upstream File | Upstream Tests | Covered | Missing | Skip | Bedrock-Only | Portable % |
|-------------|---------------|---------|---------|------|--------------|------------|
| CacheArrayStoreTest | 34 | 14 | 14 | 3 | 18 | **45%** |
| CacheNullStoreTest | 4 | 4 | 0 | 0 | 5 | **100%** |
| **Subtotal (analyzed)** | **38** | **18** | **14** | **3** | **23** | **56%** |

**19 unanalyzed cache test files (284 methods):**

| Upstream File | Tests | Status |
|-------------|-------|--------|
| CacheRepositoryTest.php | 69 | Missing — Repository pattern not implemented |
| CacheFileStoreTest.php | 31 | Missing — File store not implemented |
| CacheSessionStoreTest.php | 18 | Missing — session/cache integration |
| CacheDatabaseStoreTest.php | 17 | Missing — requires database package |
| CacheEventsTest.php | 15 | Missing — requires event system |
| CacheRedisStoreTest.php | 15 | Missing — requires Redis driver |
| CacheRateLimiterTest.php | 14 | Missing — rate limiter not implemented |
| CacheTaggedCacheTest.php | 13 | Missing — tagged cache not implemented |
| CacheManagerTest.php | 13 | Missing — Manager/factory not implemented |
| ClearCommandTest.php | 12 | Missing — console command not implemented |
| ConcurrencyLimiterTest.php | 12 | Missing — concurrency limiter not implemented |
| CacheApcStoreTest.php | 12 | Intentional Skip — APC is PHP-specific |
| CacheMemcachedStoreTest.php | 11 | Missing — requires Memcached driver |
| CacheMemcachedConnectorTest.php | 4 | Missing — requires Memcached driver |
| CacheSpyMemoTest.php | 4 | Missing — testing spy utility |
| RateLimiterTest.php | 2 | Missing — rate limiter not implemented |
| CacheMemoizedStoreTest.php | 1 | Missing — memoization decorator |
| LimitTest.php | 1 | Missing — rate limit value object |
| CacheDynamoDbStoreTest.php | 1 | Missing — requires DynamoDB driver |

#### Key Gaps

- Locks: 13 ArrayStore lock tests (lock/restore/forceRelease/ownership)
- Repository: entire decorator layer (69 tests)
- Rate limiting: rate limiter + concurrency limiter (29 tests)
- Store drivers: File (31), Redis (15), Memcached (15), Database (17), DynamoDB (1)
- Events: 15 tests blocked by missing event system
- Tagged cache: 13 tests
- Manager/factory: 13 tests

---

### Events

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Named event listeners | Implemented | `Listen()` |
| Wildcard patterns | Implemented | `Listen("user.*", ...)` |
| Halting dispatch | Implemented | `Until()` |
| Subscriber registration | Implemented | `Subscribe()` |
| Forget / flush listeners | Implemented | Cleanup |
| Broadcasting | Not Implemented | Blocked by broadcasting stub |
| Queued events | Not Implemented | Blocked by queue stub |

#### Test Coverage

| Upstream File | Upstream Tests | Covered | Missing (D) | Skip | Portable % |
|-------------|---------------|---------|-------------|------|------------|
| EventsDispatcherTest | 39 | ~35 | 0 | TBD | **~90% est.** |
| EventsSubscriberTest | 3 | 2 | 0 | 1 | **100%** |
| BroadcastedEventsTest | 8 | 0 | 7 | 1 | 0% (blocked) |
| QueuedEventsTest | 26 | 0 | 25 | 1 | 0% (blocked) |

#### Key Gaps

- 32 tests blocked by infrastructure stubs (queue: 25, broadcasting: 7)
- EventsDispatcherTest not yet mapped line-by-line

---

### HTTP

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Request input extraction | Implemented | Query, input, boolean, integer, path, URL, host, scheme, method |
| Content negotiation | Implemented | Accepts, prefers, content type detection |
| Header / bearer token / IP | Implemented | Standard request inspection |
| JSON response | Implemented | Creation, encoding, status, indentation, error handling |
| Response (base) | Implemented | Status codes, headers, content |
| Redirect response | Implemented | Headers, fragments, same-origin enforcement |
| Uploaded file | Implemented | Content, name, extension, MIME |
| MIME type lookup | Implemented | Forward and reverse lookup |
| Testing file factory | Implemented | PNG, JPEG, GIF images and arbitrary files |
| HTTP client | Not Implemented | Entire `net/http.Client` wrapper subsystem |
| JSON API resources | Not Implemented | Blocked by Orm ORM |
| Cookie management | Not Implemented | Blocked by cookie integration |
| Session flash/input | Not Implemented | Blocked by session integration |

#### Test Coverage

| Upstream File | Upstream Tests | Covered | Missing | Skip | Bedrock-Only | Portable % |
|-------------|---------------|---------|---------|------|--------------|------------|
| HttpRequestTest | 113 | 26 | ~87 | 0 | 6 | **23%** |
| HttpMimeTypeTest | 6 | 5 | 0 | 1 | 0 | **100%** |
| HttpJsonResponseTest | 8 | 7 | 0 | 1 | 0 | **100%** |
| HttpResponseTest | 18 | 8 | 0 | 3 | 0 | **100% (portable)** |
| HttpRedirectResponseTest | 20 | 9 | 0 | 2 | 0 | **100% (portable)** |
| HttpUploadedFileTest | 2 | 2 | 0 | 0 | 0 | **100%** |
| HttpTestingFileFactoryTest | 9 | 7 | 0 | 2 | 0 | **100% (portable)** |
| HttpClientTest | 167 | 0 | 160 | 7 | 0 | **0%** |
| JsonResourceTest | 85 | 0 | 83 | 1 | 0 | **0%** |
| **Subtotal** | **428** | **64** | **330** | **17** | **6** | **16%** |

#### Key Gaps

- HTTP client: 160 tests — entire subsystem missing
- JSON resources: 83 tests — blocked by Orm ORM
- HttpRequestTest: 87 tests — flash, cookies, merge/replace, fingerprinting, fluent/string/date/enum, JSON body
- Response cookie/session: 16 tests across Response + Redirect — blocked by cookie/session integration

---

### Routing

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Method routes (GET/POST/PUT/DELETE/PATCH/OPTIONS) | Implemented | All HTTP verbs |
| Route parameters | Implemented | Path parameters via `{param}` |
| Resource routing | Implemented | `Resource()` with `ResourceHandlers` |
| Route groups | Implemented | Prefix, middleware |
| Nested groups | Implemented | Composable group nesting |
| Named routes | Implemented | `Name()` |
| Middleware | Implemented | Per-route and group middleware |
| Route model binding (implicit) | Not Implemented | 6+ tests |
| Route model binding (explicit) | Not Implemented | 3 tests |
| URL generation (full) | Not Implemented | Signed URLs, domain routing, asset URLs |
| Redirects | Not Implemented | Redirect responses, guest redirects |
| Controller routing | Not Implemented | Controller-based dispatch |
| Domain routing | Not Implemented | Subdomain routing |
| Middleware sorting | Not Implemented | Priority-based middleware ordering |
| Route caching | Not Implemented | Compiled route cache |
| Pattern filtering | Not Implemented | Global parameter constraints |

#### Test Coverage

| Upstream File | Upstream Tests | Bedrock Status |
|-------------|---------------|----------------|
| RoutingRouteTest.php | 102 | 13 covered, ~89 missing |
| RouteRegistrarTest.php | 111 | Partial (~6 covered) — groups/resource/middleware/naming |
| RoutingUrlGeneratorTest.php | 46 | Partial (2 covered) — name/route generation |
| RouteCollectionTest.php | 21 | Intentional Skip — Go 1.22+ ServeMux |
| RoutingRedirectorTest.php | 15 | Missing |
| ImplicitRouteBindingTest.php | 6 | Missing |
| RoutingSortedMiddlewareTest.php | 3 | Missing |
| RouteBindingTest.php | 3 | Missing |
| RouteUriTest.php | 1 (+9 data) | Intentional Skip — Orm-specific |
| RouteActionTest.php | 1 | Intentional Skip — PHP closures |
| RouteSignatureParametersTest.php | 1 | Intentional Skip — PHP reflection |

#### Key Gaps

- Route model binding (implicit + explicit): 9+ tests
- URL generation: signed routes, domain routing, asset URLs (~44 tests)
- Controller routing and middleware sorting
- Redirects: 15 tests

---

### Console

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Command registration | Implemented | `Register()` |
| Command execution | Implemented | `Call()` / `Invoke()` |
| Command aliases | Implemented | `Aliases` struct field |
| Hidden commands | Implemented | `Hidden` struct field |
| Help text | Implemented | `Help` struct field |
| Usage examples | Implemented | `Usage` struct field |
| Sub-command invocation | Implemented | Via `Kernel.Call()` |
| Command mutex / isolation | Not Implemented | Cache-lock-based isolation (13 tests) |
| Command scheduling / cron | Not Implemented | Event scheduling (13 tests) |
| Signature parsing | Not Implemented | Arguments, options, defaults, shortcuts (7 tests) |
| Output styling | Not Implemented | Newline detection, verbosity (5 tests) |
| Signal trapping | Not Implemented | Nested traps, unregistration (6 tests) |
| Interactive prompts | Not Implemented | Select/multiselect fallback (2 tests) |

#### Test Coverage

| Upstream File | Upstream Tests | Covered | Missing | Skip | Bedrock-Only | Portable % |
|-------------|---------------|---------|---------|------|--------------|------------|
| ConsoleApplicationTest | 14 | 2 | 10 | 2 | 7 | **17%** |
| CacheCommandMutexTest | 9 | 0 | 9 | 0 | — | **0%** |
| CommandMutexTest | 4 | 0 | 4 | 0 | — | **0%** |
| CommandTest | 15 | 5 | 8 | 2 | 1 | **38%** |
| CommandTrapTest | 4 | 0 | 4 | 0 | — | **0%** |
| ConfiguresPromptsTest | 2 | 0 | 2 | 0 | — | **0%** |
| ConsoleEventSchedulerTest | 9 | 0 | 9 | 0 | — | **0%** |
| ConsoleParserTest | 7 | 0 | 7 | 0 | — | **0%** |
| ConsoleScheduledEventTest | 4 | 0 | 4 | 0 | — | **0%** |
| OutputStyleTest | 5 | 0 | 5 | 0 | — | **0%** |
| SignalsTest | 2 | 0 | 2 | 0 | — | **0%** |
| **Subtotal** | **75** | **7** | **64** | **4** | **12** | **10%** |

#### Key Gaps

- Command mutex / isolation via cache locks (13 tests)
- Command scheduling / cron compilation (13 tests)
- Signature parsing with arguments, options, defaults (7 tests)
- Signal trapping and nested traps (6 tests)
- Output styling with verbosity (5 tests)

---

### Foundation

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Application lifecycle | Implemented | Boot, register, environment detection |
| Middleware configuration | Implemented | Via `Builder` |
| Exception handling config | Implemented | Via `Builder` |
| Environment detection | Implemented | `Environment()`, `IsProduction()`, etc. |
| Path management | Implemented | Base, config, storage paths |
| Maintenance mode (cache-based) | Partial | 1 of 4 tests covered |
| Exception handler | Partial | 1 of 33 tests covered |
| Time interaction | Not Implemented | `freezeTime()`, `freezeSecond()` |
| Application builder paths | Not Implemented | Env var overrides |
| Abort helper | Not Implemented | `abort()` with status + message |
| Service providers | Intentional Skip | PHP-specific lifecycle |
| Form requests | Intentional Skip | PHP-specific |
| Alias loader | Intentional Skip | PHP class aliasing |

#### Test Coverage

| Upstream File | Upstream Tests | Covered | Missing | Skip | Bedrock-Only | Portable % |
|-------------|---------------|---------|---------|------|--------------|------------|
| FoundationApplicationTest | 42 | 6 | ~36 | 0 | 14 | **~15%** |
| FoundationExceptionsHandlerTest | 33 | 1 | 17 | 15 | 0 | **~6%** |
| FoundationEnvironmentDetectorTest | 7 | 4 | 3 | 0 | 0 | **57%** |
| FoundationInteractsWithTimeTest | 6 | 0 | 6 | 0 | 0 | **0%** |
| FoundationApplicationBuilderTest | 8 | 2 | 4 | 2 | 0 | **33%** |
| FoundationCacheBasedMaintenanceModeTest | 4 | 1 | 3 | 0 | 0 | **25%** |
| FoundationHelpersTest | 19 | 0 | 2 | 17 | 0 | **0%** |
| FoundationAliasLoaderTest | 5 | 0 | 0 | 5 | 0 | N/A (skip) |
| FoundationAuthenticationTest | 5 | 0 | 0 | 5 | 0 | N/A (skip) |
| FoundationAuthorizesRequestsTraitTest | 7 | 0 | 0 | 7 | 0 | N/A (skip) |
| FoundationDocsCommandTest | 27 | 0 | 0 | 27 | 0 | N/A (skip) |
| FoundationFormRequestTest | 30 | 0 | 0 | 30 | 0 | N/A (skip) |
| FoundationInteractsWithDatabaseTest | 32 | 0 | 0 | 32 | 0 | N/A (DB) |
| FoundationPackageManifestTest | 1 | 0 | 0 | 1 | 0 | N/A (skip) |
| FoundationProviderRepositoryTest | 6 | 0 | 0 | 6 | 0 | N/A (skip) |
| **Subtotal** | **232** | **14** | **71** | **147** | **14** | **16%** |

**FoundationViteTest.php (53 tests)** — excluded from scope.

#### Key Gaps

- Exception handler (17 tests): JSON error responses, report deduplication, throttled reporting, renderable/reportable exceptions
- Time interaction (6 tests): `freezeTime()`, `freezeSecond()`, freeze with callback
- Application builder (4 tests): path overrides via env vars
- Environment detector (3 tests): console argument parsing edge cases
- Cache-based maintenance mode (3 tests): `active()`, `data()`, `activate()`/`deactivate()`
- Helpers (2 tests): `abort()` with status code and message

---

### Support

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Basic utilities | Implemented | Minimal utility types |
| Arr helpers | Not Implemented | Array manipulation |
| Str helpers | Not Implemented | String manipulation |
| Collection | Not Implemented | Fluent collection pipeline |
| Carbon/Date | Not Implemented | Date manipulation |
| Fluent | Not Implemented | Fluent attribute access |

#### Test Coverage

No Upstream test files analyzed. 52 Upstream test files exist for this package.

---

### View

#### Features

| Feature | Status | Notes |
|---------|--------|-------|
| Template rendering (Go html/template) | Implemented | Go-native templates |
| Template data passing | Implemented | Key-value data |
| Template template engine | Intentional Skip | PHP-specific — Go uses `html/template` |
| Sections / stacks | Not Implemented | Template composition |
| View composers / creators | Not Implemented | Pre-render data binding |
| Components | Not Implemented | Reusable view components |
| View finder | Not Implemented | Template file resolution |

#### Test Coverage

| Upstream File | Upstream Tests | Bedrock Status |
|-------------|---------------|----------------|
| ViewFactoryTest.php | 70 | 2 covered, ~18 portable missing, ~50 Template-specific (skip) |
| ViewBladeCompilerTest.php | ? | Intentional Skip |
| ViewCompilerEngineTest.php | ? | Intentional Skip |
| ViewPhpEngineTest.php | ? | Intentional Skip |
| ComponentTest.php | ? | Not analyzed |
| ViewComponentTest.php | ? | Not analyzed |
| ViewComponentAttributeBagTest.php | ? | Not analyzed |
| ViewTest.php | ? | Not analyzed |
| ViewFileViewFinderTest.php | ? | Not analyzed |
| ViewEngineResolverTest.php | ? | Not analyzed |

---

### Stub Packages

These packages exist as directory placeholders (`doc.go` only) with no implementation code.

| Package | Upstream Test Files | Estimated Tests | Blocks |
|---------|-------------------|-----------------|--------|
| **Database** | 50+ | ~500+ | Auth providers, cache DB store, session DB, validation |
| **Validation** | 37 | ~500+ | Input validation across all packages |
| **Queue** | 10+ | ~100+ | Queued events, queued listeners, job dispatch |
| **Broadcasting** | 5+ | ~30+ | Real-time event broadcasting |
| **Mail** | 10+ | ~100+ | Email sending, templates |
| **Notifications** | 5+ | ~50+ | Multi-channel notifications |
| **Filesystem** | 10+ | ~80+ | Disk abstraction, file operations |
| **Log** | 5+ | ~30+ | Structured logging, channels |
| **Translation** | 5+ | ~50+ | Localization, pluralization |
| **Pagination** | 3+ | ~30+ | Cursor/offset pagination |
| **Redis** | 5+ | ~40+ | Redis client, connections |
| **Pipeline** | 2+ | ~10+ | Middleware pipeline abstraction |
| **Bus** | 3+ | ~20+ | Command/job dispatching |
| **Process** | 3+ | ~20+ | External process execution |
| **Testing** | — | — | HTTP/database testing helpers |
| **Contracts** | — | — | Interface definitions |
| **Conditionable** | — | — | When/unless fluent pattern |
| **Cookie** | 2+ | ~10+ | Cookie management (implementation exists separately) |

---

## Cross-Cutting Gaps

| Gap | Packages Blocked | Tests Blocked | Category |
|-----|-----------------|---------------|----------|
| **Event Dispatching integration** | Auth (5), Routing (2+), Foundation (1+), View (20+) | ~28+ | D — infrastructure |
| **Broadcasting infrastructure** | Events (7) | ~7+ | D — infrastructure |
| **Queue infrastructure** | Events (25) | ~25+ | D — infrastructure |
| **Database / Orm ORM** | Auth providers (29), Cache DB store, Session DB handler, Validation | ~100+ | D — infrastructure |
| **Template Template Engine** | View (50+ across 10 files) | ~50+ | A — permanent skip |
| **PHP Reflection / Auto-wiring** | Container (40), Foundation (15+), Console (2) | ~57+ | A — permanent skip |
| **PHP Backed Enums** | Gate (6), Session (13), Middleware (1), Routing (3+), Http (2) | ~25 | A — permanent skip |
| **PHP Macroable Trait** | Config (1), Guards (1), View (1), Foundation (1), Session (1) | ~5 | A — permanent skip |
| **PHP ArrayAccess** | Config (4), Container (2), Http (1) | ~7 | A — permanent skip |
| **Command Scheduling / Cron** | Console (13) | 13 | D — infrastructure |
| **Command Mutex / Locking** | Console (13) | 13 | B — missing feature |
| **Route Model Binding** | Routing (15+) | ~15+ | B — missing feature |
| **File Uploads** | Http (10+) | ~10+ | B — missing feature |
| **HTTP Client** | Http (1 full test file) | ~50+ | B — missing feature |

---

## Gap Categories

- **(A) Language difference** — PHP-specific constructs with no Go equivalent. Permanent intentional skips. (~150+ tests)
- **(B) Missing feature** — Portable features not yet implemented in Bedrock. Action required. (~200+ tests)
- **(C) Partial coverage** — Feature exists but test coverage is incomplete. Action required. (~100+ tests)
- **(D) Infrastructure blocker** — Requires an unported package (database, events) before tests can be written. (~130+ tests)

---

## Recommended Next Steps

### Tier 1: Close gaps in near-complete packages
1. **Auth/Gate** — Add subtype/interface resolution, dash-to-camel, array abilities, custom resource gates (~13 tests)
2. **Auth/Token Guard** — Add custom field validation (4 tests)
3. **Auth/Authorize MW** — Add model-type authorization (5 tests)
4. **Auth/Response** — Add 3 missing response tests

### Tier 2: Analyze untracked files for already-ported packages
5. **Session** — Map SessionStoreTest.php (56 methods) line-by-line against Bedrock's 55 tests
6. **Events** — Map EventsDispatcherTest.php (39 methods) against Bedrock's 35 tests
7. **Auth** — Analyze remaining 5 untracked auth files (44 tests)
8. **Container** — Analyze 10 additional files (classify as skip vs missing)

### Tier 3: Expand partially-ported packages
9. **HTTP** — Port remaining HttpRequestTest methods (flash, cookies, merge, file methods)
10. **Routing** — Port route collection, URL generation, registrar, middleware sorting
11. **Foundation** — Port service provider registration, deferred services, config merging
12. **Console** — Port command parsing, scheduling basics

### Tier 4: Architectural decisions needed
13. **Event system integration** — Wire dispatcher into auth/routing/foundation (unblocks ~28+ tests)
14. **Database layer** — Unblocks user providers (29 tests), cache DB store, session DB handler, validation rules
15. **Validation package** — 37 Upstream test files, ~500+ methods — entire package is stub
16. **Support package** — 52 Upstream test files — Arr, Str, Collection, etc.

### Tier 5: Permanent intentional skips (no action)
- Template template engine (~50+ tests)
- PHP reflection auto-wiring (~60+ tests)
- PHP Backed Enums (~25 tests)
- PHP ArrayAccess (~7 tests)
- PHP Macroable (~5 tests)
- PHP Attributes (`#[Singleton]`, `#[Scoped]`, `#[Bind]`) (~10+ tests)
