package prompts

import (
	"strings"
	"unicode/utf8"
)

// TextareaOption configures a TextareaPrompt.
type TextareaOption func(*TextareaPrompt)

// TextareaWithPlaceholder sets placeholder text.

// TextareaWithDefault sets the default value.

// TextareaWithRequired makes the prompt required.

// TextareaWithValidate sets validation.

// TextareaWithHint sets hint text.

// TextareaWithRows sets the number of visible rows.

// TextareaWithTransform sets the transform function.

// TextareaPrompt handles multi-line text input.
type TextareaPrompt struct {
	prompt      *Prompt
	TypedValue  TypedValue
	label       string
	placeholder string
	hint        string
	rows        int
}

func TextareaWithPlaceholder(s string) TextareaOption {
	return func(p *TextareaPrompt) { p.placeholder = s }
}

func TextareaWithDefault(s string) TextareaOption {
	return func(p *TextareaPrompt) { p.TypedValue.SetValue(s) }
}

func TextareaWithRequired(v any) TextareaOption {
	return func(p *TextareaPrompt) { p.prompt.required = v }
}

func TextareaWithValidate(fn ValidateFunc) TextareaOption {
	return func(p *TextareaPrompt) { p.prompt.validate = fn }
}

func TextareaWithHint(s string) TextareaOption { return func(p *TextareaPrompt) { p.hint = s } }

func TextareaWithRows(n int) TextareaOption { return func(p *TextareaPrompt) { p.rows = n } }

func TextareaWithTransform(fn TransformFunc) TextareaOption {
	return func(p *TextareaPrompt) { p.prompt.transform = fn }
}

// Textarea displays a multi-line text input prompt.
func Textarea(label string, opts ...TextareaOption) (string, error) {
	p := &TextareaPrompt{
		prompt: newPrompt(),
		label:  label,
		rows:   5,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.prompt.valueFn = p.TypedValue.Value
	p.prompt.renderer = func(state State) string {
		return getTheme().TextareaRenderer(p, state)
	}

	// Register key handler for newline on Enter.
	p.prompt.On("key", func(key string) {
		if p.prompt.state != StateActive && p.prompt.state != StateError {
			return
		}

		switch {
		case key == KeyEnter:
			// In textarea, Enter inserts a newline.
			p.TypedValue.insertAtCursor("\n")
		case key == KeyTab:
			// Ctrl+D or Tab submits.
			return
		case key == KeyCtrlD:
			p.prompt.submit()
		case key == KeyBackspace || key == KeyCtrlH:
			p.TypedValue.deleteCharBackward()
		case key == KeyDelete:
			p.TypedValue.deleteCharForward()
		case key == KeyOptionBackspace:
			p.TypedValue.deleteWordBackward()
		case key == KeyCtrlU:
			p.TypedValue.typedValue = runeSlice(p.TypedValue.typedValue, p.TypedValue.cursorPosition, utf8.RuneCountInString(p.TypedValue.typedValue))
			p.TypedValue.cursorPosition = 0
		case IsLeftKey(key):
			if p.TypedValue.cursorPosition > 0 {
				p.TypedValue.cursorPosition--
			}
		case IsRightKey(key):
			if p.TypedValue.cursorPosition < utf8.RuneCountInString(p.TypedValue.typedValue) {
				p.TypedValue.cursorPosition++
			}
		case IsUpKey(key):
			p.moveCursorUp()
		case IsDownKey(key):
			p.moveCursorDown()
		case key == KeyCtrlA || OneOfKey(KeyHome, key):
			p.TypedValue.cursorPosition = 0
		case key == KeyCtrlE || OneOfKey(KeyEnd, key):
			p.TypedValue.cursorPosition = utf8.RuneCountInString(p.TypedValue.typedValue)
		default:
			if len(key) > 0 && key[0] >= 32 && key[0] != 127 && !strings.HasPrefix(key, "\x1b") {
				p.TypedValue.insertAtCursor(key)
			}
		}
	})

	return p.prompt.run()
}

func (p *TextareaPrompt) moveCursorUp() {
	lines := strings.Split(p.TypedValue.typedValue, "\n")
	pos := p.TypedValue.cursorPosition
	lineStart := 0

	for _, line := range lines {
		lineLen := utf8.RuneCountInString(line) + 1

		if pos < lineStart+lineLen {
			col := pos - lineStart

			if lineStart == 0 {
				return // Already on first line.
			}
			// Find previous line start.
			prevLineStart := 0

			for _, l := range lines {
				pll := utf8.RuneCountInString(l) + 1

				if prevLineStart+pll >= lineStart {
					break
				}

				prevLineStart += pll
			}

			prevLineLen := lineStart - prevLineStart - 1

			if col > prevLineLen {
				col = prevLineLen
			}

			p.TypedValue.cursorPosition = prevLineStart + col

			return
		}

		lineStart += lineLen
	}
}

func (p *TextareaPrompt) moveCursorDown() {
	lines := strings.Split(p.TypedValue.typedValue, "\n")
	pos := p.TypedValue.cursorPosition
	lineStart := 0

	for i, line := range lines {
		lineLen := utf8.RuneCountInString(line) + 1

		if pos < lineStart+lineLen {
			col := pos - lineStart

			if i >= len(lines)-1 {
				return // Already on last line.
			}

			nextLineStart := lineStart + lineLen
			nextLineLen := utf8.RuneCountInString(lines[i+1])

			if col > nextLineLen {
				col = nextLineLen
			}

			p.TypedValue.cursorPosition = nextLineStart + col

			return
		}

		lineStart += lineLen
	}
}

// Value returns the current typed value.
func (p *TextareaPrompt) Value() string { return p.TypedValue.Value() }
