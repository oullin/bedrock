package mcp

import "strings"

// maxCompletionValues is the maximum number of values returned in a single
// completion response. Upstream caps this at 100.

// CompletionResult holds the completion suggestions for a prompt or resource
// argument. It mirrors Upstream's CompletionResponse.
type CompletionResult struct {
	Values  []string
	HasMore bool
	Total   int
}

const maxCompletionValues = 100

// toMap serialises the result for the completion/complete response.
func (r *CompletionResult) toMap() map[string]any {
	values := r.Values

	if values == nil {
		values = []string{}
	}

	return map[string]any{
		"values":  values,
		"total":   r.Total,
		"hasMore": r.HasMore,
	}
}

// EmptyCompletion returns a CompletionResult with no suggestions.
func EmptyCompletion() *CompletionResult {
	return &CompletionResult{Values: []string{}}
}

// EnumCompletion returns the full enum list, truncated to maxCompletionValues.
func EnumCompletion(values []string) *CompletionResult {
	total := len(values)
	hasMore := false

	if total > maxCompletionValues {
		values = values[:maxCompletionValues]
		hasMore = true
	}

	return &CompletionResult{Values: values, Total: total, HasMore: hasMore}
}

// MatchCompletion filters values whose prefix matches input (case-insensitive)
// and returns up to maxCompletionValues results.
func MatchCompletion(values []string, input string) *CompletionResult {
	lower := strings.ToLower(input)
	matched := make([]string, 0, len(values))

	for _, v := range values {
		if strings.HasPrefix(strings.ToLower(v), lower) {
			matched = append(matched, v)
		}
	}

	total := len(matched)
	hasMore := false

	if total > maxCompletionValues {
		matched = matched[:maxCompletionValues]
		hasMore = true
	}

	return &CompletionResult{Values: matched, Total: total, HasMore: hasMore}
}
