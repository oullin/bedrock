package logtail

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// RenderOptions configures CLI-style rendering of parsed log entries.
type RenderOptions struct {
	Verbose bool
	Columns int
}

// Render formats multiple entries using the same box-drawing layout used by the
// upstream CLI tests.

// RenderEntry formats one entry into a boxed CLI block.

type logtailOrigin struct {
	Type      string
	Command   string
	Method    string
	Path      string
	AuthID    string
	AuthEmail string
	Queue     string
	Job       string
}

const defaultRenderColumns = 50

func Render(entries []Entry, opts RenderOptions) string {
	if len(entries) == 0 {
		return ""
	}

	blocks := make([]string, 0, len(entries))

	for _, entry := range entries {
		blocks = append(blocks, RenderEntry(entry, opts))
	}

	return strings.Join(blocks, "\n")
}

func RenderEntry(entry Entry, opts RenderOptions) string {
	columns := renderColumns(opts.Columns)
	header := renderHeader(entry, opts.Verbose)
	footer := renderFooter(entry)

	lines := make([]string, 0, 2+len(entry.Trace))
	lines = append(lines, frameLine("┌", header, "┐", columns, '─'))

	for _, line := range strings.Split(entry.Message, "\n") {
		lines = append(lines, framedContent(line, columns))
	}

	if opts.Verbose {
		for i, trace := range entry.Trace {
			lines = append(lines, framedContent(fmt.Sprintf("%d. %s", i+1, trace), columns))
		}
	}

	lines = append(lines, renderFooterLines(footer, columns)...)

	return strings.Join(lines, "\n")
}

func renderColumns(requested int) int {
	if requested > 0 {
		return requested
	}

	if raw := strings.TrimSpace(os.Getenv("COLUMNS")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			return parsed
		}
	}

	return defaultRenderColumns
}

func renderHeader(entry Entry, verbose bool) string {
	timePart := entry.Timestamp.Format("15:04:05")

	if verbose && !entry.Timestamp.IsZero() {
		timePart = entry.Timestamp.Format("2006-01-02 15:04:05")
	}

	if class, location, ok := exceptionMetadata(entry.Context); ok && class != "" {
		parts := []string{timePart, class}

		if location != "" {
			parts = append(parts, location)
		}

		return strings.Join(parts, " ")
	}

	level := strings.ToUpper(strings.TrimSpace(entry.Level))

	if level == "" {
		level = "INFO"
	}

	return strings.TrimSpace(timePart + " " + level)
}

func renderFooter(entry Entry) string {
	parts := make([]string, 0, 1+len(entry.Context))

	if summary := originSummary(entry.Context); summary != "" {
		parts = append(parts, summary)
	} else {
		parts = append(parts, "cli eval")
	}

	entries := contextSummaryParts(entry.Context)

	if len(entries) > 0 {
		parts = append(parts, entries...)
	}

	return strings.Join(parts, " • ")
}

func renderFooterLines(footer string, columns int) []string {
	if footer == "" {
		return []string{frameLine("└", "", "┘", columns, '─')}
	}

	parts := strings.Split(footer, "\n")
	lines := make([]string, 0, len(parts))
	lines = append(lines, frameLine("└", parts[0], "┘", columns, '─'))

	for _, part := range parts[1:] {
		lines = append(lines, framedContent(part, columns))
	}

	return lines
}

func frameLine(left, text, right string, columns int, fill rune) string {
	innerWidth := int(math.Max(float64(columns-4), 0))
	text = truncateToWidth(text, innerWidth)
	padding := innerWidth - utf8.RuneCountInString(text)

	if padding < 0 {
		padding = 0
	}

	return left + " " + text + strings.Repeat(string(fill), padding) + " " + right
}

func framedContent(text string, columns int) string {
	innerWidth := int(math.Max(float64(columns-4), 0))
	text = truncateToWidth(text, innerWidth)
	padding := innerWidth - utf8.RuneCountInString(text)

	if padding < 0 {
		padding = 0
	}

	return "│ " + text + strings.Repeat(" ", padding) + " │"
}

func truncateToWidth(text string, width int) string {
	if width <= 0 {
		return ""
	}

	if utf8.RuneCountInString(text) <= width {
		return text
	}

	if width == 1 {
		return "…"
	}

	runes := []rune(text)

	return string(runes[:width-1]) + "…"
}

func exceptionMetadata(ctx map[string]any) (class string, location string, ok bool) {
	if ctx == nil {
		return "", "", false
	}

	raw, ok := ctx["exception"]

	if !ok {
		return "", "", false
	}

	exception, ok := raw.(map[string]any)

	if !ok {
		return "", "", false
	}

	class, _ = exception["class"].(string)
	location = exceptionLocation(exception)

	return class, location, true
}

func exceptionLocation(exception map[string]any) string {
	file, _ := exception["file"].(string)

	if file == "" {
		return ""
	}

	normalized := strings.ReplaceAll(file, "\\", "/")

	if idx := strings.LastIndex(normalized, ":"); idx > strings.LastIndex(normalized, "/") {
		pathPart := normalized[:idx]
		linePart := normalized[idx:]
		pathPart = strings.TrimPrefix(pathPart, string(filepath.Separator))

		if pathPart == "" {
			return normalized
		}

		return pathPart + linePart
	}

	base := filepath.Base(normalized)

	if base == "." || base == string(filepath.Separator) {
		return normalized
	}

	return base
}

func originSummary(ctx map[string]any) string {
	origin := originContext(ctx)

	if origin == nil {
		return ""
	}

	switch strings.ToLower(origin.Type) {
	case "console":
		if origin.Command != "" {
			return "cli " + origin.Command
		}
	case "http":
		method := origin.Method
		path := origin.Path
		auth := "guest"

		if origin.AuthID != "" {
			auth = origin.AuthID

			if origin.AuthEmail != "" {
				auth += " (" + origin.AuthEmail + ")"
			}
		}

		parts := make([]string, 0, 2)

		if method != "" && path != "" {
			parts = append(parts, method+": "+path)
		} else if method != "" {
			parts = append(parts, method)
		} else if path != "" {
			parts = append(parts, path)
		}

		parts = append(parts, "Auth ID: "+auth)

		return strings.Join(parts, " • ")
	case "queue":
		parts := make([]string, 0, 3)

		if origin.Command != "" {
			parts = append(parts, origin.Command)
		}

		if origin.Queue != "" {
			parts = append(parts, origin.Queue)
		}

		if origin.Job != "" {
			parts = append(parts, origin.Job)
		}

		return strings.Join(parts, " ")
	}

	return ""
}

func traceFromContext(ctx map[string]any) []string {
	if ctx == nil {
		return nil
	}

	raw, ok := ctx["trace"]

	if !ok {
		return nil
	}

	switch trace := raw.(type) {
	case []string:
		out := make([]string, 0, len(trace))

		for _, item := range trace {
			if item != "" {
				out = append(out, item)
			}
		}

		return out
	case []any:
		out := make([]string, 0, len(trace))

		for _, item := range trace {
			if text := strings.TrimSpace(formatScalar(item)); text != "" {
				out = append(out, text)
			}
		}

		return out
	case string:
		if trace == "" {
			return nil
		}

		return []string{trace}
	}

	rv := reflect.ValueOf(raw)

	if !rv.IsValid() {
		return nil
	}

	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil
	}

	out := make([]string, 0, rv.Len())

	for i := 0; i < rv.Len(); i++ {
		if text := strings.TrimSpace(formatScalar(rv.Index(i).Interface())); text != "" {
			out = append(out, text)
		}
	}

	return out
}

func originContext(ctx map[string]any) *logtailOrigin {
	if ctx == nil {
		return nil
	}

	raw, ok := ctx["__logtail"]

	if !ok {
		return nil
	}

	data, ok := raw.(map[string]any)

	if !ok {
		return nil
	}

	rawOrigin, ok := data["origin"]

	if !ok {
		return nil
	}

	originMap, ok := rawOrigin.(map[string]any)

	if !ok {
		return nil
	}

	origin := &logtailOrigin{}
	origin.Type, _ = originMap["type"].(string)
	origin.Command, _ = originMap["command"].(string)
	origin.Method, _ = originMap["method"].(string)
	origin.Path, _ = originMap["path"].(string)
	origin.Queue, _ = originMap["queue"].(string)
	origin.Job, _ = originMap["job"].(string)

	if authID, ok := originMap["auth_id"]; ok && authID != nil {
		origin.AuthID = formatScalar(authID)
	}

	origin.AuthEmail, _ = originMap["auth_email"].(string)

	return origin
}

func contextSummaryParts(ctx map[string]any) []string {
	if len(ctx) == 0 {
		return nil
	}

	keys := make([]string, 0, len(ctx))

	for key := range ctx {
		if key == "__logtail" || key == "exception" || key == "trace" {
			continue
		}

		keys = append(keys, key)
	}

	sort.Strings(keys)

	parts := make([]string, 0, len(keys))

	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s: %s", key, formatValue(ctx[key])))
	}

	return parts
}

func formatValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case string:
		return quoteScalar(v)
	case []string:
		items := make([]string, 0, len(v))

		for i, item := range v {
			items = append(items, fmt.Sprintf("%d => %s", i, quoteScalar(item)))
		}

		return "array ( " + strings.Join(items, ", ") + ", )"
	case []any:
		items := make([]string, 0, len(v))

		for i, item := range v {
			items = append(items, fmt.Sprintf("%d => %s", i, formatValue(item)))
		}

		return "array ( " + strings.Join(items, ", ") + ", )"
	case map[string]any:
		keys := make([]string, 0, len(v))

		for key := range v {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		items := make([]string, 0, len(keys))

		for _, key := range keys {
			items = append(items, fmt.Sprintf("%s => %s", quoteScalar(key), formatValue(v[key])))
		}

		return "array ( " + strings.Join(items, ", ") + ", )"
	}

	rv := reflect.ValueOf(value)

	if !rv.IsValid() {
		return "null"
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		items := make([]string, 0, rv.Len())

		for i := 0; i < rv.Len(); i++ {
			items = append(items, fmt.Sprintf("%d => %s", i, formatValue(rv.Index(i).Interface())))
		}

		return "array ( " + strings.Join(items, ", ") + ", )"
	case reflect.Map:
		keys := rv.MapKeys()

		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface())
		})

		items := make([]string, 0, len(keys))

		for _, key := range keys {
			items = append(items, fmt.Sprintf("%s => %s", quoteScalar(fmt.Sprint(key.Interface())), formatValue(rv.MapIndex(key).Interface())))
		}

		return "array ( " + strings.Join(items, ", ") + ", )"
	}

	return fmt.Sprint(value)
}

func formatScalar(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.Itoa(int(v))
	case int16:
		return strconv.Itoa(int(v))
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprint(value)
	}
}

func quoteScalar(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "\\'") + "'"
}
