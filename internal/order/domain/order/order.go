package order

import (
	"fmt"
	"github.com/SimonMorphy/gorder/order/entity"

	"github.com/go-playground/validator/v10"
	"github.com/stripe/stripe-go/v81"
)

type Order struct {
	ID          string         `validate:"required"`
	CustomerID  string         `validate:"required"`
	Status      string         `validate:"required"`
	PaymentLink string         `validate:"required"`
	Items       []*entity.Item `validate:"required,min=1"`
}

func NewOrder(ID string, customerID string, status string, paymentLink string, items []*entity.Item) (*Order, error) {
	order := Order{ID: ID, CustomerID: customerID, Status: status, PaymentLink: paymentLink, Items: items}
	validate := validator.New()
	err := validate.Struct(order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *Order) IsPaid() error {
	if o.Status == string(stripe.CheckoutSessionPaymentStatusPaid) {
		return nil
	}
	return fmt.Errorf("order not paid,id=%s,status=%s", o.ID, o.Status)
}

type NotFundError struct {
	OrderId string
}
type StockLackError struct {
}

func (s StockLackError) Error() string {
	return fmt.Sprintf("item  in lack of quantity")
}

func (n NotFundError) Error() string {
	return fmt.Sprintf("order '%s' not found ", n.OrderId)
}
