# Upstream Compliance Report

Generated: 2026-04-22T02:04:07Z

Source of truth: services/compliance

## Compliance Workflow Requirements

- Every mapped Bedrock package compliance pass must include an upstream feature audit: `make sure we also have all the upstream features for <package>`.
- The audit compares upstream source contracts, public APIs, middleware, events, and runtime behavior against the Bedrock surface.
- Bedrock-equivalent features need executable parity tests; PHP-only or intentionally different behavior belongs in `services/compliance/divergences.yml`.
- `services/compliance/features.yml` must contain feature audit coverage for every mapped Bedrock surface.

## Compliance Command Center

Compliance target: classified parity. An upstream item is compliant when it is ported, adapted with a divergence rationale, or excluded at an approved product boundary.

| Area                  |   Classified | Missing | Compliance Status | Fastest Next Move                                                                         |
| --------------------- | -----------: | ------: | ----------------- | ----------------------------------------------------------------------------------------- |
| Tests                 |  738 / 13066 |   12328 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations.               |
| Documentation         |    367 / 875 |     508 | Not compliant     | Port relevant sections, adapt Go-specific sections, or exclude product-boundary sections. |
| Upstream skeleton demo |      49 / 49 |       0 | Complete          | Keep complete while upstream skeleton changes.                                            |
| Feature audits        |      59 / 62 |       3 | Not compliant     | Add audits for missing surfaces, then move partial surfaces to ported as parity closes.   |
| Overall               | 1213 / 14052 |   12839 | Not compliant     | Start with the Critical Path below; close the largest missing surfaces first.             |

## Critical Path

| Rank | Surface                                  | Missing | Classified | Status        | Next Action                                                                               |
| ---: | ---------------------------------------- | ------: | ---------: | ------------- | ----------------------------------------------------------------------------------------- |
|    1 | `Framework\Database` -> `database`      |    2626 |   0 / 2626 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    2 | `Framework\Support` -> `support`        |    1425 | 129 / 1554 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    3 | `upstream/ai` -> `ai/sdk`                 |     774 |    0 / 774 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    4 | `upstream/boost` -> `ai/boost`            |     671 |    0 / 671 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    5 | `upstream/mcp` -> `ai/mcp`                |     624 |    0 / 624 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    6 | `Framework\Collections` -> `collection` |     608 |   21 / 629 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    7 | `Framework\Validation` -> `validation`  |     563 |   34 / 597 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    8 | Upstream documentation                    |     508 |  367 / 875 | Not compliant | Port relevant sections, adapt Go-specific sections, or exclude product-boundary sections. |
|    9 | `Framework\Http` -> `httpx`             |     488 |    0 / 488 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   10 | `Framework\Testing`                     |     381 |    0 / 381 | Not compliant | Implement the surface or convert it to a permanent exclusion.                             |
|   11 | `Framework\Routing` -> `routing`        |     315 |    0 / 315 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   12 | `Framework\Cache` -> `cache`            |     313 |    0 / 313 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   13 | `upstream/prompts` -> `prompts`           |     310 |    3 / 313 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   14 | `upstream/search` -> `search`               |     252 |    0 / 252 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   15 | `upstream/websockets` -> `websockets`             |     232 |    0 / 232 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |

## Inventory Compliance Ledger

These inventory paths are tracking files, not compliant code paths. A row is compliant only when `Missing Tests` is 0.

| Inventory Path                            | What This Path Is                                      | Source                                                 | Classified | Missing Tests | Compliance Status | Next Action                                                                 |
| ----------------------------------------- | ------------------------------------------------------ | ------------------------------------------------------ | ---------: | ------------: | ----------------- | --------------------------------------------------------------------------- |
| inventories/framework-auth.txt            | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Auth`                    |  257 / 257 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-broadcasting.txt    | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Broadcasting`            |    55 / 55 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-bus.txt             | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Bus`                     |     0 / 66 |            66 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-cache.txt           | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Cache`                   |    0 / 313 |           313 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-concurrency.txt     | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Integration/Concurrency` |     0 / 12 |            12 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-conditionable.txt   | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Conditionable`           |      0 / 2 |             2 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-config.txt          | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Config`                  |     0 / 33 |            33 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-console.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Console`                 |    0 / 156 |           156 | Not compliant     | Implement the surface or convert it to a permanent exclusion.               |
| inventories/framework-container.txt       | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Container`               |    0 / 193 |           193 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-cookie.txt          | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Cookie`                  |     0 / 26 |            26 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-database.txt        | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Database`                |   0 / 2626 |          2626 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-encryption.txt      | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Encryption`              |     0 / 27 |            27 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-events.txt          | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Events`                  |     0 / 76 |            76 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-filesystem.txt      | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Filesystem`              |    0 / 134 |           134 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-hashing.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Hashing`                 |     0 / 13 |            13 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-http.txt            | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Http`                    |    0 / 488 |           488 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-json-schema.txt     | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/JsonSchema`              |     0 / 47 |            47 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-log.txt             | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Log`                     |     0 / 98 |            98 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-mail.txt            | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Mail`                    |    0 / 135 |           135 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-notifications.txt   | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Notifications`           |     0 / 71 |            71 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-pagination.txt      | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Pagination`              |     0 / 48 |            48 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-pipeline.txt        | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Pipeline`                |     0 / 22 |            22 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-process.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Process`                 |     0 / 54 |            54 | Not compliant     | Implement the surface or convert it to a permanent exclusion.               |
| inventories/framework-queue.txt           | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Queue`                   |  161 / 204 |            43 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-redis.txt           | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Redis`                   |    0 / 103 |           103 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-routing.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Routing`                 |    0 / 315 |           315 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-session.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Session`                 |     0 / 93 |            93 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-support.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Support`                 | 129 / 1554 |          1425 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-support-str.txt     | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Support`                 |   69 / 115 |            46 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-support-lottery.txt | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Support`                 |     0 / 14 |            14 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-testing.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Testing`                 |    0 / 381 |           381 | Not compliant     | Implement the surface or convert it to a permanent exclusion.               |
| inventories/framework-translation.txt     | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Translation`             |     9 / 52 |            43 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-validation.txt      | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Validation`              |   34 / 597 |           563 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-collections.txt     | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Support`                 |   21 / 629 |           608 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-ai.txt                | Tracking file for upstream tests, not a compliant path | `upstream/ai@0.x:tests`                                 |    0 / 774 |           774 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-boost.txt             | Tracking file for upstream tests, not a compliant path | `upstream/boost@main:tests`                             |    0 / 671 |           671 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-authflows.txt           | Tracking file for upstream tests, not a compliant path | `upstream/authflows@1.x:tests`                            |     0 / 91 |            91 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-authkit.txt         | Tracking file for upstream tests, not a compliant path | `upstream/authkit@5.x:tests`                          |     0 / 43 |            43 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-mcp.txt               | Tracking file for upstream tests, not a compliant path | `upstream/mcp@main:tests`                               |    0 / 624 |           624 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-oauthserver.txt          | Tracking file for upstream tests, not a compliant path | `upstream/oauthserver@13.x:tests`                          |    0 / 211 |           211 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-featureflags.txt           | Tracking file for upstream tests, not a compliant path | `upstream/featureflags@1.x:tests`                            |    0 / 186 |           186 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-prompts.txt           | Tracking file for upstream tests, not a compliant path | `upstream/prompts@main:tests`                           |    3 / 313 |           310 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-websockets.txt            | Tracking file for upstream tests, not a compliant path | `upstream/websockets@main:tests`                            |    0 / 232 |           232 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-search.txt             | Tracking file for upstream tests, not a compliant path | `upstream/search@11.x:tests`                             |    0 / 252 |           252 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-socialauth.txt         | Tracking file for upstream tests, not a compliant path | `upstream/socialauth@5.x:tests`                          |     0 / 40 |            40 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-billing.txt             | Tracking file for upstream tests, not a compliant path | `/Users/gocanto/Sites/madora/tests/Feature/Billing`    |    0 / 166 |           166 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-debugbar.txt         | Tracking file for upstream tests, not a compliant path | `upstream/debugbar@5.x:tests`                          |    0 / 113 |           113 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-jobqueue.txt           | Tracking file for upstream tests, not a compliant path | `upstream/jobqueue@5.x:tests`                            |    0 / 180 |           180 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-logtail.txt              | Tracking file for upstream tests, not a compliant path | `upstream/logtail@main:tests`                              |     0 / 38 |            38 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-remotetasks.txt             | Tracking file for upstream tests, not a compliant path | `upstream/remotetasks@2.x:tests`                              |     0 / 20 |            20 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-routegen.txt         | Tracking file for upstream tests, not a compliant path | `upstream/routegen@main:tests`                         |    0 / 103 |           103 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |

## Framework Coverage

| Upstream Source               | Bedrock Surface | Port Status          | Tests                                                                                 | Docs                        |
| ---------------------------- | --------------- | -------------------- | ------------------------------------------------------------------------------------- | --------------------------- |
| `Framework\Auth`            | `auth`          | `mapped`             | Ported tests: 54 / 257 (21.0%); Missing tests: 0 (0.0%); Adapted tests: 203 (79.0%)   | `packages/auth.md`          |
| `Framework\Broadcasting`    | `broadcasting`  | `mapped`             | Ported tests: 55 / 55 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)      | `packages/broadcasting.md`  |
| `Framework\Bus`             | `bus`           | `mapped`             | Ported tests: 0 / 66 (0.0%); Missing tests: 66 (100.0%); Adapted tests: 0 (0.0%)      | `packages/bus.md`           |
| `Framework\Cache`           | `cache`         | `mapped`             | Ported tests: 0 / 313 (0.0%); Missing tests: 313 (100.0%); Adapted tests: 0 (0.0%)    | `packages/cache.md`         |
| `Framework\Concurrency`     | `concurrency`   | `mapped`             | Ported tests: 0 / 12 (0.0%); Missing tests: 12 (100.0%); Adapted tests: 0 (0.0%)      | `packages/concurrency.md`   |
| `Framework\Conditionable`   | `conditionable` | `mapped`             | Ported tests: 0 / 2 (0.0%); Missing tests: 2 (100.0%); Adapted tests: 0 (0.0%)        | `packages/conditionable.md` |
| `Framework\Config`          | `config`        | `mapped`             | Ported tests: 0 / 33 (0.0%); Missing tests: 33 (100.0%); Adapted tests: 0 (0.0%)      | `packages/config.md`        |
| `Framework\Console`         | n/a             | `missing`            | Ported tests: 0 / 156 (0.0%); Missing tests: 156 (100.0%); Adapted tests: 0 (0.0%)    | n/a                         |
| `Framework\Container`       | `container`     | `mapped`             | Ported tests: 0 / 193 (0.0%); Missing tests: 193 (100.0%); Adapted tests: 0 (0.0%)    | `packages/container.md`     |
| `Framework\Contracts`       | `contracts`     | `mapped`             | no inventory                                                                          | `packages/contracts.md`     |
| `Framework\Cookie`          | `cookie`        | `mapped`             | Ported tests: 0 / 26 (0.0%); Missing tests: 26 (100.0%); Adapted tests: 0 (0.0%)      | `packages/cookie.md`        |
| `Framework\Database`        | `database`      | `mapped`             | Ported tests: 0 / 2626 (0.0%); Missing tests: 2626 (100.0%); Adapted tests: 0 (0.0%)  | `packages/database.md`      |
| `Framework\Encryption`      | `encryption`    | `mapped`             | Ported tests: 0 / 27 (0.0%); Missing tests: 27 (100.0%); Adapted tests: 0 (0.0%)      | `packages/encryption.md`    |
| `Framework\Events`          | `events`        | `mapped`             | Ported tests: 0 / 76 (0.0%); Missing tests: 76 (100.0%); Adapted tests: 0 (0.0%)      | `packages/events.md`        |
| `Framework\Filesystem`      | `filesystem`    | `mapped`             | Ported tests: 0 / 134 (0.0%); Missing tests: 134 (100.0%); Adapted tests: 0 (0.0%)    | `packages/filesystem.md`    |
| `Framework\Foundation`      | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                         |
| `Framework\Hashing`         | `hashing`       | `mapped`             | Ported tests: 0 / 13 (0.0%); Missing tests: 13 (100.0%); Adapted tests: 0 (0.0%)      | `packages/hashing.md`       |
| `Framework\Http`            | `httpx`         | `mapped`             | Ported tests: 0 / 488 (0.0%); Missing tests: 488 (100.0%); Adapted tests: 0 (0.0%)    | `packages/httpx.md`         |
| `Framework\JsonSchema`      | `jsonx`         | `mapped`             | Ported tests: 0 / 47 (0.0%); Missing tests: 47 (100.0%); Adapted tests: 0 (0.0%)      | `packages/jsonx.md`         |
| `Framework\Log`             | `log`           | `mapped`             | Ported tests: 0 / 98 (0.0%); Missing tests: 98 (100.0%); Adapted tests: 0 (0.0%)      | `packages/log.md`           |
| `Framework\Mail`            | `mailx`         | `mapped`             | Ported tests: 0 / 135 (0.0%); Missing tests: 135 (100.0%); Adapted tests: 0 (0.0%)    | `packages/mailx.md`         |
| `Framework\Notifications`   | `notifications` | `mapped`             | Ported tests: 0 / 71 (0.0%); Missing tests: 71 (100.0%); Adapted tests: 0 (0.0%)      | `packages/notifications.md` |
| `Framework\Pagination`      | `pagination`    | `mapped`             | Ported tests: 0 / 48 (0.0%); Missing tests: 48 (100.0%); Adapted tests: 0 (0.0%)      | `packages/pagination.md`    |
| `Framework\Pipeline`        | `pipeline`      | `mapped`             | Ported tests: 0 / 22 (0.0%); Missing tests: 22 (100.0%); Adapted tests: 0 (0.0%)      | `packages/pipeline.md`      |
| `Framework\Process`         | n/a             | `missing`            | Ported tests: 0 / 54 (0.0%); Missing tests: 54 (100.0%); Adapted tests: 0 (0.0%)      | n/a                         |
| `Framework\Queue`           | `queue`         | `mapped`             | Ported tests: 161 / 204 (78.9%); Missing tests: 43 (21.1%); Adapted tests: 0 (0.0%)   | `packages/queue.md`         |
| `Framework\Redis`           | `redis`         | `mapped`             | Ported tests: 0 / 103 (0.0%); Missing tests: 103 (100.0%); Adapted tests: 0 (0.0%)    | `packages/redis.md`         |
| `Framework\Routing`         | `routing`       | `mapped`             | Ported tests: 0 / 315 (0.0%); Missing tests: 315 (100.0%); Adapted tests: 0 (0.0%)    | `packages/routing.md`       |
| `Framework\Session`         | `session`       | `mapped`             | Ported tests: 0 / 93 (0.0%); Missing tests: 93 (100.0%); Adapted tests: 0 (0.0%)      | `packages/session.md`       |
| `Framework\Support`         | `support`       | `mapped`             | Ported tests: 129 / 1554 (8.3%); Missing tests: 1425 (91.7%); Adapted tests: 0 (0.0%) | `packages/support.md`       |
| `Framework\Support\Str`     | `str`           | `mapped`             | Ported tests: 69 / 115 (60.0%); Missing tests: 46 (40.0%); Adapted tests: 0 (0.0%)    | `packages/str.md`           |
| `Framework\Support\Lottery` | `lottery`       | `mapped`             | Ported tests: 0 / 14 (0.0%); Missing tests: 14 (100.0%); Adapted tests: 0 (0.0%)      | `packages/lottery.md`       |
| `Framework\Testing`         | n/a             | `missing`            | Ported tests: 0 / 381 (0.0%); Missing tests: 381 (100.0%); Adapted tests: 0 (0.0%)    | n/a                         |
| `Framework\Translation`     | `translation`   | `mapped`             | Ported tests: 9 / 52 (17.3%); Missing tests: 43 (82.7%); Adapted tests: 0 (0.0%)      | `packages/translation.md`   |
| `Framework\Validation`      | `validation`    | `mapped`             | Ported tests: 34 / 597 (5.7%); Missing tests: 563 (94.3%); Adapted tests: 0 (0.0%)    | `packages/validation.md`    |
| `Framework\Collections`     | `collection`    | `mapped`             | Ported tests: 21 / 629 (3.3%); Missing tests: 608 (96.7%); Adapted tests: 0 (0.0%)    | `packages/collection.md`    |
| `Framework\Macroable`       | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                         |
| `Framework\Reflection`      | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                         |
| `Framework\View`            | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                         |

## Package Coverage

| Upstream Package               | Bedrock Surface | Port Status          | Tests                                                                              | Docs                       |
| ----------------------------- | --------------- | -------------------- | ---------------------------------------------------------------------------------- | -------------------------- |
| `upstream/ai`                  | `ai/sdk`        | `mapped`             | Ported tests: 0 / 774 (0.0%); Missing tests: 774 (100.0%); Adapted tests: 0 (0.0%) | `packages/ai/sdk.md`       |
| `upstream/boost`               | `ai/boost`      | `mapped`             | Ported tests: 0 / 671 (0.0%); Missing tests: 671 (100.0%); Adapted tests: 0 (0.0%) | `packages/ai/boost.md`     |
| `upstream/broadcastclient`                | `broadcastclient`          | `mapped`             | no inventory                                                                       | `packages/broadcastclient.md`         |
| `upstream/authflows`             | `inception`     | `mapped`             | Ported tests: 0 / 91 (0.0%); Missing tests: 91 (100.0%); Adapted tests: 0 (0.0%)   | `packages/inception.md`    |
| `upstream/authkit`           | `inception`     | `mapped`             | Ported tests: 0 / 43 (0.0%); Missing tests: 43 (100.0%); Adapted tests: 0 (0.0%)   | `packages/inception.md`    |
| `upstream/mcp`                 | `ai/mcp`        | `mapped`             | Ported tests: 0 / 624 (0.0%); Missing tests: 624 (100.0%); Adapted tests: 0 (0.0%) | `packages/ai/mcp.md`       |
| `upstream/oauthserver`            | `oauthserver`      | `mapped`             | Ported tests: 0 / 211 (0.0%); Missing tests: 211 (100.0%); Adapted tests: 0 (0.0%) | `packages/oauthserver.md`     |
| `upstream/featureflags`             | `featureflags`       | `mapped`             | Ported tests: 0 / 186 (0.0%); Missing tests: 186 (100.0%); Adapted tests: 0 (0.0%) | `packages/featureflags.md`      |
| `upstream/httppreview`        | `httppreview`  | `mapped`             | no inventory                                                                       | `packages/httppreview.md` |
| `upstream/prompts`             | `prompts`       | `mapped`             | Ported tests: 3 / 313 (1.0%); Missing tests: 310 (99.0%); Adapted tests: 0 (0.0%)  | `packages/prompts.md`      |
| `upstream/websockets`              | `websockets`        | `mapped`             | Ported tests: 0 / 232 (0.0%); Missing tests: 232 (100.0%); Adapted tests: 0 (0.0%) | `packages/websockets.md`       |
| `upstream/search`               | `search`         | `mapped`             | Ported tests: 0 / 252 (0.0%); Missing tests: 252 (100.0%); Adapted tests: 0 (0.0%) | `packages/search.md`        |
| `upstream/socialauth`           | `socialauth`     | `mapped`             | Ported tests: 0 / 40 (0.0%); Missing tests: 40 (100.0%); Adapted tests: 0 (0.0%)   | `packages/socialauth.md`    |
| `upstream/billing`               | `billing`         | `mapped`             | Ported tests: 0 / 166 (0.0%); Missing tests: 166 (100.0%); Adapted tests: 0 (0.0%) | `packages/billing.md`        |
| `upstream/debugbar`           | `debugbar`     | `mapped`             | Ported tests: 0 / 113 (0.0%); Missing tests: 113 (100.0%); Adapted tests: 0 (0.0%) | `packages/debugbar.md`    |
| `upstream/jobqueue`             | `jobqueue`       | `mapped`             | Ported tests: 0 / 180 (0.0%); Missing tests: 180 (100.0%); Adapted tests: 0 (0.0%) | `packages/jobqueue.md`      |
| `upstream/logtail`                | `logtail`          | `mapped`             | Ported tests: 0 / 38 (0.0%); Missing tests: 38 (100.0%); Adapted tests: 0 (0.0%)   | `packages/logtail.md`         |
| `upstream/remotetasks`               | `remotetasks`         | `mapped`             | Ported tests: 0 / 20 (0.0%); Missing tests: 20 (100.0%); Adapted tests: 0 (0.0%)   | `packages/remotetasks.md`        |
| `upstream/routegen`           | `routegen`     | `mapped`             | Ported tests: 0 / 103 (0.0%); Missing tests: 103 (100.0%); Adapted tests: 0 (0.0%) | `packages/routegen.md`    |
| `upstream/tokens`             | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `upstream/octane`              | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `upstream/pulse`               | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `upstream/folio`               | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `upstream/dusk`                | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `upstream/browser-kit-testing` | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `upstream/sail`                | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `upstream/installer`           | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `upstream/valet`               | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `upstream/vite-plugin`         | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `starter-kits`                | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `upstream/workos`              | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |

## Test Porting Summary

| Scope                 | Inventories | Upstream Tests | Ported Tests | Pending / Missing Tests | Adapted Tests |
| --------------------- | ----------: | -------------: | -----------: | ----------------------: | ------------: |
| All inventories       |          51 |          13066 |   535 (4.1%) |           12328 (94.4%) |    203 (1.6%) |
| Framework inventories |          34 |           9009 |   532 (5.9%) |            8274 (91.8%) |    203 (2.3%) |
| Package inventories   |          17 |           4057 |     3 (0.1%) |            4054 (99.9%) |      0 (0.0%) |

## Documentation Porting Summary

| Scope        | Doc Inventories | Upstream Sections | Ported Sections | Pending / Missing Sections | Adapted Sections | Excluded Sections |
| ------------ | --------------: | ----------------: | --------------: | -------------------------: | ---------------: | ----------------: |
| Upstream docs |               1 |               875 |     197 (22.5%) |                508 (58.1%) |        10 (1.1%) |       160 (18.3%) |

## Upstream Skeleton Demo Summary

| Scope           | Skeleton Inventories | Upstream Files | Ported Files | Pending / Missing Files | Adapted Files | Excluded Files |
| --------------- | -------------------: | -------------: | -----------: | ----------------------: | ------------: | -------------: |
| upstream/upstream |                    1 |             49 |   20 (40.8%) |                0 (0.0%) |    10 (20.4%) |     19 (38.8%) |

## Documentation Detail

### Top Pending Documentation Files

| Upstream Docs File           | Pending Sections |
| --------------------------- | ---------------: |
| `orm.md`               |               20 |
| `queries.md`                |               16 |
| `oauthserver.md`               |               16 |
| `cashier-paddle.md`         |               15 |
| `prompts.md`                |               14 |
| `billing.md`                |               14 |
| `featureflags.md`                |               11 |
| `notifications.md`          |               10 |
| `migrations.md`             |               10 |
| `mcp.md`                    |               10 |
| `orm-relationships.md` |               10 |
| `contributions.md`          |               10 |
| `boost.md`                  |               10 |
| `cli.md`                |               10 |
| `mail.md`                   |                9 |
| `http-tests.md`             |                9 |
| `deployment.md`             |                9 |
| `ai-sdk.md`                 |                9 |
| `queues.md`                 |                8 |
| `packages.md`               |                8 |

### Classified Documentation Sections

| Upstream Docs File | Status   | Sections |
| ----------------- | -------- | -------: |
| `installation.md` | Adapted  |        4 |
| `lifecycle.md`    | Adapted  |        1 |
| `structure.md`    | Adapted  |        3 |
| `testing.md`      | Adapted  |        2 |
| `template.md`        | Excluded |       13 |
| `dusk.md`         | Excluded |       17 |
| `folio.md`        | Excluded |       17 |
| `homestead.md`    | Excluded |       11 |
| `mix.md`          | Excluded |        2 |
| `octane.md`       | Excluded |       12 |
| `pint.md`         | Excluded |        6 |
| `pulse.md`        | Excluded |        7 |
| `sail.md`         | Excluded |       25 |
| `tokens.md`      | Excluded |        8 |
| `starter-kits.md` | Excluded |        8 |
| `valet.md`        | Excluded |       12 |
| `views.md`        | Excluded |        6 |
| `vite.md`         | Excluded |       16 |

## Skeleton Detail

| Upstream Skeleton File                    | Status   |
| ---------------------------------------- | -------- |
| `.editorconfig`                          | Excluded |
| `.env.example`                           | Adapted  |
| `.gitattributes`                         | Excluded |
| `.github/workflows/issues.yml`           | Excluded |
| `.github/workflows/pull-requests.yml`    | Excluded |
| `.github/workflows/tests.yml`            | Adapted  |
| `.github/workflows/update-changelog.yml` | Excluded |
| `.gitignore`                             | Adapted  |
| `.npmrc`                                 | Adapted  |
| `.styleci.yml`                           | Excluded |
| `CHANGELOG.md`                           | Excluded |
| `README.md`                              | Adapted  |
| `app/Http/Controllers/Controller.php`    | Excluded |
| `cli`                                | Excluded |
| `bootstrap/providers.php`                | Adapted  |
| `composer.json`                          | Excluded |
| `config/auth.php`                        | Adapted  |
| `config/mail.php`                        | Adapted  |
| `config/services.php`                    | Adapted  |
| `package.json`                           | Excluded |
| `phpunit.xml`                            | Excluded |
| `public/.htaccess`                       | Excluded |
| `public/favicon.ico`                     | Excluded |
| `public/robots.txt`                      | Excluded |
| `resources/css/app.css`                  | Excluded |
| `resources/js/app.js`                    | Excluded |
| `resources/views/welcome.template.php`      | Excluded |
| `routes/console.php`                     | Adapted  |
| `vite.config.js`                         | Excluded |

## Feature Coverage

| Feature                      | Source                      | Bedrock Surface | Status     | Docs                        | Notes                                                                                                                                                                                                                                        |
| ---------------------------- | --------------------------- | --------------- | ---------- | --------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Authentication guards        | `framework.auth`            | `auth`          | `ported`   | `packages/auth.md`          | Session, token, request, basic-auth, provider, and event surfaces are complete under classified parity; PHP runtime mechanics are covered by Go adaptation rules in divergences.yml.                                                         |
| Password broker              | `framework.auth`            | `auth`          | `ported`   | `packages/auth.md`          | Password reset tokens, throttling, callback notification override, notification sender hook, and reset-link events are complete under classified parity; PHP manager/config mechanics are covered by Go adaptation rules in divergences.yml. |
| Broadcasting contracts       | `framework.broadcasting`    | `broadcasting`  | `ported`   | `packages/broadcasting.md`  | Server-side channel auth, broadcast event payload handling, Pusher/Redis/Ably broadcaster responses, and Pusher channel conventions are ported; BroadcastClient remains the client receiver package.                                                    |
| Command bus dispatch         | `framework.bus`             | `bus`           | `partial`  | `packages/bus.md`           | Bus package exists with command dispatch and middleware; upstream test parity is not complete.                                                                                                                                               |
| Cache repository             | `framework.cache`           | `cache`         | `partial`  | `packages/cache.md`         | Cache abstractions are present, but upstream test parity is still missing.                                                                                                                                                                   |
| Concurrent task execution    | `framework.concurrency`     | `concurrency`   | `partial`  | `packages/concurrency.md`   | Bedrock has a concurrency package; upstream concurrency tests have not been matched yet.                                                                                                                                                     |
| Conditional fluent execution | `framework.conditionable`   | `conditionable` | `partial`  | `packages/conditionable.md` | Go fluent proxy exists; PHP dynamic behaviour is adapted to Go.                                                                                                                                                                              |
| Configuration repository     | `framework.config`          | `config`        | `partial`  | `packages/config.md`        | Config package exists; upstream test parity is not complete.                                                                                                                                                                                 |
| Console commands             | `framework.console`         | n/a             | `missing`  | n/a                         | No Bedrock console package exists yet.                                                                                                                                                                                                       |
| Container bindings           | `framework.container`       | `container`     | `partial`  | `packages/container.md`     | Container package exists; upstream container tests have not been matched yet.                                                                                                                                                                |
| Shared interfaces            | `framework.contracts`       | `contracts`     | `ported`   | `packages/contracts.md`     | Interface compliance is tracked through concrete package inventories.                                                                                                                                                                        |
| Cookie handling              | `framework.cookie`          | `cookie`        | `partial`  | `packages/cookie.md`        | Cookie package exists; upstream test parity is not complete.                                                                                                                                                                                 |
| Database query builder       | `framework.database`        | `database`      | `partial`  | `packages/database.md`      | Query, schema, migrations, and Orm surfaces exist; broad upstream parity remains incomplete.                                                                                                                                            |
| Orm-style ORM           | `framework.database`        | `database`      | `partial`  | `packages/database.md`      | ORM primitives exist with Go adaptations for attributes, relations, and errors.                                                                                                                                                              |
| Encrypter                    | `framework.encryption`      | `encryption`    | `partial`  | `packages/encryption.md`    | CBC and GCM encryption exist; upstream encryption tests are not fully matched.                                                                                                                                                               |
| Event dispatcher             | `framework.events`          | `events`        | `partial`  | `packages/events.md`        | Event dispatching exists; upstream test parity remains incomplete.                                                                                                                                                                           |
| Local filesystem             | `framework.filesystem`      | `filesystem`    | `partial`  | `packages/filesystem.md`    | Local filesystem operations are implemented; upstream parity tests are not complete.                                                                                                                                                         |
| Application foundation       | `framework.foundation`      | n/a             | `excluded` | n/a                         | Full Upstream application kernel is outside Bedrock product boundaries.                                                                                                                                                                       |
| Password hashing             | `framework.hashing`         | `hashing`       | `partial`  | `packages/hashing.md`       | Hashing drivers exist; upstream test parity is not complete.                                                                                                                                                                                 |
| HTTP primitives              | `framework.http`            | `httpx`         | `partial`  | `packages/httpx.md`         | HTTP primitives and test helpers exist; upstream test parity remains incomplete.                                                                                                                                                             |
| JSON schema builder          | `framework.json_schema`     | `jsonx`         | `partial`  | `packages/jsonx.md`         | Fluent JSON schema builder exists; upstream test parity is not complete.                                                                                                                                                                     |
| Logging channels             | `framework.log`             | `log`           | `partial`  | `packages/log.md`           | Driver-based structured logging exists; upstream log tests are not matched yet.                                                                                                                                                              |
| Mail messages                | `framework.mail`            | `mailx`         | `partial`  | `packages/mailx.md`         | Mail package exists; upstream mail test parity is incomplete.                                                                                                                                                                                |
| Notification channels        | `framework.notifications`   | `notifications` | `partial`  | `packages/notifications.md` | Multi-channel notifications exist; upstream parity remains incomplete.                                                                                                                                                                       |
| Pagination                   | `framework.pagination`      | `pagination`    | `partial`  | `packages/pagination.md`    | Offset and cursor paginators exist; upstream parity is incomplete.                                                                                                                                                                           |
| Pipeline                     | `framework.pipeline`        | `pipeline`      | `partial`  | `packages/pipeline.md`      | Pipeline package exists; upstream test parity is not complete.                                                                                                                                                                               |
| Process runner               | `framework.process`         | n/a             | `missing`  | n/a                         | No Bedrock process package exists yet.                                                                                                                                                                                                       |
| Queue workers and drivers    | `framework.queue`           | `queue`         | `partial`  | `packages/queue.md`         | Queue has the strongest upstream test match so far, but missing tests remain.                                                                                                                                                                |
| Redis commands               | `framework.redis`           | `redis`         | `partial`  | `packages/redis.md`         | Redis command surface exists; upstream Redis tests are not matched yet.                                                                                                                                                                      |
| Routing                      | `framework.routing`         | `routing`       | `partial`  | `packages/routing.md`       | Routing package exists; upstream routing tests remain incomplete.                                                                                                                                                                            |
| Session stores               | `framework.session`         | `session`       | `partial`  | `packages/session.md`       | Session stores exist; upstream session tests are not matched yet.                                                                                                                                                                            |
| Support helpers              | `framework.support`         | `support`       | `partial`  | `packages/support.md`       | Several support helpers are ported; upstream support coverage is still partial.                                                                                                                                                              |
| String helpers               | `framework.support_str`     | `str`           | `partial`  | `packages/str.md`           | String helper parity has many ported tests, but missing cases remain.                                                                                                                                                                        |
| Lottery helper               | `framework.support_lottery` | `lottery`       | `partial`  | `packages/lottery.md`       | Lottery package exists; upstream Lottery tests are not matched yet.                                                                                                                                                                          |
| Testing utilities            | `framework.testing`         | n/a             | `missing`  | n/a                         | No standalone Bedrock testing package exists yet.                                                                                                                                                                                            |
| Translation loader           | `framework.translation`     | `translation`   | `partial`  | `packages/translation.md`   | Translation has some matched upstream tests; full parity remains incomplete.                                                                                                                                                                 |
| Validation rules             | `framework.validation`      | `validation`    | `partial`  | `packages/validation.md`    | Rule engine and many rules exist; upstream test parity is still broad and incomplete.                                                                                                                                                        |
| Collections                  | `framework.collections`     | `collection`    | `partial`  | `packages/collection.md`    | Fluent slice, key-value, lazy, arr, and kv packages are imported; legacy support helper cleanup is deferred.                                                                                                                                 |
| Macroable runtime            | `framework.macroable`       | n/a             | `excluded` | n/a                         | PHP runtime macro behaviour does not map cleanly to Go static method sets.                                                                                                                                                                   |
| Reflection helpers           | `framework.reflection`      | n/a             | `excluded` | n/a                         | PHP reflection helpers are not a Bedrock product surface.                                                                                                                                                                                    |
| Template views                  | `framework.view`            | n/a             | `excluded` | n/a                         | Bedrock does not ship Upstream Template or the view engine.                                                                                                                                                                                      |
| AI SDK                       | `package.ai`                | `ai/sdk`        | `partial`  | `packages/ai/sdk.md`        | Bedrock has an AI SDK package; upstream package parity is not complete.                                                                                                                                                                      |
| Agent tooling                | `package.boost`             | `ai/boost`      | `partial`  | `packages/ai/boost.md`      | Boost package exists; upstream package parity remains incomplete.                                                                                                                                                                            |
| BroadcastClient client                  | `package.broadcastclient`              | `broadcastclient`          | `partial`  | `packages/broadcastclient.md`          | Tracked by feature inventory because there is no stable upstream test inventory configured.                                                                                                                                                  |
| AuthFlows auth flows           | `package.authflows`           | `inception`     | `partial`  | `packages/authflows.md`       | Exposed through Inception with AuthFlows-specific documentation.                                                                                                                                                                               |
| AuthKit teams              | `package.authkit`         | `inception`     | `partial`  | `packages/authkit.md`     | Exposed through Inception with AuthKit-specific documentation.                                                                                                                                                                             |
| MCP server                   | `package.mcp`               | `ai/mcp`        | `partial`  | `packages/ai/mcp.md`        | MCP package exists; upstream package parity remains incomplete.                                                                                                                                                                              |
| OAuth server                 | `package.oauthserver`          | `oauthserver`      | `partial`  | `packages/oauthserver.md`      | OAuthServer package and docs exist; upstream package parity is not complete.                                                                                                                                                                    |
| Feature flags                | `package.featureflags`           | `featureflags`       | `partial`  | `packages/featureflags.md`       | FeatureFlags package exists; upstream package parity remains incomplete.                                                                                                                                                                          |
| HTTPPreview client          | `package.httppreview`      | `httppreview`  | `partial`  | `packages/httppreview.md`  | Tracked by feature inventory because there is no stable upstream test inventory configured.                                                                                                                                                  |
| CLI prompts                  | `package.prompts`           | `prompts`       | `partial`  | `packages/prompts.md`       | Prompt package exists with a small number of matched upstream tests.                                                                                                                                                                         |
| WebSockets websocket server      | `package.websockets`            | `websockets`        | `partial`  | `packages/websockets.md`        | WebSockets package exists; upstream package parity remains incomplete.                                                                                                                                                                           |
| Search indexing              | `package.search`             | `search`         | `partial`  | `packages/search.md`         | Search package exists; upstream package parity remains incomplete.                                                                                                                                                                            |
| Social OAuth                 | `package.socialauth`         | `socialauth`     | `partial`  | `packages/socialauth.md`     | SocialAuth package exists; upstream package parity remains incomplete.                                                                                                                                                                        |
| Subscription billing         | `package.billing`             | `billing`         | `partial`  | `packages/billing.md`         | Billing is tracked against local Madora billing tests.                                                                                                                                                                                         |
| DebugBar observability      | `package.debugbar`         | `debugbar`     | `partial`  | `packages/debugbar.md`     | DebugBar package exists; upstream package parity remains incomplete.                                                                                                                                                                        |
| RouteGen route helpers      | `package.routegen`         | `routegen`     | `partial`  | `packages/routegen.md`     | RouteGen package exists; upstream package parity remains incomplete.                                                                                                                                                                        |
| Queue monitoring             | `package.jobqueue`           | `jobqueue`       | `partial`  | `packages/jobqueue.md`       | JobQueue is tracked as queue monitoring primitives; the browser dashboard is not ported.                                                                                                                                                      |
| Log tailing                  | `package.logtail`              | `logtail`          | `partial`  | `packages/logtail.md`          | LogTail is tracked as log parsing and filtering primitives for Go readers.                                                                                                                                                                      |
| Remote task runner           | `package.remotetasks`             | `remotetasks`         | `partial`  | `packages/remotetasks.md`         | RemoteTasks is tracked as task planning and command execution primitives with injectable runners.                                                                                                                                                  |
| Tokens tokens               | `package.tokens`           | n/a             | `excluded` | n/a                         | Tokens is intentionally not part of Bedrock.                                                                                                                                                                                                |
| Browser automation           | `package.dusk`              | n/a             | `excluded` | n/a                         | Browser automation package is excluded; browser work should use Agent Browser.                                                                                                                                                               |

## Sources Without Test Inventories

| Source                 | Bedrock        | Tracking                     | Reason                                                                                                                    |
| ---------------------- | -------------- | ---------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `Framework\Contracts` | `contracts`    | Concrete package inventories | Interface-only component; compliance is verified through concrete package inventories.                                    |
| `upstream/broadcastclient`         | `broadcastclient`         | Feature inventory            | TypeScript client package; compliance is tracked by feature coverage until a stable upstream test inventory is available. |
| `upstream/httppreview` | `httppreview` | Feature inventory            | JavaScript package; compliance is tracked by feature coverage until a stable upstream test inventory is available.        |

## Permanent Exclusions

- `Framework\Foundation` - Product boundary; Bedrock does not ship Upstream's full application kernel/Foundation layer.
- `Framework\Macroable` - PHP runtime macro behaviour does not map cleanly to Go's static method set.
- `Framework\Reflection` - PHP reflection helpers are not a Bedrock product surface.
- `Framework\View` - Product boundary; Bedrock does not ship Upstream's Template/view engine.
- `upstream/tokens` - Product boundary; Tokens is intentionally not part of Bedrock.
- `upstream/octane` - Product boundary; PHP worker/runtime management does not map to Bedrock.
- `upstream/pulse` - Product boundary; application monitoring dashboard is intentionally excluded.
- `upstream/folio` - Product boundary; page-based PHP routing is intentionally excluded.
- `upstream/dusk` - Product boundary; browser automation package is excluded. Any future browser work must use Agent Browser.
- `upstream/browser-kit-testing` - Product boundary; legacy BrowserKit compatibility is intentionally excluded.
- `upstream/sail` - Product boundary; Docker development environment tooling is intentionally excluded.
- `upstream/installer` - Product boundary; Upstream application installer tooling is intentionally excluded.
- `upstream/valet` - Product boundary; macOS local development service management is intentionally excluded.
- `upstream/vite-plugin` - Product boundary; Vite integration tooling is intentionally excluded.
- `starter-kits` - Product boundary; application starter kits are intentionally excluded.
- `upstream/workos` - Product boundary; WorkOS starter-kit integration is intentionally excluded.
