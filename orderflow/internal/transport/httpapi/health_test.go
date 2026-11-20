package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"orderflow/internal/order"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h := NewHandler(order.NewService(order.NewRepo(), NewFakeStock()))

	h.healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("wrong status code")
	}

	var response healthResponse

	err := json.NewDecoder(rec.Result().Body).Decode(&response)
	if err != nil {
		t.Fatalf("invalid body")
	}

	if response.Status != "ok" {
		t.Errorf("wrong response status")
	}

	if rec.Header().Get("Content-type") != "application/json" {
		t.Errorf("wrong header")
	}
}

func TestHealthHandlerMethodNotAllow(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()

	h := NewHandler(order.NewService(order.NewRepo(), NewFakeStock()))

	h.healthHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code")
	}

	if rec.Header().Get("Allow") != http.MethodGet {
		t.Errorf("wrong response header")
	}
}
