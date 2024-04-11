package middlewares

import (
	"errors"
	lib "go-gateway/lib/const"
	"go-gateway/public"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if adminInfo, ok := session.Get(lib.AdminSessionInfoKey).(string); !ok || adminInfo == "" {
			public.ResponseError(c, 500, errors.New("user not login"))
			c.Abort()
			return
		}
		c.Next()
	}
}
