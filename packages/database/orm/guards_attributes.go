package orm

// GuardsAttributes provides mass assignment protection.
type GuardsAttributes struct {
	fillable  []string
	guarded   []string
	unguarded bool
}

// SetFillable sets the mass-assignable attributes.
func (g *GuardsAttributes) SetFillable(fillable []string) {
	g.fillable = fillable
}

// GetFillable returns the mass-assignable attributes.
func (g *GuardsAttributes) GetFillable() []string {
	return g.fillable
}

// MergeFillable adds to the fillable attributes.
func (g *GuardsAttributes) MergeFillable(fillable []string) {
	g.fillable = append(g.fillable, fillable...)
}

// SetGuarded sets the non-mass-assignable attributes.
func (g *GuardsAttributes) SetGuarded(guarded []string) {
	g.guarded = guarded
}

// GetGuarded returns the non-mass-assignable attributes.
func (g *GuardsAttributes) GetGuarded() []string {
	if g.guarded == nil {
		return []string{"*"}
	}

	return g.guarded
}

// IsFillable checks if an attribute is mass-assignable.
func (g *GuardsAttributes) IsFillable(key string) bool {
	if g.unguarded {
		return true
	}

	if len(g.fillable) > 0 {
		for _, f := range g.fillable {
			if f == key {
				return true
			}
		}

		return false
	}

	return !g.IsGuarded(key)
}

// IsGuarded checks if an attribute is guarded from mass assignment.
func (g *GuardsAttributes) IsGuarded(key string) bool {
	if g.unguarded {
		return false
	}

	for _, grd := range g.GetGuarded() {
		if grd == "*" || grd == key {
			return true
		}
	}

	return false
}

// TotallyGuarded checks if all attributes are guarded.
func (g *GuardsAttributes) TotallyGuarded() bool {
	if g.unguarded {
		return false
	}

	return len(g.fillable) == 0 && len(g.guarded) == 1 && g.guarded[0] == "*"
}

// Unguard disables mass assignment protection.
func (g *GuardsAttributes) Unguard() {
	g.unguarded = true
}

// Reguard re-enables mass assignment protection.
func (g *GuardsAttributes) Reguard() {
	g.unguarded = false
}

// IsUnguarded returns whether mass assignment protection is disabled.
func (g *GuardsAttributes) IsUnguarded() bool {
	return g.unguarded
}

// Fill fills the model with an array of attributes, respecting mass assignment.
func (g *GuardsAttributes) Fill(attrs *HasAttributes, values map[string]any) error {
	for key, value := range values {
		if g.IsFillable(key) {
			attrs.SetAttribute(key, value)
		} else if !g.unguarded {
			return ErrMassAssignment
		}
	}

	return nil
}

// ForceFill fills the model ignoring mass assignment protection.
func (g *GuardsAttributes) ForceFill(attrs *HasAttributes, values map[string]any) {
	for key, value := range values {
		attrs.SetAttribute(key, value)
	}
}
