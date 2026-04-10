package auth

import "strings"

// Recaller parses the remember-me cookie value in the format "id|token|hash".
type Recaller struct {
	parts []string
}

// NewRecaller parses a remember-me cookie value. Returns nil if invalid.
func NewRecaller(value string) *Recaller {
	parts := strings.Split(value, "|")

	if len(parts) != 3 {
		return nil
	}

	return &Recaller{parts: parts}
}

// ID returns the user ID from the remember cookie.
func (r *Recaller) ID() string { return r.parts[0] }

// Token returns the remember token from the cookie.
func (r *Recaller) Token() string { return r.parts[1] }

// Hash returns the hash segment from the cookie.
func (r *Recaller) Hash() string { return r.parts[2] }

// Valid reports whether all three segments are non-empty.
func (r *Recaller) Valid() bool {
	return r.parts[0] != "" && r.parts[1] != "" && r.parts[2] != ""
}
