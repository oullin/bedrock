package notifications

// DatabaseMessage represents a notification payload to be stored in the
// database.
type DatabaseMessage struct {
	// Data is the notification data to persist.
	Data map[string]any
}

// NewDatabaseMessage creates a DatabaseMessage with the given data.
func NewDatabaseMessage(data ...map[string]any) *DatabaseMessage {
	m := &DatabaseMessage{}

	if len(data) > 0 {
		m.Data = data[0]
	}

	return m
}
