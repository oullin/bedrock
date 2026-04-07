package spark_test

import (
	"testing"

	"github.com/bedrock/packages/spark"
)

func TestSubscriptionPlanValid(t *testing.T) {
	tests := []struct {
		plan spark.SubscriptionPlan
		want bool
	}{
		{spark.PlanStarter, true},
		{spark.PlanPro, true},
		{spark.PlanEnterprise, true},
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
	all := []spark.SubscriptionStatus{
		spark.StatusPending, spark.StatusAwaitingPayment, spark.StatusActive,
		spark.StatusTrialing, spark.StatusPastDue, spark.StatusPaused,
		spark.StatusCanceled, spark.StatusExpired,
	}

	for _, s := range all {
		if !s.Valid() {
			t.Errorf("SubscriptionStatus(%q).Valid() = false, want true", s)
		}
	}

	if spark.SubscriptionStatus("bogus").Valid() {
		t.Error("SubscriptionStatus(bogus).Valid() = true, want false")
	}
}

func TestSubscriptionStatusGrantsAccess(t *testing.T) {
	grant := []spark.SubscriptionStatus{spark.StatusActive, spark.StatusTrialing, spark.StatusPastDue}
	deny := []spark.SubscriptionStatus{
		spark.StatusPending, spark.StatusAwaitingPayment,
		spark.StatusPaused, spark.StatusCanceled, spark.StatusExpired,
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
	if !spark.StatusCanceled.IsTerminal() {
		t.Error("StatusCanceled.IsTerminal() = false, want true")
	}

	if !spark.StatusExpired.IsTerminal() {
		t.Error("StatusExpired.IsTerminal() = false, want true")
	}

	if spark.StatusActive.IsTerminal() {
		t.Error("StatusActive.IsTerminal() = true, want false")
	}
}

func TestBillingPeriodDisplayLabel(t *testing.T) {
	tests := []struct {
		period spark.BillingPeriod
		want   string
	}{
		{spark.PeriodMonthly, "Billed monthly"},
		{spark.PeriodYearly, "Billed annually"},
		{spark.PeriodFree, "Free"},
		{spark.PeriodCustom, "Custom"},
	}

	for _, tt := range tests {
		if got := tt.period.DisplayLabel(); got != tt.want {
			t.Errorf("BillingPeriod(%q).DisplayLabel() = %q, want %q", tt.period, got, tt.want)
		}
	}
}

func TestPlanPricingModeDisplayLabel(t *testing.T) {
	tests := []struct {
		mode spark.PlanPricingMode
		want string
	}{
		{spark.PricingModeMoney, "Money"},
		{spark.PricingModeFree, "Free"},
		{spark.PricingModeCustom, "Custom"},
	}

	for _, tt := range tests {
		if got := tt.mode.DisplayLabel(); got != tt.want {
			t.Errorf("PlanPricingMode(%q).DisplayLabel() = %q, want %q", tt.mode, got, tt.want)
		}
	}
}

func TestSubscriptionFeatureCodeValid(t *testing.T) {
	valid := []spark.SubscriptionFeatureCode{
		spark.FeatureActiveRoleCapacity, spark.FeatureBookingDispatch,
		spark.FeatureCalendarIntegrations, spark.FeatureViewLLMAnalysis,
	}

	for _, c := range valid {
		if !c.Valid() {
			t.Errorf("SubscriptionFeatureCode(%q).Valid() = false, want true", c)
		}
	}

	if spark.SubscriptionFeatureCode("nope").Valid() {
		t.Error("SubscriptionFeatureCode(nope).Valid() = true, want false")
	}
}

func TestProrationBehaviorValid(t *testing.T) {
	valid := []spark.ProrationBehavior{
		spark.ProratedNextBillingPeriod, spark.FullNextBillingPeriod,
		spark.ProratedImmediately, spark.FullImmediately, spark.DoNotBill,
	}

	for _, b := range valid {
		if !b.Valid() {
			t.Errorf("ProrationBehavior(%q).Valid() = false, want true", b)
		}
	}

	if spark.ProrationBehavior("invalid").Valid() {
		t.Error("ProrationBehavior(invalid).Valid() = true, want false")
	}
}

func TestTransactionStatusValid(t *testing.T) {
	valid := []spark.TransactionStatus{
		spark.TransactionDraft, spark.TransactionReady, spark.TransactionBilled,
		spark.TransactionPaid, spark.TransactionCompleted, spark.TransactionCanceled,
		spark.TransactionPastDue,
	}

	for _, s := range valid {
		if !s.Valid() {
			t.Errorf("TransactionStatus(%q).Valid() = false, want true", s)
		}
	}

	if spark.TransactionStatus("fake").Valid() {
		t.Error("TransactionStatus(fake).Valid() = true, want false")
	}
}
