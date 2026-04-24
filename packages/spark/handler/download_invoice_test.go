package handler_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx/routingx"
	"github.com/bedrock/packages/routing"
	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/handler"
)

type invoiceTransactionStore struct {
	transactions map[string]*spark.Transaction
	err          error
}

type invoiceDownloader struct {
	download *spark.InvoiceDownload
	err      error
}

func (s invoiceTransactionStore) FindByProviderID(_ context.Context, providerID string) (*spark.Transaction, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.transactions[providerID], nil
}

func (s invoiceTransactionStore) FindByBillable(context.Context, string, int64, int) ([]spark.Transaction, error) {
	return nil, nil
}

func (s invoiceTransactionStore) Create(context.Context, *spark.Transaction) error { return nil }
func (s invoiceTransactionStore) Save(context.Context, *spark.Transaction) error   { return nil }

func (d invoiceDownloader) DownloadInvoice(context.Context, *spark.Transaction) (*spark.InvoiceDownload, error) {
	if d.err != nil {
		return nil, d.err
	}

	if d.download != nil {
		return d.download, nil
	}

	return &spark.InvoiceDownload{
		FileName:    "invoice.pdf",
		ContentType: "application/pdf",
		Body:        io.NopCloser(strings.NewReader("%PDF")),
	}, nil
}

// AgreementControllerTest::test_invoice_download_route_exists
// AgreementControllerTest::test_invoice_download_is_scoped_to_the_current_team
// AgreementControllerTest::test_invoice_download_ignores_a_tampered_team_id_request_attribute
func TestDownloadInvoiceHandlerScopesInvoicesToResolvedBillable(t *testing.T) {
	store := invoiceTransactionStore{transactions: map[string]*spark.Transaction{
		"txn_current": {PaddleID: "txn_current", BillableType: "team", BillableID: 10},
		"txn_foreign": {PaddleID: "txn_foreign", BillableType: "team", BillableID: 20},
	}}
	billable := &stubBillable{id: 10, btype: "team"}
	resolver := func(r *http.Request) (spark.Billable, error) {
		return billable, nil
	}

	invoices := handler.NewDownloadInvoiceHandler(store, resolver, invoiceDownloader{})
	router := routing.NewRouter(nil, nil)
	router.Get("/spark/{type}/{id}/invoices/{transaction}/download", invoices.Download).Name(spark.RouteInvoiceDownload)
	dispatcher := routingx.NewHandler(router)

	req := httptest.NewRequest(http.MethodGet, "/spark/team/10/invoices/txn_current/download", nil)
	rec := httptest.NewRecorder()
	dispatcher.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("current team invoice status = %d, want 200", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/spark/team/10/invoices/txn_foreign/download", nil)
	rec = httptest.NewRecorder()
	dispatcher.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("foreign invoice status = %d, want 404", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/spark/team/20/invoices/txn_foreign/download", nil)
	req.Header.Set("X-Team-ID", "20")
	rec = httptest.NewRecorder()
	dispatcher.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("tampered team header status = %d, want 404", rec.Code)
	}
}

func TestDownloadInvoiceHandler_ResolverError(t *testing.T) {
	h := handler.NewDownloadInvoiceHandler(
		invoiceTransactionStore{},
		func(*http.Request) (spark.Billable, error) { return nil, errors.New("no") },
		invoiceDownloader{},
	)

	req := httptest.NewRequest(http.MethodGet, "/spark/team/1/invoices/tx/download", nil)
	rec := httptest.NewRecorder()
	h.Download(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("code = %d, want 401", rec.Code)
	}
}

func TestDownloadInvoiceHandler_EmptyTransactionID(t *testing.T) {
	h := handler.NewDownloadInvoiceHandler(
		invoiceTransactionStore{},
		func(*http.Request) (spark.Billable, error) { return &stubBillable{id: 1, btype: "team"}, nil },
		invoiceDownloader{},
	)

	req := httptest.NewRequest(http.MethodGet, "/spark/download", nil)
	rec := httptest.NewRecorder()
	h.Download(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404", rec.Code)
	}
}

func TestDownloadInvoiceHandler_StoreError(t *testing.T) {
	h := handler.NewDownloadInvoiceHandler(
		invoiceTransactionStore{err: errors.New("boom")},
		func(*http.Request) (spark.Billable, error) { return &stubBillable{id: 1, btype: "team"}, nil },
		invoiceDownloader{},
	)

	req := httptest.NewRequest(http.MethodGet, "/spark/download", nil)
	req.SetPathValue("transaction", "tx_1")
	rec := httptest.NewRecorder()
	h.Download(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404", rec.Code)
	}
}

func TestDownloadInvoiceHandler_DownloaderError(t *testing.T) {
	store := invoiceTransactionStore{transactions: map[string]*spark.Transaction{
		"tx_1": {PaddleID: "tx_1", BillableType: "team", BillableID: 10},
	}}

	h := handler.NewDownloadInvoiceHandler(
		store,
		func(*http.Request) (spark.Billable, error) { return &stubBillable{id: 10, btype: "team"}, nil },
		invoiceDownloader{err: errors.New("down")},
	)

	req := httptest.NewRequest(http.MethodGet, "/spark/download", nil)
	req.SetPathValue("transaction", "tx_1")
	rec := httptest.NewRecorder()
	h.Download(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("code = %d, want 500", rec.Code)
	}
}

func TestDownloadInvoiceHandler_DefaultContentTypeAndFileName(t *testing.T) {
	store := invoiceTransactionStore{transactions: map[string]*spark.Transaction{
		"tx_1": {PaddleID: "tx_1", BillableType: "team", BillableID: 10},
	}}

	h := handler.NewDownloadInvoiceHandler(
		store,
		func(*http.Request) (spark.Billable, error) { return &stubBillable{id: 10, btype: "team"}, nil },
		invoiceDownloader{download: &spark.InvoiceDownload{
			Body: io.NopCloser(strings.NewReader("%PDF")),
		}},
	)

	req := httptest.NewRequest(http.MethodGet, "/spark/download", nil)
	req.SetPathValue("transaction", "tx_1")
	rec := httptest.NewRecorder()
	h.Download(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200", rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/pdf" {
		t.Fatalf("content-type = %q, want application/pdf", got)
	}

	if got := rec.Header().Get("Content-Disposition"); !strings.Contains(got, "invoice-tx_1.pdf") {
		t.Fatalf("content-disposition = %q", got)
	}
}
