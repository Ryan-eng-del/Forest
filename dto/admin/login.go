package adminDto

import (
	lib "go-gateway/lib/func"

	"github.com/gin-gonic/gin"
)


type AdminLoginInput struct {
	Username string `json:"username" form:"username" comment:"姓名"  example:"admin" validate:"required,valid_username" en_comment:"username"`
	Password string `json:"password" form:"password" comment:"密码"  example:"123456" validate:"required" en_comment:"password"`
}

func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return lib.ValidateParams(c, param)
}