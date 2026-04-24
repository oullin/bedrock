package handler_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/handler"
)

type invoiceTransactionStore struct {
	transactions map[string]*billing.Transaction
}

type invoiceDownloader struct{}

func (s invoiceTransactionStore) FindByProviderID(_ context.Context, providerID string) (*billing.Transaction, error) {
	return s.transactions[providerID], nil
}

func (s invoiceTransactionStore) FindByBillable(context.Context, string, int64, int) ([]billing.Transaction, error) {
	return nil, nil
}

func (s invoiceTransactionStore) Create(context.Context, *billing.Transaction) error { return nil }
func (s invoiceTransactionStore) Save(context.Context, *billing.Transaction) error   { return nil }

func (invoiceDownloader) DownloadInvoice(context.Context, *billing.Transaction) (*billing.InvoiceDownload, error) {
	return &billing.InvoiceDownload{
		FileName:    "invoice.pdf",
		ContentType: "application/pdf",
		Body:        io.NopCloser(strings.NewReader("%PDF")),
	}, nil
}

// AgreementControllerTest::test_invoice_download_route_exists
// AgreementControllerTest::test_invoice_download_is_scoped_to_the_current_team
// AgreementControllerTest::test_invoice_download_ignores_a_tampered_team_id_request_attribute
func TestDownloadInvoiceHandlerScopesInvoicesToResolvedBillable(t *testing.T) {
	store := invoiceTransactionStore{transactions: map[string]*billing.Transaction{
		"txn_current": {PaddleID: "txn_current", BillableType: "team", BillableID: 10},
		"txn_foreign": {PaddleID: "txn_foreign", BillableType: "team", BillableID: 20},
	}}
	billable := &stubBillable{id: 10, btype: "team"}
	resolver := func(r *http.Request) (billing.Billable, error) {
		return billable, nil
	}

	invoices := handler.NewDownloadInvoiceHandler(store, resolver, invoiceDownloader{})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /agreement/invoices/{transaction}/download", invoices.Download)

	req := httptest.NewRequest(http.MethodGet, "/agreement/invoices/txn_current/download", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("current team invoice status = %d, want 200", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/agreement/invoices/txn_foreign/download", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("foreign invoice status = %d, want 404", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/agreement/invoices/txn_foreign/download", nil)
	req.Header.Set("X-Team-ID", "20")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("tampered team header status = %d, want 404", rec.Code)
	}
}
