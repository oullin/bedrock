package ai

// EstimateTokens returns a coarse OpenAI-style token estimate (4 chars/tok).
// Suitable for budgeting decisions where exactness is not required. A proper
// tokenizer can land later behind the same signature.
func EstimateTokens(s string) int {
	if s == "" {
		return 0
	}

	return (len(s) + 3) / 4
}
