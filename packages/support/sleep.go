package support

import (
	"sync"
	"testing"
	"time"
)

// sleepFn is the underlying sleep function, injectable for testing.
var (
	sleepMu sync.Mutex
	sleepFn  = time.Sleep
)

// Sleep pauses the current goroutine for the given duration.
// In tests, the sleep can be intercepted via FakeSleepWith.
// Mirrors Sleep::for()->seconds() etc. (simplified to a single function).
func Sleep(d time.Duration) {
	sleepMu.Lock()
	fn := sleepFn
	sleepMu.Unlock()
	fn(d)
}

// SleepUntil pauses until the given time.
func SleepUntil(t time.Time) {
	d := time.Until(t)
	if d > 0 {
		Sleep(d)
	}
}

// FakeSleep records sleep calls instead of actually sleeping.
// Use FakeSleepWith to install it as the active sleep implementation.
// Mirrors Sleep::fake().
type FakeSleep struct {
	mu    sync.Mutex
	calls []time.Duration
}

// Sleep records the duration without actually sleeping.
func (f *FakeSleep) Sleep(d time.Duration) {
	f.mu.Lock()
	f.calls = append(f.calls, d)
	f.mu.Unlock()
}

// TotalSlept returns the total duration across all recorded sleep calls.
func (f *FakeSleep) TotalSlept() time.Duration {
	f.mu.Lock()
	defer f.mu.Unlock()
	var total time.Duration
	for _, d := range f.calls {
		total += d
	}
	return total
}

// SleptTimes returns a copy of all recorded sleep durations.
func (f *FakeSleep) SleptTimes() []time.Duration {
	f.mu.Lock()
	defer f.mu.Unlock()
	result := make([]time.Duration, len(f.calls))
	copy(result, f.calls)
	return result
}

// AssertSlept asserts that Sleep was called the given number of times with the given duration.
func (f *FakeSleep) AssertSlept(t *testing.T, d time.Duration, times int) {
	t.Helper()
	count := 0
	for _, call := range f.SleptTimes() {
		if call == d {
			count++
		}
	}
	if count != times {
		t.Errorf("expected %d sleep call(s) of %v, got %d", times, d, count)
	}
}

// AssertNeverSlept asserts that Sleep was never called.
func (f *FakeSleep) AssertNeverSlept(t *testing.T) {
	t.Helper()
	if len(f.SleptTimes()) > 0 {
		t.Errorf("expected no sleep calls, but got %d", len(f.SleptTimes()))
	}
}

// AssertSleptAtLeast asserts that the total sleep was at least d.
func (f *FakeSleep) AssertSleptAtLeast(t *testing.T, d time.Duration) {
	t.Helper()
	total := f.TotalSlept()
	if total < d {
		t.Errorf("expected total sleep of at least %v, got %v", d, total)
	}
}

// AssertSleptTimes asserts that Sleep was called exactly n times.
func (f *FakeSleep) AssertSleptTimes(t *testing.T, n int) {
	t.Helper()
	actual := len(f.SleptTimes())
	if actual != n {
		t.Errorf("expected %d sleep calls, got %d", n, actual)
	}
}

// AssertSequence asserts that Sleep was called with the given sequence of durations.
func (f *FakeSleep) AssertSequence(t *testing.T, expected []time.Duration) {
	t.Helper()
	actual := f.SleptTimes()
	if len(actual) != len(expected) {
		t.Errorf("expected sleep sequence of length %d, got %d", len(expected), len(actual))
		return
	}
	for i, d := range expected {
		if actual[i] != d {
			t.Errorf("sleep call %d: expected %v, got %v", i, d, actual[i])
		}
	}
}

// FakeSleepWith installs the given FakeSleep as the active sleep implementation.
// Returns a cleanup function that restores normal sleep behaviour.
// Call via defer: defer FakeSleepWith(fake)()
// Mirrors Sleep::fake().
func FakeSleepWith(f *FakeSleep) func() {
	sleepMu.Lock()
	prev := sleepFn
	sleepFn = f.Sleep
	sleepMu.Unlock()

	return func() {
		sleepMu.Lock()
		sleepFn = prev
		sleepMu.Unlock()
	}
}
