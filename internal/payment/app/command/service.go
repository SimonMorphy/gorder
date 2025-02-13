package command

import (
	"context"

	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
