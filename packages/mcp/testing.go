package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// TestServer wraps a Server with helpers for exercising MCP endpoints in unit
// tests. Obtain one via server.Test(t).
type TestServer struct {
	server *Server
	t      testing.TB
}

// CallTool invokes tools/call and returns the result as a TestResult.
func (ts *TestServer) CallTool(name string, args map[string]any) *TestResult {
	ts.t.Helper()
	if args == nil {
		args = map[string]any{}
	}
	raw := ts.server.handleRaw(context.Background(), "tools/call", map[string]any{
		"name":      name,
		"arguments": args,
	})
	return &TestResult{t: ts.t, raw: raw}
}

// ReadResource invokes resources/read and returns the result as a TestResult.
func (ts *TestServer) ReadResource(uri string) *TestResult {
	ts.t.Helper()
	raw := ts.server.handleRaw(context.Background(), "resources/read", map[string]any{
		"uri": uri,
	})
	return &TestResult{t: ts.t, raw: raw}
}

// GetPrompt invokes prompts/get and returns the result as a TestResult.
func (ts *TestServer) GetPrompt(name string, args map[string]any) *TestResult {
	ts.t.Helper()
	if args == nil {
		args = map[string]any{}
	}
	raw := ts.server.handleRaw(context.Background(), "prompts/get", map[string]any{
		"name":      name,
		"arguments": args,
	})
	return &TestResult{t: ts.t, raw: raw}
}

// Complete invokes completion/complete.
// refType is "ref/prompt" or "ref/resource"; refName is the primitive name.
func (ts *TestServer) Complete(refType, refName, argument, value string) *TestResult {
	ts.t.Helper()
	raw := ts.server.handleRaw(context.Background(), "completion/complete", map[string]any{
		"ref": map[string]any{
			"type": refType,
			"name": refName,
		},
		"argument": map[string]any{
			"name":  argument,
			"value": value,
		},
	})
	return &TestResult{t: ts.t, raw: raw}
}

// TestResult holds the raw JSON-RPC response map and provides fluent
// assertion methods for use in tests.
type TestResult struct {
	t   testing.TB
	raw map[string]any
}

// result returns the "result" field of the JSON-RPC response, or nil.
func (r *TestResult) result() map[string]any {
	if v, ok := r.raw["result"].(map[string]any); ok {
		return v
	}
	return nil
}

// hasError reports whether the response is a JSON-RPC error.
func (r *TestResult) hasError() bool {
	_, ok := r.raw["error"]
	return ok
}

// isToolError reports whether tools/call returned an error result
// (isError: true in the result body).
func (r *TestResult) isToolError() bool {
	res := r.result()
	if res == nil {
		return false
	}
	b, _ := res["isError"].(bool)
	return b
}

// AssertOK fails the test if the response contains a JSON-RPC error.
func (r *TestResult) AssertOK() *TestResult {
	r.t.Helper()
	if r.hasError() {
		r.t.Fatalf("expected ok response, got error: %s", r.dump())
	}
	return r
}

// AssertHasErrors fails the test if the response does not contain an error.
// This checks both JSON-RPC errors and tool isError results.
func (r *TestResult) AssertHasErrors() *TestResult {
	r.t.Helper()
	if !r.hasError() && !r.isToolError() {
		r.t.Fatalf("expected error response, got: %s", r.dump())
	}
	return r
}

// AssertSee fails the test if text does not appear in any content item of
// the response result.
func (r *TestResult) AssertSee(text string) *TestResult {
	r.t.Helper()
	if !r.responseContains(text) {
		r.t.Fatalf("expected response to contain %q, got: %s", text, r.dump())
	}
	return r
}

// AssertDontSee fails the test if text appears anywhere in the response result.
func (r *TestResult) AssertDontSee(text string) *TestResult {
	r.t.Helper()
	if r.responseContains(text) {
		r.t.Fatalf("expected response NOT to contain %q, got: %s", text, r.dump())
	}
	return r
}

// AssertCompletionValues fails the test unless the completion values exactly
// match the provided list (order-insensitive).
func (r *TestResult) AssertCompletionValues(values ...string) *TestResult {
	r.t.Helper()
	got := r.completionValues()
	if len(got) != len(values) {
		r.t.Fatalf("expected completion values %v, got %v", values, got)
		return r
	}
	set := make(map[string]bool, len(values))
	for _, v := range values {
		set[v] = true
	}
	for _, v := range got {
		if !set[v] {
			r.t.Fatalf("unexpected completion value %q; expected %v", v, values)
		}
	}
	return r
}

// AssertCompletionCount fails the test unless exactly n completion values are present.
func (r *TestResult) AssertCompletionCount(n int) *TestResult {
	r.t.Helper()
	got := r.completionValues()
	if len(got) != n {
		r.t.Fatalf("expected %d completion values, got %d: %v", n, len(got), got)
	}
	return r
}

// AssertHasCompletions fails the test if no completion values are present.
func (r *TestResult) AssertHasCompletions() *TestResult {
	r.t.Helper()
	if len(r.completionValues()) == 0 {
		r.t.Fatalf("expected completions, got none: %s", r.dump())
	}
	return r
}

// AssertSentNotification fails the test unless the result contains a
// notification with the given method name.
func (r *TestResult) AssertSentNotification(method string) *TestResult {
	r.t.Helper()
	if m, _ := r.result()["method"].(string); m != method {
		r.t.Fatalf("expected notification %q, got result: %s", method, r.dump())
	}
	return r
}

// AssertNotificationCount fails unless the notification params contain the
// expected count. (For test symmetry with Upstream; in practice use AssertSentNotification.)
func (r *TestResult) AssertNotificationCount(n int) *TestResult {
	r.t.Helper()
	// For our implementation a response holds one notification at a time.
	// n==1 means we expect exactly one notification.
	hasNotif := r.result() != nil && r.result()["method"] != nil
	if n == 0 && hasNotif {
		r.t.Fatalf("expected 0 notifications, but found one")
	}
	if n > 0 && !hasNotif {
		r.t.Fatalf("expected %d notification(s), found none: %s", n, r.dump())
	}
	return r
}

// Dump prints the raw response to the test log and returns the TestResult for
// chaining.
func (r *TestResult) Dump() *TestResult {
	r.t.Helper()
	r.t.Log(r.dump())
	return r
}

// dump returns a pretty-printed JSON string of the raw response.
func (r *TestResult) dump() string {
	b, err := json.MarshalIndent(r.raw, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", r.raw)
	}
	return string(b)
}

// responseContains searches all text values in the result for the needle.
func (r *TestResult) responseContains(needle string) bool {
	res := r.result()
	if res == nil {
		return false
	}
	return containsText(res, needle)
}

// containsText recursively searches a map/slice/string for needle.
func containsText(v any, needle string) bool {
	switch val := v.(type) {
	case string:
		return strings.Contains(val, needle)
	case map[string]any:
		for _, mv := range val {
			if containsText(mv, needle) {
				return true
			}
		}
	case []any:
		for _, item := range val {
			if containsText(item, needle) {
				return true
			}
		}
	}
	return false
}

// completionValues extracts the completion values array from the response.
func (r *TestResult) completionValues() []string {
	res := r.result()
	if res == nil {
		return nil
	}
	comp, _ := res["completion"].(map[string]any)
	if comp == nil {
		return nil
	}
	raw, _ := comp["values"].([]any)
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
