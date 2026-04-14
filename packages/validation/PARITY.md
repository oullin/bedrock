# Upstream Validation Parity — Adaptation Rules

This document governs how `packages/validation` maintains behavioural parity with
`upstream/framework` 13.x `Framework\Validation` and its PHPUnit test suite.

**Source of truth:** https://github.com/upstream/framework/tree/13.x/src/Framework/Validation

**Test baseline:** `validator_laravel_test.go` — every test begins with a comment
referencing its PHP source.

---

## 1. Test naming

| PHP style                                              | Go style                                 |
| ------------------------------------------------------ | ---------------------------------------- |
| `public function testFooBarBaz()`                      | `func TestFooBarBaz(t *testing.T)`      |
| `public function test_foo_bar()`                       | `func TestFooBar(t *testing.T)`         |

Every ported Go test **must** begin with:

```go
// Port of Framework\Tests\Validation\ValidationValidatorTest::testFoo
func TestFoo(t *testing.T) { ... }
```

---

## 2. What cannot port verbatim

| Upstream feature                                     | Go adaptation                                                           |
| --------------------------------------------------- | ----------------------------------------------------------------------- |
| IoC container                                       | Not needed — rules are plain functions.                                  |
| Orm `exists`/`unique` DB integration           | `PresenceVerifier` interface — test with a fake verifier.               |
| PHP exceptions                                      | `*ValidationException` error value wrapping a `*MessageBag`.            |
| `Closure`-based rules                               | `ValidationRule` interface (Validate method) or `RuleFunc`.             |
| PHP `enum` validation                               | Not implemented — Go has no built-in enum type.                         |
| `SerializableClosure`                               | Out of scope.                                                           |
| `$validator->sometimes(attr, rules, closure)`       | Use `"sometimes"` as a rule marker prefix (field skipped when absent).  |
| Breach-check for Password rule                      | Not implemented — requires external HTTP call; marked optional.         |
| `ActiveUrl` (DNS lookup)                            | Implemented via `net.LookupHost` — may fail in offline environments.    |

---

## 3. Coverage

Ported test groups (in `validator_laravel_test.go`):

| PHP test method                              | Go test function               | Status |
| -------------------------------------------- | ------------------------------ | ------ |
| testPassesReturnsTrueIfNoFailingRules        | TestPassesReturnsTrueIfNoFailingRules | ✅ |
| testFailsReturnsFalseIfFailingRules          | TestFailsReturnsFalseIfFailingRules   | ✅ |
| testHasFailedRules                           | TestHasFailedRules                    | ✅ |
| testValidateRequired                         | TestValidateRequired                  | ✅ |
| testValidateRequiredIf                       | TestValidateRequiredIf                | ✅ |
| testValidateRequiredUnless                   | TestValidateRequiredUnless            | ✅ |
| testValidateRequiredWith                     | TestValidateRequiredWith              | ✅ |
| testValidateRequiredWithout                  | TestValidateRequiredWithout           | ✅ |
| testValidatePresent                          | TestValidatePresent                   | ✅ |
| testValidateFilled                           | TestValidateFilled                    | ✅ |
| testValidateMissing                          | TestValidateMissing                   | ✅ |
| testValidateProhibited                       | TestValidateProhibited                | ✅ |
| testValidateProhibitedIf                     | TestValidateProhibitedIf              | ✅ |
| testValidateAccepted                         | TestValidateAccepted                  | ✅ |
| testValidateIn                               | TestValidateIn                        | ✅ |
| testValidateNotIn                            | TestValidateNotIn                     | ✅ |
| testValidateMin                              | TestValidateMin                       | ✅ |
| testValidateMax                              | TestValidateMax                       | ✅ |
| testValidateBetween                          | TestValidateBetween                   | ✅ |
| testValidateSize                             | TestValidateSize                      | ✅ |
| testValidateEmail                            | TestValidateEmail                     | ✅ |
| testValidateUrl                              | TestValidateUrl                       | ✅ |
| testValidateIp                               | TestValidateIp                        | ✅ |
| testValidateAlpha                            | TestValidateAlpha                     | ✅ |
| testValidateAlphaDash                        | TestValidateAlphaDash                 | ✅ |
| testValidateAlphaNum                         | TestValidateAlphaNum                  | ✅ |
| testValidateNumeric                          | TestValidateNumeric                   | ✅ |
| testValidateInteger                          | TestValidateInteger                   | ✅ |
| testValidateBoolean                          | TestValidateBoolean                   | ✅ |
| testValidateDate                             | TestValidateDate                      | ✅ |
| testValidateDateFormat                       | TestValidateDateFormat                | ✅ |
| testValidateBefore                           | TestValidateBefore                    | ✅ |
| testValidateAfter                            | TestValidateAfter                     | ✅ |
| testValidateSame                             | TestValidateSame                      | ✅ |
| testValidateDifferent                        | TestValidateDifferent                 | ✅ |
| testValidateConfirmed                        | TestValidateConfirmed                 | ✅ |
| testValidateDistinct                         | TestValidateDistinct                  | ✅ |
| testValidateArray                            | TestValidateArray                     | ✅ |
| testValidateBail                             | TestValidateBail                      | ✅ |
| testValidateNullable                         | TestValidateNullable                  | ✅ |
| testSometimesWorksOnNestedArrays             | TestSometimesWorksOnNestedArrays      | ✅ |
| testCustomValidationRules                    | TestCustomValidationRules             | ✅ |
| testWildcardNestedRules                      | TestWildcardNestedRules               | ✅ |
| testConditionalRules                         | TestConditionalRules                  | ✅ |
| testValidateRegex                            | TestValidateRegex                     | ✅ |
| testValidateUUID                             | TestValidateUUID                      | ✅ |
| testValidateJson                             | TestValidateJson                      | ✅ |
| testValidateStartsWith                       | TestValidateStartsWith                | ✅ |
| testValidateEndsWith                         | TestValidateEndsWith                  | ✅ |
| testValidateInArray                          | TestValidateInArray                   | ✅ |
| testValidateHexColor                         | TestValidateHexColor                  | ✅ |
| testValidateTimezone                         | TestValidateTimezone                  | ✅ |

---

## 4. Known semantic divergences

| Divergence                                              | Rationale                                                           |
| ------------------------------------------------------- | ------------------------------------------------------------------- |
| Go errors are values, not thrown exceptions.            | `Validate()` returns `*ValidationException` (implements `error`).   |
| No PHP `enum` type validation rule.                     | Go has `iota` enums but no runtime reflection equivalent.           |
| Password breach-check (`Uncompromised`) is not wired.   | Requires external HTTP — omitted by default; caller can extend.     |
| `ActiveUrl` uses `net.LookupHost`, not `checkdnsrr`.   | Equivalent semantics; may differ on edge cases.                     |
| Hex color #RGBA and #RRGGBBAA are valid (4 / 8 chars). | Matches Upstream 13.x — both RGB and RGBA shorthand accepted.        |

_Update this section every time a deliberate skew is introduced._
