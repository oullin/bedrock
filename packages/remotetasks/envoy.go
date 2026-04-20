package remotetasks

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"sort"
	"strings"
	"text/template"
)

// Task describes commands to run against one or more hosts.
type Task struct {
	Name      string
	Hosts     []string
	Commands  []string
	Variables map[string]any
}

// Step is one rendered command for one target host.
type Step struct {
	Task    string
	Host    string
	Command string
}

// Result is the output of running one step.
type Result struct {
	Step   Step
	Stdout string
	Stderr string
}

// Runner executes one step.
type Runner interface {
	Run(ctx context.Context, step Step) (Result, error)
}

// RunnerFunc adapts a function into a Runner.
type RunnerFunc func(ctx context.Context, step Step) (Result, error)

// Run executes step by calling f.

// Plan renders task commands for each host.

// Run plans task and executes every step with runner.

// LocalRunner runs commands on the current machine through the system shell.
type LocalRunner struct{}

func (f RunnerFunc) Run(ctx context.Context, step Step) (Result, error) {
	return f(ctx, step)
}

func Plan(task Task, variables map[string]any) ([]Step, error) {
	if task.Name == "" {
		return nil, ErrMissingTaskName
	}

	if len(task.Commands) == 0 {
		return nil, ErrNoCommands
	}

	merged := map[string]any{}

	for key, value := range task.Variables {
		merged[key] = value
	}

	for key, value := range variables {
		merged[key] = value
	}

	hosts := append([]string(nil), task.Hosts...)

	if len(hosts) == 0 {
		hosts = []string{""}
	} else {
		sort.Strings(hosts)
	}

	steps := make([]Step, 0, len(hosts)*len(task.Commands))

	for _, host := range hosts {
		for _, command := range task.Commands {
			rendered, err := render(command, merged)

			if err != nil {
				return nil, err
			}

			steps = append(steps, Step{
				Task:    task.Name,
				Host:    host,
				Command: rendered,
			})
		}
	}

	return steps, nil
}

func Run(ctx context.Context, runner Runner, task Task, variables map[string]any) ([]Result, error) {
	if runner == nil {
		return nil, ErrMissingRunner
	}

	steps, err := Plan(task, variables)

	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(steps))

	for _, step := range steps {
		result, err := runner.Run(ctx, step)

		if err != nil {
			return results, err
		}

		results = append(results, result)
	}

	return results, nil
}

// Run executes step.Command locally.
func (LocalRunner) Run(ctx context.Context, step Step) (Result, error) {
	command := exec.CommandContext(ctx, "sh", "-c", step.Command)

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()

	return Result{
		Step:   step,
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, err
}

func render(command string, variables map[string]any) (string, error) {
	tpl, err := template.New("command").Option("missingkey=error").Parse(command)

	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	if err := tpl.Execute(&output, variables); err != nil {
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}

var (
	// ErrMissingTaskName is returned when a task has no name.
	ErrMissingTaskName = errors.New("remotetasks: missing task name")
	// ErrNoCommands is returned when a task has no commands.
	ErrNoCommands = errors.New("remotetasks: no commands configured")
	// ErrMissingRunner is returned when Run has no runner.
	ErrMissingRunner = errors.New("remotetasks: missing runner")
)
