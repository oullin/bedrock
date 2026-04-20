# Laravel Compliance Report

Generated: 2026-04-20T04:16:33Z

Source of truth: services/compliance

## Inventories

| Inventory | Total | Ported | Adapted | Missing |
| --- | ---: | ---: | ---: | ---: |
| inventories/framework-queue.txt | 203 | 161 | 0 | 42 |
| inventories/framework-translation.txt | 57 | 9 | 0 | 48 |

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
