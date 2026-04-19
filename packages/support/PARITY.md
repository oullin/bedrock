# Laravel Support Parity — Split Package Notes

This document tracks how Bedrock's split support packages map onto
`laravel/framework` 13.x `Illuminate\Support` and its PHPUnit test suite.
The implementation now lives across three modules:

- `packages/support` for helpers, `Fluent`, `Optional`, `MessageBag`, `Sleep`, and `Timebox`
- `packages/str` for the `Str*` API and `StringBuilder`
- `packages/lottery` for the `Lottery` API

**Source of truth:** https://github.com/laravel/framework/tree/13.x/src/Illuminate/Support

**Test baseline:** Every PHP test method must have a Go counterpart with a
`// Port of Illuminate\Tests\Support\SupportXxxTest::testMethodName` header comment.

---

## 1. Test naming

Laravel's PHPUnit tests use two naming styles. Both map deterministically to Go:

| PHP style                                                | Go style                                                  |
| -------------------------------------------------------- | --------------------------------------------------------- |
| `public function testFluentAttributesSetByConstructor()` | `func TestFluentAttributesSetByConstructor(t *testing.T)` |
| `public function test_it_can_freeze_uuids()`             | `func TestItCanFreezeUuids(t *testing.T)`                 |

Rules:

1. Strip the `test_` / `test` prefix, then CamelCase the remainder.
2. Tests prefixed with `test_` in PHP (snake_case) become Go tests with underscores removed and each segment capitalised.
3. Every test has `t.Parallel()` unless it mutates package-level global state (UUID/ULID/random factory, sleep function). Those are marked with `// NOT parallel`.

## 2. PHP → Go semantic differences

| PHP feature                                         | Go adaptation                                                                                                                                       |
| --------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Str::` static class methods                        | Package-level functions prefixed `Str*` in `packages/str`                                                                                           |
| `Stringable` class (`__toString`, fluent chain)     | `str.StringBuilder` wrapping a string; `Of(value)` constructor; methods return `*StringBuilder`; `String()` implements `fmt.Stringer`               |
| `Str::freezeUuids()` global state                   | Package-level mutex-guarded `uuidFactory` var; tests using freeze must NOT call `t.Parallel()` at the top level                                     |
| `Str::createUuidsUsingSequence()`                   | `CreateUuidsUsingSequence([]string)` returns cleanup func; same mutex pattern                                                                       |
| `optional($x)->method()` null proxy                 | `support.Optional[T]` generic struct; `Some(v)`, `None[T]()`, `Opt(ptr)`; methods: `Get`, `OrElse`, `IsPresent`, `IfPresent`, `Map`, `Filter`       |
| `Fluent->__get($key)` / `__set($key, $value)` magic | Explicit `support.Fluent` `Get(key)` / `Set(key, value)` methods; no dynamic property access                                                        |
| `Fluent::fill()` vs `Fluent::merge()`               | `Fill` overwrites existing keys; `Merge` skips existing keys                                                                                        |
| `MessageBag::has()` with wildcard                   | `support.MessageBag` uses `filepath.Match` for `*` glob patterns                                                                                    |
| `Lottery::alwaysWin()` / `alwaysLose()`             | `lottery.Always()` / `Never()` (also `ForceWin()` / `ForceLose()` aliases)                                                                          |
| `Lottery::fix(sequence)`                            | `lottery.Fix([]bool)` returns a `*fixedLottery` with predetermined outcomes                                                                         |
| `Sleep::fake()` static global                       | `support.FakeSleepWith(fake)` installs and returns a cleanup closure                                                                                |
| `Str::plural()` non-English                         | English-only via `jinzhu/inflection`; known mismatches documented below                                                                             |
| `Str::mask()` with negative index                   | Negative index counts from end: `-n` starts at `len - n`                                                                                            |
| `Macroable` trait                                   | Not ported; Go has no dynamic method registration                                                                                                   |
| PHP exceptions                                      | Go `error` values; no exception hierarchy                                                                                                           |
| `Retry` with sleep                                  | `Retry(times int, fn func(attempt int) error, sleep ...any)` where sleep may be an int (ms), `[]int` backoff schedule, or `func(int) time.Duration` |

## 3. Known semantic divergences

### `str.StrPlural` — English-only

`StrPlural` / `StrSingular` use `jinzhu/inflection` which is English-only. Non-English
pluralisation (Laravel supports custom language files) is not supported.

Additionally, `jinzhu/inflection` applies a broad "-man" → "-men" suffix rule that
incorrectly pluralises "human" → "humen". This is corrected by a local `pluralOverrides`
map in `str_plural.go`. Additional overrides should be added to that map as needed.

### Fluent dot-notation

Laravel's `Fluent` supports dot-notation for nested access (e.g. `$f->get('a.b.c')`).
The Go port's `Get` / `Set` use simple flat keys only. Dot-notation navigation is not
implemented because the primary consumers in this codebase use flat maps.

### Optional[T] type parameter

PHP's `optional()` is a runtime null-proxy that intercepts any method call. Go's
`Optional[T any]` is a compile-time generic wrapper. The type must be declared at
construction time (`Some[int](5)`). Cross-type mapping requires explicit `Map`.

### UUID/ULID ordered generation

`StrOrderedUuid()` uses `uuid.NewV7()` (time-ordered). PHP uses a custom algorithm
seeded from `Carbon::now()`. The v7 UUID is compliant with RFC 9562 and lexicographically
sortable, satisfying the same purpose.

### Markdown rendering

`StrMarkdown` / `StrInlineMarkdown` use `yuin/goldmark` with the GFM extension.
PHP uses `league/commonmark`. Output is semantically identical for standard Markdown;
minor whitespace differences may exist in edge cases.

### Str::ascii — transliteration map

The ASCII transliteration map in `str_ascii.go` covers Latin-1, Latin Extended-A,
Greek, Cyrillic, and common symbols. It does not include the full multi-language
map from Laravel's `lang/` directory. Language-specific overrides are not supported.

## 4. Test file map

| PHP test class                                   | Go test file                                                                 |
| ------------------------------------------------ | ---------------------------------------------------------------------------- |
| `Illuminate\Tests\Support\SupportStrTest`        | `packages/str/str_test.go`, `packages/str/str_uuid_test.go`                  |
| `Illuminate\Tests\Support\SupportFluentTest`     | `packages/support/fluent_test.go`                                            |
| `Illuminate\Tests\Support\SupportOptionalTest`   | `packages/support/optional_test.go`                                          |
| `Illuminate\Tests\Support\SupportMessageBagTest` | `packages/support/message_bag_test.go`                                       |
| `Illuminate\Tests\Support\SupportHelpersTest`    | `packages/support/helpers_test.go`, `packages/support/helpers_retry_test.go` |
| `Illuminate\Tests\Support\LotteryTest`           | `packages/lottery/lottery_test.go`                                           |
| `Illuminate\Tests\Support\SleepTest`             | `packages/support/sleep_test.go`                                             |
| `Illuminate\Tests\Support\TimeboxTest`           | `packages/support/timebox_test.go`                                           |

## 5. Coverage rule

A test is ported **only** when its Go counterpart asserts the same _observable behaviour_
as the PHP test. It is **not enough** to have a test that compiles — the assertion must
cover the same logical invariant.

Every ported test must include a comment header:

```
// Port of Illuminate\Tests\Support\SupportXxxTest::testMethodName
```

Tests without this header are either:

- New Go-specific tests (e.g. `TestOptionalFilter`, `TestLotterySequenceExhausted`)
- Tests covering Go-specific adaptation details not present in the PHP suite

## 6. Excluded components

The following Laravel `Illuminate\Support` components are **explicitly excluded** from
this port by user decision:

| Excluded                      | Reason                                                              |
| ----------------------------- | ------------------------------------------------------------------- |
| `Number`                      | Number formatting; use `fmt` / `golang.org/x/text/message` directly |
| `Once`                        | Memoization; Go's `sync.Once` is the idiomatic equivalent           |
| `Benchmark`                   | Benchmarking utility; Go's `testing.B` is the idiomatic equivalent  |
| Facades / ServiceProvider     | Handled by `packages/container`                                     |
| `Carbon` / DateFactory        | Go's `time` package                                                 |
| `Macroable`                   | PHP-specific dynamic dispatch; no Go equivalent                     |
| `ProcessUtils`, `BinaryCodec` | PHP ecosystem specific                                              |
