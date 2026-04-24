package handler

import (
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/state"
	"github.com/bedrock/packages/routegen"
)

// PortalHandler renders the billing portal frontend state.
// Mirrors Billing\Http\Controllers\BillingPortalController.
type PortalHandler struct {
	manager  *billing.Manager
	frontend *state.FrontendState
	resolver billing.ResolverFunc
	routes   *routegen.Registry
}

// NewPortalHandler creates a PortalHandler.
func NewPortalHandler(mgr *billing.Manager, fs *state.FrontendState, resolver billing.ResolverFunc) *PortalHandler {
	routes := billing.NewRouteRegistry()
	fs.WithRoutes(routes)

	return &PortalHandler{manager: mgr, frontend: fs, resolver: resolver, routes: routes}
}

// WithRoutes sets the route registry used by the portal shell.
func (h *PortalHandler) WithRoutes(routes *routegen.Registry) *PortalHandler {
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
		errorResponse(w, http.StatusBadRequest, billing.ErrBillableRequired.Error())

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
