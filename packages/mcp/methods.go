package mcp

import (
	"context"
	"fmt"
)

// supportedVersions lists the MCP protocol versions this server supports,
// in preference order (most recent first).
var supportedVersions = []string{
	"2025-11-25",
	"2025-06-18",
	"2025-03-26",
	"2024-11-05",
}

// versionsWithoutInstructions are protocol versions that must not include
// the instructions field in the initialize response.
var versionsWithoutInstructions = map[string]bool{
	"2024-11-05": true,
	"2025-03-26": true,
}

// methodHandler is the function signature for all JSON-RPC method handlers.
type methodHandler func(ctx context.Context, req *JsonRpcRequest, sc *ServerContext) (map[string]any, error)

// dispatchTable maps JSON-RPC method names to handler functions.
var dispatchTable = map[string]methodHandler{
	"initialize":               handleInitialize,
	"ping":                     handlePing,
	"tools/list":               handleToolsList,
	"tools/call":               handleToolsCall,
	"resources/list":           handleResourcesList,
	"resources/templates/list": handleResourcesTemplatesList,
	"resources/read":           handleResourcesRead,
	"prompts/list":             handlePromptsList,
	"prompts/get":              handlePromptsGet,
	"completion/complete":      handleCompletionComplete,
}

// handleInitialize handles the initialize JSON-RPC method.
// It negotiates the protocol version and returns server capabilities.
func handleInitialize(_ context.Context, req *JsonRpcRequest, sc *ServerContext) (map[string]any, error) {
	// Negotiate protocol version.
	requested, _ := req.Get("protocolVersion").(string)
	version := negotiateVersion(requested)

	capabilities := buildCapabilities(sc)

	result := map[string]any{
		"protocolVersion": version,
		"capabilities":    capabilities,
		"serverInfo": map[string]any{
			"name":    sc.ServerName,
			"version": sc.ServerVersion,
		},
	}

	// Older protocol versions do not include the instructions field.
	if sc.Instructions != "" && !versionsWithoutInstructions[version] {
		result["instructions"] = sc.Instructions
	}

	return result, nil
}

// negotiateVersion selects the best supported protocol version. If the
// requested version is supported it is used; otherwise the most recent
// supported version is returned.
func negotiateVersion(requested string) string {
	for _, v := range supportedVersions {
		if v == requested {
			return v
		}
	}
	return supportedVersions[0]
}

// buildCapabilities constructs the capabilities map for the initialize response.
func buildCapabilities(sc *ServerContext) map[string]any {
	caps := map[string]any{}
	if len(sc.tools) > 0 {
		caps["tools"] = map[string]any{"listChanged": false}
	}
	if len(sc.resources) > 0 {
		caps["resources"] = map[string]any{"listChanged": false, "subscribe": false}
	}
	if len(sc.prompts) > 0 {
		caps["prompts"] = map[string]any{"listChanged": false}
	}
	if sc.HasCompletions() {
		caps["completions"] = map[string]any{}
	}
	return caps
}

// handlePing handles the ping method.
func handlePing(_ context.Context, _ *JsonRpcRequest, _ *ServerContext) (map[string]any, error) {
	return map[string]any{}, nil
}

// handleToolsList handles tools/list with cursor pagination.
func handleToolsList(_ context.Context, req *JsonRpcRequest, sc *ServerContext) (map[string]any, error) {
	tools := sc.ToolsList()
	items := make([]any, len(tools))
	for i, t := range tools {
		items[i] = toolToMap(t)
	}

	perPage := sc.PerPage(0)
	pager := NewCursorPaginator(items, perPage, req.Cursor())
	return pager.Paginate("tools"), nil
}

// handleToolsCall handles tools/call.
func handleToolsCall(ctx context.Context, req *JsonRpcRequest, sc *ServerContext) (map[string]any, error) {
	name, _ := req.Get("name").(string)
	tool, ok := sc.FindTool(name)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrToolNotFound, name)
	}

	mcpReq := req.ToRequest()
	resp, err := tool.Handle(ctx, mcpReq)
	if err != nil {
		// Surface handler errors as MCP error results (isError: true).
		return Error(err.Error()).toToolResult(), nil
	}

	return resp.toToolResult(), nil
}

// handleResourcesList handles resources/list with cursor pagination.
func handleResourcesList(_ context.Context, req *JsonRpcRequest, sc *ServerContext) (map[string]any, error) {
	resources := sc.ResourcesList()
	items := make([]any, len(resources))
	for i, r := range resources {
		items[i] = resourceToMap(r)
	}

	perPage := sc.PerPage(0)
	pager := NewCursorPaginator(items, perPage, req.Cursor())
	return pager.Paginate("resources"), nil
}

// handleResourcesTemplatesList handles resources/templates/list.
func handleResourcesTemplatesList(_ context.Context, _ *JsonRpcRequest, sc *ServerContext) (map[string]any, error) {
	templates := sc.ResourceTemplates()
	items := make([]map[string]any, len(templates))
	for i, rt := range templates {
		items[i] = resourceTemplateToMap(rt)
	}
	return map[string]any{"resourceTemplates": items}, nil
}

// handleResourcesRead handles resources/read.
func handleResourcesRead(ctx context.Context, req *JsonRpcRequest, sc *ServerContext) (map[string]any, error) {
	uri, _ := req.Get("uri").(string)
	resource, vars, ok := sc.FindResource(uri)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrResourceNotFound, uri)
	}

	mcpReq := req.ToRequest()
	mcpReq.URIVars = vars

	resp, err := resource.Read(ctx, mcpReq)
	if err != nil {
		return nil, err
	}

	contents := resp.toResourceResult(uri)
	return map[string]any{"contents": contents}, nil
}

// handlePromptsList handles prompts/list with cursor pagination.
func handlePromptsList(_ context.Context, req *JsonRpcRequest, sc *ServerContext) (map[string]any, error) {
	prompts := sc.PromptsList()
	items := make([]any, len(prompts))
	for i, p := range prompts {
		items[i] = promptToMap(p)
	}

	perPage := sc.PerPage(0)
	pager := NewCursorPaginator(items, perPage, req.Cursor())
	return pager.Paginate("prompts"), nil
}

// handlePromptsGet handles prompts/get.
func handlePromptsGet(ctx context.Context, req *JsonRpcRequest, sc *ServerContext) (map[string]any, error) {
	name, _ := req.Get("name").(string)
	prompt, ok := sc.FindPrompt(name)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrPromptNotFound, name)
	}

	mcpReq := req.ToRequest()
	// The arguments sub-object is already extracted by ToRequest; also
	// accept a top-level arguments map from prompts/get params.
	if a, ok := req.Params["arguments"].(map[string]any); ok {
		mcpReq.Arguments = a
	}

	messages, err := prompt.Invoke(ctx, mcpReq)
	if err != nil {
		return nil, err
	}

	msgs := make([]map[string]any, 0, len(messages))
	for _, m := range messages {
		msgs = append(msgs, m.ToMap())
	}

	return map[string]any{
		"description": prompt.Description(),
		"messages":    msgs,
	}, nil
}

// handleCompletionComplete handles completion/complete.
func handleCompletionComplete(ctx context.Context, req *JsonRpcRequest, sc *ServerContext) (map[string]any, error) {
	ref, _ := req.Params["ref"].(map[string]any)
	if ref == nil {
		return map[string]any{"completion": EmptyCompletion().toMap()}, nil
	}

	refType, _ := ref["type"].(string)
	refName, _ := ref["name"].(string)

	argMap, _ := req.Params["argument"].(map[string]any)
	argName, _ := argMap["name"].(string)
	argValue, _ := argMap["value"].(string)

	var completable Completable
	switch refType {
	case "ref/prompt":
		p, ok := sc.FindPrompt(refName)
		if !ok {
			return map[string]any{"completion": EmptyCompletion().toMap()}, nil
		}
		completable, _ = p.(Completable)
	case "ref/resource":
		r, _, ok := sc.FindResource(refName)
		if !ok {
			return map[string]any{"completion": EmptyCompletion().toMap()}, nil
		}
		completable, _ = r.(Completable)
	}

	if completable == nil {
		return map[string]any{"completion": EmptyCompletion().toMap()}, nil
	}

	result := completable.Complete(ctx, argName, argValue)
	return map[string]any{"completion": result.toMap()}, nil
}
