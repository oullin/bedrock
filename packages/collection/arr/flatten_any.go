package arr

// FlattenAny flattens nested []any values up to the requested depth.
func FlattenAny(items []any, depth int) []any {
	if depth <= 0 {
		return append([]any(nil), items...)
	}

	result := make([]any, 0, len(items))

	for _, item := range items {
		if nested, ok := item.([]any); ok {
			result = append(result, FlattenAny(nested, depth-1)...)

			continue
		}

		result = append(result, item)
	}

	return result
}
