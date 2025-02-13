package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SimonMorphy/gorder/common/broker"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/payment/app"
	"github.com/SimonMorphy/gorder/payment/app/command"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
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
	queue, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("fail to consume: queue=%s, err=%v", queue.Name, err)
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
	tracer := otel.Tracer("rabbitMQ")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("rabbitMQ.%s.consume", queue.Name))
	defer span.End()
	logrus.Infof("Payment recieve a message from %s, msg= %v", queue.Name, string(msg.Body))
	o := &orderpb.Order{}
	err := json.Unmarshal(msg.Body, o)
	if err != nil {
		logrus.Warnf("failed to unmarshall msg to order ,err=%v", err)
		_ = msg.Nack(false, false)
		return
	}
	_, err = c.App.Commands.CreatePayment.Handle(ctx, command.CreatePayment{
		Order: o,
	})
	if err != nil {
		//TODO: retry
		logrus.Warnf("failed to create payment,err=%v", err)
		_ = msg.Nack(false, false)
		return
	}
	_ = msg.Ack(false)
	span.AddEvent("payment.created")
	logrus.Info("consume success")
}
