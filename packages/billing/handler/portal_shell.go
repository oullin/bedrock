package handler

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
)

//go:embed templates/portal.html
var portalTemplates embed.FS

type portalShellData struct {
	Title        string
	StatePath    string
	AssetBaseURL string
	InitialState template.JS
}

var portalTemplate = template.Must(template.ParseFS(portalTemplates, "templates/portal.html"))

func renderPortalShell(w http.ResponseWriter, r *http.Request, title string, state map[string]any) {
	payload, err := json.Marshal(state)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	portalTemplate.Execute(w, portalShellData{
		Title:        title,
		StatePath:    "/billing/state",
		AssetBaseURL: "/billing/assets/",
		InitialState: template.JS(payload),
	})
}
