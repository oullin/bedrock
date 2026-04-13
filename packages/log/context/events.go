package context

// ContextDehydrating is dispatched before context data is serialized.
type ContextDehydrating struct {
	Context *Repository
}

// ContextHydrated is dispatched after context data has been restored.
type ContextHydrated struct {
	Context *Repository
}
