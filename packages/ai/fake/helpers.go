package fake

import "github.com/bedrock/packages/ai/data"

// DataUsage returns a zero-value Usage for use in test fixtures.
func DataUsage() data.Usage { return data.Usage{} }

// DataMeta returns a zero-value Meta for use in test fixtures.
func DataMeta() data.Meta { return data.Meta{Citations: []any{}} }
