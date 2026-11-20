package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"orderflow/internal/order"
	"strconv"
)

func (h *Handler) orderPreviewHandler(w http.ResponseWriter, r *http.Request) {
	var prevReq previewOrderRequest

	if r.Method != http.MethodPost {
		w.Header().Add("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&prevReq)
	if err != nil {
		fmt.Println("fail decode body")
		http.Error(w, "Fail decode body", http.StatusBadRequest)
		return
	}

	orderItems := make([]order.OrderItem, 0, len(prevReq.Items))

	for _, item := range prevReq.Items {
		orderItem := order.OrderItem{
			Product: order.Product{
				Name:  item.Name,
				Price: item.Price,
			},
			Quantity: item.Quantity,
		}

		orderItems = append(orderItems, orderItem)
	}

	createdOrder, err := order.New(orderItems)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
		return
	}

	preResponse := previewOrderResponse{
		Total: createdOrder.Total(),
	}

	w.Header().Add("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(preResponse)
	if err != nil {
		fmt.Println("fail encode")
		http.Error(w, "fail encode", http.StatusBadGateway)
		return
	}
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var prevReq previewOrderRequest

	if r.Method != http.MethodPost {
		w.Header().Add("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&prevReq)
	if err != nil {
		fmt.Println("fail decode body")
		http.Error(w, "Fail decode body", http.StatusBadRequest)
		return
	}

	orderItems := make([]order.OrderItem, 0, len(prevReq.Items))

	for _, item := range prevReq.Items {
		orderItem := order.OrderItem{
			Product: order.Product{
				Name:  item.Name,
				Price: item.Price,
			},
			Quantity: item.Quantity,
		}

		orderItems = append(orderItems, orderItem)
	}

	newOrder, err := h.service.Create(r.Context(), orderItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdOrderResponse := orderResponse{
		ID:     newOrder.ID,
		Status: string(newOrder.Status),
		Total:  newOrder.Total(),
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(createdOrderResponse)
	if err != nil {
		fmt.Println("fail encode")
		http.Error(w, "fail encode", http.StatusBadGateway)
		return
	}
}

func (h *Handler) getOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Fail decode path", http.StatusBadRequest)
		return
	}

	ord, err := h.service.GetByID(r.Context(), orderID)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	createdOrderResponse := orderResponse{
		ID:     ord.ID,
		Status: string(ord.Status),
		Total:  ord.Total(),
	}

	w.Header().Add("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(createdOrderResponse)
	if err != nil {
		fmt.Println("fail encode")
		http.Error(w, "fail encode", http.StatusBadGateway)
		return
	}
}

func (h *Handler) confirmOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Add("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Fail decode path", http.StatusBadRequest)
		return
	}

	ord, err := h.service.Confirm(r.Context(), orderID)
	if err != nil {
		switch {
		case errors.Is(err, order.ErrStockUnavailable):
			http.Error(w, "stock unavailable", http.StatusServiceUnavailable)
		case errors.Is(err, order.ErrInsufficientStock):
			http.Error(w, "insufficient stock", http.StatusConflict)
		case errors.Is(err, order.ErrorderNotFound):
			http.Error(w, "order not found", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	createdOrderResponse := orderResponse{
		ID:     ord.ID,
		Status: string(ord.Status),
		Total:  ord.Total(),
	}

	w.Header().Add("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(createdOrderResponse)
	if err != nil {
		fmt.Println("fail encode")
		http.Error(w, "fail encode", http.StatusBadGateway)
		return
	}
}
