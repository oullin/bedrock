package notifications

// SimpleMessage is a basic notification message with level, subject, greeting,
// body lines, an optional action, and a salutation.
type SimpleMessage struct {
	level      string
	subject    string
	greeting   string
	salutation string
	introLines []string
	outroLines []string
	actionText string
	actionURL  string
	mailerName string
}

// NewSimpleMessage creates a SimpleMessage with default info level.
func NewSimpleMessage() *SimpleMessage {
	return &SimpleMessage{level: "info"}
}

// Success sets the message level to success.
func (m *SimpleMessage) Success() *SimpleMessage {
	m.level = "success"

	return m
}

// Error sets the message level to error.
func (m *SimpleMessage) Error() *SimpleMessage {
	m.level = "error"

	return m
}

// Level sets the message level.
func (m *SimpleMessage) Level(level string) *SimpleMessage {
	m.level = level

	return m
}

// GetLevel returns the message level.
func (m *SimpleMessage) GetLevel() string { return m.level }

// Subject sets the message subject.
func (m *SimpleMessage) Subject(subject string) *SimpleMessage {
	m.subject = subject

	return m
}

// GetSubject returns the message subject.
func (m *SimpleMessage) GetSubject() string { return m.subject }

// Greeting sets the greeting line.
func (m *SimpleMessage) Greeting(greeting string) *SimpleMessage {
	m.greeting = greeting

	return m
}

// GetGreeting returns the greeting line.
func (m *SimpleMessage) GetGreeting() string { return m.greeting }

// Salutation sets the closing salutation.
func (m *SimpleMessage) Salutation(salutation string) *SimpleMessage {
	m.salutation = salutation

	return m
}

// GetSalutation returns the salutation.
func (m *SimpleMessage) GetSalutation() string { return m.salutation }

// Line appends a line to the message body. Lines added before an Action call
// are intro lines; lines added after are outro lines.
func (m *SimpleMessage) Line(line string) *SimpleMessage {
	if m.actionText == "" {
		m.introLines = append(m.introLines, line)
	} else {
		m.outroLines = append(m.outroLines, line)
	}

	return m
}

// LineIf appends a line only if the condition is true.
func (m *SimpleMessage) LineIf(condition bool, line string) *SimpleMessage {
	if condition {
		return m.Line(line)
	}

	return m
}

// Lines appends multiple lines to the message body.
func (m *SimpleMessage) Lines(lines []string) *SimpleMessage {
	for _, line := range lines {
		m.Line(line)
	}

	return m
}

// LinesIf appends multiple lines only if the condition is true.
func (m *SimpleMessage) LinesIf(condition bool, lines []string) *SimpleMessage {
	if condition {
		return m.Lines(lines)
	}

	return m
}

// With is an alias for Line.
func (m *SimpleMessage) With(line string) *SimpleMessage {
	return m.Line(line)
}

// Action sets the call-to-action button text and URL.
func (m *SimpleMessage) Action(text, url string) *SimpleMessage {
	m.actionText = text
	m.actionURL = url

	return m
}

// GetActionText returns the action text.
func (m *SimpleMessage) GetActionText() string { return m.actionText }

// GetActionURL returns the action URL.
func (m *SimpleMessage) GetActionURL() string { return m.actionURL }

// Mailer sets the mailer name to use for delivery.
func (m *SimpleMessage) Mailer(name string) *SimpleMessage {
	m.mailerName = name

	return m
}

// GetMailer returns the mailer name.
func (m *SimpleMessage) GetMailer() string { return m.mailerName }

// GetIntroLines returns the intro lines.
func (m *SimpleMessage) GetIntroLines() []string { return m.introLines }

// GetOutroLines returns the outro lines.
func (m *SimpleMessage) GetOutroLines() []string { return m.outroLines }

// ToMap returns the message as a map representation.
func (m *SimpleMessage) ToMap() map[string]any {
	result := map[string]any{
		"level":                m.level,
		"subject":              m.subject,
		"greeting":             m.greeting,
		"salutation":           m.salutation,
		"introLines":           m.introLines,
		"outroLines":           m.outroLines,
		"actionText":           m.actionText,
		"actionUrl":            m.actionURL,
		"displayableActionUrl": m.actionURL,
	}

	return result
}
