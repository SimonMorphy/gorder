package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"io"
	"net/http"

	"github.com/SimonMorphy/gorder/common/broker"
	"github.com/SimonMorphy/gorder/common/config"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/payment/domain"
	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

type PaymentHandler struct {
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", h.handleWebHook)
}

func (h *PaymentHandler) handleWebHook(c *gin.Context) {
	logrus.Info("Got webhook from stripe")
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf("Error reading request body: %v\n", err)
		c.Writer.WriteHeader(http.StatusServiceUnavailable)
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return
	}

	endpointSecret := viper.GetString("endpoint-stripe-secret")
	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		logrus.Errorf("Error verifying webhook signature: %v\n", err)
		c.Writer.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			logrus.Errorf("error unmarshal event.data.raw into session, err=%v", err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			logrus.Infof("payment for checkout session %v success!", session.ID)
			ctx, cancelFunc := context.WithCancel(context.TODO())
			defer cancelFunc()
			var items []*orderpb.Item
			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)
			o, err := json.Marshal(&domain.Order{
				ID:          session.Metadata["orderID"],
				CustomerID:  session.Metadata["customerID"],
				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
				PaymentLink: session.Metadata["paymentLink"],
				Items:       items,
			})
			if err != nil {
				logrus.Infof("error ocured in marshal domain order")
				return
			}
			tracer := otel.Tracer("rabbitMQ")
			ctx, span := tracer.Start(ctx, fmt.Sprintf("rabbitMQ.%s.publish", broker.EventOrderCreated))
			defer span.End()
			headers := broker.InjectRabbitMQHeaders(ctx)
			_ = config.RabbitMQ.PublishWithContext(ctx, broker.EventOrderPaid, "", false, false, amqp091.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp091.Persistent,
				Body:         o,
				Headers:      headers,
			})
			logrus.Infof("message published to %s, body:%s", broker.EventOrderPaid, string(o))
		}
	}
	c.JSON(http.StatusOK, nil)
}
