package translation

// ArrayLoader is an in-memory Loader implementation.  Messages are stored as
// messages[namespace][locale][group] → map[string]any, mirroring Laravel's
// Illuminate\Translation\ArrayLoader.
type ArrayLoader struct {
	messages   map[string]map[string]map[string]map[string]any
	namespaces map[string]string
}

var _ Loader = (*ArrayLoader)(nil)

// NewArrayLoader returns an initialised *ArrayLoader.
func NewArrayLoader() *ArrayLoader {
	return &ArrayLoader{
		messages:   make(map[string]map[string]map[string]map[string]any),
		namespaces: make(map[string]string),
	}
}

// Load returns the message map for locale/group/namespace.
// A nil namespace is treated as the global namespace ("*").
// Always returns a non-nil map.
func (l *ArrayLoader) Load(locale, group string, namespace *string) map[string]any {
	ns := globalNS

	if namespace != nil && *namespace != "" {
		ns = *namespace
	}

	if byLocale, ok := l.messages[ns]; ok {
		if byGroup, ok := byLocale[locale]; ok {
			if msgs, ok := byGroup[group]; ok {
				return msgs
			}
		}
	}

	return map[string]any{}
}

// AddMessages stores messages for the given locale/group/namespace.
// A nil namespace defaults to "*".  Returns the loader for fluent chaining.
func (l *ArrayLoader) AddMessages(locale, group string, messages map[string]any, namespace *string) *ArrayLoader {
	ns := globalNS

	if namespace != nil && *namespace != "" {
		ns = *namespace
	}

	if l.messages[ns] == nil {
		l.messages[ns] = make(map[string]map[string]map[string]any)
	}

	if l.messages[ns][locale] == nil {
		l.messages[ns][locale] = make(map[string]map[string]any)
	}

	existing := l.messages[ns][locale][group]

	if existing == nil {
		l.messages[ns][locale][group] = messages
	} else {
		l.messages[ns][locale][group] = mergeMaps(existing, messages)
	}

	return l
}

// AddNamespace records the namespace in the registry (no-op for lookups).
func (l *ArrayLoader) AddNamespace(namespace, hint string) {
	l.namespaces[namespace] = hint
}

// AddJsonPath is a no-op for ArrayLoader (satisfies Loader interface).
func (l *ArrayLoader) AddJsonPath(_ string) {}

// Namespaces returns a shallow copy of the registered namespace map.
func (l *ArrayLoader) Namespaces() map[string]string {
	cp := make(map[string]string, len(l.namespaces))

	for k, v := range l.namespaces {
		cp[k] = v
	}

	return cp
}

const globalNS = "*"
