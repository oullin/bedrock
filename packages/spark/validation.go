package spark

import "strings"

// CheckoutInput holds the validated input for a checkout request.
type CheckoutInput struct {
	Plan   SubscriptionPlan
	Period BillingPeriod
}

// Validate checks that the checkout input fields are valid.
func (c *CheckoutInput) Validate() ValidationErrors {
	var errs ValidationErrors

	if !c.Plan.Valid() {
		errs = append(errs, ValidationError{Field: "plan", Message: "The selected plan is not available."})
	}

	if !c.Period.Valid() {
		errs = append(errs, ValidationError{Field: "period", Message: "The selected billing period is not available."})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// CustomPlanInquiryInput holds the validated input for a custom plan inquiry.
type CustomPlanInquiryInput struct {
	Name    string
	Email   string
	Company string
	Message string
}

// Validate checks that the inquiry input fields are valid.
func (c *CustomPlanInquiryInput) Validate() ValidationErrors {
	var errs ValidationErrors

	if strings.TrimSpace(c.Name) == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "Your name is required."})
	} else if len(c.Name) > 150 {
		errs = append(errs, ValidationError{Field: "name", Message: "Your name must not be greater than 150 characters."})
	}

	if strings.TrimSpace(c.Email) == "" {
		errs = append(errs, ValidationError{Field: "email", Message: "An email address is required."})
	} else if len(c.Email) > 255 || !strings.Contains(c.Email, "@") {
		errs = append(errs, ValidationError{Field: "email", Message: "Please provide a valid email address."})
	}

	if len(c.Company) > 150 {
		errs = append(errs, ValidationError{Field: "company", Message: "Company name must not be greater than 150 characters."})
	}

	if strings.TrimSpace(c.Message) == "" {
		errs = append(errs, ValidationError{Field: "message", Message: "A message is required."})
	} else if len(c.Message) > 2000 {
		errs = append(errs, ValidationError{Field: "message", Message: "A message must not be greater than 2000 characters."})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// ValidPlan checks whether the given plan name exists and is active in the
// manager for the specified billable type.
func ValidPlan(manager *Manager, billableType string, planID string) bool {
	for _, p := range manager.Plans(billableType) {
		if p.ID == planID && p.Active {
			return true
		}
	}

	return false
}
