package tcpMiddlewares

import (
	"bytes"
	"fmt"
	"go-gateway/handler"
	libConst "go-gateway/lib/const"
	"go-gateway/model"
	"strings"
)

func TCPFlowLimitMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}

		serviceDetail := serverInterface.(*model.ServiceDetail)
		serviceFlowNum := serviceDetail.AccessControl.ServiceFlowLimit
		serviceFlowType := serviceDetail.AccessControl.ServiceFlowType

		if serviceFlowNum > 0 {		
			serviceLimiter, err := handler.FlowLimiterHandler.GetLimiter(
				libConst.FlowServicePrefix+serviceDetail.Info.ServiceName, float64(serviceFlowNum), serviceFlowType, true)
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("service flow limit %v", serviceFlowNum)))
				c.Abort()
				return
			}
    }

		clientFlowNum := serviceDetail.AccessControl.ClientIPFlowLimit
		clientFlowType := serviceDetail.AccessControl.ClientFlowType

		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}

		if clientFlowNum > 0 {
			cLimiterBuffer := bytes.NewBufferString(libConst.FlowServicePrefix)
			cLimiterBuffer.WriteString(serviceDetail.Info.ServiceName)
			cLimiterBuffer.WriteString("_")
			cLimiterBuffer.WriteString(clientIP)
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(cLimiterBuffer.String(), float64(clientFlowNum), int(clientFlowType), true)

			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if clientLimiter == nil {
				c.conn.Write([]byte(fmt.Sprintf("%v flow limit %v", clientIP, clientFlowNum)))
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("%v flow limit %v", clientIP, clientFlowNum)))
				c.Abort()
				return
			}
		}
		c.Next()
  }
}