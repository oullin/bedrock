package orm

import "time"

// HasTimestamps manages created_at and updated_at columns automatically.
type HasTimestamps struct {
	timestamps bool
	createdAt  string
	updatedAt  string
}

// InitTimestamps sets up default timestamp column names.
func (h *HasTimestamps) InitTimestamps() {
	h.timestamps = true
	h.createdAt = "created_at"
	h.updatedAt = "updated_at"
}

// UsesTimestamps returns whether the model uses timestamps.
func (h *HasTimestamps) UsesTimestamps() bool {
	return h.timestamps
}

// SetTimestamps enables or disables automatic timestamps.
func (h *HasTimestamps) SetTimestamps(enabled bool) {
	h.timestamps = enabled
}

// GetCreatedAtColumn returns the created_at column name.
func (h *HasTimestamps) GetCreatedAtColumn() string {
	return h.createdAt
}

// GetUpdatedAtColumn returns the updated_at column name.
func (h *HasTimestamps) GetUpdatedAtColumn() string {
	return h.updatedAt
}

// SetCreatedAtColumn sets the created_at column name.
func (h *HasTimestamps) SetCreatedAtColumn(column string) {
	h.createdAt = column
}

// SetUpdatedAtColumn sets the updated_at column name.
func (h *HasTimestamps) SetUpdatedAtColumn(column string) {
	h.updatedAt = column
}

// FreshTimestamp returns the current time for timestamp columns.
func (h *HasTimestamps) FreshTimestamp() time.Time {
	return time.Now()
}

// FreshTimestampString returns the current time formatted as a string.
func (h *HasTimestamps) FreshTimestampString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// UpdateTimestamps sets the timestamp attributes on the given attribute map.
func (h *HasTimestamps) UpdateTimestamps(attrs map[string]any, creating bool) {
	now := h.FreshTimestampString()
	if h.updatedAt != "" {
		attrs[h.updatedAt] = now
	}
	if creating && h.createdAt != "" {
		attrs[h.createdAt] = now
	}
}
