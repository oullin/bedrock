package process

import (
	"context"
	"time"
)

// Result is the completed result of a process.
type Result interface {
	Successful() bool
	Failed() bool
	ExitCode() int
	Output() string
	ErrorOutput() string
	Command() Command
	Throw() error
}

// InvokedProcess is a started process that can be waited on.
type InvokedProcess interface {
	Wait(context.Context) (Result, error)
	WaitUntil(context.Context, func(output string, errorOutput string) bool) (Result, error)
	Output() string
	LatestOutput() string
	ErrorOutput() string
}

// PendingProcess configures and runs a command.
type PendingProcess interface {
	Path(string) PendingProcess
	Env(map[string]string) PendingProcess
	Input(string) PendingProcess
	Timeout(time.Duration) PendingProcess
	Quietly() PendingProcess
	Throw() PendingProcess
	Run(context.Context) (Result, error)
	Start(context.Context) (InvokedProcess, error)
}

// Runner executes process commands.
type Runner interface {
	Run(context.Context, Command) (Result, error)
	Start(context.Context, Command) (InvokedProcess, error)
	PreventStrayProcesses() Runner
	AssertRan(Command) error
	AssertNothingRan() error
}
