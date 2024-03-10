package middlewares

import (
	"errors"
	"fmt"
	libConf "go-gateway/lib/conf"
	libLog "go-gateway/lib/log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)


func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func ()  {
			if err := recover(); err != nil {
				libLog.ComLogNotice(ctx, "_com_panic", map[string]interface{}{
					"error": fmt.Sprint(err),
					"stack": string(debug.Stack()),
				})

				if libConf.BaseConfInstance.DebugMode != "debug" {
					libLog.ResponseError(ctx, 500, errors.New("内部错误"))
					return
				} else {
					libLog.ResponseError(ctx, 500, errors.New(fmt.Sprint(err)))
					return
				}
			}
		}()

		ctx.Next()
	}
}