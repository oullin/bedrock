package billing

import (
	"fmt"
	"net/http"
	"sync"
)

// Manager is the central registry for billable types, plans, and
// callbacks. Thread-safe. Mirrors Billing\BillingManager.
type Manager struct {
	mu sync.RWMutex

	resolvers            map[string]ResolverFunc
	authorizers          map[string]AuthorizerFunc
	eligibilityChecks    map[string][]EligibilityFunc
	seatNames            map[string]string
	seatCountCallbacks   map[string]SeatCountFunc
	plans                map[string][]*Plan
	checkoutOptions      map[string]func(Billable) map[string]any
	paymentMethodOptions map[string]func(Billable) map[string]any
	billableConfigs      map[string]BillableConfig
	prorates             bool
}

// NewManager creates an empty Manager.

// RegisterBillable registers a billable type with its configuration.

// Billable returns a fluent builder for configuring a billable type.

// ResolveBillable resolves the billable from a request using the
// registered resolver for the given type.

// IsAuthorized checks whether the billable may view the billing portal.

// EnsurePlanEligibility runs all eligibility checks for the billable's
// type. Returns the first error encountered, or nil.

// AddPlan registers a plan for a billable type.

// Plan creates and registers a new plan for the given billable type.

// Plans returns all plans registered for a billable type.

// ChargesPerSeat reports whether the billable type uses per-seat billing.

// SeatCount returns the number of seats for the billable.

// SeatName returns the human-readable seat name for a billable type.

// SetProrates sets the global proration preference.

// Prorates reports whether subscription changes should be prorated.

// DefaultBillableType returns the first registered billable type, or
// "user" if none are registered.

// BillableConfig returns the config for a billable type.

// SetCheckoutSessionOptions registers a callback that provides custom
// checkout session options for a billable type.

// CheckoutSessionOptions returns the registered checkout options
// callback, or nil.

// SetPaymentMethodSessionOptions registers a callback for payment
// method session options.

// PaymentMethodSessionOptions returns the registered payment method
// options callback, or nil.

// ValidPlan reports whether the given plan ID exists for the billable type.

// BillableConfigBuilder provides a fluent API for configuring a
// billable type on the Manager. Mirrors Billing\BillableConfigurationBuilder.
type BillableConfigBuilder struct {
	manager      *Manager
	billableType string
}

func NewManager() *Manager {
	return &Manager{
		resolvers:            make(map[string]ResolverFunc),
		authorizers:          make(map[string]AuthorizerFunc),
		eligibilityChecks:    make(map[string][]EligibilityFunc),
		seatNames:            make(map[string]string),
		seatCountCallbacks:   make(map[string]SeatCountFunc),
		plans:                make(map[string][]*Plan),
		checkoutOptions:      make(map[string]func(Billable) map[string]any),
		paymentMethodOptions: make(map[string]func(Billable) map[string]any),
		billableConfigs:      make(map[string]BillableConfig),
		prorates:             true,
	}
}

func (m *Manager) RegisterBillable(cfg BillableConfig) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.billableConfigs[cfg.Model] = cfg
}

func (m *Manager) Billable(billableType string) *BillableConfigBuilder {
	return &BillableConfigBuilder{manager: m, billableType: billableType}
}

func (m *Manager) ResolveBillable(billableType string, r *http.Request) (Billable, error) {
	m.mu.RLock()
	resolver, ok := m.resolvers[billableType]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no resolver registered for billable type %q", billableType)
	}

	return resolver(r)
}

func (m *Manager) IsAuthorized(billable Billable, r *http.Request) bool {
	m.mu.RLock()
	auth, ok := m.authorizers[billable.BillableType()]
	m.mu.RUnlock()

	if !ok {
		return true
	}

	return auth(billable, r)
}

func (m *Manager) EnsurePlanEligibility(billable Billable, plan Plan) error {
	m.mu.RLock()
	checks := m.eligibilityChecks[billable.BillableType()]
	m.mu.RUnlock()

	for _, check := range checks {
		if err := check(billable, plan); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) AddPlan(billableType string, plan *Plan) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.plans[billableType] = append(m.plans[billableType], plan)
}

func (m *Manager) Plan(billableType, name, id string) *Plan {
	plan := NewPlan(name, id)
	m.AddPlan(billableType, plan)

	return plan
}

func (m *Manager) Plans(billableType string) []*Plan {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.plans[billableType]
}

func (m *Manager) ChargesPerSeat(billableType string) bool {
	m.mu.RLock()

	defer m.mu.RUnlock()

	_, ok := m.seatCountCallbacks[billableType]

	return ok
}

func (m *Manager) SeatCount(billableType string, billable Billable) int {
	m.mu.RLock()
	cb, ok := m.seatCountCallbacks[billableType]
	m.mu.RUnlock()

	if !ok {
		return 0
	}

	return cb(billable)
}

func (m *Manager) SeatName(billableType string) string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.seatNames[billableType]
}

func (m *Manager) SetProrates(v bool) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.prorates = v
}

func (m *Manager) Prorates() bool {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.prorates
}

func (m *Manager) DefaultBillableType() string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	for k := range m.billableConfigs {
		return k
	}

	return "user"
}

func (m *Manager) BillableConfig(billableType string) (BillableConfig, bool) {
	m.mu.RLock()

	defer m.mu.RUnlock()

	cfg, ok := m.billableConfigs[billableType]

	return cfg, ok
}

func (m *Manager) SetCheckoutSessionOptions(billableType string, fn func(Billable) map[string]any) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.checkoutOptions[billableType] = fn
}

func (m *Manager) CheckoutSessionOptions(billableType string) func(Billable) map[string]any {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.checkoutOptions[billableType]
}

func (m *Manager) SetPaymentMethodSessionOptions(billableType string, fn func(Billable) map[string]any) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.paymentMethodOptions[billableType] = fn
}

func (m *Manager) PaymentMethodSessionOptions(billableType string) func(Billable) map[string]any {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.paymentMethodOptions[billableType]
}

func ValidPlan(manager *Manager, billableType, planID string) bool {
	plans := manager.Plans(billableType)

	for _, p := range plans {
		if p.ID == planID && p.Active {
			return true
		}
	}

	return false
}

// Resolve registers the resolver callback for this billable type.
func (b *BillableConfigBuilder) Resolve(fn ResolverFunc) *BillableConfigBuilder {
	b.manager.mu.Lock()

	defer b.manager.mu.Unlock()

	b.manager.resolvers[b.billableType] = fn

	return b
}

// Authorize registers the authorization callback.
func (b *BillableConfigBuilder) Authorize(fn AuthorizerFunc) *BillableConfigBuilder {
	b.manager.mu.Lock()

	defer b.manager.mu.Unlock()

	b.manager.authorizers[b.billableType] = fn

	return b
}

// CheckPlanEligibility registers a plan eligibility check.
func (b *BillableConfigBuilder) CheckPlanEligibility(fn EligibilityFunc) *BillableConfigBuilder {
	b.manager.mu.Lock()

	defer b.manager.mu.Unlock()

	b.manager.eligibilityChecks[b.billableType] = append(
		b.manager.eligibilityChecks[b.billableType], fn,
	)

	return b
}

// ChargePerSeat configures per-seat billing for this billable type.
func (b *BillableConfigBuilder) ChargePerSeat(seatName string, fn SeatCountFunc) *BillableConfigBuilder {
	b.manager.mu.Lock()

	defer b.manager.mu.Unlock()

	b.manager.seatNames[b.billableType] = seatName
	b.manager.seatCountCallbacks[b.billableType] = fn

	return b
}

// Plan creates and registers a new plan for this billable type.
func (b *BillableConfigBuilder) Plan(name, id string) *Plan {
	return b.manager.Plan(b.billableType, name, id)
}
