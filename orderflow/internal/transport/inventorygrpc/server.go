package inventorygrpc

import (
	"context"
	"orderflow/internal/gen/inventory"
)

type Server struct {
	inventory.UnimplementedInventoryServiceServer
	stock map[string]int32
}

func NewServer(stock map[string]int32) *Server {
	return &Server{
		stock: stock,
	}
}

func (s *Server) GetStock(
	ctx context.Context,
	req *inventory.GetStockRequest,
) (*inventory.GetStockResponse, error) {

	return &inventory.GetStockResponse{
		AvailableQuantity: s.stock[req.ProductName],
	}, nil
}
