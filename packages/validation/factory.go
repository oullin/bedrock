package validation

import "github.com/bedrock/packages/validation/rules"

// Factory creates Validator instances.  It holds shared extensions and
// replacers that are propagated to every validator it creates.
type Factory struct {
	extensions  map[string]rules.RuleFunc
	implicitExt map[string]bool
	pv          rules.PresenceVerifier
}

// NewFactory returns a ready-to-use Factory.
func NewFactory() *Factory {
	return &Factory{
		extensions:  make(map[string]rules.RuleFunc),
		implicitExt: make(map[string]bool),
	}
}

// Make creates a new Validator without running it.
//
//   - data:       the input to validate (map[string]any)
//   - ruleMap:    rules keyed by attribute name; values may be strings,
//                 []string, []any, or ValidationRule objects
//   - messages:   custom error messages (may be nil)
//   - attributes: human-readable attribute names (may be nil)
func (f *Factory) Make(
	data map[string]any,
	ruleMap map[string]any,
	messages map[string]string,
	attributes map[string]string,
) *Validator {
	// Copy extensions so each validator gets its own independent map
	ext := make(map[string]rules.RuleFunc, len(f.extensions))
	for k, v := range f.extensions {
		ext[k] = v
	}

	impl := make(map[string]bool, len(f.implicitExt))
	for k, v := range f.implicitExt {
		impl[k] = v
	}

	return newValidator(data, ruleMap, messages, attributes, ext, impl, f.pv)
}

// Validate creates a Validator and immediately runs it, returning the
// validated data on success or a *ValidationException on failure.
func (f *Factory) Validate(
	data map[string]any,
	ruleMap map[string]any,
	messages map[string]string,
	attributes map[string]string,
) (map[string]any, error) {
	v := f.Make(data, ruleMap, messages, attributes)
	return v.Validated()
}

// Extend registers a custom validation rule function.  Rules registered here
// are available to every Validator created by this factory.
func (f *Factory) Extend(name string, fn rules.RuleFunc) {
	f.extensions[StudlyCase(name)] = fn
}

// ExtendImplicit registers a custom implicit validation rule — one that runs
// even when the field is absent.
func (f *Factory) ExtendImplicit(name string, fn rules.RuleFunc) {
	studly := StudlyCase(name)
	f.extensions[studly] = fn
	f.implicitExt[studly] = true
}

// SetPresenceVerifier sets the database presence verifier for all validators
// created by this factory.
func (f *Factory) SetPresenceVerifier(pv rules.PresenceVerifier) {
	f.pv = pv
}

