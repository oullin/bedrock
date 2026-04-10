package spark_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/billing"
)

func TestSubscriptionIsActive(t *testing.T) {
	tests := []struct {
		status billing.SubscriptionStatus
		want   bool
	}{
		{billing.StatusActive, true},
		{billing.StatusTrialing, true},
		{billing.StatusPastDue, true},
		{billing.StatusCanceled, false},
		{billing.StatusPending, false},
	}

	for _, tt := range tests {
		s := &billing.Subscription{Status: tt.status}

		if got := s.IsActive(); got != tt.want {
			t.Errorf("Subscription{Status: %q}.IsActive() = %v, want %v", tt.status, got, tt.want)
		}
	}
}

func TestSubscriptionIsNew(t *testing.T) {
	if !(&billing.Subscription{}).IsNew() {
		t.Error("zero-ID subscription should be new")
	}

	if (&billing.Subscription{ID: 1}).IsNew() {
		t.Error("non-zero-ID subscription should not be new")
	}
}

func TestSubscriptionOnTrial(t *testing.T) {
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	clock := fixedClock{now: now}

	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)

	s := &billing.Subscription{TrialEndsAt: &future}

	if !s.OnTrial(clock) {
		t.Error("should be on trial with future end date")
	}

	s.TrialEndsAt = &past

	if s.OnTrial(clock) {
		t.Error("should not be on trial with past end date")
	}

	s.TrialEndsAt = nil

	if s.OnTrial(clock) {
		t.Error("should not be on trial with nil end date")
	}
}

func TestSubscriptionOnGracePeriod(t *testing.T) {
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	clock := fixedClock{now: now}
	future := now.Add(24 * time.Hour)

	s := &billing.Subscription{Status: billing.StatusCanceled, EndsAt: &future}

	if !s.OnGracePeriod(clock) {
		t.Error("canceled subscription with future EndsAt should be on grace period")
	}

	s.Status = billing.StatusActive

	if s.OnGracePeriod(clock) {
		t.Error("active subscription should not be on grace period")
	}
}

func TestSubscriptionOnPausedGracePeriod(t *testing.T) {
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	clock := fixedClock{now: now}
	future := now.Add(24 * time.Hour)

	s := &billing.Subscription{Status: billing.StatusPaused, PausedAt: &future}

	if !s.OnPausedGracePeriod(clock) {
		t.Error("paused subscription with future PausedAt should be on paused grace period")
	}

	s.PausedAt = nil

	if s.OnPausedGracePeriod(clock) {
		t.Error("paused subscription with nil PausedAt should not be on paused grace period")
	}
}

func TestSubscriptionHasProduct(t *testing.T) {
	s := &billing.Subscription{
		Items: []billing.SubscriptionItem{
			{ProductID: "prod_abc", PriceID: "pri_123"},
			{ProductID: "prod_def", PriceID: "pri_456"},
		},
	}

	if !s.HasProduct("prod_abc") {
		t.Error("should have product prod_abc")
	}

	if s.HasProduct("prod_xyz") {
		t.Error("should not have product prod_xyz")
	}
}

func TestSubscriptionHasPrice(t *testing.T) {
	s := &billing.Subscription{
		Items: []billing.SubscriptionItem{
			{PriceID: "pri_123"},
		},
	}

	if !s.HasPrice("pri_123") {
		t.Error("should have price pri_123")
	}

	if s.HasPrice("pri_999") {
		t.Error("should not have price pri_999")
	}
}

func TestSubscriptionValid(t *testing.T) {
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	clock := fixedClock{now: now}

	active := &billing.Subscription{Status: billing.StatusActive}

	if !active.Valid(clock, false) {
		t.Error("active subscription should be valid")
	}

	pastDue := &billing.Subscription{Status: billing.StatusPastDue}

	if pastDue.Valid(clock, false) {
		t.Error("past_due subscription should not be valid without keepPastDueActive")
	}

	if !pastDue.Valid(clock, true) {
		t.Error("past_due subscription should be valid with keepPastDueActive")
	}

	future := now.Add(24 * time.Hour)
	trialing := &billing.Subscription{Status: billing.StatusTrialing, TrialEndsAt: &future}

	if !trialing.Valid(clock, false) {
		t.Error("trialing subscription with future trial should be valid")
	}
}

func TestSubscriptionProration(t *testing.T) {
	s := &billing.Subscription{}

	if s.ProrationBehavior() != billing.ProratedNextBillingPeriod {
		t.Errorf("default proration should be ProratedNextBillingPeriod, got %q", s.ProrationBehavior())
	}

	s.NoProrate()

	if s.ProrationBehavior() != billing.FullNextBillingPeriod {
		t.Errorf("after NoProrate, got %q", s.ProrationBehavior())
	}

	s.ProrateImmediately()

	if s.ProrationBehavior() != billing.ProratedImmediately {
		t.Errorf("after ProrateImmediately, got %q", s.ProrationBehavior())
	}
}
