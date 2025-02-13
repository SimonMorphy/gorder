package ports

import (
	"context"
	"github.com/SimonMorphy/gorder/common/tracing"

	"github.com/SimonMorphy/gorder/common/genproto/stockpb"
	"github.com/SimonMorphy/gorder/stock/app"
	"github.com/SimonMorphy/gorder/stock/app/query"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	Application app.Application
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	ctx, span := tracing.Start(ctx, "stock.grpc.getItems")
	defer span.End()
	logrus.Info("rpc_request In stock.GetItem")
	defer func() {
		logrus.Info("rpc_request out stock.GetItem")
	}()
	items, err := G.Application.Queries.GetItemHandler.Handle(ctx, query.GetItem{Ids: request.ItemIds})
	response := &stockpb.GetItemsResponse{
		Items: items,
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return response, nil
}

func (G GRPCServer) CheckItemInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	ctx, span := tracing.Start(ctx, "stock.grpc.checkItem")
	defer span.End()
	logrus.Info("rpc_request In stock.CheckItem")
	defer func() {
		logrus.Info("rpc_request out stock.CheckItem")
	}()
	items, err := G.Application.Queries.CheckItemInStockHandler.Handle(ctx, query.CheckItemInStock{
		Items: request.Items,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	response := &stockpb.CheckIfItemsInStockResponse{
		Items: items,
	}
	return response, nil
}

func NewGRPCServer(application app.Application) *GRPCServer {
	return &GRPCServer{
		Application: application,
	}
}
