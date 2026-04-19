package inertia

import (
	"fmt"

	"github.com/bedrock/packages/config"
	"github.com/bedrock/packages/inertia/protocol"
)

// DefaultHead returns a Head with sensible defaults: lang "en", robots
// "index, follow", and placeholder slots for common meta and link tags.
// Empty Content values are skipped during rendering.
func DefaultHead() protocol.Head {
	head := protocol.Head{
		Lang: "en",
		Meta: []protocol.MetaTag{
			{Name: "description", Content: ""},
			{Name: "keywords", Content: ""},
			{Name: "robots", Content: "index, follow"},
			{Property: "og:title", Content: ""},
			{Property: "og:description", Content: ""},
			{Property: "og:image", Content: ""},
			{Property: "og:url", Content: ""},
			{Property: "og:type", Content: "website"},
			{Property: "og:site_name", Content: ""},
			{Property: "og:locale", Content: "en_US"},
			{Name: "twitter:card", Content: ""},
			{Name: "twitter:title", Content: ""},
			{Name: "twitter:description", Content: ""},
			{Name: "twitter:image", Content: ""},
			{Name: "twitter:site", Content: ""},
		},
		Links: []protocol.LinkTag{
			{Rel: "canonical", Href: ""},
			{Rel: "alternate", Href: "", HrefLang: ""},
		},
	}

	head.ApplyEnv()

	return head
}

// LoadHead reads a YAML head/SEO config file. Defaults are applied first,
// then the file values are merged on top, and finally environment variable
// overrides (INERTIA_SEO_*) are applied.
func LoadHead(path string) (protocol.Head, error) {
	repo := config.NewWithDefaults(map[string]any{"lang": "en"})

	v := repo.Viper()

	v.SetEnvPrefix("INERTIA_SEO")

	v.AutomaticEnv()

	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return protocol.Head{}, fmt.Errorf("head: read config: %w", err)
	}

	var override protocol.Head

	if err := v.Unmarshal(&override); err != nil {
		return protocol.Head{}, fmt.Errorf("head: parse config: %w", err)
	}

	head := protocol.MergeHead(DefaultHead(), override)

	head.ApplyEnv()

	return head, nil
}
