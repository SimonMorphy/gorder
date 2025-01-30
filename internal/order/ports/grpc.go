package ports

import (
	"context"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/order/app"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	Application app.Application
}

func NewGRPCServer(application app.Application) *GRPCServer {
	return &GRPCServer{
		Application: application,
	}

}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
