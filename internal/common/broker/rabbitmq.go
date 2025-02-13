package broker

import (
	"context"
	"github.com/SimonMorphy/gorder/common/config/models"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

const (
	DLX = "dlx_order"
	DLQ = "dlq"
)

func Connect(rv *models.RabbitMQVoucher) (*amqp091.Channel, func() error) {
	address := rv.DSN()
	conn, err := amqp091.Dial(address)
	if err != nil {
		logrus.Fatal(err)
	}
	channel, err := conn.Channel()
	if err != nil {
		logrus.Fatal(err)
	}
	err = channel.ExchangeDeclare(EventOrderCreated, "direct", true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	err = channel.ExchangeDeclare(EventOrderPaid, "fanout", true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	if err = createDLX(channel); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("success to connect to RabbitMQ")
	return channel, channel.Close
}

func createDLX(channel *amqp091.Channel) error {
	q, err := channel.QueueDeclare("share.queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = channel.ExchangeDeclare(DLX, "fanout", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = channel.QueueBind(q.Name, "", DLX, false, nil)
	if err != nil {
		return err
	}
	_, err = channel.QueueDeclare(DLQ, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

type RabbitMQHeaderCarrier map[string]interface{}

func (r RabbitMQHeaderCarrier) Get(key string) string {
	v, e := r[key]
	if !e {
		return ""
	}
	return v.(string)
}

func (r RabbitMQHeaderCarrier) Set(key string, value string) {
	r[key] = value
}

func (r RabbitMQHeaderCarrier) Keys() []string {
	keys := make([]string, 0)
	for key := range r {
		keys = append(keys, key)
	}
	return keys
}

func InjectRabbitMQHeaders(ctx context.Context) map[string]interface{} {
	carrier := make(RabbitMQHeaderCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier
}

func ExtractRabbitMQHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, RabbitMQHeaderCarrier(headers))
}
