package logtail

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
	"time"
)

// Entry is a parsed log line.
type Entry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Context   map[string]any
	Trace     []string
	Raw       string
}

// Filter selects entries while tailing.
type Filter struct {
	Levels   []string
	Contains string
	Since    time.Time
	Until    time.Time
}

const laravelTimestampLayout = "2006-01-02 15:04:05"

// Match reports whether entry satisfies the filter.
func (f Filter) Match(entry Entry) bool {
	if len(f.Levels) > 0 {
		matched := false

		for _, level := range f.Levels {
			if strings.EqualFold(entry.Level, level) {
				matched = true

				break
			}
		}

		if !matched {
			return false
		}
	}

	if f.Contains != "" {
		needle := strings.ToLower(f.Contains)

		if !strings.Contains(strings.ToLower(entry.Message), needle) &&
			!strings.Contains(strings.ToLower(entry.Level), needle) &&
			!strings.Contains(strings.ToLower(entry.Raw), needle) {
			return false
		}
	}

	if !f.Since.IsZero() && entry.Timestamp.Before(f.Since) {
		return false
	}

	if !f.Until.IsZero() && entry.Timestamp.After(f.Until) {
		return false
	}

	return true
}

// ParseLine parses a Upstream-style log line.
func ParseLine(line string) Entry {
	entry := Entry{
		Message: line,
		Raw:     line,
	}

	if len(line) < len("[2006-01-02 15:04:05]") || line[0] != '[' {
		return entry
	}

	end := strings.IndexByte(line, ']')

	if end == -1 {
		return entry
	}

	if timestamp, err := time.ParseInLocation(laravelTimestampLayout, line[1:end], time.Local); err == nil {
		entry.Timestamp = timestamp
	}

	rest := strings.TrimSpace(line[end+1:])
	colon := strings.Index(rest, ":")

	if colon == -1 {
		entry.Message = rest

		return entry
	}

	channelLevel := strings.TrimSpace(rest[:colon])
	message := strings.TrimSpace(rest[colon+1:])

	if dot := strings.LastIndex(channelLevel, "."); dot >= 0 && dot < len(channelLevel)-1 {
		entry.Level = strings.ToLower(channelLevel[dot+1:])
	} else {
		entry.Level = strings.ToLower(channelLevel)
	}

	entry.Message = message

	if contextStart := strings.LastIndex(message, " {"); contextStart >= 0 {
		rawContext := strings.TrimSpace(message[contextStart+1:])
		context := map[string]any{}

		if err := json.Unmarshal([]byte(rawContext), &context); err == nil {
			entry.Message = strings.TrimSpace(message[:contextStart])
			entry.Context = context
			entry.Trace = traceFromContext(context)
		}
	}

	return entry
}

// Collect reads log lines and returns entries that match filter.
func Collect(r io.Reader, filter Filter) ([]Entry, error) {
	scanner := bufio.NewScanner(r)
	entries := make([]Entry, 0)

	for scanner.Scan() {
		entry := ParseLine(scanner.Text())

		if filter.Match(entry) {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
