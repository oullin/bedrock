package spark

// BillableConfigBuilder provides a fluent API for configuring a billable type.
type BillableConfigBuilder struct {
	manager      *Manager
	billableType string
}

// Resolve registers a callback to resolve the billable from an HTTP request.
func (b *BillableConfigBuilder) Resolve(fn ResolverFunc) *BillableConfigBuilder {
	b.manager.mu.Lock()
	defer b.manager.mu.Unlock()

	entry := b.ensureEntry()
	entry.resolver = fn

	return b
}

// Authorize registers a callback that checks portal access.
func (b *BillableConfigBuilder) Authorize(fn AuthorizerFunc) *BillableConfigBuilder {
	b.manager.mu.Lock()
	defer b.manager.mu.Unlock()

	entry := b.ensureEntry()
	entry.authorizer = fn

	return b
}

// CheckPlanEligibility registers a callback that validates plan eligibility.
func (b *BillableConfigBuilder) CheckPlanEligibility(fn EligibilityFunc) *BillableConfigBuilder {
	b.manager.mu.Lock()
	defer b.manager.mu.Unlock()

	entry := b.ensureEntry()
	entry.eligibility = fn

	return b
}

// ChargePerSeat configures per-seat billing for this billable type.
func (b *BillableConfigBuilder) ChargePerSeat(name string, fn SeatCountFunc) *BillableConfigBuilder {
	b.manager.mu.Lock()
	defer b.manager.mu.Unlock()

	entry := b.ensureEntry()
	entry.seatName = name
	entry.seatCount = fn

	return b
}

// ensureEntry returns the existing entry or creates a new one. Caller must
// hold manager.mu.
func (b *BillableConfigBuilder) ensureEntry() *billableEntry {
	entry, ok := b.manager.billables[b.billableType]
	if !ok {
		entry = &billableEntry{config: BillableConfig{ModelName: b.billableType}}
		b.manager.billables[b.billableType] = entry
	}

	return entry
}
