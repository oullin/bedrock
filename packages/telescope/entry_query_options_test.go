package telescope_test

import (
	"testing"

	"github.com/bedrock/packages/telescope"
)

func TestDefaultQueryOptionsHasLimit50(t *testing.T) {
	t.Parallel()

	opts := telescope.DefaultQueryOptions()

	if opts.Limit != 50 {
		t.Fatalf("expected default limit=50, got %d", opts.Limit)
	}
}

func TestEntryQueryOptionsForBatchID(t *testing.T) {
	t.Parallel()

	opts := telescope.DefaultQueryOptions().ForBatchID("batch-xyz")

	if opts.BatchID != "batch-xyz" {
		t.Fatalf("expected batch_id=batch-xyz, got %q", opts.BatchID)
	}

	// Original should not be mutated.
	original := telescope.DefaultQueryOptions()

	if original.BatchID != "" {
		t.Fatal("DefaultQueryOptions should return a fresh zero-value struct")
	}
}

func TestEntryQueryOptionsWithTag(t *testing.T) {
	t.Parallel()

	opts := telescope.DefaultQueryOptions().WithTag("billing")

	if opts.Tag != "billing" {
		t.Fatalf("expected tag=billing, got %q", opts.Tag)
	}
}

func TestEntryQueryOptionsWithLimit(t *testing.T) {
	t.Parallel()

	opts := telescope.DefaultQueryOptions().WithLimit(100)

	if opts.Limit != 100 {
		t.Fatalf("expected limit=100, got %d", opts.Limit)
	}
}

func TestEntryQueryOptionsWithFamilyHash(t *testing.T) {
	t.Parallel()

	opts := telescope.DefaultQueryOptions().WithFamilyHash("abc123")

	if opts.FamilyHash != "abc123" {
		t.Fatalf("expected family_hash=abc123, got %q", opts.FamilyHash)
	}
}

func TestEntryQueryOptionsBefore(t *testing.T) {
	t.Parallel()

	opts := telescope.DefaultQueryOptions().Before(999)

	if opts.BeforeSequence != 999 {
		t.Fatalf("expected before_sequence=999, got %d", opts.BeforeSequence)
	}
}

func TestEntryQueryOptionsImmutable(t *testing.T) {
	t.Parallel()

	base := telescope.DefaultQueryOptions()
	derived := base.WithTag("foo").WithLimit(25).ForBatchID("b1")

	// base must remain unchanged.
	if base.Tag != "" || base.Limit != 50 || base.BatchID != "" {
		t.Fatal("DefaultQueryOptions should not be mutated by fluent methods")
	}

	// derived must have all three values.
	if derived.Tag != "foo" || derived.Limit != 25 || derived.BatchID != "b1" {
		t.Fatalf("derived options incorrect: %+v", derived)
	}
}
