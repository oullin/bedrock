package billing

import "time"

// Currency represents an ISO 4217 currency with minor unit information.
type Currency struct {
	ID          int64
	UUID        string
	Code        string // ISO 4217 code (e.g. "USD").
	NumericCode int
	MinorUnit   int // Number of decimal places (e.g. 2 for USD).
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
