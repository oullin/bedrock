# Upstream Compliance Report

Generated: 2026-04-20T04:16:33Z

Source of truth: services/compliance

## Inventories

| Inventory | Total | Ported | Adapted | Missing |
| --- | ---: | ---: | ---: | ---: |
| inventories/framework-queue.txt | 203 | 161 | 0 | 42 |
| inventories/framework-translation.txt | 57 | 9 | 0 | 48 |

## Permanent Exclusions

- `Framework\Foundation` - Product boundary; Bedrock does not ship Upstream's full application kernel/Foundation layer.
- `Framework\Collections` - Product boundary; collection behaviour is not a standalone Bedrock package target.
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
