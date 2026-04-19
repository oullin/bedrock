package enums

// MessageRole identifies the participant of a conversation message.
// Mirrors Laravel\Ai\Messages\MessageRole.
type MessageRole string

const (
	RoleAssistant  MessageRole = "assistant"
	RoleUser       MessageRole = "user"
	RoleToolResult MessageRole = "tool_result"
)

// String returns the raw string value.
func (r MessageRole) String() string { return string(r) }

// Valid reports whether the role is a known value.
func (r MessageRole) Valid() bool {
	switch r {
	case RoleAssistant, RoleUser, RoleToolResult:
		return true
	}

	return false
}

// TryFrom parses a string into a MessageRole, returning false if unrecognised.
func TryFromMessageRole(s string) (MessageRole, bool) {
	r := MessageRole(s)

	return r, r.Valid()
}
