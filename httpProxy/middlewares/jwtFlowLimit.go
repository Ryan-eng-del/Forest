package httpMiddlewares

import (
	"bytes"
	"fmt"
	"go-gateway/handler"
	lib "go-gateway/lib/const"
	"go-gateway/model"
	"go-gateway/public"

	"github.com/gin-gonic/gin"
)

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}

		appInfo := appInterface.(*model.App)

		if appInfo.Qps > 0 {
			cLimiterBuffer := bytes.NewBufferString(lib.FlowAppPrefix)
			cLimiterBuffer.WriteString(appInfo.AppID)
			cLimiterBuffer.WriteString("_")
			//cLimiterBuffer.WriteString(c.ClientIP())
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(cLimiterBuffer.String(), float64(appInfo.Qps), 0, true)
			if err != nil {
				public.ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				public.ResponseError(c, 5002, fmt.Errorf("%v flow limit %v", c.ClientIP(), appInfo.Qps), )
				c.Abort()
				return
			}
		
			c.Next()
		}

	}
}