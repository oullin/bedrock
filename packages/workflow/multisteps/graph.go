package multisteps

// compiledGraph captures the dependency edges and topological order of a workflow.
type compiledGraph struct {
	jobs     []JobSpec
	byName   map[string]JobSpec
	indexOf  map[string]int
	parents  map[string][]string
	children map[string][]string
	roots    []string
}

func compileGraph(workflowName string, jobs []JobSpec) (*compiledGraph, error) {
	if len(jobs) == 0 {
		return nil, &CompileError{Reason: "workflow " + workflowName + " declares no jobs"}
	}

	byName := make(map[string]JobSpec, len(jobs))
	indexOf := make(map[string]int, len(jobs))
	parents := make(map[string][]string, len(jobs))
	children := make(map[string][]string, len(jobs))

	for i, spec := range jobs {
		if spec.name == "" {
			return nil, &CompileError{Reason: "job at index has empty name"}
		}

		if spec.handler == nil {
			return nil, &CompileError{Reason: "job " + spec.name + " has no handler"}
		}

		if _, exists := byName[spec.name]; exists {
			return nil, &CompileError{Reason: "duplicate job name: " + spec.name}
		}

		byName[spec.name] = spec
		indexOf[spec.name] = i
	}

	for _, spec := range jobs {
		deps := map[string]struct{}{}

		for _, arg := range spec.args {
			if dep := arg.DependsOnJob(); dep != "" {
				deps[dep] = struct{}{}
			}
		}

		for _, dep := range spec.dependsOn {
			deps[dep] = struct{}{}
		}

		for dep := range deps {
			if _, ok := byName[dep]; !ok {
				return nil, &CompileError{Reason: "job " + spec.name + " depends on unknown job " + dep}
			}

			parents[spec.name] = append(parents[spec.name], dep)
			children[dep] = append(children[dep], spec.name)
		}
	}

	roots := make([]string, 0)

	for _, spec := range jobs {
		if len(parents[spec.name]) == 0 {
			roots = append(roots, spec.name)
		}
	}

	g := &compiledGraph{
		jobs:     jobs,
		byName:   byName,
		indexOf:  indexOf,
		parents:  parents,
		children: children,
		roots:    roots,
	}

	if err := g.detectCycle(); err != nil {
		return nil, err
	}

	return g, nil
}

// detectCycle returns a CompileError if the graph contains any cycle.
func (g *compiledGraph) detectCycle() error {
	const (
		unvisited = 0
		visiting  = 1
		done      = 2
	)

	state := make(map[string]int, len(g.jobs))

	var visit func(name string) error

	visit = func(name string) error {
		switch state[name] {
		case visiting:
			return &CompileError{Reason: "cycle detected involving job " + name}
		case done:
			return nil
		}

		state[name] = visiting

		for _, child := range g.children[name] {
			if err := visit(child); err != nil {
				return err
			}
		}

		state[name] = done

		return nil
	}

	for _, spec := range g.jobs {
		if err := visit(spec.name); err != nil {
			return err
		}
	}

	return nil
}
