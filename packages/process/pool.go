package process

import (
	"context"
	"fmt"

	cprocess "github.com/bedrock/packages/contracts/process"
)

type poolEntry struct {
	name    string
	command Command
}

// Pool runs multiple commands and returns their results by name.
type Pool struct {
	manager *Manager
	entries []poolEntry
}

// Pool creates a process pool.
func (m *Manager) Pool(commands ...Command) *Pool {
	p := &Pool{manager: m}

	for i, command := range commands {
		p.entries = append(p.entries, poolEntry{name: fmt.Sprintf("%d", i), command: command})
	}

	return p
}

// Command adds a named command to the pool.
func (p *Pool) Command(name string, command Command) *Pool {
	p.entries = append(p.entries, poolEntry{name: name, command: command})

	return p
}

// Run runs the pool.
func (p *Pool) Run(ctx context.Context) (map[string]cprocess.Result, error) {
	results := make(map[string]cprocess.Result, len(p.entries))

	for _, entry := range p.entries {
		result, err := p.manager.Run(ctx, entry.command)
		results[entry.name] = result

		if err != nil {
			return results, err
		}
	}

	return results, nil
}
