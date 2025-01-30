package main

import (
	"context"
	"github.com/SimonMorphy/gorder/common/config"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/common/server"
	"github.com/SimonMorphy/gorder/order/ports"
	"github.com/SimonMorphy/gorder/order/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")
	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()
		application := service.NewApplication(ctx)
		grpcServer := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, grpcServer)
	})
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, NewHTTPServer(), ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
