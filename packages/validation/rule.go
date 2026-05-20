package validation

import (
	"reflect"
	"strings"

	contract "github.com/bedrock/packages/contracts/validation"
	"github.com/bedrock/packages/validation/rules"
)

// Rule provides static factory methods for building rule objects, mirroring
// the upstream Illuminate\Validation\Rule facade.

type ruleBuilder struct{}

// In creates an InRule that passes when the value is one of the given values.

// NotIn creates a NotInRule that passes when the value is NOT one of the given values.

// When creates a ConditionalRule that applies trueRules when condition is true,
// or falseRules (if supplied) when condition is false.

// Unless is the inverse of When.

// RequiredIf creates a rule that makes the field required when the callback
// returns true.

// ExcludeIf creates a rule that excludes the field when the callback returns true.

// ProhibitedIf creates a rule that prohibits the field when the callback returns true.

// Password creates a PasswordRule for fluent password validation.

// Array creates an ArrayRule that validates allowed keys for associative arrays.

// ─── InRule ────────────────────────────────────────────────────────────────────

// InRule validates that the value is in a fixed set.
type InRule struct {
	values []string
}

// ─── NotInRule ─────────────────────────────────────────────────────────────────

// NotInRule validates that the value is NOT in a fixed set.
type NotInRule struct {
	values []string
}

// ─── ArrayRule ────────────────────────────────────────────────────────────────

// ArrayRule validates that the value is an array and, when configured, that a
// map contains only the listed keys.
type ArrayRule struct {
	keys []string
}

// ─── ConditionalRule ───────────────────────────────────────────────────────────

// ConditionalRule applies one set of rules when a condition is true and an
// optional fallback set when it is false.  It implements ValidationRule by
// delegating to one of the two rule sets.
type ConditionalRule struct {
	condition    bool
	trueRules    any
	defaultRules any
}

// ActiveRules returns the rule set that should apply given the condition.

// ConditionalRule is resolved by the Validator during rule expansion;
// it is not directly invoked as a ValidationRule.

// ─── CallbackRequiredIfRule ────────────────────────────────────────────────────

// CallbackRequiredIfRule makes the field required when a callback returns true.
type CallbackRequiredIfRule struct {
	cb func() bool
}

// ─── CallbackExcludeIfRule ─────────────────────────────────────────────────────

// CallbackExcludeIfRule excludes the field from validated output when the
// callback returns true.
type CallbackExcludeIfRule struct {
	cb func() bool
}

// ShouldExclude reports whether the field should be excluded.

// ─── CallbackProhibitedIfRule ──────────────────────────────────────────────────

// CallbackProhibitedIfRule prohibits the field when the callback returns true.
type CallbackProhibitedIfRule struct {
	cb func() bool
}

// ─── PasswordRule ──────────────────────────────────────────────────────────────

// PasswordRule provides a fluent interface for password validation requirements.
type PasswordRule struct {
	opts rules.PasswordOptions
}

var Rule = ruleBuilder{}

func (ruleBuilder) In(values ...any) *InRule {
	strs := make([]string, 0, len(values))

	for _, v := range values {
		strs = append(strs, stringify(v))
	}

	return &InRule{values: strs}
}

func (ruleBuilder) NotIn(values ...any) *NotInRule {
	strs := make([]string, 0, len(values))

	for _, v := range values {
		strs = append(strs, stringify(v))
	}

	return &NotInRule{values: strs}
}

func (ruleBuilder) Array(keys ...any) *ArrayRule {
	strs := make([]string, 0, len(keys))

	for _, k := range keys {
		strs = append(strs, stringify(k))
	}

	return &ArrayRule{keys: strs}
}

func (ruleBuilder) When(condition bool, trueRules any, falseRules ...any) *ConditionalRule {
	var defRules any

	if len(falseRules) > 0 {
		defRules = falseRules[0]
	}

	return &ConditionalRule{
		condition:    condition,
		trueRules:    trueRules,
		defaultRules: defRules,
	}
}

func (ruleBuilder) Unless(condition bool, trueRules any, falseRules ...any) *ConditionalRule {
	return Rule.When(!condition, trueRules, falseRules...)
}

func (ruleBuilder) RequiredIf(cb func() bool) *CallbackRequiredIfRule {
	return &CallbackRequiredIfRule{cb: cb}
}

func (ruleBuilder) ExcludeIf(cb func() bool) *CallbackExcludeIfRule {
	return &CallbackExcludeIfRule{cb: cb}
}

func (ruleBuilder) ProhibitedIf(cb func() bool) *CallbackProhibitedIfRule {
	return &CallbackProhibitedIfRule{cb: cb}
}

func (ruleBuilder) Password() *PasswordRule {
	return &PasswordRule{opts: rules.PasswordOptions{Min: 8}}
}

func (r *InRule) Validate(attribute string, value any, fail func(message string)) {
	s := stringify(value)

	for _, v := range r.values {
		if s == v {
			return
		}
	}

	fail("The selected " + attribute + " is invalid.")
}

func (r *InRule) String() string {
	return "in:" + strings.Join(r.values, ",")
}

func (r *NotInRule) Validate(attribute string, value any, fail func(message string)) {
	s := stringify(value)

	for _, v := range r.values {
		if s == v {
			fail("The selected " + attribute + " is invalid.")

			return
		}
	}
}

func (r *NotInRule) String() string {
	return "not_in:" + strings.Join(r.values, ",")
}

func (r *ArrayRule) Validate(attribute string, value any, fail func(message string)) {
	if value == nil {
		fail("The " + attribute + " field must be an array.")

		return
	}

	rv := reflect.ValueOf(value)

	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Map {
		if _, ok := value.([]any); !ok {
			if _, ok := value.(map[string]any); !ok {
				fail("The " + attribute + " field must be an array.")

				return
			}
		}
	}

	if len(r.keys) > 0 {
		m, ok := value.(map[string]any)

		if !ok {
			return
		}

		for k := range m {
			allowed := false

			for _, key := range r.keys {
				if key == k {
					allowed = true

					break
				}
			}

			if !allowed {
				fail("The " + attribute + " field must be an array.")

				return
			}
		}
	}

	if rv.Kind() == reflect.Map || rv.Kind() == reflect.Slice {
		return
	}

	if _, ok := value.([]any); ok {
		return
	}

	if _, ok := value.(map[string]any); ok {
		return
	}

	if len(r.keys) == 0 {
		fail("The " + attribute + " field must be an array.")
	}
}

func (r *ArrayRule) String() string {
	return "array:" + strings.Join(r.keys, ",")
}

func (r *ConditionalRule) ActiveRules() []ParsedRule {
	if r.condition {
		return Explode(r.trueRules)
	}

	if r.defaultRules != nil {
		return Explode(r.defaultRules)
	}

	return nil
}

func (r *ConditionalRule) Validate(attribute string, value any, fail func(message string)) {

}

var _ contract.ImplicitRule = (*CallbackRequiredIfRule)(nil)

func (r *CallbackRequiredIfRule) IsImplicit() bool { return true }

func (r *CallbackRequiredIfRule) Validate(attribute string, value any, fail func(message string)) {
	if r.cb() {
		if isBlank(value) {
			fail("The " + attribute + " field is required.")
		}
	}
}

func (r *CallbackExcludeIfRule) Validate(_ string, _ any, _ func(message string)) {}

func (r *CallbackExcludeIfRule) ShouldExclude() bool {
	return r.cb()
}

var _ contract.ImplicitRule = (*CallbackProhibitedIfRule)(nil)

func (r *CallbackProhibitedIfRule) IsImplicit() bool { return true }

func (r *CallbackProhibitedIfRule) Validate(attribute string, value any, fail func(message string)) {
	if r.cb() && !isBlank(value) {
		fail("The " + attribute + " field is prohibited.")
	}
}

// Min sets the minimum character length.
func (r *PasswordRule) Min(length int) *PasswordRule {
	r.opts.Min = length

	return r
}

// Max sets the maximum character length (0 = no limit).
func (r *PasswordRule) Max(length int) *PasswordRule {
	r.opts.Max = length

	return r
}

// Letters requires at least one letter.
func (r *PasswordRule) Letters() *PasswordRule {
	r.opts.Letters = true

	return r
}

// Numbers requires at least one digit.
func (r *PasswordRule) Numbers() *PasswordRule {
	r.opts.Numbers = true

	return r
}

// Symbols requires at least one non-letter, non-digit character.
func (r *PasswordRule) Symbols() *PasswordRule {
	r.opts.Symbols = true

	return r
}

// MixedCase requires both upper- and lower-case letters.
func (r *PasswordRule) MixedCase() *PasswordRule {
	r.opts.Mixed = true

	return r
}

func (r *PasswordRule) Validate(attribute string, value any, fail func(message string)) {
	s, ok := value.(string)

	if !ok {
		fail("The " + attribute + " must be a string.")

		return
	}

	if !rules.CheckPassword(s, r.opts) {
		fail("The " + attribute + " field format is invalid.")
	}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func isBlank(value any) bool {
	if value == nil {
		return true
	}

	if s, ok := value.(string); ok {
		return strings.TrimSpace(s) == ""
	}

	return false
}
