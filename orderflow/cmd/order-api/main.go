package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"orderflow/internal/order"
	"orderflow/internal/storage/postgres"
	"orderflow/internal/storage/redis"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	goredis "github.com/redis/go-redis/v9"
)

const RedisTTL = 1 * time.Minute

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

type api struct {
	service *order.Service
}

func main() {
	ctx := context.Background()
	_ = godotenv.Load("../.env")

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

	if err := pool.Ping(ctx); err != nil {
		fmt.Printf("fail to ping db: %v", err)
		return
	}

	defer pool.Close()

	redisCli := goredis.NewClient(&goredis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	defer redisCli.Close()

	redisRepo := redis.New(postgres.NewRepo(pool), redisCli, RedisTTL)

	api := api{
		service: order.NewService(redisRepo),
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/orders/preview", api.orderPreviewHandler)
	http.HandleFunc("/orders", api.createOrder)
	http.HandleFunc("/orders/{id}", api.getOrderByID)

	fmt.Println("OrderFlow server started on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("failed to start server: %v", err)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ok := healthResponse{
		Status: "ok",
	}
	w.Header().Add("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(ok)
	if err != nil {
		fmt.Printf("fail to encode: %v", err)
		return
	}
}

func (api *api) orderPreviewHandler(w http.ResponseWriter, r *http.Request) {
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

func (api *api) createOrder(w http.ResponseWriter, r *http.Request) {
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

	newOrder, err := api.service.Create(r.Context(), orderItems)
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

func (api *api) getOrderByID(w http.ResponseWriter, r *http.Request) {
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

	order, err := api.service.GetByID(r.Context(), orderID)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	createdOrderResponse := orderResponse{
		ID:     order.ID,
		Status: string(order.Status),
		Total:  order.Total(),
	}

	w.Header().Add("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(createdOrderResponse)
	if err != nil {
		fmt.Println("fail encode")
		http.Error(w, "fail encode", http.StatusBadGateway)
		return
	}
}
