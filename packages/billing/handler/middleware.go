// Package handler provides HTTP handlers for the Billing billing system.
package handler

import (
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/service"
	"github.com/bedrock/packages/httpx"
)

// VerifyBillableIsSubscribed is middleware that checks whether the
// resolved billable has a valid subscription. If not, it redirects
// HTML requests to the billing portal or returns 402 for JSON/XHR.
// app/Http/Middleware/EnsureTeamSubscribed.
func VerifyBillableIsSubscribed(
	manager *billing.Manager,
	billing *service.BillingService,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			billableType := manager.DefaultBillableType()

			billable, err := manager.ResolveBillable(billableType, r)

			if err != nil {
				redirectToBilling(w, r)

				return
			}

			ctx := r.Context()
			subscribed, err := billing.IsSubscribedToAnyProvider(ctx, billable.BillableType(), billable.BillableID())

			if err != nil || !subscribed {
				redirectToBilling(w, r)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func redirectToBilling(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") == "application/json" ||
		r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		jsonResponse(w, http.StatusPaymentRequired, map[string]string{"message": "Payment required"})

		return
	}

	_ = httpx.NewRedirectResponse(w, r, billing.NewRouteRegistry().URL(billing.RoutePortal, nil), http.StatusFound).Send()
}
