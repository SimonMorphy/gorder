package ports

import (
	"context"
	"github.com/SimonMorphy/gorder/order/convertor"

	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/order/app"
	"github.com/SimonMorphy/gorder/order/app/command"
	"github.com/SimonMorphy/gorder/order/app/query"
	"github.com/SimonMorphy/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	_, err := G.Application.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerId: request.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	result, err := G.Application.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerID: request.CustomerID,
		OrderID:    request.OrderId,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return convertor.NewOrderConverter().EntityToProto(result), nil
}

func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (*emptypb.Empty, error) {
	logrus.Infof("order_grpc || request in || request= %v", request)
	newOrder, err := order.NewOrder(request.ID, request.CustomerId, request.Status, request.PaymentLink, convertor.NewItemConvertor().ProtosToEntities(request.Items))
	if err != nil {
		logrus.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	_, err = G.Application.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: newOrder,
		UpdateFn: func(ctx context.Context, order *order.Order) (*order.Order, error) {
			return newOrder, nil
		},
	})
	return nil, err
}
