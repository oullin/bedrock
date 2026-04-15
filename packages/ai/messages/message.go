// Package messages contains the conversation message types for the AI package.
// Mirrors the Laravel\Ai\Messages namespace.
package messages

import (
	"fmt"

	"github.com/bedrock/packages/ai/enums"
)

// Message is the base conversation message type.
// Mirrors Laravel\Ai\Messages\Message.
type Message struct {
	Role    enums.MessageRole
	Content *string
}

// NewMessage constructs a Message with the given role and optional content.
func NewMessage(role enums.MessageRole, content *string) *Message {
	return &Message{Role: role, Content: content}
}

// TryFrom constructs a Message from a map or struct with role/content fields.
// Accepts:
//   - *Message (returned as-is)
//   - map[string]any with "role" and optional "content" keys
func TryFrom(v any) (*Message, error) {
	switch val := v.(type) {
	case *Message:
		return val, nil
	case Message:
		return &val, nil
	case map[string]any:
		roleRaw, ok := val["role"]
		if !ok {
			return nil, fmt.Errorf("messages: missing 'role' key in map")
		}
		roleStr, ok := roleRaw.(string)
		if !ok {
			return nil, fmt.Errorf("messages: 'role' must be a string")
		}
		role, valid := enums.TryFromMessageRole(roleStr)
		if !valid {
			return nil, fmt.Errorf("messages: unknown role %q", roleStr)
		}
		var content *string
		if raw, exists := val["content"]; exists && raw != nil {
			if s, ok := raw.(string); ok {
				content = &s
			}
		}
		return NewMessage(role, content), nil
	default:
		return nil, fmt.Errorf("messages: unsupported type %T", v)
	}
}
