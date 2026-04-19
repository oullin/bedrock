package mcp

import "context"

// Role is the MCP message role ("user" or "assistant").
type Role string

// Message is a role/content pair returned by a Prompt.
type Message struct {
	Role    Role
	Content Content
}

// ToMap serialises the message for inclusion in a prompts/get response.

// UserMessage is a convenience constructor for a user-role message.

// AssistantMessage is a convenience constructor for an assistant-role message.

// Prompt is an MCP prompt template that can be invoked with arguments to
// produce a sequence of role/content messages.
type Prompt interface {
	Name() string
	Description() string
	Arguments() []*Argument
	Invoke(ctx context.Context, req *Request) ([]*Message, error)
}

// funcPrompt implements Prompt from plain functions.
type funcPrompt struct {
	name        string
	description string
	arguments   []*Argument
	invoker     func(ctx context.Context, req *Request) ([]*Message, error)
}

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

func (m *Message) ToMap() map[string]any {
	return map[string]any{
		"role":    string(m.Role),
		"content": m.Content.ToPrompt(),
	}
}

func UserMessage(content Content) *Message {
	return &Message{Role: RoleUser, Content: content}
}

func AssistantMessage(content Content) *Message {
	return &Message{Role: RoleAssistant, Content: content}
}

func (p *funcPrompt) Name() string           { return p.name }
func (p *funcPrompt) Description() string    { return p.description }
func (p *funcPrompt) Arguments() []*Argument { return p.arguments }
func (p *funcPrompt) Invoke(ctx context.Context, req *Request) ([]*Message, error) {
	return p.invoker(ctx, req)
}

// NewPrompt creates a Prompt from plain functions. args may be nil.
func NewPrompt(
	name, description string,
	args []*Argument,
	invoker func(ctx context.Context, req *Request) ([]*Message, error),
) Prompt {
	if args == nil {
		args = []*Argument{}
	}

	return &funcPrompt{
		name:        name,
		description: description,
		arguments:   args,
		invoker:     invoker,
	}
}

// promptToMap serialises a Prompt for inclusion in a prompts/list response.
func promptToMap(p Prompt) map[string]any {
	args := make([]map[string]any, 0, len(p.Arguments()))

	for _, a := range p.Arguments() {
		args = append(args, a.ToMap())
	}

	return map[string]any{
		"name":        p.Name(),
		"description": p.Description(),
		"arguments":   args,
	}
}
