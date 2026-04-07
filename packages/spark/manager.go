package spark

import (
	"fmt"
	"net/http"
	"sync"
)

// Manager is the central registry for billable types, plans, and callbacks.
// It mirrors Spark's SparkManager.
type Manager struct {
	mu         sync.RWMutex
	billables  map[string]*billableEntry
	plans      map[string][]SparkPlan
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

// NewManager creates an empty Spark manager.
func NewManager() *Manager {
	return &Manager{
		billables:  make(map[string]*billableEntry),
		plans:      make(map[string][]SparkPlan),
		prorations: true,
	}
}

// RegisterBillable adds a billable type configuration.
func (m *Manager) RegisterBillable(cfg BillableConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.billables[cfg.ModelName] = &billableEntry{config: cfg}
}

// Billable returns a configuration builder for the given billable type.
func (m *Manager) Billable(billableType string) *BillableConfigBuilder {
	return &BillableConfigBuilder{manager: m, billableType: billableType}
}

// Plans returns the registered plans for a billable type.
func (m *Manager) Plans(billableType string) []SparkPlan {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.plans[billableType]
}

// AddPlan registers a plan for a billable type.
func (m *Manager) AddPlan(billableType string, plan SparkPlan) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.plans[billableType] = append(m.plans[billableType], plan)
}

// ResolveBillable resolves the current billable entity from the request.
func (m *Manager) ResolveBillable(billableType string, r *http.Request) (Billable, error) {
	m.mu.RLock()
	entry, ok := m.billables[billableType]
	m.mu.RUnlock()

	if !ok || entry.resolver == nil {
		return nil, fmt.Errorf("spark: no resolver registered for billable type %q", billableType)
	}

	return entry.resolver(r)
}

// IsAuthorized checks whether the billable is authorized to view the billing
// portal.
func (m *Manager) IsAuthorized(billable Billable, r *http.Request) bool {
	m.mu.RLock()
	entry, ok := m.billables[billable.BillableType()]
	m.mu.RUnlock()

	if !ok || entry.authorizer == nil {
		return true
	}

	return entry.authorizer(billable, r)
}

// EnsurePlanEligibility checks that the billable may subscribe to the plan.
func (m *Manager) EnsurePlanEligibility(billable Billable, plan SparkPlan) error {
	m.mu.RLock()
	entry, ok := m.billables[billable.BillableType()]
	m.mu.RUnlock()

	if !ok || entry.eligibility == nil {
		return nil
	}

	return entry.eligibility(billable, plan)
}

// ChargesPerSeat reports whether the billable type uses per-seat billing.
func (m *Manager) ChargesPerSeat(billableType string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.billables[billableType]

	return ok && entry.seatCount != nil
}

// SeatCount returns the current seat count for the billable.
func (m *Manager) SeatCount(billableType string, billable Billable) int {
	m.mu.RLock()
	entry, ok := m.billables[billableType]
	m.mu.RUnlock()

	if !ok || entry.seatCount == nil {
		return 1
	}

	return entry.seatCount(billable)
}

// SeatName returns the seat display name for the billable type.
func (m *Manager) SeatName(billableType string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.billables[billableType]
	if !ok {
		return ""
	}

	return entry.seatName
}

// DefaultBillableType returns the default billable type. If only one is
// registered it returns that; otherwise it returns "user".
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

// BillableConfig returns the configuration for a billable type.
func (m *Manager) BillableConfig(billableType string) (BillableConfig, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.billables[billableType]
	if !ok {
		return BillableConfig{}, false
	}

	return entry.config, true
}

// SetProrates configures whether plan changes are prorated.
func (m *Manager) SetProrates(v bool) { m.prorations = v }

// Prorates reports whether plan changes are prorated.
func (m *Manager) Prorates() bool { return m.prorations }
