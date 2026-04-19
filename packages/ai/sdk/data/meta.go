package data

// Meta holds provider and model metadata attached to a response.
// Mirrors Upstream\Ai\Responses\Data\Meta.
type Meta struct {
	Provider  *string `json:"provider"`
	Model     *string `json:"model"`
	Citations []any   `json:"citations"`
}

// NewMeta constructs a Meta with optional provider and model strings.
func NewMeta(provider, model *string) Meta {
	return Meta{
		Provider:  provider,
		Model:     model,
		Citations: []any{},
	}
}

// ToMap returns a map representation.
func (m Meta) ToMap() map[string]any {
	return map[string]any{
		"provider":  m.Provider,
		"model":     m.Model,
		"citations": m.Citations,
	}
}
