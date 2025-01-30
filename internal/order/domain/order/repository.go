package order

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	Get(ctx context.Context, id, customerId string) (*Order, error)
	Update(
		ctx context.Context,
		o *Order,
		updateFn func(context.Context, *Order) (*Order, error),
	) error
}
