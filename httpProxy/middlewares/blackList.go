package httpMiddlewares

import (
	"fmt"
	lib "go-gateway/lib/func"
	"go-gateway/model"
	"go-gateway/public"
	"strings"

	"github.com/gin-gonic/gin"
)

func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}

		whiteList := serviceDetail.AccessControl.WhiteList
		blackList := serviceDetail.AccessControl.BlackList

		if whiteList == "" && lib.IsInArrayString(c.ClientIP(), strings.Split(blackList, ",")) {
			public.ResponseError(c, 3001, fmt.Errorf("%s in black ip list", c.ClientIP()))
			c.Abort()
			return
		}
		c.Next()
	}

}