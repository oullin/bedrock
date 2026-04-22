# Upstream Compliance Report

Generated: 2026-04-22T12:52:03Z

Source of truth: services/compliance

## Compliance Workflow Requirements

- Every mapped Bedrock package compliance pass must include an upstream feature audit: `make sure we also have all the upstream features for <package>`.
- The audit compares upstream source contracts, public APIs, middleware, events, and runtime behavior against the Bedrock surface.
- Bedrock-equivalent features need executable parity tests; PHP-only or intentionally different behavior belongs in `services/compliance/divergences.yml`.
- `services/compliance/features.yml` must contain feature audit coverage for every mapped Bedrock surface.

## Compliance Command Center

Compliance target: classified parity. An upstream item is compliant when it is ported, adapted with a divergence rationale, or excluded at an approved product boundary.

| Area                  |    Classified | Missing | Compliance Status | Fastest Next Move                                                                         |
| --------------------- | ------------: | ------: | ----------------- | ----------------------------------------------------------------------------------------- |
| Tests                 | 10574 / 13066 |    2492 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations.               |
| Documentation         |     370 / 875 |     505 | Not compliant     | Port relevant sections, adapt Go-specific sections, or exclude product-boundary sections. |
| Upstream skeleton demo |       49 / 49 |       0 | Complete          | Keep complete while upstream skeleton changes.                                            |
| Feature audits        |       62 / 62 |       0 | Complete          | Add audits for missing surfaces, then move partial surfaces to ported as parity closes.   |
| Overall               | 11055 / 14052 |    2997 | Not compliant     | Start with the Critical Path below; close the largest missing surfaces first.             |

## Critical Path

| Rank | Surface                                 | Missing | Classified | Status        | Next Action                                                                               |
| ---: | --------------------------------------- | ------: | ---------: | ------------- | ----------------------------------------------------------------------------------------- |
|    1 | `upstream/ai` -> `ai/sdk`                |     551 |  223 / 774 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    2 | `upstream/boost` -> `ai/boost`           |     518 |  153 / 671 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    3 | Upstream documentation                   |     505 |  370 / 875 | Not compliant | Port relevant sections, adapt Go-specific sections, or exclude product-boundary sections. |
|    4 | `Framework\Validation` -> `validation` |     411 |  186 / 597 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    5 | `upstream/mcp` -> `ai/mcp`               |     410 |  214 / 624 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    6 | `Framework\Routing` -> `routing`       |     173 |  142 / 315 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    7 | `upstream/search` -> `search`              |     170 |   82 / 252 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    8 | `upstream/websockets` -> `websockets`            |      90 |  142 / 232 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|    9 | `upstream/oauthserver` -> `oauthserver`        |      76 |  135 / 211 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   10 | `upstream/prompts` -> `prompts`          |      49 |  264 / 313 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |
|   11 | `upstream/featureflags` -> `featureflags`          |      44 |  142 / 186 | Not compliant | Port executable equivalents or add divergence overrides for Go adaptations.               |

## Inventory Compliance Ledger

These inventory paths are tracking files, not compliant code paths. A row is compliant only when `Missing Tests` is 0.

| Inventory Path                            | What This Path Is                                      | Source                                                 |  Classified | Missing Tests | Compliance Status | Next Action                                                                 |
| ----------------------------------------- | ------------------------------------------------------ | ------------------------------------------------------ | ----------: | ------------: | ----------------- | --------------------------------------------------------------------------- |
| inventories/framework-auth.txt            | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Auth`                    |   257 / 257 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-broadcasting.txt    | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Broadcasting`            |     55 / 55 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-bus.txt             | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Bus`                     |     66 / 66 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-cache.txt           | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Cache`                   |   313 / 313 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-concurrency.txt     | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Integration/Concurrency` |     12 / 12 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-conditionable.txt   | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Conditionable`           |       2 / 2 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-config.txt          | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Config`                  |     33 / 33 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-console.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Console`                 |   156 / 156 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-container.txt       | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Container`               |   193 / 193 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-cookie.txt          | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Cookie`                  |     26 / 26 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-database.txt        | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Database`                | 2626 / 2626 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-encryption.txt      | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Encryption`              |     27 / 27 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-events.txt          | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Events`                  |     76 / 76 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-filesystem.txt      | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Filesystem`              |   134 / 134 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-hashing.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Hashing`                 |     13 / 13 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-http.txt            | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Http`                    |   488 / 488 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-json-schema.txt     | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/JsonSchema`              |     47 / 47 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-log.txt             | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Log`                     |     98 / 98 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-mail.txt            | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Mail`                    |   135 / 135 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-notifications.txt   | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Notifications`           |     71 / 71 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-pagination.txt      | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Pagination`              |     48 / 48 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-pipeline.txt        | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Pipeline`                |     22 / 22 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-process.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Process`                 |     54 / 54 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-queue.txt           | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Queue`                   |   204 / 204 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-redis.txt           | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Redis`                   |   103 / 103 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-routing.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Routing`                 |   142 / 315 |           173 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-session.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Session`                 |     93 / 93 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-support.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Support`                 | 1554 / 1554 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-support-str.txt     | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Support`                 |   115 / 115 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-support-lottery.txt | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Support`                 |     14 / 14 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-testing.txt         | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Testing`                 |   381 / 381 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-translation.txt     | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Translation`             |     52 / 52 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/framework-validation.txt      | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Validation`              |   186 / 597 |           411 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/framework-collections.txt     | Tracking file for upstream tests, not a compliant path | `upstream/framework@13.x:tests/Support`                 |   629 / 629 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/package-ai.txt                | Tracking file for upstream tests, not a compliant path | `upstream/ai@0.x:tests`                                 |   223 / 774 |           551 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-boost.txt             | Tracking file for upstream tests, not a compliant path | `upstream/boost@main:tests`                             |   153 / 671 |           518 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-authflows.txt           | Tracking file for upstream tests, not a compliant path | `upstream/authflows@1.x:tests`                            |     91 / 91 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/package-authkit.txt         | Tracking file for upstream tests, not a compliant path | `upstream/authkit@5.x:tests`                          |     43 / 43 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/package-mcp.txt               | Tracking file for upstream tests, not a compliant path | `upstream/mcp@main:tests`                               |   214 / 624 |           410 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-oauthserver.txt          | Tracking file for upstream tests, not a compliant path | `upstream/oauthserver@13.x:tests`                          |   135 / 211 |            76 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-featureflags.txt           | Tracking file for upstream tests, not a compliant path | `upstream/featureflags@1.x:tests`                            |   142 / 186 |            44 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-prompts.txt           | Tracking file for upstream tests, not a compliant path | `upstream/prompts@main:tests`                           |   264 / 313 |            49 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-websockets.txt            | Tracking file for upstream tests, not a compliant path | `upstream/websockets@main:tests`                            |   142 / 232 |            90 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-search.txt             | Tracking file for upstream tests, not a compliant path | `upstream/search@11.x:tests`                             |    82 / 252 |           170 | Not compliant     | Port executable equivalents or add divergence overrides for Go adaptations. |
| inventories/package-socialauth.txt         | Tracking file for upstream tests, not a compliant path | `upstream/socialauth@5.x:tests`                          |     40 / 40 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/package-billing.txt             | Tracking file for upstream tests, not a compliant path | `/Users/gocanto/Sites/madora/tests/Feature/Billing`    |   166 / 166 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/package-debugbar.txt         | Tracking file for upstream tests, not a compliant path | `upstream/debugbar@5.x:tests`                          |   113 / 113 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/package-jobqueue.txt           | Tracking file for upstream tests, not a compliant path | `upstream/jobqueue@5.x:tests`                            |   180 / 180 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/package-logtail.txt              | Tracking file for upstream tests, not a compliant path | `upstream/logtail@main:tests`                              |     38 / 38 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/package-remotetasks.txt             | Tracking file for upstream tests, not a compliant path | `upstream/remotetasks@2.x:tests`                              |     20 / 20 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |
| inventories/package-routegen.txt         | Tracking file for upstream tests, not a compliant path | `upstream/routegen@main:tests`                         |   103 / 103 |             0 | Complete          | No action; every upstream entry is classified as ported or adapted.         |

## Framework Coverage

| Upstream Source               | Bedrock Surface | Port Status          | Tests                                                                                  | Docs                        |
| ---------------------------- | --------------- | -------------------- | -------------------------------------------------------------------------------------- | --------------------------- |
| `Framework\Auth`            | `auth`          | `mapped`             | Ported tests: 54 / 257 (21.0%); Missing tests: 0 (0.0%); Adapted tests: 203 (79.0%)    | `packages/auth.md`          |
| `Framework\Broadcasting`    | `broadcasting`  | `mapped`             | Ported tests: 55 / 55 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)       | `packages/broadcasting.md`  |
| `Framework\Bus`             | `bus`           | `mapped`             | Ported tests: 54 / 66 (81.8%); Missing tests: 0 (0.0%); Adapted tests: 12 (18.2%)      | `packages/bus.md`           |
| `Framework\Cache`           | `cache`         | `mapped`             | Ported tests: 217 / 313 (69.3%); Missing tests: 0 (0.0%); Adapted tests: 96 (30.7%)    | `packages/cache.md`         |
| `Framework\Concurrency`     | `concurrency`   | `mapped`             | Ported tests: 5 / 12 (41.7%); Missing tests: 0 (0.0%); Adapted tests: 7 (58.3%)        | `packages/concurrency.md`   |
| `Framework\Conditionable`   | `conditionable` | `mapped`             | Ported tests: 2 / 2 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)         | `packages/conditionable.md` |
| `Framework\Config`          | `config`        | `mapped`             | Ported tests: 32 / 33 (97.0%); Missing tests: 0 (0.0%); Adapted tests: 1 (3.0%)        | `packages/config.md`        |
| `Framework\Console`         | `console`       | `mapped`             | Ported tests: 156 / 156 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)     | `packages/console.md`       |
| `Framework\Container`       | `container`     | `mapped`             | Ported tests: 102 / 193 (52.8%); Missing tests: 0 (0.0%); Adapted tests: 91 (47.2%)    | `packages/container.md`     |
| `Framework\Contracts`       | `contracts`     | `mapped`             | no inventory                                                                           | `packages/contracts.md`     |
| `Framework\Cookie`          | `cookie`        | `mapped`             | Ported tests: 24 / 26 (92.3%); Missing tests: 0 (0.0%); Adapted tests: 2 (7.7%)        | `packages/cookie.md`        |
| `Framework\Database`        | `database`      | `mapped`             | Ported tests: 0 / 2626 (0.0%); Missing tests: 0 (0.0%); Adapted tests: 2626 (100.0%)   | `packages/database.md`      |
| `Framework\Encryption`      | `encryption`    | `mapped`             | Ported tests: 26 / 27 (96.3%); Missing tests: 0 (0.0%); Adapted tests: 1 (3.7%)        | `packages/encryption.md`    |
| `Framework\Events`          | `events`        | `mapped`             | Ported tests: 61 / 76 (80.3%); Missing tests: 0 (0.0%); Adapted tests: 15 (19.7%)      | `packages/events.md`        |
| `Framework\Filesystem`      | `filesystem`    | `mapped`             | Ported tests: 54 / 134 (40.3%); Missing tests: 0 (0.0%); Adapted tests: 80 (59.7%)     | `packages/filesystem.md`    |
| `Framework\Foundation`      | n/a             | `excluded-permanent` | no inventory                                                                           | n/a                         |
| `Framework\Hashing`         | `hashing`       | `mapped`             | Ported tests: 12 / 13 (92.3%); Missing tests: 0 (0.0%); Adapted tests: 1 (7.7%)        | `packages/hashing.md`       |
| `Framework\Http`            | `httpx`         | `mapped`             | Ported tests: 488 / 488 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)     | `packages/httpx.md`         |
| `Framework\JsonSchema`      | `jsonx`         | `mapped`             | Ported tests: 47 / 47 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)       | `packages/jsonx.md`         |
| `Framework\Log`             | `log`           | `mapped`             | Ported tests: 98 / 98 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)       | `packages/log.md`           |
| `Framework\Mail`            | `mailx`         | `mapped`             | Ported tests: 135 / 135 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)     | `packages/mailx.md`         |
| `Framework\Notifications`   | `notifications` | `mapped`             | Ported tests: 71 / 71 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)       | `packages/notifications.md` |
| `Framework\Pagination`      | `pagination`    | `mapped`             | Ported tests: 48 / 48 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)       | `packages/pagination.md`    |
| `Framework\Pipeline`        | `pipeline`      | `mapped`             | Ported tests: 22 / 22 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)       | `packages/pipeline.md`      |
| `Framework\Process`         | `process`       | `mapped`             | Ported tests: 54 / 54 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)       | `packages/process.md`       |
| `Framework\Queue`           | `queue`         | `mapped`             | Ported tests: 195 / 204 (95.6%); Missing tests: 0 (0.0%); Adapted tests: 9 (4.4%)      | `packages/queue.md`         |
| `Framework\Redis`           | `redis`         | `mapped`             | Ported tests: 79 / 103 (76.7%); Missing tests: 0 (0.0%); Adapted tests: 24 (23.3%)     | `packages/redis.md`         |
| `Framework\Routing`         | `routing`       | `mapped`             | Ported tests: 142 / 315 (45.1%); Missing tests: 173 (54.9%); Adapted tests: 0 (0.0%)   | `packages/routing.md`       |
| `Framework\Session`         | `session`       | `mapped`             | Ported tests: 65 / 93 (69.9%); Missing tests: 0 (0.0%); Adapted tests: 28 (30.1%)      | `packages/session.md`       |
| `Framework\Support`         | `support`       | `mapped`             | Ported tests: 1006 / 1554 (64.7%); Missing tests: 0 (0.0%); Adapted tests: 548 (35.3%) | `packages/support.md`       |
| `Framework\Support\Str`     | `str`           | `mapped`             | Ported tests: 113 / 115 (98.3%); Missing tests: 0 (0.0%); Adapted tests: 2 (1.7%)      | `packages/str.md`           |
| `Framework\Support\Lottery` | `lottery`       | `mapped`             | Ported tests: 14 / 14 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)       | `packages/lottery.md`       |
| `Framework\Testing`         | `testing`       | `mapped`             | Ported tests: 359 / 381 (94.2%); Missing tests: 0 (0.0%); Adapted tests: 22 (5.8%)     | missing                     |
| `Framework\Translation`     | `translation`   | `mapped`             | Ported tests: 52 / 52 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)       | `packages/translation.md`   |
| `Framework\Validation`      | `validation`    | `mapped`             | Ported tests: 161 / 597 (27.0%); Missing tests: 411 (68.8%); Adapted tests: 25 (4.2%)  | `packages/validation.md`    |
| `Framework\Collections`     | `collection`    | `mapped`             | Ported tests: 489 / 629 (77.7%); Missing tests: 0 (0.0%); Adapted tests: 140 (22.3%)   | `packages/collection.md`    |
| `Framework\Macroable`       | n/a             | `excluded-permanent` | no inventory                                                                           | n/a                         |
| `Framework\Reflection`      | n/a             | `excluded-permanent` | no inventory                                                                           | n/a                         |
| `Framework\View`            | n/a             | `excluded-permanent` | no inventory                                                                           | n/a                         |

## Package Coverage

| Upstream Package               | Bedrock Surface | Port Status          | Tests                                                                                 | Docs                       |
| ----------------------------- | --------------- | -------------------- | ------------------------------------------------------------------------------------- | -------------------------- |
| `upstream/ai`                  | `ai/sdk`        | `mapped`             | Ported tests: 211 / 774 (27.3%); Missing tests: 551 (71.2%); Adapted tests: 12 (1.6%) | `packages/ai/sdk.md`       |
| `upstream/boost`               | `ai/boost`      | `mapped`             | Ported tests: 147 / 671 (21.9%); Missing tests: 518 (77.2%); Adapted tests: 6 (0.9%)  | `packages/ai/boost.md`     |
| `upstream/broadcastclient`                | `broadcastclient`          | `mapped`             | no inventory                                                                          | `packages/broadcastclient.md`         |
| `upstream/authflows`             | `inception`     | `mapped`             | Ported tests: 63 / 91 (69.2%); Missing tests: 0 (0.0%); Adapted tests: 28 (30.8%)     | `packages/inception.md`    |
| `upstream/authkit`           | `inception`     | `mapped`             | Ported tests: 28 / 43 (65.1%); Missing tests: 0 (0.0%); Adapted tests: 15 (34.9%)     | `packages/inception.md`    |
| `upstream/mcp`                 | `ai/mcp`        | `mapped`             | Ported tests: 201 / 624 (32.2%); Missing tests: 410 (65.7%); Adapted tests: 13 (2.1%) | `packages/ai/mcp.md`       |
| `upstream/oauthserver`            | `oauthserver`      | `mapped`             | Ported tests: 116 / 211 (55.0%); Missing tests: 76 (36.0%); Adapted tests: 19 (9.0%)  | `packages/oauthserver.md`     |
| `upstream/featureflags`             | `featureflags`       | `mapped`             | Ported tests: 129 / 186 (69.4%); Missing tests: 44 (23.7%); Adapted tests: 13 (7.0%)  | `packages/featureflags.md`      |
| `upstream/httppreview`        | `httppreview`  | `mapped`             | no inventory                                                                          | `packages/httppreview.md` |
| `upstream/prompts`             | `prompts`       | `mapped`             | Ported tests: 200 / 313 (63.9%); Missing tests: 49 (15.7%); Adapted tests: 64 (20.4%) | `packages/prompts.md`      |
| `upstream/websockets`              | `websockets`        | `mapped`             | Ported tests: 142 / 232 (61.2%); Missing tests: 90 (38.8%); Adapted tests: 0 (0.0%)   | `packages/websockets.md`       |
| `upstream/search`               | `search`         | `mapped`             | Ported tests: 65 / 252 (25.8%); Missing tests: 170 (67.5%); Adapted tests: 17 (6.7%)  | `packages/search.md`        |
| `upstream/socialauth`           | `socialauth`     | `mapped`             | Ported tests: 40 / 40 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)      | `packages/socialauth.md`    |
| `upstream/billing-paddle`        | `billing`         | `mapped`             | Ported tests: 92 / 166 (55.4%); Missing tests: 0 (0.0%); Adapted tests: 74 (44.6%)    | `packages/billing.md`        |
| `upstream/debugbar`           | `debugbar`     | `mapped`             | Ported tests: 89 / 113 (78.8%); Missing tests: 0 (0.0%); Adapted tests: 24 (21.2%)    | `packages/debugbar.md`    |
| `upstream/jobqueue`             | `jobqueue`       | `mapped`             | Ported tests: 97 / 180 (53.9%); Missing tests: 0 (0.0%); Adapted tests: 83 (46.1%)    | `packages/jobqueue.md`      |
| `upstream/logtail`                | `logtail`          | `mapped`             | Ported tests: 38 / 38 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)      | `packages/logtail.md`         |
| `upstream/remotetasks`               | `remotetasks`         | `mapped`             | Ported tests: 16 / 20 (80.0%); Missing tests: 0 (0.0%); Adapted tests: 4 (20.0%)      | `packages/remotetasks.md`        |
| `upstream/routegen`           | `routegen`     | `mapped`             | Ported tests: 103 / 103 (100.0%); Missing tests: 0 (0.0%); Adapted tests: 0 (0.0%)    | `packages/routegen.md`    |
| `upstream/tokens`             | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `upstream/octane`              | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `upstream/pulse`               | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `upstream/folio`               | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `upstream/dusk`                | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `upstream/browser-kit-testing` | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `upstream/sail`                | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `upstream/installer`           | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `upstream/valet`               | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `upstream/vite-plugin`         | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `starter-kits`                | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |
| `upstream/workos`              | n/a             | `excluded-permanent` | no inventory                                                                          | n/a                        |

## Test Porting Summary

| Scope                 | Inventories | Upstream Tests | Ported Tests | Pending / Missing Tests | Adapted Tests |
| --------------------- | ----------: | -------------: | -----------: | ----------------------: | ------------: |
| All inventories       |          51 |          13066 | 6269 (48.0%) |            2492 (19.1%) |  4305 (32.9%) |
| Framework inventories |          34 |           9009 | 4492 (49.9%) |              584 (6.5%) |  3933 (43.7%) |
| Package inventories   |          17 |           4057 | 1777 (43.8%) |            1908 (47.0%) |    372 (9.2%) |

## Documentation Porting Summary

| Scope        | Doc Inventories | Upstream Sections | Ported Sections | Pending / Missing Sections | Adapted Sections | Excluded Sections |
| ------------ | --------------: | ----------------: | --------------: | -------------------------: | ---------------: | ----------------: |
| Upstream docs |               1 |               875 |     200 (22.9%) |                505 (57.7%) |        10 (1.1%) |       160 (18.3%) |

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
| `billing.md`                |               14 |
| `prompts.md`                |               13 |
| `featureflags.md`                |               11 |
| `notifications.md`          |               10 |
| `migrations.md`             |               10 |
| `mcp.md`                    |               10 |
| `orm-relationships.md` |               10 |
| `contributions.md`          |               10 |
| `boost.md`                  |               10 |
| `mail.md`                   |                9 |
| `http-tests.md`             |                9 |
| `deployment.md`             |                9 |
| `cli.md`                |                9 |
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

| Feature                      | Source                      | Bedrock Surface | Status     | Docs                        | Notes                                                                                                                                                                                                                                                                                                                                         |
| ---------------------------- | --------------------------- | --------------- | ---------- | --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Authentication guards        | `framework.auth`            | `auth`          | `ported`   | `packages/auth.md`          | Session, token, request, basic-auth, provider, and event surfaces are complete under classified parity; PHP runtime mechanics are covered by Go adaptation rules in divergences.yml.                                                                                                                                                          |
| Password broker              | `framework.auth`            | `auth`          | `ported`   | `packages/auth.md`          | Password reset tokens, throttling, callback notification override, notification sender hook, and reset-link events are complete under classified parity; PHP manager/config mechanics are covered by Go adaptation rules in divergences.yml.                                                                                                  |
| Broadcasting contracts       | `framework.broadcasting`    | `broadcasting`  | `ported`   | `packages/broadcasting.md`  | Server-side channel auth, broadcast event payload handling, Pusher/Redis/Ably broadcaster responses, and Pusher channel conventions are ported; BroadcastClient remains the client receiver package.                                                                                                                                                     |
| Command bus dispatch         | `framework.bus`             | `bus`           | `ported`   | `packages/bus.md`           | Command dispatch, queue routing, after-response dispatch, chains, batches, lifecycle events, and typed failure callbacks are complete under classified parity; PHP runtime mechanics are covered by Go adaptation rules in divergences.yml.                                                                                                   |
| Cache repository             | `framework.cache`           | `cache`         | `partial`  | `packages/cache.md`         | Cache abstractions are present, but upstream test parity is still missing.                                                                                                                                                                                                                                                                    |
| Concurrent task execution    | `framework.concurrency`     | `concurrency`   | `ported`   | `packages/concurrency.md`   | Executable Go equivalents cover distributed task execution, error propagation, typed error inspection, and result ordering; PHP process and serialized-closure mechanics are adapted to explicit Go tasks, goroutines, error returns, and panic recovery.                                                                                     |
| Conditional fluent execution | `framework.conditionable`   | `conditionable` | `ported`   | `packages/conditionable.md` | Executable Go equivalents cover when and unless; PHP dynamic no-argument and magic behavior is represented through explicit Go proxy constructors.                                                                                                                                                                                            |
| Configuration repository     | `framework.config`          | `config`        | `ported`   | `packages/config.md`        | Dot-notation reads, literal-dot key precedence, explicit nil values, write helpers, typed getters, collection wrapping, and explicit Go offset adaptations are complete under classified parity; PHP Macroable runtime mechanics are covered by Go adaptation rules in divergences.yml.                                                       |
| Console commands             | `framework.console`         | `console`       | `ported`   | `packages/console.md`       | Command registration and dispatch, signature parsing, typed input and output helpers, signal traps, scheduling primitives, and explicit command/scheduling mutexes are complete under classified parity; Symfony/PHP attribute runtime and dynamic Macroable/callable-string mechanics are covered by Go adaptation rules in divergences.yml. |
| Container bindings           | `framework.container`       | `container`     | `ported`   | `packages/container.md`     | Bindings, shared and scoped lifecycles, aliases, contextual bindings, tags, extenders, method bindings, and resolving callbacks are complete under classified parity; PHP reflection, attributes, callable strings, ArrayAccess, and environment selectors are covered by Go adaptation rules in divergences.yml.                             |
| Shared interfaces            | `framework.contracts`       | `contracts`     | `ported`   | `packages/contracts.md`     | Interface compliance is tracked through concrete package inventories.                                                                                                                                                                                                                                                                         |
| Cookie handling              | `framework.cookie`          | `cookie`        | `ported`   | `packages/cookie.md`        | Cookie creation, queueing, response attachment, and encrypted-cookie middleware are complete under classified parity; PHP mixed arguments and Macroable runtime mechanics are covered by Go adaptation rules in divergences.yml.                                                                                                              |
| Database query builder       | `framework.database`        | `database`      | `ported`   | `packages/database.md`      | Query builders, SQL grammars, connection management, transactions, schema, migrations, and seeding are complete under classified parity; PHP facades, static model calls, exceptions, and unsupported SQL Server behavior are covered by Go adaptation rules in divergences.yml.                                                              |
| Orm-style ORM           | `framework.database`        | `database`      | `ported`   | `packages/database.md`      | Model, attribute, timestamp, scope, soft-delete, visibility, collection, and builder primitives are complete under classified parity; PHP trait, magic property, relation, and recursive runtime mechanics are covered by Go adaptation rules in divergences.yml.                                                                             |
| Encrypter                    | `framework.encryption`      | `encryption`    | `ported`   | `packages/encryption.md`    | CBC and GCM encryption, previous-key decryption, MAC/tag validation, key parsing, cipher validation, tamper rejection, and encrypted-value detection are complete under classified parity; PHP serialization and mixed runtime input checks are covered by Go adaptation rules in divergences.yml.                                            |
| Event dispatcher             | `framework.events`          | `events`        | `ported`   | `packages/events.md`        | Dispatch, defer, pushed events, wildcards, subscribers, queued listener options, and broadcast event integration are complete under classified parity; PHP container, interface, enum, named-argument, and unique-lock mechanics are covered by Go adaptation rules in divergences.yml.                                                       |
| Local filesystem             | `framework.filesystem`      | `filesystem`    | `ported`   | `packages/filesystem.md`    | Local file reads, writes, metadata, directory operations, hashing, globbing, and locking-oriented reads are complete under classified parity; Flysystem adapter/manager, PHP include, runtime macro, and helper mechanics are covered by Go adaptation rules in divergences.yml.                                                              |
| Application foundation       | `framework.foundation`      | n/a             | `excluded` | n/a                         | Full Upstream application kernel is outside Bedrock product boundaries.                                                                                                                                                                                                                                                                        |
| Password hashing             | `framework.hashing`         | `hashing`       | `ported`   | `packages/hashing.md`       | Bcrypt, Argon2i, Argon2id, verification, hash detection, long-password rejection, and algorithm mismatch behavior are complete under classified parity; PHP null input is covered by Go adaptation rules in divergences.yml.                                                                                                                  |
| HTTP primitives              | `framework.http`            | `httpx`         | `ported`   | `packages/httpx.md`         | Requests, responses, middleware, client, resources, uploaded files, and testing helpers are complete under classified parity; PHP runtime and framework integration mechanics are represented by explicit Go APIs.                                                                                                                            |
| JSON schema builder          | `framework.json_schema`     | `jsonx`         | `ported`   | `packages/jsonx.md`         | Fluent schema builders, serialization, type metadata, defaults, enum, required, and additional-property behavior are complete under classified parity.                                                                                                                                                                                        |
| Logging channels             | `framework.log`             | `log`           | `ported`   | `packages/log.md`           | Channel managers, stacks, handlers, levels, structured records, context repository, and formatting behavior are complete under classified parity.                                                                                                                                                                                             |
| Mail messages                | `framework.mail`            | `mailx`         | `ported`   | `packages/mailx.md`         | Mailables, messages, headers, attachments, mailers, lifecycle events, array/log/SMTP transports, failover, round-robin, and typed API-client transports are complete under classified parity.                                                                                                                                                 |
| Notification channels        | `framework.notifications`   | `notifications` | `ported`   | `packages/notifications.md` | Manager/channel resolution, sender flow, mail, database, broadcast, custom channels, queued dispatch, anonymous routes, and message helpers are complete under classified parity.                                                                                                                                                             |
| Pagination                   | `framework.pagination`      | `pagination`    | `ported`   | `packages/pagination.md`    | Offset, length-aware, and cursor paginators, cursor encoding, URL windows, and JSON/map serialization are complete under classified parity; Orm loadMorph helpers are adapted because Bedrock paginators operate on caller-owned Go slices.                                                                                              |
| Pipeline                     | `framework.pipeline`        | `pipeline`      | `ported`   | `packages/pipeline.md`      | Core pipe flow, string pipe parsing, via methods, finally callbacks, carry/error hooks, conditional pipes, hub execution, and transaction wrapping are complete under classified parity.                                                                                                                                                      |
| Process runner               | `framework.process`         | `process`       | `ported`   | `packages/process.md`       | Process execution, async invocation, pools, pipes, fakes, sequences, stray-process prevention, timeouts, input, output, and assertions are complete under classified parity.                                                                                                                                                                  |
| Queue workers and drivers    | `framework.queue`           | `queue`         | `ported`   | `packages/queue.md`         | Queue workers, drivers, payload hooks, routing, FIFO metadata, delay options, after-commit precedence, and selected Go adaptations are complete under classified parity.                                                                                                                                                                      |
| Redis commands               | `framework.redis`           | `redis`         | `partial`  | `packages/redis.md`         | Redis command surface exists; upstream Redis tests are not matched yet.                                                                                                                                                                                                                                                                       |
| Routing                      | `framework.routing`         | `routing`       | `partial`  | `packages/routing.md`       | Routing package exists; upstream routing tests remain incomplete.                                                                                                                                                                                                                                                                             |
| Session stores               | `framework.session`         | `session`       | `partial`  | `packages/session.md`       | Session stores exist; upstream session tests are not matched yet.                                                                                                                                                                                                                                                                             |
| Support helpers              | `framework.support`         | `support`       | `partial`  | `packages/support.md`       | Several support helpers are ported; upstream support coverage is still partial.                                                                                                                                                                                                                                                               |
| String helpers               | `framework.support_str`     | `str`           | `partial`  | `packages/str.md`           | String helper parity has many ported tests, but missing cases remain.                                                                                                                                                                                                                                                                         |
| Lottery helper               | `framework.support_lottery` | `lottery`       | `partial`  | `packages/lottery.md`       | Lottery package exists; upstream Lottery tests are not matched yet.                                                                                                                                                                                                                                                                           |
| Testing utilities            | `framework.testing`         | `testing`       | `partial`  | missing                     | Initial TestResponse-style assertions exist under packages/testing for HTTP status, headers, redirects, JSON, fluent JSON, cookies, session validation errors, streams, and downloads.                                                                                                                                                        |
| Translation loader           | `framework.translation`     | `translation`   | `partial`  | `packages/translation.md`   | Translation has some matched upstream tests; full parity remains incomplete.                                                                                                                                                                                                                                                                  |
| Validation rules             | `framework.validation`      | `validation`    | `partial`  | `packages/validation.md`    | Rule engine and many rules exist; upstream test parity is still broad and incomplete.                                                                                                                                                                                                                                                         |
| Collections                  | `framework.collections`     | `collection`    | `partial`  | `packages/collection.md`    | Fluent slice, key-value, lazy, arr, and kv packages are imported; legacy support helper cleanup is deferred.                                                                                                                                                                                                                                  |
| Macroable runtime            | `framework.macroable`       | n/a             | `excluded` | n/a                         | PHP runtime macro behaviour does not map cleanly to Go static method sets.                                                                                                                                                                                                                                                                    |
| Reflection helpers           | `framework.reflection`      | n/a             | `excluded` | n/a                         | PHP reflection helpers are not a Bedrock product surface.                                                                                                                                                                                                                                                                                     |
| Template views                  | `framework.view`            | n/a             | `excluded` | n/a                         | Bedrock does not ship Upstream Template or the view engine.                                                                                                                                                                                                                                                                                       |
| AI SDK                       | `package.ai`                | `ai/sdk`        | `partial`  | `packages/ai/sdk.md`        | Bedrock has an AI SDK package; upstream package parity is not complete.                                                                                                                                                                                                                                                                       |
| Agent tooling                | `package.boost`             | `ai/boost`      | `partial`  | `packages/ai/boost.md`      | Boost package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                             |
| BroadcastClient client                  | `package.broadcastclient`              | `broadcastclient`          | `partial`  | `packages/broadcastclient.md`          | Tracked by feature inventory because there is no stable upstream test inventory configured.                                                                                                                                                                                                                                                   |
| AuthFlows auth flows           | `package.authflows`           | `inception`     | `partial`  | `packages/authflows.md`       | Exposed through Inception with AuthFlows-specific documentation.                                                                                                                                                                                                                                                                                |
| AuthKit teams              | `package.authkit`         | `inception`     | `partial`  | `packages/authkit.md`     | Exposed through Inception with AuthKit-specific documentation.                                                                                                                                                                                                                                                                              |
| MCP server                   | `package.mcp`               | `ai/mcp`        | `partial`  | `packages/ai/mcp.md`        | MCP package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                               |
| OAuth server                 | `package.oauthserver`          | `oauthserver`      | `partial`  | `packages/oauthserver.md`      | OAuthServer package and docs exist; upstream package parity is not complete.                                                                                                                                                                                                                                                                     |
| Feature flags                | `package.featureflags`           | `featureflags`       | `partial`  | `packages/featureflags.md`       | FeatureFlags package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                           |
| HTTPPreview client          | `package.httppreview`      | `httppreview`  | `partial`  | `packages/httppreview.md`  | Tracked by feature inventory because there is no stable upstream test inventory configured.                                                                                                                                                                                                                                                   |
| CLI prompts                  | `package.prompts`           | `prompts`       | `partial`  | `packages/prompts.md`       | Prompt package exists with a small number of matched upstream tests.                                                                                                                                                                                                                                                                          |
| WebSockets websocket server      | `package.websockets`            | `websockets`        | `partial`  | `packages/websockets.md`        | WebSockets package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                            |
| Search indexing              | `package.search`             | `search`         | `partial`  | `packages/search.md`         | Search package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                             |
| Social OAuth                 | `package.socialauth`         | `socialauth`     | `partial`  | `packages/socialauth.md`     | SocialAuth package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                         |
| Subscription billing         | `package.billing`             | `billing`         | `partial`  | `packages/billing.md`         | Local Madora composer resolution of upstream/billing-paddle:^6.0 from https://billing.upstream.com at /Users/gocanto/Sites/madora/vendor/upstream/billing-paddle; upstream/cashier-paddle is its dependency, and the Billing tests in tests/Feature/Billing provide the local product integration coverage.                                             |
| DebugBar observability      | `package.debugbar`         | `debugbar`     | `partial`  | `packages/debugbar.md`     | DebugBar package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                         |
| RouteGen route helpers      | `package.routegen`         | `routegen`     | `partial`  | `packages/routegen.md`     | RouteGen package exists; upstream package parity remains incomplete.                                                                                                                                                                                                                                                                         |
| Queue monitoring             | `package.jobqueue`           | `jobqueue`       | `partial`  | `packages/jobqueue.md`       | JobQueue is tracked as queue monitoring primitives; the browser dashboard is not ported.                                                                                                                                                                                                                                                       |
| Log tailing                  | `package.logtail`              | `logtail`          | `partial`  | `packages/logtail.md`          | LogTail is tracked as log parsing and filtering primitives for Go readers.                                                                                                                                                                                                                                                                       |
| Remote task runner           | `package.remotetasks`             | `remotetasks`         | `partial`  | `packages/remotetasks.md`         | RemoteTasks is tracked as task planning and command execution primitives with injectable runners.                                                                                                                                                                                                                                                   |
| Tokens tokens               | `package.tokens`           | n/a             | `excluded` | n/a                         | Tokens is intentionally not part of Bedrock.                                                                                                                                                                                                                                                                                                 |
| Browser automation           | `package.dusk`              | n/a             | `excluded` | n/a                         | Browser automation package is excluded; browser work should use Agent Browser.                                                                                                                                                                                                                                                                |

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
