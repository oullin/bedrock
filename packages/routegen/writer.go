package routegen

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed resources/routegen.ts
var routegenTSContent []byte

// writeRouteGenTS copies the embedded routegen.ts runtime utility to
// {dir}/index.ts, creating the directory if necessary.
func writeRouteGenTS(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "index.ts"), routegenTSContent, 0644)
}

// writeFile writes content to path, creating intermediate directories.
func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}
