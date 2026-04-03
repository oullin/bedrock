package authflows

// RoutePath resolves configurable AuthFlows route paths.
type RoutePath struct {
	paths map[string]string
}

// NewRoutePath creates a route path resolver.
func NewRoutePath(paths map[string]string) RoutePath {
	cloned := make(map[string]string, len(paths))

	for key, value := range paths {
		cloned[key] = value
	}

	return RoutePath{paths: cloned}
}

// For returns the configured path or the supplied default.
func (r RoutePath) For(name string, defaultPath string) string {
	if path, ok := r.paths[name]; ok && path != "" {
		return path
	}

	return defaultPath
}
