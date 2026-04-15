# Upstream Queue Parity — Adaptation Rules

This document governs how `packages/queue` maintains behavioural parity with
`upstream/framework` 13.x `Framework\Queue` and its PHPUnit test suite.

**Source of truth:** https://github.com/upstream/framework/tree/13.x/src/Framework/Queue

**Test baseline:** `testdata/upstream-queue-test-inventory.txt` — one line per PHP
test method. Every entry must have a Go counterpart.

---

## 1. Test naming

Upstream's PHPUnit tests use two naming styles. Both map deterministically to Go:

| PHP style                                                        | Go style                                                   |
| ---------------------------------------------------------------- | ---------------------------------------------------------- |
| `public function testJobCanBeFired()`                            | `func TestJobCanBeFired(t *testing.T)`                     |
| `public function test_it_can_create_timeout_exception_for_job()` | `func TestItCanCreateTimeoutExceptionForJob(t *testing.T)` |

Rules:

1. Strip the `test_` / `test` prefix, then CamelCase the remainder.
2. Tests prefixed with `test_` in PHP (snake_case) become Go tests with the
   underscores removed and each segment capitalised.
3. The Go test file lives in the location mapped in `packages/queue/PARITY.md`
   (see §3) and has `t.Parallel()` unless the test mutates shared global state.

## 2. What cannot port verbatim

Upstream's queue tests rely on features Go has no direct equivalent for. These
are the adaptation rules:

| Upstream feature                                               | Go adaptation                                                                                                                                                |
| ------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ | --- | ---------------------------------------- |
| IoC container / `App::make`                                   | `HandlerRegistry` keyed by name (see `registry.go`).                                                                                                         |
| Orm model serialisation (`SerializesModels`)             | Out of scope. Jobs carry plain structs / maps.                                                                                                               |
| Closures via `SerializableClosure`                            | Replaced by registered handler names.                                                                                                                        |
| `#[Tries]`, `#[Backoff]`, `#[Timeout]`, `#[Queue]` attributes | Go struct tags: `queue:"tries=3,backoff=1s                                                                                                                   | 5s  | 10s,timeout=60s"`parsed via`options.go`. |
| `config('queue.*')` helpers                                   | Explicit config maps passed to `Manager.SetConfig`.                                                                                                          |
| `Event::fake()` / `Bus::fake()`                               | Test helpers in `internal/testutil` (created as needed).                                                                                                     |
| Database tests hitting real MySQL/SQLite                      | Ported to table-driven unit tests with mock executors. Tests that MUST hit a real DB get a `_integration_test.go` suffix and a `//go:build integration` tag. |
| Redis tests hitting real Redis                                | Mock Redis client (`drivers/drivers_test.go`) grown to support `EVAL`/`EVALSHA`.                                                                             |
| SQS tests hitting real AWS                                    | Mock SQS client — same pattern as existing `drivers/sqs_test.go`.                                                                                            |
| PHPUnit `$this->expectException(...)` on specific type        | `errors.As(err, &target)` with matching error struct.                                                                                                        |
| PHP exceptions extending other exceptions                     | Go error structs that embed the parent error via pointer (see `timeout_exceeded_error.go`).                                                                  |
| `resolveName()` on the PHP `Job` class                        | `ResolveName() string` method on the Go `Job` interface — introduced in Step 8. Until then, a local `ResolveNamer` interface is used internally.             |

## 3. Where each PHP test file lives in Go

See the "Test file map" table in `/Users/gocanto/.claude/plans/jazzy-finding-island.md`.

Tests that do not have an obvious package home (e.g. cross-cutting concerns)
live at the top level of `packages/queue`.

## 4. Coverage rule

A test is ported **only** when its Go counterpart asserts the same
_observable behaviour_ as the PHP test. It is **not enough** to have a test
with a matching name that trivially passes. CI enforces this two ways:

1. **Inventory check** (`scripts/queue-parity.sh`): the set of Go `Test*`
   functions must be a superset of the inventory file.
2. **Behaviour check** (manual review): on every PR that changes behaviour,
   the reviewer confirms the ported Go test asserts the same thing as the PHP
   test named in its comment header.

Every ported Go test MUST begin with a comment:

```go
// Port of Framework\Tests\Queue\QueueExceptionTest::test_it_can_create_timeout_exception_for_job
func TestItCanCreateTimeoutExceptionForJob(t *testing.T) { ... }
```

The `scripts/queue-parity.sh` script greps for these headers and flags any
inventory entry without a matching comment.

## 5. When Upstream changes

When Upstream releases a new patch of 13.x:

1. Re-run the inventory generator (see `scripts/queue-parity.sh --refresh`).
2. Diff `testdata/upstream-queue-test-inventory.txt` against its previous
   version; any added entry becomes a new task in the backlog.
3. Removed entries stay in Go (we do not delete tests for defensive reasons),
   but they are marked `// Removed upstream in upstream/framework@<sha>`.
4. Port semantics changes promptly; skew is tracked in this file under §6.

## 6. Known semantic divergences

| Divergence                                                                                         | Rationale                                                                         |
| -------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| Go errors are values, not thrown exceptions.                                                       | Language difference. Behaviour preserved via `errors.As` / `errors.Is`.           |
| No transaction-aware dispatch by default — consumer must inject a `TxCallback` (Step 14).          | Go has no global DB facade analogous to Upstream's `DB::transaction`.              |
| No runtime memory cap in core Worker — `MaxMemoryMiB` is advisory, checked via `runtime.MemStats`. | Go's GC makes PHP-style hard memory caps meaningless; Upstream behaviour imitated. |
| `#[UniqueFor]` requires a `CacheLock` interface; default implementation is in-memory.              | Upstream uses the cache store directly; we avoid a hard dep.                       |

_Update this section every time a deliberate skew is introduced._
