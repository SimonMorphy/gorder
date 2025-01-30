package ports

import (
	"context"
	"github.com/SimonMorphy/gorder/common/genproto/stockpb"
	"github.com/SimonMorphy/gorder/stock/app"
)

type GRPCServer struct {
	Application app.Application
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) CheckItemInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewGRPCServer(application app.Application) *GRPCServer {
	return &GRPCServer{
		Application: application,
	}
}
