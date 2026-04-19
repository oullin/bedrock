package prompts

// Clear clears the terminal screen.
func Clear() {
	w := getWriter()
	w.Write("\x1b[2J\x1b[H")
}
