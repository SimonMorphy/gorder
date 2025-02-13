package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SimonMorphy/gorder/common/broker"
	"github.com/SimonMorphy/gorder/common/decorator"
	"github.com/SimonMorphy/gorder/order/app/query"
	"github.com/SimonMorphy/gorder/order/convertor"
	domain "github.com/SimonMorphy/gorder/order/domain/order"
	"github.com/SimonMorphy/gorder/order/entity"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type CreateOrder struct {
	CustomerId string
	Items      []*entity.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderId string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
	channel   *amqp091.Channel
}

func (c createOrderHandler) Handle(ctx context.Context, query CreateOrder) (*CreateOrderResult, error) {
	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Warn("queue declare failed")
	}
	tracer := otel.Tracer("rabbitMQ")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("rabbitMQ.%s.publish", q.Name))
	defer span.End()
	validateItems, err := c.validate(ctx, query.Items)
	if err != nil {
		return nil, err
	}
	create, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: query.CustomerId,
		Items:      validateItems,
	})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(create)
	if err != nil {
		return nil, err
	}
	header := broker.InjectRabbitMQHeaders(ctx)
	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp091.Persistent,
		Body:         data,
		Headers:      header,
	})
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{OrderId: create.ID}, nil
}

func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have at least one item")
	}
	items = packItems(items)
	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
	if err != nil {
		return nil, err
	}
	logrus.Info("Items:", resp.Items,
		"itemsLen", len(resp.Items))
	var ids []string
	for _, item := range items {
		ids = append(ids, item.Id)
	}
	if len(ids) == len(resp.Items) {
		return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
	}
	return nil, domain.StockLackError{}
}

func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.Id] += item.Quantity
	}
	var res []*entity.ItemWithQuantity
	for id, quan := range merged {
		res = append(res, &entity.ItemWithQuantity{
			Id:       id,
			Quantity: quan,
		})
	}
	return res
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	StockService query.StockService,
	logger *logrus.Entry,
	client decorator.MetricsClient,
	channel *amqp091.Channel,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	if channel == nil {
		panic("nil channel")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{
			orderRepo: orderRepo,
			stockGRPC: StockService,
			channel:   channel,
		},
		logger,
		client,
	)
}
