package translation

// Loader loads translation messages for a given locale and group.
type Loader interface {
	// Load returns the message map for the given locale/group/namespace triple.
	// A nil namespace is treated as the global namespace ("*").
	// Implementations must return a non-nil map (empty map on miss).
	Load(locale, group string, namespace *string) map[string]any

	// AddNamespace registers a package namespace with a hint path.
	AddNamespace(namespace, hint string)

	// AddJsonPath registers an additional directory to scan for flat JSON
	// locale files (e.g. en.json).
	AddJsonPath(path string)

	// Namespaces returns all registered namespace-to-hint mappings.
	Namespaces() map[string]string
}

// Countable allows custom types to report a collection size for Choice
type Countable interface {
	Len() int
}
