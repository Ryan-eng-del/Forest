package httpMiddlewares

import (
	"errors"
	"go-gateway/model"
	"go-gateway/public"
	"strings"

	"github.com/gin-gonic/gin"
)

func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			public.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}

		serviceDetail := serverInterface.(*model.ServiceDetail)
		for _, item := range strings.Split(serviceDetail.HTTPRule.HeaderTransfor, ",") {
			items:=strings.Split(item," ")
			if len(items)!=3{
				continue
			}

			if items[0]=="add" || items[0]=="edit"{
				c.Request.Header.Set(items[1],items[2])
			}
			if items[0]=="del"{
				c.Request.Header.Del(items[1])
			}
		}
		c.Next()
	}
}