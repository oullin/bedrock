package mcp_test

import (
	"encoding/json"
	"testing"
)

func decodeJSON(t testing.TB, raw string) map[string]any {
	t.Helper()

	var out map[string]any

	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("invalid JSON %q: %v", raw, err)
	}

	return out
}
