package service

import (
	"context"

	grpcClient "github.com/SimonMorphy/gorder/common/client"
	"github.com/SimonMorphy/gorder/common/config"
	"github.com/SimonMorphy/gorder/common/metrics"
	"github.com/SimonMorphy/gorder/order/adapters"
	"github.com/SimonMorphy/gorder/order/adapters/grpc"
	"github.com/SimonMorphy/gorder/order/app"
	"github.com/SimonMorphy/gorder/order/app/command"
	"github.com/SimonMorphy/gorder/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func() error) {
	GRPCClient, f, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	client := grpc.NewStockGPRC(GRPCClient)

	return newApplication(ctx, client), func() error {
		_ = f()
		return nil
	}
}

func newApplication(_ context.Context, service query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NewTodoMetrics()
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(
				orderRepo, service, logger, metricsClient, config.RabbitMQ,
			),
			UpdateOrder: command.NewUpdateOrderHandler(
				orderRepo, logger, metricsClient,
			),
		},
		Queries: app.Queries{GetCustomerOrder: query.NewGetCustomerOrderHandler(
			orderRepo, logger, metricsClient)},
	}
}
