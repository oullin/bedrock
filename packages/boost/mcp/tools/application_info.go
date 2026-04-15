package tools

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

// ApplicationInfo returns comprehensive application information: Go version,
// module name, OS, architecture, and Go dependencies from go.mod.
// Mirrors Upstream\Boost\Mcp\Tools\ApplicationInfo.
// Tagged IsReadOnly.
type ApplicationInfo struct {
	// ModFilePath is the path to the go.mod file to inspect.
	// Defaults to "go.mod" when empty.
	ModFilePath string
}

func (t *ApplicationInfo) Name() string { return "application_info" }
func (t *ApplicationInfo) IsReadOnly() bool { return true }

func (t *ApplicationInfo) Description() string {
	return "Get comprehensive application information including Go version, module name, " +
		"OS, architecture, and all module dependencies with their versions. " +
		"Use this tool at the start of each session to understand the project context."
}

func (t *ApplicationInfo) Schema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	}
}

// Handle returns application info.
func (t *ApplicationInfo) Handle(_ McpRequest) (McpResponse, error) {
	modPath := t.ModFilePath
	if modPath == "" {
		modPath = "go.mod"
	}

	moduleName, deps := parseGoMod(modPath)

	return OkResponse(map[string]any{
		"go_version":  runtime.Version(),
		"module_name": moduleName,
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"packages":    deps,
	}), nil
}

// parseGoMod extracts the module name and require dependencies from go.mod.
func parseGoMod(path string) (string, []map[string]string) {
	f, err := os.Open(path)
	if err != nil {
		return "", nil
	}
	defer f.Close()

	var moduleName string
	var deps []map[string]string
	inRequire := false
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "module ") {
			moduleName = strings.TrimPrefix(line, "module ")
			continue
		}

		if line == "require (" {
			inRequire = true
			continue
		}

		if inRequire && line == ")" {
			inRequire = false
			continue
		}

		if inRequire || strings.HasPrefix(line, "require ") {
			entry := strings.TrimPrefix(line, "require ")
			entry = strings.TrimSuffix(entry, " // indirect")
			parts := strings.Fields(entry)

			if len(parts) == 2 {
				deps = append(deps, map[string]string{
					"package": parts[0],
					"version": parts[1],
				})
			}
		}
	}

	return moduleName, deps
}
