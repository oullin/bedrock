package mail

import "fmt"

// Address represents an email address with an optional display name.
type Address struct {
	Name  string
	Email string
}

// String returns the RFC 5322 formatted address.
func (a Address) String() string {
	if a.Name == "" {
		return a.Email
	}

	return fmt.Sprintf("%s <%s>", a.Name, a.Email)
}
