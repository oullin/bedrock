package watchers

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bedrock/packages/telescope"
)

// milliseconds

// QueryWatcher monitors database query execution and records entries with SQL,
// bindings, execution time, and slow-query tagging. It mirrors the upstream // QueryWatcher class.
//
// Options:
//   - "slow" (float64): threshold in milliseconds above which a query is
//     tagged as slow (default 100 ms).
type QueryWatcher struct {
	telescope.BaseWatcher
}

const defaultSlowQueryThreshold = 100.0

// NewQueryWatcher creates a QueryWatcher with the given options.
func NewQueryWatcher(t *telescope.Telescope, options map[string]any) *QueryWatcher {
	w := &QueryWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

// Register is a no-op for QueryWatcher; callers drive it by calling Record
// directly. When integrated with a database driver, call Record after each
// query execution.
func (w *QueryWatcher) Register(_ any) error { return nil }

// slowThreshold returns the configured slow-query threshold in milliseconds.
func (w *QueryWatcher) slowThreshold() float64 {
	t := w.Float64Option("slow")

	if t <= 0 {
		return defaultSlowQueryThreshold
	}

	return t
}

// Record records a database query entry. sql is the raw SQL statement,
// bindings contains the parameter values (positional or named), duration is
// the execution time, and connection is the database connection name.
// callerFile and callerLine identify the call site (optional; pass "" and 0
// to omit).
func (w *QueryWatcher) Record(sql string, bindings []any, duration time.Duration, connection string, callerFile string, callerLine int) {
	durationMs := float64(duration.Microseconds()) / 1000.0
	slow := durationMs >= w.slowThreshold()

	stmt := ReplaceBindings(sql, bindings)

	// Family hash groups identical queries together (keyed on raw SQL so
	// different bindings for the same statement share the same hash).
	hash := fmt.Sprintf("%x", md5.Sum([]byte(sql)))

	content := map[string]any{
		"connection": connection,
		"sql":        stmt,
		"raw_sql":    sql,
		"bindings":   bindings,
		"time":       durationMs,
		"slow":       slow,
	}

	if callerFile != "" {
		content["file"] = callerFile
		content["line"] = callerLine
	}

	entry := telescope.NewEntry(telescope.EntryTypeQuery, content).
		WithFamilyHash(hash)

	if slow {
		entry.AddTags("slow")
	}

	w.Scope().RecordQuery(entry)
}

// ─── Binding replacement ─────────────────────────────────────────────────────

// positionalPlaceholder matches a single ? in SQL.
var positionalPlaceholder = regexp.MustCompile(`\?`)

// namedPlaceholder matches :name style placeholders in SQL.
var namedPlaceholder = regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9_]*)`)

// ReplaceBindings substitutes SQL parameter placeholders with their actual
// values, mirroring QueryWatcher::replaceBindings().
//
// It supports both positional (?) and named (:key) parameter styles.
// String values are single-quoted; nil is rendered as NULL.
func ReplaceBindings(sql string, bindings []any) string {
	if len(bindings) == 0 {
		return sql
	}

	// Detect style: if the first binding is a map key we have named params.
	// In practice Go drivers use positional params, but we support both.
	sql = replacePositional(sql, bindings)

	return sql
}

// replacePositional replaces ? placeholders with binding values.
func replacePositional(sql string, bindings []any) string {
	idx := 0

	return positionalPlaceholder.ReplaceAllStringFunc(sql, func(_ string) string {
		if idx >= len(bindings) {
			idx++

			return "?"
		}

		v := bindings[idx]
		idx++

		return quoteBinding(v)
	})
}

// ReplaceNamedBindings replaces :key placeholders with values from a map.
// Exported so tests can call it directly.
func ReplaceNamedBindings(sql string, bindings map[string]any) string {
	return namedPlaceholder.ReplaceAllStringFunc(sql, func(match string) string {
		key := strings.TrimPrefix(match, ":")

		v, ok := bindings[key]

		if !ok {
			return match
		}

		return quoteBinding(v)
	})
}

// quoteBinding formats a single binding value for inline SQL substitution,
// mirroring QueryWatcher::quoteStringBinding().
func quoteBinding(v any) string {
	if v == nil {
		return "NULL"
	}

	switch val := v.(type) {
	case bool:
		if val {
			return "1"
		}

		return "0"
	case string:
		// Escape single quotes by doubling them.
		escaped := strings.ReplaceAll(val, "'", "''")

		return "'" + escaped + "'"
	case []byte:
		escaped := strings.ReplaceAll(string(val), "'", "''")

		return "'" + escaped + "'"
	case time.Time:
		return "'" + val.Format("2006-01-02 15:04:05") + "'"
	default:
		return fmt.Sprintf("%v", val)
	}
}
