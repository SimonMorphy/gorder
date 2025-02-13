package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, query C) (R, error)
}

func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, client MetricsClient) CommandHandler[C, R] {
	return QueryLoggingDecorator[C, R]{
		logger: logger,
		base: QueryMetricsDecorator[C, R]{
			base:   handler,
			client: client,
		},
	}
}
