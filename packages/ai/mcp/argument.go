package mcp

// Argument describes a named parameter that a Prompt accepts.
type Argument struct {
	Name        string
	Description string
	Required    bool
}

// NewArgument creates an Argument. Pass true as the optional required
// parameter to mark the argument as mandatory.
func NewArgument(name, description string, required ...bool) *Argument {
	req := len(required) > 0 && required[0]

	return &Argument{Name: name, Description: description, Required: req}
}

// ToMap serialises the argument for inclusion in a prompts/list response.
func (a *Argument) ToMap() map[string]any {
	return map[string]any{
		"name":        a.Name,
		"description": a.Description,
		"required":    a.Required,
	}
}
