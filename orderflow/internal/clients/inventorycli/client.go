package inventorycli

import (
	"context"
	inventorypb "orderflow/internal/gen/inventory"
)

type InventoryClient struct {
	cli inventorypb.InventoryServiceClient
}

func NewClient(cli inventorypb.InventoryServiceClient) *InventoryClient {
	return &InventoryClient{
		cli: cli,
	}
}

func (ic *InventoryClient) GetStock(ctx context.Context, productName string) (int32, error) {
	req := inventorypb.GetStockRequest{
		ProductName: productName,
	}

	res, err := ic.cli.GetStock(ctx, &req)
	if err != nil {
		return 0, err
	}

	return res.AvailableQuantity, nil
}
