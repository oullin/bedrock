package handler

import "github.com/bedrock/packages/spark"

// CheckoutInput holds the validated input for a checkout request.
type CheckoutInput struct {
	Plan   spark.SubscriptionPlan
	Period spark.BillingPeriod
}

// Validate checks that the checkout input fields are valid.
func (c *CheckoutInput) Validate() spark.ValidationErrors {
	var errs spark.ValidationErrors

	if !c.Plan.Valid() {
		errs = append(errs, spark.ValidationError{Field: "plan", Message: "The selected plan is not available."})
	}

	if !c.Period.Valid() {
		errs = append(errs, spark.ValidationError{Field: "period", Message: "The selected billing period is not available."})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
