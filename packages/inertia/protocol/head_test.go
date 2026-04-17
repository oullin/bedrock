package protocol_test

import (
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
)

// --- RenderHTML ---

func TestRenderHTML_Empty(t *testing.T) {
	t.Parallel()

	h := protocol.Head{}

	if got := h.RenderHTML(); strings.TrimSpace(got) != "" {
		t.Errorf("empty Head should render empty string, got: %q", got)
	}
}

func TestRenderHTML_TitleOnly(t *testing.T) {
	t.Parallel()

	h := protocol.Head{Title: "My Page"}
	got := h.RenderHTML()

	if !strings.Contains(got, "<title>My Page</title>") {
		t.Errorf("expected title tag, got: %q", got)
	}
}

func TestRenderHTML_TitleEscapesHTML(t *testing.T) {
	t.Parallel()

	h := protocol.Head{Title: `Page "with" <special> & chars`}
	got := h.RenderHTML()

	if strings.Contains(got, "<special>") {
		t.Errorf("title should be escaped, got: %q", got)
	}

	if !strings.Contains(got, "&amp;") {
		t.Errorf("ampersand should be escaped, got: %q", got)
	}
}

func TestRenderHTML_MetaName(t *testing.T) {
	t.Parallel()

	h := protocol.Head{
		Meta: []protocol.MetaTag{
			{Name: "description", Content: "A test page"},
		},
	}
	got := h.RenderHTML()

	if !strings.Contains(got, `<meta name="description" content="A test page">`) {
		t.Errorf("expected meta name tag, got: %q", got)
	}
}

func TestRenderHTML_MetaProperty(t *testing.T) {
	t.Parallel()

	h := protocol.Head{
		Meta: []protocol.MetaTag{
			{Property: "og:title", Content: "OG Title"},
		},
	}
	got := h.RenderHTML()

	if !strings.Contains(got, `<meta property="og:title" content="OG Title">`) {
		t.Errorf("expected meta property tag, got: %q", got)
	}
}

func TestRenderHTML_MetaEmptyContentSkipped(t *testing.T) {
	t.Parallel()

	h := protocol.Head{
		Title: "Page",
		Meta: []protocol.MetaTag{
			{Name: "keywords", Content: ""},
			{Name: "description", Content: "Has content"},
		},
	}
	got := h.RenderHTML()

	if strings.Contains(got, "keywords") {
		t.Errorf("empty content meta should be skipped, got: %q", got)
	}

	if !strings.Contains(got, "description") {
		t.Errorf("non-empty content meta should be rendered, got: %q", got)
	}
}

func TestRenderHTML_MetaEscapesContent(t *testing.T) {
	t.Parallel()

	h := protocol.Head{
		Meta: []protocol.MetaTag{
			{Name: "description", Content: `He said "hello" & <goodbye>`},
		},
	}
	got := h.RenderHTML()

	if strings.Contains(got, `"hello"`) {
		t.Errorf("quotes should be escaped, got: %q", got)
	}
}

func TestRenderHTML_LinkTag(t *testing.T) {
	t.Parallel()

	h := protocol.Head{
		Links: []protocol.LinkTag{
			{Rel: "canonical", Href: "https://example.com/page"},
		},
	}
	got := h.RenderHTML()

	if !strings.Contains(got, `<link rel="canonical" href="https://example.com/page">`) {
		t.Errorf("expected canonical link, got: %q", got)
	}
}

func TestRenderHTML_LinkHreflang(t *testing.T) {
	t.Parallel()

	h := protocol.Head{
		Links: []protocol.LinkTag{
			{Rel: "alternate", Href: "/es/page", HrefLang: "es"},
		},
	}
	got := h.RenderHTML()

	if !strings.Contains(got, `hreflang="es"`) {
		t.Errorf("expected hreflang attribute, got: %q", got)
	}
}

func TestRenderHTML_LinkEmptyHrefSkipped(t *testing.T) {
	t.Parallel()

	h := protocol.Head{
		Title: "Page",
		Links: []protocol.LinkTag{
			{Rel: "canonical", Href: ""},
			{Rel: "alternate", Href: "/es", HrefLang: "es"},
		},
	}
	got := h.RenderHTML()

	if strings.Contains(got, "canonical") {
		t.Errorf("empty href link should be skipped, got: %q", got)
	}

	if !strings.Contains(got, "alternate") {
		t.Errorf("non-empty href link should be rendered, got: %q", got)
	}
}

func TestRenderHTML_LinkWithType(t *testing.T) {
	t.Parallel()

	h := protocol.Head{
		Links: []protocol.LinkTag{
			{Rel: "alternate", Href: "/feed", Type: "application/rss+xml"},
		},
	}
	got := h.RenderHTML()

	if !strings.Contains(got, `type="application/rss+xml"`) {
		t.Errorf("expected type attribute, got: %q", got)
	}
}

func TestRenderHTML_MultipleTags(t *testing.T) {
	t.Parallel()

	h := protocol.Head{
		Title: "Full Page",
		Meta: []protocol.MetaTag{
			{Name: "description", Content: "Desc"},
			{Property: "og:title", Content: "OG"},
		},
		Links: []protocol.LinkTag{
			{Rel: "canonical", Href: "https://example.com"},
		},
	}
	got := h.RenderHTML()

	if !strings.Contains(got, "<title>Full Page</title>") {
		t.Error("missing title")
	}

	if !strings.Contains(got, `name="description"`) {
		t.Error("missing description meta")
	}

	if !strings.Contains(got, `property="og:title"`) {
		t.Error("missing og:title meta")
	}

	if !strings.Contains(got, `rel="canonical"`) {
		t.Error("missing canonical link")
	}
}

// --- MergeHead ---

func TestMergeHead_EmptyBaseReturnsOverride(t *testing.T) {
	t.Parallel()

	override := protocol.Head{Title: "Override", Lang: "es"}
	result := protocol.MergeHead(protocol.Head{}, override)

	if result.Title != "Override" {
		t.Errorf("title = %q, want %q", result.Title, "Override")
	}

	if result.Lang != "es" {
		t.Errorf("lang = %q, want %q", result.Lang, "es")
	}
}

func TestMergeHead_EmptyOverrideReturnsBase(t *testing.T) {
	t.Parallel()

	base := protocol.Head{Title: "Base", Lang: "en", Direction: "ltr"}
	result := protocol.MergeHead(base, protocol.Head{})

	if result.Title != "Base" {
		t.Errorf("title = %q, want %q", result.Title, "Base")
	}

	if result.Lang != "en" {
		t.Errorf("lang = %q, want %q", result.Lang, "en")
	}
}

func TestMergeHead_OverrideTitleReplacesBase(t *testing.T) {
	t.Parallel()

	base := protocol.Head{Title: "Site Name"}
	override := protocol.Head{Title: "Page - Site Name"}
	result := protocol.MergeHead(base, override)

	if result.Title != "Page - Site Name" {
		t.Errorf("title = %q", result.Title)
	}
}

func TestMergeHead_OverrideLangReplacesBase(t *testing.T) {
	t.Parallel()

	base := protocol.Head{Lang: "en"}
	override := protocol.Head{Lang: "es"}
	result := protocol.MergeHead(base, override)

	if result.Lang != "es" {
		t.Errorf("lang = %q", result.Lang)
	}
}

func TestMergeHead_OverrideDirectionReplacesBase(t *testing.T) {
	t.Parallel()

	base := protocol.Head{Direction: "ltr"}
	override := protocol.Head{Direction: "rtl"}
	result := protocol.MergeHead(base, override)

	if result.Direction != "rtl" {
		t.Errorf("direction = %q", result.Direction)
	}
}

func TestMergeHead_MetaReplaceBySameName(t *testing.T) {
	t.Parallel()

	base := protocol.Head{
		Meta: []protocol.MetaTag{
			{Name: "description", Content: "Base desc"},
		},
	}
	override := protocol.Head{
		Meta: []protocol.MetaTag{
			{Name: "description", Content: "Override desc"},
		},
	}
	result := protocol.MergeHead(base, override)

	if len(result.Meta) != 1 {
		t.Fatalf("expected 1 meta tag, got %d", len(result.Meta))
	}

	if result.Meta[0].Content != "Override desc" {
		t.Errorf("content = %q", result.Meta[0].Content)
	}
}

func TestMergeHead_MetaReplaceBySameProperty(t *testing.T) {
	t.Parallel()

	base := protocol.Head{
		Meta: []protocol.MetaTag{
			{Property: "og:title", Content: "Base OG"},
		},
	}
	override := protocol.Head{
		Meta: []protocol.MetaTag{
			{Property: "og:title", Content: "Override OG"},
		},
	}
	result := protocol.MergeHead(base, override)

	if len(result.Meta) != 1 {
		t.Fatalf("expected 1 meta tag, got %d", len(result.Meta))
	}

	if result.Meta[0].Content != "Override OG" {
		t.Errorf("content = %q", result.Meta[0].Content)
	}
}

func TestMergeHead_MetaAppendNew(t *testing.T) {
	t.Parallel()

	base := protocol.Head{
		Meta: []protocol.MetaTag{
			{Name: "description", Content: "Desc"},
		},
	}
	override := protocol.Head{
		Meta: []protocol.MetaTag{
			{Name: "robots", Content: "noindex"},
		},
	}
	result := protocol.MergeHead(base, override)

	if len(result.Meta) != 2 {
		t.Fatalf("expected 2 meta tags, got %d", len(result.Meta))
	}
}

func TestMergeHead_LinkReplaceBySameKey(t *testing.T) {
	t.Parallel()

	base := protocol.Head{
		Links: []protocol.LinkTag{
			{Rel: "canonical", Href: "https://old.com"},
		},
	}
	override := protocol.Head{
		Links: []protocol.LinkTag{
			{Rel: "canonical", Href: "https://new.com"},
		},
	}
	result := protocol.MergeHead(base, override)

	if len(result.Links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(result.Links))
	}

	if result.Links[0].Href != "https://new.com" {
		t.Errorf("href = %q", result.Links[0].Href)
	}
}

func TestMergeHead_LinkHreflangMergesByLang(t *testing.T) {
	t.Parallel()

	base := protocol.Head{
		Links: []protocol.LinkTag{
			{Rel: "alternate", Href: "/en", HrefLang: "en"},
			{Rel: "alternate", Href: "/es", HrefLang: "es"},
		},
	}
	override := protocol.Head{
		Links: []protocol.LinkTag{
			{Rel: "alternate", Href: "/es-mx", HrefLang: "es"},
		},
	}
	result := protocol.MergeHead(base, override)

	if len(result.Links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(result.Links))
	}

	for _, link := range result.Links {
		if link.HrefLang == "es" && link.Href != "/es-mx" {
			t.Errorf("es link href = %q, want /es-mx", link.Href)
		}
	}
}

// --- ApplyEnv ---

func TestApplyEnv_OverridesTitle(t *testing.T) {
	t.Setenv("INERTIA_SEO_TITLE", "Env Title")

	h := protocol.Head{Title: "YAML Title"}

	h.ApplyEnv()

	if h.Title != "Env Title" {
		t.Errorf("title = %q, want %q", h.Title, "Env Title")
	}
}

func TestApplyEnv_OverridesLang(t *testing.T) {
	t.Setenv("INERTIA_SEO_LANG", "fr")

	h := protocol.Head{Lang: "en"}

	h.ApplyEnv()

	if h.Lang != "fr" {
		t.Errorf("lang = %q, want %q", h.Lang, "fr")
	}
}

func TestApplyEnv_OverridesExistingMeta(t *testing.T) {
	t.Setenv("INERTIA_SEO_DESCRIPTION", "Env description")

	h := protocol.Head{
		Meta: []protocol.MetaTag{
			{Name: "description", Content: "YAML description"},
		},
	}

	h.ApplyEnv()

	if h.Meta[0].Content != "Env description" {
		t.Errorf("description = %q", h.Meta[0].Content)
	}
}

func TestApplyEnv_AddsNewMetaWhenMissing(t *testing.T) {
	t.Setenv("INERTIA_SEO_OG_TITLE", "Env OG Title")

	h := protocol.Head{}

	h.ApplyEnv()

	found := false

	for _, tag := range h.Meta {
		if tag.Property == "og:title" && tag.Content == "Env OG Title" {
			found = true
		}
	}

	if !found {
		t.Error("expected og:title to be added from env")
	}
}

func TestApplyEnv_OverridesExistingPropertyMeta(t *testing.T) {
	t.Setenv("INERTIA_SEO_OG_TITLE", "Env OG Override")

	h := protocol.Head{
		Meta: []protocol.MetaTag{
			{Property: "og:title", Content: "YAML OG Title"},
		},
	}

	h.ApplyEnv()

	if h.Meta[0].Content != "Env OG Override" {
		t.Errorf("og:title = %q, want %q", h.Meta[0].Content, "Env OG Override")
	}
}

func TestApplyEnv_NoOpWhenEnvEmpty(t *testing.T) {
	t.Setenv("INERTIA_SEO_TITLE", "")

	h := protocol.Head{Title: "Keep"}

	h.ApplyEnv()

	if h.Title != "Keep" {
		t.Errorf("title should not change, got %q", h.Title)
	}
}
