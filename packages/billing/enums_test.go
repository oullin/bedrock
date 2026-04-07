package spark_test

import (
	"testing"

	"github.com/bedrock/packages/billing"
)

func TestSubscriptionPlanValid(t *testing.T) {
	tests := []struct {
		plan billing.SubscriptionPlan
		want bool
	}{
		{billing.PlanStarter, true},
		{billing.PlanPro, true},
		{billing.PlanEnterprise, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := tt.plan.Valid(); got != tt.want {
			t.Errorf("SubscriptionPlan(%q).Valid() = %v, want %v", tt.plan, got, tt.want)
		}
	}
}

func TestSubscriptionStatusValid(t *testing.T) {
	all := []billing.SubscriptionStatus{
		billing.StatusPending, billing.StatusAwaitingPayment, billing.StatusActive,
		billing.StatusTrialing, billing.StatusPastDue, billing.StatusPaused,
		billing.StatusCanceled, billing.StatusExpired,
	}

	for _, s := range all {
		if !s.Valid() {
			t.Errorf("SubscriptionStatus(%q).Valid() = false, want true", s)
		}
	}

	if billing.SubscriptionStatus("bogus").Valid() {
		t.Error("SubscriptionStatus(bogus).Valid() = true, want false")
	}
}

func TestSubscriptionStatusGrantsAccess(t *testing.T) {
	grant := []billing.SubscriptionStatus{billing.StatusActive, billing.StatusTrialing, billing.StatusPastDue}
	deny := []billing.SubscriptionStatus{
		billing.StatusPending, billing.StatusAwaitingPayment,
		billing.StatusPaused, billing.StatusCanceled, billing.StatusExpired,
	}

	for _, s := range grant {
		if !s.GrantsAccess() {
			t.Errorf("SubscriptionStatus(%q).GrantsAccess() = false, want true", s)
		}
	}

	for _, s := range deny {
		if s.GrantsAccess() {
			t.Errorf("SubscriptionStatus(%q).GrantsAccess() = true, want false", s)
		}
	}
}

func TestSubscriptionStatusIsTerminal(t *testing.T) {
	if !billing.StatusCanceled.IsTerminal() {
		t.Error("StatusCanceled.IsTerminal() = false, want true")
	}

	if !billing.StatusExpired.IsTerminal() {
		t.Error("StatusExpired.IsTerminal() = false, want true")
	}

	if billing.StatusActive.IsTerminal() {
		t.Error("StatusActive.IsTerminal() = true, want false")
	}
}

func TestBillingPeriodDisplayLabel(t *testing.T) {
	tests := []struct {
		period billing.BillingPeriod
		want   string
	}{
		{billing.PeriodMonthly, "Billed monthly"},
		{billing.PeriodYearly, "Billed annually"},
		{billing.PeriodFree, "Free"},
		{billing.PeriodCustom, "Custom"},
	}

	for _, tt := range tests {
		if got := tt.period.DisplayLabel(); got != tt.want {
			t.Errorf("BillingPeriod(%q).DisplayLabel() = %q, want %q", tt.period, got, tt.want)
		}
	}
}

func TestPlanPricingModeDisplayLabel(t *testing.T) {
	tests := []struct {
		mode billing.PlanPricingMode
		want string
	}{
		{billing.PricingModeMoney, "Money"},
		{billing.PricingModeFree, "Free"},
		{billing.PricingModeCustom, "Custom"},
	}

	for _, tt := range tests {
		if got := tt.mode.DisplayLabel(); got != tt.want {
			t.Errorf("PlanPricingMode(%q).DisplayLabel() = %q, want %q", tt.mode, got, tt.want)
		}
	}
}

func TestSubscriptionFeatureCodeValid(t *testing.T) {
	valid := []billing.SubscriptionFeatureCode{
		billing.FeatureActiveRoleCapacity, billing.FeatureBookingDispatch,
		billing.FeatureCalendarIntegrations, billing.FeatureViewLLMAnalysis,
	}

	for _, c := range valid {
		if !c.Valid() {
			t.Errorf("SubscriptionFeatureCode(%q).Valid() = false, want true", c)
		}
	}

	if billing.SubscriptionFeatureCode("nope").Valid() {
		t.Error("SubscriptionFeatureCode(nope).Valid() = true, want false")
	}
}

func TestProrationBehaviorValid(t *testing.T) {
	valid := []billing.ProrationBehavior{
		billing.ProratedNextBillingPeriod, billing.FullNextBillingPeriod,
		billing.ProratedImmediately, billing.FullImmediately, billing.DoNotBill,
	}

	for _, b := range valid {
		if !b.Valid() {
			t.Errorf("ProrationBehavior(%q).Valid() = false, want true", b)
		}
	}

	if billing.ProrationBehavior("invalid").Valid() {
		t.Error("ProrationBehavior(invalid).Valid() = true, want false")
	}
}

func TestTransactionStatusValid(t *testing.T) {
	valid := []billing.TransactionStatus{
		billing.TransactionDraft, billing.TransactionReady, billing.TransactionBilled,
		billing.TransactionPaid, billing.TransactionCompleted, billing.TransactionCanceled,
		billing.TransactionPastDue,
	}

	for _, s := range valid {
		if !s.Valid() {
			t.Errorf("TransactionStatus(%q).Valid() = false, want true", s)
		}
	}

	if billing.TransactionStatus("fake").Valid() {
		t.Error("TransactionStatus(fake).Valid() = true, want false")
	}
}
