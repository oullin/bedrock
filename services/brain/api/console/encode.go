package console

import (
	"encoding/json"
	"fmt"

	"github.com/bedrock/packages/filesystem"
)

// writeJSON is the shared encoder for graph JSON outputs. It marshals with
// 2-space indent and writes via packages/filesystem so missing parents are
// created automatically.
func writeJSON(path string, v any) error {
	body, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	if err := filesystem.New().Put(path, append(body, '\n')); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
