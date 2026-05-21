package eloquent

import (
	"encoding/json"
	"time"
)

// HasAttributes provides attribute storage, retrieval, dirty tracking, and
// casting. It is the Go port of @bedrock\Database\Eloquent\Concerns\HasAttributes.
type HasAttributes struct {
	attributes map[string]any
	original   map[string]any
	casts      map[string]string
	appends    []string
	dateFormat string
}

// InitAttributes initializes the attribute maps.
func (h *HasAttributes) InitAttributes() {
	if h.attributes == nil {
		h.attributes = make(map[string]any)
	}

	if h.original == nil {
		h.original = make(map[string]any)
	}

	if h.casts == nil {
		h.casts = make(map[string]string)
	}
}

// GetAttribute returns the value of an attribute, applying casts.
func (h *HasAttributes) GetAttribute(key string) any {
	value, ok := h.attributes[key]

	if !ok {
		return nil
	}

	return h.castAttribute(key, value)
}

// SetAttribute sets an attribute value.
func (h *HasAttributes) SetAttribute(key string, value any) {
	h.InitAttributes()
	h.attributes[key] = value
}

// GetAttributes returns all attributes.
func (h *HasAttributes) GetAttributes() map[string]any {
	return h.attributes
}

// SetRawAttributes sets all attributes and syncs originals.
func (h *HasAttributes) SetRawAttributes(attributes map[string]any, sync bool) {
	h.attributes = attributes

	if sync {
		h.SyncOriginal()
	}
}

// GetOriginal returns the original attribute values.
func (h *HasAttributes) GetOriginal() map[string]any {
	return h.original
}

// GetOriginalAttribute returns the original value of a specific attribute.
func (h *HasAttributes) GetOriginalAttribute(key string) any {
	if h.original == nil {
		return nil
	}

	return h.original[key]
}

// SyncOriginal copies current attributes to originals.
func (h *HasAttributes) SyncOriginal() {
	h.original = make(map[string]any, len(h.attributes))

	for k, v := range h.attributes {
		h.original[k] = v
	}
}

// SyncOriginalAttribute syncs a single attribute.
func (h *HasAttributes) SyncOriginalAttribute(key string) {
	if h.original == nil {
		h.original = make(map[string]any)
	}

	h.original[key] = h.attributes[key]
}

// IsDirty checks if any of the given attributes have been modified.
// If no attributes are given, checks if any attribute is dirty.
func (h *HasAttributes) IsDirty(attributes ...string) bool {
	dirty := h.GetDirty()

	if len(attributes) == 0 {
		return len(dirty) > 0
	}

	for _, attr := range attributes {
		if _, ok := dirty[attr]; ok {
			return true
		}
	}

	return false
}

// IsClean checks if the given attributes have NOT been modified.
func (h *HasAttributes) IsClean(attributes ...string) bool {
	return !h.IsDirty(attributes...)
}

// WasChanged checks if any attributes were changed when the model was last saved.
func (h *HasAttributes) WasChanged(attributes ...string) bool {
	return h.IsDirty(attributes...)
}

// GetDirty returns the attributes that have been modified.
func (h *HasAttributes) GetDirty() map[string]any {
	dirty := make(map[string]any)

	for key, value := range h.attributes {
		original, exists := h.original[key]

		if !exists || value != original {
			dirty[key] = value
		}
	}

	return dirty
}

// GetChanges returns all changed attributes (alias for GetDirty).
func (h *HasAttributes) GetChanges() map[string]any {
	return h.GetDirty()
}

// HasAttribute checks if an attribute exists.
func (h *HasAttributes) HasAttribute(key string) bool {
	_, ok := h.attributes[key]

	return ok
}

// SetCasts sets the attribute casting map.
func (h *HasAttributes) SetCasts(casts map[string]string) {
	h.casts = casts
}

// GetCasts returns the attribute casting map.
func (h *HasAttributes) GetCasts() map[string]string {
	return h.casts
}

// HasCast checks if an attribute has a cast.
func (h *HasAttributes) HasCast(key string) bool {
	_, ok := h.casts[key]

	return ok
}

// SetAppends sets the attributes to append during serialization.
func (h *HasAttributes) SetAppends(appends []string) {
	h.appends = appends
}

// GetAppends returns the appended attributes.
func (h *HasAttributes) GetAppends() []string {
	return h.appends
}

// SetDateFormat sets the date serialization format.
func (h *HasAttributes) SetDateFormat(format string) {
	h.dateFormat = format
}

// ToMap returns the attributes as a map.
func (h *HasAttributes) ToMap() map[string]any {
	m := make(map[string]any, len(h.attributes))

	for k, v := range h.attributes {
		m[k] = h.castAttribute(k, v)
	}

	return m
}

// ToJSON serializes the attributes as JSON.
func (h *HasAttributes) ToJSON() ([]byte, error) {
	return json.Marshal(h.ToMap())
}

// castAttribute applies the registered cast for a key.
func (h *HasAttributes) castAttribute(key string, value any) any {
	castType, ok := h.casts[key]

	if !ok {
		return value
	}

	switch castType {
	case "int", "integer":
		return toInt(value)
	case "float", "double", "real":
		return toFloat(value)
	case "string":
		return toString(value)
	case "bool", "boolean":
		return toBoolVal(value)
	case "json", "array", "object":
		return value
	case "date", "datetime", "timestamp":
		return toTime(value)
	default:
		return value
	}
}

func toInt(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	case string:
		var i int64

		json.Unmarshal([]byte(n), &i)

		return i
	default:
		return 0
	}
}

func toFloat(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int64:
		return float64(n)
	case int:
		return float64(n)
	default:
		return 0
	}
}

func toString(v any) string {
	switch s := v.(type) {
	case string:
		return s
	case []byte:
		return string(s)
	default:
		b, _ := json.Marshal(v)

		return string(b)
	}
}

func toBoolVal(v any) bool {
	switch b := v.(type) {
	case bool:
		return b
	case int64:
		return b != 0
	case int:
		return b != 0
	case float64:
		return b != 0
	case string:
		return b == "1" || b == "true"
	default:
		return false
	}
}

func toTime(v any) time.Time {
	switch t := v.(type) {
	case time.Time:
		return t
	case string:
		parsed, err := time.Parse(time.RFC3339, t)

		if err != nil {
			parsed, _ = time.Parse("2006-01-02 15:04:05", t)
		}

		return parsed
	default:
		return time.Time{}
	}
}
