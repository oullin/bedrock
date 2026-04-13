package notifications

// Action represents a clickable action in a notification message.
type Action struct {
	// Text is the display text for the action.
	Text string
	// URL is the target URL for the action.
	URL string
}

// NewAction creates an Action with the given text and URL.
func NewAction(text, url string) Action {
	return Action{Text: text, URL: url}
}
