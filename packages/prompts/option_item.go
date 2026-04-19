package prompts

import "fmt"

// OptionItem represents a single option in a select/multiselect/search prompt.
type OptionItem struct {
	Key   string
	Label string
}

// resolveOptions converts various option formats to a standardized slice.
// Accepts []string or map[string]string. For maps, insertion order is not
// guaranteed — use []string for ordered options.
func resolveOptions(opts any) ([]OptionItem, error) {
	switch v := opts.(type) {
	case []string:
		items := make([]OptionItem, len(v))

		for i, s := range v {
			items[i] = OptionItem{Key: s, Label: s}
		}

		return items, nil
	case map[string]string:
		items := make([]OptionItem, 0, len(v))

		for k, label := range v {
			items = append(items, OptionItem{Key: k, Label: label})
		}

		return items, nil
	case []OptionItem:
		return v, nil
	default:
		return nil, fmt.Errorf("%w: expected []string, map[string]string, or []OptionItem, got %T", ErrInvalidOptions, opts)
	}
}
