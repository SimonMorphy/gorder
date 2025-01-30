package main

import (
	"context"
	"github.com/SimonMorphy/gorder/common/config"
	"github.com/SimonMorphy/gorder/common/genproto/stockpb"
	"github.com/SimonMorphy/gorder/common/server"
	"github.com/SimonMorphy/gorder/stock/ports"
	"github.com/SimonMorphy/gorder/stock/service"
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
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application := service.NewApplication(ctx)
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
