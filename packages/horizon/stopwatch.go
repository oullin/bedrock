package horizon

import "time"

// Stopwatch measures durations between checks.
type Stopwatch struct {
	last time.Time
}

// NewStopwatch starts a stopwatch at the given time.
func NewStopwatch(start time.Time) *Stopwatch {
	return &Stopwatch{last: start}
}

// Check returns the elapsed time since the previous check and advances the clock.
func (s *Stopwatch) Check(now time.Time) time.Duration {
	elapsed := now.Sub(s.last)
	s.last = now

	return elapsed
}
