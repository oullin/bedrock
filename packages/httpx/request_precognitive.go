package httpx

import (
	"regexp"
	"strings"
)

// IsPrecognitive returns true when the request carries a HTTPPreview header
// with a truthy value.
func (r *Request) IsPrecognitive() bool {
	return isTruthy(r.raw.Header.Get("HTTPPreview"))
}

// IsAttemptingHTTPPreview returns true when the request carries a
// HTTPPreview header with the exact value "true". This checks the client's
// intent to make a precognitive request.
//
// Ref: @bedrock/code-0222
func (r *Request) IsAttemptingHTTPPreview() bool {
	return r.raw.Header.Get("HTTPPreview") == "true"
}

// PrecognitiveValidateOnly returns the comma-separated list of fields that the
// client wants validated without a full submission.
func (r *Request) PrecognitiveValidateOnly() []string {
	header := r.raw.Header.Get("HTTPPreview-Validate-Only")

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
// fields listed in the HTTPPreview-Validate-Only header. If the header is not
// present the full rule set is returned unchanged.
//
// Patterns support wildcards: "address.*" matches "address.street" and
// "address.city" but not "address.street.line".
//
// Ref: @bedrock/code-0222
func (r *Request) FilterPrecognitiveRules(rules map[string]any) map[string]any {
	if r.raw.Header.Get("HTTPPreview-Validate-Only") == "" {
		return rules
	}

	validateOnly := r.PrecognitiveValidateOnly()
	filtered := make(map[string]any)

	for attr, v := range rules {
		if shouldValidatePrecognitiveAttribute(attr, validateOnly) {
			filtered[attr] = v
		}
	}

	return filtered
}

// shouldValidatePrecognitiveAttribute reports whether the given attribute
// should be validated based on the HTTPPreview-Validate-Only patterns.
// Each pattern is converted to a regex where * is replaced with [^.]+ to
// match a single dot-separated segment.
//
// Ref: @bedrock/code-0222
func shouldValidatePrecognitiveAttribute(attribute string, validateOnly []string) bool {
	for _, pattern := range validateOnly {
		escaped := regexp.QuoteMeta(pattern)
		escaped = strings.ReplaceAll(escaped, `\*`, `[^.]+`)
		re, err := regexp.Compile(`^` + escaped + `$`)

		if err != nil {
			continue
		}

		if re.MatchString(attribute) {
			return true
		}
	}

	return false
}

func isTruthy(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "on", "yes":
		return true
	}

	return false
}
