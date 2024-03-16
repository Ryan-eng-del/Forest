package adminController

import (
	adminDto "go-gateway/dto/admin"
	commonDto "go-gateway/dto/common"
	_jwt "go-gateway/lib/jwt"
	libMysql "go-gateway/lib/mysql"
	"go-gateway/middlewares"
	"go-gateway/model"
	"log"

	"github.com/gin-gonic/gin"
)

// gin-swagger middleware
// swagger embed files

type AdminLoginController struct {}

func Register(group *gin.RouterGroup) {
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
// @Success 200 {object} middlewares.Response{data=commonDto.TokensOutput} "success"
// @Router /admin_login/login [post]	
func (adminLogin *AdminLoginController) Login (c *gin.Context) {
	params := adminDto.AdminLoginInput{}

	if err := params.BindValidParam(c); err != nil {
		middlewares.ResponseError(c, 2000, err)
		return
	}

	tx, err := libMysql.GetGormPool("default")
	if err != nil {
		middlewares.ResponseError(c, 2001, err)
		return
	}
	

	admin := &model.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middlewares.ResponseError(c, 2002, err)
		return
	}

	token, err := _jwt.NewJWT().GenerateTokenWithUserID(admin.ID)
	log.Println(token, "token")
	if err != nil {
		middlewares.ResponseError(c, 2004, err)
		return
	}

	output := &commonDto.TokensOutput{
		ExpiresIn: 3600 * 7,
		TokenType:"Bearer",
		AccessToken: token,
		Scope: "all",
	}

	middlewares.ResponseSuccess(c, output)
}

func (adminLogin *AdminLoginController) LoginOut (c *gin.Context) {

}