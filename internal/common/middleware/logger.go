package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func StructuredLog(entry *logrus.Entry) gin.HandlerFunc {
	return func(context *gin.Context) {
		now := time.Now()
		context.Next()
		elapsed := time.Since(now)
		entry.WithFields(logrus.Fields{
			"time_elapsed_ms": elapsed.Milliseconds(),
			"request_url":     context.Request.RequestURI,
			"client_ip":       context.ClientIP(),
			"full_path":       context.FullPath(),
		}).Info("request_out")
	}
}
