package auth

import "strings"

// Recaller parses a remember-me cookie payload.
type Recaller struct {
	value string
}

// NewRecaller creates a new recaller parser.
func NewRecaller(value string) Recaller {
	return Recaller{value: value}
}

// ID returns the recaller user identifier.
func (r Recaller) ID() string {
	return r.segment(0)
}

// Token returns the remember token.
func (r Recaller) Token() string {
	return r.segment(1)
}

// Hash returns the password MAC stored in the cookie.
func (r Recaller) Hash() string {
	return r.segment(2)
}

// Valid reports whether the payload is well-formed.
func (r Recaller) Valid() bool {
	segments := strings.Split(r.value, "|")

	return len(segments) >= 3 && strings.TrimSpace(segments[0]) != "" && strings.TrimSpace(segments[1]) != ""
}

func (r Recaller) segment(index int) string {
	segments := strings.Split(r.value, "|")

	if index >= len(segments) {
		return ""
	}

	return segments[index]
}
