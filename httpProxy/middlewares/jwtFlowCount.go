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
func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*model.App)
		countBuffer := bytes.NewBufferString(lib.FlowAppPrefix)
		countBuffer.WriteString(appInfo.AppID)

		appCounter, err := handler.ServerCountHandler.GetCounter(countBuffer.String())

		if err != nil {
			public.ResponseError(c, 2002, err)
			c.Abort()
			return
		}

		appCounter.Increase()

		if appInfo.Qpd > 0 && appCounter.TotalCount > int64(appInfo.Qpd) {
			public.ResponseError(c, 2003, fmt.Errorf("租户日请求量限流 limit:%v current:%v", appInfo.Qpd, appCounter.TotalCount))
			c.Abort()
			return
		}
		c.Next()
	}
}