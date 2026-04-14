package queue_test

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
)

// Full port of Illuminate\Tests\Queue\QueueListenerTest (5 / 5).
//
// Laravel's test uses Mockery to partial-mock Symfony's Process and
// the Listener itself, then asserts that makeProcess() returns a
// process with the right command-line string, working directory, and
// timeout. The Go equivalent:
//
//   - substitutes Symfony Process with the package-local ProcessRunner
//     interface (see listener.go);
//   - uses a recording fakeListenerProcess fixture for the two
//     runProcess tests (drop-in replacement for the Mockery mock);
//   - asserts on the []string command parts instead of a shell-
//     escaped string — same observable shape, no escaping quirks.

// --- fixtures ---------------------------------------------------------

type fakeListenerProcess struct {
	mu       sync.Mutex
	runCalls int
	runErr   error
}

func (p *fakeListenerProcess) Run() error {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.runCalls++

	return p.runErr
}

func (p *fakeListenerProcess) CommandParts() []string   { return nil }
func (p *fakeListenerProcess) WorkingDirectory() string { return "" }
func (p *fakeListenerProcess) Timeout() time.Duration   { return 0 }

// --- ports ------------------------------------------------------------

// Port of Illuminate\Tests\Queue\QueueListenerTest::testRunProcessCallsProcess
func TestRunProcessCallsProcess(t *testing.T) {
	t.Parallel()

	listener := queue.NewListener(".")
	listener.MemoryExceededFunc = func(limit int) bool {
		if limit != 1 {
			t.Errorf("memoryExceeded called with %d, want 1", limit)
		}

		return false
	}

	stopCalled := false
	listener.StopFunc = func() { stopCalled = true }

	process := &fakeListenerProcess{}

	if err := listener.RunProcess(process, 1); err != nil {
		t.Fatalf("RunProcess: %v", err)
	}

	if process.runCalls != 1 {
		t.Errorf("process.Run calls: got %d, want 1", process.runCalls)
	}

	if stopCalled {
		t.Error("expected Stop NOT to be called when memory is under the limit")
	}
}

// Port of Illuminate\Tests\Queue\QueueListenerTest::testListenerStopsWhenMemoryIsExceeded
func TestListenerStopsWhenMemoryIsExceeded(t *testing.T) {
	t.Parallel()

	listener := queue.NewListener(".")
	listener.MemoryExceededFunc = func(limit int) bool {
		if limit != 1 {
			t.Errorf("memoryExceeded called with %d, want 1", limit)
		}

		return true
	}

	stopCalled := false
	listener.StopFunc = func() { stopCalled = true }

	process := &fakeListenerProcess{}

	if err := listener.RunProcess(process, 1); err != nil {
		t.Fatalf("RunProcess: %v", err)
	}

	if process.runCalls != 1 {
		t.Errorf("process.Run calls: got %d, want 1", process.runCalls)
	}

	if !stopCalled {
		t.Error("expected Stop to be called when memory is exceeded")
	}
}

// Port of Illuminate\Tests\Queue\QueueListenerTest::testMakeProcessCorrectlyFormatsCommandLine
func TestMakeProcessCorrectlyFormatsCommandLine(t *testing.T) {
	t.Parallel()

	workdir := "/tmp/bedrock-listener-test"
	listener := queue.NewListener(workdir)

	opts := queue.NewListenerOptions()
	opts.Backoff = 1
	opts.Memory = 2
	opts.Timeout = 3 * time.Second

	process := listener.MakeProcess("connection", "queue", opts)

	if process.WorkingDirectory() != workdir {
		t.Errorf("WorkingDirectory: got %q, want %q", process.WorkingDirectory(), workdir)
	}

	if process.Timeout() != 3*time.Second {
		t.Errorf("Timeout: got %s, want 3s", process.Timeout())
	}

	// Laravel asserts a shell-escaped single-line string. The Go port
	// asserts the argv slice directly — same observable command shape
	// without shell-escaping quirks.
	want := []string{
		"php",
		"artisan",
		"queue:work",
		"connection",
		"--once",
		"--name=default",
		"--queue=queue",
		"--backoff=1",
		"--memory=2",
		"--sleep=3",
		"--tries=1",
	}

	if got := process.CommandParts(); !reflect.DeepEqual(got, want) {
		t.Errorf("CommandParts:\n got  %v\n want %v", got, want)
	}
}

// Port of Illuminate\Tests\Queue\QueueListenerTest::testMakeProcessCorrectlyFormatsCommandLineWithAnEnvironmentSpecified
func TestMakeProcessCorrectlyFormatsCommandLineWithAnEnvironmentSpecified(t *testing.T) {
	t.Parallel()

	workdir := "/tmp/bedrock-listener-test"
	listener := queue.NewListener(workdir)

	opts := queue.NewListenerOptionsWithEnv("default", "test")
	opts.Backoff = 1
	opts.Memory = 2
	opts.Timeout = 3 * time.Second

	process := listener.MakeProcess("connection", "queue", opts)

	if process.WorkingDirectory() != workdir {
		t.Errorf("WorkingDirectory: got %q, want %q", process.WorkingDirectory(), workdir)
	}

	if process.Timeout() != 3*time.Second {
		t.Errorf("Timeout: got %s, want 3s", process.Timeout())
	}

	want := []string{
		"php",
		"artisan",
		"queue:work",
		"connection",
		"--once",
		"--name=default",
		"--queue=queue",
		"--backoff=1",
		"--memory=2",
		"--sleep=3",
		"--tries=1",
		"--env=test",
	}

	if got := process.CommandParts(); !reflect.DeepEqual(got, want) {
		t.Errorf("CommandParts:\n got  %v\n want %v", got, want)
	}
}

// Port of Illuminate\Tests\Queue\QueueListenerTest::testMakeProcessCorrectlyFormatsCommandLineWhenTheConnectionIsNotSpecified
func TestMakeProcessCorrectlyFormatsCommandLineWhenTheConnectionIsNotSpecified(t *testing.T) {
	t.Parallel()

	workdir := "/tmp/bedrock-listener-test"
	listener := queue.NewListener(workdir)

	opts := queue.NewListenerOptionsWithEnv("default", "test")
	opts.Backoff = 1
	opts.Memory = 2
	opts.Timeout = 3 * time.Second

	process := listener.MakeProcess("", "queue", opts)

	if process.WorkingDirectory() != workdir {
		t.Errorf("WorkingDirectory: got %q, want %q", process.WorkingDirectory(), workdir)
	}

	if process.Timeout() != 3*time.Second {
		t.Errorf("Timeout: got %s, want 3s", process.Timeout())
	}

	// When connection is empty the argv slot is omitted — matching
	// Laravel's array_filter pass over the command list.
	want := []string{
		"php",
		"artisan",
		"queue:work",
		"--once",
		"--name=default",
		"--queue=queue",
		"--backoff=1",
		"--memory=2",
		"--sleep=3",
		"--tries=1",
		"--env=test",
	}

	if got := process.CommandParts(); !reflect.DeepEqual(got, want) {
		t.Errorf("CommandParts:\n got  %v\n want %v", got, want)
	}
}
