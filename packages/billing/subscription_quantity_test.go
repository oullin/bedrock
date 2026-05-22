package billing_test

import (
	"testing"

	"github.com/bedrock/packages/billing"
)

func TestSubscriptionQuantity(t *testing.T) {
	tests := []struct {
		name  string
		items []billing.SubscriptionItem
		want  int
	}{
		{name: "no items", items: nil, want: 1},
		{name: "zero quantity", items: []billing.SubscriptionItem{{Quantity: 0}}, want: 1},
		{name: "negative quantity", items: []billing.SubscriptionItem{{Quantity: -2}}, want: 1},
		{name: "positive quantity", items: []billing.SubscriptionItem{{Quantity: 5}}, want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &billing.Subscription{Items: tt.items}

			if got := sub.Quantity(); got != tt.want {
				t.Fatalf("Quantity() = %d, want %d", got, tt.want)
			}
		})
	}
}
