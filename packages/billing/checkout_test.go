package billing

import (
	"encoding/json"
	"testing"
)

func TestCheckoutItemJSONUsesPaddleKeys(t *testing.T) {
	payload, err := json.Marshal(GuestCheckout([]CheckoutItem{
		{PriceID: "pri_123", Quantity: 2},
	}).ToMap())

	if err != nil {
		t.Fatalf("marshal checkout: %v", err)
	}

	var decoded struct {
		Items []map[string]any `json:"items"`
	}

	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal checkout: %v", err)
	}

	if len(decoded.Items) != 1 {
		t.Fatalf("items = %#v, want one item", decoded.Items)
	}

	if decoded.Items[0]["priceId"] != "pri_123" {
		t.Fatalf("priceId = %v, want pri_123", decoded.Items[0]["priceId"])
	}

	if decoded.Items[0]["quantity"] != float64(2) {
		t.Fatalf("quantity = %v, want 2", decoded.Items[0]["quantity"])
	}

	if _, ok := decoded.Items[0]["PriceID"]; ok {
		t.Fatalf("checkout item used legacy PriceID key: %s", payload)
	}
}
