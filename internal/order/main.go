package main

import (
	"context"
	"github.com/SimonMorphy/gorder/common/tracing"

	"github.com/SimonMorphy/gorder/common/config"
	"github.com/SimonMorphy/gorder/common/discovery"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/common/server"
	"github.com/SimonMorphy/gorder/order/infrastructure/consumer"
	"github.com/SimonMorphy/gorder/order/ports"
	"github.com/SimonMorphy/gorder/order/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cleanUps []func() error

func init() {
	config.NewLogrusConfig(config.WithLevel(logrus.InfoLevel), config.WithServiceName("order"))
	f, err := config.InitConfig()
	if err != nil {
		logrus.Fatal(err, "初始化失败！")
	}
	cleanUps = append(cleanUps, f)

}

func main() {
	serviceName := viper.GetString("order.service-name")
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	defer func() {
		for _, up := range cleanUps {
			_ = up()
		}
	}()
	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)
	application, f := service.NewApplication(ctx)
	defer f()
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	go consumer.NewConsumer(application).Listen(config.RabbitMQ)
	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		grpcServer := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, grpcServer)
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		router.StaticFile("/success", "../../public/success.html")
		ports.RegisterHandlersWithOptions(router, HTTPServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
