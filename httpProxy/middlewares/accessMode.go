package httpMiddlewares

import (
	"go-gateway/handler"
	"go-gateway/public"

	"github.com/gin-gonic/gin"
)


func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, err := handler.ServiceManagerHandler.HTTPAccessMode(ctx)

		if err != nil {
			public.ResponseError(ctx, 2001, err)
			ctx.Abort()
			return
		}

		ctx.Set("service", service)
		ctx.Next()
	}
}