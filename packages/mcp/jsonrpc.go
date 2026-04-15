package mcp

import (
	"encoding/json"
	"fmt"
)

// JSON-RPC 2.0 version string.
const jsonrpcVersion = "2.0"

// Standard JSON-RPC 2.0 error codes.
const (
	CodeParseError     = -32700
	CodeInvalidRequest = -32600
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternalError  = -32603
)

// JsonRpcRequest represents a decoded JSON-RPC 2.0 request.
type JsonRpcRequest struct {
	JSONRPC   string         `json:"jsonrpc"`
	ID        any            `json:"id"`
	Method    string         `json:"method"`
	Params    map[string]any `json:"params"`
	SessionID string         `json:"-"`
}

// ParseJsonRpcRequest decodes a raw JSON message into a JsonRpcRequest.
// sessionID is the MCP session identifier carried out-of-band (e.g. via an
// HTTP header).
func ParseJsonRpcRequest(data []byte, sessionID string) (*JsonRpcRequest, error) {
	var r JsonRpcRequest
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseError, err)
	}
	if r.Method == "" {
		return nil, ErrInvalidRequest
	}
	r.SessionID = sessionID
	if r.Params == nil {
		r.Params = map[string]any{}
	}
	return &r, nil
}

// Get returns the value of a params key. Returns fallback (or nil) when
// the key is absent.
func (r *JsonRpcRequest) Get(key string, fallback ...any) any {
	if v, ok := r.Params[key]; ok {
		return v
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return nil
}

// Cursor returns the pagination cursor from params, if present.
func (r *JsonRpcRequest) Cursor() string {
	if c, ok := r.Params["cursor"].(string); ok {
		return c
	}
	return ""
}

// Meta returns the _meta map from params, if present.
func (r *JsonRpcRequest) Meta() map[string]any {
	if m, ok := r.Params["_meta"].(map[string]any); ok {
		return m
	}
	return nil
}

// ToRequest converts the JSON-RPC request into an MCP Request suitable for
// passing to tool/resource/prompt handlers.
func (r *JsonRpcRequest) ToRequest() *Request {
	args := make(map[string]any, len(r.Params))
	for k, v := range r.Params {
		if k != "_meta" {
			args[k] = v
		}
	}
	// arguments sub-object takes precedence when present (tools/call)
	if a, ok := r.Params["arguments"].(map[string]any); ok {
		args = a
	}
	uri, _ := r.Params["uri"].(string)
	return &Request{
		Arguments: args,
		SessionID: r.SessionID,
		Meta:      r.Meta(),
		URI:       uri,
	}
}

// JsonRpcResponse represents a JSON-RPC 2.0 response message.
type JsonRpcResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      any           `json:"id,omitempty"`
	Result  any           `json:"result,omitempty"`
	Error   *JsonRpcError `json:"error,omitempty"`
}

// JsonRpcError is the error object inside a JSON-RPC error response.
type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ResultResponse builds a successful JSON-RPC response.
func ResultResponse(id, result any) *JsonRpcResponse {
	return &JsonRpcResponse{JSONRPC: jsonrpcVersion, ID: id, Result: result}
}

// ErrorResponse builds an error JSON-RPC response.
func ErrorResponse(id any, code int, message string, data ...any) *JsonRpcResponse {
	e := &JsonRpcError{Code: code, Message: message}
	if len(data) > 0 {
		e.Data = data[0]
	}
	return &JsonRpcResponse{JSONRPC: jsonrpcVersion, ID: id, Error: e}
}

// NotificationResponse builds a JSON-RPC notification (no ID).
func NotificationResponse(method string, params map[string]any) *JsonRpcResponse {
	if params == nil {
		params = map[string]any{}
	}
	return &JsonRpcResponse{
		JSONRPC: jsonrpcVersion,
		Result:  map[string]any{"method": method, "params": params},
	}
}

// ToJSON serialises the response to a JSON byte slice.
func (r *JsonRpcResponse) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
