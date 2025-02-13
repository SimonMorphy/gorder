package main

import (
	"context"
	"github.com/SimonMorphy/gorder/common/tracing"

	"github.com/SimonMorphy/gorder/common/config"
	"github.com/SimonMorphy/gorder/common/discovery"
	"github.com/SimonMorphy/gorder/common/server"
	"github.com/SimonMorphy/gorder/payment/infrastructure/consumer"
	"github.com/SimonMorphy/gorder/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var cleanUps []func() error

func init() {
	config.NewLogrusConfig(config.WithServiceName("payment"))
	f, err := config.InitConfig()
	if err != nil {
		logrus.Fatal(err, "初始化失败！")
	}
	cleanUps = append(cleanUps, f)

}

func main() {
	serviceName := viper.GetString("payment.service-name")
	serverType := viper.GetString("payment.server-to-run")

	ctx, cancelFunc := context.WithCancel(context.Background())
	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)
	application, cleanUp := service.NewApplication(ctx)
	defer func() {
		_ = cleanUp()
		cancelFunc()
		for _, up := range cleanUps {
			_ = up()
		}
	}()
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	cleanUps = append(cleanUps, deregisterFunc)

	go consumer.NewConsumer(application).Listen(config.RabbitMQ)

	handler := NewPaymentHandler()
	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, handler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported server type")
	default:
		logrus.Panic("unsupported server type")
	}

}
