package spark_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/billing"
)

// Mirrors OrderTest::test_order_can_be_marked_as_completed
func TestOrder_MarkAsCompleted(t *testing.T) {
	order := &billing.Order{
		ID:     1,
		Status: billing.OrderStatusPending,
	}

	if !order.IsPending() {
		t.Fatal("new order should be pending")
	}

	now := time.Now()
	order.MarkAsCompleted("txn_test_123", now)

	if order.Status != billing.OrderStatusCompleted {
		t.Errorf("status = %q, want %q", order.Status, billing.OrderStatusCompleted)
	}

	if order.ProviderTransactionID != "txn_test_123" {
		t.Errorf("ProviderTransactionID = %q, want %q", order.ProviderTransactionID, "txn_test_123")
	}

	if order.CompletedAt == nil {
		t.Fatal("CompletedAt should not be nil after completion")
	}

	if !order.IsCompleted() {
		t.Error("IsCompleted() = false after MarkAsCompleted")
	}
}

// Mirrors OrderTest::test_team_can_check_if_product_purchased
// This tests the Order status checks used by the store's
// HasCompletedForProduct method.
func TestOrder_HasPurchased(t *testing.T) {
	pendingOrder := &billing.Order{
		TeamID:    1,
		ProductID: 10,
		Status:    billing.OrderStatusPending,
	}

	// Pending order should not count as purchased.
	if pendingOrder.IsCompleted() {
		t.Error("pending order should not be completed")
	}

	completedOrder := &billing.Order{
		TeamID:    1,
		ProductID: 10,
		Status:    billing.OrderStatusCompleted,
	}

	if !completedOrder.IsCompleted() {
		t.Error("completed order should be completed")
	}
}
