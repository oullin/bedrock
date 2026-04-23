package process

import (
	"bytes"
	"sync"
)

// safeBuffer is an io.Writer whose contents can be safely read while
// another goroutine is writing. Used to capture live stdout/stderr
// from exec.Cmd.
type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *safeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()

	defer b.mu.Unlock()

	return b.buf.Write(p)
}

func (b *safeBuffer) String() string {
	b.mu.Lock()

	defer b.mu.Unlock()

	return b.buf.String()
}
