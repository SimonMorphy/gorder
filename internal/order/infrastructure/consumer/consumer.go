package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"

	"github.com/SimonMorphy/gorder/common/broker"
	"github.com/SimonMorphy/gorder/order/app"
	"github.com/SimonMorphy/gorder/order/app/command"
	domain "github.com/SimonMorphy/gorder/order/domain/order"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	App app.Application
}

func NewConsumer(application app.Application) *Consumer {
	return &Consumer{
		App: application,
	}
}

func (c *Consumer) Listen(ch *amqp091.Channel) {
	queue, err := ch.QueueDeclare(broker.EventOrderPaid, true, false, true, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	err = ch.QueueBind(queue.Name, "", broker.EventOrderPaid, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.handleMessage(msg, queue, ch)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(msg amqp091.Delivery, queue amqp091.Queue, ch *amqp091.Channel) {
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	t := otel.Tracer("rabbitMQ")
	_, span := t.Start(ctx, fmt.Sprintf("rabbitMQ.%s.consume", queue.Name))
	defer span.End()
	o := &domain.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Error(err.Error())
		_ = msg.Nack(false, false)
	}
	_, err := c.App.Commands.UpdateOrder.Handle(context.Background(), command.UpdateOrder{
		Order: o,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := o.IsPaid(); err != nil {
				logrus.Error(err.Error())
				return nil, err
			}
			return order, nil
		},
	})
	if err != nil {
		logrus.Errorf("error order:%s updating ,err=%v", o.ID, err)
		//TODO:RETRY
		return
	}
	span.AddEvent("order.updated")
	_ = msg.Ack(false)
	logrus.Infof("order consume paid event success!")
}
