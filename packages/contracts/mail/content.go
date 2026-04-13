package mail

// Content describes the body of an email message. At least one of HTML,
// Text, Markdown, or HTMLString should be set.
type Content struct {
	// HTML is a template name or literal HTML when HTMLString is true.
	HTML string
	// Text is a template name for the plain-text alternative.
	Text string
	// Markdown is raw Markdown source converted to HTML before sending.
	Markdown string
	// HTMLString when true indicates that HTML contains literal HTML,
	// not a template name.
	HTMLString bool
	// With holds data passed to templates during rendering.
	With map[string]any
}
