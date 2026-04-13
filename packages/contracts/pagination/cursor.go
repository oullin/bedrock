package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// ErrInvalidCursor is returned when a cursor string cannot be decoded.

// Cursor encapsulates the parameters used for cursor-based pagination.
type Cursor struct {
	parameters        map[string]string
	pointsToNextItems bool
}

var ErrInvalidCursor = errors.New("pagination: invalid cursor")

// NewCursor creates a new Cursor with the given parameters and direction.
func NewCursor(parameters map[string]string, pointsToNextItems bool) *Cursor {
	return &Cursor{
		parameters:        parameters,
		pointsToNextItems: pointsToNextItems,
	}
}

// DecodeCursor decodes a base64-encoded cursor string into a Cursor.
// Returns nil, nil if the encoded string is empty.
func DecodeCursor(encoded string) (*Cursor, error) {
	if encoded == "" {
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		return nil, ErrInvalidCursor
	}

	var payload struct {
		Parameters        map[string]string `json:"parameters"`
		PointsToNextItems bool              `json:"points_to_next_items"`
	}

	if err := json.Unmarshal(decoded, &payload); err != nil {
		return nil, ErrInvalidCursor
	}

	return NewCursor(payload.Parameters, payload.PointsToNextItems), nil
}

// Parameter returns the value for a single cursor parameter.
func (c *Cursor) Parameter(name string) string {
	if c == nil {
		return ""
	}

	return c.parameters[name]
}

// Parameters returns all cursor parameters.
func (c *Cursor) Parameters() map[string]string {
	if c == nil {
		return nil
	}

	result := make(map[string]string, len(c.parameters))

	for k, v := range c.parameters {
		result[k] = v
	}

	return result
}

// PointsToNextItems reports whether this cursor points forward.
func (c *Cursor) PointsToNextItems() bool {
	if c == nil {
		return true
	}

	return c.pointsToNextItems
}

// PointsToPreviousItems reports whether this cursor points backward.
func (c *Cursor) PointsToPreviousItems() bool {
	return !c.PointsToNextItems()
}

// Encode returns the base64-encoded string representation of the cursor.
func (c *Cursor) Encode() string {
	if c == nil {
		return ""
	}

	payload := struct {
		Parameters        map[string]string `json:"parameters"`
		PointsToNextItems bool              `json:"points_to_next_items"`
	}{
		Parameters:        c.parameters,
		PointsToNextItems: c.pointsToNextItems,
	}

	b, _ := json.Marshal(payload)

	return base64.StdEncoding.EncodeToString(b)
}

// ToMap returns the cursor as a map suitable for JSON serialization.
func (c *Cursor) ToMap() map[string]any {
	if c == nil {
		return nil
	}

	return map[string]any{
		"parameters":           c.parameters,
		"points_to_next_items": c.pointsToNextItems,
	}
}
