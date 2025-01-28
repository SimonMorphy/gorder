package main

import (
	"github.com/SimonMorphy/gorder/common/config"
	"github.com/SimonMorphy/gorder/common/server"
	"github.com/SimonMorphy/gorder/order/ports"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})
}
