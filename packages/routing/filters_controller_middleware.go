package routing

// same name. It exposes the static helper used by [ControllerDispatcher] and
// [ResourceRegistrar] to decide whether a middleware applies to a given
// controller method based on its only/except options.
type FiltersControllerMiddleware struct{}

// MethodExcludedByOptions reports whether method is excluded from the option
// set's "only"/"except" filters.
//
// The rules mirror PHP exactly:
//   - If "only" is present and method is not in it → excluded.
//   - If "except" is non-empty and method is in it → excluded.
func MethodExcludedByOptions(method string, options map[string]any) bool {
	if only, ok := options["only"]; ok {
		if list, ok := toStringSlice(only); ok {
			if !containsString(list, method) {
				return true
			}
		}
	}

	if except, ok := options["except"]; ok {
		if list, ok := toStringSlice(except); ok && len(list) > 0 {
			if containsString(list, method) {
				return true
			}
		}
	}

	return false
}

// MethodExcludedByOptions is the method form of the package-level helper, kept
// for parity with the PHP trait's static method dispatch.
func (FiltersControllerMiddleware) MethodExcludedByOptions(method string, options map[string]any) bool {
	return MethodExcludedByOptions(method, options)
}

func toStringSlice(v any) ([]string, bool) {
	switch x := v.(type) {
	case string:
		return []string{x}, true
	case []string:
		return x, true
	case []any:
		out := make([]string, 0, len(x))

		for _, e := range x {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}

		return out, true
	}

	return nil, false
}

func containsString(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}

	return false
}
