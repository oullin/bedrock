package process

import (
	"context"

	cprocess "github.com/bedrock/packages/contracts/process"
)

// Pipe runs commands sequentially, passing each stdout to the next stdin.
type Pipe struct {
	manager  *Manager
	commands []Command
}

// Pipe creates a process pipe.
func (m *Manager) Pipe(commands ...Command) *Pipe {
	return &Pipe{manager: m, commands: append([]Command(nil), commands...)}
}

// Run runs the command pipe.
func (p *Pipe) Run(ctx context.Context) (cprocess.Result, error) {
	var input string

	var result cprocess.Result

	for _, command := range p.commands {
		var err error

		result, err = p.manager.Command(command).Input(input).Run(ctx)

		if err != nil {
			return result, err
		}

		if result.Failed() {
			return result, result.Throw()
		}

		input = result.Output()
	}

	return result, nil
}
