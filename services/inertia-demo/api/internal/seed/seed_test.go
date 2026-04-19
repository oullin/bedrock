package seed

import (
	"testing"

	"github.com/bedrock/services/inertia-demo/api/internal/database"
)

func TestRunSeedsDatabase(t *testing.T) {
	t.Parallel()

	db, err := database.Open(":memory:")

	if err != nil {
		t.Fatalf("Open(:memory:) error = %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if err := Run(db); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := database.GetCounter(db, "priority_escalations"); got != 18 {
		t.Fatalf("GetCounter(priority_escalations) = %d, want 18", got)
	}

	if got := database.InviteCount(db); got != 20 {
		t.Fatalf("InviteCount() = %d, want 20", got)
	}

	if got := database.UploadCount(db); got != 15 {
		t.Fatalf("UploadCount() = %d, want 15", got)
	}

	if got := database.ApprovalCount(db); got != 15 {
		t.Fatalf("ApprovalCount() = %d, want 15", got)
	}

	users := []struct {
		email string
		name  string
	}{
		{email: "test@example.com", name: "Demo User"},
		{email: "alice@example.com", name: "Alice Manager"},
	}

	for _, want := range users {
		user, err := database.FindUserByEmail(db, want.email)

		if err != nil {
			t.Fatalf("FindUserByEmail(%q) error = %v", want.email, err)
		}

		if user == nil || user.Name != want.name {
			t.Fatalf("FindUserByEmail(%q) = %#v, want name %q", want.email, user, want.name)
		}
	}

	contacts, err := database.ListContacts(db, "", false)

	if err != nil {
		t.Fatalf("ListContacts() error = %v", err)
	}

	if len(contacts) == 0 {
		t.Fatal("ListContacts() returned no contacts")
	}

	organizations, err := database.ListOrganizations(db, "")

	if err != nil {
		t.Fatalf("ListOrganizations() error = %v", err)
	}

	if len(organizations) == 0 {
		t.Fatal("ListOrganizations() returned no organizations")
	}

	notes, err := database.ListRecentNotes(db, 5)

	if err != nil {
		t.Fatalf("ListRecentNotes() error = %v", err)
	}

	if len(notes) == 0 {
		t.Fatal("ListRecentNotes() returned no notes")
	}

	if invites := database.ListInvites(db); len(invites) != 20 {
		t.Fatalf("ListInvites() len = %d, want 20", len(invites))
	}

	if uploads := database.ListUploads(db); len(uploads) != 15 {
		t.Fatalf("ListUploads() len = %d, want 15", len(uploads))
	}

	if approvals := database.ListApprovals(db); len(approvals) != 15 {
		t.Fatalf("ListApprovals() len = %d, want 15", len(approvals))
	}
}

func TestRunTruncatesExistingData(t *testing.T) {
	t.Parallel()

	db, err := database.Open(":memory:")

	if err != nil {
		t.Fatalf("Open(:memory:) error = %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if _, err := database.CreateUser(db, "Extra User", "extra@example.com", "password", nil); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if err := database.CreateInvite(db, "legacy", "Legacy User", "legacy@example.com", "Owner", "Queued"); err != nil {
		t.Fatalf("CreateInvite() error = %v", err)
	}

	if err := Run(db); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	user, err := database.FindUserByEmail(db, "extra@example.com")

	if err != nil {
		t.Fatalf("FindUserByEmail(extra@example.com) error = %v", err)
	}

	if user != nil {
		t.Fatalf("FindUserByEmail(extra@example.com) = %#v, want nil after truncate", user)
	}

	if invites := database.ListInvites(db); len(invites) != 20 {
		t.Fatalf("ListInvites() len = %d, want 20 after reseed", len(invites))
	}
}
