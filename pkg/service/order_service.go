package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/galkinallan/L0/internal/cache"
	"github.com/galkinallan/L0/internal/db"
	"github.com/galkinallan/L0/internal/models"
)

type OrderService struct {
	cache *cache.MemCache
	db    db.Storage
}

func CreateOrderService(db *db.Storage, cache *cache.MemCache) *OrderService {
	return &OrderService{
		db:    *db,
		cache: cache,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, orderUid string, order models.Order) error {

	s.cache.Set(&order)
	return s.db.CreateOrder(ctx, orderUid, order)
}

func (s *OrderService) GetById(ctx context.Context, orderUid string) (*models.Order, error) {
	data := s.cache.Get(orderUid)
	if data != nil {
		return data, nil
	}

	return s.db.GetOrderById(ctx, orderUid)
}

func (s *OrderService) GetHttpHandle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUid := r.URL.Path[1:]
		order, err := s.GetById(context.Background(), orderUid)
		if err != nil {
			fmt.Fprintf(w, "Error %v", err)
		}

		fmt.Fprintf(w, "order %v", order)
	}
}
