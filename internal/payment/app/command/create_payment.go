package command

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/SimonMorphy/gorder/common/decorator"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/payment/domain"
	"github.com/sirupsen/logrus"
)

type CreatePayment struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]

type createPaymentHandler struct {
	processor domain.Processor
	orderGrpc OrderService
}

func (c createPaymentHandler) Handle(ctx context.Context, query CreatePayment) (string, error) {
	tracer := otel.Tracer("payment")
	ctx, span := tracer.Start(ctx, "create_payment")
	defer span.End()
	link, err := c.processor.CreatePaymentLink(ctx, query.Order)
	if err != nil {
		return "", err
	}
	logrus.Infof("create payment link for %v success,payment link :%s", query.Order, link)
	o := &orderpb.Order{
		ID:          query.Order.ID,
		CustomerId:  query.Order.CustomerId,
		Status:      "waiting_for_payment",
		Items:       query.Order.Items,
		PaymentLink: link,
	}
	err = c.orderGrpc.UpdateOrder(ctx, o)
	if err != nil {
		logrus.Errorf("failed to update order,err=%v", err)
		return "", err
	}
	return link, nil
}

func NewCreatePaymentHandler(
	processor domain.Processor,
	service OrderService,
	logger *logrus.Entry,
	client decorator.MetricsClient,
) CreatePaymentHandler {
	return decorator.ApplyCommandDecorators[CreatePayment, string](
		createPaymentHandler{
			processor: processor,
			orderGrpc: service,
		}, logger, client,
	)
}
