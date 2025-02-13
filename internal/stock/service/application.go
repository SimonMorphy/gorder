package service

import (
	"context"

	"github.com/SimonMorphy/gorder/common/metrics"
	"github.com/SimonMorphy/gorder/stock/adapters"
	"github.com/SimonMorphy/gorder/stock/app"
	"github.com/SimonMorphy/gorder/stock/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	stockRepo := adapters.NewStockMemoryRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NewTodoMetrics()

	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckItemInStockHandler: query.NewCheckItemHandler(
				stockRepo, logger, metricsClient,
			),
			GetItemHandler: query.NewGetItemHandler(
				stockRepo, logger, metricsClient,
			),
		},
	}
}
