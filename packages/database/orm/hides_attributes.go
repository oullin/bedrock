package orm

// HidesAttributes controls which attributes are visible during serialization.
type HidesAttributes struct {
	hidden  []string
	visible []string
}

// SetHidden sets the attributes that should be hidden during serialization.
func (h *HidesAttributes) SetHidden(hidden []string) {
	h.hidden = hidden
}

// GetHidden returns the hidden attributes.
func (h *HidesAttributes) GetHidden() []string {
	return h.hidden
}

// SetVisible sets the attributes that should be visible during serialization.
func (h *HidesAttributes) SetVisible(visible []string) {
	h.visible = visible
}

// GetVisible returns the visible attributes.
func (h *HidesAttributes) GetVisible() []string {
	return h.visible
}

// MakeVisible makes the given attributes visible.
func (h *HidesAttributes) MakeVisible(attributes ...string) {
	h.hidden = filterStrings(h.hidden, attributes)
}

// MakeHidden hides the given attributes.
func (h *HidesAttributes) MakeHidden(attributes ...string) {
	h.hidden = append(h.hidden, attributes...)
}

// FilterAttributes filters a map of attributes based on hidden/visible rules.
func (h *HidesAttributes) FilterAttributes(attrs map[string]any) map[string]any {
	if len(h.visible) > 0 {
		filtered := make(map[string]any, len(h.visible))

		for _, key := range h.visible {
			if v, ok := attrs[key]; ok {
				filtered[key] = v
			}
		}

		return filtered
	}

	if len(h.hidden) > 0 {
		filtered := make(map[string]any, len(attrs))

		for k, v := range attrs {
			if !containsString(h.hidden, k) {
				filtered[k] = v
			}
		}

		return filtered
	}

	return attrs
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

func filterStrings(slice []string, exclude []string) []string {
	var result []string

	for _, s := range slice {
		if !containsString(exclude, s) {
			result = append(result, s)
		}
	}

	return result
}
