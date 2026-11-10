package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"orderflow/internal/order"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

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

	healthHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code")
	}

	if rec.Header().Get("Allow") != http.MethodGet {
		t.Errorf("wrong response header")
	}
}

func TestOrderPreviewHandler(t *testing.T) {
	var resp previewOrderResponse
	api := api{
		service: order.NewService(order.NewRepo()),
	}
	req := httptest.NewRequest(
		http.MethodPost,
		"/orders/preview",
		strings.NewReader(`{"items":[{"name":"Keyboard","price":1500,"quantity":2}]}`),
	)
	rec := httptest.NewRecorder()

	api.orderPreviewHandler(rec, req)

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("wrong content-type header")
	}

	if rec.Code != http.StatusOK {
		t.Errorf("wrong status code")
	}

	if err := json.NewDecoder(rec.Result().Body).Decode(&resp); err != nil {
		t.Fatalf("fail to decode response: %v", err)
	}

	if resp.Total != 3000 {
		t.Errorf("wrong answer for total")
	}

}

func TestWrongOrderInvalidJSONPreviewHandler(t *testing.T) {
	api := api{
		service: order.NewService(order.NewRepo()),
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders/preview",
		strings.NewReader(`{"items":[{"name":"Keyboard"]}`),
	)
	rec := httptest.NewRecorder()

	api.orderPreviewHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("wrong status code")
	}
}

func TestWrongOrderPreviewHandlerAllowMethod(t *testing.T) {
	api := api{
		service: order.NewService(order.NewRepo()),
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/preview",
		strings.NewReader(`{"items":[{"name":"Keyboard","price":1500,"quantity":2}]}`),
	)
	rec := httptest.NewRecorder()

	api.orderPreviewHandler(rec, req)

	if rec.Header().Get("Allow") != http.MethodPost {
		t.Errorf("wrong allow method")
	}

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code")
	}
}

func TestWrongOrderInvalidDataPreviewHandler(t *testing.T) {
	api := api{
		service: order.NewService(order.NewRepo()),
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders/preview",
		strings.NewReader(`{"items":[{"name":"Keyboard","price":1500,"quantity":0}]}`),
	)
	rec := httptest.NewRecorder()

	api.orderPreviewHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("wrong status code")
	}
}

func TestCreateOrder(t *testing.T) {
	var respData orderResponse

	api := api{
		service: order.NewService(order.NewRepo()),
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		strings.NewReader(`{"items":[{"name":"Keyboard","price":1500,"quantity":2}]}`),
	)
	rec := httptest.NewRecorder()

	api.createOrder(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("wrong response code")
	}

	if err := json.NewDecoder(rec.Result().Body).Decode(&respData); err != nil {
		t.Fatalf("invalid response")
	}

	if respData.ID != 1 || respData.Total != 3000 {
		t.Errorf("wrong data in response")
	}
}

func TestGetOrderByID(t *testing.T) {
	var getOrderResp orderResponse
	ctx := context.Background()

	api := api{
		service: order.NewService(order.NewRepo()),
	}

	items := []order.OrderItem{
		{
			Product: order.Product{
				Name:  "Test1",
				Price: 100,
			},
			Quantity: 1,
		},
		{
			Product: order.Product{
				Name:  "Test2",
				Price: 50,
			},
			Quantity: 2,
		},
	}

	_, err := api.service.Create(ctx, items)
	if err != nil {
		t.Fatalf("fail to create order")
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/1",
		nil,
	)

	rec := httptest.NewRecorder()
	req.SetPathValue("id", "1")

	api.getOrderByID(rec, req)

	err = json.NewDecoder(rec.Result().Body).Decode(&getOrderResp)
	if err != nil {
		t.Fatalf("fail to parse response")
	}

	if getOrderResp.Total != 200 || getOrderResp.ID != 1 {
		t.Errorf("wrong resp data")
	}
}

func TestWrongIdGetOrder(t *testing.T) {
	api := api{
		service: order.NewService(order.NewRepo()),
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/3",
		nil,
	)

	rec := httptest.NewRecorder()
	req.SetPathValue("id", "3")

	api.getOrderByID(rec, req)

	if rec.Result().StatusCode != http.StatusNotFound {
		t.Errorf("wrong status code")
	}
}
