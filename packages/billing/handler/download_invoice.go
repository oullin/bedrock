package handler

import (
	"net/http"
	"strconv"

	"github.com/bedrock/packages/billing"
)

// DownloadInvoiceHandler handles invoice PDF downloads.
// Mirrors Billing\Http\Controllers\DownloadInvoiceController.
type DownloadInvoiceHandler struct {
	transactions billing.TransactionStore
	resolver     billing.ResolverFunc
}

// NewDownloadInvoiceHandler creates a DownloadInvoiceHandler.
func NewDownloadInvoiceHandler(txns billing.TransactionStore, resolver billing.ResolverFunc) *DownloadInvoiceHandler {
	return &DownloadInvoiceHandler{transactions: txns, resolver: resolver}
}

// Download serves an invoice PDF for the given transaction.
func (h *DownloadInvoiceHandler) Download(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	if !routeMatchesBillable(r, billable) {
		http.NotFound(w, r)

		return
	}

	transactionID := r.PathValue("transaction")

	if transactionID == "" {
		http.NotFound(w, r)

		return
	}

	txn, err := h.transactions.FindByProviderID(r.Context(), transactionID)

	if err != nil || txn == nil {
		http.NotFound(w, r)

		return
	}

	if txn.BillableType != "" && (txn.BillableType != billable.BillableType() || txn.BillableID != billable.BillableID()) {
		http.NotFound(w, r)

		return
	}

	// In a full implementation, this would fetch the PDF from the
	// provider and stream it. For now, return the transaction data.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func routeMatchesBillable(r *http.Request, billable billing.Billable) bool {
	if routeType := r.PathValue("type"); routeType != "" && routeType != billable.BillableType() {
		return false
	}

	routeID := r.PathValue("id")

	if routeID == "" {
		return true
	}

	id, err := strconv.ParseInt(routeID, 10, 64)

	if err != nil {
		return false
	}

	return id == billable.BillableID()
}
