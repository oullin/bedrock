package crm

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/inertia/validation"
	"github.com/bedrock/services/inertia-demo/api/internal/database"
)

type contactForm struct {
	OrganizationID string `json:"organization_id"`
	FirstName      string `json:"first_name" validate:"required,max=255"`
	LastName       string `json:"last_name" validate:"required,max=255"`
	Email          string `json:"email" validate:"required,email,max=255"`
	Phone          string `json:"phone" validate:"omitempty,max=255"`
}

type organizationForm struct {
	Name string `json:"name" validate:"required,max=255"`
}

func newContactForm(r *http.Request) contactForm {
	return contactForm{
		OrganizationID: strings.TrimSpace(r.FormValue("organization_id")),
		FirstName:      strings.TrimSpace(r.FormValue("first_name")),
		LastName:       strings.TrimSpace(r.FormValue("last_name")),
		Email:          strings.TrimSpace(r.FormValue("email")),
		Phone:          strings.TrimSpace(r.FormValue("phone")),
	}
}

func newContactFormFromContact(contact database.Contact) contactForm {
	form := contactForm{
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
	}

	if contact.OrganizationID != nil {
		form.OrganizationID = fmt.Sprintf("%d", *contact.OrganizationID)
	}

	return form
}

func emptyContactForm() contactForm {
	return contactForm{}
}

func (f contactForm) validate() protocol.ValidationErrors {
	errors := validation.Validate(f)

	if strings.TrimSpace(f.OrganizationID) != "" {
		if _, err := strconv.ParseInt(f.OrganizationID, 10, 64); err != nil {
			if errors == nil {
				errors = make(protocol.ValidationErrors)
			}

			errors["organization_id"] = "The organization id field must be a valid identifier."
		}
	}

	return errors
}

func (f contactForm) record() database.Contact {
	return database.Contact{
		OrganizationID: parseOrganizationID(f.OrganizationID),
		FirstName:      f.FirstName,
		LastName:       f.LastName,
		Email:          f.Email,
		Phone:          f.Phone,
	}
}

func newOrganizationForm(r *http.Request) organizationForm {
	return organizationForm{
		Name: strings.TrimSpace(r.FormValue("name")),
	}
}

func (f organizationForm) validate() protocol.ValidationErrors {
	return validation.Validate(f)
}

func parseOrganizationID(raw string) *int64 {
	raw = strings.TrimSpace(raw)

	if strings.TrimSpace(raw) == "" {
		return nil
	}

	id, err := strconv.ParseInt(raw, 10, 64)

	if err != nil {
		panic(fmt.Sprintf("parseOrganizationID: invalid input %q should have been caught by validation", raw))
	}

	return &id
}
