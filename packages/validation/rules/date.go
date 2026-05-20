package rules

import (
	"strings"
	"time"
)

// Common date layouts that upstream recognises (subset).
var commonDateLayouts = []string{
	"2006-01-02",
	"2006-01-02 15:04:05",
	"02/01/2006",
	"01/02/2006",
	"2006-01-02T15:04:05Z07:00",
	time.RFC3339,
	"January 2, 2006",
	"Jan 2, 2006",
	"2 January 2006",
}

func init_date() {
	Register("Date", validateDate)
	Register("DateFormat", validateDateFormat)
	Register("DateEquals", validateDateEquals)
	Register("Before", validateBefore)
	Register("BeforeOrEqual", validateBeforeOrEqual)
	Register("After", validateAfter)
	Register("AfterOrEqual", validateAfterOrEqual)
	Register("Timezone", validateTimezone)
}

func parseDate(value any) (time.Time, bool) {
	s, ok := value.(string)

	if !ok {
		return time.Time{}, false
	}

	s = strings.TrimSpace(s)

	switch strings.ToLower(s) {
	case "today":
		now := time.Now().UTC()

		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), true
	case "tomorrow":
		now := time.Now().UTC().AddDate(0, 0, 1)

		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), true
	case "yesterday":
		now := time.Now().UTC().AddDate(0, 0, -1)

		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), true
	}

	for _, layout := range commonDateLayouts {
		t, err := time.Parse(layout, s)

		if err == nil {
			return t, true
		}
	}

	return time.Time{}, false
}

func parseDateWithFormat(s, format string) (time.Time, bool) {
	// Convert PHP date format to Go layout
	layout := phpFormatToGo(format)
	t, err := time.Parse(layout, strings.TrimSpace(s))

	return t, err == nil
}

// phpFormatToGo converts a PHP date format string to a Go time layout.
// Only a practical subset is supported.
func phpFormatToGo(format string) string {
	replacer := strings.NewReplacer(
		"Y", "2006",
		"y", "06",
		"m", "01",
		"d", "02",
		"H", "15",
		"i", "04",
		"s", "05",
		"n", "1",
		"j", "2",
		"g", "3",
		"A", "PM",
		"a", "pm",
	)

	return replacer.Replace(format)
}

func validateDate(_ string, value any, _ []string, _ RuleContext) bool {
	_, ok := parseDate(value)

	return ok
}

func validateDateFormat(_ string, value any, params []string, _ RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	s, ok := value.(string)

	if !ok {
		return false
	}

	_, valid := parseDateWithFormat(s, params[0])

	return valid
}

// resolveDate resolves a value that might be a field reference or a literal date.
func resolveDate(param string, ctx RuleContext) (time.Time, bool) {
	// Check if param is a field name
	other := ctx.GetValue(param)

	if other != nil {
		return parseDate(other)
	}

	// Try as a literal date string
	return parseDate(param)
}

func validateDateEquals(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	t1, ok1 := parseDate(value)

	if !ok1 {
		return false
	}

	t2, ok2 := resolveDate(params[0], ctx)

	if !ok2 {
		return false
	}

	return t1.Equal(t2)
}

func validateBefore(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	t1, ok1 := parseDate(value)

	if !ok1 {
		return false
	}

	t2, ok2 := resolveDate(params[0], ctx)

	if !ok2 {
		return false
	}

	return t1.Before(t2)
}

func validateBeforeOrEqual(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	t1, ok1 := parseDate(value)

	if !ok1 {
		return false
	}

	t2, ok2 := resolveDate(params[0], ctx)

	if !ok2 {
		return false
	}

	return !t1.After(t2)
}

func validateAfter(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	t1, ok1 := parseDate(value)

	if !ok1 {
		return false
	}

	t2, ok2 := resolveDate(params[0], ctx)

	if !ok2 {
		return false
	}

	return t1.After(t2)
}

func validateAfterOrEqual(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	t1, ok1 := parseDate(value)

	if !ok1 {
		return false
	}

	t2, ok2 := resolveDate(params[0], ctx)

	if !ok2 {
		return false
	}

	return !t1.Before(t2)
}

// validateTimezone checks if value is a valid IANA timezone identifier.
func validateTimezone(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	_, err := time.LoadLocation(strings.TrimSpace(s))

	return err == nil
}
