package database

import (
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"
)

type fakeResult struct {
	rows int64
	err  error
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := Open(":memory:")

	if err != nil {
		t.Fatalf("Open(:memory:) error = %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func mustCreateUser(t *testing.T, db *sql.DB, name, email string) int64 {
	t.Helper()

	id, err := CreateUser(db, name, email, "password", nil)

	if err != nil {
		t.Fatalf("CreateUser(%q) error = %v", email, err)
	}

	return id
}

func mustCreateOrganization(t *testing.T, db *sql.DB, name string) int64 {
	t.Helper()

	id, err := CreateOrganization(db, name)

	if err != nil {
		t.Fatalf("CreateOrganization(%q) error = %v", name, err)
	}

	return id
}

func mustCreateContact(t *testing.T, db *sql.DB, orgID *int64, firstName, lastName, email string, favorite bool) int64 {
	t.Helper()

	id, err := CreateContact(db, Contact{
		OrganizationID: orgID,
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		Phone:          "+1 555 0100",
		IsFavorite:     favorite,
	})

	if err != nil {
		t.Fatalf("CreateContact(%q) error = %v", email, err)
	}

	return id
}

func (r fakeResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (r fakeResult) RowsAffected() (int64, error) {
	return r.rows, r.err
}

func TestOpenAndTruncate(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	orgID := mustCreateOrganization(t, db, "Acme")

	if _, err := CreateContact(db, Contact{
		OrganizationID: &orgID,
		FirstName:      "Ana",
		LastName:       "Able",
		Email:          "ana@example.com",
		Phone:          "555-0101",
	}); err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}

	if err := CreateInvite(db, "invite_1", "Invited", "invite@example.com", "Manager", "Queued"); err != nil {
		t.Fatalf("CreateInvite() error = %v", err)
	}

	if err := CreateUpload(db, "upload_1", "Demo upload", "demo.csv", "10 KB", "Ready"); err != nil {
		t.Fatalf("CreateUpload() error = %v", err)
	}

	if err := CreateApproval(db, "approval_1", "Demo approval", "Queued"); err != nil {
		t.Fatalf("CreateApproval() error = %v", err)
	}

	if err := SetCounter(db, "counter", 4); err != nil {
		t.Fatalf("SetCounter() error = %v", err)
	}

	if err := Truncate(db); err != nil {
		t.Fatalf("Truncate() error = %v", err)
	}

	if got, err := CountOrganizations(db); err != nil || got != 0 {
		t.Fatalf("CountOrganizations() = %d, %v, want 0, nil", got, err)
	}

	if got, err := CountContacts(db); err != nil || got != 0 {
		t.Fatalf("CountContacts() = %d, %v, want 0, nil", got, err)
	}

	if got := InviteCount(db); got != 0 {
		t.Fatalf("InviteCount() = %d, want 0", got)
	}

	if got := UploadCount(db); got != 0 {
		t.Fatalf("UploadCount() = %d, want 0", got)
	}

	if got := ApprovalCount(db); got != 0 {
		t.Fatalf("ApprovalCount() = %d, want 0", got)
	}

	if got := GetCounter(db, "counter"); got != 0 {
		t.Fatalf("GetCounter(counter) = %d, want 0", got)
	}

	id, err := CreateOrganization(db, "Reset Sequence")

	if err != nil {
		t.Fatalf("CreateOrganization() after truncate error = %v", err)
	}

	if id != 1 {
		t.Fatalf("CreateOrganization() after truncate id = %d, want 1", id)
	}
}

func TestInviteUploadApprovalAndCounterQueries(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	now := time.Now()

	if err := CreateInviteAt(db, "invite_2", "Pat Invite", "pat@example.com", "Owner", "Accepted", now.Add(-2*time.Hour)); err != nil {
		t.Fatalf("CreateInviteAt() error = %v", err)
	}

	if err := CreateInvite(db, "invite_3", "Robin Invite", "robin@example.com", "Analyst", "Queued"); err != nil {
		t.Fatalf("CreateInvite() error = %v", err)
	}

	if invites := ListInvites(db); len(invites) != 2 {
		t.Fatalf("ListInvites() len = %d, want 2", len(invites))
	}

	if got := InviteCount(db); got != 2 {
		t.Fatalf("InviteCount() = %d, want 2", got)
	}

	if err := CreateUploadAt(db, "upload_2", "Import", "import.csv", "24 KB", "Processed", now.Add(-3*time.Hour)); err != nil {
		t.Fatalf("CreateUploadAt() error = %v", err)
	}

	if err := CreateUpload(db, "upload_3", "Export", "export.csv", "40 KB", "Queued"); err != nil {
		t.Fatalf("CreateUpload() error = %v", err)
	}

	uploads := ListUploads(db)

	if len(uploads) != 2 {
		t.Fatalf("ListUploads() len = %d, want 2", len(uploads))
	}

	if got := UploadCount(db); got != 2 {
		t.Fatalf("UploadCount() = %d, want 2", got)
	}

	if err := CreateApprovalAt(db, "approval_2", "Approve launch", "Synced", now.Add(-time.Hour)); err != nil {
		t.Fatalf("CreateApprovalAt() error = %v", err)
	}

	if err := CreateApproval(db, "approval_3", "Approve pricing", "Queued"); err != nil {
		t.Fatalf("CreateApproval() error = %v", err)
	}

	if approvals := ListApprovals(db); len(approvals) != 2 {
		t.Fatalf("ListApprovals() len = %d, want 2", len(approvals))
	}

	if got := ApprovalCount(db); got != 2 {
		t.Fatalf("ApprovalCount() = %d, want 2", got)
	}

	if got := GetCounter(db, "jobs"); got != 0 {
		t.Fatalf("GetCounter(jobs) = %d, want 0", got)
	}

	if err := IncrementCounter(db, "jobs"); err != nil {
		t.Fatalf("IncrementCounter() error = %v", err)
	}

	if err := IncrementCounter(db, "jobs"); err != nil {
		t.Fatalf("IncrementCounter() second call error = %v", err)
	}

	if got := GetCounter(db, "jobs"); got != 2 {
		t.Fatalf("GetCounter(jobs) = %d, want 2", got)
	}

	if err := SetCounter(db, "jobs", 7); err != nil {
		t.Fatalf("SetCounter() error = %v", err)
	}

	if got := GetCounter(db, "jobs"); got != 7 {
		t.Fatalf("GetCounter(jobs) after SetCounter = %d, want 7", got)
	}
}

func TestTimeAgo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		when time.Time
		want string
	}{
		{name: "just now", when: time.Now().Add(-10 * time.Second), want: "Just now"},
		{name: "minutes", when: time.Now().Add(-5 * time.Minute), want: "5m ago"},
		{name: "hours", when: time.Now().Add(-2 * time.Hour), want: "2h ago"},
		{name: "one day", when: time.Now().Add(-24 * time.Hour), want: "1d ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeAgo(tt.when); got != tt.want {
				t.Fatalf("TimeAgo() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestScanRowsAndHelpers(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)

	rows, err := db.Query(`SELECT value FROM (SELECT 1 AS value UNION ALL SELECT 'bad' AS value UNION ALL SELECT 3 AS value)`)

	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}

	t.Cleanup(func() {
		rows.Close()
	})

	values := scanRows(rows, func(scan func(...any) error) (int, error) {
		var value int

		if err := scan(&value); err != nil {
			return 0, err
		}

		return value, nil
	})

	if len(values) != 2 || values[0] != 1 || values[1] != 3 {
		t.Fatalf("scanRows() = %#v, want [1 3]", values)
	}

	var id int64 = 42

	if got := nullableInt64(&id); got != int64(42) {
		t.Fatalf("nullableInt64(&id) = %v, want 42", got)
	}

	if got := nullableInt64(nil); got != nil {
		t.Fatalf("nullableInt64(nil) = %v, want nil", got)
	}

	if got := boolToInt(true); got != 1 {
		t.Fatalf("boolToInt(true) = %d, want 1", got)
	}

	if got := boolToInt(false); got != 0 {
		t.Fatalf("boolToInt(false) = %d, want 0", got)
	}
}

func TestCheckRowsAffected(t *testing.T) {
	t.Parallel()

	if err := checkRowsAffected(fakeResult{rows: 1}); err != nil {
		t.Fatalf("checkRowsAffected() error = %v, want nil", err)
	}

	if err := checkRowsAffected(fakeResult{rows: 0}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("checkRowsAffected() error = %v, want %v", err, ErrNotFound)
	}

	wantErr := errors.New("rows affected failed")

	if err := checkRowsAffected(fakeResult{err: wantErr}); !errors.Is(err, wantErr) {
		t.Fatalf("checkRowsAffected() error = %v, want %v", err, wantErr)
	}
}

func TestOpenCreatesSchema(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)

	var count int

	err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'users'`).Scan(&count)

	if err != nil {
		t.Fatalf("sqlite_master query error = %v", err)
	}

	if count != 1 {
		t.Fatalf("users table count = %d, want 1", count)
	}
}

func TestCreateOrganizationAndUserTrimInput(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)

	orgID, err := CreateOrganization(db, "  Trimmed Org  ")

	if err != nil {
		t.Fatalf("CreateOrganization() error = %v", err)
	}

	org, err := GetOrganization(db, orgID)

	if err != nil {
		t.Fatalf("GetOrganization() error = %v", err)
	}

	if org == nil || org.Name != "Trimmed Org" {
		t.Fatalf("GetOrganization() = %#v, want trimmed name", org)
	}

	userID, err := CreateUser(db, "Demo", "  User@Example.COM  ", "password", nil)

	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	user, err := FindUserByID(db, userID)

	if err != nil {
		t.Fatalf("FindUserByID() error = %v", err)
	}

	if user == nil || user.Email != "user@example.com" {
		t.Fatalf("FindUserByID() = %#v, want normalized email", user)
	}
}

func TestListHelpersReturnHumanTimes(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	now := time.Now().Add(-2 * time.Hour)

	if err := CreateInviteAt(db, "invite_human", "Human Invite", "human@example.com", "Owner", "Queued", now); err != nil {
		t.Fatalf("CreateInviteAt() error = %v", err)
	}

	if err := CreateUploadAt(db, "upload_human", "Human Upload", "human.csv", "42 KB", "Ready", now); err != nil {
		t.Fatalf("CreateUploadAt() error = %v", err)
	}

	if err := CreateApprovalAt(db, "approval_human", "Human Approval", "Approved", now); err != nil {
		t.Fatalf("CreateApprovalAt() error = %v", err)
	}

	for _, item := range []map[string]any{
		ListInvites(db)[0],
		ListUploads(db)[0],
		ListApprovals(db)[0],
	} {
		timeValue, ok := item["time"].(string)

		if !ok || strings.TrimSpace(timeValue) == "" {
			t.Fatalf("time field = %#v, want non-empty string", item["time"])
		}
	}
}

func TestClosedDatabaseErrorPaths(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)

	db.Close()

	if _, err := CreateUser(db, "Closed", "closed@example.com", "password", nil); err == nil {
		t.Fatal("CreateUser() error = nil, want error on closed DB")
	}

	if _, err := FindUserByEmail(db, "closed@example.com"); err == nil {
		t.Fatal("FindUserByEmail() error = nil, want error on closed DB")
	}

	if _, err := FindUserByID(db, 1); err == nil {
		t.Fatal("FindUserByID() error = nil, want error on closed DB")
	}

	if _, err := ListOrganizations(db, ""); err == nil {
		t.Fatal("ListOrganizations() error = nil, want error on closed DB")
	}

	if _, err := GetOrganization(db, 1); err == nil {
		t.Fatal("GetOrganization() error = nil, want error on closed DB")
	}

	if _, err := CreateOrganization(db, "Closed"); err == nil {
		t.Fatal("CreateOrganization() error = nil, want error on closed DB")
	}

	if err := UpdateOrganization(db, 1, "Closed"); err == nil {
		t.Fatal("UpdateOrganization() error = nil, want error on closed DB")
	}

	if _, err := CountOrganizations(db); err == nil {
		t.Fatal("CountOrganizations() error = nil, want error on closed DB")
	}

	if _, err := ListContacts(db, "", false); err == nil {
		t.Fatal("ListContacts() error = nil, want error on closed DB")
	}

	if _, err := ListRecentContacts(db, 1); err == nil {
		t.Fatal("ListRecentContacts() error = nil, want error on closed DB")
	}

	if _, err := ListContactsByOrganization(db, 1); err == nil {
		t.Fatal("ListContactsByOrganization() error = nil, want error on closed DB")
	}

	if _, err := GetContact(db, 1); err == nil {
		t.Fatal("GetContact() error = nil, want error on closed DB")
	}

	if _, err := CreateContact(db, Contact{FirstName: "Closed", LastName: "DB", Email: "closed@example.com"}); err == nil {
		t.Fatal("CreateContact() error = nil, want error on closed DB")
	}

	if err := UpdateContact(db, 1, Contact{}); err == nil {
		t.Fatal("UpdateContact() error = nil, want error on closed DB")
	}

	if err := ToggleContactFavorite(db, 1); err == nil {
		t.Fatal("ToggleContactFavorite() error = nil, want error on closed DB")
	}

	if err := DeleteContact(db, 1); err == nil {
		t.Fatal("DeleteContact() error = nil, want error on closed DB")
	}

	if _, err := ListContactsPaginated(db, "", false, nil, "next", 2); err == nil {
		t.Fatal("ListContactsPaginated() error = nil, want error on closed DB")
	}

	if _, err := ListOrganizationsPaginated(db, "", 1, 2); err == nil {
		t.Fatal("ListOrganizationsPaginated() error = nil, want error on closed DB")
	}

	if _, err := ListContactsByOrgPaginated(db, 1, nil, "next", 2); err == nil {
		t.Fatal("ListContactsByOrgPaginated() error = nil, want error on closed DB")
	}

	if _, err := CreateNote(db, 1, 1, "Closed"); err == nil {
		t.Fatal("CreateNote() error = nil, want error on closed DB")
	}

	if _, err := CreateNoteAt(db, 1, 1, "Closed", time.Now()); err == nil {
		t.Fatal("CreateNoteAt() error = nil, want error on closed DB")
	}

	if _, err := ListContactNotes(db, 1); err == nil {
		t.Fatal("ListContactNotes() error = nil, want error on closed DB")
	}

	if _, err := ListRecentNotes(db, 1); err == nil {
		t.Fatal("ListRecentNotes() error = nil, want error on closed DB")
	}

	if _, err := CountContacts(db); err == nil {
		t.Fatal("CountContacts() error = nil, want error on closed DB")
	}

	if _, err := CountNotes(db); err == nil {
		t.Fatal("CountNotes() error = nil, want error on closed DB")
	}

	if invites := ListInvites(db); invites != nil {
		t.Fatalf("ListInvites() = %#v, want nil on closed DB", invites)
	}

	if uploads := ListUploads(db); uploads != nil {
		t.Fatalf("ListUploads() = %#v, want nil on closed DB", uploads)
	}

	if approvals := ListApprovals(db); approvals != nil {
		t.Fatalf("ListApprovals() = %#v, want nil on closed DB", approvals)
	}
}
