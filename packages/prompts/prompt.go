package prompts

import (
	"fmt"
	"strings"
	"sync"
)

// State represents the current state of a prompt.
type State string

// ValidateFunc is a validation function. Returns "" on success, or an error
// message string on failure.
type ValidateFunc func(value string) string

// TransformFunc transforms the prompt value before returning it.
type TransformFunc func(value string) string

// RenderFunc renders a prompt to a string given the current state.
type RenderFunc func(state State) string

// Prompt is the base type embedded by all interactive prompts. It manages
// the input loop, rendering cycle, validation, and terminal state.
type Prompt struct {
	mu            sync.Mutex
	state         State
	errorMsg      string
	cancelMessage string
	required      any // bool or string
	transform     TransformFunc
	validate      ValidateFunc
	defaultValue  string
	prevFrame     string
	newLines      int
	cursorHidden  bool
	listeners     map[string][]func(string)
	renderer      RenderFunc
	valueFn       func() string
	terminal      Terminal
	writer        Writer
}

const (
	StateInitial State = "initial"
	StateActive  State = "active"
	StateError   State = "error"
	StateSubmit  State = "submit"
	StateCancel  State = "cancel"
)

// newPrompt creates a new base Prompt with defaults.
func newPrompt() *Prompt {
	return &Prompt{
		state:     StateInitial,
		listeners: make(map[string][]func(string)),
	}
}

// On registers a key event handler. The handler receives the key input.
func (p *Prompt) On(event string, handler func(key string)) {
	p.listeners[event] = append(p.listeners[event], handler)
}

// emit dispatches an event to all registered handlers.
func (p *Prompt) emit(event string, key string) {
	for _, h := range p.listeners[event] {
		h(key)
	}
}

// State returns the current prompt state.
func (p *Prompt) State() State {
	return p.state
}

// Error returns the current error message.
func (p *Prompt) Error() string {
	return p.errorMsg
}

// run executes the prompt's input loop. It sets up the terminal, renders
// the prompt, reads keys, and dispatches events until the prompt is
// submitted or cancelled.
func (p *Prompt) run() (string, error) {
	if !isInteractive() {
		return p.runNonInteractive()
	}

	p.terminal = getTerminal()
	p.writer = getWriter()

	if err := p.terminal.SetTty("-icanon -broadcastclient"); err != nil {
		return "", err
	}

	defer func() { _ = p.terminal.RestoreTty() }()

	p.hideCursor()

	defer p.showCursor()

	p.state = StateActive
	p.render()

	for {
		key, err := p.terminal.Read()

		if err != nil {
			return "", err
		}

		if key == KeyCtrlC {
			p.state = StateCancel
			p.render()

			return "", ErrCancelled
		}

		p.emit("key", key)
		p.render()

		if p.state == StateSubmit {
			val := ""

			if p.valueFn != nil {
				val = p.valueFn()
			}

			if p.transform != nil {
				val = p.transform(val)
			}

			return val, nil
		}
	}
}

// submit validates the current value and transitions to StateSubmit.
func (p *Prompt) submit() {
	val := ""

	if p.valueFn != nil {
		val = p.valueFn()
	}

	// Check required.
	if p.required != nil {
		switch r := p.required.(type) {
		case bool:
			if r && strings.TrimSpace(val) == "" {
				p.state = StateError
				p.errorMsg = "Required."

				return
			}
		case string:
			if strings.TrimSpace(val) == "" {
				p.state = StateError

				if r != "" {
					p.errorMsg = r
				} else {
					p.errorMsg = "Required."
				}

				return
			}
		}
	}

	// Run validation.
	if p.validate != nil {
		if msg := p.validate(val); msg != "" {
			p.state = StateError
			p.errorMsg = msg

			return
		}
	}

	p.state = StateSubmit
}

// render calls the renderer and writes the frame to the terminal.
func (p *Prompt) render() {
	if p.renderer == nil {
		return
	}

	frame := p.renderer(p.state)

	// On first render, just write the frame.
	if p.prevFrame == "" {
		p.writer.Write(frame)
		p.prevFrame = frame
		p.newLines = strings.Count(frame, "\n")

		return
	}

	// Erase previous frame and write new one.
	if p.newLines > 0 {
		EraseLines(p.writer, p.newLines+1)
	}

	p.writer.Write(frame)
	p.prevFrame = frame
	p.newLines = strings.Count(frame, "\n")
}

func (p *Prompt) hideCursor() {
	if !p.cursorHidden {
		HideCursor(p.writer)
		p.cursorHidden = true
	}
}

func (p *Prompt) showCursor() {
	if p.cursorHidden {
		ShowCursor(p.writer)
		p.cursorHidden = false
	}
}

// runNonInteractive returns the default value without rendering.
// If required and no default, returns ErrNonInteractive.
// If validation fails on the default, returns ErrValidation.
func (p *Prompt) runNonInteractive() (string, error) {
	val := p.defaultValue

	if p.valueFn != nil && val == "" {
		val = p.valueFn()
	}

	// Check required.
	if p.required != nil {
		switch r := p.required.(type) {
		case bool:
			if r && strings.TrimSpace(val) == "" {
				return "", fmt.Errorf("%w: value is required", ErrNonInteractive)
			}
		case string:
			if strings.TrimSpace(val) == "" {
				msg := "value is required"

				if r != "" {
					msg = r
				}

				return "", fmt.Errorf("%w: %s", ErrNonInteractive, msg)
			}
		}
	}

	// Validate.
	if p.validate != nil {
		if msg := p.validate(val); msg != "" {
			return "", fmt.Errorf("%w: %s", ErrValidation, msg)
		}
	}

	// Transform.
	if p.transform != nil {
		val = p.transform(val)
	}

	return val, nil
}
