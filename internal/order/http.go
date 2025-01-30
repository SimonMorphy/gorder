package main

import "github.com/gin-gonic/gin"

type HTTPServer struct {
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerId string) {

}

func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {

}
