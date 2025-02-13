package grpc

import (
	"context"

	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

type StockGPRC struct {
	client stockpb.StockServiceClient
}

func NewStockGPRC(client stockpb.StockServiceClient) *StockGPRC {
	return &StockGPRC{client: client}
}

func (s StockGPRC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
	resp, err := s.client.CheckItemInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
	logrus.Info("stock_grpc response:", resp)
	return resp, err
}

func (s StockGPRC) GetItems(ctx context.Context, itemIds []string) ([]*orderpb.Item, error) {
	items, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIds: itemIds})
	if err != nil {
		return nil, err
	}
	return items.Items, nil
}
