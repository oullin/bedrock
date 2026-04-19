package debugbar

// Watcher is implemented by all DebugBar monitoring components. Each watcher
// hooks into one aspect of the application (HTTP requests, DB queries, events,
// etc.) and forwards captured data to a DebugBar instance.
//
// Register is called once during application boot with the application
// container so the watcher can subscribe to the events or middleware it needs.
type Watcher interface {
	Register(app any) error
}
