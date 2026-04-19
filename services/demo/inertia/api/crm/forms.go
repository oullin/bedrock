package crm

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/validation"
	"github.com/bedrock/services/demo/inertia/api/internal/database"
)

type contactForm struct {
	OrganizationID string `json:"organization_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
}

type organizationForm struct {
	Name string `json:"name"`
}

var validatorFactory = validation.NewFactory()

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
	errors := runValidation(
		map[string]any{
			"first_name": f.FirstName,
			"last_name":  f.LastName,
			"email":      f.Email,
			"phone":      f.Phone,
		},
		map[string]any{
			"first_name": "required|max:255",
			"last_name":  "required|max:255",
			"email":      "required|email|max:255",
			"phone":      "max:255",
		},
	)

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
	return runValidation(
		map[string]any{"name": f.Name},
		map[string]any{"name": "required|max:255"},
	)
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

// runValidation runs the bedrock validator and projects the MessageBag into
// the field-keyed string map Inertia sends to the client. Only the first
// message per field is kept, matching the precognition error contract.
func runValidation(data, rules map[string]any) protocol.ValidationErrors {
	err := validatorFactory.Make(data, rules, nil, nil).Validate()

	if err == nil {
		return nil
	}

	ve, ok := err.(*validation.ValidationException)

	if !ok {
		return protocol.ValidationErrors{"_error": err.Error()}
	}

	bag := ve.Bag.ToMap()

	if len(bag) == 0 {
		return nil
	}

	out := make(protocol.ValidationErrors, len(bag))

	for key, msgs := range bag {
		if len(msgs) > 0 {
			out[key] = msgs[0]
		}
	}

	return out
}
