// Package responses contains all response types returned by the AI package.
// They mirror the Upstream\Ai\Responses namespace.
package responses

import (
	"github.com/bedrock/packages/ai/data"
	"github.com/bedrock/packages/ai/messages"
)

// TextResponse is the base response for any text generation call.
// Mirrors Upstream\Ai\Responses\TextResponse.
type TextResponse struct {
	Text        string
	Usage       data.Usage
	Meta        data.Meta
	Messages    []any // slice of *messages.Message subtypes
	ToolCalls   []data.ToolCall
	ToolResults []data.ToolResult
	Steps       []data.Step
}

// WithMessages attaches the full message history and extracts tool call/result data.
func (r *TextResponse) WithMessages(msgs []any) *TextResponse {
	r.Messages = msgs

	var calls []data.ToolCall

	var results []data.ToolResult

	for _, m := range msgs {
		switch v := m.(type) {
		case *messages.AssistantMessage:
			calls = append(calls, v.ToolCalls...)
		case *messages.ToolResultMessage:
			results = append(results, v.ToolResults...)
		}
	}

	return r.WithToolCallsAndResults(calls, results)
}

// WithToolCallsAndResults sets tool call and result collections, filtering
// Anthropic-style structured output pseudo-calls.
func (r *TextResponse) WithToolCallsAndResults(calls []data.ToolCall, results []data.ToolResult) *TextResponse {
	filtered := calls[:0]

	for _, c := range calls {
		if c.Name != "output_structured_data" {
			filtered = append(filtered, c)
		}
	}

	r.ToolCalls = filtered
	r.ToolResults = results

	return r
}

// WithSteps attaches the reasoning steps.
func (r *TextResponse) WithSteps(steps []data.Step) *TextResponse {
	r.Steps = steps

	return r
}

// String returns the text content, satisfying fmt.Stringer.
func (r *TextResponse) String() string { return r.Text }
