package prompts

// Note displays a styled notification message.
func Note(message string) {
	w := getWriter()
	w.Write(getTheme().NoteRenderer(message, ""))
}

// Error displays an error message.
func Error(message string) {
	w := getWriter()
	w.Write(getTheme().NoteRenderer(message, "error"))
}

// Warning displays a warning message.
func Warning(message string) {
	w := getWriter()
	w.Write(getTheme().NoteRenderer(message, "warning"))
}

// Info displays an info message.
func Info(message string) {
	w := getWriter()
	w.Write(getTheme().NoteRenderer(message, "info"))
}

// Alert displays an alert message.
func Alert(message string) {
	w := getWriter()
	w.Write(getTheme().NoteRenderer(message, "alert"))
}

// Intro displays an introduction message.
func Intro(message string) {
	w := getWriter()
	w.Write(getTheme().NoteRenderer(message, "intro"))
}

// Outro displays a closing message.
func Outro(message string) {
	w := getWriter()
	w.Write(getTheme().NoteRenderer(message, "outro"))
}
