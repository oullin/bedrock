package prompts

import "fmt"

// Title sets the terminal window title.
func Title(title string) {
	w := getWriter()
	w.Write(fmt.Sprintf("\x1b]0;%s\x07", title))
}
