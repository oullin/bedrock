package workflow

// Transition is an arc in the Petri-Net moving tokens from `From` places to `To` places.
type Transition struct {
	Name     string
	From     []string
	To       []string
	Metadata map[string]any
}

func NewTransition(name string, from []string, to []string) Transition {
	return Transition{
		Name: name,
		From: append([]string(nil), from...),
		To:   append([]string(nil), to...),
	}
}
