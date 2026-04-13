package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/service"
)

// GatewayHandler handles the payment provider selection page.
// Mirrors app/Http/Controllers/BillingGatewayController.
type GatewayHandler struct {
	manager  *billing.Manager
	billing  *service.BillingService
	products billing.ProductStore
	resolver billing.ResolverFunc
}

// NewGatewayHandler creates a GatewayHandler.
func NewGatewayHandler(
	mgr *billing.Manager,
	billing *service.BillingService,
	products billing.ProductStore,
	resolver billing.ResolverFunc,
) *GatewayHandler {
	return &GatewayHandler{
		manager:  mgr,
		billing:  billing,
		products: products,
		resolver: resolver,
	}
}

// Show returns the provider selection page data as JSON, including
// available subscription and one-time products.
func (h *GatewayHandler) Show(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	ctx := r.Context()

	// If already subscribed, redirect to the provider's portal.
	subscribed, _ := h.billing.IsSubscribedToAnyProvider(ctx, billable.BillableType(), billable.BillableID())

	if subscribed {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"redirect": "/billing"})

		return
	}

	subProducts, _ := h.products.ActiveSubscriptions(ctx)
	oneTimeProducts, _ := h.products.ActiveOneTime(ctx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"subscriptionProducts": subProducts,
		"oneTimeProducts":      oneTimeProducts,
	})
}

// Store validates and persists the chosen payment provider.
func (h *GatewayHandler) Store(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	var input struct {
		Provider string `json:"provider"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	if !billing.IsValidProvider(input.Provider) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": map[string][]string{
				"provider": {"The selected provider is invalid."},
			},
		})

		return
	}

	if pc, ok := billable.(billing.ProviderConfigurable); ok {
		pc.SetPaymentProvider(input.Provider)
	}

	redirect := "/billing"

	if input.Provider == "paddle" {
		redirect = "/billing/paddle"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"redirect": redirect})
}
