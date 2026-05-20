package ai

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	contractsai "github.com/bedrock/packages/contracts/ai"
)

// SubAgent adapts any Promptable into a contractsai.Tool so it can be returned
// from a parent agent's Tools() and invoked by the LLM as a delegated task.
// Mirrors Upstream's sub-agent feature (upstream/ai 0.x).
//
// Each SubAgent.Handle call invokes the wrapped agent in isolation: the parent
// agent's conversation history is not passed in. Callers must include all
// necessary context in the "task" argument.
type SubAgent struct {
	agent       contractsai.Promptable
	name        string
	description string
	schema      map[string]any
}

// Compile-time interface checks.
var (
	_ contractsai.Tool         = (*SubAgent)(nil)
	_ contractsai.CanActAsTool = (*SubAgent)(nil)
)

// AsTool wraps a Promptable so it can be added to another agent's tools.
//
// If the wrapped agent implements contractsai.CanActAsTool, its Name() and
// Description() are used. Otherwise a snake_cased Go type name and a generic
// description are derived via reflection.
//
// Panics if agent is nil — a nil sub-agent would only fail at Handle() time
// with a nil-pointer dereference, which is harder to diagnose.
func AsTool(agent contractsai.Promptable) *SubAgent {
	if agent == nil {
		panic("ai: AsTool received a nil Promptable")
	}

	name, description := resolveSubAgentIdentity(agent)

	return &SubAgent{
		agent:       agent,
		name:        name,
		description: description,
		schema:      defaultSubAgentSchema(),
	}
}

// WithName overrides the auto-derived tool name.
func (s *SubAgent) WithName(name string) *SubAgent {
	s.name = name

	return s
}

// WithDescription overrides the auto-derived tool description.
func (s *SubAgent) WithDescription(description string) *SubAgent {
	s.description = description

	return s
}

// Name returns the tool identifier surfaced to the parent LLM.
func (s *SubAgent) Name() string { return s.name }

// Description returns the human-readable description surfaced to the parent LLM.
func (s *SubAgent) Description() string { return s.description }

// Schema returns the JSON schema describing the tool's input. The argument is
// accepted for contractsai.Tool symmetry but ignored; sub-agents always accept
// a single string "task" property.
func (s *SubAgent) Schema(_ contractsai.JsonSchema) map[string]any {
	return s.schema
}

// Handle delegates the LLM's tool invocation to the wrapped agent's Prompt,
// returning the assistant text as the tool result. The wrapped agent runs in
// isolation: no parent messages are forwarded.
func (s *SubAgent) Handle(ctx context.Context, req *contractsai.ToolRequest) (any, error) {
	task, _ := req.Arguments["task"].(string)

	if task == "" {
		return nil, fmt.Errorf("ai: sub-agent %q invoked without a 'task' argument", s.name)
	}

	resp, err := s.agent.Prompt(ctx, task)

	if err != nil {
		return nil, err
	}

	return resp.GetText(), nil
}

// Agent returns the wrapped Promptable. Useful for tests and callers that need
// to inspect or further configure the underlying agent.
func (s *SubAgent) Agent() contractsai.Promptable { return s.agent }

func defaultSubAgentSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"task": map[string]any{
				"type":        "string",
				"description": "Self-contained task description; the sub-agent receives no parent history.",
			},
		},
		"required": []string{"task"},
	}
}

func resolveSubAgentIdentity(agent contractsai.Promptable) (string, string) {
	if a, ok := agent.(contractsai.CanActAsTool); ok {
		return a.Name(), a.Description()
	}

	name := snakeCase(typeName(agent))

	if name == "" {
		name = "sub_agent"
	}

	return name, "Delegate a task to the " + name + " sub-agent."
}

func typeName(v any) string {
	t := reflect.TypeOf(v)

	if t == nil {
		return ""
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

// snakeCase converts CamelCase or PascalCase to snake_case. Empty input yields "".
func snakeCase(s string) string {
	if s == "" {
		return ""
	}

	var b strings.Builder

	runes := []rune(s)

	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			next := rune(0)

			if i+1 < len(runes) {
				next = runes[i+1]
			}
			// Insert underscore between a lowercase/digit and an uppercase,
			// or between two uppercases followed by a lowercase (e.g. "URLLoader" -> "url_loader").
			if unicode.IsLower(prev) || unicode.IsDigit(prev) || (unicode.IsUpper(prev) && next != 0 && unicode.IsLower(next)) {
				b.WriteRune('_')
			}
		}

		b.WriteRune(unicode.ToLower(r))
	}

	return b.String()
}
