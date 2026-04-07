package spark_test

import (
	"testing"

	"github.com/bedrock/packages/billing"
)

func TestCheckoutInputValidate(t *testing.T) {
	valid := &billing.CheckoutInput{Plan: billing.PlanPro, Period: billing.PeriodMonthly}
	if errs := valid.Validate(); errs != nil {
		t.Errorf("valid input returned errors: %v", errs)
	}

	invalid := &billing.CheckoutInput{Plan: "bogus", Period: "nope"}
	errs := invalid.Validate()

	if errs == nil {
		t.Fatal("expected validation errors for invalid input")
	}

	if !errs.HasField("plan") {
		t.Error("expected error on field 'plan'")
	}

	if !errs.HasField("period") {
		t.Error("expected error on field 'period'")
	}
}

func TestCustomPlanInquiryInputValidate(t *testing.T) {
	valid := &billing.CustomPlanInquiryInput{
		Name:    "Jane",
		Email:   "jane@example.com",
		Message: "I need a custom plan",
	}

	if errs := valid.Validate(); errs != nil {
		t.Errorf("valid input returned errors: %v", errs)
	}

	empty := &billing.CustomPlanInquiryInput{}
	errs := empty.Validate()

	if errs == nil {
		t.Fatal("expected validation errors for empty input")
	}

	if !errs.HasField("name") {
		t.Error("expected error on field 'name'")
	}

	if !errs.HasField("email") {
		t.Error("expected error on field 'email'")
	}

	if !errs.HasField("message") {
		t.Error("expected error on field 'message'")
	}
}

func TestValidPlan(t *testing.T) {
	m := billing.NewManager()
	m.AddPlan("team", billing.BillingPlan{ID: "pri_pro", Name: "Pro", Active: true})
	m.AddPlan("team", billing.BillingPlan{ID: "pri_archived", Name: "Old", Active: false})

	if !billing.ValidPlan(m, "team", "pri_pro") {
		t.Error("ValidPlan should return true for active plan")
	}

	if billing.ValidPlan(m, "team", "pri_archived") {
		t.Error("ValidPlan should return false for archived plan")
	}

	if billing.ValidPlan(m, "team", "pri_missing") {
		t.Error("ValidPlan should return false for missing plan")
	}
}
