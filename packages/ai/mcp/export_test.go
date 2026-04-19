package mcp

// ExportToolResult exposes the internal toToolResult method for use in
// package-level tests.
func ExportToolResult(r *Response) map[string]any {
	return r.toToolResult()
}
