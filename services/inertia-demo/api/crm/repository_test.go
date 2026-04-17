package crm

import (
	"testing"

	"github.com/bedrock/services/inertia-demo/api/internal/database"
)

func TestRepositoryMethods(t *testing.T) {
	t.Parallel()

	h := newCRMHarness(t)
	repo, err := newRepository(h.db)

	if err != nil {
		t.Fatalf("newRepository() error = %v", err)
	}

	if _, err := newRepository(nil); err == nil {
		t.Fatal("newRepository(nil) error = nil, want error")
	}

	contactCount, err := repo.CountContacts()

	if err != nil {
		t.Fatalf("CountContacts() error = %v", err)
	}

	if contactCount == 0 {
		t.Fatal("CountContacts() = 0, want > 0")
	}

	orgCount, err := repo.CountOrganizations()

	if err != nil {
		t.Fatalf("CountOrganizations() error = %v", err)
	}

	if orgCount == 0 {
		t.Fatal("CountOrganizations() = 0, want > 0")
	}

	noteCount, err := repo.CountNotes()

	if err != nil {
		t.Fatalf("CountNotes() error = %v", err)
	}

	if noteCount == 0 {
		t.Fatal("CountNotes() = 0, want > 0")
	}

	recentNotes, err := repo.ListRecentNotes(3)

	if err != nil {
		t.Fatalf("ListRecentNotes() error = %v", err)
	}

	if len(recentNotes) == 0 || len(recentNotes) > 3 {
		t.Fatalf("ListRecentNotes(3) len = %d, want 1..3", len(recentNotes))
	}

	contactsPage, err := repo.ListContactsPaginated("", false, nil, "next", 5)

	if err != nil {
		t.Fatalf("ListContactsPaginated() error = %v", err)
	}

	if len(contactsPage.Data) == 0 {
		t.Fatal("ListContactsPaginated() returned empty page, want > 0")
	}

	contact, err := repo.GetContact(1)

	if err != nil {
		t.Fatalf("GetContact() error = %v", err)
	}

	if contact == nil || contact.ID != 1 {
		t.Fatalf("GetContact(1) = %#v, want contact with ID 1", contact)
	}

	contactNotes, err := repo.ListContactNotes(1)

	if err != nil {
		t.Fatalf("ListContactNotes() error = %v", err)
	}

	if contactNotes == nil {
		t.Fatal("ListContactNotes(1) returned nil, want non-nil slice")
	}

	orgs, err := repo.ListOrganizations("")

	if err != nil {
		t.Fatalf("ListOrganizations() error = %v", err)
	}

	if len(orgs) == 0 {
		t.Fatal("ListOrganizations() returned empty, want > 0")
	}

	orgsPage, err := repo.ListOrganizationsPaginated("", 1, 5)

	if err != nil {
		t.Fatalf("ListOrganizationsPaginated() error = %v", err)
	}

	if len(orgsPage.Data) == 0 || orgsPage.Total == 0 {
		t.Fatalf("ListOrganizationsPaginated() = %d items, total %d, want > 0", len(orgsPage.Data), orgsPage.Total)
	}

	org, err := repo.GetOrganization(1)

	if err != nil {
		t.Fatalf("GetOrganization() error = %v", err)
	}

	if org == nil || org.ID != 1 {
		t.Fatalf("GetOrganization(1) = %#v, want org with ID 1", org)
	}

	orgContacts, err := repo.ListContactsByOrganization(1)

	if err != nil {
		t.Fatalf("ListContactsByOrganization() error = %v", err)
	}

	if len(orgContacts) == 0 {
		t.Fatal("ListContactsByOrganization(1) returned empty, want > 0")
	}

	orgContactsPage, err := repo.ListContactsByOrgPaginated(1, nil, "next", 5)

	if err != nil {
		t.Fatalf("ListContactsByOrgPaginated() error = %v", err)
	}

	if len(orgContactsPage.Data) == 0 {
		t.Fatal("ListContactsByOrgPaginated(1) returned empty page, want > 0")
	}

	id, err := repo.CreateContact(database.Contact{
		FirstName: "Repo",
		LastName:  "Created",
		Email:     "repo.created@example.com",
		Phone:     "555-0199",
	})

	if err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}

	if id == 0 {
		t.Fatal("CreateContact() returned zero id")
	}

	if err := repo.UpdateContact(id, database.Contact{
		FirstName: "Repo",
		LastName:  "Updated",
		Email:     "repo.updated@example.com",
		Phone:     "555-0200",
	}); err != nil {
		t.Fatalf("UpdateContact() error = %v", err)
	}

	updated, err := repo.GetContact(id)

	if err != nil {
		t.Fatalf("GetContact(updated) error = %v", err)
	}

	if updated == nil || updated.LastName != "Updated" || updated.Email != "repo.updated@example.com" {
		t.Fatalf("GetContact(updated) = %#v, want LastName=Updated, Email=repo.updated@example.com", updated)
	}

	if err := repo.ToggleContactFavorite(id); err != nil {
		t.Fatalf("ToggleContactFavorite() error = %v", err)
	}

	toggled, err := repo.GetContact(id)

	if err != nil {
		t.Fatalf("GetContact(toggled) error = %v", err)
	}

	if toggled == nil || !toggled.IsFavorite {
		t.Fatalf("GetContact(toggled) = %#v, want IsFavorite=true", toggled)
	}

	noteID, err := repo.CreateNote(id, h.user.ID, "Repository note")

	if err != nil {
		t.Fatalf("CreateNote() error = %v", err)
	}

	if noteID == 0 {
		t.Fatal("CreateNote() returned zero id")
	}

	notes, err := repo.ListContactNotes(id)

	if err != nil {
		t.Fatalf("ListContactNotes(created) error = %v", err)
	}

	if len(notes) != 1 || notes[0].Body != "Repository note" {
		t.Fatalf("ListContactNotes(created) = %#v, want single note with body 'Repository note'", notes)
	}

	if err := repo.DeleteContact(id); err != nil {
		t.Fatalf("DeleteContact() error = %v", err)
	}

	deleted, err := repo.GetContact(id)

	if err != nil {
		t.Fatalf("GetContact(deleted) error = %v", err)
	}

	if deleted != nil {
		t.Fatalf("GetContact(deleted) = %#v, want nil", deleted)
	}

	if err := repo.UpdateOrganization(1, "Repository Updated Org"); err != nil {
		t.Fatalf("UpdateOrganization() error = %v", err)
	}

	updatedOrg, err := repo.GetOrganization(1)

	if err != nil {
		t.Fatalf("GetOrganization(updated) error = %v", err)
	}

	if updatedOrg == nil || updatedOrg.Name != "Repository Updated Org" {
		t.Fatalf("GetOrganization(updated) = %#v, want Name='Repository Updated Org'", updatedOrg)
	}
}
