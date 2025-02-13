package command

import (
	"context"

	"github.com/SimonMorphy/gorder/common/decorator"
	domain "github.com/SimonMorphy/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
)

type UpdateOrder struct {
	Order    *domain.Order
	UpdateFn func(ctx context.Context, order *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]

type updateOrderHandler struct {
	orderRepo domain.Repository
}

func (u updateOrderHandler) Handle(ctx context.Context, query UpdateOrder) (interface{}, error) {
	if query.UpdateFn == nil {
		logrus.Warnf("updateOrderHanlder got nil updateFn, order = %v", query.Order)
		query.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		}
	}
	if err := u.orderRepo.Update(ctx, query.Order, query.UpdateFn); err != nil {
		return nil, err
	}
	return nil, nil
}

func NewUpdateOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	client decorator.MetricsClient,
) UpdateOrderHandler {
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		updateOrderHandler{
			orderRepo: orderRepo,
		},
		logger,
		client,
	)
}
