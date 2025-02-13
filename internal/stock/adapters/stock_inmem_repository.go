package adapters

import (
	"context"
	"sync"

	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/stock/domain/stock"
)

var stub = map[string]*orderpb.Item{
	"item-1": {
		ID:       "item-1",
		Name:     "stub_item",
		Quantity: 1000,
		PriceID:  "stub-item-price-id",
	},
	"item-2": {
		ID:       "item-2",
		Name:     "stub_item-2",
		Quantity: 1000,
		PriceID:  "stub-item-price-id-2",
	},
	"item-3": {
		ID:       "item-3",
		Name:     "stub_item-3",
		Quantity: 1230,
		PriceID:  "stub-item-price-id-3",
	},
}

type StockMemoryRepository struct {
	lock  *sync.RWMutex
	items map[string]*orderpb.Item
}

func NewStockMemoryRepository() *StockMemoryRepository {
	return &StockMemoryRepository{lock: &sync.RWMutex{}, items: stub}
}

func (s *StockMemoryRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var (
		items   []*orderpb.Item
		missing []string
	)
	for _, id := range ids {
		if item, exist := s.items[id]; exist {
			items = append(items, item)
		} else {
			missing = append(missing, id)
		}
	}
	if len(missing) == 0 {
		return items, nil
	}
	return nil, stock.NotFountError{MissingIds: missing}
}
