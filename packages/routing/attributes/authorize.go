package attributes

// Authorize is the parity counterpart of the PHP #[Authorize] attribute.
//
// Ability is the gate name that must be satisfied; Models is the optional
// list of route parameter names whose values should be passed to the gate.
// In PHP this is realized as a class-level attribute; the Go form is a value
// object that callers attach via [controllers.HasMiddleware] returning a
// middleware that performs the gate check.
type Authorize struct {
	Ability string
	Models  []string
}
