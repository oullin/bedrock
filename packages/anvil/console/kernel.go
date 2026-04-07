package console

import (
	"context"
	"fmt"
	"io"
	"sort"
)

// Invocation describes a console command invocation.
type Invocation struct {
	Args   []string
	Stdout io.Writer
	Stderr io.Writer
}

// Comment writes a user-facing console line.
func (i *Invocation) Comment(message string) {
	_, _ = fmt.Fprintln(i.Stdout, message)
}

// Command describes a registered console command.
type Command struct {
	Name    string
	Purpose string
	Run     func(context.Context, *Invocation) error
}

// Kernel stores and executes console commands.
type Kernel struct {
	commands map[string]Command
}

// New creates a new console kernel.
func New() *Kernel {
	return &Kernel{commands: map[string]Command{}}
}

// Command registers a named command.
func (k *Kernel) Command(name string, purpose string, run func(context.Context, *Invocation) error) {
	k.commands[name] = Command{
		Name:    name,
		Purpose: purpose,
		Run:     run,
	}
}

// Commands returns all commands sorted by name.
func (k *Kernel) Commands() []Command {
	names := make([]string, 0, len(k.commands))

	for name := range k.commands {
		names = append(names, name)
	}

	sort.Strings(names)

	commands := make([]Command, 0, len(names))
	for _, name := range names {
		commands = append(commands, k.commands[name])
	}

	return commands
}

// Handle executes a registered command and returns a shell exit code.
func (k *Kernel) Handle(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stdout, "available commands:")
		for _, command := range k.Commands() {
			_, _ = fmt.Fprintf(stdout, "  %s\t%s\n", command.Name, command.Purpose)
		}

		return 0
	}

	command, ok := k.commands[args[0]]
	if !ok {
		_, _ = fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		return 1
	}

	if err := command.Run(context.Background(), &Invocation{
		Args:   args[1:],
		Stdout: stdout,
		Stderr: stderr,
	}); err != nil {
		_, _ = fmt.Fprintln(stderr, err.Error())
		return 1
	}

	return 0
}
