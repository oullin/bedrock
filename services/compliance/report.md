# Laravel Compliance Report

Generated: 2026-04-21T05:08:32Z

Source of truth: services/compliance

## Compliance Workflow Requirements

- Every mapped Bedrock package compliance pass must include an upstream feature audit: `make sure we also have all the upstream features for <package>`.
- The audit compares upstream source contracts, public APIs, middleware, events, and runtime behavior against the Bedrock surface.
- Bedrock-equivalent features need executable parity tests; PHP-only or intentionally different behavior belongs in `services/compliance/divergences.yml`.
- `services/compliance/features.yml` must contain feature audit coverage for every mapped Bedrock surface.

## Test Porting Summary

| Scope                 | Inventories | Upstream Tests | Ported Tests | Pending / Missing Tests | Adapted Tests |
| --------------------- | ----------: | -------------: | -----------: | ----------------------: | ------------: |
| All inventories       |          51 |          13012 |   535 (4.1%) |           12274 (94.3%) |    203 (1.6%) |
| Framework inventories |          34 |           8975 |   532 (5.9%) |            8240 (91.8%) |    203 (2.3%) |
| Package inventories   |          17 |           4037 |     3 (0.1%) |            4034 (99.9%) |      0 (0.0%) |

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

## Inventories

| Inventory                                 | Upstream Tests | Ported Tests | Missing Tests | Adapted Tests |
| ----------------------------------------- | -------------: | -----------: | ------------: | ------------: |
| inventories/framework-auth.txt            |            257 |   54 (21.0%) |      0 (0.0%) |   203 (79.0%) |
| inventories/framework-broadcasting.txt    |             55 |  55 (100.0%) |      0 (0.0%) |      0 (0.0%) |
| inventories/framework-bus.txt             |             66 |     0 (0.0%) |   66 (100.0%) |      0 (0.0%) |
| inventories/framework-cache.txt           |            313 |     0 (0.0%) |  313 (100.0%) |      0 (0.0%) |
| inventories/framework-concurrency.txt     |             12 |     0 (0.0%) |   12 (100.0%) |      0 (0.0%) |
| inventories/framework-conditionable.txt   |              2 |     0 (0.0%) |    2 (100.0%) |      0 (0.0%) |
| inventories/framework-config.txt          |             33 |     0 (0.0%) |   33 (100.0%) |      0 (0.0%) |
| inventories/framework-console.txt         |            156 |     0 (0.0%) |  156 (100.0%) |      0 (0.0%) |
| inventories/framework-container.txt       |            193 |     0 (0.0%) |  193 (100.0%) |      0 (0.0%) |
| inventories/framework-cookie.txt          |             26 |     0 (0.0%) |   26 (100.0%) |      0 (0.0%) |
| inventories/framework-database.txt        |           2626 |     0 (0.0%) | 2626 (100.0%) |      0 (0.0%) |
| inventories/framework-encryption.txt      |             27 |     0 (0.0%) |   27 (100.0%) |      0 (0.0%) |
| inventories/framework-events.txt          |             76 |     0 (0.0%) |   76 (100.0%) |      0 (0.0%) |
| inventories/framework-filesystem.txt      |            134 |     0 (0.0%) |  134 (100.0%) |      0 (0.0%) |
| inventories/framework-hashing.txt         |             13 |     0 (0.0%) |   13 (100.0%) |      0 (0.0%) |
| inventories/framework-http.txt            |            472 |     0 (0.0%) |  472 (100.0%) |      0 (0.0%) |
| inventories/framework-json-schema.txt     |             47 |     0 (0.0%) |   47 (100.0%) |      0 (0.0%) |
| inventories/framework-log.txt             |             85 |     0 (0.0%) |   85 (100.0%) |      0 (0.0%) |
| inventories/framework-mail.txt            |            130 |     0 (0.0%) |  130 (100.0%) |      0 (0.0%) |
| inventories/framework-notifications.txt   |             71 |     0 (0.0%) |   71 (100.0%) |      0 (0.0%) |
| inventories/framework-pagination.txt      |             48 |     0 (0.0%) |   48 (100.0%) |      0 (0.0%) |
| inventories/framework-pipeline.txt        |             22 |     0 (0.0%) |   22 (100.0%) |      0 (0.0%) |
| inventories/framework-process.txt         |             54 |     0 (0.0%) |   54 (100.0%) |      0 (0.0%) |
| inventories/framework-queue.txt           |            204 |  161 (78.9%) |    43 (21.1%) |      0 (0.0%) |
| inventories/framework-redis.txt           |            103 |     0 (0.0%) |  103 (100.0%) |      0 (0.0%) |
| inventories/framework-routing.txt         |            315 |     0 (0.0%) |  315 (100.0%) |      0 (0.0%) |
| inventories/framework-session.txt         |             93 |     0 (0.0%) |   93 (100.0%) |      0 (0.0%) |
| inventories/framework-support.txt         |           1554 |   129 (8.3%) |  1425 (91.7%) |      0 (0.0%) |
| inventories/framework-support-str.txt     |            115 |   69 (60.0%) |    46 (40.0%) |      0 (0.0%) |
| inventories/framework-support-lottery.txt |             14 |     0 (0.0%) |   14 (100.0%) |      0 (0.0%) |
| inventories/framework-testing.txt         |            381 |     0 (0.0%) |  381 (100.0%) |      0 (0.0%) |
| inventories/framework-translation.txt     |             52 |    9 (17.3%) |    43 (82.7%) |      0 (0.0%) |
| inventories/framework-validation.txt      |            597 |    34 (5.7%) |   563 (94.3%) |      0 (0.0%) |
| inventories/framework-collections.txt     |            629 |    21 (3.3%) |   608 (96.7%) |      0 (0.0%) |
| inventories/package-ai.txt                |            770 |     0 (0.0%) |  770 (100.0%) |      0 (0.0%) |
| inventories/package-boost.txt             |            671 |     0 (0.0%) |  671 (100.0%) |      0 (0.0%) |
| inventories/package-fortify.txt           |             91 |     0 (0.0%) |   91 (100.0%) |      0 (0.0%) |
| inventories/package-jetstream.txt         |             43 |     0 (0.0%) |   43 (100.0%) |      0 (0.0%) |
| inventories/package-mcp.txt               |            624 |     0 (0.0%) |  624 (100.0%) |      0 (0.0%) |
| inventories/package-passport.txt          |            211 |     0 (0.0%) |  211 (100.0%) |      0 (0.0%) |
| inventories/package-pennant.txt           |            186 |     0 (0.0%) |  186 (100.0%) |      0 (0.0%) |
| inventories/package-prompts.txt           |            308 |     3 (1.0%) |   305 (99.0%) |      0 (0.0%) |
| inventories/package-reverb.txt            |            232 |     0 (0.0%) |  232 (100.0%) |      0 (0.0%) |
| inventories/package-scout.txt             |            252 |     0 (0.0%) |  252 (100.0%) |      0 (0.0%) |
| inventories/package-socialite.txt         |             40 |     0 (0.0%) |   40 (100.0%) |      0 (0.0%) |
| inventories/package-spark.txt             |            166 |     0 (0.0%) |  166 (100.0%) |      0 (0.0%) |
| inventories/package-telescope.txt         |            113 |     0 (0.0%) |  113 (100.0%) |      0 (0.0%) |
| inventories/package-horizon.txt           |            169 |     0 (0.0%) |  169 (100.0%) |      0 (0.0%) |
| inventories/package-pail.txt              |             38 |     0 (0.0%) |   38 (100.0%) |      0 (0.0%) |
| inventories/package-envoy.txt             |             20 |     0 (0.0%) |   20 (100.0%) |      0 (0.0%) |
| inventories/package-wayfinder.txt         |            103 |     0 (0.0%) |  103 (100.0%) |      0 (0.0%) |

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
| `Illuminate\Http`            | `httpx`         | `mapped`             | Ported tests: 0 / 472 (0.0%); Missing tests: 472 (100.0%); Adapted tests: 0 (0.0%)    | `packages/httpx.md`         |
| `Illuminate\JsonSchema`      | `jsonx`         | `mapped`             | Ported tests: 0 / 47 (0.0%); Missing tests: 47 (100.0%); Adapted tests: 0 (0.0%)      | `packages/jsonx.md`         |
| `Illuminate\Log`             | `log`           | `mapped`             | Ported tests: 0 / 85 (0.0%); Missing tests: 85 (100.0%); Adapted tests: 0 (0.0%)      | `packages/log.md`           |
| `Illuminate\Mail`            | `mailx`         | `mapped`             | Ported tests: 0 / 130 (0.0%); Missing tests: 130 (100.0%); Adapted tests: 0 (0.0%)    | `packages/mailx.md`         |
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
