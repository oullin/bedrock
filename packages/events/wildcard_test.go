package events_test

import (
	"testing"

	"github.com/bedrock/packages/events"
)

func TestMatchesWildcard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		pattern string
		name    string
		want    bool
	}{
		// Trailing wildcard matches any suffix.
		{"order.*", "order.created", true},
		{"order.*", "order.shipped", true},
		{"order.*", "order.item.created", true},

		// Leading wildcard matches any prefix segment.
		{"*.created", "order.created", true},
		{"*.created", "user.created", true},

		// Middle wildcard.
		{"order.*.shipped", "order.item.shipped", true},
		{"order.*.shipped", "order.item.created", false},

		// No match.
		{"order.*", "user.created", false},
		{"order.created", "order.shipped", false},

		// Exact match (no wildcard).
		{"order.created", "order.created", true},

		// Single wildcard matches everything.
		{"*", "order.created", true},
		{"*", "anything", true},

		// Segment count mismatch (non-trailing wildcard).
		{"order.*.shipped", "order.shipped", false},

		// Empty pattern/name.
		{"", "", true},
		{"order.*", "", false},
	}

	for _, tt := range tests {
		got := events.ExportMatchesWildcard(tt.pattern, tt.name)

		if got != tt.want {
			t.Errorf("matchesWildcard(%q, %q) = %v, want %v", tt.pattern, tt.name, got, tt.want)
		}
	}
}
