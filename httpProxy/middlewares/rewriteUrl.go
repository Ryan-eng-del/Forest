package httpMiddlewares

import (
	"go-gateway/model"
	"go-gateway/public"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)


func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}

		for _, item := range strings.Split(serviceDetail.HTTPRule.UrlRewrite, ",") {
			items := strings.Split(item, " ")
			if len(items)!=2 {
				continue
			}

			regexp, err := regexp.Compile(items[0])

			if err != nil {
				continue
			}

			requestPath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
			c.Request.URL.Path = string(requestPath)
		}
		c.Next()
	}

}