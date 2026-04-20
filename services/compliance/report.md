# Upstream Compliance Report

Generated: 2026-04-20T06:26:43Z

Source of truth: services/compliance

## Inventories

| Inventory                                 | Total | Ported | Adapted | Missing |
| ----------------------------------------- | ----: | -----: | ------: | ------: |
| inventories/framework-auth.txt            |   257 |      0 |       0 |     257 |
| inventories/framework-broadcasting.txt    |    55 |      0 |       0 |      55 |
| inventories/framework-bus.txt             |    66 |      0 |       0 |      66 |
| inventories/framework-cache.txt           |   313 |      0 |       0 |     313 |
| inventories/framework-concurrency.txt     |    12 |      0 |       0 |      12 |
| inventories/framework-conditionable.txt   |     2 |      0 |       0 |       2 |
| inventories/framework-config.txt          |    33 |      0 |       0 |      33 |
| inventories/framework-console.txt         |   156 |      0 |       0 |     156 |
| inventories/framework-container.txt       |   193 |      0 |       0 |     193 |
| inventories/framework-cookie.txt          |    26 |      0 |       0 |      26 |
| inventories/framework-database.txt        |  2626 |      0 |       0 |    2626 |
| inventories/framework-encryption.txt      |    27 |      0 |       0 |      27 |
| inventories/framework-events.txt          |    76 |      0 |       0 |      76 |
| inventories/framework-filesystem.txt      |   134 |      0 |       0 |     134 |
| inventories/framework-hashing.txt         |    13 |      0 |       0 |      13 |
| inventories/framework-http.txt            |   472 |      0 |       0 |     472 |
| inventories/framework-json-schema.txt     |    47 |      0 |       0 |      47 |
| inventories/framework-log.txt             |    85 |      0 |       0 |      85 |
| inventories/framework-mail.txt            |   130 |      0 |       0 |     130 |
| inventories/framework-notifications.txt   |    71 |      0 |       0 |      71 |
| inventories/framework-pagination.txt      |    48 |      0 |       0 |      48 |
| inventories/framework-pipeline.txt        |    22 |      0 |       0 |      22 |
| inventories/framework-process.txt         |    54 |      0 |       0 |      54 |
| inventories/framework-queue.txt           |   204 |    161 |       0 |      43 |
| inventories/framework-redis.txt           |   103 |      0 |       0 |     103 |
| inventories/framework-routing.txt         |   315 |      0 |       0 |     315 |
| inventories/framework-session.txt         |    93 |      0 |       0 |      93 |
| inventories/framework-support.txt         |  1554 |    129 |       0 |    1425 |
| inventories/framework-support-str.txt     |   115 |     69 |       0 |      46 |
| inventories/framework-support-lottery.txt |    14 |      0 |       0 |      14 |
| inventories/framework-testing.txt         |   381 |      0 |       0 |     381 |
| inventories/framework-translation.txt     |    52 |      9 |       0 |      43 |
| inventories/framework-validation.txt      |   597 |     34 |       0 |     563 |
| inventories/framework-collections.txt     |   629 |     21 |       0 |     608 |
| inventories/package-ai.txt                |   770 |      0 |       0 |     770 |
| inventories/package-boost.txt             |   671 |      0 |       0 |     671 |
| inventories/package-authflows.txt           |    91 |      0 |       0 |      91 |
| inventories/package-authkit.txt         |    43 |      0 |       0 |      43 |
| inventories/package-mcp.txt               |   624 |      0 |       0 |     624 |
| inventories/package-oauthserver.txt          |   211 |      0 |       0 |     211 |
| inventories/package-featureflags.txt           |   186 |      0 |       0 |     186 |
| inventories/package-prompts.txt           |   308 |      3 |       0 |     305 |
| inventories/package-websockets.txt            |   232 |      0 |       0 |     232 |
| inventories/package-search.txt             |   252 |      0 |       0 |     252 |
| inventories/package-socialauth.txt         |    40 |      0 |       0 |      40 |
| inventories/package-billing.txt             |   160 |      0 |       0 |     160 |
| inventories/package-debugbar.txt         |   113 |      0 |       0 |     113 |
| inventories/package-routegen.txt         |   103 |      0 |       0 |     103 |

## Framework Coverage

| Upstream Source               | Bedrock Surface | Port Status          | Tests                                    | Docs                        |
| ---------------------------- | --------------- | -------------------- | ---------------------------------------- | --------------------------- |
| `Framework\Auth`            | `auth`          | `mapped`             | 0/257 ported, 0 adapted, 257 missing     | `packages/auth.md`          |
| `Framework\Broadcasting`    | `broadcastclient`          | `missing`            | 0/55 ported, 0 adapted, 55 missing       | `packages/broadcastclient.md`          |
| `Framework\Bus`             | `bus`           | `mapped`             | 0/66 ported, 0 adapted, 66 missing       | `packages/bus.md`           |
| `Framework\Cache`           | `cache`         | `mapped`             | 0/313 ported, 0 adapted, 313 missing     | `packages/cache.md`         |
| `Framework\Concurrency`     | `concurrency`   | `mapped`             | 0/12 ported, 0 adapted, 12 missing       | `packages/concurrency.md`   |
| `Framework\Conditionable`   | `conditionable` | `mapped`             | 0/2 ported, 0 adapted, 2 missing         | `packages/conditionable.md` |
| `Framework\Config`          | `config`        | `mapped`             | 0/33 ported, 0 adapted, 33 missing       | `packages/config.md`        |
| `Framework\Console`         | n/a             | `missing`            | 0/156 ported, 0 adapted, 156 missing     | n/a                         |
| `Framework\Container`       | `container`     | `mapped`             | 0/193 ported, 0 adapted, 193 missing     | `packages/container.md`     |
| `Framework\Contracts`       | `contracts`     | `mapped`             | no inventory                             | `packages/contracts.md`     |
| `Framework\Cookie`          | `cookie`        | `mapped`             | 0/26 ported, 0 adapted, 26 missing       | `packages/cookie.md`        |
| `Framework\Database`        | `database`      | `mapped`             | 0/2626 ported, 0 adapted, 2626 missing   | `packages/database.md`      |
| `Framework\Encryption`      | `encryption`    | `mapped`             | 0/27 ported, 0 adapted, 27 missing       | `packages/encryption.md`    |
| `Framework\Events`          | `events`        | `mapped`             | 0/76 ported, 0 adapted, 76 missing       | `packages/events.md`        |
| `Framework\Filesystem`      | `filesystem`    | `mapped`             | 0/134 ported, 0 adapted, 134 missing     | `packages/filesystem.md`    |
| `Framework\Foundation`      | n/a             | `excluded-permanent` | no inventory                             | n/a                         |
| `Framework\Hashing`         | `hashing`       | `mapped`             | 0/13 ported, 0 adapted, 13 missing       | `packages/hashing.md`       |
| `Framework\Http`            | `httpx`         | `mapped`             | 0/472 ported, 0 adapted, 472 missing     | `packages/httpx.md`         |
| `Framework\JsonSchema`      | `jsonx`         | `mapped`             | 0/47 ported, 0 adapted, 47 missing       | `packages/jsonx.md`         |
| `Framework\Log`             | `log`           | `mapped`             | 0/85 ported, 0 adapted, 85 missing       | `packages/log.md`           |
| `Framework\Mail`            | `mailx`         | `mapped`             | 0/130 ported, 0 adapted, 130 missing     | `packages/mailx.md`         |
| `Framework\Notifications`   | `notifications` | `mapped`             | 0/71 ported, 0 adapted, 71 missing       | `packages/notifications.md` |
| `Framework\Pagination`      | `pagination`    | `mapped`             | 0/48 ported, 0 adapted, 48 missing       | `packages/pagination.md`    |
| `Framework\Pipeline`        | `pipeline`      | `mapped`             | 0/22 ported, 0 adapted, 22 missing       | `packages/pipeline.md`      |
| `Framework\Process`         | n/a             | `missing`            | 0/54 ported, 0 adapted, 54 missing       | n/a                         |
| `Framework\Queue`           | `queue`         | `mapped`             | 161/204 ported, 0 adapted, 43 missing    | `packages/queue.md`         |
| `Framework\Redis`           | `redis`         | `mapped`             | 0/103 ported, 0 adapted, 103 missing     | `packages/redis.md`         |
| `Framework\Routing`         | `routing`       | `mapped`             | 0/315 ported, 0 adapted, 315 missing     | `packages/routing.md`       |
| `Framework\Session`         | `session`       | `mapped`             | 0/93 ported, 0 adapted, 93 missing       | `packages/session.md`       |
| `Framework\Support`         | `support`       | `mapped`             | 129/1554 ported, 0 adapted, 1425 missing | `packages/support.md`       |
| `Framework\Support\Str`     | `str`           | `mapped`             | 69/115 ported, 0 adapted, 46 missing     | `packages/str.md`           |
| `Framework\Support\Lottery` | `lottery`       | `mapped`             | 0/14 ported, 0 adapted, 14 missing       | `packages/lottery.md`       |
| `Framework\Testing`         | n/a             | `missing`            | 0/381 ported, 0 adapted, 381 missing     | n/a                         |
| `Framework\Translation`     | `translation`   | `mapped`             | 9/52 ported, 0 adapted, 43 missing       | `packages/translation.md`   |
| `Framework\Validation`      | `validation`    | `mapped`             | 34/597 ported, 0 adapted, 563 missing    | `packages/validation.md`    |
| `Framework\Collections`     | `collection`    | `mapped`             | 21/629 ported, 0 adapted, 608 missing    | `packages/collection.md`    |
| `Framework\Macroable`       | n/a             | `excluded-permanent` | no inventory                             | n/a                         |
| `Framework\Reflection`      | n/a             | `excluded-permanent` | no inventory                             | n/a                         |
| `Framework\View`            | n/a             | `excluded-permanent` | no inventory                             | n/a                         |

## Feature Coverage

| Feature                      | Source                      | Bedrock Surface | Status     | Docs                        | Notes                                                                                                        |
| ---------------------------- | --------------------------- | --------------- | ---------- | --------------------------- | ------------------------------------------------------------------------------------------------------------ |
| Authentication guards        | `framework.auth`            | `auth`          | `partial`  | `packages/auth.md`          | Session, token, and request guards exist; upstream test inventory is still mostly missing.                   |
| Password broker              | `framework.auth`            | `auth`          | `partial`  | `packages/auth.md`          | Password management is ported through auth; upstream test parity remains incomplete.                         |
| Broadcasting contracts       | `framework.broadcasting`    | `broadcastclient`          | `missing`  | `packages/broadcastclient.md`          | BroadcastClient and WebSockets pieces exist, but no framework broadcasting contract inventory is ported yet.                |
| Command bus dispatch         | `framework.bus`             | `bus`           | `partial`  | `packages/bus.md`           | Bus package exists with command dispatch and middleware; upstream test parity is not complete.               |
| Cache repository             | `framework.cache`           | `cache`         | `partial`  | `packages/cache.md`         | Cache abstractions are present, but upstream test parity is still missing.                                   |
| Concurrent task execution    | `framework.concurrency`     | `concurrency`   | `partial`  | `packages/concurrency.md`   | Bedrock has a concurrency package; upstream concurrency tests have not been matched yet.                     |
| Conditional fluent execution | `framework.conditionable`   | `conditionable` | `partial`  | `packages/conditionable.md` | Go fluent proxy exists; PHP dynamic behaviour is adapted to Go.                                              |
| Configuration repository     | `framework.config`          | `config`        | `partial`  | `packages/config.md`        | Config package exists; upstream test parity is not complete.                                                 |
| Console commands             | `framework.console`         | n/a             | `missing`  | n/a                         | No Bedrock console package exists yet.                                                                       |
| Container bindings           | `framework.container`       | `container`     | `partial`  | `packages/container.md`     | Container package exists; upstream container tests have not been matched yet.                                |
| Shared interfaces            | `framework.contracts`       | `contracts`     | `ported`   | `packages/contracts.md`     | Interface compliance is tracked through concrete package inventories.                                        |
| Cookie handling              | `framework.cookie`          | `cookie`        | `partial`  | `packages/cookie.md`        | Cookie package exists; upstream test parity is not complete.                                                 |
| Database query builder       | `framework.database`        | `database`      | `partial`  | `packages/database.md`      | Query, schema, migrations, and Orm surfaces exist; broad upstream parity remains incomplete.            |
| Orm-style ORM           | `framework.database`        | `database`      | `partial`  | `packages/database.md`      | ORM primitives exist with Go adaptations for attributes, relations, and errors.                              |
| Encrypter                    | `framework.encryption`      | `encryption`    | `partial`  | `packages/encryption.md`    | CBC and GCM encryption exist; upstream encryption tests are not fully matched.                               |
| Event dispatcher             | `framework.events`          | `events`        | `partial`  | `packages/events.md`        | Event dispatching exists; upstream test parity remains incomplete.                                           |
| Local filesystem             | `framework.filesystem`      | `filesystem`    | `partial`  | `packages/filesystem.md`    | Local filesystem operations are implemented; upstream parity tests are not complete.                         |
| Application foundation       | `framework.foundation`      | n/a             | `excluded` | n/a                         | Full Upstream application kernel is outside Bedrock product boundaries.                                       |
| Password hashing             | `framework.hashing`         | `hashing`       | `partial`  | `packages/hashing.md`       | Hashing drivers exist; upstream test parity is not complete.                                                 |
| HTTP primitives              | `framework.http`            | `httpx`         | `partial`  | `packages/httpx.md`         | HTTP primitives and test helpers exist; upstream test parity remains incomplete.                             |
| JSON schema builder          | `framework.json_schema`     | `jsonx`         | `partial`  | `packages/jsonx.md`         | Fluent JSON schema builder exists; upstream test parity is not complete.                                     |
| Logging channels             | `framework.log`             | `log`           | `partial`  | `packages/log.md`           | Driver-based structured logging exists; upstream log tests are not matched yet.                              |
| Mail messages                | `framework.mail`            | `mailx`         | `partial`  | `packages/mailx.md`         | Mail package exists; upstream mail test parity is incomplete.                                                |
| Notification channels        | `framework.notifications`   | `notifications` | `partial`  | `packages/notifications.md` | Multi-channel notifications exist; upstream parity remains incomplete.                                       |
| Pagination                   | `framework.pagination`      | `pagination`    | `partial`  | `packages/pagination.md`    | Offset and cursor paginators exist; upstream parity is incomplete.                                           |
| Pipeline                     | `framework.pipeline`        | `pipeline`      | `partial`  | `packages/pipeline.md`      | Pipeline package exists; upstream test parity is not complete.                                               |
| Process runner               | `framework.process`         | n/a             | `missing`  | n/a                         | No Bedrock process package exists yet.                                                                       |
| Queue workers and drivers    | `framework.queue`           | `queue`         | `partial`  | `packages/queue.md`         | Queue has the strongest upstream test match so far, but missing tests remain.                                |
| Redis commands               | `framework.redis`           | `redis`         | `partial`  | `packages/redis.md`         | Redis command surface exists; upstream Redis tests are not matched yet.                                      |
| Routing                      | `framework.routing`         | `routing`       | `partial`  | `packages/routing.md`       | Routing package exists; upstream routing tests remain incomplete.                                            |
| Session stores               | `framework.session`         | `session`       | `partial`  | `packages/session.md`       | Session stores exist; upstream session tests are not matched yet.                                            |
| Support helpers              | `framework.support`         | `support`       | `partial`  | `packages/support.md`       | Several support helpers are ported; upstream support coverage is still partial.                              |
| String helpers               | `framework.support_str`     | `str`           | `partial`  | `packages/str.md`           | String helper parity has many ported tests, but missing cases remain.                                        |
| Lottery helper               | `framework.support_lottery` | `lottery`       | `partial`  | `packages/lottery.md`       | Lottery package exists; upstream Lottery tests are not matched yet.                                          |
| Testing utilities            | `framework.testing`         | n/a             | `missing`  | n/a                         | No standalone Bedrock testing package exists yet.                                                            |
| Translation loader           | `framework.translation`     | `translation`   | `partial`  | `packages/translation.md`   | Translation has some matched upstream tests; full parity remains incomplete.                                 |
| Validation rules             | `framework.validation`      | `validation`    | `partial`  | `packages/validation.md`    | Rule engine and many rules exist; upstream test parity is still broad and incomplete.                        |
| Collections                  | `framework.collections`     | `collection`    | `partial`  | `packages/collection.md`    | Fluent slice, key-value, lazy, arr, and kv packages are imported; legacy support helper cleanup is deferred. |
| Macroable runtime            | `framework.macroable`       | n/a             | `excluded` | n/a                         | PHP runtime macro behaviour does not map cleanly to Go static method sets.                                   |
| Reflection helpers           | `framework.reflection`      | n/a             | `excluded` | n/a                         | PHP reflection helpers are not a Bedrock product surface.                                                    |
| Template views                  | `framework.view`            | n/a             | `excluded` | n/a                         | Bedrock does not ship Upstream Template or the view engine.                                                      |
| AI SDK                       | `package.ai`                | `ai/sdk`        | `partial`  | `packages/ai/sdk.md`        | Bedrock has an AI SDK package; upstream package parity is not complete.                                      |
| Agent tooling                | `package.boost`             | `ai/boost`      | `partial`  | `packages/ai/boost.md`      | Boost package exists; upstream package parity remains incomplete.                                            |
| BroadcastClient client                  | `package.broadcastclient`              | `broadcastclient`          | `partial`  | `packages/broadcastclient.md`          | Tracked by feature inventory because there is no stable upstream test inventory configured.                  |
| AuthFlows auth flows           | `package.authflows`           | `inception`     | `partial`  | `packages/authflows.md`       | Exposed through Inception with AuthFlows-specific documentation.                                               |
| AuthKit teams              | `package.authkit`         | `inception`     | `partial`  | `packages/authkit.md`     | Exposed through Inception with AuthKit-specific documentation.                                             |
| MCP server                   | `package.mcp`               | `ai/mcp`        | `partial`  | `packages/ai/mcp.md`        | MCP package exists; upstream package parity remains incomplete.                                              |
| OAuth server                 | `package.oauthserver`          | `oauthserver`      | `partial`  | `packages/oauthserver.md`      | OAuthServer package and docs exist; upstream package parity is not complete.                                    |
| Feature flags                | `package.featureflags`           | `featureflags`       | `partial`  | `packages/featureflags.md`       | FeatureFlags package exists; upstream package parity remains incomplete.                                          |
| HTTPPreview client          | `package.httppreview`      | `httppreview`  | `partial`  | `packages/httppreview.md`  | Tracked by feature inventory because there is no stable upstream test inventory configured.                  |
| CLI prompts                  | `package.prompts`           | `prompts`       | `partial`  | `packages/prompts.md`       | Prompt package exists with a small number of matched upstream tests.                                         |
| WebSockets websocket server      | `package.websockets`            | `websockets`        | `partial`  | `packages/websockets.md`        | WebSockets package exists; upstream package parity remains incomplete.                                           |
| Search indexing              | `package.search`             | `search`         | `partial`  | `packages/search.md`         | Search package exists; upstream package parity remains incomplete.                                            |
| Social OAuth                 | `package.socialauth`         | `socialauth`     | `partial`  | `packages/socialauth.md`     | SocialAuth package exists; upstream package parity remains incomplete.                                        |
| Subscription billing         | `package.billing`             | `billing`         | `partial`  | `packages/billing.md`         | Billing is tracked against local Madora billing tests.                                                         |
| DebugBar observability      | `package.debugbar`         | `debugbar`     | `partial`  | `packages/debugbar.md`     | DebugBar package exists; upstream package parity remains incomplete.                                        |
| RouteGen route helpers      | `package.routegen`         | `routegen`     | `partial`  | `packages/routegen.md`     | RouteGen package exists; upstream package parity remains incomplete.                                        |
| Tokens tokens               | `package.tokens`           | n/a             | `excluded` | n/a                         | Tokens is intentionally not part of Bedrock.                                                                |
| JobQueue dashboard            | `package.jobqueue`           | n/a             | `excluded` | n/a                         | Queue monitoring UI and runtime are outside Bedrock product boundaries.                                      |
| Browser automation           | `package.dusk`              | n/a             | `excluded` | n/a                         | Browser automation package is excluded; browser work should use Agent Browser.                               |

## Mapped Sources Without Inventories

| Source                 | Bedrock        | Reason                                           |
| ---------------------- | -------------- | ------------------------------------------------ |
| `Framework\Contracts` | `contracts`    | No generated upstream test inventory configured. |
| `upstream/broadcastclient`         | `broadcastclient`         | No generated upstream test inventory configured. |
| `upstream/httppreview` | `httppreview` | No generated upstream test inventory configured. |

## Permanent Exclusions

- `Framework\Foundation` - Product boundary; Bedrock does not ship Upstream's full application kernel/Foundation layer.
- `Framework\Macroable` - PHP runtime macro behaviour does not map cleanly to Go's static method set.
- `Framework\Reflection` - PHP reflection helpers are not a Bedrock product surface.
- `Framework\View` - Product boundary; Bedrock does not ship Upstream's Template/view engine.
- `upstream/tokens` - Product boundary; Tokens is intentionally not part of Bedrock.
- `upstream/jobqueue` - Product boundary; queue monitoring UI/runtime is intentionally excluded.
- `upstream/octane` - Product boundary; PHP worker/runtime management does not map to Bedrock.
- `upstream/pulse` - Product boundary; application monitoring dashboard is intentionally excluded.
- `upstream/logtail` - Product boundary; log-tail CLI is intentionally excluded.
- `upstream/folio` - Product boundary; page-based PHP routing is intentionally excluded.
- `upstream/dusk` - Product boundary; browser automation package is excluded. Any future browser work must use Agent Browser.
- `upstream/browser-kit-testing` - Product boundary; legacy BrowserKit compatibility is intentionally excluded.
- `upstream/sail` - Product boundary; Docker development environment tooling is intentionally excluded.
- `upstream/installer` - Product boundary; Upstream application installer tooling is intentionally excluded.
- `upstream/valet` - Product boundary; macOS local development service management is intentionally excluded.
- `upstream/remotetasks` - Product boundary; SSH task runner is intentionally excluded.
- `upstream/vite-plugin` - Product boundary; Vite integration tooling is intentionally excluded.
- `starter-kits` - Product boundary; application starter kits are intentionally excluded.
- `upstream/workos` - Product boundary; WorkOS starter-kit integration is intentionally excluded.
