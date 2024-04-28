package tcpMiddlewares

import (
	"go-gateway/handler"
	libConst "go-gateway/lib/const"
)

func TCPFlowCountMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serviceDetail,err := c.GetServiceDetail()
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}

		totalCounter, err := handler.ServerCountHandler.GetCounter(libConst.FlowTotal)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}

		totalCounter.Increase()
		serviceCounter, err := handler.ServerCountHandler.GetCounter(libConst.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		serviceCounter.Increase()
		c.Next()
	}

}