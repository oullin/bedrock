package prompts

import "strings"

// ParseAnsi splits an ANSI-formatted string into segments, each with its
// styling information. This is used for precise text measurement and
// manipulation while preserving color/style information.
type AnsiSegment struct {
	Text  string
	Style string
}

// ParseAnsiText parses a string containing ANSI escape sequences into a
// slice of segments, each containing the text and the currently active
// ANSI styling. Mirrors Laravel Prompts' parseAnsiText.
func ParseAnsiText(s string) []AnsiSegment {
	var segments []AnsiSegment

	var currentStyle strings.Builder

	var currentText strings.Builder

	runes := []rune(s)
	i := 0

	for i < len(runes) {
		if runes[i] == '\x1b' && i+1 < len(runes) && runes[i+1] == '[' {
			// Flush current text as a segment if non-empty.
			if currentText.Len() > 0 {
				segments = append(segments, AnsiSegment{
					Text:  currentText.String(),
					Style: currentStyle.String(),
				})
				currentText.Reset()
			}

			// Parse the escape sequence.
			var seq strings.Builder

			seq.WriteRune(runes[i])
			i++
			seq.WriteRune(runes[i])
			i++

			for i < len(runes) && !((runes[i] >= 'A' && runes[i] <= 'Z') || (runes[i] >= 'a' && runes[i] <= 'z')) {
				seq.WriteRune(runes[i])
				i++
			}

			if i < len(runes) {
				seq.WriteRune(runes[i])
				i++
			}

			seqStr := seq.String()

			if seqStr == resetCode {
				currentStyle.Reset()
			} else {
				currentStyle.WriteString(seqStr)
			}

			continue
		}

		currentText.WriteRune(runes[i])
		i++
	}

	// Flush remaining text.
	if currentText.Len() > 0 {
		segments = append(segments, AnsiSegment{
			Text:  currentText.String(),
			Style: currentStyle.String(),
		})
	}

	return segments
}

// StripAnsi removes all ANSI escape sequences from a string.
// This is the public version of stripAnsi, available for external use.
func StripAnsi(s string) string {
	return stripAnsi(s)
}

// VisibleWidth returns the visible width of an ANSI string.
func VisibleWidth(s string) int {
	return visibleLen(s)
}
