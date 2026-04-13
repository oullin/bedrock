package spark_test

import (
	"testing"

	"github.com/bedrock/packages/spark"
)

// Mirrors BillingServiceTest::test_product_returns_correct_price_id_for_stripe
func TestProduct_PriceIDForProvider_Stripe(t *testing.T) {
	product := &spark.Product{
		StripePriceID: "price_stripe_123",
		PaddlePriceID: "pri_paddle_456",
	}

	got := product.PriceIDForProvider("stripe")

	if got != "price_stripe_123" {
		t.Errorf("PriceIDForProvider(stripe) = %q, want %q", got, "price_stripe_123")
	}
}

// Mirrors BillingServiceTest::test_product_returns_correct_price_id_for_paddle
func TestProduct_PriceIDForProvider_Paddle(t *testing.T) {
	product := &spark.Product{
		StripePriceID: "price_stripe_123",
		PaddlePriceID: "pri_paddle_456",
	}

	got := product.PriceIDForProvider("paddle")

	if got != "pri_paddle_456" {
		t.Errorf("PriceIDForProvider(paddle) = %q, want %q", got, "pri_paddle_456")
	}
}

// Mirrors BillingServiceTest::test_product_returns_null_for_unknown_provider
func TestProduct_PriceIDForProvider_Unknown(t *testing.T) {
	product := &spark.Product{
		StripePriceID: "price_stripe_123",
		PaddlePriceID: "pri_paddle_456",
	}

	got := product.PriceIDForProvider("unknown")

	if got != "" {
		t.Errorf("PriceIDForProvider(unknown) = %q, want empty string", got)
	}
}

// Mirrors BillingServiceTest::test_product_scopes_filter_correctly
func TestProduct_TypeChecksAndScopes(t *testing.T) {
	sub := &spark.Product{Type: spark.ProductTypeSubscription, Active: true}
	oneTime := &spark.Product{Type: spark.ProductTypeOneTime, Active: true}
	inactive := &spark.Product{Type: spark.ProductTypeSubscription, Active: false}

	if !sub.IsSubscription() {
		t.Error("subscription product: IsSubscription() = false, want true")
	}

	if sub.IsOneTime() {
		t.Error("subscription product: IsOneTime() = true, want false")
	}

	if !oneTime.IsOneTime() {
		t.Error("one-time product: IsOneTime() = false, want true")
	}

	if oneTime.IsSubscription() {
		t.Error("one-time product: IsSubscription() = true, want false")
	}

	if inactive.Active {
		t.Error("inactive product: Active = true, want false")
	}

	// Verify active subscription count: only sub is active subscription
	products := []*spark.Product{sub, oneTime, inactive}

	var activeSubs, activeOneTime, active int

	for _, p := range products {
		if p.Active && p.IsSubscription() {
			activeSubs++
		}

		if p.Active && p.IsOneTime() {
			activeOneTime++
		}

		if p.Active {
			active++
		}
	}

	if activeSubs != 1 {
		t.Errorf("active subscriptions = %d, want 1", activeSubs)
	}

	if activeOneTime != 1 {
		t.Errorf("active one-time = %d, want 1", activeOneTime)
	}

	if active != 2 {
		t.Errorf("total active = %d, want 2", active)
	}
}
