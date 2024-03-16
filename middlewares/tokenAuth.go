package middlewares

import (
	"go-gateway/lib/jwt"
	"go-gateway/public"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTTokenAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := strings.Replace(ctx.GetHeader("Authorization"), "Bearer ", "", 1)
		if token != "" {
			claims, err := jwt.NewJWT().ParseJWT(token)
			if err != nil || claims.UserId == 0{
				public.ResponseError(ctx, 2003, err)
				ctx.Abort()
				return
			}
			ctx.Set("UserID", claims.UserId) 
			ctx.Next()
		}
	}	
}