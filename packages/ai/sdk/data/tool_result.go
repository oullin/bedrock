package data

// ToolResult holds the result of a tool invocation returned to the LLM.
// Mirrors Upstream\Ai\Responses\Data\ToolResult.
type ToolResult struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content any    `json:"content"`
}

// ToMap returns a map representation.
func (tr ToolResult) ToMap() map[string]any {
	return map[string]any{
		"id":      tr.ID,
		"name":    tr.Name,
		"content": tr.Content,
	}
}
