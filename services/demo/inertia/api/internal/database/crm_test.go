package database

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestCreateUserHashesPassword(t *testing.T) {
	t.Parallel()

	db, err := Open(":memory:")

	if err != nil {
		t.Fatalf("Open(:memory:) error = %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if _, err := CreateUser(db, "Demo User", "test@example.com", "password", nil); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	user, err := FindUserByEmail(db, "test@example.com")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	if user == nil {
		t.Fatal("FindUserByEmail() returned nil user")
	}

	if user.PasswordHash == "password" {
		t.Fatal("CreateUser() stored plaintext password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password")); err != nil {
		t.Fatalf("stored password hash did not verify: %v", err)
	}
}

func TestCreateUserNormalisesEmail(t *testing.T) {
	t.Parallel()

	db, err := Open(":memory:")

	if err != nil {
		t.Fatalf("Open(:memory:) error = %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if _, err := CreateUser(db, "Demo User", "  User@Example.COM  ", "password", nil); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	user, err := FindUserByEmail(db, "USER@EXAMPLE.COM")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	if user == nil {
		t.Fatal("FindUserByEmail() returned nil for differently-cased email")
	}

	if user.Email != "user@example.com" {
		t.Fatalf("stored email = %q, want %q", user.Email, "user@example.com")
	}
}

func TestOrganizationQueriesAndPagination(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	alphaID := mustCreateOrganization(t, db, "Alpha Corp")
	bravoID := mustCreateOrganization(t, db, "Bravo LLC")
	charlieID := mustCreateOrganization(t, db, "Charlie Inc")

	_ = bravoID

	mustCreateContact(t, db, &alphaID, "Dan", "Drew", "dan@example.com", false)
	mustCreateContact(t, db, &alphaID, "Eve", "Ellis", "eve@example.com", false)
	mustCreateContact(t, db, &charlieID, "Fay", "Ford", "fay@example.com", false)

	orgs, err := ListOrganizations(db, "Alpha")

	if err != nil {
		t.Fatalf("ListOrganizations() error = %v", err)
	}

	if len(orgs) != 1 || orgs[0].ID != alphaID || orgs[0].ContactsCount != 2 {
		t.Fatalf("ListOrganizations() = %#v, want Alpha with 2 contacts", orgs)
	}

	org, err := GetOrganization(db, bravoID)

	if err != nil {
		t.Fatalf("GetOrganization() error = %v", err)
	}

	if org == nil || org.Name != "Bravo LLC" {
		t.Fatalf("GetOrganization() = %#v, want Bravo LLC", org)
	}

	if missing, err := GetOrganization(db, 9999); err != nil || missing != nil {
		t.Fatalf("GetOrganization(9999) = %#v, %v, want nil, nil", missing, err)
	}

	if err := UpdateOrganization(db, bravoID, "  Bravo Updated  "); err != nil {
		t.Fatalf("UpdateOrganization() error = %v", err)
	}

	org, err = GetOrganization(db, bravoID)

	if err != nil {
		t.Fatalf("GetOrganization() after update error = %v", err)
	}

	if org == nil || org.Name != "Bravo Updated" {
		t.Fatalf("GetOrganization() after update = %#v, want trimmed name", org)
	}

	if err := UpdateOrganization(db, 9999, "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("UpdateOrganization() error = %v, want %v", err, ErrNotFound)
	}

	if got, err := CountOrganizations(db); err != nil || got != 3 {
		t.Fatalf("CountOrganizations() = %d, %v, want 3, nil", got, err)
	}

	page, err := ListOrganizationsPaginated(db, "", 1, 2)

	if err != nil {
		t.Fatalf("ListOrganizationsPaginated() error = %v", err)
	}

	if len(page.Data) != 2 || page.Total != 3 || page.CurrentPage != 1 || page.LastPage != 2 {
		t.Fatalf("ListOrganizationsPaginated() = %#v", page)
	}

	lastPage, err := ListOrganizationsPaginated(db, "", 9, 2)

	if err != nil {
		t.Fatalf("ListOrganizationsPaginated(last page) error = %v", err)
	}

	if lastPage.CurrentPage != 2 {
		t.Fatalf("ListOrganizationsPaginated(last page).CurrentPage = %d, want 2", lastPage.CurrentPage)
	}

	if _, err := ListOrganizationsPaginated(db, "", 1, 0); err == nil {
		t.Fatal("ListOrganizationsPaginated() error = nil, want error for perPage <= 0")
	}
}

func TestContactQueriesAndPagination(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	orgID := mustCreateOrganization(t, db, "Acme")

	firstID := mustCreateContact(t, db, &orgID, "Alice", "Able", "alice@example.com", false)
	secondID := mustCreateContact(t, db, &orgID, "Bob", "Baker", "bob@example.com", true)
	thirdID := mustCreateContact(t, db, nil, "Cara", "Cole", "cara@example.com", false)

	contacts, err := ListContacts(db, "acme", false)

	if err != nil {
		t.Fatalf("ListContacts(search) error = %v", err)
	}

	if len(contacts) != 2 {
		t.Fatalf("ListContacts(search) len = %d, want 2", len(contacts))
	}

	favorites, err := ListContacts(db, "", true)

	if err != nil {
		t.Fatalf("ListContacts(favorites) error = %v", err)
	}

	if len(favorites) != 1 || favorites[0].ID != secondID {
		t.Fatalf("ListContacts(favorites) = %#v, want Bob only", favorites)
	}

	recent, err := ListRecentContacts(db, 2)

	if err != nil {
		t.Fatalf("ListRecentContacts() error = %v", err)
	}

	if len(recent) != 2 {
		t.Fatalf("ListRecentContacts() len = %d, want 2", len(recent))
	}

	contact, err := GetContact(db, thirdID)

	if err != nil {
		t.Fatalf("GetContact() error = %v", err)
	}

	if contact == nil || contact.ID != thirdID || contact.OrganizationID != nil {
		t.Fatalf("GetContact() = %#v, want contact without org", contact)
	}

	if missing, err := GetContact(db, 9999); err != nil || missing != nil {
		t.Fatalf("GetContact(9999) = %#v, %v, want nil, nil", missing, err)
	}

	if err := UpdateContact(db, firstID, Contact{
		OrganizationID: nil,
		FirstName:      "Alice",
		LastName:       "Updated",
		Email:          "alice.updated@example.com",
		Phone:          "555-0109",
	}); err != nil {
		t.Fatalf("UpdateContact() error = %v", err)
	}

	updated, err := GetContact(db, firstID)

	if err != nil {
		t.Fatalf("GetContact(updated) error = %v", err)
	}

	if updated == nil || updated.LastName != "Updated" || updated.OrganizationID != nil {
		t.Fatalf("GetContact(updated) = %#v, want updated contact", updated)
	}

	if err := UpdateContact(db, 9999, Contact{}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("UpdateContact(missing) error = %v, want %v", err, ErrNotFound)
	}

	if err := ToggleContactFavorite(db, firstID); err != nil {
		t.Fatalf("ToggleContactFavorite() error = %v", err)
	}

	toggled, err := GetContact(db, firstID)

	if err != nil {
		t.Fatalf("GetContact(toggled) error = %v", err)
	}

	if toggled == nil || !toggled.IsFavorite {
		t.Fatalf("GetContact(toggled) = %#v, want favorite=true", toggled)
	}

	if err := ToggleContactFavorite(db, 9999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("ToggleContactFavorite(missing) error = %v, want %v", err, ErrNotFound)
	}

	byOrg, err := ListContactsByOrganization(db, orgID)

	if err != nil {
		t.Fatalf("ListContactsByOrganization() error = %v", err)
	}

	if len(byOrg) != 1 || byOrg[0].ID != secondID {
		t.Fatalf("ListContactsByOrganization() = %#v, want Bob only after org removal", byOrg)
	}

	page, err := ListContactsPaginated(db, "", false, nil, "next", 2)

	if err != nil {
		t.Fatalf("ListContactsPaginated() error = %v", err)
	}

	if len(page.Data) != 2 || page.NextCursor == nil || page.PrevCursor != nil {
		t.Fatalf("ListContactsPaginated() = %#v, want next cursor only", page)
	}

	nextPage, err := ListContactsPaginated(db, "", false, page.NextCursor, "next", 2)

	if err != nil {
		t.Fatalf("ListContactsPaginated(next page) error = %v", err)
	}

	if len(nextPage.Data) != 1 || nextPage.PrevCursor == nil {
		t.Fatalf("ListContactsPaginated(next page) = %#v, want remaining page with prev cursor", nextPage)
	}

	prevPage, err := ListContactsPaginated(db, "", false, nextPage.PrevCursor, "prev", 1)

	if err != nil {
		t.Fatalf("ListContactsPaginated(prev page) error = %v", err)
	}

	if len(prevPage.Data) != 1 {
		t.Fatalf("ListContactsPaginated(prev page) = %#v, want one row", prevPage)
	}

	orgPage, err := ListContactsByOrgPaginated(db, orgID, nil, "next", 1)

	if err != nil {
		t.Fatalf("ListContactsByOrgPaginated() error = %v", err)
	}

	if len(orgPage.Data) != 1 || orgPage.Data[0].ID != secondID {
		t.Fatalf("ListContactsByOrgPaginated() = %#v, want Bob page", orgPage)
	}

	if got, err := CountContacts(db); err != nil || got != 3 {
		t.Fatalf("CountContacts() = %d, %v, want 3, nil", got, err)
	}

	if err := DeleteContact(db, thirdID); err != nil {
		t.Fatalf("DeleteContact() error = %v", err)
	}

	if missing, err := GetContact(db, thirdID); err != nil || missing != nil {
		t.Fatalf("GetContact(deleted) = %#v, %v, want nil, nil", missing, err)
	}

	if err := DeleteContact(db, 9999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("DeleteContact(missing) error = %v, want %v", err, ErrNotFound)
	}

	if _, err := ListContactsPaginated(db, "", false, nil, "next", 0); err == nil {
		t.Fatal("ListContactsPaginated() error = nil, want error for perPage <= 0")
	}

	invalidCursor := "bad"

	if _, err := ListContactsPaginated(db, "", false, &invalidCursor, "next", 2); err == nil {
		t.Fatal("ListContactsPaginated() error = nil, want invalid cursor error")
	}

	if _, err := ListContactsByOrgPaginated(db, orgID, &invalidCursor, "next", 1); err == nil {
		t.Fatal("ListContactsByOrgPaginated() error = nil, want invalid cursor error")
	}
}

func TestUserAndNoteQueries(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	verifiedAt := time.Now().Add(-time.Hour)
	userID, err := CreateUser(db, "Dana Demo", "dana@example.com", "password", &verifiedAt)

	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	user, err := FindUserByID(db, userID)

	if err != nil {
		t.Fatalf("FindUserByID() error = %v", err)
	}

	if user == nil || user.VerifiedAt == nil {
		t.Fatalf("FindUserByID() = %#v, want verified user", user)
	}

	userByEmail, err := FindUserByEmail(db, " DANA@example.com ")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	if userByEmail == nil || userByEmail.ID != userID {
		t.Fatalf("FindUserByEmail() = %#v, want same user", userByEmail)
	}

	if missing, err := FindUserByID(db, 9999); err != nil || missing != nil {
		t.Fatalf("FindUserByID(9999) = %#v, %v, want nil, nil", missing, err)
	}

	orgID := mustCreateOrganization(t, db, "Acme")
	contactID := mustCreateContact(t, db, &orgID, "Nina", "North", "nina@example.com", false)

	firstNoteAt := time.Date(2024, time.January, 2, 10, 0, 0, 0, time.UTC)
	secondNoteAt := time.Date(2024, time.January, 2, 11, 0, 0, 0, time.UTC)

	if _, err := CreateNoteAt(db, contactID, userID, "First note", firstNoteAt); err != nil {
		t.Fatalf("CreateNoteAt() error = %v", err)
	}

	noteID, err := CreateNoteAt(db, contactID, userID, "Second note", secondNoteAt)

	if err != nil {
		t.Fatalf("CreateNoteAt(second) error = %v", err)
	}

	if noteID == 0 {
		t.Fatal("CreateNoteAt(second) returned zero id")
	}

	thirdNoteID, err := CreateNote(db, contactID, userID, "Third note")

	if err != nil {
		t.Fatalf("CreateNote() error = %v", err)
	}

	if thirdNoteID == 0 {
		t.Fatal("CreateNote() returned zero id")
	}

	notes, err := ListContactNotes(db, contactID)

	if err != nil {
		t.Fatalf("ListContactNotes() error = %v", err)
	}

	if len(notes) != 3 || notes[1].Body != "Second note" || notes[2].Body != "First note" {
		t.Fatalf("ListContactNotes() = %#v, want explicit notes in descending order after newest current note", notes)
	}

	recent, err := ListRecentNotes(db, 2)

	if err != nil {
		t.Fatalf("ListRecentNotes() error = %v", err)
	}

	if len(recent) != 2 || recent[1].Body != "Second note" {
		t.Fatalf("ListRecentNotes() = %#v, want explicit note ordering after newest current note", recent)
	}

	if got, err := CountNotes(db); err != nil || got != 3 {
		t.Fatalf("CountNotes() = %d, %v, want 3, nil", got, err)
	}

	if err := DeleteContact(db, contactID); err != nil {
		t.Fatalf("DeleteContact() error = %v", err)
	}

	if got, err := CountNotes(db); err != nil || got != 0 {
		t.Fatalf("CountNotes() after delete = %d, %v, want 0, nil", got, err)
	}
}
