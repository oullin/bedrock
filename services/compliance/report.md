# Laravel Compliance Report

Generated: 2026-04-20T05:02:48Z

Source of truth: services/compliance

## Inventories

| Inventory | Total | Ported | Adapted | Missing |
| --- | ---: | ---: | ---: | ---: |
| inventories/framework-auth.txt | 257 | 0 | 0 | 257 |
| inventories/framework-broadcasting.txt | 55 | 0 | 0 | 55 |
| inventories/framework-bus.txt | 66 | 0 | 0 | 66 |
| inventories/framework-cache.txt | 313 | 0 | 0 | 313 |
| inventories/framework-concurrency.txt | 12 | 0 | 0 | 12 |
| inventories/framework-conditionable.txt | 2 | 0 | 0 | 2 |
| inventories/framework-config.txt | 33 | 0 | 0 | 33 |
| inventories/framework-console.txt | 156 | 0 | 0 | 156 |
| inventories/framework-container.txt | 193 | 0 | 0 | 193 |
| inventories/framework-cookie.txt | 26 | 0 | 0 | 26 |
| inventories/framework-database.txt | 2626 | 0 | 0 | 2626 |
| inventories/framework-encryption.txt | 27 | 0 | 0 | 27 |
| inventories/framework-events.txt | 76 | 0 | 0 | 76 |
| inventories/framework-filesystem.txt | 134 | 0 | 0 | 134 |
| inventories/framework-hashing.txt | 13 | 0 | 0 | 13 |
| inventories/framework-http.txt | 472 | 0 | 0 | 472 |
| inventories/framework-json-schema.txt | 47 | 0 | 0 | 47 |
| inventories/framework-log.txt | 85 | 0 | 0 | 85 |
| inventories/framework-mail.txt | 130 | 0 | 0 | 130 |
| inventories/framework-notifications.txt | 71 | 0 | 0 | 71 |
| inventories/framework-pagination.txt | 48 | 0 | 0 | 48 |
| inventories/framework-pipeline.txt | 22 | 0 | 0 | 22 |
| inventories/framework-process.txt | 54 | 0 | 0 | 54 |
| inventories/framework-queue.txt | 204 | 161 | 0 | 43 |
| inventories/framework-redis.txt | 103 | 0 | 0 | 103 |
| inventories/framework-routing.txt | 315 | 0 | 0 | 315 |
| inventories/framework-session.txt | 93 | 0 | 0 | 93 |
| inventories/framework-support.txt | 1554 | 129 | 0 | 1425 |
| inventories/framework-support-str.txt | 115 | 69 | 0 | 46 |
| inventories/framework-support-lottery.txt | 14 | 0 | 0 | 14 |
| inventories/framework-testing.txt | 381 | 0 | 0 | 381 |
| inventories/framework-translation.txt | 52 | 9 | 0 | 43 |
| inventories/framework-validation.txt | 597 | 34 | 0 | 563 |
| inventories/package-ai.txt | 770 | 0 | 0 | 770 |
| inventories/package-boost.txt | 671 | 0 | 0 | 671 |
| inventories/package-fortify.txt | 91 | 0 | 0 | 91 |
| inventories/package-jetstream.txt | 43 | 0 | 0 | 43 |
| inventories/package-mcp.txt | 624 | 0 | 0 | 624 |
| inventories/package-passport.txt | 211 | 0 | 0 | 211 |
| inventories/package-pennant.txt | 186 | 0 | 0 | 186 |
| inventories/package-prompts.txt | 308 | 3 | 0 | 305 |
| inventories/package-reverb.txt | 232 | 0 | 0 | 232 |
| inventories/package-scout.txt | 252 | 0 | 0 | 252 |
| inventories/package-socialite.txt | 40 | 0 | 0 | 40 |
| inventories/package-spark.txt | 160 | 0 | 0 | 160 |
| inventories/package-telescope.txt | 113 | 0 | 0 | 113 |
| inventories/package-wayfinder.txt | 103 | 0 | 0 | 103 |

## Framework Coverage

| Laravel Source | Bedrock Surface | Port Status | Tests | Docs |
| --- | --- | --- | --- | --- |
| `Illuminate\Auth` | `auth` | `mapped` | 0/257 ported, 0 adapted, 257 missing | `packages/auth.md` |
| `Illuminate\Broadcasting` | `echo` | `missing` | 0/55 ported, 0 adapted, 55 missing | `packages/echo.md` |
| `Illuminate\Bus` | `bus` | `mapped` | 0/66 ported, 0 adapted, 66 missing | `packages/bus.md` |
| `Illuminate\Cache` | `cache` | `mapped` | 0/313 ported, 0 adapted, 313 missing | `packages/cache.md` |
| `Illuminate\Concurrency` | `concurrency` | `mapped` | 0/12 ported, 0 adapted, 12 missing | `packages/concurrency.md` |
| `Illuminate\Conditionable` | `conditionable` | `mapped` | 0/2 ported, 0 adapted, 2 missing | `packages/conditionable.md` |
| `Illuminate\Config` | `config` | `mapped` | 0/33 ported, 0 adapted, 33 missing | `packages/config.md` |
| `Illuminate\Console` | n/a | `missing` | 0/156 ported, 0 adapted, 156 missing | n/a |
| `Illuminate\Container` | `container` | `mapped` | 0/193 ported, 0 adapted, 193 missing | `packages/container.md` |
| `Illuminate\Contracts` | `contracts` | `mapped` | no inventory | `packages/contracts.md` |
| `Illuminate\Cookie` | `cookie` | `mapped` | 0/26 ported, 0 adapted, 26 missing | `packages/cookie.md` |
| `Illuminate\Database` | `database` | `mapped` | 0/2626 ported, 0 adapted, 2626 missing | `packages/database.md` |
| `Illuminate\Encryption` | `encryption` | `mapped` | 0/27 ported, 0 adapted, 27 missing | `packages/encryption.md` |
| `Illuminate\Events` | `events` | `mapped` | 0/76 ported, 0 adapted, 76 missing | `packages/events.md` |
| `Illuminate\Filesystem` | `filesystem` | `mapped` | 0/134 ported, 0 adapted, 134 missing | `packages/filesystem.md` |
| `Illuminate\Foundation` | n/a | `excluded-permanent` | no inventory | n/a |
| `Illuminate\Hashing` | `hashing` | `mapped` | 0/13 ported, 0 adapted, 13 missing | `packages/hashing.md` |
| `Illuminate\Http` | `httpx` | `mapped` | 0/472 ported, 0 adapted, 472 missing | `packages/httpx.md` |
| `Illuminate\JsonSchema` | `jsonx` | `mapped` | 0/47 ported, 0 adapted, 47 missing | `packages/jsonx.md` |
| `Illuminate\Log` | `log` | `mapped` | 0/85 ported, 0 adapted, 85 missing | `packages/log.md` |
| `Illuminate\Mail` | `mailx` | `mapped` | 0/130 ported, 0 adapted, 130 missing | `packages/mailx.md` |
| `Illuminate\Notifications` | `notifications` | `mapped` | 0/71 ported, 0 adapted, 71 missing | `packages/notifications.md` |
| `Illuminate\Pagination` | `pagination` | `mapped` | 0/48 ported, 0 adapted, 48 missing | `packages/pagination.md` |
| `Illuminate\Pipeline` | `pipeline` | `mapped` | 0/22 ported, 0 adapted, 22 missing | `packages/pipeline.md` |
| `Illuminate\Process` | n/a | `missing` | 0/54 ported, 0 adapted, 54 missing | n/a |
| `Illuminate\Queue` | `queue` | `mapped` | 161/204 ported, 0 adapted, 43 missing | `packages/queue.md` |
| `Illuminate\Redis` | `redis` | `mapped` | 0/103 ported, 0 adapted, 103 missing | `packages/redis.md` |
| `Illuminate\Routing` | `routing` | `mapped` | 0/315 ported, 0 adapted, 315 missing | `packages/routing.md` |
| `Illuminate\Session` | `session` | `mapped` | 0/93 ported, 0 adapted, 93 missing | `packages/session.md` |
| `Illuminate\Support` | `support` | `mapped` | 129/1554 ported, 0 adapted, 1425 missing | `packages/support.md` |
| `Illuminate\Support\Str` | `str` | `mapped` | 69/115 ported, 0 adapted, 46 missing | `packages/str.md` |
| `Illuminate\Support\Lottery` | `lottery` | `mapped` | 0/14 ported, 0 adapted, 14 missing | `packages/lottery.md` |
| `Illuminate\Testing` | n/a | `missing` | 0/381 ported, 0 adapted, 381 missing | n/a |
| `Illuminate\Translation` | `translation` | `mapped` | 9/52 ported, 0 adapted, 43 missing | `packages/translation.md` |
| `Illuminate\Validation` | `validation` | `mapped` | 34/597 ported, 0 adapted, 563 missing | `packages/validation.md` |
| `Illuminate\Collections` | n/a | `excluded-permanent` | no inventory | n/a |
| `Illuminate\Macroable` | n/a | `excluded-permanent` | no inventory | n/a |
| `Illuminate\Reflection` | n/a | `excluded-permanent` | no inventory | n/a |
| `Illuminate\View` | n/a | `excluded-permanent` | no inventory | n/a |

## Feature Coverage

| Feature | Source | Bedrock Surface | Status | Docs | Notes |
| --- | --- | --- | --- | --- | --- |
| Authentication guards | `framework.auth` | `auth` | `partial` | `packages/auth.md` | Session, token, and request guards exist; upstream test inventory is still mostly missing. |
| Password broker | `framework.auth` | `auth` | `partial` | `packages/auth.md` | Password management is ported through auth; upstream test parity remains incomplete. |
| Broadcasting contracts | `framework.broadcasting` | `echo` | `missing` | `packages/echo.md` | Echo and Reverb pieces exist, but no framework broadcasting contract inventory is ported yet. |
| Command bus dispatch | `framework.bus` | `bus` | `partial` | `packages/bus.md` | Bus package exists with command dispatch and middleware; upstream test parity is not complete. |
| Cache repository | `framework.cache` | `cache` | `partial` | `packages/cache.md` | Cache abstractions are present, but upstream test parity is still missing. |
| Concurrent task execution | `framework.concurrency` | `concurrency` | `partial` | `packages/concurrency.md` | Bedrock has a concurrency package; upstream concurrency tests have not been matched yet. |
| Conditional fluent execution | `framework.conditionable` | `conditionable` | `partial` | `packages/conditionable.md` | Go fluent proxy exists; PHP dynamic behaviour is adapted to Go. |
| Configuration repository | `framework.config` | `config` | `partial` | `packages/config.md` | Config package exists; upstream test parity is not complete. |
| Console commands | `framework.console` | n/a | `missing` | n/a | No Bedrock console package exists yet. |
| Container bindings | `framework.container` | `container` | `partial` | `packages/container.md` | Container package exists; upstream container tests have not been matched yet. |
| Shared interfaces | `framework.contracts` | `contracts` | `ported` | `packages/contracts.md` | Interface compliance is tracked through concrete package inventories. |
| Cookie handling | `framework.cookie` | `cookie` | `partial` | `packages/cookie.md` | Cookie package exists; upstream test parity is not complete. |
| Database query builder | `framework.database` | `database` | `partial` | `packages/database.md` | Query, schema, migrations, and Eloquent surfaces exist; broad upstream parity remains incomplete. |
| Eloquent-style ORM | `framework.database` | `database` | `partial` | `packages/database.md` | ORM primitives exist with Go adaptations for attributes, relations, and errors. |
| Encrypter | `framework.encryption` | `encryption` | `partial` | `packages/encryption.md` | CBC and GCM encryption exist; upstream encryption tests are not fully matched. |
| Event dispatcher | `framework.events` | `events` | `partial` | `packages/events.md` | Event dispatching exists; upstream test parity remains incomplete. |
| Local filesystem | `framework.filesystem` | `filesystem` | `partial` | `packages/filesystem.md` | Local filesystem operations are implemented; upstream parity tests are not complete. |
| Application foundation | `framework.foundation` | n/a | `excluded` | n/a | Full Laravel application kernel is outside Bedrock product boundaries. |
| Password hashing | `framework.hashing` | `hashing` | `partial` | `packages/hashing.md` | Hashing drivers exist; upstream test parity is not complete. |
| HTTP primitives | `framework.http` | `httpx` | `partial` | `packages/httpx.md` | HTTP primitives and test helpers exist; upstream test parity remains incomplete. |
| JSON schema builder | `framework.json_schema` | `jsonx` | `partial` | `packages/jsonx.md` | Fluent JSON schema builder exists; upstream test parity is not complete. |
| Logging channels | `framework.log` | `log` | `partial` | `packages/log.md` | Driver-based structured logging exists; upstream log tests are not matched yet. |
| Mail messages | `framework.mail` | `mailx` | `partial` | `packages/mailx.md` | Mail package exists; upstream mail test parity is incomplete. |
| Notification channels | `framework.notifications` | `notifications` | `partial` | `packages/notifications.md` | Multi-channel notifications exist; upstream parity remains incomplete. |
| Pagination | `framework.pagination` | `pagination` | `partial` | `packages/pagination.md` | Offset and cursor paginators exist; upstream parity is incomplete. |
| Pipeline | `framework.pipeline` | `pipeline` | `partial` | `packages/pipeline.md` | Pipeline package exists; upstream test parity is not complete. |
| Process runner | `framework.process` | n/a | `missing` | n/a | No Bedrock process package exists yet. |
| Queue workers and drivers | `framework.queue` | `queue` | `partial` | `packages/queue.md` | Queue has the strongest upstream test match so far, but missing tests remain. |
| Redis commands | `framework.redis` | `redis` | `partial` | `packages/redis.md` | Redis command surface exists; upstream Redis tests are not matched yet. |
| Routing | `framework.routing` | `routing` | `partial` | `packages/routing.md` | Routing package exists; upstream routing tests remain incomplete. |
| Session stores | `framework.session` | `session` | `partial` | `packages/session.md` | Session stores exist; upstream session tests are not matched yet. |
| Support helpers | `framework.support` | `support` | `partial` | `packages/support.md` | Several support helpers are ported; upstream support coverage is still partial. |
| String helpers | `framework.support_str` | `str` | `partial` | `packages/str.md` | String helper parity has many ported tests, but missing cases remain. |
| Lottery helper | `framework.support_lottery` | `lottery` | `partial` | `packages/lottery.md` | Lottery package exists; upstream Lottery tests are not matched yet. |
| Testing utilities | `framework.testing` | n/a | `missing` | n/a | No standalone Bedrock testing package exists yet. |
| Translation loader | `framework.translation` | `translation` | `partial` | `packages/translation.md` | Translation has some matched upstream tests; full parity remains incomplete. |
| Validation rules | `framework.validation` | `validation` | `partial` | `packages/validation.md` | Rule engine and many rules exist; upstream test parity is still broad and incomplete. |
| Collections | `framework.collections` | n/a | `excluded` | n/a | Collections are not a standalone Bedrock package target. |
| Macroable runtime | `framework.macroable` | n/a | `excluded` | n/a | PHP runtime macro behaviour does not map cleanly to Go static method sets. |
| Reflection helpers | `framework.reflection` | n/a | `excluded` | n/a | PHP reflection helpers are not a Bedrock product surface. |
| Blade views | `framework.view` | n/a | `excluded` | n/a | Bedrock does not ship Laravel Blade or the view engine. |
| AI SDK | `package.ai` | `ai/sdk` | `partial` | `packages/ai/sdk.md` | Bedrock has an AI SDK package; upstream package parity is not complete. |
| Agent tooling | `package.boost` | `ai/boost` | `partial` | `packages/ai/boost.md` | Boost package exists; upstream package parity remains incomplete. |
| Echo client | `package.echo` | `echo` | `partial` | `packages/echo.md` | Tracked by feature inventory because there is no stable upstream test inventory configured. |
| Fortify auth flows | `package.fortify` | `inception` | `partial` | `packages/fortify.md` | Exposed through Inception with Fortify-specific documentation. |
| Jetstream teams | `package.jetstream` | `inception` | `partial` | `packages/jetstream.md` | Exposed through Inception with Jetstream-specific documentation. |
| MCP server | `package.mcp` | `ai/mcp` | `partial` | `packages/ai/mcp.md` | MCP package exists; upstream package parity remains incomplete. |
| OAuth server | `package.passport` | `passport` | `partial` | `packages/passport.md` | Passport package and docs exist; upstream package parity is not complete. |
| Feature flags | `package.pennant` | `pennant` | `partial` | `packages/pennant.md` | Pennant package exists; upstream package parity remains incomplete. |
| Precognition client | `package.precognition` | `precognition` | `partial` | `packages/precognition.md` | Tracked by feature inventory because there is no stable upstream test inventory configured. |
| CLI prompts | `package.prompts` | `prompts` | `partial` | `packages/prompts.md` | Prompt package exists with a small number of matched upstream tests. |
| Reverb websocket server | `package.reverb` | `reverb` | `partial` | `packages/reverb.md` | Reverb package exists; upstream package parity remains incomplete. |
| Search indexing | `package.scout` | `scout` | `partial` | `packages/scout.md` | Scout package exists; upstream package parity remains incomplete. |
| Social OAuth | `package.socialite` | `socialite` | `partial` | `packages/socialite.md` | Socialite package exists; upstream package parity remains incomplete. |
| Subscription billing | `package.spark` | `spark` | `partial` | `packages/spark.md` | Spark is tracked against local Madora billing tests. |
| Telescope observability | `package.telescope` | `telescope` | `partial` | `packages/telescope.md` | Telescope package exists; upstream package parity remains incomplete. |
| Wayfinder route helpers | `package.wayfinder` | `wayfinder` | `partial` | `packages/wayfinder.md` | Wayfinder package exists; upstream package parity remains incomplete. |
| Sanctum tokens | `package.sanctum` | n/a | `excluded` | n/a | Sanctum is intentionally not part of Bedrock. |
| Horizon dashboard | `package.horizon` | n/a | `excluded` | n/a | Queue monitoring UI and runtime are outside Bedrock product boundaries. |
| Browser automation | `package.dusk` | n/a | `excluded` | n/a | Browser automation package is excluded; browser work should use Agent Browser. |

## Mapped Sources Without Inventories

| Source | Bedrock | Reason |
| --- | --- | --- |
| `Illuminate\Contracts` | `contracts` | No generated upstream test inventory configured. |
| `laravel/echo` | `echo` | No generated upstream test inventory configured. |
| `laravel/precognition` | `precognition` | No generated upstream test inventory configured. |

## Permanent Exclusions

- `Illuminate\Foundation` - Product boundary; Bedrock does not ship Laravel's full application kernel/Foundation layer.
- `Illuminate\Collections` - Product boundary; collection behaviour is not a standalone Bedrock package target.
- `Illuminate\Macroable` - PHP runtime macro behaviour does not map cleanly to Go's static method set.
- `Illuminate\Reflection` - PHP reflection helpers are not a Bedrock product surface.
- `Illuminate\View` - Product boundary; Bedrock does not ship Laravel's Blade/view engine.
- `laravel/sanctum` - Product boundary; Sanctum is intentionally not part of Bedrock.
- `laravel/horizon` - Product boundary; queue monitoring UI/runtime is intentionally excluded.
- `laravel/octane` - Product boundary; PHP worker/runtime management does not map to Bedrock.
- `laravel/pulse` - Product boundary; application monitoring dashboard is intentionally excluded.
- `laravel/pail` - Product boundary; log-tail CLI is intentionally excluded.
- `laravel/folio` - Product boundary; page-based PHP routing is intentionally excluded.
- `laravel/dusk` - Product boundary; browser automation package is excluded. Any future browser work must use Agent Browser.
- `laravel/browser-kit-testing` - Product boundary; legacy BrowserKit compatibility is intentionally excluded.
- `laravel/sail` - Product boundary; Docker development environment tooling is intentionally excluded.
- `laravel/installer` - Product boundary; Laravel application installer tooling is intentionally excluded.
- `laravel/valet` - Product boundary; macOS local development service management is intentionally excluded.
- `laravel/envoy` - Product boundary; SSH task runner is intentionally excluded.
- `laravel/vite-plugin` - Product boundary; Vite integration tooling is intentionally excluded.
- `starter-kits` - Product boundary; application starter kits are intentionally excluded.
- `laravel/workos` - Product boundary; WorkOS starter-kit integration is intentionally excluded.
