package cache

import (
	"fmt"
	"sync"

	"github.com/galkinallan/L0/internal/models"
)

type MemCache struct {
	sync.RWMutex
	Cache map[string]*models.Order
}

func CreateCache() *MemCache {
	newMemCache := MemCache{}
	newMemCache.Cache = make(map[string]*models.Order)

	return &newMemCache
}

func (cache *MemCache) Set(order *models.Order) {
	cache.Lock()
	defer cache.Unlock()

	cache.Cache[order.OrderUID] = order
}

func (cache *MemCache) Get(orderId string) *models.Order {
	cache.RLock()
	defer cache.RUnlock()

	order, _ := cache.Cache[orderId]

	return order
}

func (cache *MemCache) PrintKeys() {
	fmt.Println("Data from Cache")
	for k, _ := range cache.Cache {
		fmt.Printf("OrderUid: %v\n", k)
	}
}
