package spark_test

import (
	"testing"

	"github.com/bedrock/packages/spark"
)

func TestCheckoutInputValidate(t *testing.T) {
	valid := &spark.CheckoutInput{Plan: spark.PlanPro, Period: spark.PeriodMonthly}
	if errs := valid.Validate(); errs != nil {
		t.Errorf("valid input returned errors: %v", errs)
	}

	invalid := &spark.CheckoutInput{Plan: "bogus", Period: "nope"}
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
	valid := &spark.CustomPlanInquiryInput{
		Name:    "Jane",
		Email:   "jane@example.com",
		Message: "I need a custom plan",
	}

	if errs := valid.Validate(); errs != nil {
		t.Errorf("valid input returned errors: %v", errs)
	}

	empty := &spark.CustomPlanInquiryInput{}
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
	m := spark.NewManager()
	m.AddPlan("team", spark.SparkPlan{ID: "pri_pro", Name: "Pro", Active: true})
	m.AddPlan("team", spark.SparkPlan{ID: "pri_archived", Name: "Old", Active: false})

	if !spark.ValidPlan(m, "team", "pri_pro") {
		t.Error("ValidPlan should return true for active plan")
	}

	if spark.ValidPlan(m, "team", "pri_archived") {
		t.Error("ValidPlan should return false for archived plan")
	}

	if spark.ValidPlan(m, "team", "pri_missing") {
		t.Error("ValidPlan should return false for missing plan")
	}
}
