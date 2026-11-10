package postgres

import (
	"context"
	"fmt"
	"orderflow/internal/order"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Save(ctx context.Context, ord *order.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("fail to start tx: %w", err)
	}

	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "insert into orders (status) values ($1) returning id", ord.Status).Scan(&ord.ID)
	if err != nil {
		return fmt.Errorf("fail insert order: %w", err)
	}

	for _, item := range ord.Items {
		_, err := tx.Exec(
			ctx,
			"insert into order_items (order_id, product_name, price, quantity) values ($1, $2, $3, $4)",
			ord.ID, item.Product.Name, item.Product.Price, item.Quantity,
		)
		if err != nil {
			return fmt.Errorf("fail insert order_item: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("fail commit tx: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	curOrd := order.Order{}
	err := r.pool.QueryRow(ctx, "select id, status from orders where id = $1", id).Scan(&curOrd.ID, &curOrd.Status)
	if err != nil {
		return nil, fmt.Errorf("select db err: %w", err)
	}

	var pdtItems []order.OrderItem
	rows, err := r.pool.Query(ctx, "select product_name, price, quantity from order_items where order_id = $1 order by id", curOrd.ID)
	if err != nil {
		return nil, fmt.Errorf("select db err: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		item := order.OrderItem{}
		err := rows.Scan(&item.Product.Name, &item.Product.Price, &item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("read row err: %w", err)
		}

		pdtItems = append(pdtItems, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("db error after read rows: %w", err)
	}

	curOrd.Items = pdtItems

	return &curOrd, nil
}
