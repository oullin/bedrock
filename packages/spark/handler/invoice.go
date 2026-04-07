package handler

import (
	"net/http"

	"github.com/bedrock/packages/spark"
)

// InvoiceHandler handles invoice downloads.
type InvoiceHandler struct {
	manager      *spark.Manager
	transactions spark.TransactionStore
	provider     spark.ProviderTransactionManager
}

// NewInvoiceHandler creates an InvoiceHandler.
func NewInvoiceHandler(manager *spark.Manager, transactions spark.TransactionStore, provider spark.ProviderTransactionManager) *InvoiceHandler {
	return &InvoiceHandler{manager: manager, transactions: transactions, provider: provider}
}

// Download redirects to the invoice PDF URL.
func (h *InvoiceHandler) Download(w http.ResponseWriter, r *http.Request) {
	transactionID := r.PathValue("transaction")
	if transactionID == "" {
		http.Error(w, "transaction ID required", http.StatusBadRequest)

		return
	}

	ctx := r.Context()
	txn, err := h.transactions.FindByProviderID(ctx, transactionID)
	if err != nil || txn == nil {
		http.NotFound(w, r)

		return
	}

	pdfURL, err := h.provider.InvoicePDFURL(ctx, txn.ProviderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, pdfURL, http.StatusFound)
}
