package crm

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/bedrock/packages/inertia"
	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/props"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/services/inertia-demo/api/internal/database"
)

const contactsPerPage = 15

func (a app) contactsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.listContactsHandler(w, r)
	case http.MethodPost:
		a.storeContactHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) contactsCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	orgs, err := a.repo.ListOrganizations("")

	if err != nil {
		slog.Error("list organizations", "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	form := emptyContactForm()

	if orgID := r.URL.Query().Get("organization_id"); strings.TrimSpace(orgID) != "" {
		form.OrganizationID = orgID
	}

	a.container.Render(w, r, "Contacts/Create", protocol.Props{
		"form":          contactFormProps(form),
		"organizations": organizationOptions(orgs),
	})
}

func (a app) contactByIDHandler(w http.ResponseWriter, r *http.Request) {
	contactID, action, ok := routeIDAndAction(r.URL.Path, "/contacts/")

	if !ok {
		http.NotFound(w, r)

		return
	}

	if strings.TrimSpace(action) == "" {
		switch r.Method {
		case http.MethodGet:
			a.showContactHandler(w, r, contactID)
		case http.MethodPost, http.MethodPut:
			a.updateContactHandler(w, r, contactID)
		case http.MethodDelete:
			a.deleteContactHandler(w, r, contactID)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		return
	}

	switch action {
	case "edit":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		a.editContactHandler(w, r, contactID)
	case "favorite":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		a.toggleFavoriteHandler(w, r, contactID)
	case "notes":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		a.storeNoteHandler(w, r, contactID)
	default:
		http.NotFound(w, r)
	}
}

func (a app) listContactsHandler(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	favoriteOnly := r.URL.Query().Get("favorite") == "true"

	var cursor *string

	if c := r.URL.Query().Get("cursor"); strings.TrimSpace(c) != "" {
		cursor = &c
	}

	direction := r.URL.Query().Get("direction")

	page, err := a.repo.ListContactsPaginated(search, favoriteOnly, cursor, direction, contactsPerPage)

	if err != nil {
		slog.Error("list contacts", "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	a.container.Render(w, r, "Contacts/Index", protocol.Props{
		"filters": map[string]any{
			"search":   search,
			"favorite": favoriteOnly,
		},
		"contacts": cursorContactsProps(page),
	})
}

func (a app) showContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	contact, err := a.repo.GetContact(contactID)

	if err != nil {
		slog.Error("get contact", "id", contactID, "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if contact == nil {
		http.NotFound(w, r)

		return
	}

	a.container.Render(w, r, "Contacts/Show", protocol.Props{
		"contact": contactProp(*contact),
		"notes": props.Defer(func() any {
			notes, err := a.repo.ListContactNotes(contactID)

			if err != nil {
				slog.Error("list contact notes", "error", err)
			}

			return notesProps(notes)
		}),
	})
}

func (a app) deleteContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := a.repo.DeleteContact(contactID); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.NotFound(w, r)

			return
		}

		slog.Error("delete contact", "id", contactID, "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if err := a.container.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Contact deleted",
		Message: "The contact has been removed.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.container.Redirect(w, r, a.container.RouteURL("contacts.index", nil))
}

func (a app) editContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	contact, err := a.repo.GetContact(contactID)

	if err != nil {
		slog.Error("get contact", "id", contactID, "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if contact == nil {
		http.NotFound(w, r)

		return
	}

	orgs, err := a.repo.ListOrganizations("")

	if err != nil {
		slog.Error("list organizations", "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	a.container.Render(w, r, "Contacts/Edit", protocol.Props{
		"contact":       contactProp(*contact),
		"form":          contactFormProps(newContactFormFromContact(*contact)),
		"organizations": organizationOptions(orgs),
	})
}

func (a app) storeContactHandler(w http.ResponseWriter, r *http.Request) {
	if err := protocol.ParseForm(r); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	form := newContactForm(r)
	errors := form.validate()

	if len(errors) > 0 {
		orgs, err := a.repo.ListOrganizations("")

		if err != nil {
			slog.Error("list organizations", "error", err)

			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		}

		ctx := inertia.SetValidationErrors(r.Context(), errors)

		a.container.Render(w, r.WithContext(ctx), "Contacts/Create", protocol.Props{
			"form":          contactFormProps(form),
			"organizations": organizationOptions(orgs),
		})

		return
	}

	id, err := a.repo.CreateContact(form.record())

	if err != nil {
		slog.Error("create contact", "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if err := a.container.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Contact created",
		Message: "The CRM record is ready for follow-up.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.container.Redirect(w, r, a.container.RouteURL("contacts.show", map[string]string{"contact": strconv.FormatInt(id, 10)}))
}

func (a app) updateContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := protocol.ParseForm(r); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	form := newContactForm(r)
	validationErrors := form.validate()

	if len(validationErrors) > 0 {
		existing, err := a.repo.GetContact(contactID)

		if err != nil {
			slog.Error("get contact", "id", contactID, "error", err)

			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		}

		if existing == nil {
			http.NotFound(w, r)

			return
		}

		orgs, err := a.repo.ListOrganizations("")

		if err != nil {
			slog.Error("list organizations", "error", err)

			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		}

		ctx := inertia.SetValidationErrors(r.Context(), validationErrors)

		a.container.Render(w, r.WithContext(ctx), "Contacts/Edit", protocol.Props{
			"contact":       contactPropValue(existing),
			"form":          contactFormProps(form),
			"organizations": organizationOptions(orgs),
		})

		return
	}

	if err := a.repo.UpdateContact(contactID, form.record()); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.NotFound(w, r)

			return
		}

		slog.Error("update contact", "id", contactID, "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if err := a.container.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Contact updated",
		Message: "The CRM record now reflects the latest details.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.container.Redirect(w, r, a.container.RouteURL("contacts.show", map[string]string{"contact": strconv.FormatInt(contactID, 10)}))
}

func (a app) toggleFavoriteHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := a.repo.ToggleContactFavorite(contactID); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.NotFound(w, r)

			return
		}

		slog.Error("toggle favorite", "id", contactID, "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if err := a.container.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Favorite updated",
		Message: "The contact pin state changed successfully.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.container.Redirect(w, r, a.container.RouteURL("contacts.show", map[string]string{"contact": strconv.FormatInt(contactID, 10)}))
}

func (a app) storeNoteHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := protocol.ParseForm(r); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	body := strings.TrimSpace(r.FormValue("body"))

	if strings.TrimSpace(body) == "" {
		ctx := inertia.SetValidationErrors(r.Context(), protocol.ValidationErrors{
			"body": "Add note content before saving.",
		})

		a.showContactHandler(w, r.WithContext(ctx), contactID)

		return
	}

	user := a.container.CurrentUser(r)

	if user == nil {
		a.container.Redirect(w, r, a.container.RouteURL("login", nil))

		return
	}

	if _, err := a.repo.CreateNote(contactID, user.ID, body); err != nil {
		slog.Error("create note", "contact_id", contactID, "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if err := a.container.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Note added",
		Message: "Recent activity has been updated.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.container.Redirect(w, r, a.container.RouteURL("contacts.show", map[string]string{"contact": strconv.FormatInt(contactID, 10)}))
}

func routeIDAndAction(path, prefix string) (int64, string, bool) {
	trimmed := strings.TrimPrefix(path, prefix)
	parts := strings.Split(strings.Trim(trimmed, "/"), "/")

	if len(parts) == 0 || strings.TrimSpace(parts[0]) == "" {
		return 0, "", false
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)

	if err != nil {
		return 0, "", false
	}

	if len(parts) == 1 {
		return id, "", true
	}

	if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
		return id, parts[1], true
	}

	return 0, "", false
}
