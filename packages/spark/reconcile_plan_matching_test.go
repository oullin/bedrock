package spark_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/spark"
)

// ReconcileSubscriptionPlanMatchingTest::test_reconciles_correct_pre_paddle_subscription_when_multiple_exist
// ReconcileSubscriptionPlanMatchingTest::test_reconciles_when_the_pending_price_snapshot_is_missing
// ReconcileSubscriptionPlanMatchingTest::test_skips_reconciliation_when_no_price_ids_available
func TestMadoraReconcileSubscriptionPlanMatching(t *testing.T) {
	t.Run("selects the pending subscription matched by price", func(t *testing.T) {
		matching := &spark.Subscription{
			ID:           10,
			BillableType: "team",
			BillableID:   42,
			Plan:         "starter",
			Status:       spark.StatusAwaitingPayment,
			Items:        []spark.SubscriptionItem{{PriceID: "pri_match"}},
		}
		other := &spark.Subscription{
			ID:           11,
			BillableType: "team",
			BillableID:   42,
			Plan:         "pro",
			Status:       spark.StatusAwaitingPayment,
			Items:        []spark.SubscriptionItem{{PriceID: "pri_other"}},
		}
		cashier := &spark.Subscription{
			ID:           12,
			BillableType: "team",
			BillableID:   42,
			PaddleID:     "sub_real",
			Plan:         "starter",
			Status:       spark.StatusActive,
			Items:        []spark.SubscriptionItem{{PriceID: "pri_match"}},
		}
		store := &madoraSubStore{subs: []*spark.Subscription{other, matching, cashier}}
		events := &madoraDispatcher{}

		reconciled, err := spark.ReconcileSubscriptionAfterCheckout(context.Background(), store, "team", 42, cashier, events)
		if err != nil {
			t.Fatalf("reconcile checkout: %v", err)
		}
		if !reconciled {
			t.Fatalf("expected reconciliation to match the pending price snapshot")
		}
		if len(store.saved) != 1 || store.saved[0].ID != cashier.ID {
			t.Fatalf("saved subscriptions = %#v, want cashier", store.saved)
		}
		if len(store.deleted) != 1 || store.deleted[0] != matching.ID {
			t.Fatalf("deleted subscriptions = %#v, want matching pending row", store.deleted)
		}
	})

	t.Run("falls back to the plan when the pending price snapshot is missing", func(t *testing.T) {
		pending := &spark.Subscription{
			ID:           20,
			BillableType: "team",
			BillableID:   42,
			Plan:         "pro",
			Status:       spark.StatusPending,
		}
		cashier := &spark.Subscription{
			ID:           21,
			BillableType: "team",
			BillableID:   42,
			PaddleID:     "sub_real",
			Plan:         "pro",
			Status:       spark.StatusActive,
			Items:        []spark.SubscriptionItem{{PriceID: "pri_pro_monthly"}},
		}
		store := &madoraSubStore{subs: []*spark.Subscription{pending, cashier}}

		reconciled, err := spark.ReconcileSubscriptionAfterCheckout(context.Background(), store, "team", 42, cashier, nil)
		if err != nil {
			t.Fatalf("reconcile checkout: %v", err)
		}
		if !reconciled {
			t.Fatalf("expected reconciliation to fall back to plan matching")
		}
		if len(store.deleted) != 1 || store.deleted[0] != pending.ID {
			t.Fatalf("deleted subscriptions = %#v, want pending row", store.deleted)
		}
	})

	t.Run("does not reconcile when neither a price nor a plan match exists", func(t *testing.T) {
		pending := &spark.Subscription{
			ID:           30,
			BillableType: "team",
			BillableID:   42,
			Status:       spark.StatusPending,
		}
		cashier := &spark.Subscription{
			ID:           31,
			BillableType: "team",
			BillableID:   42,
			PaddleID:     "sub_real",
			Plan:         "pro",
			Status:       spark.StatusActive,
			Items:        []spark.SubscriptionItem{{PriceID: "pri_pro_monthly"}},
		}
		store := &madoraSubStore{subs: []*spark.Subscription{pending, cashier}}

		reconciled, err := spark.ReconcileSubscriptionAfterCheckout(context.Background(), store, "team", 42, cashier, nil)
		if err != nil {
			t.Fatalf("reconcile checkout: %v", err)
		}
		if reconciled {
			t.Fatalf("expected reconciliation to skip when no price ids are available")
		}
		if len(store.deleted) != 0 || len(store.saved) != 0 {
			t.Fatalf("store mutated despite no matching pending row: saved=%#v deleted=%#v", store.saved, store.deleted)
		}
	})
}
