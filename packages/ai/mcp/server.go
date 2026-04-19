package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// Server is an MCP server that handles JSON-RPC 2.0 messages and exposes
// tools, resources, and prompts to AI agents.
type Server struct {
	name         string
	version      string
	description  string
	instructions string

	defaultPaginationLength int
	maxPaginationLength     int

	tools     []Tool
	resources []Resource
	prompts   []Prompt
}

// Option is a functional option for configuring a Server.
type Option func(*Server)

const (
	defaultPagination = 15
	maxPagination     = 50
)

// WithDescription sets a human-readable description for the server.
func WithDescription(d string) Option {
	return func(s *Server) { s.description = d }
}

// WithInstructions sets the system instructions included in the initialize
// response (MCP protocol versions ≥ 2025-06-18).
func WithInstructions(i string) Option {
	return func(s *Server) { s.instructions = i }
}

// WithPagination configures the default and maximum pagination page sizes.
func WithPagination(defaultN, maxN int) Option {
	return func(s *Server) {
		s.defaultPaginationLength = defaultN
		s.maxPaginationLength = maxN
	}
}

// NewServer creates a new MCP server with the given name and version.
func NewServer(name, version string, opts ...Option) *Server {
	s := &Server{
		name:                    name,
		version:                 version,
		defaultPaginationLength: defaultPagination,
		maxPaginationLength:     maxPagination,
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

// AddTool registers one or more tools with the server.
func (s *Server) AddTool(tools ...Tool) *Server {
	s.tools = append(s.tools, tools...)

	return s
}

// AddResource registers one or more resources (or resource templates) with
// the server.
func (s *Server) AddResource(resources ...Resource) *Server {
	s.resources = append(s.resources, resources...)

	return s
}

// AddPrompt registers one or more prompts with the server.
func (s *Server) AddPrompt(prompts ...Prompt) *Server {
	s.prompts = append(s.prompts, prompts...)

	return s
}

// context builds a ServerContext snapshot from the current server state.
func (s *Server) context() *ServerContext {
	return &ServerContext{
		SupportedVersions:       supportedVersions,
		ServerName:              s.name,
		ServerVersion:           s.version,
		Description:             s.description,
		Instructions:            s.instructions,
		MaxPaginationLength:     s.maxPaginationLength,
		DefaultPaginationLength: s.defaultPaginationLength,
		tools:                   s.tools,
		resources:               s.resources,
		prompts:                 s.prompts,
	}
}

// Handle processes a raw JSON-RPC message string and returns the serialised
// JSON-RPC response. sessionID is the MCP session identifier, if any.
func (s *Server) Handle(ctx context.Context, message, sessionID string) (string, error) {
	req, err := ParseJsonRpcRequest([]byte(message), sessionID)

	if err != nil {
		var resp *JsonRpcResponse

		if err == ErrInvalidRequest {
			resp = ErrorResponse(nil, CodeInvalidRequest, err.Error())
		} else {
			resp = ErrorResponse(nil, CodeParseError, err.Error())
		}

		b, _ := resp.ToJSON()

		return string(b), nil
	}

	sc := s.context()
	handler, ok := dispatchTable[req.Method]

	if !ok {
		resp := ErrorResponse(req.ID, CodeMethodNotFound, fmt.Sprintf("method not found: %q", req.Method))
		b, _ := resp.ToJSON()

		return string(b), nil
	}

	result, err := func() (result map[string]any, rerr error) {
		defer func() {
			if r := recover(); r != nil {
				rerr = fmt.Errorf("internal panic: %v", r)
			}
		}()

		return handler(ctx, req, sc)
	}()

	if err != nil {
		resp := ErrorResponse(req.ID, CodeInternalError, err.Error())
		b, _ := resp.ToJSON()

		return string(b), nil
	}

	resp := ResultResponse(req.ID, result)
	b, err := resp.ToJSON()

	if err != nil {
		errResp := ErrorResponse(req.ID, CodeInternalError, "failed to serialise response")
		b, _ = errResp.ToJSON()
	}

	return string(b), nil
}

// ServeHTTP implements http.Handler. Register it on any ServeMux:
//
//	mux.Handle("/mcp", server)
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := NewHttpTransport(w, r)
	t.OnReceive(s.Handle)
	t.Run(r.Context()) //nolint:errcheck
}

// ServeStdio starts the stdio transport and blocks until ctx is cancelled or
// EOF is reached on stdin.
func (s *Server) ServeStdio(ctx context.Context) error {
	t := NewStdioTransport()
	t.OnReceive(s.Handle)

	return t.Run(ctx)
}

// Test returns a TestServer for exercising the server in unit tests.
func (s *Server) Test(t testing.TB) *TestServer {
	return &TestServer{server: s, t: t}
}

// handleRaw is a lower-level helper used by TestServer: it calls Handle and
// unmarshals the JSON result into a map.
func (s *Server) handleRaw(ctx context.Context, method string, params map[string]any) map[string]any {
	if params == nil {
		params = map[string]any{}
	}

	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}
	b, _ := json.Marshal(req)
	respStr, _ := s.Handle(ctx, string(b), "test-session")

	var out map[string]any

	json.Unmarshal([]byte(respStr), &out) //nolint:errcheck

	return out
}
