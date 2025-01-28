package main

import (
	"github.com/SimonMorphy/gorder/common/server"
	"github.com/SimonMorphy/gorder/order/ports"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	switch serverType {

	}
}
