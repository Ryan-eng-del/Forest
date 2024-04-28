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

func HttpFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}
		
		serviceFlowNum := serviceDetail.AccessControl.ServiceFlowLimit
		serviceFlowType := serviceDetail.AccessControl.ServiceFlowType

		if serviceFlowNum > 0 {			
			limiterBuffer := bytes.NewBufferString(lib.FlowServicePrefix)
			limiterBuffer.WriteString(serviceDetail.Info.ServiceName)

			serviceLimiter, err := handler.FlowLimiterHandler.GetLimiter(limiterBuffer.String(), float64(serviceFlowNum), int(serviceFlowType), true)
			
			if err != nil {
				public.ResponseError(c, 5001, err)
				c.Abort()
				return
			}

			if !serviceLimiter.Allow() {
				public.ResponseError(c, 5002, fmt.Errorf("service flow limit %v", serviceFlowNum))
				c.Abort()
				return
			}
		}



		clientFlowNum := serviceDetail.AccessControl.ClientIPFlowLimit
		clientFlowType := serviceDetail.AccessControl.ClientFlowType


		if clientFlowNum > 0 {
			cLimiterBuffer := bytes.NewBufferString(lib.FlowServicePrefix)
			cLimiterBuffer.WriteString(serviceDetail.Info.ServiceName)
			cLimiterBuffer.WriteString("_")
			cLimiterBuffer.WriteString(c.ClientIP())
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(cLimiterBuffer.String(), float64(clientFlowNum), int(clientFlowType), true)
			if err != nil {
				public.ResponseError(c, 5003, err)
				c.Abort()
				return
			}
			if clientLimiter == nil {
				public.ResponseError(c, 5002, fmt.Errorf("clientLimiter is nil"))
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				public.ResponseError(c, 5002, fmt.Errorf("%v flow limit %v", c.ClientIP(), clientFlowNum))
				c.Abort()
				return
			}
		}
	}
}