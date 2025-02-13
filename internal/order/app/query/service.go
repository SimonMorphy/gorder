package query

import (
	"context"

	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/common/genproto/stockpb"
)

type StockService interface {
	CheckIfItemsInStock(ctx context.Context, itemIds []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
	GetItems(ctx context.Context, itemIds []string) ([]*orderpb.Item, error)
}
