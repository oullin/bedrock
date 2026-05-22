package handler

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/httpx"
	"github.com/bedrock/packages/wayfinder"
)

//go:embed templates/portal.html
var portalTemplates embed.FS

type portalShellData struct {
	Title        string
	StatePath    string
	AssetBaseURL string
	Routes       template.JS
	InitialState template.JS
}

var portalTemplate = template.Must(template.ParseFS(portalTemplates, "templates/portal.html"))

func renderPortalShell(w http.ResponseWriter, r *http.Request, routes *wayfinder.Registry, title string, state map[string]any) {
	payload, err := json.Marshal(state)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	if routes == nil {
		routes = billing.NewRouteRegistry()
	}

	manifest, err := json.Marshal(routes.Export())

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	var body bytes.Buffer

	if err := portalTemplate.Execute(&body, portalShellData{
		Title:        title,
		StatePath:    routes.URL(billing.RouteState, nil),
		AssetBaseURL: "/billing/assets/",
		Routes:       template.JS(manifest),
		InitialState: template.JS(payload),
	}); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	_ = httpx.NewResponse(w).
		Header("Content-Type", "text/html; charset=utf-8").
		Send(body.Bytes())
}
