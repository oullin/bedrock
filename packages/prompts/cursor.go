package prompts

import "fmt"

// ANSI escape sequences for cursor control.
const (
	cursorHide = "\x1b[?25l"
	cursorShow = "\x1b[?25h"
)

// HideCursor writes the ANSI sequence to hide the cursor.
func HideCursor(w Writer) {
	w.Write(cursorHide)
}

// ShowCursor writes the ANSI sequence to show the cursor.
func ShowCursor(w Writer) {
	w.Write(cursorShow)
}

// MoveCursor moves the cursor by the given x (columns) and y (rows) offset.
// Positive x moves right, negative left. Positive y moves down, negative up.
func MoveCursor(w Writer, x, y int) {
	if x > 0 {
		w.Write(fmt.Sprintf("\x1b[%dC", x))
	} else if x < 0 {
		w.Write(fmt.Sprintf("\x1b[%dD", -x))
	}

	if y > 0 {
		w.Write(fmt.Sprintf("\x1b[%dB", y))
	} else if y < 0 {
		w.Write(fmt.Sprintf("\x1b[%dA", -y))
	}
}

// MoveCursorToColumn moves the cursor to the given absolute column (1-based).
func MoveCursorToColumn(w Writer, col int) {
	w.Write(fmt.Sprintf("\x1b[%dG", col))
}

// MoveCursorUp moves the cursor up by the given number of lines.
func MoveCursorUp(w Writer, lines int) {
	if lines > 0 {
		w.Write(fmt.Sprintf("\x1b[%dA", lines))
	}
}

// MoveCursorDown moves the cursor down by the given number of lines.
func MoveCursorDown(w Writer, lines int) {
	if lines > 0 {
		w.Write(fmt.Sprintf("\x1b[%dB", lines))
	}
}
