# inertia

Server-side Go adapter for the Inertia.js protocol.

## Overview

The `inertia` package renders Inertia pages — returning JSON on XHR visits and
full HTML on the first request. It merges shared and per-request props, manages
the `<head>` (title, meta, links), and integrates with CSRF, i18n, precognition,
and flash middleware.

**Module:** `github.com/bedrock/packages/inertia`

```bash
go get github.com/bedrock/packages/inertia@latest
```

## Features

- Automatic XHR vs. full-page response detection
- Shared props merged into every page response
- Head management (title, meta tags, link tags)
- CSRF token injection
- Precognition middleware integration
- Flash message propagation

## Coming Soon

Full documentation is in progress.
