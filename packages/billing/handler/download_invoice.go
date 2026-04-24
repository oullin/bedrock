package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bedrock/packages/billing"
)

// DownloadInvoiceHandler handles invoice PDF downloads.
// Mirrors Billing\Http\Controllers\DownloadInvoiceController.
type DownloadInvoiceHandler struct {
	transactions billing.TransactionStore
	resolver     billing.ResolverFunc
	downloader   billing.InvoiceDownloader
}

// NewDownloadInvoiceHandler creates a DownloadInvoiceHandler.
func NewDownloadInvoiceHandler(txns billing.TransactionStore, resolver billing.ResolverFunc, downloader billing.InvoiceDownloader) *DownloadInvoiceHandler {
	return &DownloadInvoiceHandler{transactions: txns, resolver: resolver, downloader: downloader}
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

	download, err := h.downloader.DownloadInvoice(r.Context(), txn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	defer download.Body.Close()

	contentType := download.ContentType

	if contentType == "" {
		contentType = "application/pdf"
	}

	fileName := download.FileName

	if fileName == "" {
		fileName = fmt.Sprintf("invoice-%s.pdf", txn.PaddleID)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, download.Body); err != nil {
		return
	}
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
