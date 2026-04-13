package routing

// GroupAttributes holds the configuration for a route group.
type GroupAttributes struct {
	Prefix            string
	Domain            string
	NamePrefix        string
	Middleware        []MiddlewareFunc
	WithoutMiddleware []MiddlewareFunc
	Where             map[string]string
}

// MergeGroup merges parent and child group attributes following Upstream
// conventions: prefixes stack, middleware inherits, domain overrides,
// where constraints merge, and name prefixes concatenate with a dot.
func MergeGroup(parent, child GroupAttributes) GroupAttributes {
	merged := GroupAttributes{
		Prefix:     parent.Prefix + child.Prefix,
		Domain:     child.Domain,
		NamePrefix: parent.NamePrefix,
		Where:      make(map[string]string),
	}

	if merged.Domain == "" {
		merged.Domain = parent.Domain
	}

	if child.NamePrefix != "" {
		if merged.NamePrefix != "" {
			merged.NamePrefix += "." + child.NamePrefix
		} else {
			merged.NamePrefix = child.NamePrefix
		}
	}

	merged.Middleware = append(
		append([]MiddlewareFunc(nil), parent.Middleware...),
		child.Middleware...,
	)

	merged.WithoutMiddleware = append(
		append([]MiddlewareFunc(nil), parent.WithoutMiddleware...),
		child.WithoutMiddleware...,
	)

	for k, v := range parent.Where {
		merged.Where[k] = v
	}

	for k, v := range child.Where {
		merged.Where[k] = v
	}

	return merged
}
