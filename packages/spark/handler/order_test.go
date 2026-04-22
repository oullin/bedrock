package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/handler"
)

type orderStoreWithData struct {
	orders []spark.Order
}

func (s *orderStoreWithData) FindByID(_ context.Context, _ int64) (*spark.Order, error) {
	return nil, nil
}
func (s *orderStoreWithData) FindByBillable(_ context.Context, _ int64, _ int) ([]spark.Order, error) {
	return s.orders, nil
}
func (s *orderStoreWithData) Create(_ context.Context, _ *spark.Order) error { return nil }
func (s *orderStoreWithData) Save(_ context.Context, _ *spark.Order) error   { return nil }
func (s *orderStoreWithData) HasCompletedForProduct(_ context.Context, _ int64, _ int64) (bool, error) {
	return false, nil
}

// OrderTest::test_guest_cannot_access_orders
func TestOrderHandler_Index_UnauthenticatedReturns401(t *testing.T) {
	resolver := func(r *http.Request) (spark.Billable, error) {
		return nil, errors.New("unauthenticated")
	}

	h := handler.NewOrderHandler(&orderStoreWithData{}, resolver)

	req := httptest.NewRequest("GET", "/billing/orders", nil)
	rec := httptest.NewRecorder()
	h.Index(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// OrderTest::test_authenticated_user_can_view_orders
func TestOrderHandler_Index_AuthenticatedReturnsOK(t *testing.T) {
	billable := &stubBillable{id: 1, btype: "team"}
	resolver := func(r *http.Request) (spark.Billable, error) {
		return billable, nil
	}

	h := handler.NewOrderHandler(&orderStoreWithData{}, resolver)

	req := httptest.NewRequest("GET", "/billing/orders", nil)
	rec := httptest.NewRecorder()
	h.Index(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

// OrderTest::test_orders_page_shows_team_orders
func TestOrderHandler_Index_ReturnsTeamOrders(t *testing.T) {
	billable := &stubBillable{id: 1, btype: "team"}
	resolver := func(r *http.Request) (spark.Billable, error) {
		return billable, nil
	}

	store := &orderStoreWithData{
		orders: []spark.Order{
			{ID: 1, TeamID: 1, Status: spark.OrderStatusCompleted},
			{ID: 2, TeamID: 1, Status: spark.OrderStatusPending},
		},
	}

	h := handler.NewOrderHandler(store, resolver)

	req := httptest.NewRequest("GET", "/billing/orders", nil)
	rec := httptest.NewRecorder()
	h.Index(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp map[string]any

	json.NewDecoder(rec.Body).Decode(&resp)

	orders, ok := resp["orders"].(map[string]any)

	if !ok {
		t.Fatal("response missing 'orders' key")
	}

	data, ok := orders["data"].([]any)

	if !ok {
		t.Fatal("response missing 'orders.data' key")
	}

	if len(data) != 2 {
		t.Errorf("orders.data count = %d, want 2", len(data))
	}
}
