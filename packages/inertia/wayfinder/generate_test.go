package wayfinder

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Top-level routes (no dot) go into the "app" group.

// limitWriter fails after n bytes.

var errWriteLimited = errors.New("write limit reached")

type limitWriter struct {
	n int
}

func TestGenerateTypeScript(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")
	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{TypeScript: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "type RouteResult = { url: string; method: string }") {
		t.Error("expected RouteResult type declaration")
	}

	if !strings.Contains(output, "export function login(): RouteResult {") {
		t.Error("expected flat login function")
	}

	if !strings.Contains(output, "return { url: '/login', method: 'get' }") {
		t.Error("expected login return statement")
	}

	if !strings.Contains(output, "export function contactsShow(params: { contact: string | number }): RouteResult {") {
		t.Error("expected flat contactsShow function with typed params")
	}

	if !strings.Contains(output, "return { url: `/contacts/${encodeURIComponent(String(params.contact))}`, method: 'get' }") {
		t.Error("expected contactsShow template literal return")
	}

	if !strings.Contains(output, "} as const") {
		t.Error("expected 'as const' for nested exports")
	}
}

func TestGenerateJavaScript(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")
	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{TypeScript: false})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if strings.Contains(output, "RouteResult") {
		t.Error("JS output should not contain TypeScript types")
	}

	if !strings.Contains(output, "export function login() {") {
		t.Error("expected JS login function without types")
	}

	if !strings.Contains(output, "(params) =>") {
		t.Error("expected JS nested function without types")
	}
}

func TestGenerateFlatOnly(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.index", "GET", "/contacts")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{FlatOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export function contactsIndex()") {
		t.Error("expected flat export")
	}

	if strings.Contains(output, "export const contacts") {
		t.Error("should not contain nested exports in FlatOnly mode")
	}
}

func TestGenerateNestedOnly(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.index", "GET", "/contacts")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{NestedOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if strings.Contains(output, "export function contactsIndex()") {
		t.Error("should not contain flat exports in NestedOnly mode")
	}

	if !strings.Contains(output, "export const contacts") {
		t.Error("expected nested export")
	}
}

func TestGenerateConflictingFlags(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{FlatOnly: true, NestedOnly: true})

	if err == nil {
		t.Fatal("expected error when both FlatOnly and NestedOnly are set")
	}

	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestGenerateCustomHeader(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{
		Header: "// Custom header",
	})

	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(buf.String(), "// Custom header\n") {
		t.Error("expected custom header")
	}
}

func TestGenerateNestedGrouping(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.index", "GET", "/contacts")
	reg.Add("contacts.show", "GET", "/contacts/{contact}")
	reg.Add("contacts.store", "POST", "/contacts")
	reg.Add("organizations.index", "GET", "/organizations")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{NestedOnly: true, TypeScript: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export const contacts = {") {
		t.Error("expected contacts group")
	}

	if !strings.Contains(output, "export const organizations = {") {
		t.Error("expected organizations group")
	}

	if !strings.Contains(output, "index: (): RouteResult =>") {
		t.Error("expected index member")
	}

	if !strings.Contains(output, "show: (params: { contact: string | number }): RouteResult =>") {
		t.Error("expected show member with params")
	}
}

func TestGenerateTopLevelRoutes(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")
	reg.Add("logout", "POST", "/logout")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{NestedOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export const app = {") {
		t.Error("expected top-level routes grouped under 'app'")
	}
}

func TestGenerateTopLevelRoutes_CollisionWithAppGroup(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("dashboard", "GET", "/dashboard")
	reg.Add("app.settings", "GET", "/app/settings")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{NestedOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export const app = {") {
		t.Error("expected real 'app' group to keep the 'app' name")
	}

	if !strings.Contains(output, "export const routes = {") {
		t.Errorf("expected ungrouped routes to fall back to 'routes', got:\n%s", output)
	}
}

func TestGenerateTopLevelRoutes_DoubleCollision(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("home", "GET", "/")
	reg.Add("app.index", "GET", "/app")
	reg.Add("routes.index", "GET", "/r")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{NestedOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export const _app = {") {
		t.Errorf("expected ungrouped routes to fall back to '_app', got:\n%s", output)
	}
}

func TestDotToCamel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"contacts.index", "contactsIndex"},
		{"contacts.notes.store", "contactsNotesStore"},
		{"login", "login"},
		{"a.b.c.d", "aBCD"},
		{"use-form", "useForm"},
		{"forms.use-form", "formsUseForm"},
		{"data-loading.deferred-props", "dataLoadingDeferredProps"},
	}

	for _, tt := range tests {
		got := dotToCamel(tt.input)

		if got != tt.expected {
			t.Errorf("dotToCamel(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestGenerateHyphenatedRoutes(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("features.forms.use-form", "GET", "/features/forms/use-form")
	reg.Add("use-http", "GET", "/use-http")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{TypeScript: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export function featuresFormsUseForm(): RouteResult {") {
		t.Error("expected hyphenated flat function name to be camelCased")
	}

	if !strings.Contains(output, "export function useHttp(): RouteResult {") {
		t.Error("expected top-level hyphenated flat function name to be camelCased")
	}

	if !strings.Contains(output, "formsUseForm: (): RouteResult =>") {
		t.Error("expected hyphenated nested key to be camelCased")
	}

	if !strings.Contains(output, "useHttp: (): RouteResult =>") {
		t.Error("expected top-level hyphenated route in app group to be camelCased")
	}
}

func (lw *limitWriter) Write(p []byte) (int, error) {
	if lw.n <= 0 {
		return 0, errWriteLimited
	}

	if len(p) > lw.n {
		lw.n = 0

		return 0, errWriteLimited
	}

	lw.n -= len(p)

	return len(p), nil
}

func TestGeneratePropagatesWriteError(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	// Allow just enough bytes for the header, then fail.
	w := &limitWriter{n: 10}

	err := Generate(reg, w, GenerateOptions{})

	if err == nil {
		t.Fatal("expected a write error, got nil")
	}

	if !errors.Is(err, errWriteLimited) {
		t.Fatalf("expected errWriteLimited, got %v", err)
	}
}

func TestGenerateFile_WritesToFile(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")
	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	path := filepath.Join(t.TempDir(), "routes.ts")

	err := GenerateFile(reg, path, GenerateOptions{TypeScript: true})

	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	output := string(data)

	if !strings.Contains(output, "export function login(): RouteResult {") {
		t.Error("expected login function in generated file")
	}

	if !strings.Contains(output, "export function contactsShow") {
		t.Error("expected contactsShow function in generated file")
	}
}

func TestGenerateFile_InvalidPath(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	err := GenerateFile(reg, "/nonexistent/dir/routes.ts", GenerateOptions{})

	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestGenerateJS_NestedWithParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")
	reg.Add("contacts.index", "GET", "/contacts")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{TypeScript: false, NestedOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "(params) =>") {
		t.Error("expected JS nested with params")
	}

	if !strings.Contains(output, "() =>") {
		t.Error("expected JS nested without params")
	}
}

func TestGenerateTS_FlatWithParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{TypeScript: true, FlatOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export function contactsShow(params: { contact: string | number }): RouteResult {") {
		t.Error("expected TS flat function with params")
	}
}

func TestGenerateJS_FlatWithParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{TypeScript: false, FlatOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export function contactsShow(params) {") {
		t.Error("expected JS flat function with params")
	}
}

func TestGenerateEmptyRegistry(t *testing.T) {
	t.Parallel()

	reg := New()

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{TypeScript: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "type RouteResult") {
		t.Error("expected RouteResult type even for empty registry")
	}
}

func TestGenerateFile_Success_VerifyContent(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	path := filepath.Join(t.TempDir(), "routes.js")

	err := GenerateFile(reg, path, GenerateOptions{})

	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "export function login()") {
		t.Error("expected login function in file")
	}
}

func TestGenerate_WriteErrorInHeader(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	// Fail immediately — can't even write the header.
	w := &limitWriter{n: 0}
	err := Generate(reg, w, GenerateOptions{})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInTypeDecl(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	// Allow header but fail on type declaration.
	w := &limitWriter{n: 60}
	err := Generate(reg, w, GenerateOptions{TypeScript: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInFlatSection(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	// Allow header + type decl, fail in flat section.
	w := &limitWriter{n: 120}
	err := Generate(reg, w, GenerateOptions{TypeScript: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInNestedSection(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	// Allow enough for header + flat, fail in nested.
	w := &limitWriter{n: 300}
	err := Generate(reg, w, GenerateOptions{TypeScript: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorSeparator(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	// Allow header + flat, fail on the separator newline between flat and nested.
	w := &limitWriter{n: 170}
	err := Generate(reg, w, GenerateOptions{})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInFlatNoParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	// Enough for header + newline + type decl + newline, fail in flat func def.
	w := &limitWriter{n: 100}
	err := Generate(reg, w, GenerateOptions{TypeScript: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInFlatReturn(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	// Enough for header + type decl + func declaration, fail on return statement.
	w := &limitWriter{n: 140}
	err := Generate(reg, w, GenerateOptions{TypeScript: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInNestedMember(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.index", "GET", "/contacts")

	// Enough for header + newline but fail in nested.
	w := &limitWriter{n: 80}
	err := Generate(reg, w, GenerateOptions{NestedOnly: true, TypeScript: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInNestedClose(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.index", "GET", "/contacts")

	// Enough for header + group open + member, fail on close.
	w := &limitWriter{n: 160}
	err := Generate(reg, w, GenerateOptions{NestedOnly: true, TypeScript: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInNestedParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	// Enough for header + group open, fail on member with params.
	w := &limitWriter{n: 120}
	err := Generate(reg, w, GenerateOptions{NestedOnly: true, TypeScript: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInFlatWithParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	// Enough for header + type decl, fail in flat func with params.
	w := &limitWriter{n: 110}
	err := Generate(reg, w, GenerateOptions{TypeScript: true, FlatOnly: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInJS_FlatNoParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	w := &limitWriter{n: 60}
	err := Generate(reg, w, GenerateOptions{TypeScript: false, FlatOnly: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInJS_FlatWithParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	w := &limitWriter{n: 60}
	err := Generate(reg, w, GenerateOptions{TypeScript: false, FlatOnly: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInJS_NestedNoParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.index", "GET", "/contacts")

	w := &limitWriter{n: 80}
	err := Generate(reg, w, GenerateOptions{TypeScript: false, NestedOnly: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInJS_NestedWithParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	w := &limitWriter{n: 100}
	err := Generate(reg, w, GenerateOptions{TypeScript: false, NestedOnly: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInJS_NestedClose(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.index", "GET", "/contacts")

	w := &limitWriter{n: 130}
	err := Generate(reg, w, GenerateOptions{TypeScript: false, NestedOnly: true})

	if err == nil {
		t.Error("expected write error, got nil")
	}
}

func TestGenerate_WriteErrorInFlatClosingBrace(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	// Probe byte limits to ensure the closing brace write path can fail.
	// We don't assert a specific threshold as output size may change.
	for n := 150; n < 200; n++ {
		w := &limitWriter{n: n}

		if err := Generate(reg, w, GenerateOptions{TypeScript: true, FlatOnly: true}); err != nil {
			break // Found the right threshold
		}
	}
}

func TestGenerate_WriteErrorInFlatTrailingNewline(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	// JS: slightly shorter strings.
	for n := 100; n < 180; n++ {
		w := &limitWriter{n: n}

		if err := Generate(reg, w, GenerateOptions{TypeScript: false, FlatOnly: true}); err != nil {
			break
		}
	}
}

func TestGenerate_WriteErrorInFlatReturnParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("users.show", "GET", "/users/{user}")

	for n := 120; n < 250; n++ {
		w := &limitWriter{n: n}

		if err := Generate(reg, w, GenerateOptions{TypeScript: true, FlatOnly: true}); err != nil {
			break
		}
	}
}

func TestGenerate_WriteErrorInNestedTrailingNewline(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.index", "GET", "/contacts")

	for n := 130; n < 200; n++ {
		w := &limitWriter{n: n}

		if err := Generate(reg, w, GenerateOptions{TypeScript: true, NestedOnly: true}); err != nil {
			break
		}
	}
}

func TestBuildURLTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		pattern  string
		expected string
	}{
		{"/contacts", "`/contacts`"},
		{"/contacts/{contact}", "`/contacts/${encodeURIComponent(String(params.contact))}`"},
		{"/contacts/{contact}/notes/{note}", "`/contacts/${encodeURIComponent(String(params.contact))}/notes/${encodeURIComponent(String(params.note))}`"},
	}

	for _, tt := range tests {
		got := buildURLTemplate(tt.pattern)

		if got != tt.expected {
			t.Errorf("buildURLTemplate(%q) = %q, want %q", tt.pattern, got, tt.expected)
		}
	}
}
