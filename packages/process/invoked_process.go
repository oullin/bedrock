package process

import (
	"context"
	"errors"
	"os/exec"
	"sync"
	"time"

	cprocess "github.com/bedrock/packages/contracts/process"
)

// InvokedProcess is a started process.
type InvokedProcess struct {
	command Command
	cmd     *exec.Cmd
	cancel  func()

	stdout safeBuffer
	stderr safeBuffer

	once   sync.Once
	done   chan struct{}
	result *Result
	err    error
}

var _ cprocess.InvokedProcess = (*InvokedProcess)(nil)

func newInvokedProcess(command Command, cmd *exec.Cmd, cancel func()) *InvokedProcess {
	p := &InvokedProcess{command: command, cmd: cmd, cancel: cancel, done: make(chan struct{})}
	cmd.Stdout = &p.stdout
	cmd.Stderr = &p.stderr

	return p
}

func newFakeInvokedProcess(result *Result, err error) *InvokedProcess {
	if result == nil {
		result = NewResult(Command{}, 0, "", "")
	}

	done := make(chan struct{})
	close(done)

	return &InvokedProcess{command: result.command, done: done, result: result, err: err}
}

func (p *InvokedProcess) waitInternal(ctx context.Context) {
	p.once.Do(func() {
		defer close(p.done)

		defer p.cancel()

		err := p.cmd.Wait()
		exitCode := 0

		if err != nil {
			exitCode = 1

			var exitErr *exec.ExitError

			if errors.As(err, &exitErr) {
				exitCode = exitErr.ExitCode()
			}
		}

		p.result = NewResult(p.command, exitCode, p.stdout.String(), p.stderr.String())

		if timedOut(ctx) {
			p.err = ErrProcessTimedOut

			return
		}

		p.err = err
	})
}

// Wait waits for the process to finish.
func (p *InvokedProcess) Wait(ctx context.Context) (cprocess.Result, error) {
	select {
	case <-p.done:
		return p.result, p.err
	case <-ctx.Done():
		if p.cancel != nil {
			p.cancel()
		}

		<-p.done

		return p.result, ctx.Err()
	}
}

// WaitUntil waits until the callback matches captured output or the process
// exits.
func (p *InvokedProcess) WaitUntil(ctx context.Context, fn func(output string, errorOutput string) bool) (cprocess.Result, error) {
	ticker := time.NewTicker(10 * time.Millisecond)

	defer ticker.Stop()

	for {
		if fn(p.Output(), p.ErrorOutput()) {
			return p.Wait(ctx)
		}

		select {
		case <-p.done:
			return p.result, p.err
		case <-ctx.Done():
			if p.cancel != nil {
				p.cancel()
			}

			<-p.done

			return p.result, ctx.Err()
		case <-ticker.C:
		}
	}
}

// Output returns all captured stdout.
func (p *InvokedProcess) Output() string {
	if p.cmd == nil {
		return p.result.Output()
	}

	return p.stdout.String()
}

// LatestOutput returns the latest captured stdout.
func (p *InvokedProcess) LatestOutput() string {
	return p.Output()
}

// ErrorOutput returns all captured stderr.
func (p *InvokedProcess) ErrorOutput() string {
	if p.cmd == nil {
		return p.result.ErrorOutput()
	}

	return p.stderr.String()
}
