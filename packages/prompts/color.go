package prompts

import "fmt"

// ANSI text style codes.
const (
	resetCode         = "\x1b[0m"
	boldCode          = "\x1b[1m"
	dimCode           = "\x1b[2m"
	italicCode        = "\x1b[3m"
	underlineCode     = "\x1b[4m"
	inverseCode       = "\x1b[7m"
	hiddenCode        = "\x1b[8m"
	strikethroughCode = "\x1b[9m"
)

// Reset wraps text with the ANSI reset code.
func Reset(text string) string { return resetCode + text + resetCode }

// Bold wraps text with the ANSI bold code.
func Bold(text string) string { return boldCode + text + resetCode }

// Dim wraps text with the ANSI dim code.
func Dim(text string) string { return dimCode + text + resetCode }

// Italic wraps text with the ANSI italic code.
func Italic(text string) string { return italicCode + text + resetCode }

// Underline wraps text with the ANSI underline code.
func Underline(text string) string { return underlineCode + text + resetCode }

// Inverse wraps text with the ANSI inverse code.
func Inverse(text string) string { return inverseCode + text + resetCode }

// Hidden wraps text with the ANSI hidden code.
func Hidden(text string) string { return hiddenCode + text + resetCode }

// Strikethrough wraps text with the ANSI strikethrough code.
func Strikethrough(text string) string { return strikethroughCode + text + resetCode }

// Foreground color functions.

func Black(text string) string   { return "\x1b[30m" + text + resetCode }
func Red(text string) string     { return "\x1b[31m" + text + resetCode }
func Green(text string) string   { return "\x1b[32m" + text + resetCode }
func Yellow(text string) string  { return "\x1b[33m" + text + resetCode }
func Blue(text string) string    { return "\x1b[34m" + text + resetCode }
func Magenta(text string) string { return "\x1b[35m" + text + resetCode }
func Cyan(text string) string    { return "\x1b[36m" + text + resetCode }
func White(text string) string   { return "\x1b[37m" + text + resetCode }
func Gray(text string) string    { return "\x1b[90m" + text + resetCode }

// Background color functions.

func BgBlack(text string) string   { return "\x1b[40m" + text + resetCode }
func BgRed(text string) string     { return "\x1b[41m" + text + resetCode }
func BgGreen(text string) string   { return "\x1b[42m" + text + resetCode }
func BgYellow(text string) string  { return "\x1b[43m" + text + resetCode }
func BgBlue(text string) string    { return "\x1b[44m" + text + resetCode }
func BgMagenta(text string) string { return "\x1b[45m" + text + resetCode }
func BgCyan(text string) string    { return "\x1b[46m" + text + resetCode }
func BgWhite(text string) string   { return "\x1b[47m" + text + resetCode }

// Fg256 wraps text with a 256-color foreground.
func Fg256(text string, code int) string {
	return fmt.Sprintf("\x1b[38;5;%dm%s%s", code, text, resetCode)
}

// Bg256 wraps text with a 256-color background.
func Bg256(text string, code int) string {
	return fmt.Sprintf("\x1b[48;5;%dm%s%s", code, text, resetCode)
}

// FgRGB wraps text with a true-color (24-bit) foreground.
func FgRGB(text string, r, g, b int) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s%s", r, g, b, text, resetCode)
}

// BgRGB wraps text with a true-color (24-bit) background.
func BgRGB(text string, r, g, b int) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm%s%s", r, g, b, text, resetCode)
}
