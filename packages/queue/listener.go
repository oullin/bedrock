package queue

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

// ProcessRunner is the narrow contract the Listener needs in order to
// invoke a worker subprocess. It is the Go analogue of Symfony's
// Process::run, pared down to the handful of methods the Listener and
// its tests actually consume. Production callers use SimpleProcess
// (backed by os/exec); tests supply their own recording fake.
type ProcessRunner interface {
	// Run executes the command and blocks until it exits.
	Run() error
	// CommandParts returns the argv slice the Listener built for this
	// process. Tests assert on this to verify the command shape.
	CommandParts() []string
	// WorkingDirectory returns the directory the process will run in.
	WorkingDirectory() string
	// Timeout returns the run-time ceiling for this process. Zero
	// means no timeout.
	Timeout() time.Duration
}

// SimpleProcess is the default ProcessRunner implementation. It stores
// the command, working directory, and timeout supplied by the Listener
// and shells out via os/exec when Run is called. Tests that need to
// observe calls use a recording fake instead of this type.
type SimpleProcess struct {
	Parts   []string
	Cwd     string
	TimeOut time.Duration
}

// Run executes the stored command under os/exec, honouring the
// timeout. A zero Timeout means no deadline.

// CommandParts returns the argv slice.

// WorkingDirectory returns the process working directory.

// Timeout returns the configured run-time ceiling.

// ListenerOptions configures the Listener's outer loop and the
// Ref: @bedrock/code-0265
// Go idiom (MaxTries instead of maxTries, Rest instead of $rest).
type ListenerOptions struct {
	// Name is the worker process name passed via --name. Defaults to
	// "default" when constructed via NewListenerOptions.
	Name string
	// Environment, if non-empty, is appended to the command as
	// --env={environment}.
	Environment string
	// Backoff seconds between retries, forwarded to the worker.
	Backoff int
	// Memory is the worker memory cap in MiB.
	Memory int
	// Timeout is the maximum run time per spawned subprocess.
	Timeout time.Duration
	// Sleep is the worker's inter-poll sleep in seconds. Default 3.
	Sleep int
	// MaxTries is the worker's attempt cap. Default 1.
	MaxTries int
	// Rest is the sleep (in seconds) between successive runProcess
	// calls in the outer loop.
	Rest int
	// Force adds --force to the worker command when true.
	Force bool
}

// NewListenerOptions returns a ListenerOptions initialised with the
// Upstream defaults: Name="default", Sleep=3, MaxTries=1.

// NewListenerOptionsWithEnv is the two-argument constructor that
// matches the upstream `new ListenerOptions($name, $environment)` form.

// Listener spawns and supervises worker subprocesses. It is the Go
// Ref: @bedrock/code-0264
// The Listener is deliberately transport-agnostic: it builds a command
// slice from the caller-supplied connection/queue/options tuple and
// hands it to a ProcessRunner. The default ProcessRunner (SimpleProcess)
// shells out via os/exec; tests swap in a recording fake. Similarly,
// MemoryExceeded and Stop are func hooks with sane defaults so tests
// can force either branch without touching the host process.
type Listener struct {
	commandPath string
	// WorkerBinary is the binary path invoked as the first argv slot.
	// For a PHP/upstream port it is "php"; for a Go-native port it
	// might be "go" + "run ./cmd/queue-worker" or a compiled binary.
	WorkerBinary string
	// EntryArg is the second argv slot, typically the cli script
	// or an equivalent subcommand dispatcher. Set to empty string to
	// omit entirely.
	EntryArg string
	// ProcessFactory builds a ProcessRunner for a given command. The
	// default returns a *SimpleProcess; tests can override to capture
	// construction arguments.
	ProcessFactory func(parts []string, cwd string, timeout time.Duration) ProcessRunner
	// MemoryExceededFunc overrides the default runtime.MemStats.Sys
	// check. Tests set this to return deterministic values.
	MemoryExceededFunc func(limitMiB int) bool
	// StopFunc is the callback invoked by Stop. Default is os.Exit(0);
	// tests override it to record the call without terminating the
	// test binary.
	StopFunc func()
	// OutputHandler is called for every stdout/stderr line produced
	// by a running worker subprocess.
	OutputHandler func(stream, line string)
}

func (p *SimpleProcess) Run() error {
	if len(p.Parts) == 0 {
		return fmt.Errorf("queue: SimpleProcess.Run: empty command")
	}

	cmd := exec.Command(p.Parts[0], p.Parts[1:]...)
	cmd.Dir = p.Cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if p.TimeOut == 0 {
		return cmd.Run()
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan error, 1)

	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		return err
	case <-time.After(p.TimeOut):
		_ = cmd.Process.Kill()

		return fmt.Errorf("queue: SimpleProcess.Run: timeout after %s", p.TimeOut)
	}
}

func (p *SimpleProcess) CommandParts() []string { return p.Parts }

func (p *SimpleProcess) WorkingDirectory() string { return p.Cwd }

func (p *SimpleProcess) Timeout() time.Duration { return p.TimeOut }

func NewListenerOptions() ListenerOptions {
	return ListenerOptions{Name: "default", Sleep: 3, MaxTries: 1}
}

func NewListenerOptionsWithEnv(name, environment string) ListenerOptions {
	opts := NewListenerOptions()
	opts.Name = name
	opts.Environment = environment

	return opts
}

// NewListener constructs a Listener rooted at commandPath. The worker
// binary defaults to "php" and the entry arg to "cli" so the
// command shape matches upstream 1:1. Go-native consumers should
// reassign WorkerBinary + EntryArg after construction.
func NewListener(commandPath string) *Listener {
	l := &Listener{
		commandPath:  commandPath,
		WorkerBinary: "php",
		EntryArg:     "cli",
	}

	l.ProcessFactory = func(parts []string, cwd string, timeout time.Duration) ProcessRunner {
		return &SimpleProcess{Parts: parts, Cwd: cwd, TimeOut: timeout}
	}

	return l
}

// CommandPath returns the root directory from which the Listener
// spawns worker processes.
func (l *Listener) CommandPath() string { return l.commandPath }

// MakeProcess builds the command for a worker subprocess and wraps it
// Ref: @bedrock/code-0264
func (l *Listener) MakeProcess(connection, queue string, opts ListenerOptions) ProcessRunner {
	cmd := l.createCommand(connection, queue, opts)

	if opts.Environment != "" {
		cmd = append(cmd, "--env="+opts.Environment)
	}

	return l.ProcessFactory(cmd, l.commandPath, opts.Timeout)
}

// createCommand builds the argv slice that MakeProcess hands to the
// Ref: @bedrock/code-0264
// the "drop null entries" step (Go: drop empty strings).
func (l *Listener) createCommand(connection, queue string, opts ListenerOptions) []string {
	name := opts.Name

	if name == "" {
		name = "default"
	}

	cmd := []string{l.WorkerBinary}

	if l.EntryArg != "" {
		cmd = append(cmd, l.EntryArg)
	}

	cmd = append(cmd, "queue:work")

	if connection != "" {
		cmd = append(cmd, connection)
	}

	cmd = append(cmd,
		"--once",
		"--name="+name,
		"--queue="+queue,
		"--backoff="+strconv.Itoa(opts.Backoff),
		"--memory="+strconv.Itoa(opts.Memory),
		"--sleep="+strconv.Itoa(opts.Sleep),
		"--tries="+strconv.Itoa(opts.MaxTries),
	)

	if opts.Force {
		cmd = append(cmd, "--force")
	}

	return cmd
}

// RunProcess runs a worker subprocess and, on return, checks the
// memory cap. If memory is exceeded, Stop is invoked — matching
// the upstream "kill the listener so the process manager restarts it"
// Ref: @bedrock/code-0264
func (l *Listener) RunProcess(process ProcessRunner, memoryLimitMiB int) error {
	if err := process.Run(); err != nil {
		return err
	}

	if l.MemoryExceeded(memoryLimitMiB) {
		l.Stop()
	}

	return nil
}

// Listen is the outer supervisor loop: call MakeProcess once, then
// RunProcess/Rest-sleep forever. Callers that need graceful shutdown
// should arrange for StopFunc to return so the loop exits.
// the upstream Listener::listen.
func (l *Listener) Listen(connection, queue string, opts ListenerOptions) error {
	process := l.MakeProcess(connection, queue, opts)

	for {
		if err := l.RunProcess(process, opts.Memory); err != nil {
			return err
		}

		if opts.Rest > 0 {
			time.Sleep(time.Duration(opts.Rest) * time.Second)
		}
	}
}

// MemoryExceeded reports whether the current process has exceeded
// limitMiB MiB. Delegates to MemoryExceededFunc when set, otherwise
// reads runtime.MemStats.Sys — the Go analogue of PHP's
// memory_get_usage(true).
func (l *Listener) MemoryExceeded(limitMiB int) bool {
	if l.MemoryExceededFunc != nil {
		return l.MemoryExceededFunc(limitMiB)
	}

	if limitMiB <= 0 {
		return false
	}

	var ms runtime.MemStats

	runtime.ReadMemStats(&ms)

	return ms.Sys >= uint64(limitMiB)*1024*1024
}

// Stop terminates the listener. Default implementation is
// os.Exit(0); tests override StopFunc to avoid killing the test binary.
func (l *Listener) Stop() {
	if l.StopFunc != nil {
		l.StopFunc()

		return
	}

	os.Exit(0)
}
