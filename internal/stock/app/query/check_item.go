package query

import (
	"context"

	"github.com/SimonMorphy/gorder/common/decorator"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type CheckItemInStock struct {
	Items []*orderpb.ItemWithQuantity
}

var stub = map[string]string{
	"item-1": "price_1Qo7pdBRPM1tZxlwM4jByiPE",
	"item-2": "price_1Qo7stBRPM1tZxlwIbsU5Z2G",
	"item-3": "price_1QosgYBRPM1tZxlw98pmhVLh",
}

type checkItemInStockHandler struct {
	stockRepo stock.Repository
}

func (c checkItemInStockHandler) Handle(ctx context.Context, query CheckItemInStock) ([]*orderpb.Item, error) {
	var (
		ids          []string
		idToQuantity = make(map[string]int32)
	)
	for _, item := range query.Items {
		ids = append(ids, item.Id)
		idToQuantity[item.Id] += item.Quantity
	}
	items, err := c.stockRepo.GetItems(ctx, ids)
	if err != nil {
		return nil, err
	}
	var res []*orderpb.Item
	for _, item := range items {
		need, exist := idToQuantity[item.ID]
		if !exist {
			continue
		}
		if item.Quantity >= need {
			res = append(res, &orderpb.Item{
				ID:       item.ID,
				Name:     item.Name,
				Quantity: idToQuantity[item.ID],
				//PriceID:  item.PriceID,
				PriceID: stub[item.ID],
			})
		}
	}
	return res, nil
}

type CheckItemInStockHandler decorator.QueryHandler[CheckItemInStock, []*orderpb.Item]

func NewCheckItemHandler(
	stockRepo stock.Repository,
	logger *logrus.Entry,
	client decorator.MetricsClient,
) CheckItemInStockHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	return decorator.ApplyQueryDecorators[CheckItemInStock, []*orderpb.Item](
		checkItemInStockHandler{stockRepo: stockRepo},
		logger,
		client,
	)
}
