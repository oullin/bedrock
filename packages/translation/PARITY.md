# Upstream Translation Parity — Adaptation Rules

This document governs how `packages/translation` maintains behavioural parity with
`upstream/framework` 13.x `Framework\Translation` and its PHPUnit test suite.

**Source of truth:** https://github.com/upstream/framework/tree/13.x/src/Framework/Translation

**Test baseline:** `testdata/upstream-translation-test-inventory.txt` — one line per PHP
test method. Every entry must have a Go counterpart.

---

## 1. Test naming

Upstream's PHPUnit tests use two naming styles. Both map deterministically to Go:

| PHP style                                                      | Go style                                                        |
| -------------------------------------------------------------- | --------------------------------------------------------------- |
| `public function testGetMethodProperlyLoadsAndRetrievesItem()` | `func TestGetMethodProperlyLoadsAndRetrievesItem(t *testing.T)` |
| `public function test_it_chooses_correct_plural()`             | `func TestItChoosesCorrectPlural(t *testing.T)`                 |

Rules:

1. Strip the `test_` / `test` prefix, then CamelCase the remainder.
2. Tests prefixed with `test_` in PHP (snake_case) become Go tests with the
   underscores removed and each segment capitalised.
3. Every Go test file has `t.Parallel()` unless the test mutates shared global state.

## 2. What cannot port verbatim

Upstream's translation tests rely on features Go has no direct equivalent for.

| Upstream feature                                          | Go adaptation                                                                                  |
| -------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| PHP arrays returned from `.php` files                    | JSON files (`.json`). Behaviour identical — same key/value structure.                          |
| `Filesystem` mock via Mockery                            | Real filesystem using `t.TempDir()` for isolation.                                             |
| PHP `Countable` interface                                | `translation.Countable` interface (`Len() int`).                                               |
| PHP `BackedEnum` / `UnitEnum`                            | Go named types implementing `fmt.Stringer` or registered via `Translator.Stringable()`.        |
| `Stringable` contract (PHP)                              | `fmt.Stringer` interface or `Translator.Stringable(reflect.Type, func(any) string)` registry.  |
| Upstream `MessageSelector::getPluralIndex` int parameter  | Go signature uses `float64`; modulo operations use `ni := int(n)` but equality uses raw float. |
| `config('app.locale')` / `config('app.fallback_locale')` | Constructor parameters: `NewTranslator(loader, locale, fallback)`.                             |
| PHP loose comparison (`1.2 == 1` is false)               | Go enforces strict equality: `n == 1` is false when n=1.2.                                     |

## 3. Where each PHP test file lives in Go

| PHP test file                   | Go test file                       |
| ------------------------------- | ---------------------------------- |
| `TranslationTranslatorTest.php` | `translator_laravel_test.go`       |
| `TranslationFileLoaderTest.php` | `file_loader_laravel_test.go`      |
| `MessageSelectorTest.php`       | `message_selector_laravel_test.go` |

## 4. Coverage rule

A test is ported **only** when its Go counterpart asserts the same
_observable behaviour_ as the PHP test. It is **not enough** to have a test
with a matching name that trivially passes.

Every ported Go test MUST begin with a comment:

```go
// Port of Framework\Tests\Translation\TranslationTranslatorTest::testGetMethodProperlyLoadsAndRetrievesItem
func TestGetMethodProperlyLoadsAndRetrievesItem(t *testing.T) { ... }
```

## 5. When Upstream changes

When Upstream releases a new patch of 13.x:

1. Diff `testdata/upstream-translation-test-inventory.txt` against the upstream PHP test files.
2. Any added PHP test becomes a new task — port it and add to the inventory.
3. Removed PHP tests stay in Go (do not delete tests for defensive reasons), but mark
   them `// Removed upstream in upstream/framework@<sha>`.
4. Port semantics changes promptly; skew is tracked in this file under §6.

## 6. Known semantic divergences

| Divergence                                                            | Rationale                                                                                  |
| --------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| Go errors are values, not thrown exceptions.                          | Language difference. Behaviour preserved via `errors.Is` / `errors.As`.                    |
| `ErrInvalidLocale` returned when locale contains `/` or `\`.          | PHP throws `InvalidArgumentException`; Go returns an error value from `SetLocale`.         |
| No IoC container. `Translator` is constructed directly.               | Go has no global service container. Dependencies are injected via `NewTranslator()`.       |
| `getPluralIndex` accepts `float64`, not `int`.                        | Preserves PHP loose-comparison semantics: `1.2 != 1` so `getPluralIndex("en", 1.2)` → 1.   |
| Tag replacements use `func(string) string` values in the replace map. | PHP uses `HtmlString`-wrapping callables; Go uses plain functions with the same signature. |

_Update this section every time a deliberate skew is introduced._
