package service

import (
	"context"

	"github.com/SimonMorphy/gorder/common/client"
	"github.com/SimonMorphy/gorder/common/metrics"
	"github.com/SimonMorphy/gorder/payment/adpaters/grpc"
	"github.com/SimonMorphy/gorder/payment/domain"
	"github.com/SimonMorphy/gorder/payment/infrastructure/processor"
	"github.com/spf13/viper"

	"github.com/SimonMorphy/gorder/payment/app"
	"github.com/SimonMorphy/gorder/payment/app/command"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func() error) {
	grpcClient, c, err := client.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := grpc.NewOrderGRPC(grpcClient)
	process := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, process), func() error {
		_ = c()
		return nil
	}
}

func newApplication(_ context.Context, service command.OrderService, process domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NewTodoMetrics()
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(
				process, service, logger, metricsClient),
		},
		Queries: app.Queries{},
	}
}
