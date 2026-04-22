# Laravel Compliance Report

Generated: 2026-04-22T04:37:44Z

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
| Tests                 | 1556 / 13066 |   11510 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations.               |
| Documentation         |    370 / 875 |     505 | Not compliant     | Port relevant sections, adapt Go-specific sections, or exclude product-boundary sections. |
| Laravel skeleton demo |      49 / 49 |       0 | Complete          | Keep complete while upstream skeleton changes.                                            |
| Feature audits        |      60 / 62 |       2 | Not compliant     | Add audits for missing surfaces, then move partial surfaces to ported as parity closes.   |
| Overall               | 2035 / 14052 |   12017 | Not compliant     | Start with the Critical Path below; close the largest missing surfaces first.             |

## Critical Path

| Rank | Surface                                  | Missing | Classified | Status        | Next Action                                                                               |
| ---: | ---------------------------------------- | ------: | ---------: | ------------- | ----------------------------------------------------------------------------------------- |
|    1 | `Illuminate\Database` -> `database`      |    2626 |   0 / 2626 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    2 | `Illuminate\Support` -> `support`        |    1425 | 129 / 1554 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    3 | `laravel/ai` -> `ai/sdk`                 |     774 |    0 / 774 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    4 | `laravel/boost` -> `ai/boost`            |     671 |    0 / 671 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    5 | `laravel/mcp` -> `ai/mcp`                |     624 |    0 / 624 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    6 | `Illuminate\Collections` -> `collection` |     608 |   21 / 629 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    7 | `Illuminate\Validation` -> `validation`  |     563 |   34 / 597 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    8 | Laravel documentation                    |     505 |  370 / 875 | Not compliant | Port relevant sections, adapt Go-specific sections, or exclude product-boundary sections. |
|    9 | `Illuminate\Http` -> `httpx`             |     488 |    0 / 488 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   10 | `Illuminate\Testing`                     |     381 |    0 / 381 | Not compliant | Implement the surface or convert it to a permanent exclusion.                             |
|   11 | `Illuminate\Routing` -> `routing`        |     315 |    0 / 315 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   12 | `laravel/prompts` -> `prompts`           |     310 |    3 / 313 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   13 | `laravel/scout` -> `scout`               |     252 |    0 / 252 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   14 | `laravel/reverb` -> `reverb`             |     232 |    0 / 232 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   15 | `laravel/passport` -> `passport`         |     211 |    0 / 211 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |

## Inventory Compliance Ledger

These inventory paths are tracking files, not compliant code paths. A row is compliant only when `Missing Tests` is 0.

| Inventory Path                            | What This Path Is                                      | Source                                                 | Classified | Missing Tests | Compliance Status | Next Action                                                                 |
| ----------------------------------------- | ------------------------------------------------------ | ------------------------------------------------------ | ---------: | ------------: | ----------------- | --------------------------------------------------------------------------- |
| inventories/framework-auth.txt            | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Auth`                    |  257 / 257 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-broadcasting.txt    | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Broadcasting`            |    55 / 55 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-bus.txt             | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Bus`                     |    66 / 66 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-cache.txt           | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Cache`                   |  313 / 313 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-concurrency.txt     | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Integration/Concurrency` |    12 / 12 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-conditionable.txt   | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Conditionable`           |      2 / 2 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-config.txt          | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Config`                  |    33 / 33 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-console.txt         | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Console`                 |  156 / 156 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-container.txt       | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Container`               |  193 / 193 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-cookie.txt          | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Cookie`                  |     0 / 26 |            26 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-database.txt        | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Database`                |   0 / 2626 |          2626 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-encryption.txt      | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Encryption`              |     0 / 27 |            27 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-events.txt          | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Events`                  |     0 / 76 |            76 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-filesystem.txt      | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Filesystem`              |    0 / 134 |           134 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-hashing.txt         | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Hashing`                 |     0 / 13 |            13 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-http.txt            | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Http`                    |    0 / 488 |           488 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-json-schema.txt     | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/JsonSchema`              |     0 / 47 |            47 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-log.txt             | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Log`                     |     0 / 98 |            98 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-mail.txt            | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Mail`                    |    0 / 135 |           135 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-notifications.txt   | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Notifications`           |     0 / 71 |            71 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-pagination.txt      | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Pagination`              |     0 / 48 |            48 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-pipeline.txt        | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Pipeline`                |     0 / 22 |            22 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-process.txt         | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Process`                 |     0 / 54 |            54 | Not compliant     | Implement the surface or convert it to a permanent exclusion.               |
| inventories/framework-queue.txt           | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Queue`                   |  204 / 204 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-redis.txt           | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Redis`                   |    0 / 103 |           103 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-routing.txt         | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Routing`                 |    0 / 315 |           315 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-session.txt         | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Session`                 |     0 / 93 |            93 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-support.txt         | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Support`                 | 129 / 1554 |          1425 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-support-str.txt     | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Support`                 |   69 / 115 |            46 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-support-lottery.txt | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Support`                 |     0 / 14 |            14 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-testing.txt         | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Testing`                 |    0 / 381 |           381 | Not compliant     | Implement the surface or convert it to a permanent exclusion.               |
| inventories/framework-translation.txt     | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Translation`             |     9 / 52 |            43 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-validation.txt      | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Validation`              |   34 / 597 |           563 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-collections.txt     | Tracking file for upstream tests, not a compliant path | `laravel/framework@13.x:tests/Support`                 |   21 / 629 |           608 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-ai.txt                | Tracking file for upstream tests, not a compliant path | `laravel/ai@0.x:tests`                                 |    0 / 774 |           774 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-boost.txt             | Tracking file for upstream tests, not a compliant path | `laravel/boost@main:tests`                             |    0 / 671 |           671 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-fortify.txt           | Tracking file for upstream tests, not a compliant path | `laravel/fortify@1.x:tests`                            |     0 / 91 |            91 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-jetstream.txt         | Tracking file for upstream tests, not a compliant path | `laravel/jetstream@5.x:tests`                          |     0 / 43 |            43 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-mcp.txt               | Tracking file for upstream tests, not a compliant path | `laravel/mcp@main:tests`                               |    0 / 624 |           624 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-passport.txt          | Tracking file for upstream tests, not a compliant path | `laravel/passport@13.x:tests`                          |    0 / 211 |           211 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-pennant.txt           | Tracking file for upstream tests, not a compliant path | `laravel/pennant@1.x:tests`                            |    0 / 186 |           186 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-prompts.txt           | Tracking file for upstream tests, not a compliant path | `laravel/prompts@main:tests`                           |    3 / 313 |           310 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-reverb.txt            | Tracking file for upstream tests, not a compliant path | `laravel/reverb@main:tests`                            |    0 / 232 |           232 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-scout.txt             | Tracking file for upstream tests, not a compliant path | `laravel/scout@11.x:tests`                             |    0 / 252 |           252 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-socialite.txt         | Tracking file for upstream tests, not a compliant path | `laravel/socialite@5.x:tests`                          |     0 / 40 |            40 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-spark.txt             | Tracking file for upstream tests, not a compliant path | `/Users/gocanto/Sites/madora/tests/Feature/Billing`    |    0 / 166 |           166 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-telescope.txt         | Tracking file for upstream tests, not a compliant path | `laravel/telescope@5.x:tests`                          |    0 / 113 |           113 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-horizon.txt           | Tracking file for upstream tests, not a compliant path | `laravel/horizon@5.x:tests`                            |    0 / 180 |           180 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-pail.txt              | Tracking file for upstream tests, not a compliant path | `laravel/pail@main:tests`                              |     0 / 38 |            38 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-envoy.txt             | Tracking file for upstream tests, not a compliant path | `laravel/envoy@2.x:tests`                              |     0 / 20 |            20 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-wayfinder.txt         | Tracking file for upstream tests, not a compliant path | `laravel/wayfinder@main:tests`                         |    0 / 103 |           103 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |

## Framework Coverage

| Laravel Source               | Bedrock Surface | Port Status          | Tests                                                                                 | Docs                        |
| ---------------------------- | --------------- | -------------------- | ------------------------------------------------------------------------------------- | --------------------------- |
| `Illuminate\Auth`            | `auth`          | `mapped`             | Ported tests: 54 / 257 (21.0%); Missing tests: 0 (0.0%); Adapted tests: 203 (79.0%)   | `packages/auth.md`          |
| `Illuminate\Broadcasting`    | `broadcasting`  | `mapped`             | Ported tests: 55 / 55 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)      | `packages/broadcasting.md`  |
| `Illuminate\Bus`             | `bus`           | `mapped`             | Ported tests: 54 / 66 (81.8%); Missing tests: 0 (0.0%); Adapted tests: 12 (18.2%)     | `packages/bus.md`           |
| `Illuminate\Cache`           | `cache`         | `mapped`             | Ported tests: 217 / 313 (69.3%); Missing tests: 0 (0.0%); Adapted tests: 96 (30.7%)   | `packages/cache.md`         |
| `Illuminate\Concurrency`     | `concurrency`   | `mapped`             | Ported tests: 5 / 12 (41.7%); Missing tests: 0 (0.0%); Adapted tests: 7 (58.3%)       | `packages/concurrency.md`   |
| `Illuminate\Conditionable`   | `conditionable` | `mapped`             | Ported tests: 2 / 2 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)        | `packages/conditionable.md` |
| `Illuminate\Config`          | `config`        | `mapped`             | Ported tests: 32 / 33 (97.0%); Missing tests: 0 (0.0%); Adapted tests: 1 (3.0%)       | `packages/config.md`        |
| `Illuminate\Console`         | `console`       | `mapped`             | Ported tests: 156 / 156 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)    | `packages/console.md`       |
| `Illuminate\Container`       | `container`     | `mapped`             | Ported tests: 102 / 193 (52.8%); Missing tests: 0 (0.0%); Adapted tests: 91 (47.2%)   | `packages/container.md`     |
| `Illuminate\Contracts`       | `contracts`     | `mapped`             | no inventory                                                                          | `packages/contracts.md`     |
| `Illuminate\Cookie`          | `cookie`        | `mapped`             | Ported tests: 0 / 26 (0.0%); Missing tests: 26 (100.0%); Adapted tests: 0 (0.0%)      | `packages/cookie.md`        |
| `Illuminate\Database`        | `database`      | `mapped`             | Ported tests: 0 / 2626 (0.0%); Missing tests: 2626 (100.0%); Adapted tests: 0 (0.0%)  | `packages/database.md`      |
| `Illuminate\Encryption`      | `encryption`    | `mapped`             | Ported tests: 0 / 27 (0.0%); Missing tests: 27 (100.0%); Adapted tests: 0 (0.0%)      | `packages/encryption.md`    |
| `Illuminate\Events`          | `events`        | `mapped`             | Ported tests: 0 / 76 (0.0%); Missing tests: 76 (100.0%); Adapted tests: 0 (0.0%)      | `packages/events.md`        |
| `Illuminate\Filesystem`      | `filesystem`    | `mapped`             | Ported tests: 0 / 134 (0.0%); Missing tests: 134 (100.0%); Adapted tests: 0 (0.0%)    | `packages/filesystem.md`    |
| `Illuminate\Foundation`      | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                         |
| `Illuminate\Hashing`         | `hashing`       | `mapped`             | Ported tests: 0 / 13 (0.0%); Missing tests: 13 (100.0%); Adapted tests: 0 (0.0%)      | `packages/hashing.md`       |
| `Illuminate\Http`            | `httpx`         | `mapped`             | Ported tests: 0 / 488 (0.0%); Missing tests: 488 (100.0%); Adapted tests: 0 (0.0%)    | `packages/httpx.md`         |
| `Illuminate\JsonSchema`      | `jsonx`         | `mapped`             | Ported tests: 0 / 47 (0.0%); Missing tests: 47 (100.0%); Adapted tests: 0 (0.0%)      | `packages/jsonx.md`         |
| `Illuminate\Log`             | `log`           | `mapped`             | Ported tests: 0 / 98 (0.0%); Missing tests: 98 (100.0%); Adapted tests: 0 (0.0%)      | `packages/log.md`           |
| `Illuminate\Mail`            | `mailx`         | `mapped`             | Ported tests: 0 / 135 (0.0%); Missing tests: 135 (100.0%); Adapted tests: 0 (0.0%)    | `packages/mailx.md`         |
| `Illuminate\Notifications`   | `notifications` | `mapped`             | Ported tests: 0 / 71 (0.0%); Missing tests: 71 (100.0%); Adapted tests: 0 (0.0%)      | `packages/notifications.md` |
| `Illuminate\Pagination`      | `pagination`    | `mapped`             | Ported tests: 0 / 48 (0.0%); Missing tests: 48 (100.0%); Adapted tests: 0 (0.0%)      | `packages/pagination.md`    |
| `Illuminate\Pipeline`        | `pipeline`      | `mapped`             | Ported tests: 0 / 22 (0.0%); Missing tests: 22 (100.0%); Adapted tests: 0 (0.0%)      | `packages/pipeline.md`      |
| `Illuminate\Process`         | n/a             | `missing`            | Ported tests: 0 / 54 (0.0%); Missing tests: 54 (100.0%); Adapted tests: 0 (0.0%)      | n/a                         |
| `Illuminate\Queue`           | `queue`         | `mapped`             | Ported tests: 195 / 204 (95.6%); Missing tests: 0 (0.0%); Adapted tests: 9 (4.4%)     | `packages/queue.md`         |
| `Illuminate\Redis`           | `redis`         | `mapped`             | Ported tests: 0 / 103 (0.0%); Missing tests: 103 (100.0%); Adapted tests: 0 (0.0%)    | `packages/redis.md`         |
| `Illuminate\Routing`         | `routing`       | `mapped`             | Ported tests: 0 / 315 (0.0%); Missing tests: 315 (100.0%); Adapted tests: 0 (0.0%)    | `packages/routing.md`       |
| `Illuminate\Session`         | `session`       | `mapped`             | Ported tests: 0 / 93 (0.0%); Missing tests: 93 (100.0%); Adapted tests: 0 (0.0%)      | `packages/session.md`       |
| `Illuminate\Support`         | `support`       | `mapped`             | Ported tests: 129 / 1554 (8.3%); Missing tests: 1425 (91.7%); Adapted tests: 0 (0.0%) | `packages/support.md`       |
| `Illuminate\Support\Str`     | `str`           | `mapped`             | Ported tests: 69 / 115 (60.0%); Missing tests: 46 (40.0%); Adapted tests: 0 (0.0%)    | `packages/str.md`           |
| `Illuminate\Support\Lottery` | `lottery`       | `mapped`             | Ported tests: 0 / 14 (0.0%); Missing tests: 14 (100.0%); Adapted tests: 0 (0.0%)      | `packages/lottery.md`       |
| `Illuminate\Testing`         | n/a             | `missing`            | Ported tests: 0 / 381 (0.0%); Missing tests: 381 (100.0%); Adapted tests: 0 (0.0%)    | n/a                         |
| `Illuminate\Translation`     | `translation`   | `mapped`             | Ported tests: 9 / 52 (17.3%); Missing tests: 43 (82.7%); Adapted tests: 0 (0.0%)      | `packages/translation.md`   |
| `Illuminate\Validation`      | `validation`    | `mapped`             | Ported tests: 34 / 597 (5.7%); Missing tests: 563 (94.3%); Adapted tests: 0 (0.0%)    | `packages/validation.md`    |
| `Illuminate\Collections`     | `collection`    | `mapped`             | Ported tests: 21 / 629 (3.3%); Missing tests: 608 (96.7%); Adapted tests: 0 (0.0%)    | `packages/collection.md`    |
| `Illuminate\Macroable`       | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                         |
| `Illuminate\Reflection`      | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                         |
| `Illuminate\View`            | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                         |

## Package Coverage

| Laravel Package               | Bedrock Surface | Port Status          | Tests                                                                              | Docs                       |
| ----------------------------- | --------------- | -------------------- | ---------------------------------------------------------------------------------- | -------------------------- |
| `laravel/ai`                  | `ai/sdk`        | `mapped`             | Ported tests: 0 / 774 (0.0%); Missing tests: 774 (100.0%); Adapted tests: 0 (0.0%) | `packages/ai/sdk.md`       |
| `laravel/boost`               | `ai/boost`      | `mapped`             | Ported tests: 0 / 671 (0.0%); Missing tests: 671 (100.0%); Adapted tests: 0 (0.0%) | `packages/ai/boost.md`     |
| `laravel/echo`                | `echo`          | `mapped`             | no inventory                                                                       | `packages/echo.md`         |
| `laravel/fortify`             | `inception`     | `mapped`             | Ported tests: 0 / 91 (0.0%); Missing tests: 91 (100.0%); Adapted tests: 0 (0.0%)   | `packages/inception.md`    |
| `laravel/jetstream`           | `inception`     | `mapped`             | Ported tests: 0 / 43 (0.0%); Missing tests: 43 (100.0%); Adapted tests: 0 (0.0%)   | `packages/inception.md`    |
| `laravel/mcp`                 | `ai/mcp`        | `mapped`             | Ported tests: 0 / 624 (0.0%); Missing tests: 624 (100.0%); Adapted tests: 0 (0.0%) | `packages/ai/mcp.md`       |
| `laravel/passport`            | `passport`      | `mapped`             | Ported tests: 0 / 211 (0.0%); Missing tests: 211 (100.0%); Adapted tests: 0 (0.0%) | `packages/passport.md`     |
| `laravel/pennant`             | `pennant`       | `mapped`             | Ported tests: 0 / 186 (0.0%); Missing tests: 186 (100.0%); Adapted tests: 0 (0.0%) | `packages/pennant.md`      |
| `laravel/precognition`        | `precognition`  | `mapped`             | no inventory                                                                       | `packages/precognition.md` |
| `laravel/prompts`             | `prompts`       | `mapped`             | Ported tests: 3 / 313 (1.0%); Missing tests: 310 (99.0%); Adapted tests: 0 (0.0%)  | `packages/prompts.md`      |
| `laravel/reverb`              | `reverb`        | `mapped`             | Ported tests: 0 / 232 (0.0%); Missing tests: 232 (100.0%); Adapted tests: 0 (0.0%) | `packages/reverb.md`       |
| `laravel/scout`               | `scout`         | `mapped`             | Ported tests: 0 / 252 (0.0%); Missing tests: 252 (100.0%); Adapted tests: 0 (0.0%) | `packages/scout.md`        |
| `laravel/socialite`           | `socialite`     | `mapped`             | Ported tests: 0 / 40 (0.0%); Missing tests: 40 (100.0%); Adapted tests: 0 (0.0%)   | `packages/socialite.md`    |
| `laravel/spark`               | `spark`         | `mapped`             | Ported tests: 0 / 166 (0.0%); Missing tests: 166 (100.0%); Adapted tests: 0 (0.0%) | `packages/spark.md`        |
| `laravel/telescope`           | `telescope`     | `mapped`             | Ported tests: 0 / 113 (0.0%); Missing tests: 113 (100.0%); Adapted tests: 0 (0.0%) | `packages/telescope.md`    |
| `laravel/horizon`             | `horizon`       | `mapped`             | Ported tests: 0 / 180 (0.0%); Missing tests: 180 (100.0%); Adapted tests: 0 (0.0%) | `packages/horizon.md`      |
| `laravel/pail`                | `pail`          | `mapped`             | Ported tests: 0 / 38 (0.0%); Missing tests: 38 (100.0%); Adapted tests: 0 (0.0%)   | `packages/pail.md`         |
| `laravel/envoy`               | `envoy`         | `mapped`             | Ported tests: 0 / 20 (0.0%); Missing tests: 20 (100.0%); Adapted tests: 0 (0.0%)   | `packages/envoy.md`        |
| `laravel/wayfinder`           | `wayfinder`     | `mapped`             | Ported tests: 0 / 103 (0.0%); Missing tests: 103 (100.0%); Adapted tests: 0 (0.0%) | `packages/wayfinder.md`    |
| `laravel/sanctum`             | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `laravel/octane`              | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `laravel/pulse`               | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `laravel/folio`               | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `laravel/dusk`                | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `laravel/browser-kit-testing` | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `laravel/sail`                | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `laravel/installer`           | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `laravel/valet`               | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `laravel/vite-plugin`         | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `starter-kits`                | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |
| `laravel/workos`              | n/a             | `excluded-permanent` | no inventory                                                                       | n/a                        |

## Test Porting Summary

| Scope                 | Inventories | Upstream Tests | Ported Tests | Pending / Missing Tests | Adapted Tests |
| --------------------- | ----------: | -------------: | -----------: | ----------------------: | ------------: |
| All inventories       |          51 |          13066 |  1137 (8.7%) |           11510 (88.1%) |    419 (3.2%) |
| Framework inventories |          34 |           9009 | 1134 (12.6%) |            7456 (82.8%) |    419 (4.7%) |
| Package inventories   |          17 |           4057 |     3 (0.1%) |            4054 (99.9%) |      0 (0.0%) |

## Documentation Porting Summary

| Scope        | Doc Inventories | Upstream Sections | Ported Sections | Pending / Missing Sections | Adapted Sections | Excluded Sections |
| ------------ | --------------: | ----------------: | --------------: | -------------------------: | ---------------: | ----------------: |
| Laravel docs |               1 |               875 |     200 (22.9%) |                505 (57.7%) |        10 (1.1%) |       160 (18.3%) |

## Laravel Skeleton Demo Summary

| Scope           | Skeleton Inventories | Upstream Files | Ported Files | Pending / Missing Files | Adapted Files | Excluded Files |
| --------------- | -------------------: | -------------: | -----------: | ----------------------: | ------------: | -------------: |
| laravel/laravel |                    1 |             49 |   20 (40.8%) |                0 (0.0%) |    10 (20.4%) |     19 (38.8%) |

## Documentation Detail

### Top Pending Documentation Files

| Laravel Docs File           | Pending Sections |
| --------------------------- | ---------------: |
| `eloquent.md`               |               20 |
| `queries.md`                |               16 |
| `passport.md`               |               16 |
| `cashier-paddle.md`         |               15 |
| `billing.md`                |               14 |
| `prompts.md`                |               13 |
| `pennant.md`                |               11 |
| `notifications.md`          |               10 |
| `migrations.md`             |               10 |
| `mcp.md`                    |               10 |
| `eloquent-relationships.md` |               10 |
| `contributions.md`          |               10 |
| `boost.md`                  |               10 |
| `mail.md`                   |                9 |
| `http-tests.md`             |                9 |
| `deployment.md`             |                9 |
| `artisan.md`                |                9 |
| `ai-sdk.md`                 |                9 |
| `queues.md`                 |                8 |
| `packages.md`               |                8 |

### Classified Documentation Sections

| Laravel Docs File | Status   | Sections |
| ----------------- | -------- | -------: |
| `installation.md` | Adapted  |        4 |
| `lifecycle.md`    | Adapted  |        1 |
| `structure.md`    | Adapted  |        3 |
| `testing.md`      | Adapted  |        2 |
| `blade.md`        | Excluded |       13 |
| `dusk.md`         | Excluded |       17 |
| `folio.md`        | Excluded |       17 |
| `homestead.md`    | Excluded |       11 |
| `mix.md`          | Excluded |        2 |
| `octane.md`       | Excluded |       12 |
| `pint.md`         | Excluded |        6 |
| `pulse.md`        | Excluded |        7 |
| `sail.md`         | Excluded |       25 |
| `sanctum.md`      | Excluded |        8 |
| `starter-kits.md` | Excluded |        8 |
| `valet.md`        | Excluded |       12 |
| `views.md`        | Excluded |        6 |
| `vite.md`         | Excluded |       16 |

## Skeleton Detail

| Laravel Skeleton File                    | Status   |
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
| `artisan`                                | Excluded |
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
| `resources/views/welcome.blade.php`      | Excluded |
| `routes/console.php`                     | Adapted  |
| `vite.config.js`                         | Excluded |

## Feature Coverage

| Feature                      | Source                      | Bedrock Surface | Status     | Docs                        | Notes                                                                                                                                                                                                                                                                                                                                         |
| ---------------------------- | --------------------------- | --------------- | ---------- | --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Authentication guards        | `framework.auth`            | `auth`          | `ported`   | `packages/auth.md`          | Session, token, request, basic-auth, provider, and event surfaces are complete under classified parity; PHP runtime mechanics are covered by Go adaptation rules in divergences.yml.                                                                                                                                                          |
| Password broker              | `framework.auth`            | `auth`          | `ported`   | `packages/auth.md`          | Password reset tokens, throttling, callback notification override, notification sender hook, and reset-link events are complete under classified parity; PHP manager/config mechanics are covered by Go adaptation rules in divergences.yml.                                                                                                  |
| Broadcasting contracts       | `framework.broadcasting`    | `broadcasting`  | `ported`   | `packages/broadcasting.md`  | Server-side channel auth, broadcast event payload handling, Pusher/Redis/Ably broadcaster responses, and Pusher channel conventions are ported; Echo remains the client receiver package.                                                                                                                                                     |
| Command bus dispatch         | `framework.bus`             | `bus`           | `ported`   | `packages/bus.md`           | Command dispatch, queue routing, after-response dispatch, chains, batches, lifecycle events, and typed failure callbacks are complete under classified parity; PHP runtime mechanics are covered by Go adaptation rules in divergences.yml.                                                                                                   |
| Cache repository             | `framework.cache`           | `cache`         | `partial`  | `packages/cache.md`         | Cache abstractions are present, but upstream test parity is still missing.                                                                                                                                                                                                                                                                    |
| Concurrent task execution    | `framework.concurrency`     | `concurrency`   | `ported`   | `packages/concurrency.md`   | Executable Go equivalents cover distributed task execution, error propagation, typed error inspection, and result ordering; PHP process and serialized-closure mechanics are adapted to explicit Go tasks, goroutines, error returns, and panic recovery.                                                                                     |
| Conditional fluent execution | `framework.conditionable`   | `conditionable` | `ported`   | `packages/conditionable.md` | Executable Go equivalents cover when and unless; PHP dynamic no-argument and magic behavior is represented through explicit Go proxy constructors.                                                                                                                                                                                            |
| Configuration repository     | `framework.config`          | `config`        | `ported`   | `packages/config.md`        | Dot-notation reads, literal-dot key precedence, explicit nil values, write helpers, typed getters, collection wrapping, and explicit Go offset adaptations are complete under classified parity; PHP Macroable runtime mechanics are covered by Go adaptation rules in divergences.yml.                                                       |
| Console commands             | `framework.console`         | `console`       | `ported`   | `packages/console.md`       | Command registration and dispatch, signature parsing, typed input and output helpers, signal traps, scheduling primitives, and explicit command/scheduling mutexes are complete under classified parity; Symfony/PHP attribute runtime and dynamic Macroable/callable-string mechanics are covered by Go adaptation rules in divergences.yml. |
| Container bindings           | `framework.container`       | `container`     | `ported`   | `packages/container.md`     | Bindings, shared and scoped lifecycles, aliases, contextual bindings, tags, extenders, method bindings, and resolving callbacks are complete under classified parity; PHP reflection, attributes, callable strings, ArrayAccess, and environment selectors are covered by Go adaptation rules in divergences.yml.                             |
| Shared interfaces            | `framework.contracts`       | `contracts`     | `ported`   | `packages/contracts.md`     | Interface compliance is tracked through concrete package inventories.                                                                                                                                                                                                                                                                         |
| Cookie handling              | `framework.cookie`          | `cookie`        | `partial`  | `packages/cookie.md`        | Cookie package exists; upstream test parity is not complete.                                                                                                                                                                                                                                                                                  |
| Database query builder       | `framework.database`        | `database`      | `partial`  | `packages/database.md`      | Query, schema, migrations, and Eloquent surfaces exist; broad upstream parity remains incomplete.                                                                                                                                                                                                                                             |
| Eloquent-style ORM           | `framework.database`        | `database`      | `partial`  | `packages/database.md`      | ORM primitives exist with Go adaptations for attributes, relations, and errors.                                                                                                                                                                                                                                                               |
| Encrypter                    | `framework.encryption`      | `encryption`    | `partial`  | `packages/encryption.md`    | CBC and GCM encryption exist; upstream encryption tests are not fully matched.                                                                                                                                                                                                                                                                |
| Event dispatcher             | `framework.events`          | `events`        | `partial`  | `packages/events.md`        | Event dispatching exists; upstream test parity remains incomplete.                                                                                                                                                                                                                                                                            |
| Local filesystem             | `framework.filesystem`      | `filesystem`    | `partial`  | `packages/filesystem.md`    | Local filesystem operations are implemented; upstream parity tests are not complete.                                                                                                                                                                                                                                                          |
| Application foundation       | `framework.foundation`      | n/a             | `excluded` | n/a                         | Full Laravel application kernel is outside Bedrock product boundaries.                                                                                                                                                                                                                                                                        |
| Password hashing             | `framework.hashing`         | `hashing`       | `partial`  | `packages/hashing.md`       | Hashing drivers exist; upstream test parity is not complete.                                                                                                                                                                                                                                                                                  |
| HTTP primitives              | `framework.http`            | `httpx`         | `partial`  | `packages/httpx.md`         | HTTP primitives and test helpers exist; upstream test parity remains incomplete.                                                                                                                                                                                                                                                              |
| JSON schema builder          | `framework.json_schema`     | `jsonx`         | `partial`  | `packages/jsonx.md`         | Fluent JSON schema builder exists; upstream test parity is not complete.                                                                                                                                                                                                                                                                      |
| Logging channels             | `framework.log`             | `log`           | `partial`  | `packages/log.md`           | Driver-based structured logging exists; upstream log tests are not matched yet.                                                                                                                                                                                                                                                               |
| Mail messages                | `framework.mail`            | `mailx`         | `partial`  | `packages/mailx.md`         | Mail package exists; upstream mail test parity is incomplete.                                                                                                                                                                                                                                                                                 |
| Notification channels        | `framework.notifications`   | `notifications` | `partial`  | `packages/notifications.md` | Multi-channel notifications exist; upstream parity remains incomplete.                                                                                                                                                                                                                                                                        |
| Pagination                   | `framework.pagination`      | `pagination`    | `partial`  | `packages/pagination.md`    | Offset and cursor paginators exist; upstream parity is incomplete.                                                                                                                                                                                                                                                                            |
| Pipeline                     | `framework.pipeline`        | `pipeline`      | `partial`  | `packages/pipeline.md`      | Pipeline package exists; upstream test parity is not complete.                                                                                                                                                                                                                                                                                |
| Process runner               | `framework.process`         | n/a             | `missing`  | n/a                         | No Bedrock process package exists yet.                                                                                                                                                                                                                                                                                                        |
| Queue workers and drivers    | `framework.queue`           | `queue`         | `ported`   | `packages/queue.md`         | Queue workers, drivers, payload hooks, routing, FIFO metadata, delay options, after-commit precedence, and selected Go adaptations are complete under classified parity.                                                                                                                                                                      |
| Redis commands               | `framework.redis`           | `redis`         | `partial`  | `packages/redis.md`         | Redis command surface exists; upstream Redis tests are not matched yet.                                                                                                                                                                                                                                                                       |
| Routing                      | `framework.routing`         | `routing`       | `partial`  | `packages/routing.md`       | Routing package exists; upstream routing tests remain incomplete.                                                                                                                                                                                                                                                                             |
| Session stores               | `framework.session`         | `session`       | `partial`  | `packages/session.md`       | Session stores exist; upstream session tests are not matched yet.                                                                                                                                                                                                                                                                             |
| Support helpers              | `framework.support`         | `support`       | `partial`  | `packages/support.md`       | Several support helpers are ported; upstream support coverage is still partial.                                                                                                                                                                                                                                                               |
| String helpers               | `framework.support_str`     | `str`           | `partial`  | `packages/str.md`           | String helper parity has many ported tests, but missing cases remain.                                                                                                                                                                                                                                                                         |
| Lottery helper               | `framework.support_lottery` | `lottery`       | `partial`  | `packages/lottery.md`       | Lottery package exists; upstream Lottery tests are not matched yet.                                                                                                                                                                                                                                                                           |
| Testing utilities            | `framework.testing`         | n/a             | `missing`  | n/a                         | No standalone Bedrock testing package exists yet.                                                                                                                                                                                                                                                                                             |
| Translation loader           | `framework.translation`     | `translation`   | `partial`  | `packages/translation.md`   | Translation has some matched upstream tests; full parity remains incomplete.                                                                                                                                                                                                                                                                  |
| Validation rules             | `framework.validation`      | `validation`    | `partial`  | `packages/validation.md`    | Rule engine and many rules exist; upstream test parity is still broad and incomplete.                                                                                                                                                                                                                                                         |
| Collections                  | `framework.collections`     | `collection`    | `partial`  | `packages/collection.md`    | Fluent slice, key-value, lazy, arr, and kv packages are imported; legacy support helper cleanup is deferred.                                                                                                                                                                                                                                  |
| Macroable runtime            | `framework.macroable`       | n/a             | `excluded` | n/a                         | PHP runtime macro behaviour does not map cleanly to Go static method sets.                                                                                                                                                                                                                                                                    |
| Reflection helpers           | `framework.reflection`      | n/a             | `excluded` | n/a                         | PHP reflection helpers are not a Bedrock product surface.                                                                                                                                                                                                                                                                                     |
| Blade views                  | `framework.view`            | n/a             | `excluded` | n/a                         | Bedrock does not ship Laravel Blade or the view engine.                                                                                                                                                                                                                                                                                       |
| AI SDK                       | `package.ai`                | `ai/sdk`        | `partial`  | `packages/ai/sdk.md`        | Bedrock has an AI SDK package; upstream package parity is not complete.                                                                                                                                                                                                                                                                       |
| Agent tooling                | `package.boost`             | `ai/boost`      | `partial`  | `packages/ai/boost.md`      | Boost package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                             |
| Echo client                  | `package.echo`              | `echo`          | `partial`  | `packages/echo.md`          | Tracked by feature inventory because there is no stable upstream test inventory configured.                                                                                                                                                                                                                                                   |
| Fortify auth flows           | `package.fortify`           | `inception`     | `partial`  | `packages/fortify.md`       | Exposed through Inception with Fortify-specific documentation.                                                                                                                                                                                                                                                                                |
| Jetstream teams              | `package.jetstream`         | `inception`     | `partial`  | `packages/jetstream.md`     | Exposed through Inception with Jetstream-specific documentation.                                                                                                                                                                                                                                                                              |
| MCP server                   | `package.mcp`               | `ai/mcp`        | `partial`  | `packages/ai/mcp.md`        | MCP package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                               |
| OAuth server                 | `package.passport`          | `passport`      | `partial`  | `packages/passport.md`      | Passport package and docs exist; upstream package parity is not complete.                                                                                                                                                                                                                                                                     |
| Feature flags                | `package.pennant`           | `pennant`       | `partial`  | `packages/pennant.md`       | Pennant package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                           |
| Precognition client          | `package.precognition`      | `precognition`  | `partial`  | `packages/precognition.md`  | Tracked by feature inventory because there is no stable upstream test inventory configured.                                                                                                                                                                                                                                                   |
| CLI prompts                  | `package.prompts`           | `prompts`       | `partial`  | `packages/prompts.md`       | Prompt package exists with a small number of matched upstream tests.                                                                                                                                                                                                                                                                          |
| Reverb websocket server      | `package.reverb`            | `reverb`        | `partial`  | `packages/reverb.md`        | Reverb package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                            |
| Search indexing              | `package.scout`             | `scout`         | `partial`  | `packages/scout.md`         | Scout package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                             |
| Social OAuth                 | `package.socialite`         | `socialite`     | `partial`  | `packages/socialite.md`     | Socialite package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                         |
| Subscription billing         | `package.spark`             | `spark`         | `partial`  | `packages/spark.md`         | Spark is tracked against local Madora billing tests.                                                                                                                                                                                                                                                                                          |
| Telescope observability      | `package.telescope`         | `telescope`     | `partial`  | `packages/telescope.md`     | Telescope package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                         |
| Wayfinder route helpers      | `package.wayfinder`         | `wayfinder`     | `partial`  | `packages/wayfinder.md`     | Wayfinder package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                         |
| Queue monitoring             | `package.horizon`           | `horizon`       | `partial`  | `packages/horizon.md`       | Horizon is tracked as queue monitoring primitives; the browser dashboard is not ported.                                                                                                                                                                                                                                                       |
| Log tailing                  | `package.pail`              | `pail`          | `partial`  | `packages/pail.md`          | Pail is tracked as log parsing and filtering primitives for Go readers.                                                                                                                                                                                                                                                                       |
| Remote task runner           | `package.envoy`             | `envoy`         | `partial`  | `packages/envoy.md`         | Envoy is tracked as task planning and command execution primitives with injectable runners.                                                                                                                                                                                                                                                   |
| Sanctum tokens               | `package.sanctum`           | n/a             | `excluded` | n/a                         | Sanctum is intentionally not part of Bedrock.                                                                                                                                                                                                                                                                                                 |
| Browser automation           | `package.dusk`              | n/a             | `excluded` | n/a                         | Browser automation package is excluded; browser work should use Agent Browser.                                                                                                                                                                                                                                                                |

## Sources Without Test Inventories

| Source                 | Bedrock        | Tracking                     | Reason                                                                                                                    |
| ---------------------- | -------------- | ---------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `Illuminate\Contracts` | `contracts`    | Concrete package inventories | Interface-only component; compliance is verified through concrete package inventories.                                    |
| `laravel/echo`         | `echo`         | Feature inventory            | TypeScript client package; compliance is tracked by feature coverage until a stable upstream test inventory is available. |
| `laravel/precognition` | `precognition` | Feature inventory            | JavaScript package; compliance is tracked by feature coverage until a stable upstream test inventory is available.        |

## Permanent Exclusions

- `Illuminate\Foundation` - Product boundary; Bedrock does not ship Laravel's full application kernel/Foundation layer.
- `Illuminate\Macroable` - PHP runtime macro behaviour does not map cleanly to Go's static method set.
- `Illuminate\Reflection` - PHP reflection helpers are not a Bedrock product surface.
- `Illuminate\View` - Product boundary; Bedrock does not ship Laravel's Blade/view engine.
- `laravel/sanctum` - Product boundary; Sanctum is intentionally not part of Bedrock.
- `laravel/octane` - Product boundary; PHP worker/runtime management does not map to Bedrock.
- `laravel/pulse` - Product boundary; application monitoring dashboard is intentionally excluded.
- `laravel/folio` - Product boundary; page-based PHP routing is intentionally excluded.
- `laravel/dusk` - Product boundary; browser automation package is excluded. Any future browser work must use Agent Browser.
- `laravel/browser-kit-testing` - Product boundary; legacy BrowserKit compatibility is intentionally excluded.
- `laravel/sail` - Product boundary; Docker development environment tooling is intentionally excluded.
- `laravel/installer` - Product boundary; Laravel application installer tooling is intentionally excluded.
- `laravel/valet` - Product boundary; macOS local development service management is intentionally excluded.
- `laravel/vite-plugin` - Product boundary; Vite integration tooling is intentionally excluded.
- `starter-kits` - Product boundary; application starter kits are intentionally excluded.
- `laravel/workos` - Product boundary; WorkOS starter-kit integration is intentionally excluded.
