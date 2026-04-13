// Package handler provides HTTP handlers for the Spark billing system.
package handler

import (
	"net/http"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/service"
)

// VerifyBillableIsSubscribed is middleware that checks whether the
// resolved billable has a valid subscription. If not, it redirects
// HTML requests to the billing gateway or returns 402 for JSON/XHR.
// Mirrors Spark\Http\Middleware\VerifyBillableIsSubscribed and
// app/Http/Middleware/EnsureTeamSubscribed.
func VerifyBillableIsSubscribed(
	manager *spark.Manager,
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte(`{"message":"Payment required"}`))

		return
	}

	http.Redirect(w, r, "/billing/choose-provider", http.StatusFound)
}
