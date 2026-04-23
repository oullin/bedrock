package mcp

// ExportToolResult exposes the internal toToolResult method for use in
// package-level tests.
func ExportToolResult(r *Response) map[string]any {
	return r.toToolResult()
}

// ExportServerContext exposes the server's request snapshot for tests that
// need to exercise internal pagination and capability helpers.
func ExportServerContext(s *Server) *ServerContext {
	return s.context()
}
