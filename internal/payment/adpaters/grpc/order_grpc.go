package grpc

import (
	"context"
	"github.com/SimonMorphy/gorder/common/tracing"

	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/sirupsen/logrus"
)

type OrderGRPC struct {
	client orderpb.OrderServiceClient
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{client: client}
}

func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) error {
	ctx, span := tracing.Start(ctx, "order.grpc.update_order")
	defer span.End()
	_, err := o.client.UpdateOrder(ctx, order)
	logrus.Infof("payment_adpater || uddate_order,err=%v ", err)
	return err
}
