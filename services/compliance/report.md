# Laravel Compliance Report

Generated: 2026-04-20T04:44:19Z

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
