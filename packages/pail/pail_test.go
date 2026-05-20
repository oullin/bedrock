package pail

import (
	"strings"
	"testing"
	"time"
)

func TestParseLineExtractsLogFields(t *testing.T) {
	t.Parallel()

	entry := ParseLine(`[2026-04-20 08:00:00] production.ERROR: Payment failed {"user_id":42}`)

	if entry.Level != "error" {
		t.Fatalf("Level = %q, want error", entry.Level)
	}

	if entry.Message != "Payment failed" {
		t.Fatalf("Message = %q, want Payment failed", entry.Message)
	}

	if entry.Context["user_id"].(float64) != 42 {
		t.Fatalf("Context user_id = %#v, want 42", entry.Context["user_id"])
	}
}

func TestCollectFiltersEntries(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(strings.Join([]string{
		`[2026-04-20 08:00:00] production.INFO: Started`,
		`[2026-04-20 08:01:00] production.ERROR: Payment failed`,
	}, "\n"))

	entries, err := Collect(input, Filter{
		Levels:   []string{"error"},
		Contains: "payment",
		Since:    time.Date(2026, 4, 20, 8, 0, 30, 0, time.Local),
	})

	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("len(entries) = %d, want 1", len(entries))
	}

	if entries[0].Level != "error" {
		t.Fatalf("entry level = %q, want error", entries[0].Level)
	}
}
