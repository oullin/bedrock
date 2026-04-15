# session

Session management with multiple storage handlers.

## Overview

The `session` package provides Upstream-inspired HTTP session management with
flash data, CSRF tokens, and swappable storage backends.

**Module:** `github.com/gocanto/bedrock/packages/session`

```bash
go get github.com/gocanto/bedrock/packages/session@latest
```

## Session Store

`session.Store` is the primary type. It exposes:

- `Get` / `Put` / `Forget` / `Flush`
- `Flash` / `Reflash` / `KeepFlash` — single-request carry-over values
- `Token` — CSRF token generation and retrieval
- `Regenerate` / `Invalidate` — session ID lifecycle management

## Usage

```go
handler := session.NewFileHandler("/storage/sessions", 120)
store := session.NewStore("my-session", handler, encrypter)

// Start the session (reads from storage)
store.Start()

// Read and write values
store.Put("user_id", 42)
userID := store.Get("user_id") // 42

// Flash data — available only on the next request
store.Flash("status", "Profile updated.")

// CSRF
token := store.Token()

// Regenerate ID (after login)
store.Regenerate(true)

// Persist to storage
store.Save()
```

## Handlers

| Handler      | Backend                                           |
| ------------ | ------------------------------------------------- |
| `array`      | In-memory map (testing only)                      |
| `file`       | Local filesystem                                  |
| `database`   | SQL table                                         |
| `cache`      | Any `cache.Store` backend                         |
| `cookie`     | Encrypted cookie                                  |
| `null`       | No-op — discards all data                         |
| `encrypting` | Wraps another handler with transparent encryption |
