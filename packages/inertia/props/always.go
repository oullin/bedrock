package props

// AlwaysProp is included in every response, even partial reloads
// that do not explicitly request it.
type AlwaysProp struct {
	Value any
}

// Always creates a prop always included in every response.
func Always(val any) AlwaysProp {
	return AlwaysProp{Value: val}
}

func (p AlwaysProp) Prop() any { return p.Value }
