package adapters

import (
	"context"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/stock/domain/stock"
	"sync"
)

var stub = map[string]*orderpb.Item{
	"item_id": {
		ID:       "foo_item",
		Name:     "stub_item",
		Quantity: 10000,
		PriceID:  "stub_item_price_id",
	},
}

type MemoryStockRepository struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var (
		items     []*orderpb.Item
		missingId []string
	)
	for _, id := range ids {
		if item, exist := m.store[id]; exist {
			items = append(items, item)
		} else {
			missingId = append(missingId, id)
		}
	}
	if len(items) == len(missingId) {
		return items, nil
	}
	return items, stock.NotFountError{MissingIds: missingId}
}

func NewMemoryOrderRepository() *MemoryStockRepository {
	return &MemoryStockRepository{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}
