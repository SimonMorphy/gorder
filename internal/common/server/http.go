package server

import (
	"github.com/SimonMorphy/gorder/common/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RunHTTPServer(serverName string, wrapper func(router *gin.Engine)) {

	addr := viper.Sub(serverName).GetString("http-addr")
	if addr == "" {
		panic("empty addr")
	}
	RunHTTPServerOnAddr(addr, wrapper)

}

func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	setMiddleWares(apiRouter)
	wrapper(apiRouter)
	apiRouter.Group("/api")
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}

func setMiddleWares(router *gin.Engine) {
	router.Use(middleware.StructuredLog(logrus.NewEntry(logrus.StandardLogger())))
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware("default-server"))
}
