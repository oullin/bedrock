package process

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	cprocess "github.com/bedrock/packages/contracts/process"
)

// Manager creates, runs, fakes, and asserts process commands.
type Manager struct {
	mu           sync.Mutex
	fakes        []*fakeRule
	preventStray bool
	invocations  []Command
}

// New creates a process manager.

// Command starts configuring a command.

// Run executes a command.

// Start starts a command asynchronously.

// Fake registers a fake result for matching commands. Pattern may be an exact
// command string or a filepath.Match glob such as "php *".

// Sequence registers ordered fake results for a command pattern.

// PreventStrayProcesses prevents commands without a matching fake from running.

// AssertRan verifies that a command was invoked.

// AssertNothingRan verifies that no process command was invoked.

type fakeRule struct {
	pattern  string
	results  []*Result
	sequence bool
}

var _ cprocess.Runner = (*Manager)(nil)

func New() *Manager {
	return &Manager{}
}

func (m *Manager) Command(command Command) *PendingProcess {
	return &PendingProcess{manager: m, command: command, env: map[string]string{}}
}

func (m *Manager) Run(ctx context.Context, command Command) (cprocess.Result, error) {
	return m.Command(command).Run(ctx)
}

func (m *Manager) Start(ctx context.Context, command Command) (cprocess.InvokedProcess, error) {
	return m.Command(command).Start(ctx)
}

func (m *Manager) Fake(pattern string, results ...*Result) *Manager {
	if len(results) == 0 {
		results = []*Result{NewResult(Command{}, 0, "", "")}
	}

	m.mu.Lock()

	defer m.mu.Unlock()

	m.fakes = append(m.fakes, &fakeRule{pattern: pattern, results: append([]*Result(nil), results...)})

	return m
}

func (m *Manager) Sequence(pattern string, results ...*Result) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.fakes = append(m.fakes, &fakeRule{pattern: pattern, results: append([]*Result(nil), results...), sequence: true})

	return m
}

func (m *Manager) PreventStrayProcesses() cprocess.Runner {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.preventStray = true

	return m
}

func (m *Manager) AssertRan(command Command) error {
	m.mu.Lock()

	defer m.mu.Unlock()

	want := command.String()

	for _, invocation := range m.invocations {
		if invocation.String() == want {
			return nil
		}
	}

	return fmt.Errorf("process: expected command %q to run", want)
}

func (m *Manager) AssertNothingRan() error {
	m.mu.Lock()

	defer m.mu.Unlock()

	if len(m.invocations) == 0 {
		return nil
	}

	return fmt.Errorf("process: expected nothing to run, got %d command(s)", len(m.invocations))
}

func (m *Manager) record(command Command) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.invocations = append(m.invocations, command)
}

func (m *Manager) fake(command Command) (*Result, bool, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	for _, rule := range m.fakes {
		if !rule.matches(command.String()) {
			continue
		}

		if len(rule.results) == 0 {
			return nil, true, ErrSequenceEmpty
		}

		result := *rule.results[0]
		result.command = command

		if rule.sequence {
			rule.results = rule.results[1:]
		}

		return &result, true, nil
	}

	if m.preventStray {
		return nil, false, ErrStrayProcess
	}

	return nil, false, nil
}

func (r *fakeRule) matches(command string) bool {
	if r.pattern == command {
		return true
	}

	if ok, _ := filepath.Match(r.pattern, command); ok {
		return true
	}

	if strings.HasSuffix(r.pattern, "*") {
		return strings.HasPrefix(command, strings.TrimSuffix(r.pattern, "*"))
	}

	return false
}

func throwIfRequested(result *Result, shouldThrow bool, err error) (*Result, error) {
	if err != nil {
		return result, err
	}

	if shouldThrow && result != nil && result.Failed() {
		return result, result.Throw()
	}

	return result, nil
}

func timedOut(ctx context.Context) bool {
	return errors.Is(ctx.Err(), context.DeadlineExceeded)
}
