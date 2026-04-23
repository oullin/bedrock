package fake

import "github.com/bedrock/packages/ai/sdk/data"

// DefaultTextTimeout is the SDK fallback used by fake text prompts when no
// explicit timeout is supplied.
const DefaultTextTimeout = 60

// DefaultMediaTimeout is the SDK fallback for non-text fake gateway prompts.
const DefaultMediaTimeout = 30

// DataUsage returns a zero-value Usage for use in test fixtures.
func DataUsage() data.Usage { return data.Usage{} }

// DataMeta returns a zero-value Meta for use in test fixtures.
func DataMeta() data.Meta { return data.Meta{Citations: []any{}} }

func timeoutOrDefault(timeout, fallback int) int {
	if timeout > 0 {
		return timeout
	}

	return fallback
}
