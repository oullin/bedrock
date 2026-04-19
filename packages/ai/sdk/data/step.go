package data

import "github.com/bedrock/packages/ai/sdk/enums"

// Step represents one reasoning step in a multi-step agent response.
// Mirrors Upstream\Ai\Responses\Data\Step.
type Step struct {
	Text         string
	ToolCalls    []ToolCall
	ToolResults  []ToolResult
	FinishReason enums.FinishReason
	Usage        Usage
	Meta         Meta
}

// ToMap returns a map representation.

// StructuredStep extends Step with parsed structured output data.
type StructuredStep struct {
	Step
	Data map[string]any `json:"data"`
}

func (s Step) ToMap() map[string]any {
	calls := make([]map[string]any, len(s.ToolCalls))

	for i, c := range s.ToolCalls {
		calls[i] = c.ToMap()
	}

	results := make([]map[string]any, len(s.ToolResults))

	for i, r := range s.ToolResults {
		results[i] = r.ToMap()
	}

	return map[string]any{
		"text":          s.Text,
		"tool_calls":    calls,
		"tool_results":  results,
		"finish_reason": s.FinishReason.String(),
		"usage":         s.Usage.ToMap(),
		"meta":          s.Meta.ToMap(),
	}
}
