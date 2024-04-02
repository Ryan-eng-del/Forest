package httpMiddlewares

import (
	"bytes"
	"go-gateway/handler"
	lib "go-gateway/lib/const"
	"go-gateway/model"
	"go-gateway/public"

	"github.com/gin-gonic/gin"
)

func HTTPFlowCountMiddleware () gin.HandlerFunc {
	return func (c *gin.Context) {
		serviceDetail,err:= model.GetServiceDetailFromGinContext(c)

		if err != nil{
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}

		totalCounter, err :=  handler.ServerCountHandler.GetCounter(lib.FlowTotal)
		if err != nil {
			public.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		totalCounter.Increase()

		sCounterBuffer := bytes.NewBufferString(lib.FlowServicePrefix)
		sCounterBuffer.WriteString(serviceDetail.Info.ServiceName)

		serviceCounter, err := handler.ServerCountHandler.GetCounter(sCounterBuffer.String())
		if err != nil {
			public.ResponseError(c, 4001, err)
			c.Abort()
			return
		}

		serviceCounter.Increase()
		c.Next()
	}
}