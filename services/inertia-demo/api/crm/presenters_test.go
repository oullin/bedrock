package crm

import (
	"testing"
	"time"

	"github.com/bedrock/services/inertia-demo/api/internal/database"
)

func TestContactProp(t *testing.T) {
	t.Parallel()

	orgID := int64(3)
	now := time.Date(2026, 4, 3, 10, 0, 0, 0, time.UTC)
	contact := database.Contact{
		ID:               7,
		OrganizationID:   &orgID,
		OrganizationName: "Acme",
		FirstName:        "Mina",
		LastName:         "Cole",
		Email:            "mina@example.test",
		Phone:            "+1 555 0107",
		IsFavorite:       true,
		CreatedAt:        now,
		UpdatedAt:        now.Add(time.Hour),
	}

	got := contactProp(contact)

	if got["first_name"] != "Mina" {
		t.Fatalf("first_name = %v", got["first_name"])
	}

	org, ok := got["organization"].(map[string]any)

	if !ok {
		t.Fatalf("organization = %T, want map[string]any", got["organization"])
	}

	if org["name"] != "Acme" {
		t.Fatalf("organization.name = %v", org["name"])
	}
}

func TestOrganizationOptions(t *testing.T) {
	t.Parallel()

	got := organizationOptions([]database.Organization{
		{ID: 1, Name: "Acme"},
		{ID: 2, Name: "Globex"},
	})

	if len(got) != 3 {
		t.Fatalf("len(options) = %d, want 3", len(got))
	}

	if got[0]["label"] != "No organization" {
		t.Fatalf("default label = %v", got[0]["label"])
	}

	if got[2]["value"] != "2" {
		t.Fatalf("last value = %v, want 2", got[2]["value"])
	}
}

func TestNotesProps(t *testing.T) {
	t.Parallel()

	got := notesProps([]database.Note{{
		ID:          4,
		ContactID:   8,
		ContactName: "Mina Cole",
		UserID:      2,
		UserName:    "Avery",
		Body:        "Need follow-up",
		CreatedAt:   time.Date(2026, 4, 3, 10, 0, 0, 0, time.UTC),
	}})

	if len(got) != 1 {
		t.Fatalf("len(notes) = %d, want 1", len(got))
	}

	if got[0]["body"] != "Need follow-up" {
		t.Fatalf("body = %v", got[0]["body"])
	}
}
