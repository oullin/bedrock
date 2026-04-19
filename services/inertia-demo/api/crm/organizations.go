package crm

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/bedrock/packages/inertia"
	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
)

const organizationsPerPage = 20

func (a app) organizationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	pageNum := 1

	if p := r.URL.Query().Get("page"); strings.TrimSpace(p) != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			pageNum = n
		}
	}

	page, err := a.repo.ListOrganizationsPaginated(search, pageNum, organizationsPerPage)

	if err != nil {
		slog.Error("list organizations", "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	a.container.Render(w, r, "Organizations/Index", protocol.Props{
		"filters": map[string]any{
			"search": search,
		},
		"organizations": offsetOrganizationsProps(page),
	})
}

func (a app) organizationByIDHandler(w http.ResponseWriter, r *http.Request) {
	organizationID, _, ok := routeIDAndAction(r.URL.Path, "/organizations/")

	if !ok {
		http.NotFound(w, r)

		return
	}

	switch r.Method {
	case http.MethodGet:
		a.showOrganizationHandler(w, r, organizationID)
	case http.MethodPost, http.MethodPut:
		a.updateOrganizationHandler(w, r, organizationID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) showOrganizationHandler(w http.ResponseWriter, r *http.Request, organizationID int64) {
	org, err := a.repo.GetOrganization(organizationID)

	if err != nil {
		slog.Error("get organization", "id", organizationID, "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if org == nil {
		http.NotFound(w, r)

		return
	}

	var cursor *string

	if c := r.URL.Query().Get("cursor"); strings.TrimSpace(c) != "" {
		cursor = &c
	}

	direction := r.URL.Query().Get("direction")

	contactsPage, err := a.repo.ListContactsByOrgPaginated(organizationID, cursor, direction, contactsPerPage)

	if err != nil {
		slog.Error("list contacts by org", "id", organizationID, "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	a.container.Render(w, r, "Organizations/Show", protocol.Props{
		"organization": organizationProp(*org),
		"contacts":     cursorContactsProps(contactsPage),
	})
}

func (a app) updateOrganizationHandler(w http.ResponseWriter, r *http.Request, organizationID int64) {
	if err := protocol.ParseForm(r); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	form := newOrganizationForm(r)
	errors := form.validate()

	if len(errors) > 0 {
		ctx := inertia.SetValidationErrors(r.Context(), errors)

		a.showOrganizationHandler(w, r.WithContext(ctx), organizationID)

		return
	}

	if err := a.repo.UpdateOrganization(organizationID, form.Name); err != nil {
		slog.Error("update organization", "id", organizationID, "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if err := a.container.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Organization updated",
		Message: "The company record was saved.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.container.Redirect(w, r, a.container.RouteURL("organizations.show", map[string]string{"organization": strconv.FormatInt(organizationID, 10)}))
}
