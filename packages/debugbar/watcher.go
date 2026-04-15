package debugbar

// Watcher is implemented by all DebugBar monitoring components. Each watcher
// subscribes to one aspect of the application (HTTP, DB, events, …) and
// forwards captured data to a DebugBar instance via its Record* methods.
//
// Register is called once during application boot. The app parameter is typed
// as any to avoid importing the container package; concrete watchers may
// perform a type assertion to *container.Container or accept a lighter-weight
// application interface.
type Watcher interface {
	Register(app any) error
}

// BaseWatcher is an embeddable struct that gives concrete watchers access to
// the DebugBar instance and their options map. It mirrors the PHP base
// Watcher class with its $options property.
type BaseWatcher struct {
	debugbar *DebugBar
	Options   map[string]any
}

// SetDebugBar attaches the DebugBar instance to the watcher.
func (b *BaseWatcher) SetDebugBar(t *DebugBar) {
	b.debugbar = t
}

// Scope returns the attached DebugBar instance.
func (b *BaseWatcher) Scope() *DebugBar {
	return b.debugbar
}

// Option returns a watcher option by key, or nil if absent.
func (b *BaseWatcher) Option(key string) any {
	if b.Options == nil {
		return nil
	}

	return b.Options[key]
}

// BoolOption returns a boolean watcher option, defaulting to false.
func (b *BaseWatcher) BoolOption(key string) bool {
	v, _ := b.Options[key].(bool)
	return v
}

// StringOption returns a string watcher option, defaulting to "".
func (b *BaseWatcher) StringOption(key string) string {
	v, _ := b.Options[key].(string)
	return v
}

// Float64Option returns a float64 watcher option, defaulting to 0.
func (b *BaseWatcher) Float64Option(key string) float64 {
	switch v := b.Options[key].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

// StringsOption returns a []string watcher option, defaulting to nil.
func (b *BaseWatcher) StringsOption(key string) []string {
	v, _ := b.Options[key].([]string)
	return v
}
