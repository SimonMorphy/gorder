package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type QueryLoggingDecorator[C, R any] struct {
	logger *logrus.Entry
	base   QueryHandler[C, R]
}

func (q QueryLoggingDecorator[C, R]) Handle(ctx context.Context, command C) (result R, err error) {
	logger := logrus.StandardLogger()
	logger.Debug("Executing Query")
	defer func() {
		if err == nil {
			logger.Info("Query Execute Successfully")
		} else {
			logger.Error("Fail to Execute")
		}
	}()
	return q.base.Handle(ctx, command)
}

func generateActionName(command any) string {
	return strings.Split(fmt.Sprintf("%T", command), ".")[1]
}
