package process

import cprocess "github.com/bedrock/packages/contracts/process"

// Result is the completed result of a process.
type Result struct {
	command     Command
	exitCode    int
	output      string
	errorOutput string
}

var _ cprocess.Result = (*Result)(nil)

// NewResult creates a process result.
func NewResult(command Command, exitCode int, output string, errorOutput string) *Result {
	return &Result{command: command, exitCode: exitCode, output: output, errorOutput: errorOutput}
}

// Successful reports whether the process exited with code 0.
func (r *Result) Successful() bool {
	return r != nil && r.exitCode == 0
}

// Failed reports whether the process exited unsuccessfully.
func (r *Result) Failed() bool {
	return !r.Successful()
}

// ExitCode returns the process exit code.
func (r *Result) ExitCode() int {
	if r == nil {
		return -1
	}

	return r.exitCode
}

// Output returns captured stdout.
func (r *Result) Output() string {
	if r == nil {
		return ""
	}

	return r.output
}

// ErrorOutput returns captured stderr.
func (r *Result) ErrorOutput() string {
	if r == nil {
		return ""
	}

	return r.errorOutput
}

// Command returns the process command.
func (r *Result) Command() Command {
	if r == nil {
		return Command{}
	}

	return r.command
}

// Throw returns ErrProcessFailed when the result failed.
func (r *Result) Throw() error {
	if r == nil || r.Successful() {
		return nil
	}

	return &ProcessError{result: r, err: ErrProcessFailed}
}
