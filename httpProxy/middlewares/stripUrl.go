package httpMiddlewares

import (
	lib "go-gateway/lib/const"
	"go-gateway/model"
	"go-gateway/public"
	"strings"

	"github.com/gin-gonic/gin"
)



func HTTPStripUriMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}

		if serviceDetail.HTTPRule.RuleType == lib.HTTPRuleTypePrefixURL && serviceDetail.HTTPRule.NeedStripUri == 1 {
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.HTTPRule.Rule,"",1)
		}

		c.Next()
	}
}