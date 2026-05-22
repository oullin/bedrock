package str

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// StrMarkdown converts the given CommonMark Markdown string to HTML.
// Options map keys:
//   - "html_input": "strip" (default) or "allow" — controls raw HTML in input
//   - "allow_unsafe_links": bool — whether to allow unsafe links (default false)
func StrMarkdown(str string, options ...map[string]any) string {
	opts := mergeMarkdownOptions(options)

	unsafe := false

	if allow, ok := opts["html_input"].(string); ok && allow == "allow" {
		unsafe = true
	}

	if u, ok := opts["allow_unsafe_links"].(bool); ok && u {
		unsafe = true
	}

	mdOpts := []goldmark.Option{
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	}

	if unsafe {
		mdOpts = append(mdOpts, goldmark.WithRendererOptions(html.WithUnsafe()))
	}

	md := goldmark.New(mdOpts...)

	var buf bytes.Buffer

	if err := md.Convert([]byte(str), &buf); err != nil {
		return str
	}

	return buf.String()
}

// StrInlineMarkdown converts the given CommonMark Markdown string to inline HTML.
// Unlike StrMarkdown, the wrapping <p> tag is removed for single-line content.
func StrInlineMarkdown(str string, options ...map[string]any) string {
	result := StrMarkdown(str, options...)

	// Remove wrapping <p> tags for inline content
	result = strings.TrimSpace(result)

	if strings.HasPrefix(result, "<p>") && strings.HasSuffix(result, "</p>") {
		inner := result[3 : len(result)-4]
		// Only strip if there's a single paragraph (no nested block elements)
		if !strings.Contains(inner, "<p>") {
			result = inner
		}
	}

	return result
}

func mergeMarkdownOptions(options []map[string]any) map[string]any {
	merged := map[string]any{
		"html_input":         "strip",
		"allow_unsafe_links": false,
	}

	if len(options) > 0 {
		for k, v := range options[0] {
			merged[k] = v
		}
	}

	return merged
}
