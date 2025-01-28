package ports

import (
	"context"
	"github.com/SimonMorphy/gorder/common/genproto/stockpb"
)

type GRPCServer struct {
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) CheckItemInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}
