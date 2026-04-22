package process

import cprocess "github.com/bedrock/packages/contracts/process"

// Command describes a process command.
type Command = cprocess.Command

// Shell creates a shell command.
func Shell(command string) Command {
	return Command{Shell: command}
}

// Args creates a direct executable command.
func Args(name string, args ...string) Command {
	copied := append([]string(nil), args...)

	return Command{Name: name, Args: copied}
}
