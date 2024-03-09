package middlewares

import (
	lib "go-gateway/lib/log"

	"github.com/gin-gonic/gin"
)

// traceId 跟踪请求链路 request -> biz -> orm
func GinContextWithTrace () gin.HandlerFunc {
	return func (ctx *gin.Context) {
		lib.SetGinTraceContext(ctx, lib.NewTrace())	
		ctx.Next()
	}
} 