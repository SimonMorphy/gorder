package common

import (
	"github.com/SimonMorphy/gorder/common/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseResponse struct {
}

func (base *BaseResponse) Response(ctx *gin.Context, err error, data interface{}) {
	if err != nil {
		base.Error(ctx, err)
	} else {
		base.Success(ctx, data)
	}
}

func (base *BaseResponse) Error(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, response{
		Errno:   2,
		Message: err.Error(),
		Data:    nil,
		TraceID: tracing.TraceID(ctx.Request.Context()),
	})
}

func (base *BaseResponse) Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, response{
		Errno:   0,
		Message: "success",
		Data:    data,
		TraceID: tracing.TraceID(ctx.Request.Context()),
	})
}

type response struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}
