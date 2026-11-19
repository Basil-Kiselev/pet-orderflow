package httpapi

import (
	"net/http"
	"orderflow/internal/order"
)

type Handler struct {
	service *order.Service
}

type healthResponse struct {
	Status string `json:"status"`
}

type previewOrderRequest struct {
	Items []previewOrderItem `json:"items"`
}

type previewOrderItem struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}
type previewOrderResponse struct {
	Total int `json:"total"`
}

type orderResponse struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
	Total  int    `json:"total"`
}

type getOrderByIDRequest struct {
	ID int `json:"id"`
}

func NewHandler(service *order.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.healthHandler)
	mux.HandleFunc("/orders", h.createOrder)
	mux.HandleFunc("/orders/preview", h.orderPreviewHandler)
	mux.HandleFunc("/orders/{id}", h.getOrderByID)
	mux.HandleFunc("/orders/{id}/confirm", h.confirmOrder)

	return mux
}
