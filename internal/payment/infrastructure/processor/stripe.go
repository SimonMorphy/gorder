package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SimonMorphy/gorder/common/tracing"

	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type StripeProcessor struct {
	ApiKey string
}

const (
	SuccessURL = "http://localhost:8282/success"
)

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	if s.ApiKey == "" {
		panic("Stripe Api Key is null")
	}
	stripe.Key = s.ApiKey

	ctx, span := tracing.Start(ctx, "stripe_processor.create_payment")
	defer span.End()

	var items []*stripe.CheckoutSessionLineItemParams
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	marshalItems, _ := json.Marshal(order.Items)
	metaData := map[string]string{
		"orderID":     order.ID,
		"customerID":  order.CustomerId,
		"status":      order.Status,
		"items":       string(marshalItems),
		"paymentLink": order.PaymentLink,
	}
	params := &stripe.CheckoutSessionParams{
		Metadata:   metaData,
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s?customerID=%s&orderID=%s", SuccessURL, order.CustomerId, order.ID)),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

func NewStripeProcessor(apiKey string) *StripeProcessor {
	return &StripeProcessor{ApiKey: apiKey}
}
