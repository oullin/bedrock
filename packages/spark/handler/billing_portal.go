package handler

import (
	"net/http"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/state"
	"github.com/bedrock/packages/wayfinder"
)

// PortalHandler renders the billing portal frontend state.
// Mirrors Spark\Http\Controllers\BillingPortalController.
type PortalHandler struct {
	manager  *spark.Manager
	frontend *state.FrontendState
	resolver spark.ResolverFunc
	routes   *wayfinder.Registry
}

// NewPortalHandler creates a PortalHandler.
func NewPortalHandler(mgr *spark.Manager, fs *state.FrontendState, resolver spark.ResolverFunc) *PortalHandler {
	routes := spark.NewRouteRegistry()
	fs.WithRoutes(routes)

	return &PortalHandler{manager: mgr, frontend: fs, resolver: resolver, routes: routes}
}

// WithRoutes sets the route registry used by the portal shell.
func (h *PortalHandler) WithRoutes(routes *wayfinder.Registry) *PortalHandler {
	if routes != nil {
		h.routes = routes

		if h.frontend != nil {
			h.frontend.WithRoutes(routes)
		}
	}

	return h
}

// Show returns the billing portal state as JSON.
func (h *PortalHandler) Show(w http.ResponseWriter, r *http.Request) {
	data, ok := h.portalState(w, r)

	if !ok {
		return
	}

	renderPortalShell(w, r, h.routes, "Billing", data)
}

// State returns the billing portal state as JSON for the Vue portal.
func (h *PortalHandler) State(w http.ResponseWriter, r *http.Request) {
	data, ok := h.portalState(w, r)

	if !ok {
		return
	}

	jsonResponse(w, http.StatusOK, data)
}

func (h *PortalHandler) portalState(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	billableType := h.manager.DefaultBillableType()

	billable, err := h.resolver(r)

	if err != nil {
		errorResponse(w, http.StatusBadRequest, spark.ErrBillableRequired.Error())

		return nil, false
	}

	if !h.manager.IsAuthorized(billable, r) {
		errorResponse(w, http.StatusForbidden, "Forbidden")

		return nil, false
	}

	data, err := h.frontend.Current(r.Context(), billableType, billable)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())

		return nil, false
	}

	return data, true
}
