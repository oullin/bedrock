package billing

import (
	"context"
	"time"
)

// BootstrapCustomer creates or updates the provider customer record for a
// newly-created billable, including Billing's generic trial window.
func BootstrapCustomer(
	ctx context.Context,
	store CustomerStore,
	provider ProviderOperations,
	billable Billable,
	trialDays int,
	now time.Time,
) (*Customer, error) {
	customer, err := store.FindByBillable(ctx, billable.BillableType(), billable.BillableID())

	if err != nil {
		return nil, err
	}

	var trialEndsAt *time.Time

	if trialDays > 0 {
		trial := now.AddDate(0, 0, trialDays)
		trialEndsAt = &trial
	}

	if customer == nil {
		customer, err = provider.CreateCustomer(ctx, billable, CustomerCreateOptions{TrialEndsAt: trialEndsAt})

		if err != nil {
			return nil, err
		}

		customer.BillableType = billable.BillableType()
		customer.BillableID = billable.BillableID()
		customer.Name = billable.BillableName()
		customer.Email = billable.BillableEmail()
		customer.TrialEndsAt = trialEndsAt

		if err := store.Create(ctx, customer); err != nil {
			return nil, err
		}

		return customer, nil
	}

	customer.Name = billable.BillableName()
	customer.Email = billable.BillableEmail()
	customer.TrialEndsAt = trialEndsAt

	if err := store.Save(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// BillingPlan returns the plan matching the current valid subscription.
func BillingPlan(manager *Manager, billable Billable, subscription *Subscription) *Plan {
	if subscription == nil || !subscription.Valid() {
		return nil
	}

	for _, plan := range manager.Plans(billable.BillableType()) {
		if subscription.HasPrice(plan.ID) {
			return plan
		}
	}

	return nil
}

// CanUpdateSeats reports whether Billing allows the billable to mutate seat
// quantity.
func CanUpdateSeats(customer *Customer, subscription *Subscription) bool {
	if customer != nil && customer.OnGenericTrial() {
		return true
	}

	return subscription != nil && subscription.Recurring() && !subscription.PastDue()
}

// AddSeats increments the provider subscription quantity.
func AddSeats(ctx context.Context, updater QuantityUpdater, subscription *Subscription, count int, prorates bool) error {
	return changeSeats(ctx, updater, subscription, count, prorates, func(current, delta int) int {
		return current + delta
	})
}

// RemoveSeats decrements the provider subscription quantity.
func RemoveSeats(ctx context.Context, updater QuantityUpdater, subscription *Subscription, count int, prorates bool) error {
	return changeSeats(ctx, updater, subscription, count, prorates, func(current, delta int) int {
		next := current - delta

		if next < 1 {
			return 1
		}

		return next
	})
}

// UpdateSeats sets the provider subscription quantity.
func UpdateSeats(ctx context.Context, updater QuantityUpdater, subscription *Subscription, count int, prorates bool) error {
	if subscription == nil {
		return nil
	}

	if count < 1 {
		count = 1
	}

	return updater.UpdateSubscriptionQuantity(ctx, subscription, count, seatProration(prorates))
}

func changeSeats(
	ctx context.Context,
	updater QuantityUpdater,
	subscription *Subscription,
	count int,
	prorates bool,
	next func(current, delta int) int,
) error {
	if subscription == nil {
		return nil
	}

	if count < 1 {
		count = 1
	}

	return updater.UpdateSubscriptionQuantity(ctx, subscription, next(subscription.Quantity(), count), seatProration(prorates))
}

func seatProration(prorates bool) ProrationBehavior {
	if prorates {
		return ProrateNextBilling
	}

	return FullNextBilling
}
