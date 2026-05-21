// Package state builds the billing portal frontend state.
package state

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/wayfinder"
)

// FrontendState builds the data shared with the billing portal frontend.
type FrontendState struct {
	manager       *spark.Manager
	config        *spark.Config
	routes        *wayfinder.Registry
	subscriptions spark.SubscriptionStore
	customers     spark.CustomerStore
	transactions  spark.TransactionStore
	payments      spark.PaymentReporter
}

// NewFrontendState creates a FrontendState builder.
func NewFrontendState(mgr *spark.Manager, cfg *spark.Config, subs spark.SubscriptionStore) *FrontendState {
	if cfg == nil {
		cfg = spark.DefaultConfig()
	}

	return &FrontendState{
		manager:       mgr,
		config:        cfg,
		routes:        spark.NewRouteRegistry(),
		subscriptions: subs,
	}
}

// WithRoutes sets the route registry used to generate portal URLs.
func (f *FrontendState) WithRoutes(routes *wayfinder.Registry) *FrontendState {
	if routes != nil {
		f.routes = routes
	}

	return f
}

// WithCustomerStore enables customer-backed portal state, including pending
// checkout detection.
func (f *FrontendState) WithCustomerStore(customers spark.CustomerStore) *FrontendState {
	f.customers = customers

	return f
}

// WithTransactionStore enables invoice list state.
func (f *FrontendState) WithTransactionStore(transactions spark.TransactionStore) *FrontendState {
	f.transactions = transactions

	return f
}

// WithPaymentReporter enables last and next payment state.
func (f *FrontendState) WithPaymentReporter(payments spark.PaymentReporter) *FrontendState {
	f.payments = payments

	return f
}

// Current returns the full frontend state for the billing portal.
func (f *FrontendState) Current(ctx context.Context, billableType string, billable spark.Billable) (map[string]any, error) {
	return f.CurrentAt(ctx, billableType, billable, time.Now())
}

// CurrentAt returns the full frontend state using an explicit clock. It is
// useful for billing DTOs that expose relative pending-window values.
func (f *FrontendState) CurrentAt(ctx context.Context, billableType string, billable spark.Billable, now time.Time) (map[string]any, error) {
	sub, _ := f.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())
	customer, err := f.customerForBillable(ctx, billable)

	if err != nil {
		return nil, err
	}

	plans := f.manager.Plans(billableType)

	var monthlyPlans, yearlyPlans []*spark.Plan

	for _, p := range plans {
		if !p.Active {
			continue
		}

		switch p.Interval {
		case "monthly":
			monthlyPlans = append(monthlyPlans, p)
		case "yearly":
			yearlyPlans = append(yearlyPlans, p)
		}
	}

	activePlan := planForSubscription(plans, sub)

	state, err := f.resolveState(ctx, sub, customer)

	if err != nil {
		return nil, err
	}

	subscription := f.subscriptionState(sub, activePlan)
	cta := f.ctaState(sub, activePlan, now)
	invoices, err := f.invoicesForBillable(ctx, billable)

	if err != nil {
		return nil, err
	}

	lastPayment, nextPayment, err := f.paymentState(ctx, sub)

	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"appLogo":            f.config.BrandLogo(),
		"appName":            f.appName(),
		"sandbox":            f.config.Sandbox(),
		"billableId":         billable.BillableID(),
		"billableName":       billable.BillableName(),
		"billableType":       billableType,
		"brandColor":         f.brandColor(),
		"clientSideToken":    f.config.ClientSideToken(),
		"dashboardUrl":       f.dashboardURL(),
		"defaultInterval":    f.defaultInterval(billableType),
		"genericTrialEndsAt": f.genericTrialEndsAt(customer),
		"invoices":           invoices,
		"lastPayment":        lastPayment,
		"message":            "",
		"monthlyPlans":       plansToMaps(monthlyPlans),
		"nextPayment":        nextPayment,
		"paddleSellerId":     f.config.SellerID(),
		"yearlyPlans":        plansToMaps(yearlyPlans),
		"plan":               planToMap(activePlan),
		"pwAuth":             f.config.RetainKey(),
		"pwCustomer":         providerCustomerID(customer),
		"seatName":           f.manager.SeatName(billableType),
		"sparkPath":          f.config.Path(),
		"state":              state,
		"subscription":       subscription,
		"cta":                cta,
		"termsUrl":           f.config.TermsURL(),
	}

	return data, nil
}

func (f *FrontendState) customerForBillable(ctx context.Context, billable spark.Billable) (*spark.Customer, error) {
	if f.customers == nil {
		return nil, nil
	}

	return f.customers.FindByBillable(ctx, billable.BillableType(), billable.BillableID())
}

func (f *FrontendState) invoicesForBillable(ctx context.Context, billable spark.Billable) ([]map[string]any, error) {
	if f.transactions == nil {
		return []map[string]any{}, nil
	}

	transactions, err := f.transactions.FindByBillable(ctx, billable.BillableType(), billable.BillableID(), 10)

	if err != nil {
		return nil, err
	}

	invoices := make([]map[string]any, 0, len(transactions))

	for _, transaction := range transactions {
		if transaction.Total <= 0 {
			continue
		}

		invoice := map[string]any{
			"id":          invoiceID(transaction),
			"total":       transaction.TotalFormatted(),
			"invoice_url": f.invoiceURL(billable, transaction),
		}

		if transaction.BilledAt != nil {
			invoice["billed_at"] = transaction.BilledAt.Format(f.dateFormat())
		}

		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func invoiceID(transaction spark.Transaction) string {
	if transaction.PaddleID != "" {
		return transaction.PaddleID
	}

	return transaction.InvoiceNumber
}

func (f *FrontendState) paymentState(ctx context.Context, sub *spark.Subscription) (map[string]any, map[string]any, error) {
	if f.payments == nil || !subscriptionIsActiveOrPastDue(sub) {
		return nil, nil, nil
	}

	last, err := f.payments.LastPayment(ctx, sub)

	if err != nil {
		return nil, nil, err
	}

	next, err := f.payments.NextPayment(ctx, sub)

	if err != nil {
		return nil, nil, err
	}

	return paymentToMap(last, f.dateFormat()), paymentToMap(next, f.dateFormat()), nil
}

func plansToMaps(plans []*spark.Plan) []map[string]any {
	payload := make([]map[string]any, 0, len(plans))

	for _, plan := range plans {
		if plan == nil {
			continue
		}

		payload = append(payload, plan.ToMap())
	}

	return payload
}

func planToMap(plan *spark.Plan) map[string]any {
	if plan == nil {
		return nil
	}

	return plan.ToMap()
}

func paymentToMap(payment *spark.Payment, dateFormat string) map[string]any {
	if payment == nil {
		return nil
	}

	return map[string]any{
		"amount":   payment.FormattedAmount(),
		"currency": payment.Currency,
		"date":     payment.Date.Format(dateFormat),
	}
}

func planForSubscription(plans []*spark.Plan, sub *spark.Subscription) *spark.Plan {
	if sub == nil {
		return nil
	}

	for _, p := range plans {
		if sub.HasPrice(p.ID) {
			return p
		}
	}

	for _, p := range plans {
		if slug, _ := p.Options["slug"].(string); slug != "" && slug == sub.Plan {
			return p
		}

		if p.Name == sub.Plan {
			return p
		}
	}

	return nil
}

func (f *FrontendState) resolveState(ctx context.Context, sub *spark.Subscription, customer *spark.Customer) (string, error) {
	if subscriptionIsActiveOrPastDue(sub) {
		if customer != nil && customer.PendingCheckoutID != "" {
			customer.PendingCheckoutID = ""

			if err := f.customers.Save(ctx, customer); err != nil {
				return "", err
			}
		}
	} else if customer != nil && customer.PendingCheckoutID != "" {
		return "pending", nil
	}

	if sub == nil {
		return "none", nil
	}

	if sub.OnGracePeriod() {
		return "onGracePeriod", nil
	}

	if sub.Active() || sub.OnTrial() {
		return "active", nil
	}

	if sub.PastDue() {
		return "past_due", nil
	}

	return "none", nil
}

func subscriptionIsActiveOrPastDue(sub *spark.Subscription) bool {
	if sub == nil {
		return false
	}

	return sub.Active() || sub.OnTrial() || sub.PastDue()
}

func (f *FrontendState) brandColor() string {
	if f.config.BrandColor() != "" {
		return f.config.BrandColor()
	}

	return "bg-gray-800"
}

func (f *FrontendState) appName() string {
	if f.config.AppName() != "" {
		return f.config.AppName()
	}

	return "Laravel"
}

func (f *FrontendState) dashboardURL() string {
	if f.config.DashboardURL() != "" {
		return f.config.DashboardURL()
	}

	return "/"
}

func (f *FrontendState) dateFormat() string {
	if f.config.DateFormat() != "" {
		return f.config.DateFormat()
	}

	return "January 2, 2006"
}

func (f *FrontendState) genericTrialEndsAt(customer *spark.Customer) any {
	if customer == nil || !customer.OnGenericTrial() || customer.TrialEndsAt == nil {
		return nil
	}

	return customer.TrialEndsAt.Format(f.dateFormat())
}

func providerCustomerID(customer *spark.Customer) any {
	if customer == nil || customer.PaddleID == "" {
		return nil
	}

	return customer.PaddleID
}

func (f *FrontendState) invoiceURL(billable spark.Billable, transaction spark.Transaction) string {
	return spark.InvoiceDownloadURL(f.routes, billable, transaction)
}

func (f *FrontendState) defaultInterval(billableType string) string {
	if billables := f.config.Billables(); billables != nil {
		if cfg, ok := billables[billableType]; ok && cfg.DefaultInterval != "" {
			return cfg.DefaultInterval
		}
	}

	return "monthly"
}

func (f *FrontendState) subscriptionState(sub *spark.Subscription, plan *spark.Plan) map[string]any {
	state := map[string]any{
		"status":             "",
		"plan_code":          "",
		"plan_name":          "",
		"pending_expires_at": (*time.Time)(nil),
		"payment_ready_at":   (*time.Time)(nil),
		"portal_url":         "",
		"pay_now":            false,
	}

	if sub == nil {
		return state
	}

	state["status"] = string(sub.Status)
	state["plan_code"] = sub.Plan

	if plan != nil {
		state["plan_name"] = plan.Name
	}

	state["pending_expires_at"] = sub.PendingExpiresAt
	state["payment_ready_at"] = sub.PaymentReadyAt

	portalURL := f.portalURL(sub, plan)
	state["portal_url"] = portalURL
	state["pay_now"] = sub.Status == spark.StatusAwaitingPayment && portalURL != ""

	return state
}

func (f *FrontendState) ctaState(sub *spark.Subscription, plan *spark.Plan, now time.Time) map[string]any {
	cta := map[string]any{
		"visible":        false,
		"label":          "",
		"remaining_days": 0,
	}

	if sub == nil || plan == nil || sub.Status != spark.StatusAwaitingPayment {
		return cta
	}

	portalURL := f.portalURL(sub, plan)

	if portalURL == "" {
		return cta
	}

	cta["visible"] = true
	cta["label"] = "Complete Subscription"
	cta["remaining_days"] = remainingDays(now, sub.PendingExpiresAt)

	return cta
}

func (f *FrontendState) portalURL(sub *spark.Subscription, plan *spark.Plan) string {
	if sub == nil || plan == nil {
		return ""
	}

	switch sub.Status {
	case spark.StatusAwaitingPayment:
		if !sub.HasPrice(plan.ID) {
			return ""
		}
	case spark.StatusActive, spark.StatusPastDue, spark.StatusPaused:
	default:
		return ""
	}

	path := strings.Trim(f.config.Path(), "/")

	if path == "" {
		return "/"
	}

	return "/" + path
}

func remainingDays(now time.Time, expiresAt *time.Time) int {
	if expiresAt == nil || !expiresAt.After(now) {
		return 0
	}

	return int(math.Ceil(expiresAt.Sub(now).Hours() / 24))
}
