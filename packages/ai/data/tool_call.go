package data

// ToolCall represents a tool invocation requested by the LLM.
// Mirrors Laravel\Ai\Responses\Data\ToolCall.
type ToolCall struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Arguments        map[string]any   `json:"arguments"`
	ResultID         *string          `json:"result_id,omitempty"`
	ReasoningID      *string          `json:"reasoning_id,omitempty"`
	ReasoningSummary []map[string]any `json:"reasoning_summary,omitempty"`
}

// ToMap returns a map representation.
func (tc ToolCall) ToMap() map[string]any {
	return map[string]any{
		"id":                tc.ID,
		"name":              tc.Name,
		"arguments":         tc.Arguments,
		"result_id":         tc.ResultID,
		"reasoning_id":      tc.ReasoningID,
		"reasoning_summary": tc.ReasoningSummary,
	}
}
