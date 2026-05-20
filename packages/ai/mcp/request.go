package mcp

// Request carries the arguments, session identifier, and contextual metadata
// for a single MCP handler invocation. It mirrors the upstream Request class.
type Request struct {
	// Arguments holds the caller-supplied parameters (tools/call arguments,
	// prompts/get arguments, etc.).
	Arguments map[string]any

	// SessionID is the MCP session identifier, if any.
	SessionID string

	// Meta holds the _meta object from the JSON-RPC params, if present.
	Meta map[string]any

	// URI is the resource URI supplied for resources/read requests.
	URI string

	// URIVars holds the variables extracted when a URI template matched the
	// incoming URI (e.g. {"id": "42"} for template "file://users/{id}").
	URIVars map[string]string
}

// Get returns the argument value for key. If the key is absent the first
// element of fallback is returned, or nil when no fallback is provided.
func (r *Request) Get(key string, fallback ...any) any {
	if r.Arguments != nil {
		if v, ok := r.Arguments[key]; ok {
			return v
		}
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return nil
}

// All returns a shallow copy of the arguments map.
func (r *Request) All() map[string]any {
	out := make(map[string]any, len(r.Arguments))

	for k, v := range r.Arguments {
		out[k] = v
	}

	return out
}

// Merge returns a new Request with extra entries merged into Arguments.
// Existing keys are overwritten by data.
func (r *Request) Merge(data map[string]any) *Request {
	merged := make(map[string]any, len(r.Arguments)+len(data))

	for k, v := range r.Arguments {
		merged[k] = v
	}

	for k, v := range data {
		merged[k] = v
	}

	return &Request{
		Arguments: merged,
		SessionID: r.SessionID,
		Meta:      r.Meta,
		URI:       r.URI,
		URIVars:   r.URIVars,
	}
}
