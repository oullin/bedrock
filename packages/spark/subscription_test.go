package spark_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/spark"
)

func TestSubscriptionIsActive(t *testing.T) {
	tests := []struct {
		status spark.SubscriptionStatus
		want   bool
	}{
		{spark.StatusActive, true},
		{spark.StatusTrialing, true},
		{spark.StatusPastDue, true},
		{spark.StatusCanceled, false},
		{spark.StatusPending, false},
	}

	for _, tt := range tests {
		s := &spark.Subscription{Status: tt.status}

		if got := s.IsActive(); got != tt.want {
			t.Errorf("Subscription{Status: %q}.IsActive() = %v, want %v", tt.status, got, tt.want)
		}
	}
}

func TestSubscriptionIsNew(t *testing.T) {
	if !(&spark.Subscription{}).IsNew() {
		t.Error("zero-ID subscription should be new")
	}

	if (&spark.Subscription{ID: 1}).IsNew() {
		t.Error("non-zero-ID subscription should not be new")
	}
}

func TestSubscriptionOnTrial(t *testing.T) {
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	clock := fixedClock{now: now}

	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)

	s := &spark.Subscription{TrialEndsAt: &future}

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

	s := &spark.Subscription{Status: spark.StatusCanceled, EndsAt: &future}

	if !s.OnGracePeriod(clock) {
		t.Error("canceled subscription with future EndsAt should be on grace period")
	}

	s.Status = spark.StatusActive

	if s.OnGracePeriod(clock) {
		t.Error("active subscription should not be on grace period")
	}
}

func TestSubscriptionOnPausedGracePeriod(t *testing.T) {
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	clock := fixedClock{now: now}
	future := now.Add(24 * time.Hour)

	s := &spark.Subscription{Status: spark.StatusPaused, PausedAt: &future}

	if !s.OnPausedGracePeriod(clock) {
		t.Error("paused subscription with future PausedAt should be on paused grace period")
	}

	s.PausedAt = nil

	if s.OnPausedGracePeriod(clock) {
		t.Error("paused subscription with nil PausedAt should not be on paused grace period")
	}
}

func TestSubscriptionHasProduct(t *testing.T) {
	s := &spark.Subscription{
		Items: []spark.SubscriptionItem{
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
	s := &spark.Subscription{
		Items: []spark.SubscriptionItem{
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

	active := &spark.Subscription{Status: spark.StatusActive}

	if !active.Valid(clock, false) {
		t.Error("active subscription should be valid")
	}

	pastDue := &spark.Subscription{Status: spark.StatusPastDue}

	if pastDue.Valid(clock, false) {
		t.Error("past_due subscription should not be valid without keepPastDueActive")
	}

	if !pastDue.Valid(clock, true) {
		t.Error("past_due subscription should be valid with keepPastDueActive")
	}

	future := now.Add(24 * time.Hour)
	trialing := &spark.Subscription{Status: spark.StatusTrialing, TrialEndsAt: &future}

	if !trialing.Valid(clock, false) {
		t.Error("trialing subscription with future trial should be valid")
	}
}

func TestSubscriptionProration(t *testing.T) {
	s := &spark.Subscription{}

	if s.ProrationBehavior() != spark.ProratedNextBillingPeriod {
		t.Errorf("default proration should be ProratedNextBillingPeriod, got %q", s.ProrationBehavior())
	}

	s.NoProrate()

	if s.ProrationBehavior() != spark.FullNextBillingPeriod {
		t.Errorf("after NoProrate, got %q", s.ProrationBehavior())
	}

	s.ProrateImmediately()

	if s.ProrationBehavior() != spark.ProratedImmediately {
		t.Errorf("after ProrateImmediately, got %q", s.ProrationBehavior())
	}
}
