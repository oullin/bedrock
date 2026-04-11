package httpx

import (
	"strings"
)

// IsPrecognitive returns true when the request carries a Precognition header
// with a truthy value.
func (r *Request) IsPrecognitive() bool {
	return isTruthy(r.raw.Header.Get("Precognition"))
}

// IsAttemptingPrecognition returns true when the request is precognitive and
// also carries a Precognition-Validate-Only header listing fields to validate.
func (r *Request) IsAttemptingPrecognition() bool {
	return r.IsPrecognitive() && r.raw.Header.Get("Precognition-Validate-Only") != ""
}

// PrecognitiveValidateOnly returns the comma-separated list of fields that the
// client wants validated without a full submission.
func (r *Request) PrecognitiveValidateOnly() []string {
	header := r.raw.Header.Get("Precognition-Validate-Only")

	if header == "" {
		return nil
	}

	parts := strings.Split(header, ",")
	fields := make([]string, 0, len(parts))

	for _, p := range parts {
		if f := strings.TrimSpace(p); f != "" {
			fields = append(fields, f)
		}
	}

	return fields
}

// FilterPrecognitiveRules filters a set of validation rules down to only those
// fields listed in the Precognition-Validate-Only header. If the request is not
// precognitive the full rule set is returned unchanged.
func (r *Request) FilterPrecognitiveRules(rules map[string]any) map[string]any {
	if !r.IsAttemptingPrecognition() {
		return rules
	}

	fields := r.PrecognitiveValidateOnly()
	allowed := make(map[string]struct{}, len(fields))

	for _, f := range fields {
		allowed[f] = struct{}{}
	}

	filtered := make(map[string]any, len(fields))

	for k, v := range rules {
		// Match the field name or any nested field (e.g., "address" matches "address.street").
		base := k

		if idx := strings.Index(k, "."); idx != -1 {
			base = k[:idx]
		}

		if _, ok := allowed[base]; ok {
			filtered[k] = v
		}
	}

	return filtered
}

func isTruthy(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "on", "yes":
		return true
	}

	return false
}
