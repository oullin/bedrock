package support

import (
	"errors"
	"time"
)

// Retry attempts to execute the given callback the specified number of times.
// If all attempts fail, the last error is returned.
// The sleep argument controls delay between retries and can be:
//   - int: milliseconds to sleep between attempts
//   - []int: backoff schedule in milliseconds (cycles if exhausted)
//   - func(int) time.Duration: callback receiving attempt number, returns sleep duration
//
// Mirrors Laravel's retry() helper.
func Retry(times int, fn func(attempt int) error, sleep ...any) error {
	var lastErr error

	for attempt := 1; attempt <= times; attempt++ {
		lastErr = fn(attempt)

		if lastErr == nil {
			return nil
		}

		// Don't sleep after the last attempt
		if attempt == times {
			break
		}

		if len(sleep) > 0 {
			d := resolveSleepDuration(sleep[0], attempt)

			if d > 0 {
				time.Sleep(d)
			}
		}
	}

	return lastErr
}

// RetryWhen retries the callback only when the given condition returns true.
// Mirrors Laravel's retry() with a when callback.
func RetryWhen(times int, fn func(attempt int) error, when func(error) bool, sleep ...any) error {
	var lastErr error

	for attempt := 1; attempt <= times; attempt++ {
		lastErr = fn(attempt)

		if lastErr == nil {
			return nil
		}

		if !when(lastErr) {
			return lastErr
		}

		if attempt == times {
			break
		}

		if len(sleep) > 0 {
			d := resolveSleepDuration(sleep[0], attempt)

			if d > 0 {
				time.Sleep(d)
			}
		}
	}

	return lastErr
}

func resolveSleepDuration(sleep any, attempt int) time.Duration {
	switch v := sleep.(type) {
	case int:
		return time.Duration(v) * time.Millisecond
	case []int:
		if len(v) == 0 {
			return 0
		}

		idx := attempt - 1

		if idx >= len(v) {
			idx = len(v) - 1
		}

		return time.Duration(v[idx]) * time.Millisecond
	case func(int) time.Duration:
		return v(attempt)
	case time.Duration:
		return v
	}

	return 0
}

// ThrowIf throws the given exception if the condition is true.
// Mirrors Laravel's throw_if() helper.
func ThrowIf(condition bool, err error) error {
	if condition {
		return err
	}

	return nil
}

// ThrowUnless throws the given exception unless the condition is true.
// Mirrors Laravel's throw_unless() helper.
func ThrowUnless(condition bool, err error) error {
	return ThrowIf(!condition, err)
}

// PanicIf panics with the given error if the condition is true.
func PanicIf(condition bool, err error) {
	if condition {
		panic(err)
	}
}

// PanicUnless panics with the given error unless the condition is true.
func PanicUnless(condition bool, err error) {
	PanicIf(!condition, err)
}

// Rescue executes the given callback and recovers from any panic,
// returning the error. Useful for wrapping code that may panic.
func Rescue[T any](fn func() T, def ...T) (result T, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				err = e
			case string:
				err = errors.New(e)
			default:
				err = errors.New("unknown panic")
			}

			if len(def) > 0 {
				result = def[0]
			}
		}
	}()

	return fn(), nil
}
