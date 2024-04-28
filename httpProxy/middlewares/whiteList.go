package httpMiddlewares

import (
	"fmt"
	lib "go-gateway/lib/func"
	"go-gateway/model"
	"go-gateway/public"
	"strings"

	"github.com/gin-gonic/gin"
)

func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)

		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}

		ipList := []string{}
		
		if serviceDetail.AccessControl.WhiteList!=""{
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		if len(ipList) > 0 {
			if !lib.InIPSliceStr(c.ClientIP(), ipList) {
				public.ResponseError(c, 3001, fmt.Errorf("%s not in white ip list", c.ClientIP()))
				c.Abort()
				return
			}
		}

		c.Set("ip_white_list", ipList)

		c.Next()
	}

}