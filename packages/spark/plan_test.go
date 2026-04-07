package spark_test

import (
	"testing"

	"github.com/bedrock/packages/spark"
)

func TestPlanPeriodPriceDisplayAmount(t *testing.T) {
	amount := int64(1999)

	tests := []struct {
		name string
		mode spark.PlanPricingMode
		amt  *int64
		curr string
		want string
	}{
		{"free", spark.PricingModeFree, nil, "", "Free"},
		{"custom", spark.PricingModeCustom, nil, "", "Custom"},
		{"money with amount", spark.PricingModeMoney, &amount, "USD", "1999 USD"},
		{"money nil amount", spark.PricingModeMoney, nil, "USD", "N/A"},
		{"money empty currency", spark.PricingModeMoney, &amount, "", "N/A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &spark.PlanPeriodPrice{
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
