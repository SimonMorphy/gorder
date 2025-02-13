package query

import (
	"context"

	"github.com/SimonMorphy/gorder/common/decorator"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type GetItem struct {
	Ids []string
}

type getItemHandler struct {
	stockRepo stock.Repository
}

func (g getItemHandler) Handle(ctx context.Context, query GetItem) ([]*orderpb.Item, error) {
	items, err := g.stockRepo.GetItems(ctx, query.Ids)
	if err != nil {
		return nil, err
	}
	return items, nil
}

type GetItemHandler decorator.QueryHandler[GetItem, []*orderpb.Item]

func NewGetItemHandler(
	stockRepo stock.Repository,
	logger *logrus.Entry,
	client decorator.MetricsClient,
) GetItemHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	return decorator.ApplyQueryDecorators[GetItem, []*orderpb.Item](
		getItemHandler{stockRepo: stockRepo},
		logger,
		client,
	)
}
