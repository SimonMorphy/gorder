package main

import (
	"context"
	"github.com/SimonMorphy/gorder/common/tracing"

	"github.com/SimonMorphy/gorder/common/config"
	"github.com/SimonMorphy/gorder/common/discovery"
	"github.com/SimonMorphy/gorder/common/genproto/stockpb"
	"github.com/SimonMorphy/gorder/common/server"
	"github.com/SimonMorphy/gorder/stock/ports"
	"github.com/SimonMorphy/gorder/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	config.NewLogrusConfig(config.WithServiceName("stock"))
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	logrus.Infof(serverType)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)
	application := service.NewApplication(ctx)
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()
	switch serverType {
	case "http":
	//server.RunHTTPServer()
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			grpcServer := ports.NewGRPCServer(application)
			stockpb.RegisterStockServiceServer(server, grpcServer)
		})
	default:
		panic("unexpected server type")
	}

}
