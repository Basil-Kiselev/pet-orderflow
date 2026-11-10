package postgres

import (
	"context"
	"orderflow/internal/order"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestSave(t *testing.T) {
	repo, pool := newTestRepository(t)
	ctx := context.Background()
	ord, err := order.New([]order.OrderItem{
		{
			Product: order.Product{
				Name:  "Test item",
				Price: 200,
			},
			Quantity: 4,
		},
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if err := repo.Save(ctx, ord); err != nil {
		t.Fatalf("save order: %v", err)
	}

	if ord.ID == 0 {
		t.Fatal("wrong order ID")
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, "delete from orders where id = $1", ord.ID)
	})

	var status string

	err = pool.QueryRow(ctx, "select status from orders where id = $1", ord.ID).Scan(&status)
	if err != nil {
		t.Fatalf("read saved order: %v", err)
	}
	if status != string(order.StatusNew) {
		t.Errorf("wrong order status: %s", status)
	}

	var savedItems []order.OrderItem

	rows, err := pool.Query(ctx, "select id, product_name, price, quantity from order_items where order_id = $1", ord.ID)
	if err != nil {
		t.Fatalf("read saved order items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id       int64
			name     string
			price    int64
			quantity int32
		)

		err := rows.Scan(&id, &name, &price, &quantity)
		if err != nil {
			t.Fatalf("fail scan order item: %v", err)
		}

		savedItems = append(savedItems, order.OrderItem{
			Product: order.Product{
				Name:  name,
				Price: int(price),
			},
			Quantity: int(quantity),
		})
	}

	if err = rows.Err(); err != nil {
		t.Fatalf("db rows err: %v", err)
	}

	if len(savedItems) != 1 {
		t.Fatalf("expected 1 item, got %d", len(savedItems))
	}

	item := savedItems[0]

	if item.Product.Name != "Test item" {
		t.Errorf("wrong product name: %q", item.Product.Name)
	}

	if item.Product.Price != 200 {
		t.Errorf("wrong price: %d", item.Product.Price)
	}

	if item.Quantity != 4 {
		t.Errorf("wrong quantity: %d", item.Quantity)
	}
}

func TestSaveRollback(t *testing.T) {
	repo, pool := newTestRepository(t)
	ctx := context.Background()

	var countBefore int
	err := pool.QueryRow(ctx, "select count(id) from orders").Scan(&countBefore)
	if err != nil {
		t.Fatalf("invalid select from db: %v", err)
	}

	ord := order.Order{
		Status: order.StatusNew,
		Items: []order.OrderItem{
			{
				Quantity: 1,
				Product: order.Product{
					Name:  "test1",
					Price: 500,
				},
			},
			{
				Quantity: 0,
				Product: order.Product{
					Name:  "test",
					Price: 5000,
				},
			},
		},
	}

	err = repo.Save(ctx, &ord)
	if err == nil {
		t.Errorf("wrong result, expected err - bad request to save order")
	}

	var newCount int
	err = pool.QueryRow(ctx, "select count(id) from orders").Scan(&newCount)
	if err != nil {
		t.Fatalf("invalid select from db: %v", err)
	}

	if countBefore != newCount {
		t.Errorf("invalid rollback")
	}

}

func newTestRepository(t *testing.T) (*Repository, *pgxpool.Pool) {
	t.Helper()

	_ = godotenv.Load("../../../../.env")

	databaseUrl := os.Getenv("TEST_DATABASE_URL")
	if databaseUrl == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("ping database: %v", err)
	}

	t.Cleanup(pool.Close)

	return NewRepo(pool), pool
}

func TestGetByID(t *testing.T) {
	repo, pool := newTestRepository(t)
	ctx := context.Background()

	ord, _ := order.New([]order.OrderItem{
		{
			Product: order.Product{
				Name:  "Test item",
				Price: 200,
			},
			Quantity: 4,
		},
		{
			Product: order.Product{
				Name:  "Test item2",
				Price: 200,
			},
			Quantity: 1,
		},
	})

	if err := repo.Save(ctx, ord); err != nil {
		t.Fatalf("save order: %v", err)
	}

	curOrder, err := repo.GetByID(ctx, ord.ID)
	if err != nil {
		t.Fatalf("error get order: %v", err)
	}

	if curOrder.Total() != 1000 {
		t.Errorf("wrong total price")
	}

	if curOrder.ID != ord.ID {
		t.Errorf("wrong order id")
	}

	if curOrder.Status != order.StatusNew {
		t.Errorf("wrong order status")
	}

	if len(curOrder.Items) != 2 {
		t.Errorf("wrong order items count")
	}

	t.Cleanup(func() {
		pool.Exec(ctx, "delete from orders where id = $1", ord.ID)
	})
}
