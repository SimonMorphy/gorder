package domain

import (
	"context"

	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error)
}

type Order struct {
	ID          string          `validate:"required"`
	CustomerID  string          `validate:"required"`
	Status      string          `validate:"required"`
	PaymentLink string          `validate:"required"`
	Items       []*orderpb.Item `validate:"required,min=1"`
}
