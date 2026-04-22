package testing

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// CompiledViewPath appends the parallel testing token to a compiled view path.
func CompiledViewPath(base, token string) string {
	base = strings.TrimRight(base, "/")
	if base == "" || token == "" {
		return ""
	}

	return base + "/" + token
}

// SwitchToCompiledViewPath returns a copy of the config with the compiled view
// path updated under config["view"]["compiled"] when present.
func SwitchToCompiledViewPath(config map[string]any, token string) map[string]any {
	if config == nil {
		return nil
	}

	clone := cloneValue(config).(map[string]any)
	view, ok := clone["view"].(map[string]any)
	if !ok {
		return clone
	}

	compiled, ok := view["compiled"].(string)
	if !ok {
		return clone
	}

	view["compiled"] = CompiledViewPath(compiled, token)
	return clone
}

// CachePrefix appends the parallel testing token to an existing cache prefix.
func CachePrefix(prefix, token string) string {
	if prefix == "" {
		return token
	}

	if token == "" {
		return prefix
	}

	return prefix + token
}

// SwitchToCachePrefix returns a copy of the config with config["cache"]["prefix"]
// updated when present.
func SwitchToCachePrefix(config map[string]any, token string) map[string]any {
	if config == nil {
		return nil
	}

	clone := cloneValue(config).(map[string]any)
	cache, ok := clone["cache"].(map[string]any)
	if !ok {
		return clone
	}

	prefix, ok := cache["prefix"].(string)
	if !ok {
		return clone
	}

	cache["prefix"] = CachePrefix(prefix, token)
	return clone
}

// DeleteCompiledViewDirectory removes a compiled view directory created during
// test bootstrapping.
func DeleteCompiledViewDirectory(path string) error {
	return os.RemoveAll(path)
}

// SwitchToDatabase returns a copy of the config with the database URL updated.
func SwitchToDatabase(config map[string]any, url string) map[string]any {
	if config == nil {
		return nil
	}

	clone := cloneValue(config).(map[string]any)
	database, ok := clone["database"].(map[string]any)
	if !ok {
		return clone
	}

	database["url"] = url
	return clone
}

// CastToJSONType returns the JSON column type used by a supported driver.
//
// SQL Server is intentionally unsupported in Bedrock, so callers can treat a
// false result as a divergence candidate instead of a portable equivalent.
func CastToJSONType(driver string) (string, bool) {
	switch strings.ToLower(driver) {
	case "mysql", "mariadb":
		return "json", true
	case "postgres", "postgresql":
		return "jsonb", true
	case "sqlite":
		return "text", true
	default:
		return "", false
	}
}

// BootTestCache registers a setup callback unless the caller has opted out of
// cache isolation.
func BootTestCache(state *ParallelTestingState, optOut bool, callback func()) {
	if state == nil || optOut || callback == nil {
		return
	}

	state.RegisterCallback(callback)
}

// ResolveConfigValue looks up a dotted config path using the same lookup rules
// as the response assertions.
func ResolveConfigValue(config map[string]any, key string) (any, bool) {
	return lookupMap(config, key)
}

// RenderConfigShow renders a config value or tree in a deterministic order.
func RenderConfigShow(config map[string]any, key string) (string, error) {
	value, ok := lookupMap(config, key)
	if !ok {
		if key == "" {
			value = config
		} else {
			return "", fmt.Errorf("config key %q does not exist", key)
		}
	}

	lines := make([]string, 0)
	flattenConfig(key, value, &lines)
	sort.Strings(lines)

	if len(lines) == 0 {
		return fmt.Sprint(value), nil
	}

	return strings.Join(lines, "\n"), nil
}

// DeprecationReporter records warnings only when enabled.
type DeprecationReporter struct {
	enabled  bool
	messages []string
}

// NewDeprecationReporter constructs a reporter that may suppress warnings.
func NewDeprecationReporter(enabled bool) *DeprecationReporter {
	return &DeprecationReporter{enabled: enabled}
}

// Warn stores the warning when reporting is enabled.
func (r *DeprecationReporter) Warn(message string) {
	if !r.enabled {
		return
	}

	r.messages = append(r.messages, message)
}

// Messages returns a copy of the recorded warnings.
func (r *DeprecationReporter) Messages() []string {
	return append([]string(nil), r.messages...)
}

// ParallelTestingState captures the small state needed for the parallel
// testing surface covered by the inventory.
type ParallelTestingState struct {
	token     string
	options   map[string]string
	callbacks []func()
}

// NewParallelTestingState constructs an isolated state container.
func NewParallelTestingState(token string) *ParallelTestingState {
	return &ParallelTestingState{
		token:   token,
		options: map[string]string{},
	}
}

// Token returns the active parallel testing token.
func (p *ParallelTestingState) Token() string {
	return p.token
}

// SetOption stores a parallel testing option.
func (p *ParallelTestingState) SetOption(key, value string) {
	p.options[key] = value
}

// Options returns a copy of the registered parallel testing options.
func (p *ParallelTestingState) Options() map[string]string {
	out := make(map[string]string, len(p.options))
	for key, value := range p.options {
		out[key] = value
	}

	return out
}

// RegisterCallback adds a callback to the current process state.
func (p *ParallelTestingState) RegisterCallback(callback func()) {
	p.callbacks = append(p.callbacks, callback)
}

// RunCallbacks executes every registered callback in registration order.
func (p *ParallelTestingState) RunCallbacks() {
	for _, callback := range p.callbacks {
		callback()
	}
}

func flattenConfig(prefix string, value any, out *[]string) {
	switch current := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(current))
		for key := range current {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			next := key
			if prefix != "" {
				next = prefix + "." + key
			}

			flattenConfig(next, current[key], out)
		}
	case []any:
		for i, item := range current {
			next := fmt.Sprintf("%s.%d", prefix, i)
			if prefix == "" {
				next = fmt.Sprintf("%d", i)
			}
			flattenConfig(next, item, out)
		}
	default:
		*out = append(*out, fmt.Sprintf("%s=%v", prefix, current))
	}
}

func cloneValue(value any) any {
	switch current := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(current))
		for key, item := range current {
			out[key] = cloneValue(item)
		}
		return out
	case []any:
		out := make([]any, len(current))
		for i, item := range current {
			out[i] = cloneValue(item)
		}
		return out
	default:
		return value
	}
}
