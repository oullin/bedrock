package telescope_test

import (
	"testing"

	"github.com/bedrock/packages/telescope"
)

func TestNewEntryGeneratesUUID(t *testing.T) {
	t.Parallel()

	e := telescope.NewEntry(telescope.EntryTypeLog, nil)

	if e.UUID == "" {
		t.Fatal("expected non-empty UUID")
	}
}

func TestNewEntryHasTimestamp(t *testing.T) {
	t.Parallel()

	e := telescope.NewEntry(telescope.EntryTypeLog, nil)

	if e.RecordedAt.IsZero() {
		t.Fatal("expected non-zero RecordedAt")
	}
}

func TestEntryWithTypeChangesType(t *testing.T) {
	t.Parallel()

	e := telescope.NewEntry(telescope.EntryTypeLog, nil)
	e.WithType(telescope.EntryTypeCache)

	if e.Type != telescope.EntryTypeCache {
		t.Fatalf("expected type=%q, got %q", telescope.EntryTypeCache, e.Type)
	}
}

func TestEntryWithBatchID(t *testing.T) {
	t.Parallel()

	e := telescope.NewEntry(telescope.EntryTypeLog, nil)
	e.WithBatchID("batch-abc-123")

	if e.BatchID != "batch-abc-123" {
		t.Fatalf("unexpected batch ID: %q", e.BatchID)
	}
}

func TestEntryWithFamilyHash(t *testing.T) {
	t.Parallel()

	e := telescope.NewEntry(telescope.EntryTypeQuery, nil)
	e.WithFamilyHash("abc123")

	if e.FamilyHash != "abc123" {
		t.Fatalf("unexpected family hash: %q", e.FamilyHash)
	}
}

func TestEntryAddTagsDeduplicates(t *testing.T) {
	t.Parallel()

	e := telescope.NewEntry(telescope.EntryTypeLog, nil)
	e.AddTags("foo", "bar", "foo")

	if len(e.Tags) != 2 {
		t.Fatalf("expected 2 unique tags, got %d: %v", len(e.Tags), e.Tags)
	}
}

func TestEntryWithUserAddsContentAndTag(t *testing.T) {
	t.Parallel()

	e := telescope.NewEntry(telescope.EntryTypeRequest, map[string]any{})
	e.WithUser(&telescope.EntryUser{ID: 7, Name: "Alice", Email: "alice@example.com"})

	user, ok := e.Content["user"].(map[string]any)

	if !ok {
		t.Fatal("expected user map in content")
	}

	if user["id"] != 7 {
		t.Fatalf("expected user.id=7, got %v", user["id"])
	}

	assertHasTag(t, e.Tags, "user:7")
}

func TestEntryHasMonitoredTag(t *testing.T) {
	t.Parallel()

	e := telescope.NewEntry(telescope.EntryTypeLog, nil)
	e.AddTags("billing", "admin")

	if !e.HasMonitoredTag([]string{"billing"}) {
		t.Fatal("expected HasMonitoredTag to return true for 'billing'")
	}

	if e.HasMonitoredTag([]string{"unknown"}) {
		t.Fatal("expected HasMonitoredTag to return false for 'unknown'")
	}
}

func TestEntryTypeCheckers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		entryType string
		check     func(*telescope.IncomingEntry) bool
	}{
		{telescope.EntryTypeRequest, (*telescope.IncomingEntry).IsRequest},
		{telescope.EntryTypeQuery, (*telescope.IncomingEntry).IsQuery},
		{telescope.EntryTypeEvent, (*telescope.IncomingEntry).IsEvent},
		{telescope.EntryTypeCache, (*telescope.IncomingEntry).IsCache},
		{telescope.EntryTypeException, (*telescope.IncomingEntry).IsException},
		{telescope.EntryTypeLog, (*telescope.IncomingEntry).IsLog},
		{telescope.EntryTypeGate, (*telescope.IncomingEntry).IsGate},
		{telescope.EntryTypeScheduledTask, (*telescope.IncomingEntry).IsScheduledTask},
		{telescope.EntryTypeClientRequest, (*telescope.IncomingEntry).IsClientRequest},
		{telescope.EntryTypeMail, (*telescope.IncomingEntry).IsMail},
		{telescope.EntryTypeNotification, (*telescope.IncomingEntry).IsNotification},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.entryType, func(t *testing.T) {
			t.Parallel()

			e := telescope.NewEntry(tc.entryType, nil)

			if !tc.check(e) {
				t.Fatalf("expected IsXxx() to return true for type %q", tc.entryType)
			}

			// Any other type should return false.
			other := telescope.NewEntry(telescope.EntryTypeDump, nil)

			if tc.check(other) {
				t.Fatalf("expected IsXxx() to return false for type %q when entry is %q",
					tc.entryType, other.Type)
			}
		})
	}
}

func TestEntryIsSlowQuery(t *testing.T) {
	t.Parallel()

	slow := telescope.NewEntry(telescope.EntryTypeQuery, map[string]any{"slow": true})

	if !slow.IsSlowQuery() {
		t.Fatal("expected IsSlowQuery()=true")
	}

	fast := telescope.NewEntry(telescope.EntryTypeQuery, map[string]any{"slow": false})

	if fast.IsSlowQuery() {
		t.Fatal("expected IsSlowQuery()=false for fast query")
	}
}

func TestEntryIsFailedJob(t *testing.T) {
	t.Parallel()

	failed := telescope.NewEntry(telescope.EntryTypeJob, map[string]any{"status": "failed"})

	if !failed.IsFailedJob() {
		t.Fatal("expected IsFailedJob()=true")
	}

	ok := telescope.NewEntry(telescope.EntryTypeJob, map[string]any{"status": "processed"})

	if ok.IsFailedJob() {
		t.Fatal("expected IsFailedJob()=false for processed job")
	}
}
