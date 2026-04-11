package contracts

import "time"

// Clock abstracts wall-clock time for testing.
type Clock interface {
	Now() time.Time
}
