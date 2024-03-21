package adminController

import (
	adminDto "go-gateway/dto/admin"
	commonDto "go-gateway/dto/common"
	_jwt "go-gateway/lib/jwt"
	libMysql "go-gateway/lib/mysql"
	"go-gateway/model"
	pubic "go-gateway/public"

	"github.com/gin-gonic/gin"
)


type AdminLoginController struct {}

func Register(group gin.IRoutes) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.Login)
	group.GET("/logout", adminLogin.LoginOut)
}


// GO-Gateway godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags Admin
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body adminDto.AdminLoginInput true "body"
// @Success 200 {object} public.Response{data=commonDto.TokensOutput} "success"
// @Router /admin_login/login [post]	
func (adminLogin *AdminLoginController) Login (c *gin.Context) {
	params := adminDto.AdminLoginInput{}

	if err := params.BindValidParam(c); err != nil {
		pubic.ResponseError(c, 2000, err)
		return
	}

	tx, err := libMysql.GetGormPool("default")
	if err != nil {
		pubic.ResponseError(c, 2001, err)
		return
	}

	admin := &model.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		pubic.ResponseError(c, 2002, err)
		return
	}

	token, err := _jwt.NewJWT().GenerateTokenWithUserID(admin.ID)
	if err != nil {
		pubic.ResponseError(c, 2004, err)
		return
	}

	output := &commonDto.TokensOutput{
		ExpiresIn: 3600 * 7,
		TokenType:"Bearer",
		AccessToken: token,
		Scope: "all",
	}

	pubic.ResponseSuccess(c, output)
}

func (adminLogin *AdminLoginController) LoginOut (c *gin.Context) {

}