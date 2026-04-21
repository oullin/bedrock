# Laravel Compliance Report

Generated: 2026-04-21T08:27:14Z

Source of truth: services/compliance

## Compliance Workflow Requirements

- Every mapped Bedrock package compliance pass must include an upstream feature audit: `make sure we also have all the upstream features for <package>`.
- The audit compares upstream source contracts, public APIs, middleware, events, and runtime behavior against the Bedrock surface.
- Bedrock-equivalent features need executable parity tests; PHP-only or intentionally different behavior belongs in `services/compliance/divergences.yml`.
- `services/compliance/features.yml` must contain feature audit coverage for every mapped Bedrock surface.

## Inventories

| Inventory                                 | Count Context                                                                                                                                                                                                              | Upstream Tests | Ported Tests | Missing Tests | Adapted Tests |
| ----------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------: | -----------: | ------------: | ------------: |
| inventories/framework-auth.txt            | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Auth`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                    |            257 |   54 (21.0%) |      0 (0.0%) |   203 (79.0%) |
| inventories/framework-broadcasting.txt    | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Broadcasting`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.            |             55 |  55 (100.0%) |      0 (0.0%) |      0 (0.0%) |
| inventories/framework-bus.txt             | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Bus`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                     |             66 |     0 (0.0%) |   66 (100.0%) |      0 (0.0%) |
| inventories/framework-cache.txt           | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Cache`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                   |            313 |     0 (0.0%) |  313 (100.0%) |      0 (0.0%) |
| inventories/framework-concurrency.txt     | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Integration/Concurrency`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`. |             12 |     0 (0.0%) |   12 (100.0%) |      0 (0.0%) |
| inventories/framework-conditionable.txt   | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Conditionable`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.           |              2 |     0 (0.0%) |    2 (100.0%) |      0 (0.0%) |
| inventories/framework-config.txt          | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Config`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                  |             33 |     0 (0.0%) |   33 (100.0%) |      0 (0.0%) |
| inventories/framework-console.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Console`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                 |            156 |     0 (0.0%) |  156 (100.0%) |      0 (0.0%) |
| inventories/framework-container.txt       | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Container`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.               |            193 |     0 (0.0%) |  193 (100.0%) |      0 (0.0%) |
| inventories/framework-cookie.txt          | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Cookie`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                  |             26 |     0 (0.0%) |   26 (100.0%) |      0 (0.0%) |
| inventories/framework-database.txt        | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Database`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                |           2626 |     0 (0.0%) | 2626 (100.0%) |      0 (0.0%) |
| inventories/framework-encryption.txt      | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Encryption`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.              |             27 |     0 (0.0%) |   27 (100.0%) |      0 (0.0%) |
| inventories/framework-events.txt          | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Events`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                  |             76 |     0 (0.0%) |   76 (100.0%) |      0 (0.0%) |
| inventories/framework-filesystem.txt      | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Filesystem`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.              |            134 |     0 (0.0%) |  134 (100.0%) |      0 (0.0%) |
| inventories/framework-hashing.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Hashing`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                 |             13 |     0 (0.0%) |   13 (100.0%) |      0 (0.0%) |
| inventories/framework-http.txt            | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Http`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                    |            488 |     0 (0.0%) |  488 (100.0%) |      0 (0.0%) |
| inventories/framework-json-schema.txt     | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/JsonSchema`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.              |             47 |     0 (0.0%) |   47 (100.0%) |      0 (0.0%) |
| inventories/framework-log.txt             | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Log`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                     |             98 |     0 (0.0%) |   98 (100.0%) |      0 (0.0%) |
| inventories/framework-mail.txt            | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Mail`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                    |            135 |     0 (0.0%) |  135 (100.0%) |      0 (0.0%) |
| inventories/framework-notifications.txt   | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Notifications`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.           |             71 |     0 (0.0%) |   71 (100.0%) |      0 (0.0%) |
| inventories/framework-pagination.txt      | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Pagination`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.              |             48 |     0 (0.0%) |   48 (100.0%) |      0 (0.0%) |
| inventories/framework-pipeline.txt        | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Pipeline`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                |             22 |     0 (0.0%) |   22 (100.0%) |      0 (0.0%) |
| inventories/framework-process.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Process`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                 |             54 |     0 (0.0%) |   54 (100.0%) |      0 (0.0%) |
| inventories/framework-queue.txt           | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Queue`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                   |            204 |  161 (78.9%) |    43 (21.1%) |      0 (0.0%) |
| inventories/framework-redis.txt           | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Redis`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                   |            103 |     0 (0.0%) |  103 (100.0%) |      0 (0.0%) |
| inventories/framework-routing.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Routing`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                 |            315 |     0 (0.0%) |  315 (100.0%) |      0 (0.0%) |
| inventories/framework-session.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Session`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                 |             93 |     0 (0.0%) |   93 (100.0%) |      0 (0.0%) |
| inventories/framework-support.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Support`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                 |           1554 |   129 (8.3%) |  1425 (91.7%) |      0 (0.0%) |
| inventories/framework-support-str.txt     | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Support`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                 |            115 |   69 (60.0%) |    46 (40.0%) |      0 (0.0%) |
| inventories/framework-support-lottery.txt | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Support`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                 |             14 |     0 (0.0%) |   14 (100.0%) |      0 (0.0%) |
| inventories/framework-testing.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Testing`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                 |            381 |     0 (0.0%) |  381 (100.0%) |      0 (0.0%) |
| inventories/framework-translation.txt     | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Translation`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.             |             52 |    9 (17.3%) |    43 (82.7%) |      0 (0.0%) |
| inventories/framework-validation.txt      | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Validation`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.              |            597 |    34 (5.7%) |   563 (94.3%) |      0 (0.0%) |
| inventories/framework-collections.txt     | Upstream tests are non-comment entries from `https://github.com/laravel/framework/tree/13.x/tests/Support`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                 |            629 |    21 (3.3%) |   608 (96.7%) |      0 (0.0%) |
| inventories/package-ai.txt                | Upstream tests are non-comment entries from `https://github.com/laravel/ai/tree/0.x/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                                 |            774 |     0 (0.0%) |  774 (100.0%) |      0 (0.0%) |
| inventories/package-boost.txt             | Upstream tests are non-comment entries from `https://github.com/laravel/boost/tree/main/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                             |            671 |     0 (0.0%) |  671 (100.0%) |      0 (0.0%) |
| inventories/package-fortify.txt           | Upstream tests are non-comment entries from `https://github.com/laravel/fortify/tree/1.x/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                            |             91 |     0 (0.0%) |   91 (100.0%) |      0 (0.0%) |
| inventories/package-jetstream.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/jetstream/tree/5.x/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                          |             43 |     0 (0.0%) |   43 (100.0%) |      0 (0.0%) |
| inventories/package-mcp.txt               | Upstream tests are non-comment entries from `https://github.com/laravel/mcp/tree/main/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                               |            624 |     0 (0.0%) |  624 (100.0%) |      0 (0.0%) |
| inventories/package-passport.txt          | Upstream tests are non-comment entries from `https://github.com/laravel/passport/tree/13.x/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                          |            211 |     0 (0.0%) |  211 (100.0%) |      0 (0.0%) |
| inventories/package-pennant.txt           | Upstream tests are non-comment entries from `https://github.com/laravel/pennant/tree/1.x/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                            |            186 |     0 (0.0%) |  186 (100.0%) |      0 (0.0%) |
| inventories/package-prompts.txt           | Upstream tests are non-comment entries from `https://github.com/laravel/prompts/tree/main/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                           |            313 |     3 (1.0%) |   310 (99.0%) |      0 (0.0%) |
| inventories/package-reverb.txt            | Upstream tests are non-comment entries from `https://github.com/laravel/reverb/tree/main/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                            |            232 |     0 (0.0%) |  232 (100.0%) |      0 (0.0%) |
| inventories/package-scout.txt             | Upstream tests are non-comment entries from `https://github.com/laravel/scout/tree/11.x/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                             |            252 |     0 (0.0%) |  252 (100.0%) |      0 (0.0%) |
| inventories/package-socialite.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/socialite/tree/5.x/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                          |             40 |     0 (0.0%) |   40 (100.0%) |      0 (0.0%) |
| inventories/package-spark.txt             | Upstream tests are non-comment entries from `/Users/gocanto/Sites/madora/tests/Feature/Billing`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                            |            166 |     0 (0.0%) |  166 (100.0%) |      0 (0.0%) |
| inventories/package-telescope.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/telescope/tree/5.x/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                          |            113 |     0 (0.0%) |  113 (100.0%) |      0 (0.0%) |
| inventories/package-horizon.txt           | Upstream tests are non-comment entries from `https://github.com/laravel/horizon/tree/5.x/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                            |            180 |     0 (0.0%) |  180 (100.0%) |      0 (0.0%) |
| inventories/package-pail.txt              | Upstream tests are non-comment entries from `https://github.com/laravel/pail/tree/main/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                              |             38 |     0 (0.0%) |   38 (100.0%) |      0 (0.0%) |
| inventories/package-envoy.txt             | Upstream tests are non-comment entries from `https://github.com/laravel/envoy/tree/2.x/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                              |             20 |     0 (0.0%) |   20 (100.0%) |      0 (0.0%) |
| inventories/package-wayfinder.txt         | Upstream tests are non-comment entries from `https://github.com/laravel/wayfinder/tree/main/tests`; ported tests match Go test markers; adapted tests match `services/compliance/divergences.yml`.                         |            103 |     0 (0.0%) |  103 (100.0%) |      0 (0.0%) |

## Framework Coverage

| Laravel Source               | Bedrock Surface | Port Status          | Tests                                                                                 | Docs                        |
| ---------------------------- | --------------- | -------------------- | ------------------------------------------------------------------------------------- | --------------------------- |
| `Illuminate\Auth`            | `auth`          | `mapped`             | Ported tests: 54 / 257 (21.0%); Missing tests: 0 (0.0%); Adapted tests: 203 (79.0%)   | `packages/auth.md`          |
| `Illuminate\Broadcasting`    | `broadcasting`  | `mapped`             | Ported tests: 55 / 55 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)      | `packages/broadcasting.md`  |
| `Illuminate\Bus`             | `bus`           | `mapped`             | Ported tests: 0 / 66 (0.0%); Missing tests: 66 (100.0%); Adapted tests: 0 (0.0%)      | `packages/bus.md`           |
| `Illuminate\Cache`           | `cache`         | `mapped`             | Ported tests: 0 / 313 (0.0%); Missing tests: 313 (100.0%); Adapted tests: 0 (0.0%)    | `packages/cache.md`         |
| `Illuminate\Concurrency`     | `concurrency`   | `mapped`             | Ported tests: 0 / 12 (0.0%); Missing tests: 12 (100.0%); Adapted tests: 0 (0.0%)      | `packages/concurrency.md`   |
| `Illuminate\Conditionable`   | `conditionable` | `mapped`             | Ported tests: 0 / 2 (0.0%); Missing tests: 2 (100.0%); Adapted tests: 0 (0.0%)        | `packages/conditionable.md` |
| `Illuminate\Config`          | `config`        | `mapped`             | Ported tests: 0 / 33 (0.0%); Missing tests: 33 (100.0%); Adapted tests: 0 (0.0%)      | `packages/config.md`        |
| `Illuminate\Console`         | n/a             | `missing`            | Ported tests: 0 / 156 (0.0%); Missing tests: 156 (100.0%); Adapted tests: 0 (0.0%)    | n/a                         |
| `Illuminate\Container`       | `container`     | `mapped`             | Ported tests: 0 / 193 (0.0%); Missing tests: 193 (100.0%); Adapted tests: 0 (0.0%)    | `packages/container.md`     |
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
| `Illuminate\Queue`           | `queue`         | `mapped`             | Ported tests: 161 / 204 (78.9%); Missing tests: 43 (21.1%); Adapted tests: 0 (0.0%)   | `packages/queue.md`         |
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
| All inventories       |          51 |          13066 |   535 (4.1%) |           12328 (94.4%) |    203 (1.6%) |
| Framework inventories |          34 |           9009 |   532 (5.9%) |            8274 (91.8%) |    203 (2.3%) |
| Package inventories   |          17 |           4057 |     3 (0.1%) |            4054 (99.9%) |      0 (0.0%) |

## Documentation Porting Summary

| Scope        | Doc Inventories | Upstream Sections | Ported Sections | Pending / Missing Sections | Adapted Sections | Excluded Sections |
| ------------ | --------------: | ----------------: | --------------: | -------------------------: | ---------------: | ----------------: |
| Laravel docs |               1 |               875 |     197 (22.5%) |                508 (58.1%) |        10 (1.1%) |       160 (18.3%) |

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
| `prompts.md`                |               14 |
| `billing.md`                |               14 |
| `pennant.md`                |               11 |
| `notifications.md`          |               10 |
| `migrations.md`             |               10 |
| `mcp.md`                    |               10 |
| `eloquent-relationships.md` |               10 |
| `contributions.md`          |               10 |
| `boost.md`                  |               10 |
| `artisan.md`                |               10 |
| `mail.md`                   |                9 |
| `http-tests.md`             |                9 |
| `deployment.md`             |                9 |
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

| Feature                      | Source                      | Bedrock Surface | Status     | Docs                        | Notes                                                                                                                                                                                     |
| ---------------------------- | --------------------------- | --------------- | ---------- | --------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Authentication guards        | `framework.auth`            | `auth`          | `partial`  | `packages/auth.md`          | Session, token, request, basic-auth, provider, and event surfaces exist with a current upstream feature audit; remaining gaps are PHP-specific mechanics or broader test parity.          |
| Password broker              | `framework.auth`            | `auth`          | `partial`  | `packages/auth.md`          | Password reset tokens, throttling, callback notification override, notification sender hook, and reset-link events exist; manager/config wiring remains adapted to Go.                    |
| Broadcasting contracts       | `framework.broadcasting`    | `broadcasting`  | `ported`   | `packages/broadcasting.md`  | Server-side channel auth, broadcast event payload handling, Pusher/Redis/Ably broadcaster responses, and Pusher channel conventions are ported; Echo remains the client receiver package. |
| Command bus dispatch         | `framework.bus`             | `bus`           | `partial`  | `packages/bus.md`           | Bus package exists with command dispatch and middleware; upstream test parity is not complete.                                                                                            |
| Cache repository             | `framework.cache`           | `cache`         | `partial`  | `packages/cache.md`         | Cache abstractions are present, but upstream test parity is still missing.                                                                                                                |
| Concurrent task execution    | `framework.concurrency`     | `concurrency`   | `partial`  | `packages/concurrency.md`   | Bedrock has a concurrency package; upstream concurrency tests have not been matched yet.                                                                                                  |
| Conditional fluent execution | `framework.conditionable`   | `conditionable` | `partial`  | `packages/conditionable.md` | Go fluent proxy exists; PHP dynamic behaviour is adapted to Go.                                                                                                                           |
| Configuration repository     | `framework.config`          | `config`        | `partial`  | `packages/config.md`        | Config package exists; upstream test parity is not complete.                                                                                                                              |
| Console commands             | `framework.console`         | n/a             | `missing`  | n/a                         | No Bedrock console package exists yet.                                                                                                                                                    |
| Container bindings           | `framework.container`       | `container`     | `partial`  | `packages/container.md`     | Container package exists; upstream container tests have not been matched yet.                                                                                                             |
| Shared interfaces            | `framework.contracts`       | `contracts`     | `ported`   | `packages/contracts.md`     | Interface compliance is tracked through concrete package inventories.                                                                                                                     |
| Cookie handling              | `framework.cookie`          | `cookie`        | `partial`  | `packages/cookie.md`        | Cookie package exists; upstream test parity is not complete.                                                                                                                              |
| Database query builder       | `framework.database`        | `database`      | `partial`  | `packages/database.md`      | Query, schema, migrations, and Eloquent surfaces exist; broad upstream parity remains incomplete.                                                                                         |
| Eloquent-style ORM           | `framework.database`        | `database`      | `partial`  | `packages/database.md`      | ORM primitives exist with Go adaptations for attributes, relations, and errors.                                                                                                           |
| Encrypter                    | `framework.encryption`      | `encryption`    | `partial`  | `packages/encryption.md`    | CBC and GCM encryption exist; upstream encryption tests are not fully matched.                                                                                                            |
| Event dispatcher             | `framework.events`          | `events`        | `partial`  | `packages/events.md`        | Event dispatching exists; upstream test parity remains incomplete.                                                                                                                        |
| Local filesystem             | `framework.filesystem`      | `filesystem`    | `partial`  | `packages/filesystem.md`    | Local filesystem operations are implemented; upstream parity tests are not complete.                                                                                                      |
| Application foundation       | `framework.foundation`      | n/a             | `excluded` | n/a                         | Full Laravel application kernel is outside Bedrock product boundaries.                                                                                                                    |
| Password hashing             | `framework.hashing`         | `hashing`       | `partial`  | `packages/hashing.md`       | Hashing drivers exist; upstream test parity is not complete.                                                                                                                              |
| HTTP primitives              | `framework.http`            | `httpx`         | `partial`  | `packages/httpx.md`         | HTTP primitives and test helpers exist; upstream test parity remains incomplete.                                                                                                          |
| JSON schema builder          | `framework.json_schema`     | `jsonx`         | `partial`  | `packages/jsonx.md`         | Fluent JSON schema builder exists; upstream test parity is not complete.                                                                                                                  |
| Logging channels             | `framework.log`             | `log`           | `partial`  | `packages/log.md`           | Driver-based structured logging exists; upstream log tests are not matched yet.                                                                                                           |
| Mail messages                | `framework.mail`            | `mailx`         | `partial`  | `packages/mailx.md`         | Mail package exists; upstream mail test parity is incomplete.                                                                                                                             |
| Notification channels        | `framework.notifications`   | `notifications` | `partial`  | `packages/notifications.md` | Multi-channel notifications exist; upstream parity remains incomplete.                                                                                                                    |
| Pagination                   | `framework.pagination`      | `pagination`    | `partial`  | `packages/pagination.md`    | Offset and cursor paginators exist; upstream parity is incomplete.                                                                                                                        |
| Pipeline                     | `framework.pipeline`        | `pipeline`      | `partial`  | `packages/pipeline.md`      | Pipeline package exists; upstream test parity is not complete.                                                                                                                            |
| Process runner               | `framework.process`         | n/a             | `missing`  | n/a                         | No Bedrock process package exists yet.                                                                                                                                                    |
| Queue workers and drivers    | `framework.queue`           | `queue`         | `partial`  | `packages/queue.md`         | Queue has the strongest upstream test match so far, but missing tests remain.                                                                                                             |
| Redis commands               | `framework.redis`           | `redis`         | `partial`  | `packages/redis.md`         | Redis command surface exists; upstream Redis tests are not matched yet.                                                                                                                   |
| Routing                      | `framework.routing`         | `routing`       | `partial`  | `packages/routing.md`       | Routing package exists; upstream routing tests remain incomplete.                                                                                                                         |
| Session stores               | `framework.session`         | `session`       | `partial`  | `packages/session.md`       | Session stores exist; upstream session tests are not matched yet.                                                                                                                         |
| Support helpers              | `framework.support`         | `support`       | `partial`  | `packages/support.md`       | Several support helpers are ported; upstream support coverage is still partial.                                                                                                           |
| String helpers               | `framework.support_str`     | `str`           | `partial`  | `packages/str.md`           | String helper parity has many ported tests, but missing cases remain.                                                                                                                     |
| Lottery helper               | `framework.support_lottery` | `lottery`       | `partial`  | `packages/lottery.md`       | Lottery package exists; upstream Lottery tests are not matched yet.                                                                                                                       |
| Testing utilities            | `framework.testing`         | n/a             | `missing`  | n/a                         | No standalone Bedrock testing package exists yet.                                                                                                                                         |
| Translation loader           | `framework.translation`     | `translation`   | `partial`  | `packages/translation.md`   | Translation has some matched upstream tests; full parity remains incomplete.                                                                                                              |
| Validation rules             | `framework.validation`      | `validation`    | `partial`  | `packages/validation.md`    | Rule engine and many rules exist; upstream test parity is still broad and incomplete.                                                                                                     |
| Collections                  | `framework.collections`     | `collection`    | `partial`  | `packages/collection.md`    | Fluent slice, key-value, lazy, arr, and kv packages are imported; legacy support helper cleanup is deferred.                                                                              |
| Macroable runtime            | `framework.macroable`       | n/a             | `excluded` | n/a                         | PHP runtime macro behaviour does not map cleanly to Go static method sets.                                                                                                                |
| Reflection helpers           | `framework.reflection`      | n/a             | `excluded` | n/a                         | PHP reflection helpers are not a Bedrock product surface.                                                                                                                                 |
| Blade views                  | `framework.view`            | n/a             | `excluded` | n/a                         | Bedrock does not ship Laravel Blade or the view engine.                                                                                                                                   |
| AI SDK                       | `package.ai`                | `ai/sdk`        | `partial`  | `packages/ai/sdk.md`        | Bedrock has an AI SDK package; upstream package parity is not complete.                                                                                                                   |
| Agent tooling                | `package.boost`             | `ai/boost`      | `partial`  | `packages/ai/boost.md`      | Boost package exists; upstream package parity remains incomplete.                                                                                                                         |
| Echo client                  | `package.echo`              | `echo`          | `partial`  | `packages/echo.md`          | Tracked by feature inventory because there is no stable upstream test inventory configured.                                                                                               |
| Fortify auth flows           | `package.fortify`           | `inception`     | `partial`  | `packages/fortify.md`       | Exposed through Inception with Fortify-specific documentation.                                                                                                                            |
| Jetstream teams              | `package.jetstream`         | `inception`     | `partial`  | `packages/jetstream.md`     | Exposed through Inception with Jetstream-specific documentation.                                                                                                                          |
| MCP server                   | `package.mcp`               | `ai/mcp`        | `partial`  | `packages/ai/mcp.md`        | MCP package exists; upstream package parity remains incomplete.                                                                                                                           |
| OAuth server                 | `package.passport`          | `passport`      | `partial`  | `packages/passport.md`      | Passport package and docs exist; upstream package parity is not complete.                                                                                                                 |
| Feature flags                | `package.pennant`           | `pennant`       | `partial`  | `packages/pennant.md`       | Pennant package exists; upstream package parity remains incomplete.                                                                                                                       |
| Precognition client          | `package.precognition`      | `precognition`  | `partial`  | `packages/precognition.md`  | Tracked by feature inventory because there is no stable upstream test inventory configured.                                                                                               |
| CLI prompts                  | `package.prompts`           | `prompts`       | `partial`  | `packages/prompts.md`       | Prompt package exists with a small number of matched upstream tests.                                                                                                                      |
| Reverb websocket server      | `package.reverb`            | `reverb`        | `partial`  | `packages/reverb.md`        | Reverb package exists; upstream package parity remains incomplete.                                                                                                                        |
| Search indexing              | `package.scout`             | `scout`         | `partial`  | `packages/scout.md`         | Scout package exists; upstream package parity remains incomplete.                                                                                                                         |
| Social OAuth                 | `package.socialite`         | `socialite`     | `partial`  | `packages/socialite.md`     | Socialite package exists; upstream package parity remains incomplete.                                                                                                                     |
| Subscription billing         | `package.spark`             | `spark`         | `partial`  | `packages/spark.md`         | Spark is tracked against local Madora billing tests.                                                                                                                                      |
| Telescope observability      | `package.telescope`         | `telescope`     | `partial`  | `packages/telescope.md`     | Telescope package exists; upstream package parity remains incomplete.                                                                                                                     |
| Wayfinder route helpers      | `package.wayfinder`         | `wayfinder`     | `partial`  | `packages/wayfinder.md`     | Wayfinder package exists; upstream package parity remains incomplete.                                                                                                                     |
| Queue monitoring             | `package.horizon`           | `horizon`       | `partial`  | `packages/horizon.md`       | Horizon is tracked as queue monitoring primitives; the browser dashboard is not ported.                                                                                                   |
| Log tailing                  | `package.pail`              | `pail`          | `partial`  | `packages/pail.md`          | Pail is tracked as log parsing and filtering primitives for Go readers.                                                                                                                   |
| Remote task runner           | `package.envoy`             | `envoy`         | `partial`  | `packages/envoy.md`         | Envoy is tracked as task planning and command execution primitives with injectable runners.                                                                                               |
| Sanctum tokens               | `package.sanctum`           | n/a             | `excluded` | n/a                         | Sanctum is intentionally not part of Bedrock.                                                                                                                                             |
| Browser automation           | `package.dusk`              | n/a             | `excluded` | n/a                         | Browser automation package is excluded; browser work should use Agent Browser.                                                                                                            |

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
