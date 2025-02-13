package processor

import (
	"context"

	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
)

type InMemProcessor struct {
}

func (i InMemProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	return "inMemPaymentLink", nil
}

func NewInMemProcessor() *InMemProcessor {
	return &InMemProcessor{}
}
