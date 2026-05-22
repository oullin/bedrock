package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/httpx"
)

// DownloadInvoiceHandler handles invoice PDF downloads.
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
		errorResponse(w, http.StatusUnauthorized, "Unauthorized")

		return
	}

	if !routeMatchesBillable(r, billable) {
		errorResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))

		return
	}

	transactionID := r.PathValue("transaction")

	if transactionID == "" {
		errorResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))

		return
	}

	txn, err := h.transactions.FindByProviderID(r.Context(), transactionID)

	if err != nil || txn == nil {
		errorResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))

		return
	}

	if txn.BillableType != "" && (txn.BillableType != billable.BillableType() || txn.BillableID != billable.BillableID()) {
		errorResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))

		return
	}

	download, err := h.downloader.DownloadInvoice(r.Context(), txn)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())

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

	body, err := io.ReadAll(download.Body)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	_ = httpx.NewResponse(w).
		Header("Content-Type", contentType).
		Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName)).
		Send(body)
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
