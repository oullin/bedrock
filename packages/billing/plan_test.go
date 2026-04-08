package spark_test

import (
	"testing"

	"github.com/bedrock/packages/billing"
)

func TestPlanPeriodPriceDisplayAmount(t *testing.T) {
	amount := int64(1999)

	tests := []struct {
		name string
		mode billing.PlanPricingMode
		amt  *int64
		curr string
		want string
	}{
		{"free", billing.PricingModeFree, nil, "", "Free"},
		{"custom", billing.PricingModeCustom, nil, "", "Custom"},
		{"money with amount", billing.PricingModeMoney, &amount, "USD", "1999 USD"},
		{"money nil amount", billing.PricingModeMoney, nil, "USD", "N/A"},
		{"money empty currency", billing.PricingModeMoney, &amount, "", "N/A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &billing.PlanPeriodPrice{
				PricingMode: tt.mode,
				AmountMinor: tt.amt,
				Currency:    tt.curr,
			}

			if got := p.DisplayAmount(nil); got != tt.want {
				t.Errorf("DisplayAmount() = %q, want %q", got, tt.want)
			}
		})
	}
}
