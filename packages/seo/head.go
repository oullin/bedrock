// Package seo defines the SEO/locale primitives shared across web-rendering
// packages: <head> tags, locale metadata, and the context key used to carry
// a resolved locale through a request. It has no third-party dependencies.
package seo

import (
	"html"
	"os"
	"strings"
)

// MetaTag represents a single <meta> element for server-side head rendering.
// Exactly one of Name or Property should be set to identify the tag.
type MetaTag struct {
	Name     string `json:"name,omitempty"     yaml:"name,omitempty"     mapstructure:"name"`
	Property string `json:"property,omitempty" yaml:"property,omitempty" mapstructure:"property"`
	Content  string `json:"content"            yaml:"content"            mapstructure:"content"`
}

// LinkTag represents a single <link> element for server-side head rendering.
type LinkTag struct {
	Rel      string `json:"rel"                yaml:"rel"                mapstructure:"rel"`
	Href     string `json:"href"               yaml:"href"               mapstructure:"href"`
	HrefLang string `json:"hreflang,omitempty"  yaml:"hreflang,omitempty"  mapstructure:"hreflang"`
	Type     string `json:"type,omitempty"      yaml:"type,omitempty"      mapstructure:"type"`
}

// Head holds the server-side head elements rendered into the page's <head>
// on initial page loads, plus the lang/dir attributes for the <html> element.
type Head struct {
	Title     string    `json:"title,omitempty"     yaml:"title,omitempty"     mapstructure:"title"`
	Lang      string    `json:"lang,omitempty"      yaml:"lang,omitempty"      mapstructure:"lang"`
	Direction string    `json:"direction,omitempty"  yaml:"direction,omitempty"  mapstructure:"direction"`
	Meta      []MetaTag `json:"meta,omitempty"      yaml:"meta,omitempty"      mapstructure:"meta"`
	Links     []LinkTag `json:"links,omitempty"     yaml:"links,omitempty"     mapstructure:"links"`
}

// Locale holds locale information typically resolved by an i18n middleware.
type Locale struct {
	Code      string `yaml:"-"         mapstructure:"-"`
	Name      string `yaml:"name"      mapstructure:"name"`
	Direction string `yaml:"direction" mapstructure:"direction"`
	Head      Head   `yaml:"head"      mapstructure:"head"`
}

var envBindings = []struct {
	envSuffix string
	name      string
	property  string
}{
	{"DESCRIPTION", "description", ""},
	{"KEYWORDS", "keywords", ""},
	{"ROBOTS", "robots", ""},
	{"GOOGLEBOT", "googlebot", ""},
	{"OG_TITLE", "", "og:title"},
	{"OG_DESCRIPTION", "", "og:description"},
	{"OG_IMAGE", "", "og:image"},
	{"OG_URL", "", "og:url"},
	{"OG_TYPE", "", "og:type"},
	{"OG_SITE_NAME", "", "og:site_name"},
	{"OG_LOCALE", "", "og:locale"},
	{"TWITTER_CARD", "twitter:card", ""},
	{"TWITTER_TITLE", "twitter:title", ""},
	{"TWITTER_DESCRIPTION", "twitter:description", ""},
	{"TWITTER_IMAGE", "twitter:image", ""},
	{"TWITTER_SITE", "twitter:site", ""},
}

func metaKey(tag MetaTag) string {
	if strings.TrimSpace(tag.Name) != "" {
		return "name:" + tag.Name
	}

	return "property:" + tag.Property
}

func linkKey(tag LinkTag) string {
	if strings.TrimSpace(tag.HrefLang) != "" {
		return tag.Rel + ":" + tag.HrefLang
	}

	return tag.Rel
}

// RenderHTML produces the HTML string for embedding in a page's <head>.
// Meta tags with empty Content are skipped (they serve as YAML placeholders).
// Returns "" for a zero-value Head.
func (h Head) RenderHTML() string {
	if strings.TrimSpace(h.Title) == "" && len(h.Meta) == 0 && len(h.Links) == 0 {
		return ""
	}

	var b strings.Builder

	if strings.TrimSpace(h.Title) != "" {
		b.WriteString("<title>")
		b.WriteString(html.EscapeString(h.Title))
		b.WriteString("</title>\n")
	}

	for _, tag := range h.Meta {
		if strings.TrimSpace(tag.Content) == "" {
			continue
		}

		b.WriteString("<meta ")

		switch {
		case strings.TrimSpace(tag.Name) != "":
			b.WriteString(`name="`)
			b.WriteString(html.EscapeString(tag.Name))
			b.WriteString(`"`)
		case strings.TrimSpace(tag.Property) != "":
			b.WriteString(`property="`)
			b.WriteString(html.EscapeString(tag.Property))
			b.WriteString(`"`)
		}

		b.WriteString(` content="`)
		b.WriteString(html.EscapeString(tag.Content))
		b.WriteString("\">\n")
	}

	for _, link := range h.Links {
		if strings.TrimSpace(link.Href) == "" {
			continue
		}

		b.WriteString(`<link rel="`)
		b.WriteString(html.EscapeString(link.Rel))
		b.WriteString(`" href="`)
		b.WriteString(html.EscapeString(link.Href))
		b.WriteString(`"`)

		if strings.TrimSpace(link.HrefLang) != "" {
			b.WriteString(` hreflang="`)
			b.WriteString(html.EscapeString(link.HrefLang))
			b.WriteString(`"`)
		}

		if strings.TrimSpace(link.Type) != "" {
			b.WriteString(` type="`)
			b.WriteString(html.EscapeString(link.Type))
			b.WriteString(`"`)
		}

		b.WriteString(">\n")
	}

	return b.String()
}

// MergeHead combines a base Head with an override Head. The override's
// Title, Lang, and Direction replace the base values if non-empty. Meta
// tags are merged by their identifying attribute (Name or Property);
// override tags replace base tags with the same key, and new tags are
// appended. Links are merged by Rel+HrefLang key.
func MergeHead(base, override Head) Head {
	result := Head{
		Title:     base.Title,
		Lang:      base.Lang,
		Direction: base.Direction,
	}

	if strings.TrimSpace(override.Title) != "" {
		result.Title = override.Title
	}

	if strings.TrimSpace(override.Lang) != "" {
		result.Lang = override.Lang
	}

	if strings.TrimSpace(override.Direction) != "" {
		result.Direction = override.Direction
	}

	metaSeen := make(map[string]int, len(base.Meta))
	merged := make([]MetaTag, len(base.Meta))

	copy(merged, base.Meta)

	for i, tag := range merged {
		metaSeen[metaKey(tag)] = i
	}

	for _, tag := range override.Meta {
		key := metaKey(tag)

		if idx, ok := metaSeen[key]; ok {
			merged[idx] = tag
		} else {
			metaSeen[key] = len(merged)
			merged = append(merged, tag)
		}
	}

	result.Meta = merged

	linkSeen := make(map[string]int, len(base.Links))
	mergedLinks := make([]LinkTag, len(base.Links))

	copy(mergedLinks, base.Links)

	for i, link := range mergedLinks {
		linkSeen[linkKey(link)] = i
	}

	for _, link := range override.Links {
		key := linkKey(link)

		if idx, ok := linkSeen[key]; ok {
			mergedLinks[idx] = link
		} else {
			linkSeen[key] = len(mergedLinks)
			mergedLinks = append(mergedLinks, link)
		}
	}

	result.Links = mergedLinks

	return result
}

// ApplyEnv overrides Head fields with values from environment variables
// when present. Uses the INERTIA_SEO_ prefix convention.
func (h *Head) ApplyEnv() {
	if v := os.Getenv("INERTIA_SEO_TITLE"); strings.TrimSpace(v) != "" {
		h.Title = v
	}

	if v := os.Getenv("INERTIA_SEO_LANG"); strings.TrimSpace(v) != "" {
		h.Lang = v
	}

	for _, binding := range envBindings {
		v := os.Getenv("INERTIA_SEO_" + binding.envSuffix)

		if strings.TrimSpace(v) == "" {
			continue
		}

		found := false

		for i := range h.Meta {
			if strings.TrimSpace(binding.name) != "" && h.Meta[i].Name == binding.name {
				h.Meta[i].Content = v
				found = true

				break
			}

			if strings.TrimSpace(binding.property) != "" && h.Meta[i].Property == binding.property {
				h.Meta[i].Content = v
				found = true

				break
			}
		}

		if !found {
			tag := MetaTag{Content: v}

			if strings.TrimSpace(binding.name) != "" {
				tag.Name = binding.name
			} else {
				tag.Property = binding.property
			}

			h.Meta = append(h.Meta, tag)
		}
	}
}
