package handler

import (
	"net/http"

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
	_, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

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

	// In a full implementation, this would fetch the PDF from the
	// provider and stream it. For now, return the transaction data.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
