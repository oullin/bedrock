package features

import (
	"net/http"

	"github.com/bedrock/packages/inertia/protocol"
)

func (a app) persistentLayoutsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Layouts/PersistentLayouts", protocol.Props{})
}

func (a app) persistentLayoutsPage2Handler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Layouts/PersistentLayoutsPageTwo", protocol.Props{})
}

func (a app) nestedLayoutsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Layouts/NestedLayouts", protocol.Props{
		"title": "Nested Layouts",
		"breadcrumbs": []map[string]string{
			{"title": "Features"},
			{"title": "Layouts"},
			{"title": "Nested Layouts"},
		},
	})
}

func (a app) headHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Layouts/Head", protocol.Props{})
}

func (a app) layoutPropsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Layouts/LayoutProps", protocol.Props{})
}
