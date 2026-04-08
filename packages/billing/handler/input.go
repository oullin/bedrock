package handler

import "github.com/bedrock/packages/billing"

// CheckoutInput holds the validated input for a checkout request.
type CheckoutInput struct {
	Plan   billing.SubscriptionPlan
	Period billing.BillingPeriod
}

// Validate checks that the checkout input fields are valid.
func (c *CheckoutInput) Validate() billing.ValidationErrors {
	var errs billing.ValidationErrors

	if !c.Plan.Valid() {
		errs = append(errs, billing.ValidationError{Field: "plan", Message: "The selected plan is not available."})
	}

	if !c.Period.Valid() {
		errs = append(errs, billing.ValidationError{Field: "period", Message: "The selected billing period is not available."})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
