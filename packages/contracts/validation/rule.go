package validation

// ValidationRule is the contract for custom validation rules.
// Implement this interface to create reusable rule objects.
//
// Call fail with a message to record a failure; if fail is never called the
// value is considered valid.
type ValidationRule interface {
	Validate(attribute string, value any, fail func(message string))
}

// ImplicitRule marks a ValidationRule as implicit — it is evaluated even when
// the field is absent or its value is empty.  Rules that require presence
// (e.g. Required) must implement this interface.
type ImplicitRule interface {
	ValidationRule
	IsImplicit() bool
}
