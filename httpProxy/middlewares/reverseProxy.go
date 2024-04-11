package httpMiddlewares

import (
	"go-gateway/handler"
	httpproxy "go-gateway/httpProxy"
	"go-gateway/model"
	"go-gateway/public"

	"github.com/gin-gonic/gin"
)

func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}

		lb, err :=  handler.LoadBalancerHandler.GetLoadBalancer(serviceDetail)

		if err != nil {
			public.ResponseError(c, 2002, err)
			c.Abort()
			return
		}

		trans, err := handler.TransportorHandler.GetTrans(serviceDetail)

		if err != nil {
			public.ResponseError(c, 2003, err)
			c.Abort()
			return
		}
		proxy := httpproxy.NewLoadBalanceReverseProxy(c, lb, trans)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}	
