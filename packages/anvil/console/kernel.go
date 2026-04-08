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
	Aliases []string
	Hidden  bool
	Help    string
	Usage   []string
	Run     func(context.Context, *Invocation) error
}

// Kernel stores and executes console commands.
type Kernel struct {
	commands map[string]Command
	aliases  map[string]string
}

// New creates a new console kernel.
func New() *Kernel {
	return &Kernel{
		commands: map[string]Command{},
		aliases:  map[string]string{},
	}
}

// Register adds a command to the kernel.
func (k *Kernel) Register(cmd Command) {
	k.commands[cmd.Name] = cmd

	for _, alias := range cmd.Aliases {
		k.aliases[alias] = cmd.Name
	}
}

// Command registers a named command.
func (k *Kernel) Command(name string, purpose string, run func(context.Context, *Invocation) error) {
	k.Register(Command{
		Name:    name,
		Purpose: purpose,
		Run:     run,
	})
}

// Commands returns all visible commands sorted by name.
func (k *Kernel) Commands() []Command {
	names := make([]string, 0, len(k.commands))

	for name, cmd := range k.commands {
		if !cmd.Hidden {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	commands := make([]Command, 0, len(names))
	for _, name := range names {
		commands = append(commands, k.commands[name])
	}

	return commands
}

// resolve looks up a command by name or alias.
func (k *Kernel) resolve(name string) (Command, bool) {
	if cmd, ok := k.commands[name]; ok {
		return cmd, true
	}

	if canonical, ok := k.aliases[name]; ok {
		cmd, ok := k.commands[canonical]
		return cmd, ok
	}

	return Command{}, false
}

// Call invokes a registered command by name and returns its exit code.
func (k *Kernel) Call(name string, args []string, stdout io.Writer, stderr io.Writer) int {
	return k.Handle(append([]string{name}, args...), stdout, stderr)
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

	command, ok := k.resolve(args[0])
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
