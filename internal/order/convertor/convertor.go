package convertor

import (
	client "github.com/SimonMorphy/gorder/common/client/order"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/order/domain/order"
	"github.com/SimonMorphy/gorder/order/entity"
)

type OrderConvertor struct {
}

type ItemConvertor struct {
}

type ItemWithQuantityConvertor struct {
}

func (c *ItemWithQuantityConvertor) EntitiesToProtos(items []*entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.EntityToProto(item))
	}
	return res
}

func (c *ItemWithQuantityConvertor) EntityToProto(item *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
	return &orderpb.ItemWithQuantity{
		Id:       item.Id,
		Quantity: item.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ProtosToEntities(items []*orderpb.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return
}

func (c *ItemWithQuantityConvertor) ProtoToEntity(item *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		Id:       item.Id,
		Quantity: item.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ClientToEntities(items []client.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.ClientToEntity(item))
	}
	return
}

func (c *ItemWithQuantityConvertor) ClientToEntity(item client.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		Id:       item.Id,
		Quantity: item.Quantity,
	}
}

func (c *OrderConvertor) EntityToProto(o *order.Order) *orderpb.Order {
	c.check(o)
	return &orderpb.Order{
		ID:          o.ID,
		CustomerId:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntitiesToProtos(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *order.Order {
	c.check(o)
	return &order.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerId,
		Status:      o.Status,
		Items:       NewItemConvertor().ProtosToEntities(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ClientToEntity(o *client.Order) *order.Order {
	c.check(o)
	return &order.Order{
		ID:          o.Id,
		CustomerID:  o.CustomerId,
		Status:      o.Status,
		Items:       NewItemConvertor().ClientsToEntities(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) EntityToClient(o *order.Order) *client.Order {
	c.check(o)
	return &client.Order{
		Id:          o.ID,
		CustomerId:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntitiesToClients(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) check(o interface{}) {
	if o == nil {
		panic("order is nil")
	}
}

func (c *ItemConvertor) EntitiesToProtos(items []*entity.Item) (res []*orderpb.Item) {
	for _, item := range items {
		res = append(res, c.EntityToProto(item))
	}
	return
}
func (c *ItemConvertor) ProtosToEntities(items []*orderpb.Item) (res []*entity.Item) {
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return
}

func (c *ItemConvertor) ClientsToEntities(items []client.Item) (res []*entity.Item) {
	for _, item := range items {
		res = append(res, c.ClientToEntity(item))
	}
	return
}

func (c *ItemConvertor) EntitiesToClients(items []*entity.Item) (res []client.Item) {
	for _, item := range items {
		res = append(res, c.EntityToClient(item))
	}
	return
}

func (c *ItemConvertor) EntityToProto(item *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemConvertor) ProtoToEntity(item *orderpb.Item) *entity.Item {
	return &entity.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemConvertor) ClientToEntity(item client.Item) *entity.Item {
	return &entity.Item{
		ID:       item.Id,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceId,
	}
}

func (c *ItemConvertor) EntityToClient(item *entity.Item) client.Item {
	return client.Item{
		Id:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceId:  item.PriceID,
	}
}
