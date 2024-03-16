package adminDto

import (
	lib "go-gateway/lib/func"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminInfoOutput struct {
	ID           uint     `json:"id"`
	Name         string    `json:"name"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}


type ChangePwdInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"` //密码
}

func (input *ChangePwdInput) BindValidParam (c *gin.Context) error {
	return lib.ValidateParams(c, input)
}
 