package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/galkinallan/L0/internal/cache"
	"github.com/galkinallan/L0/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func MustCreateDb(ctx context.Context) *Storage {
	db, err := sqlx.ConnectContext(ctx, "postgres", "user=root dbname=orders password=secret sslmode=disable port=5432 host=127.0.0.1")
	if err != nil {
		log.Fatalln(err)
	}

	return &Storage{
		db: db,
	}
}

func (orders *Storage) CreateOrder(ctx context.Context, orderUid string, order models.Order) error {
	newCtx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("cant marshal order %v", err)
	}

	tx, err := orders.db.BeginTx(newCtx, nil)
	if err != nil {
		return fmt.Errorf("begin tranc error %v", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(newCtx, "INSERT INTO orders (order_uid, order_info) VALUES ($1, $2)", orderUid, data)
	if err != nil {
		return fmt.Errorf("tranc execute error %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit error %v", err)
	}

	return nil
}

func (orders *Storage) GetOrderById(ctx context.Context, orderUid string) (*models.Order, error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	var order models.Order
	row := orders.db.QueryRowxContext(newCtx, "SELECT order_info  FROM orders WHERE order_uid = $1", orderUid)

	err := row.StructScan(&order)

	if err != nil {
		return nil, fmt.Errorf("no such order with that id %v", err)
	}

	fmt.Println(order.OrderUID)

	return &order, nil
}

func (orders *Storage) Close() error {
	return orders.db.Close()
}

func (orders *Storage) FillCache(ctx context.Context, cache *cache.MemCache) error {
	newCtx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	rows, err := orders.db.QueryContext(newCtx, "SELECT order_info FROM orders")
	if err != nil {
		fmt.Println("cant restore cache from postgres")
		return err
	}

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order); err != nil {
			fmt.Println("unable to scan values")
			return err
		}
		cache.Set(&order)
	}

	return nil
}
