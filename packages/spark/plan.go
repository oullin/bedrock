package spark

// Plan represents a subscription plan with fluent configuration.
// Mirrors Spark\Plan.
type Plan struct {
	ID               string
	Name             string
	Interval         string
	TrialDays        int
	Price            float64
	Currency         string
	RawPrice         int
	MonthlyIncentive string
	YearlyIncentive  string
	ShortDescription string
	Features         []string
	Options          map[string]any
	Active           bool
}

// NewPlan creates a new plan with the given name and provider price ID.
func NewPlan(name, id string) *Plan {
	return &Plan{
		ID:       id,
		Name:     name,
		Interval: "monthly",
		Active:   true,
	}
}

// SetInterval sets the plan's billing interval.
func (p *Plan) SetInterval(interval string) *Plan {
	p.Interval = interval

	return p
}

// Monthly sets the plan interval to monthly.
func (p *Plan) Monthly() *Plan {
	p.Interval = "monthly"

	return p
}

// Yearly sets the plan interval to yearly.
func (p *Plan) Yearly() *Plan {
	p.Interval = "yearly"

	return p
}

// SetTrialDays sets the number of trial days for this plan.
func (p *Plan) SetTrialDays(days int) *Plan {
	p.TrialDays = days

	return p
}

// SetIncentive sets the monthly and yearly incentive text.
func (p *Plan) SetIncentive(monthly, yearly string) *Plan {
	p.MonthlyIncentive = monthly
	p.YearlyIncentive = yearly

	return p
}

// SetShortDescription sets the plan's marketing description.
func (p *Plan) SetShortDescription(desc string) *Plan {
	p.ShortDescription = desc

	return p
}

// SetFeatures sets the plan's feature list.
func (p *Plan) SetFeatures(features []string) *Plan {
	p.Features = features

	return p
}

// SetOptions sets the plan's custom options.
func (p *Plan) SetOptions(opts map[string]any) *Plan {
	p.Options = opts

	return p
}

// SetStatus sets whether the plan is active.
func (p *Plan) SetStatus(active bool) *Plan {
	p.Active = active

	return p
}

// Archive marks the plan as inactive (archived).
func (p *Plan) Archive() *Plan {
	p.Active = false

	return p
}

// ToMap serialises the plan for JSON responses.
func (p *Plan) ToMap() map[string]any {
	return map[string]any{
		"id":        p.ID,
		"name":      p.Name,
		"interval":  p.Interval,
		"price":     p.Price,
		"currency":  p.Currency,
		"raw_price": p.RawPrice,
		"incentive": map[string]string{
			"monthly": p.MonthlyIncentive,
			"yearly":  p.YearlyIncentive,
		},
		"short_description": p.ShortDescription,
		"trial_days":        p.TrialDays,
		"features":          p.Features,
		"options":           p.Options,
		"active":            p.Active,
	}
}
