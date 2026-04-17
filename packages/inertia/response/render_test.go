package response_test

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/inertia/response"
)

var errMarshal = &marshalError{}

type testMarshaler struct{}

type failMarshaler struct{}

type marshalError struct{}

func (m *testMarshaler) Marshal(v any) ([]byte, error)   { return json.Marshal(v) }
func (m *testMarshaler) Unmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }

func (m *failMarshaler) Marshal(v any) ([]byte, error)   { return nil, errMarshal }
func (m *failMarshaler) Unmarshal(b []byte, v any) error { return errMarshal }

func (e *marshalError) Error() string { return "marshal failed" }

// --- WriteJSON ---

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	page := &response.Page{
		Component: "Users/Index",
		Props:     map[string]any{"name": "alice"},
		URL:       "/users",
		Version:   "v1",
	}

	w := httptest.NewRecorder()
	err := response.WriteJSON(w, page, &testMarshaler{})

	if err != nil {
		t.Fatal(err)
	}

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	if resp.Header.Get(protocol.HeaderInertia) != "true" {
		t.Error("missing X-Inertia header")
	}

	var got response.Page

	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got.Component != "Users/Index" {
		t.Errorf("component = %q, want %q", got.Component, "Users/Index")
	}

	if got.URL != "/users" {
		t.Errorf("url = %q, want %q", got.URL, "/users")
	}
}

func TestWriteJSON_MarshalError(t *testing.T) {
	t.Parallel()

	page := &response.Page{Component: "Page"}

	w := httptest.NewRecorder()
	err := response.WriteJSON(w, page, &failMarshaler{})

	if err == nil {
		t.Error("expected error from marshal failure")
	}
}

// --- WriteHTML ---

func TestWriteHTML(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(
		`<!DOCTYPE html><html><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`,
	))

	page := &response.Page{
		Component: "Home",
		Props:     map[string]any{"title": "Welcome"},
		URL:       "/",
		Version:   "v1",
	}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "app", Marshaler: &testMarshaler{}})

	if err != nil {
		t.Fatal(err)
	}

	resp := w.Result()

	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", ct, "text/html; charset=utf-8")
	}

	body := w.Body.String()

	if !contains(body, `<div id="app"></div>`) {
		t.Error("missing empty container div")
	}

	if !contains(body, `<script data-page="app" type="application/json">`) {
		t.Error("missing page data script element")
	}

	if !contains(body, `"component":"Home"`) {
		t.Error("missing component in script JSON")
	}
}

func TestWriteHTML_MarshalError(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))
	page := &response.Page{Component: "Page"}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "app", Marshaler: &failMarshaler{}})

	if err == nil {
		t.Error("expected error from marshal failure")
	}
}

func TestWriteHTML_ExtraTemplateData(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(
		`{{ .inertia }}|title={{ .pageTitle }}`,
	))

	page := &response.Page{
		Component: "Home",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
	}

	extra := protocol.TemplateData{"pageTitle": "My App"}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "app", Marshaler: &testMarshaler{}, ExtraData: extra})

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if !contains(body, "title=My App") {
		t.Errorf("body missing %q, got %s", "title=My App", body)
	}
}

func TestWriteHTML_ScriptClosingTagEscaped(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))

	// Props containing "</script>" should be escaped to prevent injection.
	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{"html": "</script>"},
		URL:       "/",
		Version:   "v1",
	}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "app", Marshaler: &testMarshaler{}})

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	// Go's json.Marshal escapes "<" as \u003c, but the explicit
	// strings.ReplaceAll("</" -> "<\/") acts as a safety net for
	// custom marshalers.
	if contains(body, "</script></script>") {
		t.Error("raw </script> found in output; should be escaped")
	}
}

func TestWriteHTML_CustomContainerID(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))

	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
	}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "root", Marshaler: &testMarshaler{}})

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if !contains(body, `<div id="root"></div>`) {
		t.Error("missing container div with custom ID")
	}

	if !contains(body, `<script data-page="root"`) {
		t.Error("missing script element with custom container ID")
	}
}

func TestWriteHTML_OmitsEmptyOptionalFields(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))

	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
	}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "app", Marshaler: &testMarshaler{}})

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if contains(body, `"encryptHistory"`) {
		t.Error("encryptHistory should be omitted when false")
	}

	if contains(body, `"clearHistory"`) {
		t.Error("clearHistory should be omitted when false")
	}

	if contains(body, `"mergeProps"`) {
		t.Error("mergeProps should be omitted when empty")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}

// --- WriteJSON Extended ---

func TestWriteJSON_IncludesMergePropsField(t *testing.T) {
	t.Parallel()

	page := &response.Page{
		Component:  "Page",
		Props:      map[string]any{"items": []int{1}},
		URL:        "/",
		Version:    "v1",
		MergeProps: []string{"items"},
	}

	w := httptest.NewRecorder()

	response.WriteJSON(w, page, &testMarshaler{})

	body := w.Body.String()

	if !contains(body, `"mergeProps":["items"]`) {
		t.Errorf("JSON missing %q, got %s", `"mergeProps":["items"]`, body)
	}
}

func TestWriteJSON_IncludesDeepMergePropsField(t *testing.T) {
	t.Parallel()

	page := &response.Page{
		Component:      "Page",
		Props:          map[string]any{},
		URL:            "/",
		Version:        "v1",
		DeepMergeProps: []string{"data"},
	}

	w := httptest.NewRecorder()

	response.WriteJSON(w, page, &testMarshaler{})

	body := w.Body.String()

	if !contains(body, `"deepMergeProps":["data"]`) {
		t.Errorf("JSON missing %q, got %s", `"deepMergeProps":["data"]`, body)
	}
}

func TestWriteJSON_IncludesDeferredPropsField(t *testing.T) {
	t.Parallel()

	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
		DeferredProps: map[string][]string{
			"sidebar": {"stats"},
		},
	}

	w := httptest.NewRecorder()

	response.WriteJSON(w, page, &testMarshaler{})

	body := w.Body.String()

	if !contains(body, `"deferredProps"`) {
		t.Errorf("JSON missing %q, got %s", `"deferredProps"`, body)
	}

	if !contains(body, `"sidebar"`) {
		t.Errorf("JSON missing %q, got %s", `"sidebar"`, body)
	}
}

func TestWriteJSON_IncludesScrollPropsField(t *testing.T) {
	t.Parallel()

	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
		ScrollProps: map[string]response.Scroll{
			"feed": {
				PageName:     "feedPage",
				PreviousPage: nil,
				NextPage:     2,
				CurrentPage:  1,
			},
		},
	}

	w := httptest.NewRecorder()

	response.WriteJSON(w, page, &testMarshaler{})

	body := w.Body.String()

	if !contains(body, `"scrollProps"`) {
		t.Errorf("JSON missing %q, got %s", `"scrollProps"`, body)
	}

	if !contains(body, `"feedPage"`) {
		t.Errorf("JSON missing %q, got %s", `"feedPage"`, body)
	}
}

func TestWriteJSON_IncludesOncePropsField(t *testing.T) {
	t.Parallel()

	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
		OnceProps: map[string]response.Once{
			"notes": {Prop: "notes"},
		},
	}

	w := httptest.NewRecorder()

	response.WriteJSON(w, page, &testMarshaler{})

	body := w.Body.String()

	if !contains(body, `"onceProps"`) {
		t.Errorf("JSON missing %q, got %s", `"onceProps"`, body)
	}

	// ExpiresAt is nil so should be omitted.
	if contains(body, `"expiresAt"`) {
		t.Errorf("JSON should not contain %q, got %s", `"expiresAt"`, body)
	}
}

func TestWriteJSON_OncePropsWithExpiresAt(t *testing.T) {
	t.Parallel()

	expires := int64(1700000000)

	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
		OnceProps: map[string]response.Once{
			"notes": {Prop: "notes", ExpiresAt: &expires},
		},
	}

	w := httptest.NewRecorder()

	response.WriteJSON(w, page, &testMarshaler{})

	body := w.Body.String()

	if !contains(body, `"expiresAt"`) {
		t.Errorf("JSON missing %q, got %s", `"expiresAt"`, body)
	}
}

func TestWriteJSON_EncryptHistoryTrue(t *testing.T) {
	t.Parallel()

	page := &response.Page{
		Component:      "Page",
		Props:          map[string]any{},
		URL:            "/",
		Version:        "v1",
		EncryptHistory: true,
	}

	w := httptest.NewRecorder()

	response.WriteJSON(w, page, &testMarshaler{})

	body := w.Body.String()

	if !contains(body, `"encryptHistory":true`) {
		t.Errorf("JSON missing %q, got %s", `"encryptHistory":true`, body)
	}
}

func TestWriteJSON_ClearHistoryTrue(t *testing.T) {
	t.Parallel()

	page := &response.Page{
		Component:    "Page",
		Props:        map[string]any{},
		URL:          "/",
		Version:      "v1",
		ClearHistory: true,
	}

	w := httptest.NewRecorder()

	response.WriteJSON(w, page, &testMarshaler{})

	body := w.Body.String()

	if !contains(body, `"clearHistory":true`) {
		t.Errorf("JSON missing %q, got %s", `"clearHistory":true`, body)
	}
}

func TestWriteJSON_AllFieldsOmittedWhenEmpty(t *testing.T) {
	t.Parallel()

	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
	}

	w := httptest.NewRecorder()

	response.WriteJSON(w, page, &testMarshaler{})

	body := w.Body.String()

	for _, field := range []string{
		`"encryptHistory"`, `"clearHistory"`,
		`"mergeProps"`, `"deepMergeProps"`,
		`"deferredProps"`, `"scrollProps"`, `"onceProps"`,
	} {
		if contains(body, field) {
			t.Errorf("field %s should be omitted when empty: %s", field, body)
		}
	}
}

// --- WriteHTML Extended ---

func TestWriteHTML_NilExtraDataHandled(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))
	page := &response.Page{Component: "Page", Props: map[string]any{}, URL: "/", Version: "v1"}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "app", Marshaler: &testMarshaler{}})

	if err != nil {
		t.Fatalf("WriteHTML with nil extraData should not panic: %v", err)
	}
}

func TestWriteHTML_EmptyExtraData(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))
	page := &response.Page{Component: "Page", Props: map[string]any{}, URL: "/", Version: "v1"}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "app", Marshaler: &testMarshaler{}, ExtraData: protocol.TemplateData{}})

	if err != nil {
		t.Fatalf("WriteHTML with empty extraData: %v", err)
	}
}

func TestWriteHTML_MultipleExtraDataKeys(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(
		`{{ .inertia }}|{{ .keyA }}|{{ .keyB }}`,
	))

	page := &response.Page{Component: "Page", Props: map[string]any{}, URL: "/", Version: "v1"}
	extra := protocol.TemplateData{"keyA": "valA", "keyB": "valB"}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "app", Marshaler: &testMarshaler{}, ExtraData: extra})

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if !contains(body, "valA") || !contains(body, "valB") {
		t.Errorf("body missing %q or %q, got %s", "valA", "valB", body)
	}
}

func TestWriteHTML_EmbeddedJSONIsValid(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))
	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{"key": "value"},
		URL:       "/",
		Version:   "v1",
	}

	w := httptest.NewRecorder()

	response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "app", Marshaler: &testMarshaler{}})

	body := w.Body.String()

	// Extract JSON between script tags.
	start := `<script data-page="app" type="application/json">`
	end := `</script>`

	startIdx := len(start)
	endIdx := len(body) - len(end) - len(`<div id="app"></div>`)

	if startIdx >= endIdx || !contains(body, start) {
		t.Fatalf("could not find script tags in: %s", body)
	}

	// Find actual indices.
	jsonStart := -1
	jsonEnd := -1

	for i := 0; i <= len(body)-len(start); i++ {
		if body[i:i+len(start)] == start {
			jsonStart = i + len(start)

			break
		}
	}

	for i := jsonStart; i <= len(body)-len(end); i++ {
		if body[i:i+len(end)] == end {
			jsonEnd = i

			break
		}
	}

	if jsonStart < 0 || jsonEnd < 0 {
		t.Fatalf("could not extract JSON from: %s", body)
	}

	// Unescape the <\/ back to </ for JSON parsing.
	jsonStr := body[jsonStart:jsonEnd]

	var parsed map[string]any

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Errorf("embedded JSON is invalid: %v (json: %s)", err, jsonStr)
	}
}

func TestWriteHTML_InertiaHeadIsEmpty(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(`head=[{{ .inertiaHead }}]`))
	page := &response.Page{Component: "Page", Props: map[string]any{}, URL: "/", Version: "v1"}

	w := httptest.NewRecorder()

	response.WriteHTML(w, page, response.HTMLConfig{Template: tmpl, ContainerID: "app", Marshaler: &testMarshaler{}})

	body := w.Body.String()

	if !contains(body, "head=[]") {
		t.Errorf("inertiaHead should be empty, got: %s", body)
	}
}

func TestWriteHTML_WithHead(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(
		`<head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body>`,
	))

	page := &response.Page{Component: "Page", Props: map[string]any{}, URL: "/", Version: "v1"}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{
		Template:    tmpl,
		ContainerID: "app",
		Marshaler:   &testMarshaler{},
		Head: protocol.Head{
			Title: "Test Page",
			Meta: []protocol.MetaTag{
				{Name: "description", Content: "A test"},
				{Property: "og:title", Content: "OG Test"},
			},
			Links: []protocol.LinkTag{
				{Rel: "canonical", Href: "https://example.com"},
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if !contains(body, "<title>Test Page</title>") {
		t.Errorf("body missing %q, got %s", "<title>Test Page</title>", body)
	}

	if !contains(body, `name="description"`) {
		t.Errorf("body missing %q, got %s", `name="description"`, body)
	}

	if !contains(body, `property="og:title"`) {
		t.Errorf("body missing %q, got %s", `property="og:title"`, body)
	}

	if !contains(body, `rel="canonical"`) {
		t.Errorf("body missing %q, got %s", `rel="canonical"`, body)
	}
}

func TestWriteHTML_WithLang(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(
		`<html lang="{{ .inertiaLang }}">{{ .inertia }}</html>`,
	))

	page := &response.Page{Component: "Page", Props: map[string]any{}, URL: "/", Version: "v1"}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{
		Template:    tmpl,
		ContainerID: "app",
		Marshaler:   &testMarshaler{},
		Head:        protocol.Head{Lang: "es"},
	})

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if !contains(body, `lang="es"`) {
		t.Errorf("body missing %q, got %s", `lang="es"`, body)
	}
}

func TestWriteHTML_WithDir(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(
		`<html dir="{{ .inertiaDir }}">{{ .inertia }}</html>`,
	))

	page := &response.Page{Component: "Page", Props: map[string]any{}, URL: "/", Version: "v1"}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, page, response.HTMLConfig{
		Template:    tmpl,
		ContainerID: "app",
		Marshaler:   &testMarshaler{},
		Head:        protocol.Head{Direction: "rtl"},
	})

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if !contains(body, `dir="rtl"`) {
		t.Errorf("body missing %q, got %s", `dir="rtl"`, body)
	}
}

func TestWriteHTML_EmptyHeadBackwardCompat(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("root").Parse(`head=[{{ .inertiaHead }}]`))
	page := &response.Page{Component: "Page", Props: map[string]any{}, URL: "/", Version: "v1"}

	w := httptest.NewRecorder()

	response.WriteHTML(w, page, response.HTMLConfig{
		Template:    tmpl,
		ContainerID: "app",
		Marshaler:   &testMarshaler{},
		Head:        protocol.Head{},
	})

	body := w.Body.String()

	if !contains(body, "head=[]") {
		t.Errorf("empty Head should render empty inertiaHead, got: %s", body)
	}
}
