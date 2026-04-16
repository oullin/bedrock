package prompts

// Key constants for terminal input, matching Laravel Prompts Key class.
const (
	KeyUp              = "\x1b[A"
	KeyShiftUp         = "\x1b[1;2A"
	KeyPageUp          = "\x1b[5~"
	KeyDown            = "\x1b[B"
	KeyShiftDown       = "\x1b[1;2B"
	KeyPageDown        = "\x1b[6~"
	KeyRight           = "\x1b[C"
	KeyLeft            = "\x1b[D"
	KeyUpArrow         = "\x1bOA"
	KeyDownArrow       = "\x1bOB"
	KeyRightArrow      = "\x1bOC"
	KeyLeftArrow       = "\x1bOD"
	KeyEscape          = "\x1b"
	KeyDelete          = "\x1b[3~"
	KeyBackspace       = "\x7f"
	KeyEnter           = "\n"
	KeySpace           = " "
	KeyTab             = "\t"
	KeyShiftTab        = "\x1b[Z"
	KeyCtrlC           = "\x03"
	KeyCtrlP           = "\x10"
	KeyCtrlN           = "\x0e"
	KeyCtrlF           = "\x06"
	KeyCtrlB           = "\x02"
	KeyCtrlH           = "\x08"
	KeyCtrlA           = "\x01"
	KeyCtrlD           = "\x04"
	KeyCtrlE           = "\x05"
	KeyCtrlU           = "\x15"
	KeyOptionBackspace = "\x1b\x7f"
)

// KeyHome and KeyEnd have multiple possible escape sequences.
var (
	KeyHome = []string{"\x1b[1~", "\x1bOH", "\x1b[H", "\x1b[7~"}
	KeyEnd  = []string{"\x1b[4~", "\x1bOF", "\x1b[F", "\x1b[8~"}
)

// OneOfKey checks if the given input matches any of the provided key constants.
func OneOfKey(keys []string, input string) bool {
	for _, k := range keys {
		if k == input {
			return true
		}
	}

	return false
}

// IsUpKey returns true if the input is any of the up key variants.
func IsUpKey(input string) bool {
	return input == KeyUp || input == KeyUpArrow || input == KeyCtrlP
}

// IsDownKey returns true if the input is any of the down key variants.
func IsDownKey(input string) bool {
	return input == KeyDown || input == KeyDownArrow || input == KeyCtrlN
}

// IsLeftKey returns true if the input is any of the left key variants.
func IsLeftKey(input string) bool {
	return input == KeyLeft || input == KeyLeftArrow || input == KeyCtrlB
}

// IsRightKey returns true if the input is any of the right key variants.
func IsRightKey(input string) bool {
	return input == KeyRight || input == KeyRightArrow || input == KeyCtrlF
}
