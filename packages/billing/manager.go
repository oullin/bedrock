package billing

import (
	"fmt"
	"net/http"
	"sync"
)

// Manager is the central registry for billable types, plans, and callbacks.
// It mirrors Billing's BillingManager.
type Manager struct {
	mu         sync.RWMutex
	billables  map[string]*billableEntry
	plans      map[string][]BillingPlan
	prorations bool
}

type billableEntry struct {
	config      BillableConfig
	resolver    ResolverFunc
	authorizer  AuthorizerFunc
	eligibility EligibilityFunc
	seatName    string
	seatCount   SeatCountFunc
}

// NewManager creates an empty Billing manager.

// RegisterBillable adds a billable type configuration.

// Billable returns a configuration builder for the given billable type.

// Plans returns the registered plans for a billable type.

// AddPlan registers a plan for a billable type.

// ResolveBillable resolves the current billable entity from the request.

// IsAuthorized checks whether the billable is authorized to view the billing
// portal.

// EnsurePlanEligibility checks that the billable may subscribe to the plan.

// ChargesPerSeat reports whether the billable type uses per-seat billing.

// SeatCount returns the current seat count for the billable.

// SeatName returns the seat display name for the billable type.

// DefaultBillableType returns the default billable type. If only one is
// registered it returns that; otherwise it returns "user".

// BillableConfig returns the configuration for a billable type.

// SetProrates configures whether plan changes are prorated.

// Prorates reports whether plan changes are prorated.

// ValidPlan checks whether the given plan name exists and is active in the
// manager for the specified billable type.

// BillableConfigBuilder provides a fluent API for configuring a billable type.
type BillableConfigBuilder struct {
	manager      *Manager
	billableType string
}

func NewManager() *Manager {
	return &Manager{
		billables:  make(map[string]*billableEntry),
		plans:      make(map[string][]BillingPlan),
		prorations: true,
	}
}

func (m *Manager) RegisterBillable(cfg BillableConfig) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.billables[cfg.ModelName] = &billableEntry{config: cfg}
}

func (m *Manager) Billable(billableType string) *BillableConfigBuilder {
	return &BillableConfigBuilder{manager: m, billableType: billableType}
}

func (m *Manager) Plans(billableType string) []BillingPlan {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.plans[billableType]
}

func (m *Manager) AddPlan(billableType string, plan BillingPlan) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.plans[billableType] = append(m.plans[billableType], plan)
}

func (m *Manager) ResolveBillable(billableType string, r *http.Request) (Billable, error) {
	m.mu.RLock()
	entry, ok := m.billables[billableType]
	m.mu.RUnlock()

	if !ok || entry.resolver == nil {
		return nil, fmt.Errorf("billing: no resolver registered for billable type %q", billableType)
	}

	return entry.resolver(r)
}

func (m *Manager) IsAuthorized(billable Billable, r *http.Request) bool {
	m.mu.RLock()
	entry, ok := m.billables[billable.BillableType()]
	m.mu.RUnlock()

	if !ok || entry.authorizer == nil {
		return true
	}

	return entry.authorizer(billable, r)
}

func (m *Manager) EnsurePlanEligibility(billable Billable, plan BillingPlan) error {
	m.mu.RLock()
	entry, ok := m.billables[billable.BillableType()]
	m.mu.RUnlock()

	if !ok || entry.eligibility == nil {
		return nil
	}

	return entry.eligibility(billable, plan)
}

func (m *Manager) ChargesPerSeat(billableType string) bool {
	m.mu.RLock()

	defer m.mu.RUnlock()

	entry, ok := m.billables[billableType]

	return ok && entry.seatCount != nil
}

func (m *Manager) SeatCount(billableType string, billable Billable) int {
	m.mu.RLock()
	entry, ok := m.billables[billableType]
	m.mu.RUnlock()

	if !ok || entry.seatCount == nil {
		return 1
	}

	return entry.seatCount(billable)
}

func (m *Manager) SeatName(billableType string) string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	entry, ok := m.billables[billableType]

	if !ok {
		return ""
	}

	return entry.seatName
}

func (m *Manager) DefaultBillableType() string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	if len(m.billables) == 1 {
		for k := range m.billables {
			return k
		}
	}

	return "user"
}

func (m *Manager) BillableConfig(billableType string) (BillableConfig, bool) {
	m.mu.RLock()

	defer m.mu.RUnlock()

	entry, ok := m.billables[billableType]

	if !ok {
		return BillableConfig{}, false
	}

	return entry.config, true
}

func (m *Manager) SetProrates(v bool) { m.prorations = v }

func (m *Manager) Prorates() bool { return m.prorations }

func ValidPlan(manager *Manager, billableType string, planID string) bool {
	for _, p := range manager.Plans(billableType) {
		if p.ID == planID && p.Active {
			return true
		}
	}

	return false
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
