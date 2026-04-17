package inertia

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"os"

	"github.com/bedrock/packages/inertia/config"
	"github.com/bedrock/packages/inertia/protocol"
)

// Option configures an Inertia instance during construction.

const (
	headSourceNone headSource = iota
	headSourceConfig
	headSourceExplicit
)

type headSource uint8

type Option func(*Inertia) error

// WithVersion sets a static asset version string.
func WithVersion(version string) Option {
	return func(i *Inertia) error {
		i.version = version

		return nil
	}
}

// WithVersionFromFile hashes the contents of path to produce the
// asset version string. This is useful for cache-busting when a
// manifest file (e.g. build/manifest.json) changes on deploy.
func WithVersionFromFile(path string) Option {
	return func(i *Inertia) error {
		data, err := os.ReadFile(path)

		if err != nil {
			return fmt.Errorf("inertia: version file: %w", err)
		}

		i.version = fmt.Sprintf("%x", md5.Sum(data))

		return nil
	}
}

// WithContainerID sets the HTML element ID used for the root container
// div. Defaults to "app".
func WithContainerID(id string) Option {
	return func(i *Inertia) error {
		i.containerID = id

		return nil
	}
}

// WithJSONMarshaler replaces the default encoding/json marshaler with
// a custom implementation.
func WithJSONMarshaler(m protocol.JSONMarshaler) Option {
	return func(i *Inertia) error {
		i.jsonMarshaler = m

		return nil
	}
}

// WithLogger sets a logger for diagnostic messages.
func WithLogger(l protocol.Logger) Option {
	return func(i *Inertia) error {
		i.logger = l

		return nil
	}
}

// WithTemplateFuncs registers additional template functions available
// in the root HTML template.
func WithTemplateFuncs(funcMap template.FuncMap) Option {
	return func(i *Inertia) error {
		i.templateFuncs = funcMap

		return nil
	}
}

// WithEncryptHistory enables encrypted browser history by default for
// all responses.
func WithEncryptHistory() Option {
	return func(i *Inertia) error {
		i.encryptHistory = true

		return nil
	}
}

// WithHead sets default head elements rendered into {{ .inertiaHead }} on
// every initial page load. Per-request head elements (set via SetHead,
// SetTitle, or SetMeta) override these defaults.
func WithHead(head protocol.Head) Option {
	return func(i *Inertia) error {
		i.head = head
		i.headSource = headSourceExplicit

		return nil
	}
}

// WithHeadDefaults sets default head elements from the built-in defaults.
// Environment variable overrides (INERTIA_SEO_*) are applied. Meta tags
// with empty Content serve as placeholders and are excluded from rendering.
func WithHeadDefaults() Option {
	return func(i *Inertia) error {
		if i.headSource == headSourceExplicit {
			return nil
		}

		i.head = config.DefaultHead()
		i.headSource = headSourceConfig

		return nil
	}
}

// WithHeadFromFile reads a YAML file at path and sets the default head
// elements. Defaults are applied first, then the file values are merged
// on top, and finally environment variable overrides (INERTIA_SEO_*) are
// applied. Meta tags with empty Content serve as placeholders and are
// excluded from rendering.
func WithHeadFromFile(path string) Option {
	return func(i *Inertia) error {
		if i.headSource == headSourceExplicit {
			return nil
		}

		head, err := config.LoadHead(path)

		if err != nil {
			return fmt.Errorf("inertia: %w", err)
		}

		i.head = head
		i.headSource = headSourceConfig

		return nil
	}
}
