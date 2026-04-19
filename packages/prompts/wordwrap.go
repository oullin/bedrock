package prompts

import (
	"strings"
	"unicode/utf8"
)

// WordWrap wraps text to the given width, respecting ANSI escape sequences
// (they don't count toward the visible width). Mirrors Laravel Prompts'
// AnsiWordwrap.
func WordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder
	lines := strings.Split(text, "\n")

	for li, line := range lines {
		if li > 0 {
			result.WriteByte('\n')
		}

		wrapLine(&result, line, width)
	}

	return result.String()
}

func wrapLine(result *strings.Builder, line string, width int) {
	if visibleLen(line) <= width {
		result.WriteString(line)

		return
	}

	words := strings.Fields(line)

	if len(words) == 0 {
		return
	}

	col := 0

	for i, word := range words {
		wl := visibleLen(word)

		if i > 0 {
			if col+1+wl > width {
				result.WriteByte('\n')
				col = 0
			} else {
				result.WriteByte(' ')
				col++
			}
		}

		if wl > width && col == 0 {
			// Word is too long for a single line — break it.
			breakLongWord(result, word, width, &col)
		} else if col+wl > width && col > 0 {
			result.WriteByte('\n')
			col = 0
			result.WriteString(word)
			col += wl
		} else {
			result.WriteString(word)
			col += wl
		}
	}
}

func breakLongWord(result *strings.Builder, word string, width int, col *int) {
	runes := []rune(stripAnsi(word))
	rawRunes := []rune(word)

	written := 0
	ri := 0

	for ri < len(rawRunes) {
		// Check if we're in an ANSI escape sequence.
		if rawRunes[ri] == '\x1b' {
			// Write the entire escape sequence.
			for ri < len(rawRunes) {
				result.WriteRune(rawRunes[ri])
				ri++

				if ri > 0 && ((rawRunes[ri-1] >= 'A' && rawRunes[ri-1] <= 'Z') ||
					(rawRunes[ri-1] >= 'a' && rawRunes[ri-1] <= 'z')) {
					break
				}
			}

			continue
		}

		if *col >= width {
			result.WriteByte('\n')
			*col = 0
		}

		result.WriteRune(rawRunes[ri])
		*col++
		written++
		ri++
	}

	_ = runes
}

// visibleLen returns the visible length of a string, ignoring ANSI codes.
func visibleLen(s string) int {
	return utf8.RuneCountInString(stripAnsi(s))
}

// stripAnsi removes ANSI escape sequences from a string.
func stripAnsi(s string) string {
	var result strings.Builder
	runes := []rune(s)
	i := 0

	for i < len(runes) {
		if runes[i] == '\x1b' {
			i++ // skip ESC

			if i < len(runes) && runes[i] == '[' {
				i++ // skip [
				// Skip until terminator letter.
				for i < len(runes) && !((runes[i] >= 'A' && runes[i] <= 'Z') || (runes[i] >= 'a' && runes[i] <= 'z')) {
					i++
				}

				if i < len(runes) {
					i++ // skip terminator
				}
			} else if i < len(runes) && runes[i] == ']' {
				// OSC sequence — skip until BEL or ST.
				i++

				for i < len(runes) && runes[i] != '\x07' {
					if runes[i] == '\x1b' && i+1 < len(runes) && runes[i+1] == '\\' {
						i += 2

						break
					}

					i++
				}

				if i < len(runes) && runes[i] == '\x07' {
					i++
				}
			}

			continue
		}

		result.WriteRune(runes[i])
		i++
	}

	return result.String()
}
