package support

import "time"

// Timebox executes the given callback and ensures a minimum total duration.
// If the callback completes before minDuration, the remaining time is spent sleeping.
// This is used to prevent timing-based side-channel attacks (e.g. in authentication).
// Mirrors @bedrock\Support\Timebox.
func Timebox(minDuration time.Duration, fn func()) time.Duration {
	start := time.Now()
	fn()
	elapsed := time.Since(start)

	if elapsed < minDuration {
		remaining := minDuration - elapsed
		Sleep(remaining)

		return minDuration
	}

	return elapsed
}

// TimeboxWithError executes the callback and ensures minimum duration,
// returning any error from the callback.
func TimeboxWithError(minDuration time.Duration, fn func() error) (time.Duration, error) {
	start := time.Now()
	err := fn()
	elapsed := time.Since(start)

	if elapsed < minDuration {
		remaining := minDuration - elapsed
		Sleep(remaining)

		return minDuration, err
	}

	return elapsed, err
}
