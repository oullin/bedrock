<!--
  Generated derivative report. Source of truth for status:
    services/compliance/docs-status.yml
    services/compliance/docs-inventories/laravel-docs.txt
  See "How to regenerate" below for the awk/python recipe that reproduces these
  numbers from the inputs above.
-->

# Documentation Porting Gaps

A prioritized view of every Laravel 13.x documentation section that is **not yet
classified** in `services/compliance/docs-status.yml` (i.e. neither `ported`,
`adapted`, nor `excluded`). Every entry below is therefore a real gap to either
adapt into a Bedrock doc or explicitly exclude.

This file is a derivative of the same inputs the
`services/scripts/laravel-compliance.sh` summary uses; it adds per-section
detail and priority bucketing on top of the file-level table that
`report.md → Documentation Detail → Top Pending Documentation Files` already
prints.

## Baseline

| | |
| --- | --- |
| Generated | 2026-04-27 |
| Upstream | `laravel/docs@13.x` pinned at `bfe893709` (`services/compliance/sources.lock.json`) |
| Inventory | `services/compliance/docs-inventories/laravel-docs.txt` — 875 sections |
| Status file | `services/compliance/docs-status.yml` |

## Headline Numbers

| Status | Sections | Share |
| --- | ---: | ---: |
| Ported (verbatim) | 0 | 0.0% |
| Adapted | 330 | 37.7% |
| Excluded (PHP-only / out of product scope) | 345 | 39.4% |
| **Pending — gap to close** | **200** | **22.9%** |
| Total | 875 | 100.0% |

The 200 pending sections are spread across **56 Laravel doc files**, every one
of which already has a corresponding Bedrock target page under
`services/docs/`. **No new top-level page needs to be created** to close any of
these gaps — every gap is an "expand the existing target doc" task.

## Priority Taxonomy

| Bucket | Definition |
| --- | --- |
| **P0 — Foundational** | Topics any Bedrock app touches on day one: HTTP request handling, routing, middleware, controllers, validation, container, configuration, lifecycle, database/queries/migrations, testing, auth. Closing these gaps compounds because every other doc cross-links into them. |
| **P1 — Widely used** | Capability layers that most non-trivial apps wire up: queues, events, cache, mail, notifications, sessions, logging, hashing, encryption, filesystem, pagination, collections, strings, helpers, http-client, redis, scheduling, broadcasting, contracts, localization. |
| **P2 — Specialized** | Optional products and ecosystem packages: Spark/billing, Passport, Socialite, Fortify, Horizon, Telescope, Pennant, Precognition, Prompts, Reverb, Scout, Concurrency, Envoy, Boost/AI SDK/MCP. |

## P0 — Foundational (18 files, 60 sections)

| Laravel doc | Bedrock target | Pending | Notes |
| --- | --- | ---: | --- |
| `authentication.md` | `services/docs/packages/auth.md` | 4 | Custom guards, custom user providers, page-level slug, password confirmation. |
| `authorization.md` | `services/docs/packages/auth.md` | 2 | Page-level slug + `gates`. Bedrock policies are partly adapted; gates are not yet captured. |
| `configuration.md` | `services/docs/packages/config.md` | 2 | Page-level slug + accessor surface. |
| `container.md` | `services/docs/packages/container.md` | 4 | Binding, resolving, container events, page-level slug. |
| `controllers.md` | `services/docs/basics/controllers.md` | 5 | Writing controllers, controller middleware, DI in controllers, resource controllers, page-level slug. |
| `csrf.md` | `services/docs/basics/csrf.md` | 4 | Page-level slug, prevention, `X-CSRF-TOKEN`, `X-XSRF-TOKEN` headers. |
| `database.md` | `services/docs/packages/database.md` | 3 | Getting started, transactions, running raw SQL. |
| `eloquent.md` | `services/docs/packages/database.md` | 1 | Getting-started slug for the ORM-equivalent surface. |
| `lifecycle.md` | `services/docs/concepts/request-lifecycle.md` | 3 | Lifecycle overview, request lifecycle, focus on service providers. |
| `middleware.md` | `services/docs/basics/middleware.md` | 5 | Defining, registering, parameters, terminable, page-level slug. |
| `migrations.md` | `services/docs/packages/database.md` | 1 | Page-level slug for the migration surface. |
| `queries.md` | `services/docs/packages/database.md` | 1 | Page-level slug for the query builder. |
| `requests.md` | `services/docs/packages/httpx.md` | 3 | Page-level slug, interacting with the request, files. |
| `responses.md` | `services/docs/packages/httpx.md` | 3 | Page-level slug, creating responses, redirects. |
| `routing.md` | `services/docs/packages/routing.md` | 6 | Page-level slug, basic routing, named routes, route groups, model binding, rate limiting. |
| `testing.md` | `services/docs/concepts/testing.md` | 3 | Getting-started slug, creating tests, running tests. |
| `urls.md` | `services/docs/basics/url-generation.md` | 3 | URL generation slug, default values, URLs for named routes. |
| `validation.md` | `services/docs/packages/validation.md` | 7 | Page-level slug, quickstart, manual validators, available rules, custom rules, validating arrays, working with error messages. |

## P1 — Widely used (21 files, 83 sections)

| Laravel doc | Bedrock target | Pending | Notes |
| --- | --- | ---: | --- |
| `artisan.md` | `services/docs/packages/console.md` | 1 | `writing-commands` slug. |
| `broadcasting.md` | `services/docs/packages/broadcasting.md` | 8 | Largest single-file gap. Page-level slug, server-side install, defining events, broadcasting events, authorizing channels, presence channels, client events, receiving broadcasts. |
| `cache.md` | `services/docs/packages/cache.md` | 4 | Page-level slug, atomic locks, cache tags, custom drivers. |
| `collections.md` | `services/docs/packages/collection.md` | 3 | Page-level slug, available methods, lazy collections. |
| `contracts.md` | `services/docs/packages/contracts.md` | 4 | Page-level slug, when/how to use, contract reference. |
| `encryption.md` | `services/docs/packages/encryption.md` | 3 | Page-level slug, configuration, using the encrypter. |
| `events.md` | `services/docs/packages/events.md` | 5 | Page-level slug, registering, dispatching, queued listeners, subscribers. |
| `filesystem.md` | `services/docs/packages/filesystem.md` | 6 | File storage slug, storing/retrieving/deleting files, directories, custom filesystems. |
| `hashing.md` | `services/docs/packages/hashing.md` | 4 | Page-level slug, configuration, basic usage, hash-algorithm verification. |
| `helpers.md` | `services/docs/packages/support.md` | 3 | Page-level slug, arrays/objects helpers, other utilities. |
| `http-client.md` | `services/docs/packages/httpx.md` | 2 | Page-level slug + making requests. The Go client surface lives inside `httpx`. |
| `localization.md` | `services/docs/packages/translation.md` | 3 | Page-level slug, defining/retrieving translation strings. |
| `logging.md` | `services/docs/packages/log.md` | 5 | Page-level slug, configuration, building log stacks, writing log messages, tailing via `pail`. |
| `mail.md` | `services/docs/packages/mailx.md` | 4 | Page-level slug, sending mail, custom transports, testing. |
| `notifications.md` | `services/docs/packages/notifications.md` | 5 | Page-level slug, sending, database notifications, custom channels, testing. |
| `pagination.md` | `services/docs/packages/pagination.md` | 4 | Basic usage, database pagination, cursor pagination, displaying results. |
| `queues.md` | `services/docs/packages/queue.md` | 7 | Page-level slug, dispatching jobs, running the worker, dealing with failed jobs, job middleware, job events, testing. |
| `redis.md` | `services/docs/packages/redis.md` | 4 | Page-level slug, configuration, interacting with Redis, pub/sub. |
| `scheduling.md` | `services/docs/packages/console.md` | 1 | `defining-schedules` slug. |
| `session.md` | `services/docs/packages/session.md` | 4 | HTTP session slug, interacting, custom drivers, session cache integration. |
| `strings.md` | `services/docs/packages/str.md` | 3 | Page-level slug, available methods, fluent strings. |

## P2 — Specialized (17 files, 57 sections)

| Laravel doc | Bedrock target | Pending | Notes |
| --- | --- | ---: | --- |
| `ai-sdk.md` | `services/docs/packages/ai/sdk.md` | 5 | Agents, embeddings, events, page-level slug, testing. |
| `ai.md` | `services/docs/packages/ai/boost.md` | 1 | `laravel-boost` page-level slug. |
| `billing.md` | `services/docs/packages/spark.md` | 4 | Page-level slug, customers, subscriptions, checkout. |
| `boost.md` | `services/docs/packages/ai/boost.md` | 3 | Page-level slug, MCP server, agent skills. |
| `concurrency.md` | `services/docs/packages/concurrency.md` | 3 | Page-level slug, running concurrent tasks, deferring concurrent tasks. |
| `envoy.md` | `services/docs/packages/envoy.md` | 3 | Page-level slug, writing tasks, running tasks. |
| `fortify.md` | `services/docs/packages/fortify.md` | 2 | Page-level slug, two-factor authentication. |
| `horizon.md` | `services/docs/packages/horizon.md` | 5 | Page-level slug, running, tags, metrics, notifications. |
| `mcp.md` | `services/docs/packages/ai/mcp.md` | 4 | Page-level slug, creating servers, tools, resources. |
| `passport.md` | `services/docs/packages/passport.md` | 4 | Page-level slug, authorization-code grant, client-credentials grant, token scopes. |
| `pennant.md` | `services/docs/packages/pennant.md` | 3 | Page-level slug, defining features, checking features. |
| `precognition.md` | `services/docs/packages/precognition.md` | 3 | Page-level slug, live validation, testing. |
| `prompts.md` | `services/docs/packages/prompts.md` | 4 | Page-level slug, available prompts, unsupported-environment fallbacks, testing. |
| `reverb.md` | `services/docs/packages/reverb.md` | 3 | Page-level slug, running the server, events. |
| `scout.md` | `services/docs/packages/scout.md` | 3 | Page-level slug, searching, custom engines. |
| `socialite.md` | `services/docs/packages/socialite.md` | 3 | Page-level slug, authentication, retrieving user details. |
| `telescope.md` | `services/docs/packages/telescope.md` | 4 | Page-level slug, available watchers, filtering, tagging. |

## How to Close a Gap

For each pending entry `<file>.md#<slug>` you intend to address:

1. **Decide the disposition.** Either the section maps to existing Bedrock
   functionality (→ `adapted`) or it does not (→ `excluded`). Verbatim port
   (`ported`) is reserved for unchanged textual structure and is currently
   unused; new closures should default to `adapted` unless a clear PHP-only
   exclusion already has precedent in `docs-status.yml`.
2. **Write or expand** the target Bedrock page identified above. Standard
   section order in Bedrock package docs: Overview / Source Coverage / Core
   Concepts / Public Surface / Capability Matrix / Usage / Configuration /
   API Reference. Embed the upstream pointer as an HTML comment inside the
   relevant section: `<!-- laravel-docs: <file>.md#<slug> -->`.
3. **Append the slug** to the `adapted:` or `excluded:` list in
   `services/compliance/docs-status.yml`. Group it under the relevant
   commented block (the file is organized by source page).
4. **Regenerate** `services/compliance/report.md` via the existing pipeline
   (`services/scripts/laravel-compliance.sh`); the
   "Top Pending Documentation Files" table and the headline classified counts
   will move accordingly.

## How to Regenerate This Report

The 200-section pending set is computed from two files only. Reproduce it
with:

```bash
python3 - <<'PY'
import yaml
status = yaml.safe_load(open('services/compliance/docs-status.yml'))
classified = (set(status.get('ported') or [])
              | set(status.get('adapted') or [])
              | set(status.get('excluded') or []))
inv = [l.strip() for l in open('services/compliance/docs-inventories/laravel-docs.txt')
       if l.strip() and not l.startswith('#')]
pending = sorted(set(inv) - classified)
print(f'pending={len(pending)} of {len(inv)}')
for s in pending:
    print(s)
PY
```

The same set is produced by the `pending_plain_inventory_entries` awk function
in `services/scripts/laravel-compliance.sh:621`. The priority bucketing in the
P0/P1/P2 tables above is editorial and lives only in this file; update it here
when the pending set shifts.

## Excluded Sections

For completeness: 345 inventory sections are explicitly out of scope and live
under `excluded:` in `services/compliance/docs-status.yml`. They cover Blade,
Views, Dusk, Sail, Valet, Octane, Folio, and per-package PHP-only sub-surfaces
(per-package `#events`, `#testing`, `#configuration` slugs that Bedrock
collapses into typed Go options). No work is required against these unless
product scope changes.

## Appendix: Full Pending List

Verbatim, grouped by upstream file. Use this as the checklist when working
through a closure batch.

### `ai-sdk.md`  — 5 pending

- `ai-sdk.md#agents`
- `ai-sdk.md#embeddings`
- `ai-sdk.md#events`
- `ai-sdk.md#laravel-ai-sdk`
- `ai-sdk.md#testing`

### `ai.md`  — 1 pending

- `ai.md#laravel-boost`

### `artisan.md`  — 1 pending

- `artisan.md#writing-commands`

### `authentication.md`  — 4 pending

- `authentication.md#adding-custom-guards`
- `authentication.md#adding-custom-user-providers`
- `authentication.md#authentication`
- `authentication.md#password-confirmation`

### `authorization.md`  — 2 pending

- `authorization.md#authorization`
- `authorization.md#gates`

### `billing.md`  — 4 pending

- `billing.md#checkout`
- `billing.md#customers`
- `billing.md#laravel-cashier-stripe`
- `billing.md#subscriptions`

### `boost.md`  — 3 pending

- `boost.md#agent-skills`
- `boost.md#laravel-boost`
- `boost.md#mcp-server`

### `broadcasting.md`  — 8 pending

- `broadcasting.md#authorizing-channels`
- `broadcasting.md#broadcasting`
- `broadcasting.md#broadcasting-events`
- `broadcasting.md#client-events`
- `broadcasting.md#defining-broadcast-events`
- `broadcasting.md#presence-channels`
- `broadcasting.md#receiving-broadcasts`
- `broadcasting.md#server-side-installation`

### `cache.md`  — 4 pending

- `cache.md#adding-custom-cache-drivers`
- `cache.md#atomic-locks`
- `cache.md#cache`
- `cache.md#cache-tags`

### `collections.md`  — 3 pending

- `collections.md#available-methods`
- `collections.md#collections`
- `collections.md#lazy-collections`

### `concurrency.md`  — 3 pending

- `concurrency.md#concurrency`
- `concurrency.md#deferring-concurrent-tasks`
- `concurrency.md#running-concurrent-tasks`

### `configuration.md`  — 2 pending

- `configuration.md#accessing-configuration-values`
- `configuration.md#configuration`

### `container.md`  — 4 pending

- `container.md#binding`
- `container.md#container-events`
- `container.md#resolving`
- `container.md#service-container`

### `contracts.md`  — 4 pending

- `contracts.md#contract-reference`
- `contracts.md#contracts`
- `contracts.md#how-to-use-contracts`
- `contracts.md#when-to-use-contracts`

### `controllers.md`  — 5 pending

- `controllers.md#controller-middleware`
- `controllers.md#controllers`
- `controllers.md#dependency-injection-and-controllers`
- `controllers.md#resource-controllers`
- `controllers.md#writing-controllers`

### `csrf.md`  — 4 pending

- `csrf.md#csrf-protection`
- `csrf.md#preventing-csrf-requests`
- `csrf.md#x-csrf-token`
- `csrf.md#x-xsrf-token`

### `database.md`  — 3 pending

- `database.md#database-getting-started`
- `database.md#database-transactions`
- `database.md#running-sql-queries`

### `eloquent.md`  — 1 pending

- `eloquent.md#eloquent-getting-started`

### `encryption.md`  — 3 pending

- `encryption.md#configuration`
- `encryption.md#encryption`
- `encryption.md#using-the-encrypter`

### `envoy.md`  — 3 pending

- `envoy.md#laravel-envoy`
- `envoy.md#running-tasks`
- `envoy.md#writing-tasks`

### `events.md`  — 5 pending

- `events.md#dispatching-events`
- `events.md#event-subscribers`
- `events.md#events`
- `events.md#queued-event-listeners`
- `events.md#registering-events-and-listeners`

### `filesystem.md`  — 6 pending

- `filesystem.md#custom-filesystems`
- `filesystem.md#deleting-files`
- `filesystem.md#directories`
- `filesystem.md#file-storage`
- `filesystem.md#retrieving-files`
- `filesystem.md#storing-files`

### `fortify.md`  — 2 pending

- `fortify.md#laravel-fortify`
- `fortify.md#two-factor-authentication`

### `hashing.md`  — 4 pending

- `hashing.md#basic-usage`
- `hashing.md#configuration`
- `hashing.md#hash-algorithm-verification`
- `hashing.md#hashing`

### `helpers.md`  — 3 pending

- `helpers.md#arrays-objects`
- `helpers.md#helpers`
- `helpers.md#other-utilities`

### `horizon.md`  — 5 pending

- `horizon.md#laravel-horizon`
- `horizon.md#metrics`
- `horizon.md#notifications`
- `horizon.md#running-horizon`
- `horizon.md#tags`

### `http-client.md`  — 2 pending

- `http-client.md#http-client`
- `http-client.md#making-requests`

### `lifecycle.md`  — 3 pending

- `lifecycle.md#focus-on-service-providers`
- `lifecycle.md#lifecycle-overview`
- `lifecycle.md#request-lifecycle`

### `localization.md`  — 3 pending

- `localization.md#defining-translation-strings`
- `localization.md#localization`
- `localization.md#retrieving-translation-strings`

### `logging.md`  — 5 pending

- `logging.md#building-log-stacks`
- `logging.md#configuration`
- `logging.md#logging`
- `logging.md#tailing-log-messages-using-pail`
- `logging.md#writing-log-messages`

### `mail.md`  — 4 pending

- `mail.md#custom-transports`
- `mail.md#mail`
- `mail.md#sending-mail`
- `mail.md#testing`

### `mcp.md`  — 4 pending

- `mcp.md#creating-servers`
- `mcp.md#laravel-mcp`
- `mcp.md#resources`
- `mcp.md#tools`

### `middleware.md`  — 5 pending

- `middleware.md#defining-middleware`
- `middleware.md#middleware`
- `middleware.md#middleware-parameters`
- `middleware.md#registering-middleware`
- `middleware.md#terminable-middleware`

### `migrations.md`  — 1 pending

- `migrations.md#database-migrations`

### `notifications.md`  — 5 pending

- `notifications.md#custom-channels`
- `notifications.md#database-notifications`
- `notifications.md#notifications`
- `notifications.md#sending-notifications`
- `notifications.md#testing`

### `pagination.md`  — 4 pending

- `pagination.md#basic-usage`
- `pagination.md#cursor-pagination`
- `pagination.md#database-pagination`
- `pagination.md#displaying-pagination-results`

### `passport.md`  — 4 pending

- `passport.md#authorization-code-grant`
- `passport.md#client-credentials-grant`
- `passport.md#laravel-passport`
- `passport.md#token-scopes`

### `pennant.md`  — 3 pending

- `pennant.md#checking-features`
- `pennant.md#defining-features`
- `pennant.md#laravel-pennant`

### `precognition.md`  — 3 pending

- `precognition.md#live-validation`
- `precognition.md#precognition`
- `precognition.md#testing`

### `prompts.md`  — 4 pending

- `prompts.md#available-prompts`
- `prompts.md#prompts`
- `prompts.md#testing`
- `prompts.md#unsupported-environments-and-fallbacks`

### `queries.md`  — 1 pending

- `queries.md#database-query-builder`

### `queues.md`  — 7 pending

- `queues.md#dealing-with-failed-jobs`
- `queues.md#dispatching-jobs`
- `queues.md#job-events`
- `queues.md#job-middleware`
- `queues.md#queues`
- `queues.md#running-the-queue-worker`
- `queues.md#testing`

### `redis.md`  — 4 pending

- `redis.md#configuration`
- `redis.md#interacting-with-redis`
- `redis.md#pub-sub`
- `redis.md#redis`

### `requests.md`  — 3 pending

- `requests.md#files`
- `requests.md#http-requests`
- `requests.md#interacting-with-the-request`

### `responses.md`  — 3 pending

- `responses.md#creating-responses`
- `responses.md#http-responses`
- `responses.md#redirects`

### `reverb.md`  — 3 pending

- `reverb.md#events`
- `reverb.md#laravel-reverb`
- `reverb.md#running-the-server`

### `routing.md`  — 6 pending

- `routing.md#basic-routing`
- `routing.md#named-routes`
- `routing.md#rate-limiting`
- `routing.md#route-groups`
- `routing.md#route-model-binding`
- `routing.md#routing`

### `scheduling.md`  — 1 pending

- `scheduling.md#defining-schedules`

### `scout.md`  — 3 pending

- `scout.md#custom-engines`
- `scout.md#laravel-scout`
- `scout.md#searching`

### `session.md`  — 4 pending

- `session.md#adding-custom-session-drivers`
- `session.md#http-session`
- `session.md#interacting-with-the-session`
- `session.md#session-cache`

### `socialite.md`  — 3 pending

- `socialite.md#authentication`
- `socialite.md#laravel-socialite`
- `socialite.md#retrieving-user-details`

### `strings.md`  — 3 pending

- `strings.md#available-methods`
- `strings.md#fluent-strings`
- `strings.md#strings`

### `telescope.md`  — 4 pending

- `telescope.md#available-watchers`
- `telescope.md#filtering`
- `telescope.md#laravel-telescope`
- `telescope.md#tagging`

### `testing.md`  — 3 pending

- `testing.md#creating-tests`
- `testing.md#running-tests`
- `testing.md#testing-getting-started`

### `urls.md`  — 3 pending

- `urls.md#default-values`
- `urls.md#url-generation`
- `urls.md#urls-for-named-routes`

### `validation.md`  — 7 pending

- `validation.md#available-validation-rules`
- `validation.md#custom-validation-rules`
- `validation.md#manually-creating-validators`
- `validation.md#validating-arrays`
- `validation.md#validation`
- `validation.md#validation-quickstart`
- `validation.md#working-with-error-messages`
