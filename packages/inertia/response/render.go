package response

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/bedrock/packages/inertia/protocol"
)

// HTMLConfig groups the rendering parameters for WriteHTML, keeping the
// function signature small and making it easy to add fields later.
type HTMLConfig struct {
	Template    *template.Template
	ContainerID string
	Marshaler   protocol.JSONMarshaler
	ExtraData   protocol.TemplateData
	Head        protocol.Head
}

// WriteJSON writes an Inertia JSON response for XHR visits.
// It sets the Content-Type and X-Inertia headers.
func WriteJSON(w http.ResponseWriter, page *Page, marshaler protocol.JSONMarshaler) error {
	data, err := marshaler.Marshal(page)

	if err != nil {
		return fmt.Errorf("inertia: marshal page: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(protocol.HeaderInertia, "true")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(data)

	return err
}

// WriteHTML executes the root HTML template for initial (non-XHR) page
// visits. The page object is embedded as JSON inside a
// <script type="application/json" data-page="ID"> element, followed by
// an empty container div for client-side mounting.
func WriteHTML(w http.ResponseWriter, page *Page, cfg HTMLConfig) error {
	data, err := cfg.Marshaler.Marshal(page)

	if err != nil {
		return fmt.Errorf("inertia: marshal page: %w", err)
	}

	// Escape "</" to prevent premature </script> closure. Go's
	// encoding/json already escapes "<" as \u003c, but custom
	// marshalers may not.
	safeJSON := strings.ReplaceAll(string(data), "</", `<\/`)

	inertiaHTML := fmt.Sprintf(
		`<script data-page="%s" type="application/json">%s</script><div id="%s"></div>`,
		cfg.ContainerID,
		safeJSON,
		cfg.ContainerID,
	)

	templateData := map[string]any{
		"inertia":     template.HTML(inertiaHTML),
		"inertiaHead": template.HTML(cfg.Head.RenderHTML()),
	}

	if strings.TrimSpace(cfg.Head.Lang) != "" {
		templateData["inertiaLang"] = cfg.Head.Lang
	}

	if strings.TrimSpace(cfg.Head.Direction) != "" {
		templateData["inertiaDir"] = cfg.Head.Direction
	}

	for k, v := range cfg.ExtraData {
		templateData[k] = v
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	return cfg.Template.Execute(w, templateData)
}
