package main

import (
	"github.com/SimonMorphy/gorder/common/server"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	server.RunGRPCServer(serviceName, func(server *grpc.Server) {

	})
}
