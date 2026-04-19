// Package prompts contains the request types sent through the AI pipeline.
package prompts

import "github.com/bedrock/packages/ai/messages"

// AgentPrompt carries all parameters for a single agent invocation.
// It is the "passable" value sent through the middleware pipeline.
type AgentPrompt struct {
	InvocationID string
	Text         string
	Attachments  []messages.Attachment
	Provider     *string // explicit provider override
	Model        *string // explicit model override
	Timeout      int     // seconds; 0 = provider default (60)
}

// WithProvider returns a copy of the prompt with a provider override.
func (p *AgentPrompt) WithProvider(provider string) *AgentPrompt {
	cp := *p
	cp.Provider = &provider

	return &cp
}

// WithModel returns a copy of the prompt with a model override.
func (p *AgentPrompt) WithModel(model string) *AgentPrompt {
	cp := *p
	cp.Model = &model

	return &cp
}

// WithTimeout returns a copy of the prompt with a timeout override.
func (p *AgentPrompt) WithTimeout(secs int) *AgentPrompt {
	cp := *p
	cp.Timeout = secs

	return &cp
}
