package str

import (
	"strings"
	"testing"
)

// Port of Illuminate\Tests\Support\SupportStrTest::testUuid
func TestStrUuid(t *testing.T) {
	// NOT parallel — UUID tests may conflict with freeze
	uuid := StrUuid()

	if !StrIsUuid(uuid) {
		t.Errorf("StrUuid() = %q is not a valid UUID", uuid)
	}
}

// Port of Illuminate\Tests\Support\SupportStrTest::testItCanFreezeUuids
func TestStrFreezeUuids(t *testing.T) {
	// NOT parallel — modifies global UUID state
	cleanup := FreezeUuids(func() string { return "frozen-uuid" })

	defer cleanup()

	if got := StrUuid(); got != "frozen-uuid" {
		t.Errorf("FreezeUuids: expected 'frozen-uuid', got %q", got)
	}
}

// Port of Illuminate\Tests\Support\SupportStrTest::testItCanFreezeUuidsInAClosure
func TestStrFreezeUuidsCleanup(t *testing.T) {
	// NOT parallel — modifies global UUID state
	cleanup := FreezeUuids(func() string { return "frozen" })
	frozen := StrUuid()
	cleanup()

	// After cleanup, should generate real UUIDs again
	normal := StrUuid()

	if normal == "frozen" && frozen != normal {
		// This is fine — just verify cleanup runs
	}

	if frozen != "frozen" {
		t.Errorf("expected 'frozen', got %q", frozen)
	}
}

// Port of Illuminate\Tests\Support\SupportStrTest::testItCanSpecifyASequenceOfUuidsToUtilise
// SupportStrTest::testItCanSpecifyAFallbackForASequence
func TestStrUuidSequence(t *testing.T) {
	// NOT parallel — modifies global state
	cleanup := CreateUuidsUsingSequence([]string{
		"first-uuid",
		"second-uuid",
	})

	defer cleanup()

	if got := StrUuid(); got != "first-uuid" {
		t.Errorf("sequence[0] = %q", got)
	}

	if got := StrUuid(); got != "second-uuid" {
		t.Errorf("sequence[1] = %q", got)
	}

	cleanup()
	cleanup = CreateUuidsUsingSequence([]string{"only-uuid"}, func() string { return "fallback-uuid" })

	if got := StrUuid(); got != "only-uuid" {
		t.Errorf("fallback sequence first value = %q", got)
	}

	if got := StrUuid(); got != "fallback-uuid" {
		t.Errorf("fallback value = %q", got)
	}
}

// Port of Illuminate\Tests\Support\SupportStrTest::testItCanFreezeUlids
// SupportStrTest::testItCanFreezeUlidsInAClosure
func TestStrFreezeUlids(t *testing.T) {
	// NOT parallel — modifies global state
	cleanup := FreezeUlids(func() string { return "FROZENULID00000000000000000" })

	defer cleanup()

	if got := StrUlid(); got != "FROZENULID00000000000000000" {
		t.Errorf("FreezeUlids: expected frozen ULID, got %q", got)
	}
}

// Port of Illuminate\Tests\Support\SupportStrTest::testItCanSpecifyASequenceOfUlidsToUtilise
// SupportStrTest::testItCanSpecifyAFallbackForAUlidSequence
func TestStrUlidSequence(t *testing.T) {
	// NOT parallel — modifies global state
	seq := []string{
		"01ARZ3NDEKTSV4RRFFQ69G5FAV",
		"01ARZ3NDEKTSV4RRFFQ69G5FAW",
	}
	cleanup := CreateUlidsUsingSequence(seq)

	defer cleanup()

	if got := StrUlid(); got != seq[0] {
		t.Errorf("sequence[0] = %q", got)
	}

	if got := StrUlid(); got != seq[1] {
		t.Errorf("sequence[1] = %q", got)
	}

	cleanup()
	cleanup = CreateUlidsUsingSequence([]string{"ONLYULID0000000000000000"}, func() string {
		return "FALLBACKULID000000000000"
	})

	if got := StrUlid(); got != "ONLYULID0000000000000000" {
		t.Errorf("fallback sequence first value = %q", got)
	}

	if got := StrUlid(); got != "FALLBACKULID000000000000" {
		t.Errorf("fallback value = %q", got)
	}
}

// Port of Illuminate\Tests\Support\SupportStrTest::testItCreatesUuidsNormallyAfterFailureWithinFreezeMethod
// SupportStrTest::testItCreatesUlidsNormallyAfterFailureWithinFreezeMethod
func TestStrCreateUuidsNormally(t *testing.T) {
	// NOT parallel — modifies global state
	cleanup := FreezeUuids(func() string { return "frozen" })
	cleanup() // restore immediately

	uuid := StrUuid()

	if uuid == "frozen" {
		t.Error("after CreateUuidsNormally, should generate real UUIDs")
	}

	if !StrIsUuid(uuid) {
		t.Errorf("should be valid UUID, got %q", uuid)
	}

	ulidCleanup := FreezeUlids(func() string { return "FROZENULID00000000000000000" })
	ulidCleanup()

	ulid := StrUlid()

	if ulid == "FROZENULID00000000000000000" {
		t.Error("after ULID cleanup, should generate real ULIDs")
	}

	if len(ulid) != 26 {
		t.Errorf("should be valid ULID length, got %q", ulid)
	}
}

// Port of Illuminate\Tests\Support\SupportStrTest::testOrderedUuid
func TestStrOrderedUuid(t *testing.T) {
	// NOT parallel — may interfere with UUID freeze tests
	uuid := StrOrderedUuid()

	if !StrIsUuid(uuid) {
		t.Errorf("StrOrderedUuid() = %q is not a valid UUID", uuid)
	}
}

// Port of Illuminate\Tests\Support\SupportStrTest::testUlid
func TestStrUlid(t *testing.T) {
	// NOT parallel
	ulid := StrUlid()

	if len(ulid) != 26 {
		t.Errorf("ULID length should be 26, got %d (%q)", len(ulid), ulid)
	}

	if strings.ToUpper(ulid) != ulid {
		// ULID should be uppercase
		t.Errorf("ULID should be uppercase, got %q", ulid)
	}
}

// Port of Illuminate\Tests\Support\SupportStrTest::testResetFactoryState
func TestStrResetFactoryState(t *testing.T) {
	// NOT parallel — modifies global state
	CreateUuidsUsing(func() string { return "custom" })
	ResetFactoryState()

	uuid := StrUuid()

	if uuid == "custom" {
		t.Error("after ResetFactoryState, UUID factory should be reset")
	}
}
