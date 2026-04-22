package support

// HtmlString represents a trusted HTML string.
type HtmlString struct {
	html string
}

// NewHtmlString creates a trusted HTML string wrapper.
func NewHtmlString(html string) HtmlString {
	return HtmlString{html: html}
}

// ToHTML returns the raw HTML value.
func (s HtmlString) ToHTML() string {
	return s.html
}

// String returns the raw HTML value.
func (s HtmlString) String() string {
	return s.html
}

// IsEmpty reports whether the HTML string is empty.
func (s HtmlString) IsEmpty() bool {
	return s.html == ""
}

// IsNotEmpty reports whether the HTML string is not empty.
func (s HtmlString) IsNotEmpty() bool {
	return !s.IsEmpty()
}
