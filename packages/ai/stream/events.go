// Package stream contains the streaming event types emitted during a
// streaming text generation. Events are produced by StreamText and consumed
// by StreamableAgentResponse.Each().
package stream

// Event is the common interface for all streaming events.
type Event interface {
	EventType() string
}

// StreamStart is emitted once at the beginning of a stream.
type StreamStart struct {
	InvocationID string
}

func (e StreamStart) EventType() string { return "stream.start" }

// StreamEnd is emitted once at the end of a stream.
type StreamEnd struct {
	InvocationID string
}

func (e StreamEnd) EventType() string { return "stream.end" }

// TextStart signals that a text block is beginning.
type TextStart struct{}

func (e TextStart) EventType() string { return "text.start" }

// TextDelta carries an incremental text chunk.
type TextDelta struct {
	Delta string
}

func (e TextDelta) EventType() string { return "text.delta" }

// TextEnd signals that a text block has finished.
type TextEnd struct {
	Text string // full accumulated text for this block
}

func (e TextEnd) EventType() string { return "text.end" }

// ReasoningStart signals the beginning of a reasoning block.
type ReasoningStart struct{}

func (e ReasoningStart) EventType() string { return "reasoning.start" }

// ReasoningDelta carries an incremental reasoning chunk.
type ReasoningDelta struct {
	Delta string
}

func (e ReasoningDelta) EventType() string { return "reasoning.delta" }

// ReasoningEnd signals the end of a reasoning block.
type ReasoningEnd struct {
	Text string
}

func (e ReasoningEnd) EventType() string { return "reasoning.end" }

// ToolCallEvent signals that the LLM is invoking a tool.
type ToolCallEvent struct {
	ID        string
	Name      string
	Arguments map[string]any
}

func (e ToolCallEvent) EventType() string { return "tool_call" }

// ToolResultEvent carries the result of a tool invocation back into the stream.
type ToolResultEvent struct {
	ID     string
	Name   string
	Result any
}

func (e ToolResultEvent) EventType() string { return "tool_result" }

// CitationEvent carries a source citation.
type CitationEvent struct {
	Title string
	URL   string
}

func (e CitationEvent) EventType() string { return "citation" }

// ErrorEvent signals an error that occurred during streaming.
type ErrorEvent struct {
	Err error
}

func (e ErrorEvent) EventType() string { return "error" }
