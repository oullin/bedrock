package spark_test

import (
	"testing"

	"github.com/bedrock/packages/spark"
)

func TestSubscriptionQuantity(t *testing.T) {
	tests := []struct {
		name  string
		items []spark.SubscriptionItem
		want  int
	}{
		{name: "no items", items: nil, want: 1},
		{name: "zero quantity", items: []spark.SubscriptionItem{{Quantity: 0}}, want: 1},
		{name: "negative quantity", items: []spark.SubscriptionItem{{Quantity: -2}}, want: 1},
		{name: "positive quantity", items: []spark.SubscriptionItem{{Quantity: 5}}, want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &spark.Subscription{Items: tt.items}

			if got := sub.Quantity(); got != tt.want {
				t.Fatalf("Quantity() = %d, want %d", got, tt.want)
			}
		})
	}
}
