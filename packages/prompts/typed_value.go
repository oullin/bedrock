package prompts

import (
	"strings"
	"unicode/utf8"
)

// TypedValue tracks user-typed text with cursor position. Embedded by
// text-based prompts (TextPrompt, PasswordPrompt, SuggestPrompt, etc.).
type TypedValue struct {
	typedValue     string
	cursorPosition int
}

// Value returns the current typed value.
func (tv *TypedValue) Value() string {
	return tv.typedValue
}

// SetValue sets the typed value and moves cursor to the end.
func (tv *TypedValue) SetValue(s string) {
	tv.typedValue = s
	tv.cursorPosition = utf8.RuneCountInString(s)
}

// CursorPosition returns the current cursor position (rune index).
func (tv *TypedValue) CursorPosition() int {
	return tv.cursorPosition
}

// TrackTypedValue registers key handlers on the prompt to manage text input.
func (tv *TypedValue) TrackTypedValue(p *Prompt, submitFn func()) {
	p.On("key", func(key string) {
		if p.state != StateActive && p.state != StateError {
			return
		}

		switch {
		case key == KeyEnter:
			submitFn()

		case key == KeyBackspace || key == KeyCtrlH:
			tv.deleteCharBackward()

		case key == KeyDelete:
			tv.deleteCharForward()

		case key == KeyOptionBackspace:
			tv.deleteWordBackward()

		case key == KeyCtrlU:
			tv.typedValue = runeSlice(tv.typedValue, tv.cursorPosition, utf8.RuneCountInString(tv.typedValue))
			tv.cursorPosition = 0

		case IsLeftKey(key):
			if tv.cursorPosition > 0 {
				tv.cursorPosition--
			}

		case IsRightKey(key):
			if tv.cursorPosition < utf8.RuneCountInString(tv.typedValue) {
				tv.cursorPosition++
			}

		case key == KeyCtrlA || OneOfKey(KeyHome, key):
			tv.cursorPosition = 0

		case key == KeyCtrlE || OneOfKey(KeyEnd, key):
			tv.cursorPosition = utf8.RuneCountInString(tv.typedValue)

		default:
			// Only insert printable characters.
			if len(key) > 0 && key[0] >= 32 && key[0] != 127 && !strings.HasPrefix(key, "\x1b") {
				tv.insertAtCursor(key)
			}
		}
	})
}

func (tv *TypedValue) insertAtCursor(text string) {
	runes := []rune(tv.typedValue)
	insert := []rune(text)
	pos := tv.cursorPosition

	if pos > len(runes) {
		pos = len(runes)
	}

	result := make([]rune, 0, len(runes)+len(insert))
	result = append(result, runes[:pos]...)
	result = append(result, insert...)
	result = append(result, runes[pos:]...)
	tv.typedValue = string(result)
	tv.cursorPosition = pos + len(insert)
}

func (tv *TypedValue) deleteCharBackward() {
	if tv.cursorPosition <= 0 {
		return
	}

	runes := []rune(tv.typedValue)
	tv.cursorPosition--
	runes = append(runes[:tv.cursorPosition], runes[tv.cursorPosition+1:]...)
	tv.typedValue = string(runes)
}

func (tv *TypedValue) deleteCharForward() {
	runes := []rune(tv.typedValue)

	if tv.cursorPosition >= len(runes) {
		return
	}

	runes = append(runes[:tv.cursorPosition], runes[tv.cursorPosition+1:]...)
	tv.typedValue = string(runes)
}

func (tv *TypedValue) deleteWordBackward() {
	if tv.cursorPosition <= 0 {
		return
	}

	runes := []rune(tv.typedValue)
	pos := tv.cursorPosition - 1

	// Skip trailing spaces.
	for pos > 0 && runes[pos] == ' ' {
		pos--
	}
	// Skip word characters.
	for pos > 0 && runes[pos-1] != ' ' {
		pos--
	}

	runes = append(runes[:pos], runes[tv.cursorPosition:]...)
	tv.typedValue = string(runes)
	tv.cursorPosition = pos
}

// AddCursor inserts a cursor indicator at the current position within the value,
// returning the value with a cursor character for rendering.
func (tv *TypedValue) AddCursor(value string, maxWidth int, cursorChar string) string {
	if cursorChar == "" {
		cursorChar = "│"
	}

	runes := []rune(value)
	pos := tv.cursorPosition

	if pos > len(runes) {
		pos = len(runes)
	}

	before := string(runes[:pos])
	after := ""

	if pos < len(runes) {
		after = string(runes[pos:])
	}

	return before + cursorChar + after
}

// runeSlice returns the substring from start to end rune positions.
func runeSlice(s string, start, end int) string {
	runes := []rune(s)

	if start > len(runes) {
		start = len(runes)
	}

	if end > len(runes) {
		end = len(runes)
	}

	if start < 0 {
		start = 0
	}

	return string(runes[start:end])
}
