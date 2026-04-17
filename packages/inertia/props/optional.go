package props

// OptionalProp is excluded from initial page loads and only included
// when explicitly requested via X-Inertia-Partial-Data.
type OptionalProp struct {
	Value any
}

// Optional creates a prop only included on partial reloads
// when explicitly requested.
func Optional(val any) OptionalProp {
	return OptionalProp{Value: val}
}

func (p OptionalProp) Prop() any { return p.Value }
