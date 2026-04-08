package spark_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/billing"
)

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

func TestCustomerOnGenericTrial(t *testing.T) {
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	clock := fixedClock{now: now}

	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)

	tests := []struct {
		name     string
		trialEnd *time.Time
		want     bool
	}{
		{"nil trial", nil, false},
		{"future trial", &future, true},
		{"past trial", &past, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &billing.Customer{TrialEndsAt: tt.trialEnd}
			if got := c.OnGenericTrial(clock); got != tt.want {
				t.Errorf("OnGenericTrial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomerHasExpiredGenericTrial(t *testing.T) {
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	clock := fixedClock{now: now}

	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)
	exact := now

	tests := []struct {
		name     string
		trialEnd *time.Time
		want     bool
	}{
		{"nil trial", nil, false},
		{"future trial", &future, false},
		{"past trial", &past, true},
		{"exact now", &exact, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &billing.Customer{TrialEndsAt: tt.trialEnd}
			if got := c.HasExpiredGenericTrial(clock); got != tt.want {
				t.Errorf("HasExpiredGenericTrial() = %v, want %v", got, tt.want)
			}
		})
	}
}
