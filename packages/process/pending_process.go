package process

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"time"

	cprocess "github.com/bedrock/packages/contracts/process"
)

// PendingProcess configures and executes one command.
type PendingProcess struct {
	manager     *Manager
	command     Command
	dir         string
	env         map[string]string
	input       string
	timeout     time.Duration
	quiet       bool
	shouldThrow bool
}

var _ cprocess.PendingProcess = (*PendingProcess)(nil)

// Path sets the working directory.
func (p *PendingProcess) Path(path string) cprocess.PendingProcess {
	p.dir = path

	return p
}

// Env adds environment variables.
func (p *PendingProcess) Env(env map[string]string) cprocess.PendingProcess {
	for key, value := range env {
		p.env[key] = value
	}

	return p
}

// Input sets stdin for the process.
func (p *PendingProcess) Input(input string) cprocess.PendingProcess {
	p.input = input

	return p
}

// Timeout sets a command timeout.
func (p *PendingProcess) Timeout(timeout time.Duration) cprocess.PendingProcess {
	p.timeout = timeout

	return p
}

// Quietly discards process output in the returned result.
func (p *PendingProcess) Quietly() cprocess.PendingProcess {
	p.quiet = true

	return p
}

// Throw makes Run return ErrProcessFailed for non-zero exit codes.
func (p *PendingProcess) Throw() cprocess.PendingProcess {
	p.shouldThrow = true

	return p
}

// Run executes the process and waits for completion.
func (p *PendingProcess) Run(ctx context.Context) (cprocess.Result, error) {
	if p.manager == nil {
		p.manager = New()
	}

	p.manager.record(p.command)

	if result, matched, err := p.manager.fake(p.command); matched || err != nil {
		if p.quiet && result != nil {
			result.output = ""
			result.errorOutput = ""
		}

		return throwIfRequested(result, p.shouldThrow, err)
	}

	runCtx := ctx
	cancel := func() {}

	if p.timeout > 0 {
		runCtx, cancel = context.WithTimeout(ctx, p.timeout)
	}

	defer cancel()

	cmd := p.execCommand(runCtx)

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	if !p.quiet {
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}

	if p.input != "" {
		cmd.Stdin = bytes.NewBufferString(p.input)
	}

	err := cmd.Run()
	exitCode := 0

	if err != nil {
		exitCode = 1

		var exitErr *exec.ExitError

		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	}

	result := NewResult(p.command, exitCode, stdout.String(), stderr.String())

	if timedOut(runCtx) {
		return result, ErrProcessTimedOut
	}

	return throwIfRequested(result, p.shouldThrow, nil)
}

// Start starts the process without waiting for completion.
func (p *PendingProcess) Start(ctx context.Context) (cprocess.InvokedProcess, error) {
	if p.manager == nil {
		p.manager = New()
	}

	p.manager.record(p.command)

	if result, matched, err := p.manager.fake(p.command); matched || err != nil {
		return newFakeInvokedProcess(result, err), nil
	}

	runCtx := ctx
	cancel := func() {}

	if p.timeout > 0 {
		runCtx, cancel = context.WithTimeout(ctx, p.timeout)
	}

	cmd := p.execCommand(runCtx)
	invoked := newInvokedProcess(p.command, cmd, cancel)

	if p.input != "" {
		cmd.Stdin = bytes.NewBufferString(p.input)
	}

	if err := cmd.Start(); err != nil {
		cancel()

		return nil, err
	}

	go invoked.waitInternal(runCtx)

	return invoked, nil
}

func (p *PendingProcess) execCommand(ctx context.Context) *exec.Cmd {
	var cmd *exec.Cmd

	if p.command.Shell != "" {
		if runtime.GOOS == "windows" {
			cmd = exec.CommandContext(ctx, "cmd", "/C", p.command.Shell)
		} else {
			cmd = exec.CommandContext(ctx, "sh", "-c", p.command.Shell)
		}
	} else {
		cmd = exec.CommandContext(ctx, p.command.Name, p.command.Args...)
	}

	if p.dir != "" {
		cmd.Dir = p.dir
	}

	if len(p.env) > 0 {
		cmd.Env = os.Environ()

		for key, value := range p.env {
			cmd.Env = append(cmd.Env, key+"="+value)
		}
	}

	return cmd
}
