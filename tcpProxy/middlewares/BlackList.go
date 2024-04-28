package tcpMiddlewares

import (
	"fmt"
	lib "go-gateway/lib/func"
	"strings"
)

func TCPBlackListMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serviceDetail, err := c.GetServiceDetail()

		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		
		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}

		ipList := []string{}
		blackList := []string{}
		if serviceDetail.AccessControl.WhiteList!=""{
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		if serviceDetail.AccessControl.BlackList != "" {
			blackList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}

		if len(blackList) > 0 && len(ipList) == 0{
			if !lib.InIPSliceStr(clientIP, ipList) {
				c.conn.Write([]byte(fmt.Sprintf("%s not in white ip list", clientIP)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}