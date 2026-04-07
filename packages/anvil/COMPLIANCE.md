# Upstream 13.x Test Compliance Report

> Generated: 2026-04-07
> Baseline: [upstream/framework 13.x tests](https://github.com/upstream/framework/tree/13.x/tests)
> All Bedrock tests pass (`go test ./...` green across all packages).

---

## Executive Summary

| Metric | Count |
|--------|-------|
| Ported packages with tests | 14 |
| Stub packages (no code) | 19 |
| **Bedrock Go tests (total)** | **555** |
| Upstream test files in scope (ported packages) | **210** |
| Upstream test files analyzed line-by-line | 23 |
| Upstream test files not yet analyzed | 187 |
| Upstream test methods analyzed | 615 |
| Analyzed: COVERED | 351 |
| Analyzed: MISSING | 135 |
| Analyzed: INTENTIONAL-SKIP | 78 |
| Analyzed: BEDROCK-ONLY | 75 |
| **Coverage of analyzed tests (excl. skips)** | **72%** |
| **Estimated full Upstream test surface (ported packages)** | **~2,500+** |

### What This Means

Bedrock has **555 passing Go tests** across 14 packages. For the subset of Upstream tests we've analyzed line-by-line (615 methods across 23 files), we cover 72% of the portable surface. However, the full Upstream test suite for these same packages spans **~210 test files** with an estimated **~2,500+ test methods** — most of which have not yet been compared.

Four packages have **100% portable coverage**: Encryption, Config, Hashing, Container.

---

## Per-Package Scorecard

### Fully Covered (100% portable)

| Package | Upstream File | Upstream Tests | Bedrock Tests | Portable Coverage | Intentional Skips |
|---------|-------------|---------------|---------------|-------------------|-------------------|
| **Encryption** | EncrypterTest.php (1 file) | 26 | 26 + 6 extra | **100%** | 0 |
| **Config** | RepositoryTest.php (1 file) | 33 | 27 + 12 extra | **100%** | 6 (ArrayAccess, macros, Collection) |
| **Hashing** | HasherTest.php (1 file) | 13 | 9 + 2 extra | **100%** | 3 (Go crypto always available) |
| **Container** | ContainerTest.php (1 of 11 files) | 77 | 59 + 22 extra | **100%** | 40 (PHP reflection, attributes, ArrayAccess) |

**Container note:** 10 additional Upstream test files exist (ContainerCallTest, ContainerExtendTest, ContainerTaggingTest, ContextualBindingTest, ContextualAttributeBindingTest, ResolvingCallbackTest, AfterResolvingAttributeCallbackTest, ContainerResolveNonInstantiableTest, RewindableGeneratorTest, UtilTest). These are primarily PHP reflection and attribute-based — most are expected intentional-skips for Go. Bedrock covers tagging, contextual binding, and resolving callbacks via its own test files.

---

### Auth Subsystem (203 Bedrock tests)

| Component | Upstream Files | Upstream Tests | Covered | Missing | Skip | Bedrock-Only | Portable % |
|-----------|--------------|---------------|---------|---------|------|--------------|------------|
| **Gate** | AuthAccessGateTest (92) | 92 | 48 | 13 | 14 | 20 | **79%** |
| **Response** | AuthAccessResponseTest (13) | 13 | 7 | 3 | 0 | - | **70%** |
| **Session Guard** | AuthGuardTest (43) | 43 | 33 | 9 | 1 | 31 | **79%** |
| **Token Guard** | AuthTokenGuardTest (14) | 14 | 10 | 4 | 0 | - | **71%** |
| **Password Broker** | AuthPasswordBrokerTest (9) | 9 | 8 | 0 | 1 | 10 | **100%** |
| **Token Repository** | AuthDatabaseTokenRepositoryTest (10) | 10 | 10 | 0 | 0 | - | **100%** |
| **Authenticate MW** | AuthenticateMiddlewareTest (11) | 11 | 7 | 2 | 2 | - | **78%** |
| **Authorize MW** | AuthorizeMiddlewareTest (16) | 16 | 8 | 5 | 2 | 3 | **62%** |
| **Email Verify MW** | EnsureEmailIsVerifiedTest (1) | 1 | 0 | 0 | 1 | 3 | N/A (skip) |
| **Redirect MW** | RedirectIfAuthenticatedMiddlewareTest (1) | 1 | 0 | 0 | 1 | 2 | N/A (skip) |
| **Subtotal (analyzed)** | **10 of 15 files** | **210** | **131** | **36** | **22** | **69** | **78%** |

**Auth files NOT yet analyzed (5 files, ~44 methods):**

| Upstream File | Tests | Bedrock Status |
|-------------|-------|----------------|
| AuthenticatableTest.php | 3 | Likely covered — Bedrock implements `Authenticatable` interface |
| AuthHandlesAuthorizationTest.php | 5 | Partially covered — tests allow/deny/status helpers on traits |
| AuthorizesResourcesTest.php | 6 | MISSING — tests resource authorization method mapping |
| AuthDatabaseUserProviderTest.php | 14 | MISSING — requires database layer (Category D) |
| AuthOrmUserProviderTest.php | 15 | MISSING — requires Orm ORM (Category D) |
| AuthListenersSendEmailVerification...Test.php | ~1 | MISSING — requires event/notification system |

**Key missing behaviors:**
- Gate: subtype/interface resolution, dash-to-camel conversion, class-name policies, array abilities, custom resource gates
- Guards: `logoutCurrentDevice`, cookie override on remember, event firing (5 tests blocked by missing event system)
- Token Guard: custom field validation (4 tests)
- Authorize MW: model-type authorization (5 tests)
- User Providers: entire database/Orm provider layer (29 tests, blocked by database package)

---

### Session (55 Bedrock tests)

| Upstream File | Tests | Bedrock Status |
|-------------|-------|----------------|
| SessionStoreTest.php | 56 | **Partially covered** — Bedrock has 55 session tests but not mapped 1:1 |
| ArraySessionHandlerTest.php | 10 | **Covered** — Bedrock has array handler tests |
| CacheBasedSessionHandlerTest.php | ? | MISSING — requires cache-backed session handler |
| EncryptedSessionStoreTest.php | ? | MISSING — requires encrypted session store |
| FileSessionHandlerTest.php | ? | MISSING — requires file-based session handler |

**Key gaps:** PHP Backed Enum key operations (13 tests in SessionStoreTest are PHP-enum-specific — intentional skip). CacheBasedSessionHandler, EncryptedSessionStore, and FileSessionHandler are not implemented.

---

### Cache (33 Bedrock tests)

| Upstream File | Tests | Bedrock Status |
|-------------|-------|----------------|
| CacheArrayStoreTest.php | 31 | **Partially covered** — Bedrock has 33 cache tests for ArrayStore |
| CacheRepositoryTest.php | ? | NOT ANALYZED |
| CacheManagerTest.php | ? | NOT ANALYZED |
| CacheEventsTest.php | ? | NOT ANALYZED — requires event system |
| CacheRateLimiterTest.php | ? | NOT ANALYZED |
| +16 more store-specific files | ? | NOT ANALYZED — stores not implemented (Redis, Memcached, Database, File, DynamoDB, etc.) |

**Key gaps:** Only ArrayStore is implemented. 20 of 21 Upstream cache test files are for unimplemented stores/features. Lock support exists in Bedrock tests.

---

### Events (35 Bedrock tests)

| Upstream File | Tests | Bedrock Status |
|-------------|-------|----------------|
| EventsDispatcherTest.php | 39 | **Mostly covered** — 35 Bedrock tests cover core dispatch |
| EventsSubscriberTest.php | ? | NOT ANALYZED |
| BroadcastedEventsTest.php | ? | NOT ANALYZED — requires broadcasting |
| QueuedEventsTest.php | ? | NOT ANALYZED — requires queue system |

**Key gaps:** Event subscribers, broadcasting integration, queued events. The Bedrock event dispatcher covers basic listener registration, halting, wildcard listeners, and dispatch. Container resolution of handlers and deferred events are likely missing.

---

### Routing (22 Bedrock tests)

| Upstream File | Tests | Bedrock Status |
|-------------|-------|----------------|
| RoutingRouteTest.php | 102 | **13 covered, ~89 missing** |
| RouteBindingTest.php | ? | NOT ANALYZED — route model binding not implemented |
| RouteCollectionTest.php | ? | NOT ANALYZED |
| RouteRegistrarTest.php | ? | NOT ANALYZED |
| RoutingUrlGeneratorTest.php | ? | NOT ANALYZED |
| ImplicitRouteBindingTest.php | ? | NOT ANALYZED |
| RoutingSortedMiddlewareTest.php | ? | NOT ANALYZED |
| RoutingRedirectorTest.php | ? | NOT ANALYZED |
| RouteActionTest.php | ? | NOT ANALYZED |
| RouteSignatureParametersTest.php | ? | NOT ANALYZED |
| RouteUriTest.php | ? | NOT ANALYZED |

**Covered:** Basic dispatch, middleware groups, nested groups, resource routing, named routes, Patch/Options dispatch.
**Major missing:** Route model binding, controller routing, domain routing, signed routes, redirects, URL generation, pattern filtering, implicit bindings, route caching, fluent routing.

---

### HTTP (32 Bedrock tests)

| Upstream File | Tests | Bedrock Status |
|-------------|-------|----------------|
| HttpRequestTest.php | 113 | **26 covered, ~87 missing** |
| HttpResponseTest.php | ? | NOT ANALYZED |
| HttpJsonResponseTest.php | ? | NOT ANALYZED |
| HttpClientTest.php | ? | NOT ANALYZED — HTTP client not implemented |
| HttpRedirectResponseTest.php | ? | NOT ANALYZED |
| JsonResourceTest.php | ? | NOT ANALYZED |
| HttpUploadedFileTest.php | ? | NOT ANALYZED — file uploads not implemented |
| HttpMimeTypeTest.php | ? | NOT ANALYZED |
| HttpTestingFileFactoryTest.php | ? | NOT ANALYZED |

**Covered:** Input/query/boolean/integer extraction, path/URL/host/scheme/method, content type, JSON detection, accepts/prefers, all/has/missing/only/except, header, bearer token, IP.
**Major missing:** File uploads, cookies, session integration, flash data, fingerprinting, merge/replace, old input, fluent/string/date/enum methods, JSON request body, HTTP client.

---

### View (10 Bedrock tests)

| Upstream File | Tests | Bedrock Status |
|-------------|-------|----------------|
| ViewFactoryTest.php | 70 | **2 covered, ~18 portable missing, ~50 Template-specific (skip)** |
| ViewBladeCompilerTest.php | ? | INTENTIONAL-SKIP — Go uses html/template |
| ViewCompilerEngineTest.php | ? | INTENTIONAL-SKIP |
| ComponentTest.php | ? | NOT ANALYZED |
| ViewComponentTest.php | ? | NOT ANALYZED |
| ViewComponentAttributeBagTest.php | ? | NOT ANALYZED |
| ViewTest.php | ? | NOT ANALYZED |
| ViewFileViewFinderTest.php | ? | NOT ANALYZED |
| ViewEngineResolverTest.php | ? | NOT ANALYZED |
| ViewPhpEngineTest.php | ? | INTENTIONAL-SKIP |

**Key note:** Bedrock uses Go `html/template` instead of Template. Template compiler tests (~50+) are permanent intentional skips. Portable view factory features (sections, stacks, composers, creators, loops) are mostly missing.

---

### Foundation (20 Bedrock tests)

| Upstream File | Tests | Bedrock Status |
|-------------|-------|----------------|
| FoundationApplicationTest.php | 42 | **6 covered, ~36 missing** |
| +15 more files | ? | NOT ANALYZED |

**Covered:** Environment detection, boot/booted/termination callbacks.
**Major missing:** Service providers, deferred services, config merging, namespace resolution, cache paths, route/event caching, alias loader, bootstrapping. Most require deeper IoC integration.

---

### Console (9 Bedrock tests)

| Upstream File | Tests | Bedrock Status |
|-------------|-------|----------------|
| ConsoleApplicationTest.php | 14 | **2 covered, ~10 missing** |
| +10 more files | ? | NOT ANALYZED |

**Covered:** Command execution, command registration.
**Major missing:** Interactive prompts, IoC command resolution, PHP attributes, command scheduling, signals, command mutex, output styling.

---

## Cross-Cutting Systemic Gaps

| Gap | Packages Blocked | Tests Blocked | Category |
|-----|-----------------|---------------|----------|
| **Event Dispatching integration** | Auth (5), Routing (2+), Foundation (1+), View (20+) | ~28+ | D — infrastructure |
| **Database / Orm ORM** | Auth user providers (29), Cache DB store, Session DB/file handlers, Validation exists/unique rules | ~100+ | D — infrastructure |
| **Template Template Engine** | View (50+ across 10 files) | ~50+ | A — language difference (permanent skip) |
| **PHP Reflection / Auto-wiring** | Container (40), Foundation (15+), Console (5+) | ~60+ | A — language difference (permanent skip) |
| **PHP Backed Enums** | Gate (6), Session (13), Middleware (1), Routing (3+), Http (2) | ~25 | A — language difference |
| **PHP Macroable Trait** | Config (1), Guards (1), View (1), Foundation (1), Session (1) | ~5 | A — language difference |
| **PHP ArrayAccess** | Config (4), Container (2), Http (1) | ~7 | A — language difference |
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

## Unanalyzed Upstream Test Files

These files exist in Upstream 13.x for packages Bedrock has ported, but have NOT been compared line-by-line:

### Auth (5 files)
- AuthenticatableTest.php (3 tests)
- AuthHandlesAuthorizationTest.php (5 tests)
- AuthorizesResourcesTest.php (6 tests)
- AuthDatabaseUserProviderTest.php (14 tests)
- AuthOrmUserProviderTest.php (15 tests)

### Container (10 files)
- AfterResolvingAttributeCallbackTest.php
- ContainerCallTest.php
- ContainerExtendTest.php
- ContainerResolveNonInstantiableTest.php
- ContainerTaggingTest.php
- ContextualAttributeBindingTest.php
- ContextualBindingTest.php
- ResolvingCallbackTest.php
- RewindableGeneratorTest.php
- UtilTest.php

### Session (3 files)
- CacheBasedSessionHandlerTest.php
- EncryptedSessionStoreTest.php
- FileSessionHandlerTest.php

### Cache (20 files)
- CacheRepositoryTest.php, CacheManagerTest.php, CacheEventsTest.php, CacheFileStoreTest.php, CacheDatabaseStoreTest.php, CacheRedisStoreTest.php, CacheMemcachedStoreTest.php, CacheMemcachedConnectorTest.php, CacheDynamoDbStoreTest.php, CacheNullStoreTest.php, CacheMemoizedStoreTest.php, CacheRateLimiterTest.php, CacheTaggedCacheTest.php, CacheSessionStoreTest.php, CacheSpyMemoTest.php, ClearCommandTest.php, ConcurrencyLimiterTest.php, LimitTest.php, RateLimiterTest.php, CacheApcStoreTest.php

### Events (3 files)
- EventsSubscriberTest.php, BroadcastedEventsTest.php, QueuedEventsTest.php

### Routing (11 files)
- ImplicitRouteBindingTest.php, RouteActionTest.php, RouteBindingTest.php, RouteCollectionTest.php, RouteRegistrarTest.php, RouteSignatureParametersTest.php, RouteUriTest.php, RoutingRedirectorTest.php, RoutingSortedMiddlewareTest.php, RoutingUrlGeneratorTest.php

### HTTP (9 files)
- HttpClientTest.php, HttpJsonResponseTest.php, HttpMimeTypeTest.php, HttpRedirectResponseTest.php, HttpResponseTest.php, HttpTestingFileFactoryTest.php, HttpUploadedFileTest.php, JsonResourceTest.php

### View (9 files)
- ComponentTest.php, ViewBladeCompilerTest.php, ViewCompilerEngineTest.php, ViewComponentAttributeBagTest.php, ViewComponentTest.php, ViewEngineResolverTest.php, ViewFileViewFinderTest.php, ViewPhpEngineTest.php, ViewTest.php

### Foundation (15 files)
- FoundationAliasLoaderTest.php, FoundationApplicationBuilderTest.php, FoundationAuthenticationTest.php, FoundationAuthorizesRequestsTraitTest.php, FoundationCacheBasedMaintenanceModeTest.php, FoundationDocsCommandTest.php, FoundationEnvironmentDetectorTest.php, FoundationExceptionsHandlerTest.php, FoundationFormRequestTest.php, FoundationHelpersTest.php, FoundationInteractsWithDatabaseTest.php, FoundationInteractsWithTimeTest.php, FoundationPackageManifestTest.php, FoundationProviderRepositoryTest.php, FoundationViteTest.php

### Console (10 files)
- CacheCommandMutexTest.php, CommandMutexTest.php, CommandTest.php, CommandTrapTest.php, ConfiguresPromptsTest.php, ConsoleEventSchedulerTest.php, ConsoleParserTest.php, ConsoleScheduledEventTest.php, OutputStyleTest.php, SignalsTest.php

### Support (52 files) — NOT PORTED
### Validation (37 files) — NOT PORTED

---

## Summary by Package

| Package | Bedrock Tests | Upstream Files (total) | Files Analyzed | Analyzed Coverage | Key Gaps |
|---------|--------------|----------------------|----------------|-------------------|----------|
| Encryption | 26 | 1 | 1 | **100%** | None |
| Config | 37 | 1 | 1 | **100%** (portable) | None |
| Hashing | 9 | 1 | 1 | **100%** (portable) | None |
| Container | 59 | 11 | 1 | **100%** (portable) | 10 files unanalyzed (mostly PHP-specific) |
| Auth (all) | 203 | 15 | 10 | **78%** (portable) | User providers (DB), event firing, model authorization |
| Session | 55 | 5 | 2 | ~80% est. | Encrypted/cache/file handlers |
| Cache | 33 | 21 | 1 | ~70% est. (ArrayStore) | 20 files for unimplemented stores |
| Events | 35 | 4 | 1 | ~75% est. | Subscribers, queued events, broadcasting |
| Routing | 22 | 12 | 1 | **~13%** | Model binding, controllers, URL gen, 11 files |
| HTTP | 32 | 10 | 1 | **~23%** | Uploads, cookies, sessions, 9 files |
| View | 10 | 11 | 1 | **~10%** | Template (skip), sections, composers, 9 files |
| Foundation | 20 | 16 | 1 | **~15%** | Service providers, deferred services, 15 files |
| Console | 9 | 11 | 1 | **~15%** | Scheduling, prompts, signals, 10 files |
| Support | 5 | 52 | 0 | **~1%** | Entire package mostly unported |

---

## Recommended Next Steps (Priority Order)

### Tier 1: Close gaps in near-complete packages
1. **Auth/Gate** — Add subtype/interface resolution, dash-to-camel, array abilities, custom resource gates (~13 tests)
2. **Auth/Token Guard** — Add custom field validation (4 tests)
3. **Auth/Authorize MW** — Add model-type authorization (5 tests)
4. **Auth/Response** — Add 3 missing response tests

### Tier 2: Analyze untracked files for already-ported packages
5. **Session** — Map SessionStoreTest.php (56 methods) line-by-line against Bedrock's 55 tests
6. **Cache** — Map CacheArrayStoreTest.php (31 methods) line-by-line against Bedrock's 33 tests
7. **Events** — Map EventsDispatcherTest.php (39 methods) against Bedrock's 35 tests
8. **Auth** — Analyze remaining 5 untracked auth files (44 tests)
9. **Container** — Analyze 10 additional files (classify as skip vs missing)

### Tier 3: Expand partially-ported packages
10. **HTTP** — Port remaining HttpRequestTest methods (flash, cookies, merge, file methods)
11. **Routing** — Port route collection, URL generation, registrar, middleware sorting
12. **Foundation** — Port service provider registration, deferred services, config merging
13. **Console** — Port command parsing, scheduling basics

### Tier 4: Architectural decisions needed
14. **Event system integration** — Wire dispatcher into auth/routing/foundation (unblocks ~28+ tests)
15. **Database layer** — Unblocks user providers (29 tests), cache DB store, session DB handler, validation rules
16. **Validation package** — 37 Upstream test files, ~500+ methods — entire package is stub
17. **Support package** — 52 Upstream test files — Arr, Str, Collection, etc.

### Tier 5: Permanent intentional skips (no action)
- Template template engine (~50+ tests)
- PHP reflection auto-wiring (~60+ tests)
- PHP Backed Enums (~25 tests)
- PHP ArrayAccess (~7 tests)
- PHP Macroable (~5 tests)
- PHP Attributes (#[Singleton], #[Scoped], #[Bind]) (~10+ tests)
