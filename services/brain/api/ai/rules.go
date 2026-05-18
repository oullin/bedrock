package ai

import (
	"errors"
	"path/filepath"
	"sort"

	"github.com/bedrock/packages/ai/boost"
	"github.com/bedrock/packages/filesystem"
)

// RuleTarget describes one editor-rules file written by GenerateRules.
type RuleTarget struct {
	Path  string
	Label string
}

// Targets is the canonical list of editor-rules files brain writes. It is
// derived from the SupportsGuidelines agents registered in packages/ai/boost,
// deduplicated by path and sorted lexicographically so the slice is
// deterministic across runs.
var Targets = computeTargets()

func computeTargets() []RuleTarget {
	manager := boost.New()
	registered := manager.GetAgents()

	keys := make([]string, 0, len(registered))

	for k := range registered {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	seen := make(map[string]bool, len(keys))
	out := make([]RuleTarget, 0, len(keys))

	for _, k := range keys {
		agent, ok := registered[k].(boost.SupportsGuidelines)

		if !ok {
			continue
		}

		path := agent.GuidelinesPath()

		if path == "" || seen[path] {
			continue
		}

		seen[path] = true

		out = append(out, RuleTarget{Path: path, Label: agent.DisplayName()})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })

	return out
}

// ErrConflict is returned when a target file already exists and force=false.
var ErrConflict = errors.New("rule target file already exists; pass force=true to overwrite")

// GenerateRules writes every Target under `root`. If a file exists and
// force is false, no files are written and ErrConflict is returned.
// On success, the returned slice lists the paths written.
//
// Uses packages/filesystem so directory creation and atomic-style writes
// stay consistent with the rest of bedrock.
func GenerateRules(root, body string, force bool) ([]string, error) {
	fs := filesystem.New()

	if !force {
		for _, t := range Targets {
			if fs.Exists(filepath.Join(root, t.Path)) {
				return nil, ErrConflict
			}
		}
	}

	written := make([]string, 0, len(Targets))

	for _, t := range Targets {
		full := filepath.Join(root, t.Path)

		if err := fs.Put(full, []byte(body)); err != nil {
			return written, err
		}

		written = append(written, t.Path)
	}

	return written, nil
}
