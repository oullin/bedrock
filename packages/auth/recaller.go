package auth

import "strings"

// Recaller parses a remember-me cookie payload.
type Recaller struct {
	value string
}

// NewRecaller creates a recaller parser.
func NewRecaller(value string) Recaller {
	return Recaller{value: value}
}

// ID returns the user id segment.
func (r Recaller) ID() string { return r.segment(0) }

// Token returns the remember token segment.
func (r Recaller) Token() string { return r.segment(1) }

// Hash returns the password hash segment.
func (r Recaller) Hash() string { return r.segment(2) }

// Valid reports whether the payload is usable.
func (r Recaller) Valid() bool {
	parts := strings.Split(r.value, "|")
	return len(parts) >= 3 && strings.TrimSpace(parts[0]) != "" && strings.TrimSpace(parts[1]) != ""
}

func (r Recaller) segment(index int) string {
	parts := strings.Split(r.value, "|")
	if index >= len(parts) {
		return ""
	}

	return parts[index]
}
