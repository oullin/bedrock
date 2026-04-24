package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/action"
	"github.com/bedrock/packages/spark/handler"
)

func TestNewSubscriptionHandler_ResolverError(t *testing.T) {
	mgr := spark.NewManager()
	creator := action.NewSubscriptionCreator(&testSubStore{}, mgr)
	h := handler.NewNewSubscriptionHandler(creator, mgr, func(*http.Request) (spark.Billable, error) {
		return nil, errors.New("no billable")
	})

	req := httptest.NewRequest(http.MethodPost, "/spark/subscription", bytes.NewReader([]byte(`{}`)))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), spark.ErrBillableRequired.Error()) {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestNewSubscriptionHandler_InvalidJSON(t *testing.T) {
	mgr := spark.NewManager()
	creator := action.NewSubscriptionCreator(&testSubStore{}, mgr)
	h := handler.NewNewSubscriptionHandler(creator, mgr, func(*http.Request) (spark.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/spark/subscription", bytes.NewReader([]byte("{not json")))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}
}

func TestNewSubscriptionHandler_InvalidPlanReturnsValidationError(t *testing.T) {
	mgr := spark.NewManager()
	mgr.Plan("team", "Pro", "pri_pro").Monthly()

	creator := action.NewSubscriptionCreator(&testSubStore{}, mgr)
	h := handler.NewNewSubscriptionHandler(creator, mgr, func(*http.Request) (spark.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	body, _ := json.Marshal(map[string]string{"plan": "unknown"})
	req := httptest.NewRequest(http.MethodPost, "/spark/subscription", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("code = %d, want 422", rec.Code)
	}
}

func TestNewSubscriptionHandler_HappyPath(t *testing.T) {
	mgr := spark.NewManager()
	mgr.Plan("team", "Pro", "pri_pro").Monthly()

	creator := action.NewSubscriptionCreator(&testSubStore{}, mgr)
	h := handler.NewNewSubscriptionHandler(creator, mgr, func(*http.Request) (spark.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	body, _ := json.Marshal(map[string]string{"plan": "pri_pro"})
	req := httptest.NewRequest(http.MethodPost, "/spark/subscription", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}

	var resp map[string]any

	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestNewSubscriptionHandler_CreatorError(t *testing.T) {
	mgr := spark.NewManager()
	mgr.Plan("team", "Pro", "pri_pro").Monthly()
	mgr.Billable("team").CheckPlanEligibility(func(spark.Billable, spark.Plan) error {
		return errors.New("boom")
	})

	creator := action.NewSubscriptionCreator(&testSubStore{}, mgr)
	h := handler.NewNewSubscriptionHandler(creator, mgr, func(*http.Request) (spark.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	body, _ := json.Marshal(map[string]string{"plan": "pri_pro"})
	req := httptest.NewRequest(http.MethodPost, "/spark/subscription", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("code = %d, want 500", rec.Code)
	}
}
