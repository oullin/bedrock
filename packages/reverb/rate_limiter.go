package reverb

import "time"

// Allow reports whether the connection is within its message rate limit.
//
// It uses a sliding window: if the number of messages recorded since the
// connection's msgWindowStart exceeds threshold within the decay duration,
// Allow returns false.
//
// When the current window has expired (time.Since(msgWindowStart) > decay),
// the window is reset and the message is counted as the first in the new
// window (returns true). Otherwise the in-flight counter is incremented and
// Allow returns true only when the count is still within threshold.
func Allow(conn *Conn, threshold int64, decay time.Duration) bool {
	if time.Since(conn.MessageWindowStart()) > decay {
		conn.ResetMessageWindow(time.Now())
		conn.IncrMessageCount()

		return true
	}

	count := conn.IncrMessageCount()

	return count <= threshold
}
