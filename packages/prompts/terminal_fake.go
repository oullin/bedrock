package prompts

import (
	"fmt"
	"sync"
)

// FakeTerminal is a test double that replays pre-queued keystrokes and
// records all output. Install it via [Fake].
type FakeTerminal struct {
	mu        sync.Mutex
	keys      []string
	keyPos    int
	cols      int
	lines     int
	trueColor bool
}

var _ Terminal = (*FakeTerminal)(nil)

// NewFakeTerminal creates a FakeTerminal with the given dimensions.
func NewFakeTerminal(cols, lines int) *FakeTerminal {
	return &FakeTerminal{
		cols:  cols,
		lines: lines,
	}
}

// QueueKey enqueues one or more key inputs to be returned by Read.
func (f *FakeTerminal) QueueKey(keys ...string) {
	f.mu.Lock()
	f.keys = append(f.keys, keys...)
	f.mu.Unlock()
}

// QueueKeys enqueues a slice of key inputs.
func (f *FakeTerminal) QueueKeys(keys []string) {
	f.mu.Lock()
	f.keys = append(f.keys, keys...)
	f.mu.Unlock()
}

func (f *FakeTerminal) Read() (string, error) {
	f.mu.Lock()

	defer f.mu.Unlock()

	if f.keyPos >= len(f.keys) {
		return "", fmt.Errorf("prompts: FakeTerminal: no more queued keys (read %d of %d)", f.keyPos, len(f.keys))
	}

	key := f.keys[f.keyPos]
	f.keyPos++

	return key, nil
}

func (f *FakeTerminal) SetTty(string) error { return nil }
func (f *FakeTerminal) RestoreTty() error   { return nil }

func (f *FakeTerminal) Cols() int  { return f.cols }
func (f *FakeTerminal) Lines() int { return f.lines }

func (f *FakeTerminal) Exit() {}

func (f *FakeTerminal) SupportsTrueColor() bool { return f.trueColor }

// SetTrueColor sets whether the fake terminal claims true-color support.
func (f *FakeTerminal) SetTrueColor(v bool) {
	f.mu.Lock()
	f.trueColor = v
	f.mu.Unlock()
}

// RemainingKeys returns the count of unread queued keys.
func (f *FakeTerminal) RemainingKeys() int {
	f.mu.Lock()

	defer f.mu.Unlock()

	return len(f.keys) - f.keyPos
}
