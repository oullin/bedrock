package validation

// Validator is the primary validation contract.
type Validator interface {
	// Passes runs all rules and returns true when every rule passes.
	Passes() bool

	// Fails runs all rules and returns true when at least one rule fails.
	Fails() bool

	// Validate runs all rules and returns a *ValidationException when any rule
	// fails, or nil on success.
	Validate() error

	// Validated returns the subset of input data that passed validation.
	// Returns an error if validation has not yet been run or if it failed.
	Validated() (map[string]any, error)

	// Errors returns the MessageBag of accumulated failures.
	Errors() MessageBag

	// Failed returns a map of attribute → []ruleName for every failed rule.
	Failed() map[string][]string

	// SetData replaces the data being validated.
	SetData(data map[string]any) Validator

	// SetRules replaces the rule set.
	SetRules(rules map[string]any) Validator

	// AddRules merges additional rules into the existing set.
	AddRules(rules map[string]any) Validator

	// GetData returns the data under validation.
	GetData() map[string]any

	// HasRule reports whether the given attribute has any of the named rules.
	HasRule(attribute string, rules ...string) bool

	// SetCustomMessages registers custom error messages keyed by
	// "attribute.rule" or just "rule".
	SetCustomMessages(messages map[string]string) Validator

	// SetAttributeNames registers human-readable names for attributes used in
	// error messages in place of the raw field key.
	SetAttributeNames(names map[string]string) Validator
}
