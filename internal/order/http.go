package main

import (
	"fmt"
	"github.com/SimonMorphy/gorder/common"
	"github.com/SimonMorphy/gorder/common/client/order"
	"github.com/SimonMorphy/gorder/order/app"
	"github.com/SimonMorphy/gorder/order/app/command"
	"github.com/SimonMorphy/gorder/order/app/dto"
	"github.com/SimonMorphy/gorder/order/app/query"
	"github.com/SimonMorphy/gorder/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app app.Application
	common.BaseResponse
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerId string) {
	var (
		req  order.CreateOrderRequest
		err  error
		resp dto.CreateOrderResponse
	)
	defer func() {
		H.Response(c, err, &resp)
	}()
	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerId: customerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientToEntities(req.Items),
	})
	if err != nil {
		return
	}
	resp = dto.CreateOrderResponse{
		CustomerID:  req.CustomerId,
		OrderID:     r.OrderId,
		RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderId),
	}
}

func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
	var (
		err  error
		resp struct {
			order *order.Order
		}
	)
	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		CustomerID: customerId,
		OrderID:    orderId,
	})
	defer func() {
		H.Response(c, err, resp)
	}()
	if err != nil {
		return
	}
	resp.order = convertor.NewOrderConverter().EntityToClient(o)
}
