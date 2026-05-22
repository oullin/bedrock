package billing_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/billing"
)

// SubscriptionStatusTransitionTest::test_active_to_past_due_keeps_entitlements_active
// SubscriptionStatusTransitionTest::test_past_due_to_active_preserves_entitlements
// SubscriptionStatusTransitionTest::test_active_to_paused_deactivates_entitlements
// SubscriptionStatusTransitionTest::test_paused_to_active_reactivates_entitlements
// SubscriptionStatusTransitionTest::test_active_to_canceled_deactivates_entitlements
// SubscriptionStatusTransitionTest::test_cancel_from_past_due_deactivates_entitlements
func TestMadoraSubscriptionStatusTransitions(t *testing.T) {
	now := time.Date(2026, 4, 22, 9, 0, 0, 0, time.UTC)

	t.Run("mark past due keeps the subscription valid", func(t *testing.T) {
		sub := &billing.Subscription{Status: billing.StatusActive}

		if !sub.MarkPastDue(now) {
			t.Fatalf("MarkPastDue returned false")
		}

		if sub.Status != billing.StatusPastDue {
			t.Fatalf("status = %s, want past_due", sub.Status)
		}

		if !sub.Valid() || !sub.Status.GrantsAccess() {
			t.Fatalf("past_due subscription should still grant access")
		}
	})

	t.Run("activate restores access from past due", func(t *testing.T) {
		sub := &billing.Subscription{Status: billing.StatusPastDue}

		if !sub.Activate(now) {
			t.Fatalf("Activate returned false")
		}

		if sub.Status != billing.StatusActive {
			t.Fatalf("status = %s, want active", sub.Status)
		}

		if !sub.Valid() || !sub.Status.GrantsAccess() {
			t.Fatalf("active subscription should grant access")
		}
	})

	t.Run("pause deactivates access", func(t *testing.T) {
		sub := &billing.Subscription{Status: billing.StatusActive}

		if !sub.Pause(now) {
			t.Fatalf("Pause returned false")
		}

		if sub.Status != billing.StatusPaused || sub.Valid() {
			t.Fatalf("paused subscription should be inactive: %#v", sub)
		}
	})

	t.Run("reactivating a paused subscription restores access", func(t *testing.T) {
		sub := &billing.Subscription{Status: billing.StatusActive}

		if !sub.Pause(now) {
			t.Fatalf("Pause returned false")
		}

		if !sub.Activate(now.Add(time.Hour)) {
			t.Fatalf("Activate returned false")
		}

		if sub.Status != billing.StatusActive || !sub.Valid() {
			t.Fatalf("reactivated subscription should be active: %#v", sub)
		}
	})

	t.Run("cancel deactivates an active subscription immediately", func(t *testing.T) {
		sub := &billing.Subscription{Status: billing.StatusActive}

		if !sub.Cancel(now) {
			t.Fatalf("Cancel returned false")
		}

		if sub.Status != billing.StatusCanceled || sub.Valid() || sub.OnGracePeriod() {
			t.Fatalf("canceled subscription should not grant access: %#v", sub)
		}
	})

	t.Run("cancel also deactivates a past due subscription", func(t *testing.T) {
		sub := &billing.Subscription{Status: billing.StatusPastDue}

		if !sub.Cancel(now) {
			t.Fatalf("Cancel returned false")
		}

		if sub.Status != billing.StatusCanceled || sub.Valid() || sub.OnGracePeriod() {
			t.Fatalf("canceled past_due subscription should not grant access: %#v", sub)
		}
	})
}

// BillingLifecycleTest::test_pending_subscription_is_idempotent_under_sequential_calls
// BillingLifecycleTest::test_starter_trial_is_idempotent_under_sequential_calls
func TestMadoraSubscriptionConstructionIsDeterministic(t *testing.T) {
	now := time.Date(2026, 4, 22, 10, 0, 0, 0, time.UTC)
	billable := madoraBillable{id: 10, typ: "team"}

	firstPending := billing.NewPendingSubscription(billable, "starter", now)
	secondPending := billing.NewPendingSubscription(billable, "starter", now)

	if firstPending.Status != billing.StatusPending || secondPending.Status != billing.StatusPending {
		t.Fatalf("pending subscriptions should start in pending status")
	}

	if !firstPending.PendingExpiresAt.Equal(*secondPending.PendingExpiresAt) {
		t.Fatalf("pending expiry timestamps should match: %#v %#v", firstPending.PendingExpiresAt, secondPending.PendingExpiresAt)
	}

	firstTrial := billing.NewTrialSubscription(billable, "starter", 14, now)
	secondTrial := billing.NewTrialSubscription(billable, "starter", 14, now)

	if firstTrial.Status != billing.StatusTrialing || secondTrial.Status != billing.StatusTrialing {
		t.Fatalf("trial subscriptions should start in trialing status")
	}

	if !firstTrial.TrialEndsAt.Equal(*secondTrial.TrialEndsAt) {
		t.Fatalf("trial expiry timestamps should match: %#v %#v", firstTrial.TrialEndsAt, secondTrial.TrialEndsAt)
	}
}
