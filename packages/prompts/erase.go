package prompts

import "fmt"

// EraseLines erases the given number of lines from the cursor position upward,
// moving the cursor to the beginning of the first erased line.
func EraseLines(w Writer, count int) {
	for i := 0; i < count; i++ {
		if i > 0 {
			w.Write("\x1b[1A") // move up
		}

		w.Write("\x1b[2K") // erase entire line
	}

	w.Write(fmt.Sprintf("\x1b[%dG", 1)) // move to column 1
}

// EraseDown erases from the cursor position to the end of the screen.
func EraseDown(w Writer) {
	w.Write("\x1b[J")
}

// EraseLine erases the entire current line.
func EraseLine(w Writer) {
	w.Write("\x1b[2K")
}
