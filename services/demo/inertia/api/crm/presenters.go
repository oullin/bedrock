package crm

import (
	"fmt"
	"time"

	"github.com/bedrock/services/demo/inertia/api/internal/database"
)

func contactsProps(contacts []database.Contact) []map[string]any {
	items := make([]map[string]any, 0, len(contacts))

	for _, contact := range contacts {
		items = append(items, contactProp(contact))
	}

	return items
}

func contactProp(contact database.Contact) map[string]any {
	item := map[string]any{
		"id":          contact.ID,
		"first_name":  contact.FirstName,
		"last_name":   contact.LastName,
		"email":       contact.Email,
		"phone":       contact.Phone,
		"is_favorite": contact.IsFavorite,
		"created_at":  contact.CreatedAt.Format(time.RFC3339),
		"updated_at":  contact.UpdatedAt.Format(time.RFC3339),
	}

	if contact.OrganizationID != nil {
		item["organization"] = map[string]any{
			"id":   *contact.OrganizationID,
			"name": contact.OrganizationName,
		}
	} else {
		item["organization"] = nil
	}

	return item
}

func contactPropValue(contact *database.Contact) map[string]any {
	if contact == nil {
		return nil
	}

	return contactProp(*contact)
}

func contactFormProps(form contactForm) map[string]any {
	return map[string]any{
		"organization_id": form.OrganizationID,
		"first_name":      form.FirstName,
		"last_name":       form.LastName,
		"email":           form.Email,
		"phone":           form.Phone,
	}
}

func organizationsProps(orgs []database.Organization) []map[string]any {
	items := make([]map[string]any, 0, len(orgs))

	for _, org := range orgs {
		items = append(items, organizationProp(org))
	}

	return items
}

func organizationProp(org database.Organization) map[string]any {
	return map[string]any{
		"id":             org.ID,
		"name":           org.Name,
		"contacts_count": org.ContactsCount,
	}
}

func organizationOptions(orgs []database.Organization) []map[string]any {
	options := []map[string]any{{"value": "", "label": "No organization"}}

	for _, org := range orgs {
		options = append(options, map[string]any{
			"value": fmt.Sprintf("%d", org.ID),
			"label": org.Name,
		})
	}

	return options
}

func notesProps(notes []database.Note) []map[string]any {
	items := make([]map[string]any, 0, len(notes))

	for _, note := range notes {
		items = append(items, map[string]any{
			"id":         note.ID,
			"body":       note.Body,
			"created_at": note.CreatedAt.Format(time.RFC3339),
			"user": map[string]any{
				"id":   note.UserID,
				"name": note.UserName,
			},
			"contact": map[string]any{
				"id":   note.ContactID,
				"name": note.ContactName,
			},
		})
	}

	return items
}

func recentActivityProps(notes []database.Note) []map[string]any {
	return notesProps(notes)
}

func cursorContactsProps(page database.CursorPage[database.Contact]) map[string]any {
	result := map[string]any{
		"data":        contactsProps(page.Data),
		"next_cursor": page.NextCursor,
		"prev_cursor": page.PrevCursor,
	}

	return result
}

func offsetOrganizationsProps(page database.OffsetPage[database.Organization]) map[string]any {
	return map[string]any{
		"data":         organizationsProps(page.Data),
		"total":        page.Total,
		"per_page":     page.PerPage,
		"current_page": page.CurrentPage,
		"last_page":    page.LastPage,
	}
}
