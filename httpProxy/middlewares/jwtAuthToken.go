package httpMiddlewares

import (
	"errors"
	"go-gateway/handler"
	libJwt "go-gateway/lib/jwt"
	"go-gateway/model"
	"go-gateway/public"
	"strings"

	"github.com/gin-gonic/gin"
)


func HttpAuthMiddlewareMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}

		appMatched := false
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer", "")

		if token != "" {
			jwtInstance := libJwt.NewJWT()
			appClaims, err := jwtInstance.ParseAppJWT(token)
			if err != nil {
				public.ResponseError(c, 2002, err)
				c.Abort()
				return
			}
			appList :=handler.AppManagerHandler.GetAppList()

			for _, appInfo := range appList {
				if appInfo.AppID == appClaims.AppId   {
					c.Set("app", appInfo)
					appMatched = true
					break
				}
			}

		}

		if serviceDetail.AccessControl.OpenAuth == 1 &&  !appMatched {
			public.ResponseError(c, 2003, errors.New("not match valid app"))
			c.Abort()
			return
		}

		c.Next()
	}
}