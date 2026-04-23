package jobqueue

import "time"

func microsecondsFloat(at time.Time) float64 {
	return float64(at.UnixNano()) / float64(time.Microsecond)
}
