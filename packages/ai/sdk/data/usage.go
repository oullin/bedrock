// Package data contains pure value types used throughout the AI package responses.
package data

// Usage tracks token consumption for an AI request.
type Usage struct {
	PromptTokens          int `json:"prompt_tokens"`
	CompletionTokens      int `json:"completion_tokens"`
	CacheWriteInputTokens int `json:"cache_write_input_tokens"`
	CacheReadInputTokens  int `json:"cache_read_input_tokens"`
	ReasoningTokens       int `json:"reasoning_tokens"`
}

// Add sums two Usage values and returns a new one.
func (u Usage) Add(other Usage) Usage {
	return Usage{
		PromptTokens:          u.PromptTokens + other.PromptTokens,
		CompletionTokens:      u.CompletionTokens + other.CompletionTokens,
		CacheWriteInputTokens: u.CacheWriteInputTokens + other.CacheWriteInputTokens,
		CacheReadInputTokens:  u.CacheReadInputTokens + other.CacheReadInputTokens,
		ReasoningTokens:       u.ReasoningTokens + other.ReasoningTokens,
	}
}

// ToMap returns a map representation compatible with JSON serialisation.
func (u Usage) ToMap() map[string]any {
	return map[string]any{
		"prompt_tokens":            u.PromptTokens,
		"completion_tokens":        u.CompletionTokens,
		"cache_write_input_tokens": u.CacheWriteInputTokens,
		"cache_read_input_tokens":  u.CacheReadInputTokens,
		"reasoning_tokens":         u.ReasoningTokens,
	}
}
