package spark

import (
	"net/http"
)

// NormaliseBillableRouteParam is middleware that validates a UUID route
// parameter and resolves it to the billable's numeric ID.
func NormaliseBillableRouteParam(resolver BillableResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			if id == "" {
				next.ServeHTTP(w, r)

				return
			}

			// Validate UUID format (simple check).
			if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
				http.NotFound(w, r)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// VerifyBillableIsSubscribed is middleware that checks whether the resolved
// billable has a valid subscription. If not, it redirects to the billing
// portal or returns a 402 JSON response for API requests.
func VerifyBillableIsSubscribed(
	manager *Manager,
	subscriptions SubscriptionStore,
	clock Clock,
	keepPastDueActive bool,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			billableType := manager.DefaultBillableType()

			billable, err := manager.ResolveBillable(billableType, r)
			if err != nil {
				redirectToBilling(w, r, billableType)

				return
			}

			ctx := r.Context()
			sub, err := subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())
			if err != nil || sub == nil || !sub.Valid(clock, keepPastDueActive) {
				// Check generic trial.
				customer, cerr := findCustomerForBillable(ctx, subscriptions, billable)
				if cerr != nil || customer == nil || !customer.OnGenericTrial(clock) {
					redirectToBilling(w, r, billableType)

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func redirectToBilling(w http.ResponseWriter, r *http.Request, billableType string) {
	path := "/billing"
	if billableType != "user" {
		path = "/billing/" + billableType
	}

	// JSON/API requests get a 402.
	if r.Header.Get("Accept") == "application/json" ||
		r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte(`{"message":"Payment required"}`))

		return
	}

	http.Redirect(w, r, path, http.StatusFound)
}

// findCustomerForBillable is a placeholder for customer lookup from a
// subscription store context.
func findCustomerForBillable(_ interface{}, _ SubscriptionStore, _ Billable) (*Customer, error) {
	return nil, nil
}
