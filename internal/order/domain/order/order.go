package order

import (
	"fmt"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
)

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
}

type NotFundError struct {
	OrderId string
}

func (n NotFundError) Error() string {
	return fmt.Sprintf("order '%s' not found ", n.OrderId)
}
