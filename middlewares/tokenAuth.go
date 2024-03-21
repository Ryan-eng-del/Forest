package middlewares

import (
	"errors"
	"go-gateway/lib/jwt"
	mysqlLib "go-gateway/lib/mysql"
	"go-gateway/model"
	"go-gateway/public"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JWTTokenAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := strings.Replace(ctx.GetHeader("Authorization"), "Bearer ", "", 1)
		if token == "" {
			public.ResponseError(ctx, 2003, errors.New("missing authorization token from server"))
			ctx.Abort()
			return
		}

		claims, err := jwt.NewJWT().ParseJWT(token)
		if err != nil || claims.UserId == 0{
			public.ResponseError(ctx, 2003, err)
			ctx.Abort()
			return
		}
		tx, err := mysqlLib.GetGormPool("default")
		if err != nil {
			public.ResponseError(ctx, 2004, err)
			ctx.Abort()
			return
		}

		admin, err := GetAdmin(ctx, claims.UserId, tx)

		if err != nil || claims.UserId == 0{
			public.ResponseError(ctx, 2004, errors.New("not a valid user"))
			ctx.Abort()
			return
		}
	
		ctx.Set("admin", admin)
		ctx.Next()
		} 
}

func GetAdmin(c *gin.Context, userId uint, tx *gorm.DB) (*model.Admin, error){
	admin := &model.Admin{}
	admin, err := admin.FindById(c, tx, userId)
	return admin, err
}